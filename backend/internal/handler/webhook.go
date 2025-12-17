package handler

import (
	"log"
	"net/http"
	"os"
	"io"

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
	switch e := parsedEvent.(type) {
	case *github.PullRequestEvent:
		log.Printf("🚀 PR Event: %s", e.GetAction())
	case *github.PingEvent:
		log.Println("🏓 Pong!")
	}

	c.JSON(http.StatusOK, gin.H{"message": "Received"})
}