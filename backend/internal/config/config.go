package config

import (
	// "log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL        string
	RedisAddr          string
	GithubClientID     string
	GithubClientSecret string
	WebhookSecret      string
}

func LoadConfig() (*Config, error) {
	// Try loading .env file (it might fail in Docker/Prod, which is fine if env vars are set)
	// We look for .env in the root (backend/) or parent directory
	_ = godotenv.Load() 
	// Also try loading from one level up just in case
	_ = godotenv.Load("../.env")

	cfg := &Config{
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/code_review_db"),
		RedisAddr:          getEnv("REDIS_ADDR", "localhost:6379"),
		GithubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		GithubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		WebhookSecret:      os.Getenv("GITHUB_WEBHOOK_SECRET"),
	}

	return cfg, nil
}

// Helper to read env with a fallback default
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}