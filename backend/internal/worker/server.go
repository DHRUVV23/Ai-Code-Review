package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

// StartWorker initializes the background processor
func StartWorker() {
	redisAddr := "localhost:6379"
	
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 10, // Process 10 jobs at once
		},
	)

	mux := asynq.NewServeMux()
	
	// Register the Handler: "When you see 'review:pr', run 'HandleReviewPR'"
	mux.HandleFunc(TypeReviewPR, HandleReviewPR)

	// Run in a separate "goroutine" (background thread)
	go func() {
		log.Println("👷 Worker Server Started... Waiting for jobs!")
		if err := srv.Run(mux); err != nil {
			log.Fatalf("❌ Worker failed to start: %v", err)
		}
	}()
}

// HandleReviewPR is the function that actually DOES the work
func HandleReviewPR(ctx context.Context, t *asynq.Task) error {
	var p ReviewPRPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("🤖 AI WORKING: Reviewing Repo %s (PR #%d)", p.RepoName, p.PrNumber)
	
	// TODO: Phase 3 - Call OpenAI here!
	// For now, we just pretend to work
	return nil
}