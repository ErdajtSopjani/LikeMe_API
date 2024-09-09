package account

import (
	"log"
	"net/http"
	"time"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers"
	"gorm.io/gorm"
)

// Login handles the login process
// it expects a code query parameter
// checks if it's valid and matches everything then returns a token used for auth
func Login(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Empty Code", http.StatusBadRequest)
			return
		}

		// get the two_factor record from the database and check if it's valid
		twoFactor := handlers.TwoFactor{}

		if err := db.Where("code = ?", code).First(&twoFactor).Error; err != nil {
			http.Error(w, "Invalid Code", http.StatusBadRequest)
			return
		} else if twoFactor.ExpiresAt.Before(time.Now()) {
			log.Println(twoFactor.ExpiresAt)
			http.Error(w, "Code Expired", http.StatusBadRequest)
			return
		}

		// get the user record from the database to check if record is valid and get the user id
		user := handlers.User{}
		if err := db.Where("id = ?", twoFactor.UserId).First(&user).Error; err != nil {
			http.Error(w, "Invalid Record", http.StatusBadRequest)
			return
		}

		/* Create and save auth token */
		userToken := handlers.GenerateToken()
		if userToken == "" { // if token generation fails
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// userTokenRecord represents the user_tokens table structure
		userTokenRecord := handlers.UserToken{
			Token:     userToken,
			UserId:    user.ID,
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30), // token expires in 30 days
		}
		if err := db.Create(&userTokenRecord).Error; err != nil { // if saving token fails
			http.Error(w, "Failed to save token", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(userToken))
	}
}
