package aws

import (
	"context"
	"log"
	"nbox/internal/application"
	"nbox/internal/domain"
	"nbox/internal/domain/models"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type Result models.Exchange[*ssm.PutParameterOutput, *models.Entry]

type secureParameterStore struct {
	client *ssm.Client
	config *application.Config
}

func NewSecureParameterStore(client *ssm.Client, config *application.Config) domain.SecretAdapter {
	return &secureParameterStore{client: client, config: config}
}

func (s *secureParameterStore) Upsert(ctx context.Context, entries []models.Entry) map[string]error {
	ch := make(chan Result)
	wg := sync.WaitGroup{}
	wg.Add(len(entries))

	for _, entry := range entries {
		go func(c chan Result, g *sync.WaitGroup, e models.Entry, x context.Context) {
			defer g.Done()
			c <- s.Send(x, e)
		}(ch, &wg, entry, ctx)
	}

	go func(g *sync.WaitGroup, c chan Result) {
		g.Wait()
		defer close(c)
	}(&wg, ch)

	summary := make(map[string]error)

	for result := range ch {
		if result.Err != nil {
			log.Printf("Err upsert secret[%s]. %v. %v \n", result.In.Key, result.Err, result.Out)
		}
		summary[result.In.Key] = result.Err
	}

	return summary
}

func (s *secureParameterStore) Send(ctx context.Context, entry models.Entry) Result {
	in := prepareSecret(entry, s.config.ParameterStoreDefaultTier, s.config.ParameterStoreKeyId)
	out, err := s.client.PutParameter(ctx, in)
	result := Result{Out: out, In: &entry, Err: err}

	if err != nil {
		return result
	}

	if out.Version == 1 {
		s.AddTags(ctx, in.Name)
	}

	return result
}

func (s *secureParameterStore) AddTags(ctx context.Context, key *string) {
	_, err := s.client.AddTagsToResource(ctx, &ssm.AddTagsToResourceInput{
		ResourceId:   key,
		ResourceType: types.ResourceTypeForTaggingParameter,
		Tags:         []types.Tag{{Key: aws.String("project"), Value: aws.String("nbox")}},
	})
	if err != nil {
		log.Printf("Err add tags [parameter:%s]. %v \n", *key, err)
	}
}

func prepareSecret(entry models.Entry, parameterStoreDefaultTier string, parameterStoreKeyId string) *ssm.PutParameterInput {
	key := entry.Key

	if !strings.HasPrefix(key, "/") {
		key = "/" + key
	}

	parameterInput := &ssm.PutParameterInput{
		Name:      aws.String(key),
		Value:     aws.String(entry.Value),
		Type:      types.ParameterTypeSecureString,
		Tier:      types.ParameterTierStandard,
		Overwrite: aws.Bool(true),
	}

	if parameterStoreDefaultTier != "" {
		parameterInput.Tier = types.ParameterTier(parameterStoreDefaultTier)
	}

	if parameterStoreKeyId != "" {
		parameterInput.KeyId = aws.String(parameterStoreKeyId)
	}

	return parameterInput
}
