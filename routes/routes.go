package routes

import (
	"net/http"

	_middleware "github.com/ErdajtSopjani/LikeMe_API/middleware"
	"github.com/ErdajtSopjani/LikeMe_API/pkg/config"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"github.com/ErdajtSopjani/LikeMe_API/pkg/handlers/account"
	userHandlers "github.com/ErdajtSopjani/LikeMe_API/pkg/handlers/users"
)

func Routes(app *config.AppConfig, db *gorm.DB) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)

	// TODO: use NoSurf in prod environment.
	// mux.Use(_middleware.NoSurf(app.IsProd))

	// TODO: use VerifyToken in every request except for /api/v1/register
	mux.Use(_middleware.VerifyToken(db))

	mux.Get("/is_running", greeting())

	mux.Post("/api/v1/register", handlers.RegisterUser(db))
	mux.Post("/api/v1/follow", userHandlers.FollowAccount(db))

	return mux
}

func greeting() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("yes"))
	}
}
