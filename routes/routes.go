package routes

import (
	"net/http"

	"github.com/ErdajtSopjani/LikeMe_API/internal/config"
	_middleware "github.com/ErdajtSopjani/LikeMe_API/middleware"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/account"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/email"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/users"
)

func Routes(app *config.AppConfig, db *gorm.DB) http.Handler {
	mux := chi.NewRouter()

	/* Middleware */
	// use cors protection in production
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(_middleware.NoSurf(app.IsProd))
	mux.Use(_middleware.VerifyToken(db))
	// TODO: use cors protection in production

	/* Get Requests */
	mux.Get("/is_running", greeting())
	// TODO: mux.Get("/api/v1/is_verified", userHandlers.IsVerified(db))

	/* Post Requests */
	mux.Post("/api/v1/register", account.RegisterUser(db))
	mux.Post("/api/v1/follow", users.FollowAccount(db))
	mux.Post("/api/v1/login", account.LoginUser(db))
	mux.Post("/api/v1/profile", users.CreateProfile(db))

	// TODO: mux.Post("/api/v1/email/resend/login", email.SendLoginEmail(db))
	mux.Post("/api/v1/email/resend/register", email.ResendVerificationEmail(db))
	mux.Post("/api/v1/email/verify", email.VerifyEmail(db))

	/* Delete Requests */
	// TODO: mux.Delete("/api/v1/unfollow", userHandlers.UnfollowAccount(db))
	// TODO: mux.Delete("/api/v1/user", userHandlers.DeleteUser(db))

	/* Put Requests */
	// TODO: mux.Put("/api/v1/update_user", userHandlers.UpdateUser(db))
	// TODO: mux.Post("/api/v1/profile", userHandlers.UpdateProfile(db)

	return mux
}

func greeting() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("yes"))
	}
}
