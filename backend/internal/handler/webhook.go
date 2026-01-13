package handler

import (
	"log"
	"net/http"
	"os"

	"github.com/DHRUVV23/ai-code-review/backend/internal/worker" 
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v50/github" 
	"github.com/hibiken/asynq"
)

type WebhookHandler struct {
	Client *asynq.Client
}

func (h *WebhookHandler) HandleWebhook(c *gin.Context) {
	// 1. Validate Payload (Security Check)
	webhookSecret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	payload, err := github.ValidatePayload(c.Request, []byte(webhookSecret))
	if err != nil {
		log.Printf("❌ Invalid signature: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	// 2. Parse Event
	event, err := github.ParseWebHook(github.WebHookType(c.Request), payload)
	if err != nil {
		log.Printf("❌ Could not parse webhook: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse webhook"})
		return
	}

	// 3. Handle Pull Request Events
	switch e := event.(type) {
	case *github.PullRequestEvent:
		action := e.GetAction()

		// We only care when PR is OPENED or UPDATED (Synchronize)
		if action == "opened" || action == "synchronize" {
			repo := e.GetRepo()
			prNumber := e.GetNumber()
			
			
			repoName := repo.GetName()
			repoOwner := repo.GetOwner().GetLogin()
			repoID := repo.GetID() 

			log.Printf(" Processing PR #%d for %s/%s", prNumber, repoOwner, repoName)

			// 4. Create Task for Worker
			task, err := worker.NewReviewTask(repoName, repoOwner, prNumber, repoID)
			if err != nil {
				log.Printf(" Failed to create task: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error"})
				return
			}

			// 5. Enqueue Task to Redis
			info, err := h.Client.Enqueue(task)
			if err != nil {
				log.Printf(" Failed to enqueue task: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue job"})
				return
			}

			log.Printf(" Review Job Enqueued! ID: %s", info.ID)
		}

	case *github.PingEvent:
		log.Println("🏓 GitHub Ping! Connection verified.")

	default:
		// Ignore other events (starring, forking, etc.)
		// log.Printf("Ignored event: %T", e)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event processed"})
}