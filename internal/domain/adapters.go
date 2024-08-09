package domain

import (
	"context"
	"nbox/internal/domain/models"
)

// TemplateAdapter store templates
type TemplateAdapter interface {
	UpsertBox(ctx context.Context, box *models.Box) []string
	BoxExists(ctx context.Context, service string, stage string, template string) (bool, error)
	RetrieveBox(ctx context.Context, service string, stage string, template string) ([]byte, error)
	List(ctx context.Context) ([]models.Box, error)
}

// EntryAdapter vars backend operations
type EntryAdapter interface {
	Upsert(ctx context.Context, entries []models.Entry) map[string]error
	Retrieve(ctx context.Context, key string) (*models.Entry, error)
	List(ctx context.Context, prefix string) ([]models.Entry, error)
	Delete(ctx context.Context, key string) error
	Tracking(ctx context.Context, key string) ([]models.Tracking, error)
}

// SecretAdapter vars encrypt
type SecretAdapter interface {
	Upsert(ctx context.Context, entries []models.Entry) map[string]error
}
