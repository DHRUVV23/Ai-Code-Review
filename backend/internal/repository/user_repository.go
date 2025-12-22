package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/DHRUVV23/ai-code-review/backend/internal/model"
)

type UserRepository struct {
	Pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{Pool: pool}
}

// UpsertUser saves a user.
// FIX: We now accept githubID and username to match auth.go
func (r *UserRepository) UpsertUser(ctx context.Context, githubID int64, username string, email string) (int, error) {
	var id int
	query := `
		INSERT INTO users (github_id, username, email, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		ON CONFLICT (email) DO UPDATE 
		SET github_id = EXCLUDED.github_id,
		    username = EXCLUDED.username, 
		    updated_at = NOW()
		RETURNING id`

	// We pass all 3 variables to the query now
	err := r.Pool.QueryRow(ctx, query, githubID, username, email).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetUserByID fetches a user by their internal database ID
func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	query := `SELECT id, github_id, username, email, created_at, updated_at FROM users WHERE id = $1`

	var user model.User
	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID, 
		&user.GithubID, 
		&user.Username, 
		&user.Email, 
		&user.CreatedAt, 
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}