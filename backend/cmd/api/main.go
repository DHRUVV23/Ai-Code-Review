package main

import (
	"log"
	
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv" 
    
    // MAKE SURE THIS PATH MATCHES YOUR go.mod FILE
	"github.com/DHRUVV23/ai-code-review/backend/internal/database"
)

func main() {
	// 1. Load environment variables
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Warning: .env file not found")
	}

	// 2. Initialize Database
	if err := database.InitDB(); err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer database.CloseDB()

	// 3. Start Server
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"status": "database connected",
		})
	})

	r.Run(":8080")
}