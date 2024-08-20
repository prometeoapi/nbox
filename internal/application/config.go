package application

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	BucketName                string   `pkl:"bucketName"`
	EntryTableName            string   `pkl:"entryTableName"`
	TrackingEntryTableName    string   `pkl:"trackingEntryTableName"`
	BoxTableName              string   `pkl:"boxTableName"`
	RegionName                string   `pkl:"regionName"`
	AccountId                 string   `pkl:"accountId"`
	ParameterStoreDefaultTier string   `pkl:"parameterStoreDefaultTier"`
	ParameterStoreKeyId       string   `pkl:"parameterStoreKeyId"`
	ParameterShortArn         bool     `pkl:"parameterShortArn"`
	DefaultPrefix             string   `pkl:"defaultPrefix"`
	AllowedPrefixes           []string `pkl:"allowedPrefixes"`
}

//func NewConfigFromPkl()  {
//	evaluator, err := pkl.NewEvaluator(context.Background(), pkl.PreconfiguredOptions)
//	if err != nil {
//		panic(err)
//	}
//	defer func(e pkl.Evaluator) {
//		_ = e.Close()
//	}(evaluator)
//
//	var cfg Config
//	if err = evaluator.EvaluateModule(context.Background(), pkl.FileSource("config.pkl"), &cfg); err != nil {
//		panic(err)
//	}
//	return &cfg
//}

//func NewConfigFromYaml()  {
//	config := &Config{}
//
//	viper.SetConfigName(ConfigName)
//	viper.SetConfigType("yaml")
//	viper.AddConfigPath(ConfigPath)
//
//	err := viper.ReadInConfig()
//	if err != nil { // Handle errs reading the config file
//		panic(fmt.Errorf("fatal error config file: %w", err))
//	}
//
//	err = viper.Unmarshal(&config)
//	if err != nil {
//		log.Error("Environment can't be loaded", zap.Error(err))
//	}
//
//	return config
//}

func NewConfigFromEnv() *Config {
	var prefixes []string

	prefixes = append(
		prefixes,
		strings.Split(env("NBOX_ALLOWED_PREFIXES", "development,qa,beta,staging,sandbox,production"), ",")...,
	)

	return &Config{
		BucketName:                env("NBOX_BUCKET_NAME", "nbox-store"),
		EntryTableName:            env("NBOX_ENTRIES_TABLE_NAME", "nbox-entry-table"),
		TrackingEntryTableName:    env("NBOX_TRACKING_ENTRIES_TABLE_NAME", "nbox-tracking-entry-table"),
		BoxTableName:              env("NBOX_BOX_TABLE_NAME", "nbox-box-table"),
		AccountId:                 env("ACCOUNT_ID", ""),
		RegionName:                env("AWS_REGION", "us-east-1"),
		ParameterStoreDefaultTier: env("NBOX_PARAMETER_STORE_DEFAULT_TIER", "Standard"), // Standard | Advanced
		ParameterStoreKeyId:       env("NBOX_PARAMETER_STORE_KEY_ID", ""),               // KMS KEY ID
		ParameterShortArn:         envBool("NBOX_PARAMETER_STORE_SHORT_ARN"),
		DefaultPrefix:             env("NBOX_DEFAULT_PREFIX", "global"),
		AllowedPrefixes:           prefixes,
	}
}

func env(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists || strings.TrimSpace(value) == "" {
		return defaultValue
	}
	return value
}

//func envInt(key string, defaultValue int) int {
//	value := env(key, fmt.Sprint(defaultValue))
//	valueInt, err := strconv.Atoi(value)
//	if err != nil {
//		return defaultValue
//	}
//	return valueInt
//}

func envBool(key string) bool {
	s := env(key, "false")
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return v
}
