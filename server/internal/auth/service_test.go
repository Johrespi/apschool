package auth

import (
	"context"
	"errors"
	"testing"
	"time"
)

type mockRepository struct {
	getUserByGithubIDFunc func(ctx context.Context, githubID int) (*User, error)
	getUserByIDFunc       func(ctx context.Context, id int) (*User, error)
	createUserFunc        func(ctx context.Context, username, email, avatarURL string) (*User, error)
	createGithubAuthFunc  func(ctx context.Context, userID, githubID int) error
}

func (m *mockRepository) GetUserByGithubID(ctx context.Context, githubID int) (*User, error) {
	return m.getUserByGithubIDFunc(ctx, githubID)
}

func (m *mockRepository) GetUserByID(ctx context.Context, id int) (*User, error) {
	return m.getUserByIDFunc(ctx, id)
}

func (m *mockRepository) CreateUser(ctx context.Context, username, email, avatarURL string) (*User, error) {
	return m.createUserFunc(ctx, username, email, avatarURL)
}

func (m *mockRepository) CreateGithubAuth(ctx context.Context, userID, githubID int) error {
	return m.createGithubAuthFunc(ctx, userID, githubID)
}

func TestCreateUserByGithub(t *testing.T) {
	now := time.Now()

	existingUser := &User{
		ID:        1,
		Username:  "existing",
		Email:     "existing@example.com",
		AvatarURL: "https://example.com/avatar.png",
		CreatedAt: now,
		UpdatedAt: now,
	}

	newUser := &User{
		ID:        2,
		Username:  "newuser",
		Email:     "new@example.com",
		AvatarURL: "https://example.com/new-avatar.png",
		CreatedAt: now,
		UpdatedAt: now,
	}

	tests := []struct {
		name      string
		githubID  int
		username  string
		email     string
		avatarURL string
		mock      *mockRepository
		want      *User
		wantErr   bool
	}{
		{
			name:      "user already exists",
			githubID:  12345,
			username:  "existing",
			email:     "existing@example.com",
			avatarURL: "https://example.com/avatar.png",
			mock: &mockRepository{
				getUserByGithubIDFunc: func(ctx context.Context, githubID int) (*User, error) {
					return existingUser, nil
				},
			},
			want:    existingUser,
			wantErr: false,
		},
		{
			name:      "new user created successfully",
			githubID:  67890,
			username:  "newuser",
			email:     "new@example.com",
			avatarURL: "https://example.com/new-avatar.png",
			mock: &mockRepository{
				getUserByGithubIDFunc: func(ctx context.Context, githubID int) (*User, error) {
					return nil, errors.New("not found")
				},
				createUserFunc: func(ctx context.Context, username, email, avatarURL string) (*User, error) {
					return newUser, nil
				},
				createGithubAuthFunc: func(ctx context.Context, userID, githubID int) error {
					return nil
				},
			},
			want:    newUser,
			wantErr: false,
		},
		{
			name:      "error creating user",
			githubID:  11111,
			username:  "failuser",
			email:     "fail@example.com",
			avatarURL: "https://example.com/fail.png",
			mock: &mockRepository{
				getUserByGithubIDFunc: func(ctx context.Context, githubID int) (*User, error) {
					return nil, errors.New("not found")
				},
				createUserFunc: func(ctx context.Context, username, email, avatarURL string) (*User, error) {
					return nil, errors.New("db error")
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:      "error creating github auth",
			githubID:  22222,
			username:  "authfail",
			email:     "authfail@example.com",
			avatarURL: "https://example.com/authfail.png",
			mock: &mockRepository{
				getUserByGithubIDFunc: func(ctx context.Context, githubID int) (*User, error) {
					return nil, errors.New("not found")
				},
				createUserFunc: func(ctx context.Context, username, email, avatarURL string) (*User, error) {
					return newUser, nil
				},
				createGithubAuthFunc: func(ctx context.Context, userID, githubID int) error {
					return errors.New("github auth error")
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(tt.mock)
			got, err := service.CreateUserByGithub(context.Background(), tt.githubID, tt.username, tt.email, tt.avatarURL)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUserByGithub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want == nil && got != nil {
				t.Errorf("CreateUserByGithub() = %v, want nil", got)
				return
			}

			if tt.want != nil {
				if got == nil {
					t.Errorf("CreateUserByGithub() = nil, want %v", tt.want)
					return
				}
				if got.ID != tt.want.ID {
					t.Errorf("CreateUserByGithub() ID = %v, want %v", got.ID, tt.want.ID)
				}
				if got.Username != tt.want.Username {
					t.Errorf("CreateUserByGithub() Username = %v, want %v", got.Username, tt.want.Username)
				}
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	now := time.Now()

	existingUser := &User{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		AvatarURL: "https://example.com/avatar.png",
		CreatedAt: now,
		UpdatedAt: now,
	}

	tests := []struct {
		name    string
		userID  int
		mock    *mockRepository
		want    *User
		wantErr error
	}{
		{
			name:   "user exists",
			userID: 1,
			mock: &mockRepository{
				getUserByIDFunc: func(ctx context.Context, id int) (*User, error) {
					return existingUser, nil
				},
			},
			want:    existingUser,
			wantErr: nil,
		},
		{
			name:   "user not found",
			userID: 999,
			mock: &mockRepository{
				getUserByIDFunc: func(ctx context.Context, id int) (*User, error) {
					return nil, errors.New("not found")
				},
			},
			want:    nil,
			wantErr: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(tt.mock)
			got, err := service.GetUserByID(context.Background(), tt.userID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("GetUserByID() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("GetUserByID() error = %v, wantErr nil", err)
				return
			}

			if got.ID != tt.want.ID {
				t.Errorf("GetUserByID() ID = %v, want %v", got.ID, tt.want.ID)
			}
			if got.Username != tt.want.Username {
				t.Errorf("GetUserByID() Username = %v, want %v", got.Username, tt.want.Username)
			}
		})
	}
}
