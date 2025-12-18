package worker

import (
	"log"
	"os"

	"github.com/hibiken/asynq"
)

// Global variable to hold the connection
var Client *asynq.Client

// InitClient connects to Redis so we can send jobs
func InitClient() {
	redisAddr := os.Getenv("REDIS_ADDR") // e.g., "localhost:6379"
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	Client = asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	log.Println("✅ Job Queue Client Connected")
}

// CloseClient cleans up the connection
func CloseClient() {
	if Client != nil {
		Client.Close()
	}
}