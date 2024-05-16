package models

import "time"

type VersionMetadata struct {
	//Value       []byte    `mapstructure:"version" dynamodbav:"value"`
	Version     int       `mapstructure:"version" dynamodbav:"version"`
	CreatedTime time.Time `mapstructure:"created_time" dynamodbav:"created_time,unixtime"`
}
