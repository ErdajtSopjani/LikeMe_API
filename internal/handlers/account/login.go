package account

import (
	"net/http"
	"time"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/helpers"
	"gorm.io/gorm"
)

// Login handles the login process
// it expects a code query parameter
// checks if it's valid and matches everything then returns a token used for auth
func Login(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			helpers.RespondError(w, "Empty Code", http.StatusBadRequest)
			return
		}

		// get the two_factor record from the database and check if it's valid
		twoFactor := handlers.TwoFactor{}

		if err := db.Where("code = ?", code).First(&twoFactor).Error; err != nil {
			helpers.RespondError(w, "Invalid Code", http.StatusBadRequest)
			return
		}

		if twoFactor.ExpiresAt.Before(time.Now()) {
			helpers.RespondError(w, "Code Expired", http.StatusBadRequest)
			return
		}

		// get the user record from the database to check if record is valid and get the user id
		user := handlers.User{}
		if err := db.Where("id = ?", twoFactor.UserId).First(&user).Error; err != nil {
			helpers.RespondError(w, "Invalid User Record", http.StatusBadRequest)
			return
		}

		// delete all login codes associated with the user
		if err := db.Where("user_id = ?", twoFactor.UserId).Delete(&handlers.TwoFactor{}).Error; err != nil {
			helpers.RespondError(w, "Failed to delete old login codes", http.StatusInternalServerError)
			return
		}

		/* Create and save auth token */
		userToken := helpers.GenerateToken()
		if userToken == "" { // if token generation fails
			helpers.RespondError(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// userTokenRecord represents the user_tokens table structure
		userTokenRecord := handlers.UserToken{
			Token:     userToken,
			UserId:    user.ID,
			CreatedAt: &handlers.Now,
			ExpiresAt: &handlers.TokenExpiresAt,
		}
		if err := db.Create(&userTokenRecord).Error; err != nil { // if saving token fails
			helpers.RespondError(w, "Failed to save token", http.StatusInternalServerError)
			return
		}

		helpers.RespondJSON(w, http.StatusOK, "Token: "+userToken)
	}
}
