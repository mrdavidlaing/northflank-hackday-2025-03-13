package main

import (
	"encoding/json"
	"net/http"

	"github.com/mrdavidlaing/simple-go-server/version"
)

// InfoResponse represents the JSON structure for the /info endpoint
type InfoResponse struct {
	Version string `json:"version"`
}

// infoHandler handles requests to the /info endpoint
func infoHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create response
	response := InfoResponse{
		Version: version.Version,
	}

	// Set content type
	w.Header().Set("Content-Type", "application/json")
	
	// Encode and send response
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
} 