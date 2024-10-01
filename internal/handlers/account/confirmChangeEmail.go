package account

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/helpers"
	"gorm.io/gorm"
)

type ConfirmChangeEmailRequest struct {
	Email       string `json:"email"`
	ChangeToken string `json:"token"`
}

// ConfirmChangeEmail confirms the email change based on the saved entry
func ConfirmChangeEmail(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ConfirmChangeEmailRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helpers.RespondError(w, err.Error(), http.StatusBadRequest)
			log.Printf("\n\nBAD REQUEST\n\tBad request: %s\n\tError: %s\n\n", req, err)
			return
		}

		user, err := helpers.GetUserFromToken(db, r.Header.Get("Authorization"))
		if err != nil {
			helpers.RespondError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// check if the token exists in the database
		var emailChange handlers.EmailChangeRequest
		if err := db.Where("user_id = ? AND change_token = ?", user.ID, req.ChangeToken).First(&emailChange).Error; err != nil {
			helpers.RespondError(w, "Invalid Token", http.StatusBadRequest)
			log.Printf("\n\nBAD/MALICIOUS\n\tBad credential request from %s: %v\n\tError: %s\n\n", r.RemoteAddr, req, err)
			return
		}
	}
}
