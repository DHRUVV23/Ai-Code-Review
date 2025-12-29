package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/DHRUVV23/ai-code-review/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepoRepository struct {
	Pool *pgxpool.Pool
}

func NewRepoRepository(pool *pgxpool.Pool) *RepoRepository {
	return &RepoRepository{Pool: pool}
}

// CreateRepository saves a new repo AND creates a default config
func (r *RepoRepository) CreateRepository(ctx context.Context, userID int, name, owner string) (*model.Repository, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// 1. Insert Repository
	var repoID int
	var createdAt time.Time
	
	repoQuery := `
		INSERT INTO repositories (user_id, name, owner)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`
	
	err = tx.QueryRow(ctx, repoQuery, userID, name, owner).Scan(&repoID, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert repo: %w", err)
	}

	// 2. Create Default Configuration
	// Note: We insert strict strings here because the DB column is TEXT
	configQuery := `
		INSERT INTO configurations (repository_id, review_style, ignore_patterns)
		VALUES ($1, 'concise', '*.md,*.lock')`
	
	_, err = tx.Exec(ctx, configQuery, repoID)
	if err != nil {
		return nil, fmt.Errorf("failed to create default config: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &model.Repository{
		ID:        repoID,
		UserID:    userID,
		Name:      name,
		Owner:     owner,
		CreatedAt: createdAt,
	}, nil
}

// ListRepositories gets all repos for a specific user
func (r *RepoRepository) ListRepositories(ctx context.Context, userID int) ([]model.Repository, error) {
	rows, err := r.Pool.Query(ctx, "SELECT id, user_id, name, owner, created_at FROM repositories WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repos []model.Repository
	for rows.Next() {
		var repo model.Repository
		if err := rows.Scan(&repo.ID, &repo.UserID, &repo.Name, &repo.Owner, &repo.CreatedAt); err != nil {
			return nil, err
		}
		repos = append(repos, repo)
	}
	return repos, nil
}

// GetRepositoryByID fetches a single repo
func (r *RepoRepository) GetRepositoryByID(ctx context.Context, id int) (*model.Repository, error) {
	query := `SELECT id, user_id, name, owner, created_at FROM repositories WHERE id = $1`
	row := r.Pool.QueryRow(ctx, query, id)

	var repo model.Repository
	err := row.Scan(&repo.ID, &repo.UserID, &repo.Name, &repo.Owner, &repo.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &repo, nil
}