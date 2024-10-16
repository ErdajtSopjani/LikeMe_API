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

type ChangeEmailRequest struct {
	Email string `json:"email"`
}

// ChangeEmail saves a new entry to the database with the new email and a token
func ChangeEmail(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ChangeEmailRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helpers.RespondError(w, err.Error(), http.StatusBadRequest)
			log.Printf("\n\nBAD REQUEST\n\tBad request: %s\n\tError: %s\n\n", req, err)
			return
		}

		// get the userId from the token
		var user handlers.User
		user, err := helpers.GetUserFromToken(db, r.Header.Get("Authorization"))
		if err != nil {
			helpers.RespondError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// check if the email is the same
		if user.Email == req.Email {
			helpers.RespondJSON(w, http.StatusBadRequest, "Email is the same")
			return
		}

		// check if email is already taken
		if !helpers.CheckUnique(db, "email", req.Email, "users") {
			helpers.RespondJSON(w, http.StatusBadRequest, "Email already taken")
			return
		}

		// generate a token
		changeToken := helpers.GenerateToken()
		if changeToken == "" {
			helpers.RespondError(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// save the new email and token to the database
		changeEmail := handlers.EmailChangeRequest{
			Email:       req.Email,
			ChangeToken: changeToken,
			UserId:      user.ID,
		}
		if err := db.Create(&changeEmail).Error; err != nil {
			helpers.RespondError(w, "Failed to save email change", http.StatusInternalServerError)
			return
		}

		// get the old email
		var oldEmail string
		if err := db.Table("users").Select("email").Where("id = ?", user.ID).Scan(&oldEmail).Error; err != nil {
			helpers.RespondError(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// send an email confirmation to change email
		err = email.SendChangeEmail(changeToken, oldEmail, changeEmail.Email)
		if err != nil {
			helpers.RespondError(w, "Failed to send email", http.StatusInternalServerError)
			return
		}

		helpers.RespondJSON(w, http.StatusCreated, "Email Confirmation Sent")
	}
}
