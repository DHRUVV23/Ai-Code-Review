package handler

import (
	"net/http"
	"strconv"

	"github.com/DHRUVV23/ai-code-review/backend/internal/model"
	"github.com/DHRUVV23/ai-code-review/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v50/github" // Ensure this matches your go.mod
	"golang.org/x/oauth2"
)

type RepoHandler struct {
	RepoRepository   *repository.RepoRepository
	ConfigRepository *repository.ConfigRepository
	UserRepository   *repository.UserRepository // <--- ADDED THIS FIELD
}

type AddRepoRequest struct {
	Name  string `json:"name" binding:"required"`
	Owner string `json:"owner" binding:"required"`
}

// RegisterRepository handles POST /api/v1/repositories
func (h *RepoHandler) RegisterRepository(c *gin.Context) {
	userID := getUserIDFromToken(c)
	if userID == 0 { return }

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

// ListRepositories handles GET /api/v1/user/repositories
func (h *RepoHandler) ListRepositories(c *gin.Context) {
	userID := getUserIDFromToken(c)
	if userID == 0 { return }

	repos, err := h.RepoRepository.ListRepositories(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch repos"})
		return
	}
	c.JSON(http.StatusOK, repos)
}

// GetConfig handles GET /repositories/:id
func (h *RepoHandler) GetConfig(c *gin.Context) {
	repoID, _ := strconv.Atoi(c.Param("id"))
	config, err := h.ConfigRepository.GetByRepoID(c.Request.Context(), repoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch config"})
		return
	}
	c.JSON(http.StatusOK, config)
}

// UpdateConfig handles PUT /repositories/:id/config
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

// CreateWebhook handles POST /repositories/:id/webhook
func (h *RepoHandler) CreateWebhook(c *gin.Context) {
	// 1. Get User ID
	userID := getUserIDFromToken(c)
	if userID == 0 { return }

	// 2. Get Repo ID
	repoID, _ := strconv.Atoi(c.Param("id"))
	
	// 3. Get Repository Details
	repo, err := h.RepoRepository.GetRepositoryByID(c.Request.Context(), repoID)
	if err != nil || repo == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	// 4. Get User's GitHub Token
	user, err := h.UserRepository.GetUserByID(c.Request.Context(), userID)
	if err != nil || user.AccessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "GitHub token not found. Please logout and login again."})
		return
	}

	// 5. Connect to GitHub API
	ctx := c.Request.Context()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: user.AccessToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// 6. Define the Webhook
	// FIX: Use a map[string]interface{} for Config, not a struct
	webhookURL := "https://example.com/webhooks" // Change this to your real URL later
	
	hookConfig := map[string]interface{}{
		"url":          webhookURL,
		"content_type": "json",
		// "secret": "your_secret", // Optional
	}
    
	hook := &github.Hook{
		Name:   github.String("web"),
		Active: github.Bool(true),
		Events: []string{"pull_request"},
		Config: hookConfig, // Now passing the map!
	}

	// 7. Create the Hook
	_, _, err = client.Repositories.CreateHook(ctx, repo.Owner, repo.Name, hook)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create webhook: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook created successfully!"})
}

// Helper
func getUserIDFromToken(c *gin.Context) int {
	tokenString := c.GetHeader("Authorization")
	if len(tokenString) > 7 {
		tokenString = tokenString[7:]
	}
	token, _, _ := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if idFloat, ok := claims["user_id"].(float64); ok {
			return int(idFloat)
		}
	}
	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
	return 0
}