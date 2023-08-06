package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"movix/metadata/internal/controller"
	"movix/metadata/pkg/model"
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
	if r.Method == "PUT" {
		h.PutMetadata(w, r)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		log.Println("request error: empty metadata id")
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

func (h *Handler) PutMetadata(w http.ResponseWriter, r *http.Request) {
	var metadata model.Metadata

	if err := json.NewDecoder(r.Body).Decode(&metadata); err != nil {
		log.Println("put metadata decode error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.ctrl.Put(r.Context(), metadata.ID, &metadata); err != nil {
		log.Println("put metadata error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(metadata)
}
