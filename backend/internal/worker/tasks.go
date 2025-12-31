package worker

import (
	"encoding/json"
	"github.com/hibiken/asynq"
)

// Task Name
const TypeReviewPR = "review:pr"

// Payload
type ReviewPayload struct {
	RepoName  string `json:"repo_name"`
	RepoOwner string `json:"repo_owner"`
	PRNumber  int    `json:"pr_number"`
	RepoID    int64  `json:"repo_id"`
}

// NewReviewTask creates the task (Use this name!)
func NewReviewTask(repoName, repoOwner string, prNumber int, repoID int64) (*asynq.Task, error) {
	payload, err := json.Marshal(ReviewPayload{
		RepoName:  repoName,
		RepoOwner: repoOwner,
		PRNumber:  prNumber,
		RepoID:    repoID,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeReviewPR, payload), nil
}