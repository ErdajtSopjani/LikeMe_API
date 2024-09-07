package email

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers"
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
			http.Error(w, "Internal Server Error.", http.StatusInternalServerError)
			log.Printf("\n\nERROR\n\tFailed to delete verification token: %v\n\tError: %s\n\n", req, err)
			return
		}

		err := SendConfirmation(db, req.Email, req.UserId) // send the email
		if err != nil {
			http.Error(w, "Internal Server Error.", http.StatusInternalServerError)
			log.Printf("\n\nERROR\n\tFailed to send confirmation email: %v\n\tError: %s\n\n", req, err)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Email sent"))
	}
}
