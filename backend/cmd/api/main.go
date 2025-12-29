package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/gin-contrib/cors"

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
	reviewRepo := repository.NewReviewRepository(database.Pool)

	// 4. Initialize Handlers (The "Waiters")
	// We inject the repositories into the handlers
	authHandler := &handler.AuthHandler{UserRepo: userRepo}
	repoHandler := &handler.RepoHandler{
		RepoRepository:   repoRepo,   // <--- ADDED
		ConfigRepository: configRepo, // <--- ADDED
	}
	userHandler := &handler.UserHandler{UserRepo: userRepo}
	reviewHandler := &handler.ReviewHandler{ReviewRepository: reviewRepo}

	

	// 5. Setup Router
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} 
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// 6. Register Routes
	// --- Webhooks ---
	r.POST("/webhooks/github", handler.HandleWebhook)
	r.GET("/auth/login", authHandler.GitHubLogin)
	r.GET("/auth/callback", authHandler.GitHubCallback)

	// 2. Dashboard API (Week 3 Structure)
	// We group them under "v1" so your URLs look like /api/v1/user/profile
	v1 := r.Group("/api/v1")
	{
		// User Routes
		v1.GET("/user/profile", userHandler.GetUserProfile)
		v1.PUT("/user/profile", userHandler.UpdateUserProfile)
		
		// Repository Routes (Refactoring old routes to match new plan)
		v1.GET("/user/repositories", repoHandler.ListRepositories) // Matches "GET /api/v1/user/repositories"
		v1.GET("/repositories/:id", repoHandler.GetConfig)         // Matches "GET /api/v1/repositories/:id" (We reused config for now)
		v1.PUT("/repositories/:id/config", repoHandler.UpdateConfig)
		v1.GET("/repositories/:id/reviews", reviewHandler.ListReviews) 
	}

	r.Run(":8080")
}