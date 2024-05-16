package domain

import (
	"context"
	"nbox/internal/domain/models"
)

type StoreOperations interface {
	CreateBox(*models.Box) (*models.Box, error)
	//RetrieveBox(boxName string, stage string) models.Box
	//GetOrCreateStage(box models.Box, stage string)
	//UpsertTemplate(box models.Box, template string)
	//UpsertVariable(box models.Box, value interface{})
}

type EntryAdapter interface {
	Upsert(ctx context.Context, entries []models.Entry) error
	Retrieve(ctx context.Context, key string) (*models.Entry, error)
	List(ctx context.Context, prefix string) ([]models.Entry, error)
}
