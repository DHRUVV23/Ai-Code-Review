package service

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type AIService struct {
	Client *genai.Client
}

func NewAIService() *AIService {
	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	
	// Create the client with the API Key
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		// If client fails, we log it but don't crash yet (will fail on usage)
		fmt.Printf("❌ Failed to create AI Client: %v\n", err)
		return &AIService{Client: nil}
	}

	return &AIService{Client: client}
}

// ReviewCode sends the diff to Gemini and gets feedback
func (s *AIService) ReviewCode(ctx context.Context, diff string, style string) (string, error) {
	if s.Client == nil {
		return "❌ AI Client not initialized. Check GEMINI_API_KEY.", nil
	}
	defer s.Client.Close()

	// 1. Select the Model (Gemini 1.5 Flash is fast & free)
	model := s.Client.GenerativeModel("gemini-flash-latest")

	// 2. Construct the Prompt
	prompt := fmt.Sprintf(`
	You are an expert Senior Software Engineer doing a code review.
	Review Style: %s (Be strictly professional and concise).
	
	Analyze the following Git Diff. 
	- Identify potential bugs, security flaws, or performance issues.
	- Suggest cleaner or more idiomatic Go code if applicable.
	- If the code looks good, just say "LGTM".
	
	Git Diff:
	%s
	`, style, diff)

	// 3. Send Request to AI
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("AI generation failed: %w", err)
	}

	// 4. Extract Response
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		// The response comes back as "Parts", usually Text
		if txt, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			return string(txt), nil
		}
	}

	return "⚠️ AI returned no content.", nil
}