package handler

import (
	"encoding/json"
	"errors"
	"movix/metadata/internal/controller"
	"net/http"
)

type Handler struct {
	ctrl *controller.Controller
}

func NewHandler(ctrl *controller.Controller) *Handler {
	return &Handler{
		ctrl: ctrl,
	}
}

func (h *Handler) GetMetadata(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metadata, err := h.ctrl.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, controller.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(metadata)
}
