package service

import (
	"context"
	"fmt"
	"os"
	"strings" 

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type GitHubService struct {
	Client *github.Client
}

func NewGitHubService() *GitHubService {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return &GitHubService{Client: github.NewClient(nil)}
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &GitHubService{Client: github.NewClient(tc)}
}

func (s *GitHubService) GetPullRequestDiff(ctx context.Context, owner, repo string, prNumber int) (string, error) {
	opts := github.RawOptions{Type: github.Diff}
	diff, _, err := s.Client.PullRequests.GetRaw(ctx, owner, repo, prNumber, opts)
	if err != nil {
		return "", fmt.Errorf("failed to fetch PR diff: %w", err)
	}
	return diff, nil
}

func (s *GitHubService) PostComment(ctx context.Context, owner, repo string, prNumber int, commentBody string) error {
	comment := &github.IssueComment{
		Body: &commentBody,
	}

	_, _, err := s.Client.Issues.CreateComment(ctx, owner, repo, prNumber, comment)
	if err != nil {
		return fmt.Errorf("failed to post comment: %w", err)
	}

	return nil
}

func (s *GitHubService) HasBotCommented(ctx context.Context, owner, repo string, prNumber int) (bool, error) {
	comments, _, err := s.Client.Issues.ListComments(ctx, owner, repo, prNumber, nil)
	if err != nil {
		return false, err
	}

	for _, comment := range comments {
		if comment.Body != nil && strings.Contains(*comment.Body, "## 🤖 AI Code Review") {
			return true, nil
		}
	}
	return false, nil
}