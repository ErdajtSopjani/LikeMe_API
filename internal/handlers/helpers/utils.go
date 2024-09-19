package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

// RespondError returns an error message as response
func RespondError(w http.ResponseWriter, message any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{"message": message})
	log.Printf("Error: %s, Status: %d", message, status)
}

// RespondJSON function to return JSON message as response (mainly used on success)
func RespondJSON(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}
