package service

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type GitHubService struct {
	Client *github.Client
}

func NewGitHubService() *GitHubService {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		// Return client without auth if token is missing 
		return &GitHubService{Client: github.NewClient(nil)}
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &GitHubService{Client: github.NewClient(tc)}
}

// GetPullRequestDiff fetches the raw text of the changes
func (s *GitHubService) GetPullRequestDiff(ctx context.Context, owner, repo string, prNumber int) (string, error) {
	opts := github.RawOptions{Type: github.Diff}
	diff, _, err := s.Client.PullRequests.GetRaw(ctx, owner, repo, prNumber, opts)
	if err != nil {
		return "", fmt.Errorf("failed to fetch PR diff: %w", err)
	}
	return diff, nil
}

// PostComment posts a markdown comment to the PR
func (s *GitHubService) PostComment(ctx context.Context, owner, repo string, prNumber int, commentBody string) error {
	
	comment := &github.IssueComment{
		Body: &commentBody,
	}

	// Use the library's built-in method 
	_, _, err := s.Client.Issues.CreateComment(ctx, owner, repo, prNumber, comment)
	if err != nil {
		return fmt.Errorf("failed to post comment: %w", err)
	}

	return nil
}