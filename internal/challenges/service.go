package challenges

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrChallengeNotFound = errors.New("challenge not found")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetChallenges(ctx context.Context, category string) ([]Challenge, error) {
	return s.repo.GetByCategory(ctx, category)
}

func (s *Service) GetChallengeByID(ctx context.Context, id int) (*Challenge, error) {
	challenge, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrChallengeNotFound
		}
		return nil, err
	}

	return challenge, nil
}
