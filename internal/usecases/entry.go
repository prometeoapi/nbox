package usecases

import (
	"context"
	"nbox/internal/domain/models"
)

type EntryUseCase struct {
}

func NewEntryUseCase() *EntryUseCase {
	return &EntryUseCase{}
}

func (e *EntryUseCase) Upsert(entry *models.Entry) error {
	return nil
}

func (e *EntryUseCase) Retrieve(ctx context.Context, key string) (*models.Entry, error) {
	return nil, nil
}

func (e *EntryUseCase) List(ctx context.Context, prefix string) ([]string, error) {
	return nil, nil
}
