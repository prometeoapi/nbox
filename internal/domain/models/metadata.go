package models

import "time"

type Metadata struct {
	Hash      string    `dynamodbav:"Hash,omitempty"`
	UpdatedAt time.Time `dynamodbav:"UpdatedAt,unixtime"`
	UpdatedBy string    `dynamodbav:"UpdatedBy,omitempty"`
}
