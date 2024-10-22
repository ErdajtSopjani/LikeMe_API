package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/helpers"
	"gorm.io/gorm"
)

// VerifyToken verifies if the token from each request is valid
func VerifyToken(db *gorm.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")

			// if its a register/login/email request continue to the next handler
			if strings.HasPrefix(r.URL.Path, "/api/v1/email") || r.URL.Path == "/api/v1/auth/register" || r.URL.Path == "/api/v1/auth/login" {
				next.ServeHTTP(w, r)
				return
			}

			// if no token is found
			if token == "" {
				log.Println("Unauthorized request from: ", r.RemoteAddr)
				helpers.RespondError(w, "Unauthorized: Token is missing", http.StatusUnauthorized)
				return
			}

			type UserToken struct {
				Token     string
				CreatedAt *time.Time
				ExpiresAt *time.Time
				UserId    int64
			}

			var userToken UserToken

			// check if token is valid
			if err := db.Select("token", "expires_at").Where("token = ?", token).First(&userToken).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					log.Println("Unauthorized request from: ", r.RemoteAddr)
					helpers.RespondError(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
				} else {
					helpers.RespondError(w, "Internal Server Error", http.StatusInternalServerError)
					log.Printf("\n\nERROR\n\tFailed to query database!\n\t%s\n\n", err)
				}
				return
			}

			// parse the token expires at from string to time
			tokenExpiresAt, err := time.Parse("2006-01-02 15:04:05.999999 -0700 MST", userToken.ExpiresAt.String())
			if err != nil {
				log.Println("Failed to parse token expiration time: ", err)
				helpers.RespondError(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// query database to check if token exists
			err = db.Select("token", "expires_at").Where("token = ?", token).First(&userToken).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound { // if token is not found
					log.Println("Unauthorized request from: ", r.RemoteAddr)
					helpers.RespondError(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
				} else if tokenExpiresAt.Before(time.Now()) { // if token is expired
					helpers.RespondError(w, "Unauthorized: Token has expired", http.StatusUnauthorized)
				} else {
					helpers.RespondError(w, "Internal Server Error", http.StatusInternalServerError)
					log.Printf("\n\nERROR\n\tFailed to query database!\n\t%s\n\n", err)
				}
				return
			}

			// if no errors are found that means the token is valid
			next.ServeHTTP(w, r) // continue to the next handler
		})
	}
}
