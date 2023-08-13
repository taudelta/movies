package handler

import (
	"context"
	"errors"
	"movix/api/gen"
	"movix/rating/internal/controller"
	"movix/rating/pkg/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler defines a movie rating gRPC handler.
type Handler struct {
	gen.UnimplementedRatingServiceServer
	ctrl *controller.Controller
}

// New creates a new movie rating gRPC handler.
func New(ctrl *controller.Controller) *Handler {
	return &Handler{
		ctrl: ctrl,
	}
}

// GetAggregatedRating
func (h *Handler) GetAggregatedRating(ctx context.Context, req *gen.GetAggregatedRatingRequest) (*gen.GetAggregatedRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.RecordType == "" {
		// FIXME: book - не хватает сообщения о пустом record_id, record_type
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty record_id or empty record_type")
	}

	m, err := h.ctrl.GetAggregatedRating(ctx, model.RecordType(req.RecordType), model.RecordID(req.RecordId))
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.GetAggregatedRatingResponse{
		RatingValue: float64(m),
	}, nil
}

func (h *Handler) PutRating(ctx context.Context, req *gen.PutRatingRequest) (*gen.PutRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.RecordType == "" || req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty record_id or empty record_type or empty user_id")
	}

	err := h.ctrl.Put(ctx, model.RecordType(req.RecordType), model.RecordID(req.RecordId), model.Rating{
		UserID: model.UserID(req.UserId),
		Value:  model.RatingValue(req.RatingValue),
	})
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.PutRatingResponse{}, nil
}
