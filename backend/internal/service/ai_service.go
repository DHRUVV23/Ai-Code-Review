package service

import (
	"context"
	"log"
	"time"
)

type AIService struct{}

func NewAIService() *AIService {
	return &AIService{}
}

// ReviewCode simulates an AI analyzing the code
func (s *AIService) ReviewCode(ctx context.Context, diff string, style string) (string, error) {
	log.Println("🧠 AI Service: Analyzing code diff...")
	
	// Simulate "Thinking" time (2 seconds)
	time.Sleep(2 * time.Second)

	// Return a fake review result for now
	// Later, we will replace this with the real Google Gemini / OpenAI code!
	return "✅ AI Review Result: The code looks clean, but check for hardcoded secrets.", nil
}