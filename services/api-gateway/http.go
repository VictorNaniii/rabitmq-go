package main

import (
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_client"
	"ride-sharing/shared/contracts"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var reqBody previewTripRequest

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if reqBody.UserID == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tripService, err := grpc_client.NewTripServiceClient()
	if err != nil {
		log.Fatalln(err)
	}
	defer tripService.Close()
	//tripService.Client.PreviewTrip()
	tripPreview, err := tripService.Client.PreviewTrip(r.Context(), reqBody.toProto())
	if err != nil {
		log.Println("Failed preview a trip ", err)
		http.Error(w, "Failed to preview a trip", http.StatusInternalServerError)
		return
	}
	response := contracts.APIResponse{Data: tripPreview}
	writeJSON(w, http.StatusOK, response)
}

func handleTripStart(w http.ResponseWriter, r *http.Request) {
	var reqBody startTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	tripService, err := grpc_client.NewTripServiceClient()

	if err != nil {
		log.Println(err)
	}
	defer tripService.Close()

	trip, err := tripService.Client.CreateTrip(r.Context(), reqBody.toProto())
	if err != nil {
		log.Println("Failed to start a trip:", err)
		http.Error(w, "Failed to start a trip", http.StatusBadRequest)
		return
	}
	response := contracts.APIResponse{Data: trip}
	writeJSON(w, http.StatusCreated, response)
}
