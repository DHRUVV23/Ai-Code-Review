package service

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

type GitHubService struct {
	Client *github.Client
}

func NewGitHubService() *GitHubService {
	// 1. Get the Token from .env (for now, we use a personal token)
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Println("⚠️ Warning: GITHUB_TOKEN is missing. Private repos will fail.")
		return &GitHubService{Client: github.NewClient(nil)}
	}

	// 2. Authenticate
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &GitHubService{Client: client}
}

// GetPullRequestDiff fetches the raw code changes from a PR
func (s *GitHubService) GetPullRequestDiff(ctx context.Context, owner, repo string, prNumber int) (string, error) {
	// 3. Request the "Diff" format from GitHub
	// We want the raw diff string, not the JSON object
	diff, _, err := s.Client.PullRequests.GetRaw(ctx, owner, repo, prNumber, github.RawOptions{Type: github.Diff})
	if err != nil {
		return "", fmt.Errorf("failed to fetch PR diff: %w", err)
	}

	return diff, nil
}