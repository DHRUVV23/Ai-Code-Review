package model

import "time"

type Repository struct {
	ID             int       `json:"id"`
	GithubRepoID   int64     `json:"github_repo_id"`
	InstallationID int64     `json:"installation_id"`
	Name           string    `json:"name"`
	FullName       string    `json:"full_name"`
	Private        bool      `json:"private"`
	CreatedAt      time.Time `json:"created_at"`
}