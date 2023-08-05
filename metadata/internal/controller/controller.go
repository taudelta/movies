package controller

import (
	"context"
	"errors"

	"movix/metadata/internal/repository"
	"movix/metadata/pkg/model"
)

var ErrNotFound = errors.New("not found")

type metadataRepository interface {
	Get(ctx context.Context, id string) (*model.Metadata, error)
	Put(ctx context.Context, id string, meta *model.Metadata) error
}

type Controller struct {
	repo metadataRepository
}

func New(repo metadataRepository) *Controller {
	return &Controller{
		repo: repo,
	}
}

func (c *Controller) Get(ctx context.Context, id string) (*model.Metadata, error) {
	meta, err := c.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
	}

	return meta, nil
}

func (c *Controller) Put(ctx context.Context, id string, meta *model.Metadata) error {
	err := c.repo.Put(ctx, id, meta)
	return err
}
