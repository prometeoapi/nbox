package application

import "os"

type Config struct {
	BucketName     string `pkl:"bucketName"`
	EntryTableName string `pkl:"entryTableName"`
	BoxTableName   string `pkl:"boxTableName"`
}

func NewConfig() *Config {
	//evaluator, err := pkl.NewEvaluator(context.Background(), pkl.PreconfiguredOptions)
	//if err != nil {
	//	panic(err)
	//}
	//defer func(e pkl.Evaluator) {
	//	_ = e.Close()
	//}(evaluator)
	//
	//var cfg Config
	//if err = evaluator.EvaluateModule(context.Background(), pkl.FileSource("config.pkl"), &cfg); err != nil {
	//	panic(err)
	//}
	//return &cfg
	return &Config{
		BucketName:     os.Getenv("NBOX_BUCKET_NAME"),
		EntryTableName: os.Getenv("NBOX_ENTRIES_TABLE_NAME"),
		BoxTableName:   os.Getenv("NBOX_BOX_TABLE_NAME"),
	}
}
