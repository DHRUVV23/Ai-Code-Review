package repository

import (
	"context"
	"fmt"
	"time"

	// "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID        int       `json:"id"`
	GithubID  int64     `json:"github_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRepository struct {
	Pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{Pool: pool}
}

// UpsertUser saves a user if they are new, or updates them if they exist
func (r *UserRepository) UpsertUser(ctx context.Context, githubID int64, username string, email string) (*User, error) {
	query := `
		INSERT INTO users (github_id, username, email, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (github_id) 
		DO UPDATE SET username = EXCLUDED.username, email = EXCLUDED.email
		RETURNING id, github_id, username, email, created_at
	`
	
	var user User
	err := r.Pool.QueryRow(ctx, query, githubID, username, email).Scan(
		&user.ID, &user.GithubID, &user.Username, &user.Email, &user.CreatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to upsert user: %w", err)
	}
	return &user, nil
}