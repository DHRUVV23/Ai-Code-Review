package repository

import (
	"context"
	"fmt"

	"github.com/DHRUVV23/ai-code-review/backend/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	Pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{Pool: pool}
}

// UpsertUser saves a user AND their Access Token
// Added 'accessToken' parameter
func (r *UserRepository) UpsertUser(ctx context.Context, githubID int64, username string, email string, accessToken string) (int, error) {
	var id int
	// We now insert access_token and update it on conflict
	query := `
		INSERT INTO users (github_id, username, email, access_token, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		ON CONFLICT (email) DO UPDATE 
		SET github_id = EXCLUDED.github_id,
			username = EXCLUDED.username, 
			access_token = EXCLUDED.access_token, -- Update token if it changed
			updated_at = NOW()
		RETURNING id`

	// Pass accessToken to the query
	err := r.Pool.QueryRow(ctx, query, githubID, username, email, accessToken).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("upsert user failed: %w", err)
	}
	return id, nil
}

// GetUserByID fetches a user including their Access Token
func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	// Added access_token to the SELECT list
	query := `SELECT id, github_id, username, email, access_token, created_at, updated_at FROM users WHERE id = $1`

	var user model.User
	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID, 
		&user.GithubID, 
		&user.Username, 
		&user.Email, 
		&user.AccessToken, // Scan the new field
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