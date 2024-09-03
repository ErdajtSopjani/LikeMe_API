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
			if r.URL.Path == "/api/v1/register" ||
				r.URL.Path == "/api/v1/login" ||
				r.URL.Path == "/api/v1/verify" ||
				r.URL.Path == "/api/v1/resend/verification" ||
				r.URL.Path == "/api/v1/resend/login" {
				next.ServeHTTP(w, r)
				return
			}

			// if no token is found
			if token == "" {
				log.Println("Unauthorized request from: ", r.RemoteAddr)
				http.Error(w, "Unauthorized: Token is missing", http.StatusUnauthorized)
				return
			}

			type UserToken struct {
				Token     string
				CreatedAt time.Time
				ExpiresAt time.Time
				UserId    int64
			}

			var userToken UserToken

			// query database to check if token exists
			err := db.Select("token", "expires_at").Where("token = ?", token).First(&userToken).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound { // if token is not found
					log.Println("Unauthorized request from: ", r.RemoteAddr)
					http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
				} else if userToken.ExpiresAt.Before(time.Now()) { // if token is expired
					http.Error(w, "Unauthorized: Token has expired", http.StatusUnauthorized)
				} else {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					log.Printf("\n\nERROR\n\tFailed to query database!\n\t%s\n\n", err)
				}
				return
			}

			// if no errors are found that means the token is valid
			next.ServeHTTP(w, r) // continue to the next handler
		})
	}
}
