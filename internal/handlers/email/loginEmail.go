package email

import (
	"net/http"
)

// SendLoginEmail sends an email for logging in
func SendLoginEmail(w http.ResponseWriter, r *http.Request) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not implemented yet"))
	}
}
