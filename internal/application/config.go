package application

import (
	"context"
	"github.com/apple/pkl-go/pkl"
)

type Config struct {
	BucketName string `pkl:"bucketName"`
	TableName  string `pkl:"tableName"`
}

func NewConfig() *Config {
	evaluator, err := pkl.NewEvaluator(context.Background(), pkl.PreconfiguredOptions)
	if err != nil {
		panic(err)
	}
	defer func(e pkl.Evaluator) {
		_ = e.Close()
	}(evaluator)

	var cfg Config
	if err = evaluator.EvaluateModule(context.Background(), pkl.FileSource("config.pkl"), &cfg); err != nil {
		panic(err)
	}
	return &cfg
}
