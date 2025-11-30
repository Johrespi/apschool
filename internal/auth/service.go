package auth

import (
	"context"
	"errors"
)

var ErrUserNotFound = errors.New("user not found")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUserByGithub(ctx context.Context, githubID int, username, avatarURL string) (*User, error) {

	// Verify if the user exists
	user, err := s.repo.GetUserByGithubID(ctx, githubID)
	if err == nil {
		return user, nil
	}

	// If not, then create user
	user, err = s.repo.CreateUser(ctx, username, avatarURL)
	if err != nil {
		return nil, err
	}

	// And create a Github user
	err = s.repo.CreateGithubAuth(ctx, user.ID, githubID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetUserByID(ctx context.Context, id int) (*User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}
