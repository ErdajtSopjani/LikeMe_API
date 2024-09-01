package middleware

import (
	"github.com/justinas/nosurf"
	"net/http"
)

// NoSurf adds CSRF protection to all POST requests
func NoSurf(isProd bool) func(next http.Handler) http.Handler {
	if !isProd { // skip crsf protection if in development mode
		println("skipping csrf protection")
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	return func(next http.Handler) http.Handler {
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
