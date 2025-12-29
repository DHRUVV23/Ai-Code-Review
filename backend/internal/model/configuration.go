package model

import "time"

type Configuration struct {
	ID             int       `json:"id"`
	RepositoryID   int       `json:"repository_id"`
	// Changed back to string to match Database TEXT column
	IgnorePatterns string    `json:"ignore_patterns"` 
	ReviewStyle    string    `json:"review_style"`
	CreatedAt      time.Time `json:"created_at"`
}