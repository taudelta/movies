package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"movix/rating/internal/controller"
	"movix/rating/pkg/model"
)

type Handler struct {
	ctrl *controller.Controller
}

func NewHandler(ctrl *controller.Controller) *Handler {
	return &Handler{
		ctrl: ctrl,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	recordType := r.FormValue("type")
	if recordType == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet || r.Method == http.MethodPost {
		rating, err := h.ctrl.GetAggregatedRating(r.Context(), model.RecordType(recordType), model.RecordID(id))
		if err != nil {
			if errors.Is(err, controller.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(rating)

	} else if r.Method == http.MethodPut {
		userID := model.UserID(r.FormValue("userId"))
		v, err := strconv.ParseFloat(r.FormValue("value"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = h.ctrl.Put(r.Context(), model.RecordType(recordType), model.RecordID(id),
			model.Rating{
				UserID: userID,
				Value:  model.RatingValue(v),
			})
		if err != nil {
			if errors.Is(err, controller.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
