package response

import (
	"log/slog"
	"net/http"
)

func ErrorResponse(w http.ResponseWriter, status int, message any) {
	err := WriteJSON(w, status, map[string]any{"error": message}, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ServerError(w http.ResponseWriter, r *http.Request, logger *slog.Logger, err error) {
	logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
	ErrorResponse(w, http.StatusInternalServerError, "internal server error")
}

func ValidationError(w http.ResponseWriter, errors map[string]string) {
	ErrorResponse(w, http.StatusUnprocessableEntity, errors)
}
