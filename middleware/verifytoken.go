package middleware

import (
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

// VerifyToken verifies if the token from each request is valid
func VerifyToken(db *gorm.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")              // get token from request header
			registerRequest := r.URL.Path == "/api/v1/register" // check if request is register request

			if token == "" && !registerRequest { // if token is missing and not register request
				log.Println("Unauthorized request from: ", r.RemoteAddr)
				http.Error(w, "Unauthorized: Token is missing", http.StatusUnauthorized)
				return
			}

			type User struct {
				Token          string
				TokenExpiresAt time.Time
			}

			var user User

			// query database to check if token exists
			err := db.Select("token", "token_expires_at").Where("token = ?", token).First(&user).Error

			if err != nil && !registerRequest { // if token doesn't exist

				// token not found
				if err == gorm.ErrRecordNotFound {
					log.Println("Unauthorized request from: ", r.RemoteAddr)
					http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
				} else {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					log.Fatal("failed to query database:", err)
				}

				// check if token is expired
				if user.TokenExpiresAt.Before(time.Now()) {
					http.Error(w, "Unauthorized: Token has expired", http.StatusUnauthorized)
				}
				return
			}

			// if token is valid and not expired proceed
			next.ServeHTTP(w, r)
		})
	}
}
