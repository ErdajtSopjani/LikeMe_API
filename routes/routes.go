package routes

import (
	"net/http"

	_middleware "github.com/ErdajtSopjani/LikeMe_API/middleware"
	"github.com/ErdajtSopjani/LikeMe_API/pkg/config"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"github.com/ErdajtSopjani/LikeMe_API/pkg/handlers/account"
	"github.com/ErdajtSopjani/LikeMe_API/pkg/handlers/email"
	userHandlers "github.com/ErdajtSopjani/LikeMe_API/pkg/handlers/users"
)

func Routes(app *config.AppConfig, db *gorm.DB) http.Handler {
	mux := chi.NewRouter()

	/* Middleware */
	mux.Use(middleware.Recoverer)
	mux.Use(_middleware.NoSurf(app.IsProd))
	mux.Use(_middleware.VerifyToken(db))

	/* Get Requests */
	mux.Get("/is_running", greeting())
	// TODO: mux.Get("/api/v1/is_verified", userHandlers.IsVerified(db))

	/* Post Requests */
	mux.Post("/api/v1/register", handlers.RegisterUser(db))
	mux.Post("/api/v1/follow", userHandlers.FollowAccount(db))
	// TODO: mux.Post("/api/v1/login", handlers.LoginUser(db))
	mux.Post("/api/v1/profile", userHandlers.CreateProfile(db))

	// TODO: mux.Post("/api/v1/resend/login", email.SendLoginEmail(db))
	mux.Post("/api/v1/resend/verification", email.ResendVerificationEmail(db))

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
