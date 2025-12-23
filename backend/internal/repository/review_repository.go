package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/DHRUVV23/ai-code-review/backend/internal/model"
)

type ReviewRepository struct {
	Pool *pgxpool.Pool
}

func NewReviewRepository(pool *pgxpool.Pool) *ReviewRepository {
	return &ReviewRepository{Pool: pool}
}

// CreateReview starts a new review entry (status: pending)
func (r *ReviewRepository) CreateReview(ctx context.Context, repoID int, prNumber int) (int, error) {
	var id int
	query := `INSERT INTO reviews (repository_id, pr_number, status, commit_sha, created_at) 
	          VALUES ($1, $2, 'pending', 'dummy-sha-123', NOW()) RETURNING id`
	
	err := r.Pool.QueryRow(ctx, query, repoID, prNumber).Scan(&id)
	return id, err
}

// UpdateReview saves the AI response and marks it as completed
func (r *ReviewRepository) UpdateReviewResult(ctx context.Context, id int, content string) error {
	query := `UPDATE reviews SET content = $1, status = 'completed' WHERE id = $2`
	_, err := r.Pool.Exec(ctx, query, content, id)
	return err
}

// GetReviewsByRepoID fetches all reviews for a specific project
func (r *ReviewRepository) GetReviewsByRepoID(ctx context.Context, repoID int) ([]model.Review, error) {
	query := `SELECT id, repository_id, pr_number, status, COALESCE(content, ''), created_at FROM reviews WHERE repository_id = $1 ORDER BY created_at DESC`
	
	rows, err := r.Pool.Query(ctx, query, repoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []model.Review
	for rows.Next() {
		var rev model.Review
		if err := rows.Scan(&rev.ID, &rev.RepositoryID, &rev.PRNumber, &rev.Status, &rev.Content, &rev.CreatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, rev)
	}
	return reviews, nil
}