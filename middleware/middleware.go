package middleware

import (
	"log"
	"net/http"

	"github.com/justinas/nosurf"
	"gorm.io/gorm"
)

// NoSurf adds CSRF protection to all POST requests
func NoSurf(isProd bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if !isProd { // skip csrf protection in development
			return next
		}

		csrfHandler := nosurf.New(next)

		csrfHandler.SetBaseCookie(http.Cookie{
			HttpOnly: true,
			Path:     "/",
			Secure:   isProd, // Conditionally set the Secure flag based on the environment
			SameSite: http.SameSiteLaxMode,
		})
		return csrfHandler
	}
}

func VerifyToken(db *gorm.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is_valid := db.Select("*").Where("users.token = ?", r.Header.Get("Authorization")).Error

			if is_valid != nil {
				log.Fatal("Invalid token")
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			println(db.Select("user").Where("users.token = ?", r.Header.Get("Authorization")))
			next.ServeHTTP(w, r)
		})
	}
}
