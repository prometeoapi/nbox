package aws

import (
	"context"
	"nbox/internal/domain"
	"nbox/internal/domain/models"
)

type parameterStoreBackend struct {
}

func (p parameterStoreBackend) Delete(ctx context.Context, key string) error {
	//TODO implement me
	panic("implement me")
}

func NewParameterStoreBackend() domain.EntryAdapter {
	return &parameterStoreBackend{}
}

func (p parameterStoreBackend) Upsert(ctx context.Context, entries []models.Entry) error {
	//TODO implement me
	panic("implement me")
}

func (p parameterStoreBackend) Retrieve(ctx context.Context, key string) (*models.Entry, error) {
	//TODO implement me
	panic("implement me")
}

func (p parameterStoreBackend) List(ctx context.Context, prefix string) ([]models.Entry, error) {
	//TODO implement me
	panic("implement me")
}
