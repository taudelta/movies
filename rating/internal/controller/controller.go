package controller

import (
	"context"
	"errors"

	"movix/rating/internal/repository"
	"movix/rating/pkg/model"
)

var ErrNotFound = errors.New("not found")

type ratingRepository interface {
	Get(ctx context.Context, recordType model.RecordType, id model.RecordID) ([]model.Rating, error)
	Put(ctx context.Context, recordType model.RecordType, id model.RecordID, rating model.Rating) error
}

type Controller struct {
	repo ratingRepository
}

func New(repo ratingRepository) *Controller {
	return &Controller{
		repo: repo,
	}
}

func (c *Controller) Get(ctx context.Context, recordType model.RecordType, id model.RecordID) ([]model.Rating, error) {
	ratings, err := c.repo.Get(ctx, recordType, id)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) || errors.Is(err, repository.ErrRecordTypeNotFound) {
			return nil, ErrNotFound
		}
	}

	return ratings, nil
}

func (c *Controller) GetAggregatedRating(ctx context.Context, recordType model.RecordType, id model.RecordID) (float32, error) {
	ratings, err := c.repo.Get(ctx, recordType, id)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) || errors.Is(err, repository.ErrRecordTypeNotFound) {
			return 0, ErrNotFound
		}
	}

	sum := float32(0)
	for _, r := range ratings {
		sum += float32(r.Value)
	}

	return float32(sum) / float32(len(ratings)), nil
}

func (c *Controller) Put(ctx context.Context, recordType model.RecordType, id model.RecordID, rating model.Rating) error {
	err := c.repo.Put(ctx, recordType, id, rating)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) || errors.Is(err, repository.ErrRecordTypeNotFound) {
			return ErrNotFound
		}
	}

	return nil
}
