package submissions

import (
	"context"
	"errors"
)

var ErrSubmissionNotPassed = errors.New("submission not found")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateSubmission(ctx context.Context, userID int, submission *Submission) error {
	if !submission.Passed {
		return ErrSubmissionNotPassed
	}

	submission.UserID = userID

	return s.repo.Create(ctx, submission)
}

func (s *Service) GetUserSubmission(ctx context.Context, userID, challengeID int) (*Submission, error) {
	return s.repo.GetByUserAndChallenge(ctx, userID, challengeID)
}

func (s *Service) GetUserSubmissions(ctx context.Context, userID int) ([]Submission, error) {
	return s.repo.GetByUser(ctx, userID)
}
