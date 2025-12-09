package submissions

import (
	"apschool/internal/middleware"
	"apschool/internal/response"
	"apschool/internal/validator"
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

func (h *Handler) CreateSubmissionsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	var submission Submission

	if err := response.ReadJSON(w, r, &submission); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	v := validator.New()
	v.Check(submission.ChallengeID != 0, "challenge_id", "is required")
	v.Check(validator.NotBlank(submission.Code), "code", "is required")

	if !v.Valid() {
		response.ValidationError(w, v.Errors)
		return
	}

	err := h.service.CreateSubmission(r.Context(), userID, &submission)
	if err != nil {
		if errors.Is(err, ErrSubmissionNotPassed) {
			response.BadRequest(w, "submission did not pass the tests")
			return
		}
		response.ServerError(w, r, h.logger, err)
		return
	}

	response.WriteJSON(w, http.StatusCreated, response.Envelope{"submission": submission}, nil)
}

func (h *Handler) GetSubmissionsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	submissions, err := h.service.GetUserSubmissions(r.Context(), userID)
	if err != nil {
		response.ServerError(w, r, h.logger, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Envelope{"submissions": submissions}, nil)
}

func (h *Handler) GetSubmissionHandler(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middleware.UserIDKey).(int)
	challengeIDStr := chi.URLParam(r, "challenge_id")
	challengeID, err := strconv.Atoi(challengeIDStr)
	if err != nil {
		response.BadRequest(w, "invalid challenge_id")
	}

	submission, err := h.service.GetUserSubmission(r.Context(), userID, challengeID)
	if err != nil {
		if errors.Is(err, ErrSubmissionNotFound) {
			response.NotFound(w)
			return
		}
		response.ServerError(w, r, h.logger, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Envelope{"submission": submission}, nil)
}
