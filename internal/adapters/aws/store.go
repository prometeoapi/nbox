package aws

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"nbox/internal/application"
	"nbox/internal/domain"
	"nbox/internal/domain/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type storeAdapter struct {
	s3             *s3.Client
	dynamodbClient *dynamodb.Client
	config         *application.Config
}

type BoxRecord struct {
	Service  string          `dynamodbav:"Service"`
	Stage    string          `dynamodbav:"Stage"`
	Template models.Template `dynamodbav:"Template"`
}

func NewStoreAdapter(s3 *s3.Client, config *application.Config, dynamodb *dynamodb.Client) domain.StoreOperations {
	return &storeAdapter{
		s3:             s3,
		dynamodbClient: dynamodb,
		config:         config,
	}
}

func (b *storeAdapter) store(ctx context.Context, path string, stage models.Stage) (*s3.PutObjectOutput, error) {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(stage.Template.Value), "", "  ")
	if err != nil {
		return nil, err
	}

	return b.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(b.config.BucketName),
		Key:    aws.String(path),
		Body:   bytes.NewReader(out.Bytes()),
	})
}

func (b *storeAdapter) BoxExists(ctx context.Context, service string, stage string, template string) (bool, error) {
	path := fmt.Sprintf("%s/%s/%s", service, stage, template)

	_, err := b.s3.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(b.config.BucketName),
		Key:    aws.String(path),
	})

	return err == nil, err
}

func (b *storeAdapter) RetrieveBox(ctx context.Context, service string, stage string, template string) ([]byte, error) {
	path := fmt.Sprintf("%s/%s/%s", service, stage, template)
	object, err := b.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(b.config.BucketName),
		Key:    aws.String(path),
	})

	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(object.Body)

	body, err := io.ReadAll(object.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func (b *storeAdapter) UpsertBox(ctx context.Context, box *models.Box) []string {
	result := make([]string, 0)
	var item map[string]types.AttributeValue

	for stageName, stage := range box.Stage {
		name := stage.Template.Name
		path := fmt.Sprintf("%s/%s/%s", box.Service, stageName, stage.Template.Name)

		stage.Template.Name = path
		box.Stage[stageName] = stage
		_, err := b.store(ctx, path, stage)
		fmt.Printf("ErrStore. %s", err)
		if err == nil {
			item, _ = attributevalue.MarshalMap(BoxRecord{
				Service: box.Service,
				Stage:   stageName,
				Template: models.Template{
					Name:  path,
					Value: name,
				},
			})
			_, err = b.dynamodbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
				TableName: aws.String(b.config.BoxTableName), Item: item,
			})
			fmt.Printf("ErrDbStore. %s", err)
		}

		if err == nil {
			result = append(result, path)
		}
	}
	return result
}

func (b *storeAdapter) List(ctx context.Context) ([]models.Box, error) {
	var err error
	boxes := map[string]models.Box{}
	results := make([]models.Box, 0)

	if err != nil {
		log.Printf("Err expression Builder %v \n", err)
		return nil, err
	}

	scan, err := b.dynamodbClient.Scan(ctx, &dynamodb.ScanInput{
		TableName:              aws.String(b.config.BoxTableName),
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
	})

	if err != nil {
		return nil, err
	}

	for _, i := range scan.Items {
		var record BoxRecord
		err = attributevalue.UnmarshalMap(i, &record)
		if err != nil {
			continue
		}

		_, ok := boxes[record.Service]
		if !ok {
			boxes[record.Service] = models.Box{Service: record.Service, Stage: map[string]models.Stage{}}
		}
		boxes[record.Service].Stage[record.Stage] = models.Stage{Template: record.Template}
	}

	for _, box := range boxes {
		results = append(results, box)
	}

	return results, nil
}
