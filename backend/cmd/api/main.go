package main

import (
	"log"
	// "os" // Add this

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	
	"github.com/DHRUVV23/ai-code-review/backend/internal/database"
	"github.com/DHRUVV23/ai-code-review/backend/internal/handler"
	"github.com/DHRUVV23/ai-code-review/backend/internal/repository" // Add this
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Warning: .env file not found")
	}

	if err := database.InitDB(); err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer database.CloseDB()

	// 1. Initialize Repositories
	userRepo := repository.NewUserRepository(database.Pool)

	// 2. Initialize Handlers
	authHandler := &handler.AuthHandler{UserRepo: userRepo}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// 3. Register Routes
	// Webhooks (Existing)
	r.POST("/webhooks/github", handler.HandleWebhook)
	
	// Auth (New!)
	r.GET("/auth/login", authHandler.GitHubLogin)
	r.GET("/auth/callback", authHandler.GitHubCallback)

	r.Run(":8080")
}