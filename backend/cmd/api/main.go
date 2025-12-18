package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	// Make sure these paths match your actual module name
	"github.com/DHRUVV23/ai-code-review/backend/internal/database"
	"github.com/DHRUVV23/ai-code-review/backend/internal/handler"
	"github.com/DHRUVV23/ai-code-review/backend/internal/repository"
	"github.com/DHRUVV23/ai-code-review/backend/internal/worker"
)

func main() {
	// 1. Load Environment Variables
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Warning: .env file not found")
	}

	// 2. Connect to Database
	if err := database.InitDB(); err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer database.CloseDB()

	worker.InitClient()
    defer worker.CloseClient()

	worker.StartWorker()

	// 3. Initialize Repositories (The "Stock Clerks")
	// These handle the raw SQL for each table
	userRepo := repository.NewUserRepository(database.Pool)
	repoRepo := repository.NewRepoRepository(database.Pool)     // <--- ADDED
	configRepo := repository.NewConfigRepository(database.Pool) // <--- ADDED

	// 4. Initialize Handlers (The "Waiters")
	// We inject the repositories into the handlers
	authHandler := &handler.AuthHandler{UserRepo: userRepo}
	repoHandler := &handler.RepoHandler{
		RepoRepository:   repoRepo,   // <--- ADDED
		ConfigRepository: configRepo, // <--- ADDED
	}

	

	// 5. Setup Router
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// 6. Register Routes
	// --- Webhooks ---
	r.POST("/webhooks/github", handler.HandleWebhook)

	// --- Authentication ---
	r.GET("/auth/login", authHandler.GitHubLogin)
	r.GET("/auth/callback", authHandler.GitHubCallback)

	// --- Repositories & Config (API) ---
	r.GET("/api/repos", repoHandler.ListRepositories)            // List all repos
	r.GET("/api/repos/:id/config", repoHandler.GetConfig)        // Get config for a repo
	r.POST("/api/repos/:id/config", repoHandler.UpdateConfig)    // Update config for a repo

	// 7. Start Server
	r.Run(":8080")
}