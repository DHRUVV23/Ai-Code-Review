package model

import "time"

type Review struct {
	ID           int       `json:"id"`
	RepositoryID int       `json:"repository_id"`
	PRNumber     int       `json:"pr_number"`
	Status       string    `json:"status"` // e.g., "pending", "completed", "failed"
	Content      string    `json:"content"` // The actual AI feedback
	CreatedAt    time.Time `json:"created_at"`
}