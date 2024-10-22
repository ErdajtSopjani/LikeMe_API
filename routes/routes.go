package routes

import (
	"net/http"

	"github.com/ErdajtSopjani/LikeMe_API/internal/config"
	_middleware "github.com/ErdajtSopjani/LikeMe_API/middleware"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"gorm.io/gorm"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/account"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/email/verify"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/social/follows"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/social/messages"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/social/profiles"
)

func Routes(app *config.AppConfig, db *gorm.DB) http.Handler {
	mux := chi.NewRouter()

	/* Middleware */
	// use cors protection in production
	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // Replace with your React app's URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	mux.Use(corsConfig.Handler)

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(_middleware.NoSurf(app.IsProd))
	mux.Use(_middleware.VerifyToken(db))
	// TODO: use cors protection in production

	/* Get Requests */
	mux.Get("/is_running", greeting())

	// social GET requests
	mux.Get("/api/v1/messages", messages.GetUserMessages(db))

	/* Post Requests */

	// auth Post Requests
	mux.Post("/api/v1/auth/register", account.RegisterUser(db))
	mux.Post("/api/v1/auth/login", account.Login(db))

	// social Post Requests
	mux.Post("/api/v1/follow", follows.FollowAccount(db))
	mux.Post("/api/v1/profile", profiles.ManageProfiles(db))

	// email Post Requests
	mux.Post("/api/v1/email/change", account.ChangeEmail(db))
	mux.Post("/api/v1/email/login", account.LoginUser(db))
	mux.Post("/api/v1/email/verify/resend", verify.ResendVerificationEmail(db))
	mux.Post("/api/v1/email/verify", verify.VerifyEmail(db))

	/* Delete Requests */
	mux.Delete("/api/v1/unfollow", follows.UnfollowAccount(db))
	mux.Delete("/api/v1/user", account.DeleteAccount(db))

	/* Put Requests */
	mux.Put("/api/v1/profile", profiles.ManageProfiles(db)) // same route as create

	return mux
}

func greeting() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("yes"))
	}
}
