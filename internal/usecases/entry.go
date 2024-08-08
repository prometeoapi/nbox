package usecases

import (
	"context"
	"fmt"
	"nbox/internal/application"
	"nbox/internal/domain"
	"nbox/internal/domain/models"
	"strings"
)

type EntryUseCase struct {
	entryAdapter  domain.EntryAdapter
	secretAdapter domain.SecretAdapter
	config        *application.Config
}

func NewEntryUseCase(
	entryAdapter domain.EntryAdapter,
	secretAdapter domain.SecretAdapter,
	config *application.Config,
) *EntryUseCase {
	return &EntryUseCase{entryAdapter: entryAdapter, secretAdapter: secretAdapter, config: config}
}

// Upsert
// ARN arn:aws:ssm:<REGION_NAME>:<ACCOUNT_ID>:parameter/<parameter-name>
func (e *EntryUseCase) Upsert(ctx context.Context, entries []models.Entry) map[string]error {

	result := make(map[string]error)

	secrets := make([]models.Entry, 0)
	for _, entry := range entries {
		if entry.Secure {
			secrets = append(secrets, entry)
		}
		result[entry.Key] = nil
	}

	secureResults := e.secretAdapter.Upsert(ctx, secrets)

	for i, entry := range entries {
		if entry.Secure {
			err := secureResults[entry.Key]
			entries[i].Value = ""

			if err != nil {
				result[entry.Key] = err
				continue
			}

			key := cleanedKey(entry.Key)
			entries[i].Value = e.GetParameterArn(key)
		}
	}

	updated := e.entryAdapter.Upsert(ctx, entries)

	for _, entry := range entries {
		if !entry.Secure {
			err := updated[entry.Key]
			if err != nil {
				result[entry.Key] = err
			}
		}
	}

	return result
}

func (e *EntryUseCase) GetParameterArn(key string) string {
	return fmt.Sprintf(
		"arn:aws:ssm:%s:%s:parameter/%s", e.config.RegionName, e.config.AccountId, cleanedKey(key),
	)
}

func cleanedKey(key string) string {
	return strings.TrimPrefix(key, "/")
}
