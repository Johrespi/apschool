package challenges

import "time"

type Challenge struct {
	ID          int       `json:"id"`
	Slug        string    `json:"slug,omitzero"`
	Category    string    `json:"category,omitzero"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitzero"`
	Template    string    `json:"template,omitzero"`
	TestCode    string    `json:"test_code,omitzero"`
	Hints       string    `json:"hints,omitzero"`
	IsActive    bool      `json:"-"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}
