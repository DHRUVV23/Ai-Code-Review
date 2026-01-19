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

func (r *UserRepository) UpsertUser(ctx context.Context, githubID int64, username string, email string, accessToken string) (int, error) {
	var id int

	query := `
		INSERT INTO users (github_id, username, email, access_token, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		ON CONFLICT (github_id) DO UPDATE 
		SET username = EXCLUDED.username, 
			email = EXCLUDED.email,
			access_token = EXCLUDED.access_token,
			updated_at = NOW()
		RETURNING id`

	err := r.Pool.QueryRow(ctx, query, githubID, username, email, accessToken).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("upsert user failed: %w", err)
	}
	return id, nil
}

// GetUserByID fetches a user from the 'users' table
func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*model.User, error) {

	query := `SELECT id, github_id, username, email, access_token, created_at, updated_at FROM users WHERE id = $1`

	var user model.User
	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID, 
		&user.GithubID, 
		&user.Username, 
		&user.Email, 
		&user.AccessToken,
		&user.CreatedAt, 
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	user.Name = user.Username
	return &user, nil
}