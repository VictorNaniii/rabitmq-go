package http

import (
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/types"
)

type previewTripRequest struct {
	UserID      string           `json:"userId"`
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}

type HttpControler struct {
	Service domain.TripService
}

func NewHttpControler(service domain.TripService) *HttpControler {
	return &HttpControler{Service: service}
}

func (h *HttpControler) HandleTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody previewTripRequest
	ctx := r.Context()
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Println("Invalid request body", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tripModel, err := h.Service.GetRoute(ctx, &reqBody.Pickup, &reqBody.Destination, true)

	if err != nil {
		log.Println("Failed to create trip model", err)
		http.Error(w, "Failed to create trip", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tripModel); err != nil {
		log.Println("Failed to encode response", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
