package handler

import (
	"net/http"
	"strconv"

	"github.com/DHRUVV23/ai-code-review/backend/internal/model"
	"github.com/DHRUVV23/ai-code-review/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type RepoHandler struct {
	RepoRepository   *repository.RepoRepository
	ConfigRepository *repository.ConfigRepository
}

type AddRepoRequest struct {
	Name  string `json:"name" binding:"required"`
	Owner string `json:"owner" binding:"required"`
}

// RegisterRepository handles POST /api/v1/repositories
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

// ListRepositories handles GET /api/v1/user/repositories
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

// Helper to extract UserID
func getUserIDFromToken(c *gin.Context) int {
	tokenString := c.GetHeader("Authorization")
	if len(tokenString) > 7 {
		tokenString = tokenString[7:] // Remove "Bearer "
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