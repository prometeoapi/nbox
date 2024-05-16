package aws

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"nbox/internal/application"
	"nbox/internal/domain"
	"nbox/internal/domain/models"
)

type storeAdapter struct {
	s3       *s3.Client
	dynamodb *dynamodb.Client
	config   *application.Config
}

func NewStoreAdapter(s3 *s3.Client, dynamodb *dynamodb.Client, config *application.Config) domain.StoreOperations {
	return &storeAdapter{
		s3:       s3,
		dynamodb: dynamodb,
		config:   config,
	}
}

//func (b *storeAdapter) BucketExists(bucketName string) (bool, error) {
//	return true, nil
//}

func (b *storeAdapter) save(ctx context.Context, box *models.Box) (*models.Box, error) {

	item, err := attributevalue.MarshalMap(box)
	if err != nil {
		return nil, err
	}
	out, err := b.dynamodb.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(b.config.TableName),
	})

	fmt.Printf("%+v\n", out)

	if err != nil {
		return nil, err
	}

	return box, nil
}

func (b *storeAdapter) storeTemplate(ctx context.Context, path string, stage models.Stage) error {

	respHead, err := b.s3.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(b.config.BucketName),
		Key:    aws.String(path),
	})
	fmt.Printf("%+v \n", respHead)

	if err != nil {
		var out bytes.Buffer
		err := json.Indent(&out, []byte(stage.Template.Value), "", "  ")
		if err != nil {
			return err
		}

		respPut, err := b.s3.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(b.config.BucketName),
			Key:    aws.String(path),
			Body:   bytes.NewReader(out.Bytes()),
		})

		fmt.Printf("%+v \n", respPut)

		if err != nil {
			return err
		}

	}
	return nil
}

func (b *storeAdapter) CreateBox(box *models.Box) (*models.Box, error) {
	ctx := context.Background()
	for stageName, stage := range box.Stage {
		path := fmt.Sprintf("%s/%s/%s", box.Service, stageName, stage.Template.Name)
		stage.Template.Name = path
		box.Stage[stageName] = stage
		s3Result := b.storeTemplate(ctx, path, stage)
		//if s3Result == nil {
		//	return b.save(ctx, box)
		//}

		return nil, s3Result
	}
	return nil, nil
}

//func (b storeAdapter) RetrieveBox(boxName string, stage string) models.Box {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (b storeAdapter) GetOrCreateStage(box models.Box, stage string) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (b storeAdapter) UpsertTemplate(box models.Box, template string) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (b storeAdapter) UpsertVariable(box models.Box, value interface{}) {
//	//TODO implement me
//	panic("implement me")
//}
