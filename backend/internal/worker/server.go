package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
    
	"github.com/DHRUVV23/ai-code-review/backend/internal/database"   // <--- Import DB
	"github.com/DHRUVV23/ai-code-review/backend/internal/repository"
	"github.com/DHRUVV23/ai-code-review/backend/internal/service"
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

// // HandleReviewPR is the function that actually DOES the work
// func HandleReviewPR(ctx context.Context, t *asynq.Task) error {
// 	var p ReviewPRPayload
// 	if err := json.Unmarshal(t.Payload(), &p); err != nil {
// 		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
// 	}

// 	log.Printf("🤖 AI WORKING: Reviewing Repo %s (PR #%d)", p.RepoName, p.PrNumber)
	
// 	// TODO: Phase 3 - Call OpenAI here!
// 	// For now, we just pretend to work
// 	return nil
// }

func HandleReviewPR(ctx context.Context, t *asynq.Task) error {
	var p ReviewPRPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("🤖 AI WORKING: Reviewing Repo %s (PR #%d)", p.RepoName, p.PrNumber)

	// 1. Initialize DB Connection & Service
	reviewRepo := repository.NewReviewRepository(database.Pool)
	ai := service.NewAIService()

	// 2. Create a "Pending" record in the database
	// This lets the frontend show "Review in Progress..."
	reviewID, err := reviewRepo.CreateReview(ctx, p.RepoID, p.PrNumber)
	if err != nil {
		log.Printf("❌ Failed to create review record: %v", err)
		return err
	}

	// 3. Mock Diff (Later we fetch from GitHub)
	mockDiff := `
	func login(password string) {
		log.Printf("User password: %s", password) // Security Issue!
	}
	`

	// 4. Get AI Analysis
	reviewContent, err := ai.ReviewCode(ctx, mockDiff, "concise")
	if err != nil {
		log.Printf("❌ AI Failed: %v", err)
		return err 
	}

	// 5. Update the Database with the Result
	err = reviewRepo.UpdateReviewResult(ctx, reviewID, reviewContent)
	if err != nil {
		log.Printf("❌ Failed to save review result: %v", err)
		return err
	}

	log.Printf("✅ Review Saved! ID: %d", reviewID)
	return nil
}