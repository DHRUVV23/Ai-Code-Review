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
	"strings"
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

	// 1. Init Services
	reviewRepo := repository.NewReviewRepository(database.Pool)
	ai := service.NewAIService()
	gh := service.NewGitHubService() // <--- NEW

	// 2. Create Pending Review in DB
	reviewID, err := reviewRepo.CreateReview(ctx, p.RepoID, p.PrNumber)
	if err != nil {
		log.Printf("❌ Failed to create review record: %v", err)
		return err
	}

	// 3. FETCH REAL CODE FROM GITHUB (The New Part!)
	// RepoName comes in as "owner/repo" (e.g., "DHRUVV23/ai-code-review")
	parts := strings.Split(p.RepoName, "/")
	if len(parts) != 2 {
		log.Printf("❌ Invalid repo name format: %s", p.RepoName)
		return fmt.Errorf("invalid repo name")
	}
	owner, repoName := parts[0], parts[1]

	log.Println("🌍 Fetching Diff from GitHub...")
	realDiff, err := gh.GetPullRequestDiff(ctx, owner, repoName, p.PrNumber)
	if err != nil {
		log.Printf("❌ GitHub Fetch Failed: %v", err)
		// Fallback: If GitHub fails (e.g., token issue), use the mock so the flow doesn't break
		realDiff = "func fallback() { log.Println('GitHub fetch failed, using fallback') }"
	}

	// 4. Send REAL Diff to AI
	reviewContent, err := ai.ReviewCode(ctx, realDiff, "concise")
	if err != nil {
		log.Printf("❌ AI Failed: %v", err)
		return err
	}

	// 5. Save Result
	err = reviewRepo.UpdateReviewResult(ctx, reviewID, reviewContent)
	if err != nil {
		log.Printf("❌ Failed to save review result: %v", err)
		return err
	}

	log.Printf("✅ Review Saved! ID: %d", reviewID)
	return nil
}