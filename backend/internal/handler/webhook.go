package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/DHRUVV23/ai-code-review/backend/internal/worker"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v50/github"
	"github.com/hibiken/asynq"
)

type WebhookHandler struct {
	Client *asynq.Client
}

func (h *WebhookHandler) HandleWebhook(c *gin.Context) {
	
	webhookSecret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	payload, err := github.ValidatePayload(c.Request, []byte(webhookSecret))
	if err != nil {
		log.Printf("Invalid signature: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}


	event, err := github.ParseWebHook(github.WebHookType(c.Request), payload)
	if err != nil {
		log.Printf("Could not parse webhook: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse webhook"})
		return
	}

	switch e := event.(type) {
	case *github.PullRequestEvent:
		action := e.GetAction()

	
		if action != "opened" && action != "synchronize" {
			log.Printf("Ignoring PR action: %s (we only care about 'opened' or 'synchronize')", action)
			c.JSON(http.StatusOK, gin.H{"status": "ignored"})
			return
		}

		repo := e.GetRepo()
		prNumber := e.GetNumber()
		repoName := repo.GetName()
		repoOwner := repo.GetOwner().GetLogin()
		repoID := int(repo.GetID())

	
		commitSHA := e.GetPullRequest().GetHead().GetSHA()

		log.Printf(" Processing PR #%d for %s/%s (Commit: %s)", prNumber, repoOwner, repoName, commitSHA)

	
		task, err := worker.NewReviewTask(repoName, repoOwner, prNumber, int64(repoID))
		if err != nil {
			log.Printf("Failed to create task: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error"})
			return
		}


		taskID := fmt.Sprintf("review:%s/%s:%d:%s", repoOwner, repoName, prNumber, commitSHA)

		info, err := h.Client.Enqueue(task,
			asynq.TaskID(taskID),          
			asynq.Retention(1*time.Hour),   
		)

		if err != nil {
		
			if strings.Contains(err.Error(), "task ID conflicts") {
				log.Printf(" Duplicate Review Task Ignored: %s", taskID)
				c.JSON(http.StatusOK, gin.H{"status": "duplicate_ignored"})
				return
			}

			
			log.Printf(" Failed to enqueue task: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue job"})
			return
		}

		log.Printf(" Review Job Enqueued! ID: %s", info.ID)

	case *github.PingEvent:
		log.Println(" GitHub Ping! Connection verified.")

	default:
		
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event processed"})
}