package routes

import (
	"net/http"

	_middleware "github.com/ErdajtSopjani/LikeMe_API/middleware"
	"github.com/ErdajtSopjani/LikeMe_API/pkg/config"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func Routes(app *config.AppConfig, db *gorm.DB) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(_middleware.NoSurf(app.IsProd))

	mux.Get("/is_running", greeting())

	return mux
}

func greeting() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("yes"))
	}
}
