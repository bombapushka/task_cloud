package server

import (
	"cloud/internal/config"
	"cloud/internal/server/handlers"
	"cloud/internal/server/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func SetupServer(cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middlewares.UserMiddleware)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.HomeHandler(w, r, cfg)
	})
	r.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
		handlers.UploadHandler(w, r, cfg)
	})
	r.Get("/download", func(w http.ResponseWriter, r *http.Request) {
		handlers.DownloadHandler(w, r, cfg)
	})

	r.Get("/login", handlers.LoginHandler)
	r.Post("/login", handlers.LoginHandler)
	r.Get("/register", handlers.RegisterHandler)
	r.Post("/register", handlers.RegisterHandler)

	r.Post("/logout", handlers.Logout)

	return r
}
