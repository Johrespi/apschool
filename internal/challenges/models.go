package challenges

import "time"

type Challenge struct {
	ID          int       `json:"id"`
	Slug        string    `json:"slug"`
	Category    string    `json:"category"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Template    string    `json:"template"`
	TestCode    string    `json:"test_code"`
	Hints       string    `json:"hints"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
