package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
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
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	reader := bytes.NewReader(jsonBody)
	resp, err := http.Post("http://trip-service:8083/trips/preview", "application/json", reader)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to contact trip service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	var respBody any
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		http.Error(w, "Failed to parse response body", http.StatusBadRequest)
		return
	}
	response := contracts.APIResponse{Data: respBody}
	writeJSON(w, http.StatusOK, response)
}
