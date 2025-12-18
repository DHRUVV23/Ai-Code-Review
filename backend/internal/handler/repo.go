package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/DHRUVV23/ai-code-review/backend/internal/model"
	"github.com/DHRUVV23/ai-code-review/backend/internal/repository"
)

type RepoHandler struct {
	RepoRepository   *repository.RepoRepository
	ConfigRepository *repository.ConfigRepository
}

// --- 1. List Repositories (The Missing Function) ---
func (h *RepoHandler) ListRepositories(c *gin.Context) {
	// For now, hardcode ID 1. Later we get this from the JWT token.
	repoID := 1
	repo, err := h.RepoRepository.GetRepositoryByID(c.Request.Context(), repoID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch repository"})
		return
	}

	if repo == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No repository found"})
		return
	}

	c.JSON(http.StatusOK, repo)
}

// --- 2. Get Config ---
func (h *RepoHandler) GetConfig(c *gin.Context) {
	repoID, _ := strconv.Atoi(c.Param("id"))

	config, err := h.ConfigRepository.GetByRepoID(c.Request.Context(), repoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch config"})
		return
	}
	c.JSON(http.StatusOK, config)
}

// --- 3. Update Config ---
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