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
			token := r.Header.Get("Authorization")

			// if its a register request continue to the next handler
			if r.URL.Path == "/api/v1/register" {
				next.ServeHTTP(w, r)
				return
			}

			// if no token is found
			if token == "" {
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
			if err != nil {
				if err == gorm.ErrRecordNotFound { // if token is not found
					log.Println("Unauthorized request from: ", r.RemoteAddr)
					http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
				} else if user.TokenExpiresAt.Before(time.Now()) { // if token is expired
					http.Error(w, "Unauthorized: Token has expired", http.StatusUnauthorized)
				} else {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					log.Fatal("failed to query database:", err)
				}
				return
			}

			// if no errors are found that means the token is valid
			next.ServeHTTP(w, r) // continue to the next handler
		})
	}
}
