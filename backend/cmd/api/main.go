package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"

	// ✅ IMPORTS WITH BACKEND PREFIX INCLUDED
	"github.com/DHRUVV23/ai-code-review/backend/internal/config"
	"github.com/DHRUVV23/ai-code-review/backend/internal/database"
	"github.com/DHRUVV23/ai-code-review/backend/internal/handler"
	"github.com/DHRUVV23/ai-code-review/backend/internal/repository"
	"github.com/DHRUVV23/ai-code-review/backend/internal/worker"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// 2. Connect to Database (Using postgres.go InitDB)
	if err := database.InitDB(); err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer database.CloseDB()




	// 3. Initialize Redis Client (Asynq)
	redisOpt := asynq.RedisClientOpt{Addr: cfg.RedisAddr}
	asynqClient := asynq.NewClient(redisOpt)
	defer asynqClient.Close()

	worker.StartWorker(cfg.RedisAddr)

	// 4. Initialize Repositories
	// We use database.Pool which was set by InitDB()
	userRepo := repository.NewUserRepository(database.Pool)
	repoRepo := repository.NewRepoRepository(database.Pool)
	configRepo := repository.NewConfigRepository(database.Pool)
	// reviewRepo := repository.NewReviewRepository(database.Pool) // Uncomment when created

	// 5. Initialize Handlers
	authHandler := &handler.AuthHandler{
		UserRepo: userRepo,
		Config:   cfg,
	}

	repoHandler := &handler.RepoHandler{
		RepoRepository:   repoRepo,
		ConfigRepository: configRepo,
		UserRepository:   userRepo,
	}

	webhookHandler := &handler.WebhookHandler{
		Client: asynqClient,
	}

	// 6. Setup Router
	r := gin.Default()

	// CORS Configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

	// 7. Register Routes
	
	// Public Routes
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.POST("/webhook", webhookHandler.HandleWebhook) // GitHub calls this
	r.GET("/auth/github/login", authHandler.GitHubLogin)
	r.GET("/auth/github/callback", authHandler.GitHubCallback)

	// Protected Routes (Dashboard)
	v1 := r.Group("/api/v1")
	{
		v1.GET("/user/profile", authHandler.GetUserProfile)
		v1.PUT("/user/profile", authHandler.UpdateUserProfile)

		v1.POST("/repositories", repoHandler.RegisterRepository)
		v1.GET("/user/repositories", repoHandler.ListRepositories)
		v1.GET("/repositories/:id", repoHandler.GetConfig)
		v1.PUT("/repositories/:id/config", repoHandler.UpdateConfig)
		
		// Connect to GitHub Button Endpoint
		v1.POST("/repositories/:id/webhook", repoHandler.CreateWebhook)
	}

	log.Println("🚀 Server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}