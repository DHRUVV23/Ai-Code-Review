package worker

import (
	"encoding/json"
	// "fmt"

	"github.com/hibiken/asynq"
)

// 1. Name the Job Type (Like a Subject Line)
const TypeReviewPR = "review:pr"

// 2. Define the Data Payload (What's inside the envelope)
type ReviewPRPayload struct {
	RepoID   int    `json:"repo_id"`
	PrNumber int    `json:"pr_number"`
	RepoName string `json:"repo_name"` // e.g., "owner/repo"
}

// 3. Create the Task Helper Function
func NewReviewPRTask(repoID int, prNumber int, repoName string) (*asynq.Task, error) {
	payload, err := json.Marshal(ReviewPRPayload{
		RepoID:   repoID,
		PrNumber: prNumber,
		RepoName: repoName,
	})
	if err != nil {
		return nil, err
	}
	
	// Create the task with the payload
	return asynq.NewTask(TypeReviewPR, payload), nil
}