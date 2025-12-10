// Package handler provides HTTP handlers for the notification server.
package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/alessandrolattao/lanotifica/internal/notification"
)

// Notification handles POST requests to send notifications.
func Notification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req notification.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.Message == "" {
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}

	if err := notification.Send(req); err != nil {
		log.Printf("Failed to send notification: %v", err)
		http.Error(w, "Failed to send notification", http.StatusInternalServerError)
		return
	}

	log.Printf("Notification sent: %s - %s", req.Title, req.Message)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "sent"}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
