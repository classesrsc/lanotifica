package handler

import (
	"encoding/json"
	"net/http"
)

// HealthResponse is the response for the health endpoint.
type HealthResponse struct {
	Status string `json:"status"`
}

// Health returns a handler for the health check endpoint.
func Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
}
