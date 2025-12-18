package model

import "time"

type Configuration struct {
	ID            int       `json:"id"`
	RepositoryID  int       `json:"repository_id"`
	IgnorePatterns []string `json:"ignore_patterns"` // e.g., ["*.lock", "assets/*"]
	ReviewStyle   string    `json:"review_style"`    // e.g., "concise", "detailed"
	Enabled       bool      `json:"enabled"`
	CreatedAt     time.Time `json:"created_at"`
}