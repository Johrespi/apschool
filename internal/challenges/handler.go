package challenges

import (
	"apschool/internal/response"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
	logger  *slog.Logger
}

func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) ListChallenges(w http.ResponseWriter, r *http.Request) {

	category := r.URL.Query().Get("category")
	if category == "" {
		response.BadRequest(w, "category is required")
		return
	}

	challenges, err := h.service.GetChallenges(r.Context(), category)
	if err != nil {
		response.ServerError(w, r, h.logger, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Envelope{"challenges": challenges}, nil)

}

func (h *Handler) GetChallenge(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "invalid id")
		return
	}

	challenge, err := h.service.GetChallengeByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrChallengeNotFound) {
			response.NotFound(w)
			return
		}
		response.ServerError(w, r, h.logger, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Envelope{"challenge": challenge}, nil)
}
