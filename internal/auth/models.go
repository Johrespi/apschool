package auth

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
}

type UserAuthGithub struct {
	UserID   int `json:"user_id"`
	GithubID int `json:"github_id"`
}
