package handler

import (
	// "context"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/DHRUVV23/ai-code-review/backend/internal/model"
	"github.com/DHRUVV23/ai-code-review/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v50/github" 
	"golang.org/x/oauth2"
)

type RepoHandler struct {
	RepoRepository   *repository.RepoRepository
	ConfigRepository *repository.ConfigRepository
	UserRepository   *repository.UserRepository
}

type AddRepoRequest struct {
	Name  string `json:"name" binding:"required"`
	Owner string `json:"owner" binding:"required"`
}


func (h *RepoHandler) RegisterRepository(c *gin.Context) {

	userID := getUserIDFromToken(c)
	if userID == 0 {
		return
	}

	var req AddRepoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repo, err := h.RepoRepository.CreateRepository(c.Request.Context(), userID, req.Name, req.Owner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create repository"})
		return
	}
	c.JSON(http.StatusCreated, repo)
}

func (h *RepoHandler) ListRepositories(c *gin.Context) {
	userID := getUserIDFromToken(c)
	if userID == 0 {
		return
	}

	repos, err := h.RepoRepository.ListRepositories(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch repos"})
		return
	}
	c.JSON(http.StatusOK, repos)
}

func (h *RepoHandler) GetConfig(c *gin.Context) {
	repoID, _ := strconv.Atoi(c.Param("id"))

	
	config, err := h.ConfigRepository.GetByRepoID(c.Request.Context(), repoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch config"})
		return
	}
	c.JSON(http.StatusOK, config)
}

func (h *RepoHandler) UpdateConfig(c *gin.Context) {
	repoID, _ := strconv.Atoi(c.Param("id"))
	var config model.Configuration
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	config.RepositoryID = repoID
	if err := h.ConfigRepository.UpsertConfig(c.Request.Context(), &config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config"})
		return
	}
	c.JSON(http.StatusOK, config)
}

func (h *RepoHandler) CreateWebhook(c *gin.Context) {
	// 1. Get User ID
	userID := getUserIDFromToken(c)
	if userID == 0 {
		return
	}


	repoID, _ := strconv.Atoi(c.Param("id"))

	repo, err := h.RepoRepository.GetRepositoryByID(c.Request.Context(), repoID)
	if err != nil || repo == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	user, err := h.UserRepository.GetUserByID(c.Request.Context(), userID)
	if err != nil || user.AccessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "GitHub token not found. Please logout and login again."})
		return
	}

	ctx := c.Request.Context()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: user.AccessToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	webhookURL := "https://verona-unabolished-ivy.ngrok-free.dev/webhook" 
	
	webhookSecret := os.Getenv("GITHUB_WEBHOOK_SECRET")

	hookConfig := map[string]interface{}{
		"url":          webhookURL,
		"content_type": "json",
		"secret":       webhookSecret,
	}

	hook := &github.Hook{
		Name:   github.String("web"),
		Active: github.Bool(true),
		Events: []string{"pull_request"},
		Config: hookConfig,
	}

	_, _, err = client.Repositories.CreateHook(ctx, repo.Owner, repo.Name, hook)
	if err != nil {
		if errResp, ok := err.(*github.ErrorResponse); ok && errResp.Response.StatusCode == 422 {
			log.Println(" Webhook already exists, treating as success.")
			c.JSON(http.StatusOK, gin.H{"message": "Webhook already active"})
			return
		}

		log.Printf("Failed to create webhook: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create webhook: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook created successfully!"})
}
func (h *RepoHandler) DeleteRepository(c *gin.Context) {
	userID := getUserIDFromToken(c)
	if userID == 0 { return }
	repoID, _ := strconv.Atoi(c.Param("id"))

	repo, err := h.RepoRepository.GetRepositoryByID(c.Request.Context(), repoID)
	if err != nil || repo == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	user, err := h.UserRepository.GetUserByID(c.Request.Context(), userID)
	if err != nil || user.AccessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "GitHub token invalid"})
		return
	}
	
	ctx := c.Request.Context()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: user.AccessToken})
	client := github.NewClient(oauth2.NewClient(ctx, ts))

	hooks, _, err := client.Repositories.ListHooks(ctx, repo.Owner, repo.Name, nil)
	if err == nil {
		targetURL := "ngrok-free.app"
		for _, hook := range hooks {
			config := hook.Config
			if url, ok := config["url"].(string); ok && strings.Contains(url, targetURL) {
				client.Repositories.DeleteHook(ctx, repo.Owner, repo.Name, hook.GetID())
				log.Printf("🗑️ Deleted GitHub Webhook ID: %d", hook.GetID())
				break
			}
		}
	} else {
		log.Printf("Could not list GitHub hooks (might already be deleted): %v", err)
	}

	if err := h.RepoRepository.DeleteRepository(ctx, repoID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete repository"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Repository deleted and unlinked successfully"})
}