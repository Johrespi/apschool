package auth

import (
	"context"
	"os"
	"testing"

	"apschool/internal/testutil"
)

var testDB *testutil.TestDB

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	testDB, err = testutil.SetupTestDB(ctx)
	if err != nil {
		panic("failed to setup test db: " + err.Error())
	}

	code := m.Run()

	testDB.Teardown(ctx)
	os.Exit(code)
}

func TestRepository_CreateUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	testDB.TruncateTables(t)
	repo := NewRepository(testDB.DB)

	tests := []struct {
		name      string
		username  string
		email     string
		avatarURL string
		wantErr   bool
	}{
		{
			name:      "create user successfully",
			username:  "testuser",
			email:     "test@example.com",
			avatarURL: "https://example.com/avatar.png",
			wantErr:   false,
		},
		{
			name:      "duplicate email fails",
			username:  "another",
			email:     "test@example.com",
			avatarURL: "https://example.com/avatar2.png",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.CreateUser(context.Background(), tt.username, tt.email, tt.avatarURL)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if user.ID == 0 {
					t.Error("CreateUser() returned user with ID = 0")
				}
				if user.Username != tt.username {
					t.Errorf("CreateUser() username = %v, want %v", user.Username, tt.username)
				}
				if user.Email != tt.email {
					t.Errorf("CreateUser() email = %v, want %v", user.Email, tt.email)
				}
			}
		})
	}
}

func TestRepository_GetUserByID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	testDB.TruncateTables(t)
	repo := NewRepository(testDB.DB)

	// Setup: crear un usuario
	createdUser, err := repo.CreateUser(context.Background(), "testuser", "test@example.com", "avatar.png")
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	tests := []struct {
		name    string
		userID  int
		wantErr bool
	}{
		{
			name:    "user exists",
			userID:  createdUser.ID,
			wantErr: false,
		},
		{
			name:    "user not found",
			userID:  99999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.GetUserByID(context.Background(), tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if user.ID != tt.userID {
					t.Errorf("GetUserByID() ID = %v, want %v", user.ID, tt.userID)
				}
			}
		})
	}
}

func TestRepository_CreateGithubAuth(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	testDB.TruncateTables(t)
	repo := NewRepository(testDB.DB)

	// Setup: crear un usuario
	user, err := repo.CreateUser(context.Background(), "testuser", "test@example.com", "avatar.png")
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	tests := []struct {
		name     string
		userID   int
		githubID int
		wantErr  bool
	}{
		{
			name:     "create github auth successfully",
			userID:   user.ID,
			githubID: 12345,
			wantErr:  false,
		},
		{
			name:     "duplicate github_id fails",
			userID:   user.ID,
			githubID: 12345,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateGithubAuth(context.Background(), tt.userID, tt.githubID)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateGithubAuth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetUserByGithubID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	testDB.TruncateTables(t)
	repo := NewRepository(testDB.DB)

	// Setup: crear usuario y github auth
	user, err := repo.CreateUser(context.Background(), "testuser", "test@example.com", "avatar.png")
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	githubID := 12345
	err = repo.CreateGithubAuth(context.Background(), user.ID, githubID)
	if err != nil {
		t.Fatalf("failed to create github auth: %v", err)
	}

	tests := []struct {
		name     string
		githubID int
		wantErr  bool
	}{
		{
			name:     "user with github auth exists",
			githubID: githubID,
			wantErr:  false,
		},
		{
			name:     "github user not found",
			githubID: 99999,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, err := repo.GetUserByGithubID(context.Background(), tt.githubID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByGithubID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if gotUser.ID != user.ID {
					t.Errorf("GetUserByGithubID() user ID = %v, want %v", gotUser.ID, user.ID)
				}
				if gotUser.Username != user.Username {
					t.Errorf("GetUserByGithubID() username = %v, want %v", gotUser.Username, user.Username)
				}
			}
		})
	}
}
