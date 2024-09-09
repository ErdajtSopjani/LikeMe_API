package email

import (
	"log"
	"net/http"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/helpers"
	"gorm.io/gorm"
)

// VerifyEmail awaits for the SendConfirmation email to be verified
func VerifyEmail(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token") // get token from urlQuery

		// decline request if token is missing
		if token == "" {
			log.Printf("\n\nBAD/MALICIOUS\n\tBad credential request from: %v\n\n", r.RemoteAddr)
			helpers.RespondError(w, "Invalid Token", http.StatusBadRequest)
			return
		}

		// check if token is valid
		if err := db.Where("token = ?", token).First(&handlers.VerificationToken{}).Error; err != nil {
			log.Printf("\n\nBAD/MALICIOUS\n\tBad credential request from %s: %v\n\tError: %s\n\n", r.RemoteAddr, token, err)
			helpers.RespondError(w, "Invalid Token", http.StatusBadRequest)
			return
		}

		// get the user_id from the associated token
		verificationToken := handlers.VerificationToken{}
		if err := db.Where("token = ?", token).First(&verificationToken).Error; err != nil {
			log.Printf("\n\nERROR\n\tFailed to query database!\n\t%s\n\n", err)
			helpers.RespondError(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		userId := verificationToken.UserId

		if err := db.Model(&handlers.User{}).Where("id = ?", userId).Update("verified", true).Error; err != nil {
			log.Printf("\n\nERROR\n\tFailed to update user: %v\n\tError: %s\n\n", token, err)
			helpers.RespondError(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// delete the token from the database
		if err := db.Where("token = ?", token).Delete(&handlers.VerificationToken{}).Error; err != nil {
			log.Printf("\n\nERROR\n\tFailed to delete token: %v\n\tError: %s\n\n", token, err)
			helpers.RespondError(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		helpers.RespondJSON(w, http.StatusOK, "User Email Verified")
	}
}
