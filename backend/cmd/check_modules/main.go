package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func main() {
	// 1. Load the .env file to get your KEY
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("⚠️ Warning: Could not load .env file (checking system env vars)")
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("❌ Error: GEMINI_API_KEY is empty.")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	fmt.Println("🔍 Checking available models for your API Key...")
	iter := client.ListModels(ctx)
	found := false
	for {
		m, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// Print only models that support content generation
		fmt.Printf("✅ Available: %s\n", m.Name)
		found = true
	}

	if !found {
		fmt.Println("❌ No models found! Your API Key might be invalid or has no permissions.")
	}
}