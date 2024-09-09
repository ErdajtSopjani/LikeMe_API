package account

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/email"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/helpers"
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
			helpers.RespondError(w, "Invalid Request-Format", http.StatusBadRequest)
			log.Printf("\n\nBAD REQUEST\n\tBad request: %v\n\tError: %s\n\n", req, err)
			return
		}

		// get user with the associated email
		var user handlers.User
		if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
			helpers.RespondError(w, "User not found", http.StatusBadRequest)
			log.Printf("\n\nBAD/MALICIOUS\n\tBad credential request from %s: %v\n\tError: %s\n\n", r.RemoteAddr, req, err)
			return
		} else if !user.Verified {
			helpers.RespondError(w, "Email not verified", http.StatusBadRequest)
			return
		}

		// send login email
		if err := email.SendLoginEmail(db, user.ID, user.Email); err != nil {
			helpers.RespondError(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("\n\nERROR\n\tFailed to send login email: %v\n\tError: %s\n\n", user, err)
			return
		}

		helpers.RespondJSON(w, http.StatusOK, "Login Email Sent")
	}
}
