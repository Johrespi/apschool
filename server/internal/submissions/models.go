package submissions

import (
	"time"
)

type Submission struct {
	ID          int       `json:"id"`
	UserID      int       `json:"-"`
	ChallengeID int       `json:"challenge_id"`
	Code        string    `json:"code"`
	Passed      bool      `json:"passed"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}
