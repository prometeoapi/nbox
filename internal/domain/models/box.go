package models

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Box struct {
	Service string           `json:"service" dynamodbav:"service"`
	Stage   map[string]Stage `json:"stage" dynamodbav:"stage"`
}

func (b Box) GetKey() map[string]types.AttributeValue {
	service, _ := attributevalue.Marshal(b.Service)
	return map[string]types.AttributeValue{"service": service}
}

type Stage struct {
	Template  Template            `json:"template" dynamodbav:"template"`
	Variables map[string]Variable `json:"variables" dynamodbav:"variables"`
}

type Template struct {
	Name  string `json:"name" dynamodbav:"name"` // s3 path
	Value string `json:"value" dynamodbav:"-"`
	VersionMetadata
}

type Variable struct {
	Value interface{} `json:"value" dynamodbav:"value"`
	VersionMetadata
}
