package models

import (
	"time"
)

type Metadata struct {
	Hash      string    `json:"hash" dynamodbav:"Hash,omitempty"`
	Secure    bool      `json:"secure" dynamodbav:"Secure"`
	Action    string    `json:"action" dynamodbav:"Action,omitempty"`
	UpdatedAt time.Time `json:"updatedAt" dynamodbav:"UpdatedAt,unixtime"`
	UpdatedBy string    `json:"updatedBy" dynamodbav:"UpdatedBy,omitempty"`
}
