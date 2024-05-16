package aws

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/cenkalti/backoff/v4"
	"log"
	"math"
	"nbox/internal/application"
	"nbox/internal/domain"
	"nbox/internal/domain/models"
	pkgPath "path"
	"strings"
	"time"
)

const (
	DynamoDBEmptyPath         = " "
	DynamoDBLockPrefix        = "_"
	DefaultParallelOperations = 128
)

type PermitPool struct {
	sem chan int
}

type Record struct {
	Path    string                 `dynamodbav:"Path"`
	Key     string                 `dynamodbav:"Key"`
	Value   []byte                 `dynamodbav:"Value"`
	Version models.VersionMetadata `dynamodbav:"Version"`
}

func NewPermitPool(permits int) *PermitPool {
	if permits < 1 {
		permits = DefaultParallelOperations
	}
	return &PermitPool{
		sem: make(chan int, permits),
	}
}

// Acquire returns when a permit has been acquired
func (c *PermitPool) Acquire() {
	c.sem <- 1
}

// Release returns a permit to the pool
func (c *PermitPool) Release() {
	<-c.sem
}

// CurrentPermits Get number of requests in the permit pool
func (c *PermitPool) CurrentPermits() int {
	return len(c.sem)
}

type dynamodbBackend struct {
	client     *dynamodb.Client
	config     *application.Config
	permitPool *PermitPool
}

func NewDynamodbBackend(dynamodb *dynamodb.Client, config *application.Config) domain.EntryAdapter {
	return &dynamodbBackend{
		client:     dynamodb,
		config:     config,
		permitPool: NewPermitPool(0),
	}
}

// Upsert is used to insert or update an entry
func (d *dynamodbBackend) Upsert(ctx context.Context, entries []models.Entry) error {
	var writeReqs []types.WriteRequest
	var item map[string]types.AttributeValue
	var err error

	records := map[string]Record{}
	for _, entry := range entries {
		path := recordPathForVaultKey(entry.Key)
		key := recordKeyForVaultKey(entry.Key)
		records[fmt.Sprintf("%s%s", path, key)] = Record{Path: path, Key: key, Value: entry.Value}

		for _, prefix := range Prefixes(entry.Key) {
			path = recordPathForVaultKey(prefix)
			key = fmt.Sprintf("%s/", recordKeyForVaultKey(prefix))
			records[fmt.Sprintf("%s%s", path, key)] = Record{
				Path: recordPathForVaultKey(prefix),
				Key:  fmt.Sprintf("%s/", recordKeyForVaultKey(prefix)),
			}
		}
	}

	for _, r := range records {
		item, err = attributevalue.MarshalMap(r)
		if err != nil {
			log.Printf("Err could not convert prefix record to DynamoDB item: %v", err)
			continue
		}
		writeReqs = append(
			writeReqs, types.WriteRequest{PutRequest: &types.PutRequest{Item: item}},
		)
	}

	//for _, entry := range entries {
	//	//item, err = d.marshalEntry(recordPathForVaultKey(entry.Key), recordKeyForVaultKey(entry.Key), entry.Value)
	//	item, err = attributevalue.MarshalMap(Record{
	//		Path:  recordPathForVaultKey(entry.Key),
	//		Key:   recordKeyForVaultKey(entry.Key),
	//		Value: entry.Value,
	//	})
	//	if err != nil {
	//		log.Printf("Err could not convert prefix record to DynamoDB item: %v", err)
	//		continue
	//	}
	//	writeReqs = append(
	//		writeReqs, types.WriteRequest{PutRequest: &types.PutRequest{Item: item}},
	//	)
	//
	//	for _, prefix := range Prefixes(entry.Key) {
	//		item, err = attributevalue.MarshalMap(Record{
	//			Path: recordPathForVaultKey(prefix),
	//			Key:  fmt.Sprintf("%s/", recordKeyForVaultKey(prefix)),
	//		})
	//		if err != nil {
	//			log.Printf("Err could not convert prefix record to DynamoDB item: %v", err)
	//			continue
	//		}
	//		writeReqs = append(writeReqs, types.WriteRequest{PutRequest: &types.PutRequest{Item: item}})
	//	}
	//}

	//fmt.Printf("writeReqs %+v", writeReqs)
	return d.writeReqsBatch(ctx, writeReqs)
}

func (d *dynamodbBackend) writeReqsBatch(ctx context.Context, requests []types.WriteRequest) error {
	for len(requests) > 0 {
		var err error
		batchSize := int(math.Min(float64(len(requests)), 25))
		batch := map[string][]types.WriteRequest{d.config.TableName: requests[:batchSize]}
		requests = requests[batchSize:]

		d.permitPool.Acquire()
		boff := backoff.NewExponentialBackOff()
		boff.MaxElapsedTime = 600 * time.Second

		for len(batch) > 0 {
			var output *dynamodb.BatchWriteItemOutput
			output, err = d.client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
				RequestItems: batch,
			})
			if err != nil {
				break
			}
			if len(output.UnprocessedItems) == 0 {
				break
			} else {
				duration := boff.NextBackOff()
				if duration != backoff.Stop {
					batch = output.UnprocessedItems
					time.Sleep(duration)
				} else {
					err = errors.New("dynamodb: timeout handling UnproccessedItems")
					break
				}
			}
		}
		d.permitPool.Release()
		if err != nil {
			return err
		}
	}
	return nil
}

// Retrieve Get is used to fetch an entry
func (d *dynamodbBackend) Retrieve(ctx context.Context, key string) (*models.Entry, error) {
	p, _ := attributevalue.Marshal(recordPathForVaultKey(key))
	k, _ := attributevalue.Marshal(recordKeyForVaultKey(key))

	resp, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		Key:            map[string]types.AttributeValue{"Path": p, "Key": k},
		TableName:      aws.String(d.config.TableName),
		ConsistentRead: aws.Bool(true),
	})

	if err != nil {
		return nil, err
	}
	if resp.Item == nil {
		return nil, nil
	}
	record := &Record{}
	err = attributevalue.UnmarshalMap(resp.Item, record)
	if err != nil {
		return nil, err
	}

	return &models.Entry{
		Key:   vaultKey(record),
		Value: record.Value,
	}, nil
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (d *dynamodbBackend) List(ctx context.Context, prefix string) ([]models.Entry, error) {
	prefix = strings.TrimSuffix(prefix, "/")
	entries := make([]models.Entry, 0)
	prefix = escapeEmptyPath(prefix)

	keyEx := expression.Key("Path").Equal(expression.Value(prefix))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()

	if err != nil {
		log.Printf("Err expression Builder %v \n", err)
		return nil, err
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(d.config.TableName),
		ConsistentRead:            aws.Bool(true),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}
	queryPaginator := dynamodb.NewQueryPaginator(d.client, queryInput)
	for queryPaginator.HasMorePages() {
		response, err := queryPaginator.NextPage(ctx)
		if err != nil {
			log.Printf("Err Couldn't query for records released in %v. %v\n", prefix, err)
			return nil, err
		}
		var records []Record
		err = attributevalue.UnmarshalListOfMaps(response.Items, &records)
		if err != nil {
			log.Printf("Err Couldn't unmarshal query response. %v\n", err)
			return nil, err
		}

		for _, record := range records {
			if !strings.HasPrefix(record.Key, DynamoDBLockPrefix) {
				entries = append(entries, models.Entry{
					Key:   record.Key,
					Value: record.Value,
				})
			}
		}
	}

	return entries, nil
}

// recordPathForVaultKey transforms a vault key into
// a value suitable for the `DynamoDBRecord`'s `Path`
// property. This path equals the vault key without
// its last component.
func recordPathForVaultKey(key string) string {
	if strings.Contains(key, "/") {
		return pkgPath.Dir(key)
	}
	return DynamoDBEmptyPath
}

// recordKeyForVaultKey transforms a vault key into
// a value suitable for the `DynamoDBRecord`'s `Key`
// property. This path equals the vault key's
// last component.
func recordKeyForVaultKey(key string) string {
	return pkgPath.Base(key)
}

// vaultKey returns the vault key for a given record
// from the DynamoDB table. This is the combination of
// the records Path and Key.
func vaultKey(record *Record) string {
	path := unescapeEmptyPath(record.Path)
	if path == "" {
		return record.Key
	}
	return pkgPath.Join(record.Path, record.Key)
}

// escapeEmptyPath is used to escape the root key's path
// with a value that can be stored in DynamoDB. DynamoDB
// does not allow values to be empty strings.
func escapeEmptyPath(s string) string {
	if s == "" {
		return DynamoDBEmptyPath
	}
	return s
}

// unescapeEmptyPath is the opposite of `escapeEmptyPath`.
func unescapeEmptyPath(s string) string {
	if s == DynamoDBEmptyPath {
		return ""
	}
	return s
}

// Prefixes is a shared helper function returns all parent 'folders' for a
// given vault key.
// e.g. for 'foo/bar/baz', it returns ['foo', 'foo/bar']
func Prefixes(s string) []string {
	components := strings.Split(s, "/")
	var result []string
	for i := 1; i < len(components); i++ {
		result = append(result, strings.Join(components[:i], "/"))
	}
	return result
}
