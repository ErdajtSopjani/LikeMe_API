package routes

import (
	"net/http"

	"github.com/ErdajtSopjani/LikeMe_API/internal/config"
	_middleware "github.com/ErdajtSopjani/LikeMe_API/middleware"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/account"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/email/verify"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/users/follows"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/users/profiles"
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

	/* Post Requests */
	mux.Post("/api/v1/register", account.RegisterUser(db))
	mux.Post("/api/v1/follow", follows.FollowAccount(db))
	mux.Post("/api/v1/login", account.Login(db))
	mux.Post("/api/v1/profile", profiles.CreateProfile(db))

	mux.Post("/api/v1/email/login", account.LoginUser(db))
	mux.Post("/api/v1/email/resend/register", verify.ResendVerificationEmail(db))
	mux.Post("/api/v1/email/verify", verify.VerifyEmail(db))

	/* Delete Requests */
	mux.Delete("/api/v1/unfollow", follows.UnfollowAccount(db))
	mux.Delete("/api/v1/user", account.DeleteAccount(db))

	/* Put Requests */
	// TODO: mux.Put("/api/v1/update_user", userHandlers.UpdateUser(db))
	// TODO: mux.Post("/api/v1/profile", userHandlers.UpdateProfile(db)

	return mux
}

func greeting() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("yes"))
	}
}
