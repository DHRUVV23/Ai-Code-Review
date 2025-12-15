package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/DHRUVV23/ai-code-review/backend/internal/model"
)

type RepoRepository struct {
	Pool *pgxpool.Pool
}

// NewRepoRepository creates a new instance
func NewRepoRepository(pool *pgxpool.Pool) *RepoRepository {
	return &RepoRepository{Pool: pool}
}

// Create (The 'C' in CRUD) - Saves a new repo to DB
func (r *RepoRepository) CreateRepository(ctx context.Context, repo *model.Repository) error {
	query := `
		INSERT INTO repositories (github_repo_id, installation_id, name, full_name, private)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	err := r.Pool.QueryRow(ctx, query,
		repo.GithubRepoID,
		repo.InstallationID,
		repo.Name,
		repo.FullName,
		repo.Private,
	).Scan(&repo.ID, &repo.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert repository: %w", err)
	}
	return nil
}

// GetByID (The 'R' in CRUD) - Fetches a repo by ID
func (r *RepoRepository) GetRepositoryByID(ctx context.Context, id int) (*model.Repository, error) {
	query := `SELECT id, github_repo_id, installation_id, name, full_name, private, created_at FROM repositories WHERE id = $1`
	
	row := r.Pool.QueryRow(ctx, query, id)
	
	var repo model.Repository
	err := row.Scan(
		&repo.ID, 
		&repo.GithubRepoID, 
		&repo.InstallationID, 
		&repo.Name, 
		&repo.FullName, 
		&repo.Private, 
		&repo.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}
	return &repo, nil
}