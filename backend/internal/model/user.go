package model

import "time"

type User struct {
	ID          int       `json:"id"`
	GithubID    int64     `json:"github_id"`
	Username    string    `json:"username"` // Frontend receives this as "username"
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	AccessToken string    `json:"-"`        // Security: Never send the token to the frontend
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}