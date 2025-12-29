package model

import "time"

type Repository struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	Owner     string    `json:"owner"`
	CreatedAt time.Time `json:"created_at"`
}