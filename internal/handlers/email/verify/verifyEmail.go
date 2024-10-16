package verify

import (
	"log"
	"net/http"
	"time"

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

		var verificationToken handlers.VerificationToken

		// check if token is valid
		if err := db.Where("token = ?", token).First(&verificationToken).Error; err != nil {
			log.Printf("\n\nBAD/MALICIOUS\n\tBad credential request from %s: %v\n\tError: %s\n\n", r.RemoteAddr, token, err)
			helpers.RespondError(w, "Invalid Token", http.StatusBadRequest)
			return
		}

		// check if token is expired
		if verificationToken.ExpiresAt.Before(time.Now()) {
			log.Printf("\n\nBAD/MALICIOUS\n\tBad credential request from %s: %v\n\tError: %s\n\n", r.RemoteAddr, token, "Token Expired")
			helpers.RespondError(w, "Token Expired", http.StatusBadRequest)
			return
		}

		// set the user as verified
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

		helpers.RespondJSON(w, http.StatusOK, "Email Verified")
	}
}
