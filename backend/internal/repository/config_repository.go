package repository

import (
	"context"
	"fmt"

	"github.com/DHRUVV23/ai-code-review/backend/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConfigRepository struct {
	Pool *pgxpool.Pool
}

func NewConfigRepository(pool *pgxpool.Pool) *ConfigRepository {
	return &ConfigRepository{Pool: pool}
}

// GetByRepoID fetches the config exactly as it is in the DB
func (r *ConfigRepository) GetByRepoID(ctx context.Context, repoID int) (*model.Configuration, error) {
	query := `
		SELECT id, repository_id, review_style, ignore_patterns, created_at
		FROM configurations 
		WHERE repository_id = $1`

	var config model.Configuration
	err := r.Pool.QueryRow(ctx, query, repoID).Scan(
		&config.ID, 
		&config.RepositoryID, 
		&config.ReviewStyle, 
		&config.IgnorePatterns, // Direct string scan
		&config.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	return &config, nil
}

// UpsertConfig updates the config
func (r *ConfigRepository) UpsertConfig(ctx context.Context, config *model.Configuration) error {
	query := `
		INSERT INTO configurations (repository_id, review_style, ignore_patterns, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (repository_id)
		DO UPDATE SET 
			review_style = $2, 
			ignore_patterns = $3, 
			updated_at = NOW()
		RETURNING id`

	return r.Pool.QueryRow(ctx, query, 
		config.RepositoryID, 
		config.ReviewStyle, 
		config.IgnorePatterns, // Direct string insert
	).Scan(&config.ID)
}