package handler

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/DHRUVV23/ai-code-review/backend/internal/worker"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v57/github"
)

func HandleWebhook(c *gin.Context) {
	// 1. Read Body
	_, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// 2. Validate Signature (USES webhookSecret)
	webhookSecret := os.Getenv("GITHUB_WEBHOOK_SECRET")

	// We pass 'webhookSecret' here, so now it is USED!
	event, err := github.ValidatePayload(c.Request, []byte(webhookSecret))
	if err != nil {
		log.Printf("❌ Invalid signature: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	// 3. Parse Event (USES event)
	webhookType := github.WebHookType(c.Request)
	// We pass 'event' (the raw payload) here
	parsedEvent, err := github.ParseWebHook(webhookType, event)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse webhook"})
		return
	}

	// 4. Log it
	// ... inside HandleWebhook ...

	switch e := parsedEvent.(type) {
	case *github.PullRequestEvent:
		action := e.GetAction()

		if action == "opened" || action == "synchronize" {
			repoName := e.GetRepo().GetFullName()
			prNumber := e.GetNumber()
			// We use a dummy ID '1' for now, later we lookup the real ID from DB
			repoID := 1

			log.Printf("🔔 Enqueueing Review for %s #%d", repoName, prNumber)

			// 1. Create the Task
			task, err := worker.NewReviewPRTask(repoID, prNumber, repoName)
			if err != nil {
				log.Printf("❌ Failed to create task: %v", err)
				return
			}

			// 2. Send it to the Queue!
			info, err := worker.Client.Enqueue(task)
			if err != nil {
				log.Printf("❌ Failed to enqueue task: %v", err)
				return
			}

			log.Printf("✅ Job Enqueued! ID: %s", info.ID)

		} else {
			log.Printf("💤 Ignoring PR action: %s", action)
		}

	case *github.PingEvent:
		log.Println("🏓 GitHub Ping! Webhook connection works.")

	default:
		log.Printf("Ignored event type: %s", webhookType)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event received"})
}
