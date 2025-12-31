package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	
	// Import your services
	"github.com/DHRUVV23/ai-code-review/backend/internal/service"
)

// HandleReviewTask is the logic that runs when a job is picked up
func HandleReviewTask(ctx context.Context, t *asynq.Task) error {
	// 1. Unmarshal the Payload (Defined in tasks.go)
	var payload ReviewPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("🤖 Processing Review for: %s/%s PR #%d", payload.RepoOwner, payload.RepoName, payload.PRNumber)

	// 2. Initialize Services
	// Note: These use os.Getenv inside them, as you defined
	ghService := service.NewGitHubService()
	aiService := service.NewAIService()

	// 3. Get the Code Diff from GitHub
	diff, err := ghService.GetPullRequestDiff(ctx, payload.RepoOwner, payload.RepoName, payload.PRNumber)
	if err != nil {
		log.Printf("❌ Failed to get diff: %v", err)
		return err // Retry later
	}

	if diff == "" {
		log.Println("⚠️ Diff is empty, skipping review.")
		return nil
	}

	// 4. Send to Gemini AI
	reviewJSON, err := aiService.ReviewCode(ctx, diff, "concise")
	if err != nil {
		log.Printf("❌ AI Analysis failed: %v", err)
		return err // Retry later
	}

	// 5. Post the Comment to GitHub
	// We wrap the JSON in a nice Markdown block
	commentBody := fmt.Sprintf("## 🤖 AI Code Review\n\n```json\n%s\n```", reviewJSON)
	
	if err := ghService.PostComment(ctx, payload.RepoOwner, payload.RepoName, payload.PRNumber, commentBody); err != nil {
		log.Printf("❌ Failed to post comment: %v", err)
		return err
	}

	log.Printf("✅ Review Posted for PR #%d!", payload.PRNumber)
	return nil
}

// StartWorker initializes the background processor
func StartWorker(redisAddr string) {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 10, // Process 10 reviews at once
		},
	)

	mux := asynq.NewServeMux()
	// Map the "review:pr" string to our handler function
	mux.HandleFunc(TypeReviewPR, HandleReviewTask)

	// Run in a Goroutine so it doesn't block the API Server
	go func() {
		log.Println("👷 Worker Server Started...")
		if err := srv.Run(mux); err != nil {
			log.Fatalf("❌ Worker failed to start: %v", err)
		}
	}()
}