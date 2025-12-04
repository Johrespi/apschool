package middleware

import (
	"apschool/internal/auth"
	"apschool/internal/response"
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "userID"

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Read Authorization from header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Unauthorized(w)
			return
		}

		// Get token (Bearer <token>)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(w)
			return
		}

		token := parts[1]

		// Validate token
		userID, err := auth.ValidateJWT(token)
		if err != nil {
			response.Unauthorized(w)
			return
		}

		// Save to context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
