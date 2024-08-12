package application

import "os"

type Config struct {
	BucketName                string `pkl:"bucketName"`
	EntryTableName            string `pkl:"entryTableName"`
	TrackingEntryTableName    string `pkl:"trackingEntryTableName"`
	BoxTableName              string `pkl:"boxTableName"`
	RegionName                string `pkl:"regionName"`
	AccountId                 string `pkl:"accountId"`
	ParameterStoreDefaultTier string `pkl:"parameterStoreDefaultTier"`
	ParameterStoreKeyId       string `pkl:"parameterStoreKeyId"`
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
		BucketName:                os.Getenv("NBOX_BUCKET_NAME"),
		EntryTableName:            os.Getenv("NBOX_ENTRIES_TABLE_NAME"),
		TrackingEntryTableName:    os.Getenv("NBOX_TRACKING_ENTRIES_TABLE_NAME"),
		BoxTableName:              os.Getenv("NBOX_BOX_TABLE_NAME"),
		AccountId:                 os.Getenv("ACCOUNT_ID"),
		RegionName:                os.Getenv("AWS_REGION"),
		ParameterStoreDefaultTier: os.Getenv("NBOX_PARAMETER_STORE_DEFAULT_TIER"), // Standard | Advanced
		ParameterStoreKeyId:       os.Getenv("NBOX_PARAMETER_STORE_KEY_ID"),       // KMS KEY ID
	}
}
