package submissions

import (
	"context"
	"database/sql"
	"errors"
)

var ErrSubmissionNotFound = errors.New("submission not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, s *Submission) error {

	query := `
	INSERT INTO submissions (user_id, challenge_id, code, passed)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (user_id, challenge_id)
	DO UPDATE SET code = $3, passed = $4, updated_at = NOW()
	RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(ctx, query,
		s.UserID,
		s.ChallengeID,
		s.Code,
		s.Passed,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)

}

func (r *Repository) GetByUserAndChallenge(ctx context.Context, userID, challengeID int) (*Submission, error) {

	query := `
	SELECT id, user_id, challenge_id, code, passed, created_at, updated_at
	FROM submissions
	WHERE user_id = $1 AND challenge_id = $2
	`

	var s Submission
	err := r.db.QueryRowContext(ctx, query, userID, challengeID).Scan(
		&s.ID,
		&s.UserID,
		&s.ChallengeID,
		&s.Code,
		&s.Passed,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSubmissionNotFound
		}
		return nil, err
	}

	return &s, nil
}

func (r *Repository) GetByUser(ctx context.Context, userID int) ([]Submission, error) {

	query := `
	SELECT id, user_id, challenge_id, code, passed, created_at, updated_at
	FROM submissions
	WHERE user_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []Submission
	for rows.Next() {
		var s Submission
		if err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.ChallengeID,
			&s.Code,
			&s.Passed,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		submissions = append(submissions, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return submissions, nil

}
