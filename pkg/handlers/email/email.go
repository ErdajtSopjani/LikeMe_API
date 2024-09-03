package email

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ErdajtSopjani/LikeMe_API/pkg/handlers"
	"gorm.io/gorm"
)

type ResendRequest struct {
	Email  string `json:"email"`
	UserId int64  `json:"user_id"`
}

// ResendVerificationEmail resends verification email and deletes the old one from the database
func ResendVerificationEmail(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ResendRequest
		// decode the request body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid Request-Format", http.StatusBadRequest) // if the request body is not valid
			return
		}

		// check if the email exists in the database
		if handlers.CheckUnique(db, "email", req.Email, "users") {
			http.Error(w, "Email not found", http.StatusBadRequest)
			return
		}
		// check if the email and userid match
		if err := db.Where("email = ? AND id = ?", req.Email, req.UserId).First(&handlers.User{}).Error; err != nil {
			http.Error(w, "Email and UserId do not match", http.StatusBadRequest)
			log.Printf("\n\nBAD/MALICIOUS\n\tBad credential request from %s: %v\n\tError: %s\n\n", r.RemoteAddr, req, err)
			return
		}

		// delete the old verification email from the verification_tokens table
		if err := db.Where("user_id= ?", req.UserId).Delete(&handlers.VerificationTokens{}).Error; err != nil {
			log.Panicf("failed to delete the old verification email: %v", err)
			http.Error(w, "Internal Server Error.", http.StatusInternalServerError)
			return
		}

		// send a new verification email
		SendConfirmation(req.Email)

		w.Write([]byte("Email sent"))
	}
}

// VerifyEmail awaits for the SendConfirmation email to be verified
func VerifyEmail(w http.ResponseWriter, r *http.Request) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not implemented yet"))
	}
}

// SendLoginEmail sends an email for logging in
func SendLoginEmail(w http.ResponseWriter, r *http.Request) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not implemented yet"))
	}
}
