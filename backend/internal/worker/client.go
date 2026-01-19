package worker

import (
	"log"
	"os"

	"github.com/hibiken/asynq"
)

// Global variable to hold the connection
var Client *asynq.Client

// InitClient connects to Redis 
func InitClient() {
	redisAddr := os.Getenv("REDIS_ADDR") 
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	Client = asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	log.Println("Job Queue Client Connected")
}

// CloseClient
func CloseClient() {
	if Client != nil {
		Client.Close()
	}
}