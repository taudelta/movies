package handler

import (
	"context"
	"errors"
	"movix/api/gen"
	"movix/metadata/internal/controller"
	"movix/metadata/pkg/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler defines a movie metadata gRPC handler.
type Handler struct {
	gen.UnimplementedMetadataServiceServer
	ctrl *controller.Controller
}

// New creates a new movie metadata gRPC handler.
func New(ctrl *controller.Controller) *Handler {
	return &Handler{
		ctrl: ctrl,
	}
}

// GetMetadata returns movie metadata.
func (h *Handler) GetMetadata(ctx context.Context, req *gen.GetMetadataRequest) (*gen.GetMetadataResponse, error) {
	if req == nil || req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	m, err := h.ctrl.Get(ctx, req.Id)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.GetMetadataResponse{
		Result: model.MetadataToProto(m),
	}, nil
}
