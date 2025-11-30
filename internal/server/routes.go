package server

import (
	"net/http"

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

	return r
}

func (s *Server) ping(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"message": "pong"}
	helper.WriteJSON(w, http.StatusOK, resp, nil)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	helper.WriteJSON(w, http.StatusOK, s.db.Health(), nil)
}
