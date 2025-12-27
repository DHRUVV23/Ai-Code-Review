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
	gh := service.NewGitHubService()
	parser := service.NewDiffParser() // ✅ You already have this

	// 2. Create Pending Review in DB
	reviewID, err := reviewRepo.CreateReview(ctx, p.RepoID, p.PrNumber)
	if err != nil {
		log.Printf("❌ Failed to create review record: %v", err)
		return err
	}

	// 3. FETCH REAL CODE FROM GITHUB
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
		// Keep fallback for testing stability
		realDiff = "func fallback() { log.Println('GitHub fetch failed, using fallback') }"
	}

	// ======================================================
	// 👇 NEW SECTION: PARSE & FILTER THE DIFF 👇
	// ======================================================
	log.Println("🔍 Parsing and Filtering files...")
	fileChanges := parser.Parse(realDiff)

	var filteredDiffBuilder strings.Builder
	validFilesCount := 0

	for _, file := range fileChanges {
		if !file.IsSafe {
			log.Printf("⚠️ Skipping ignored file: %s", file.Path)
			continue
		}
		
		validFilesCount++
		// Reconstruct the clean diff for the AI
		filteredDiffBuilder.WriteString(fmt.Sprintf("--- FILE: %s (%s) ---\n", file.Path, file.Language))
		filteredDiffBuilder.WriteString(file.Content)
		filteredDiffBuilder.WriteString("\n\n")
	}

	finalDiff := filteredDiffBuilder.String()

	if validFilesCount == 0 {
		log.Println("❌ No relevant code files found to review.")
		// Optional: You could save a "No code changes found" message to the DB here
		return nil
	}
	
	log.Printf("✨ Sending %d valid files to AI...", validFilesCount)
	// ======================================================
	// 👆 END OF NEW SECTION 👆
	// ======================================================

	// 4. Send CLEAN Diff to AI (Changed 'realDiff' to 'finalDiff')
	reviewContent, err := ai.ReviewCode(ctx, finalDiff, "concise")
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

	

	// ======================================================
    // 👇 NEW STEP 6: POST TO GITHUB 👇
    // ======================================================

	log.Println("💬 Posting comment to GitHub...")

    // A. Parse the JSON so we can format it nicely
    type ReviewIssue struct {
        File       string `json:"file"`
        Line       int    `json:"line"`
        Type       string `json:"type"`
        Severity   string `json:"severity"`
        Message    string `json:"message"`
        Suggestion string `json:"suggestion"`
    }

    var issues []ReviewIssue
    if err := json.Unmarshal([]byte(reviewContent), &issues); err != nil {
        log.Printf("⚠️ Could not parse JSON for GitHub comment: %v", err)
        // We don't return error here, because the review itself was successful
    } else if len(issues) > 0 {
        
        // B. Build a Markdown Table
        var commentBuilder strings.Builder
        commentBuilder.WriteString("## 🤖 AI Code Review Summary\n\n")
        commentBuilder.WriteString("I found the following issues in your changes:\n\n")
        commentBuilder.WriteString("| Severity | File | Issue | Suggestion |\n")
        commentBuilder.WriteString("|---|---|---|---|\n")

        for _, issue := range issues {
            // Pick an icon based on severity
            icon := "ℹ️"
            if strings.ToLower(issue.Severity) == "high" { icon = "🔴" }
            if strings.ToLower(issue.Severity) == "medium" { icon = "🟡" }
            if strings.ToLower(issue.Severity) == "low" { icon = "🟢" }

            // Format the row
            row := fmt.Sprintf("| %s **%s** | `%s` (Line %d) | %s | `%s` |\n",
                icon, issue.Severity, issue.File, issue.Line, issue.Message, issue.Suggestion)
            commentBuilder.WriteString(row)
        }
        
        commentBuilder.WriteString("\n*Generated by AI Code Reviewer* 🚀")

        // C. Send it!
        err = gh.PostComment(ctx, owner, repoName, p.PrNumber, commentBuilder.String())
        if err != nil {
            log.Printf("❌ Failed to post comment to GitHub: %v", err)
        } else {
            log.Printf("✅ Comment posted to GitHub successfully!")
        }
    } else {
        log.Println("✅ No issues found, skipping GitHub comment.")
    }
    // ======================================================

    log.Printf("✅ Review Saved! ID: %d", reviewID)
    return nil
}
