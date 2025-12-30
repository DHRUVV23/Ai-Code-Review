package model

import "time"

type User struct {
	ID          int       `json:"id"`
	GithubID    int64     `json:"github_id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	AccessToken string    `json:"-"` // New Field! (json:"-" hides it from API responses for security)
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}