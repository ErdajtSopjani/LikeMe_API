package account

import (
	"net/http"

	"gorm.io/gorm"
)

type ConfirmChangeEmailRequest struct {
	Email       string `json:"email"`
	ChangeToken string `json:"token"`
	UserId      int    `json:"user_id"`
}

// ConfirmChangeEmail confirms the email change based on the saved entry
func ConfirmChangeEmail(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
