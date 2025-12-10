package auth

import (
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, username, email, avatarURL string) (*User, error) {

	var user User
	query := `INSERT INTO users (username, email, avatar_url)
	VALUES ($1, $2, $3)
	RETURNING id, username, email, avatar_url, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query, username, email, avatarURL).Scan(&user.ID, &user.Username, &user.Email, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err

	}

	return &user, nil

}

func (r *Repository) GetUserByID(ctx context.Context, id int) (*User, error) {
	var user User

	query := `SELECT id, username, email, avatar_url, created_at, updated_at FROM users
	WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetUserByGithubID(ctx context.Context, githubID int) (*User, error) {

	var user User

	query := `SELECT u.id, u.username, u.email, u.avatar_url, u.created_at, u.updated_at
	FROM users u
	JOIN user_auth_github g ON u.id = g.user_id
	WHERE g.github_id = $1`

	err := r.db.QueryRowContext(ctx, query, githubID).Scan(&user.ID, &user.Username, &user.Email, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) CreateGithubAuth(ctx context.Context, userID, githubID int) error {

	q := `INSERT INTO user_auth_github (user_id, github_id)
	VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, q, userID, githubID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	// TODO: implement this method

	return nil, nil
}
