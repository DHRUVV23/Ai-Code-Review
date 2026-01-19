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
		fmt.Printf("Failed to create AI Client: %v\n", err)
		return &AIService{Client: nil}
	}

	return &AIService{Client: client}
}

// ReviewCode sends the diff to Gemini and gets feedback
func (s *AIService) ReviewCode(ctx context.Context, diff string, style string) (string, error) {
	if s.Client == nil {
		return "AI Client not initialized. Check GEMINI_API_KEY.", nil
	}
	defer s.Client.Close()

	model := s.Client.GenerativeModel("gemini-flash-latest")
	model.ResponseMIMEType = "application/json"

	//  PROMPT TEMPLATE
	prompt := fmt.Sprintf(`
	You are a Senior Code Reviewer. 
	Analyze the following Git Diff code changes.
	
	OBJECTIVE:
	Identify bugs, security vulnerabilities, performance issues, and bad practices.
	
	STRICT OUTPUT FORMAT:
	You must respond ONLY with a valid JSON array. Do not use markdown formatting.
	Use this schema:
	[
		{
			"file": "filename.ext",
			"line": 10,
			"type": "security|bug|performance|style",
			"severity": "high|medium|low",
			"message": "Concise explanation of the issue",
			"suggestion": "Code or logic to fix it"
		}
	]

	If the code is perfectly fine, return an empty array: []

	CODE CONTEXT (DIFF):
	%s
	`, diff)

	//  Send Request
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	//  Extract Response
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		part := resp.Candidates[0].Content.Parts[0]
		if txt, ok := part.(genai.Text); ok {
			return string(txt), nil
		}
	}

	return "[]", nil 
}