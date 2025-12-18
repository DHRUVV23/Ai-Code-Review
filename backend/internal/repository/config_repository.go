package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/DHRUVV23/ai-code-review/backend/internal/model"
)

type ConfigRepository struct {
	Pool *pgxpool.Pool
}

func NewConfigRepository(pool *pgxpool.Pool) *ConfigRepository {
	return &ConfigRepository{Pool: pool}
}

// GetByRepoID fetches settings for a specific repo
func (r *ConfigRepository) GetByRepoID(ctx context.Context, repoID int) (*model.Configuration, error) {
	query := `SELECT id, repository_id, ignore_patterns, review_style, enabled, created_at FROM configurations WHERE repository_id = $1`
	
	var config model.Configuration
	err := r.Pool.QueryRow(ctx, query, repoID).Scan(
		&config.ID, &config.RepositoryID, &config.IgnorePatterns, &config.ReviewStyle, &config.Enabled, &config.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			// If no config exists, return a default one
			return &model.Configuration{
				RepositoryID:   repoID,
				IgnorePatterns: []string{},
				ReviewStyle:    "concise",
				Enabled:        true,
			}, nil
		}
		return nil, fmt.Errorf("failed to get config: %w", err)
	}
	return &config, nil
}

// UpsertConfig updates settings or creates them if missing
func (r *ConfigRepository) UpsertConfig(ctx context.Context, config *model.Configuration) error {
	query := `
		INSERT INTO configurations (repository_id, ignore_patterns, review_style, enabled)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (repository_id) DO UPDATE 
		SET ignore_patterns = EXCLUDED.ignore_patterns,
		    review_style = EXCLUDED.review_style,
		    enabled = EXCLUDED.enabled
		RETURNING id, created_at
	`
	return r.Pool.QueryRow(ctx, query, 
		config.RepositoryID, config.IgnorePatterns, config.ReviewStyle, config.Enabled,
	).Scan(&config.ID, &config.CreatedAt)
}