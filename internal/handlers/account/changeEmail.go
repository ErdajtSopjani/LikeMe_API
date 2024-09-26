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

		// check if email is already taken
		if !helpers.CheckUnique(db, "email", req.Email, "users") {
			helpers.RespondError(w, "Email already taken", http.StatusBadRequest)
			return
		}

		// get the userId from the token
		var userToken handlers.UserToken
		if err := db.Select("user_id").Where("token = ?", r.Header.Get("Authorization")).First(&userToken).Error; err != nil {
			log.Printf("\n\nERROR\n\tFailed to query database!\n\t%s\n\n", err)
			helpers.RespondError(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// check if the user exists
		if helpers.CheckUnique(db, "id", userToken.UserId, "users") {
			helpers.RespondError(w, "User does not exist", http.StatusBadRequest)
			return
		}

		// generate a token
		changeToken := helpers.GenerateToken()
		if changeToken == "" {
			helpers.RespondError(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// save the new email and token to the database
		changeEmail := handlers.EmailChange{
			Email:       req.Email,
			ChangeToken: changeToken,
			UserId:      userToken.UserId,
		}
		if err := db.Create(&changeEmail).Error; err != nil {
			helpers.RespondError(w, "Failed to save email change", http.StatusInternalServerError)
			return
		}

		// get the old email
		var oldEmail string
		if err := db.Table("users").Select("email").Where("id = ?", userToken.UserId).Scan(&oldEmail).Error; err != nil {
			helpers.RespondError(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err := email.SendChangeEmail(changeEmail.Email, oldEmail)
		if err != nil {
			helpers.RespondError(w, "Failed to send email", http.StatusInternalServerError)
			return
		}

		helpers.RespondJSON(w, http.StatusCreated, "Email Confirmation Sent")
	}
}
