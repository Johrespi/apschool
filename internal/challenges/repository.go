package challenges

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

func (r *Repository) GetByCategory(ctx context.Context, category string) ([]Challenge, error) {

	query := `SELECT id, title FROM challenges
	WHERE category = $1 AND is_active = true
	ORDER BY slug`

	rows, err := r.db.QueryContext(ctx, query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var challenges []Challenge
	for rows.Next() {
		var c Challenge
		if err := rows.Scan(&c.ID, &c.Title); err != nil {
			return nil, err
		}
		challenges = append(challenges, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return challenges, nil
}

func (r *Repository) GetByID(ctx context.Context, id int) (*Challenge, error) {

	query := `SELECT id, slug, category, title, description, template, test_code, hints, is_active, created_at, updated_at
	FROM challenges
	WHERE id = $1 AND is_active = true
	`

	var c Challenge

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID,
		&c.Slug,
		&c.Category,
		&c.Title,
		&c.Description,
		&c.Template,
		&c.TestCode,
		&c.Hints,
		&c.IsActive,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &c, nil

}
