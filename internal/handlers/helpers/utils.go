package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondError(w http.ResponseWriter, message string, status int) {
	http.Error(w, message, status)
	log.Printf("Error: %s, Status: %d", message, status)
}

func RespondJSON(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}
