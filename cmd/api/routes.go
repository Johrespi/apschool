package main

import (
	"apschool/internal/response"
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/", app.ping)
	r.Get("/health", app.health)

	//Auth routes
	r.Get("/api/auth/github/login", app.auth.GithubLogin)
	r.Get("/api/auth/github/callback", app.auth.GithubCallback)

	return r
}

func (app *application) ping(w http.ResponseWriter, r *http.Request) {
	response.WriteJSON(w, http.StatusOK, response.Envelope{"message": "pong"}, nil)
}

func (app *application) health(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
	err := app.db.PingContext(ctx)
	if err != nil {
		response.WriteJSON(w, http.StatusServiceUnavailable, response.Envelope{"status": "down"}, nil)
		return
	}
	response.WriteJSON(w, http.StatusOK, response.Envelope{"status": "up"}, nil)
}
