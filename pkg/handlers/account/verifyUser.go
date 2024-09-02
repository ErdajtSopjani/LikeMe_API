package handlers

import (
	// "encoding/json"
	// "log"
	"net/http"
)

// VerifyEmail is a handler for verifying a user's email
func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

// SendConfirmation is a handler for sending an email verification
func SendConfirmation(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}
