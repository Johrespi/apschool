package server

import (
	"net/http"

	"apschool/internal/auth"
	"apschool/internal/helper"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/", s.ping)

	r.Get("/health", s.healthHandler)

	// Auth routes
	authHandler := auth.NewHandler(auth.NewService(auth.NewRepository(s.db.GetDB())))

	r.Get("/api/auth/github", authHandler.GithubLogin)
	r.Get("/api/auth/github/callback", authHandler.GithubCallback)
	r.Get("/api/auth/me", authHandler.GetMe)

	return r
}

func (s *Server) ping(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"message": "pong"}
	helper.WriteJSON(w, http.StatusOK, resp, nil)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	helper.WriteJSON(w, http.StatusOK, s.db.Health(), nil)
}
