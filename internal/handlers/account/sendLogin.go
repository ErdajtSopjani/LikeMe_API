package account

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/email"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email string `json:"email"`
}

// LoginUser is a handlers for login requests that returns the user token
func LoginUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		// validate request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid Request-Format", http.StatusBadRequest)
			log.Printf("\n\nBAD REQUEST\n\tBad request: %v\n\tError: %s\n\n", req, err)
			return
		}

		// get user with the associated email
		var user handlers.User
		if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
			http.Error(w, "User not found", http.StatusBadRequest)
			log.Printf("\n\nBAD/MALICIOUS\n\tBad credential request from %s: %v\n\tError: %s\n\n", r.RemoteAddr, req, err)
			return
		} else if !user.Verified {
			http.Error(w, "Email not verified", http.StatusBadRequest)
			return
		}

		// send login email
		if err := email.SendLoginEmail(db, user.ID, user.Email); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("\n\nERROR\n\tFailed to send login email: %v\n\tError: %s\n\n", user, err)
			return
		}

		// create and save a token for the user
		userToken := &handlers.UserTokens{
			Token:  handlers.GenerateToken(),
			UserId: user.ID,
		}
		if err := db.Create(&userToken).Error; err != nil || userToken.Token == "" {
			http.Error(w, "Internal Server Error...", http.StatusInternalServerError)
			log.Printf("\n\nERROR\n\tFailed to create token: %v\n\tError: %s\n\n", userToken, err)
			return
		}

		w.Write([]byte("Login Email Sent"))
	}
}
