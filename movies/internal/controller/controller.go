package controller

import (
	"context"
	"errors"
	"log"
	"movix/movies/pkg/model"

	metadatamodel "movix/metadata/pkg/model"
	"movix/movies/internal/gateway"
	ratingmodel "movix/rating/pkg/model"
)

var ErrNotFound = errors.New("movie metadata not found")

type Controller struct {
	ratingGateway   ratingGateway
	metadataGateway metadataGateway
}

type ratingGateway interface {
	GetAggregatedRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType) (float64, error)
	PutRating(ctx context.Context, recordID ratingmodel.
		RecordID, recordType ratingmodel.RecordType, userID ratingmodel.UserID, ratingValue ratingmodel.RatingValue) error
}

type metadataGateway interface {
	Get(ctx context.Context, id string) (*metadatamodel.Metadata, error)
}

func New(ratingGateway ratingGateway, metadataGateway metadataGateway) *Controller {
	return &Controller{
		ratingGateway,
		metadataGateway,
	}
}

func (c *Controller) Get(ctx context.Context, id string) (*model.MovieDetails, error) {
	metadata, err := c.metadataGateway.Get(ctx, id)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	details := &model.MovieDetails{
		Metadata: *metadata,
	}

	rating, err := c.ratingGateway.GetAggregatedRating(ctx, ratingmodel.RecordID(id), ratingmodel.RecordTypeMovie)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		// Just proceed in this case, it's ok not to have ratings yet.
		log.Println("rating not found for movie with id", ratingmodel.RecordID(id))
	} else if err != nil {
		return nil, err
	}

	details.Rating = &rating

	return details, nil
}
