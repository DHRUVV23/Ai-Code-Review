package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/DHRUVV23/ai-code-review/backend/internal/repository"
)

type ReviewHandler struct {
	ReviewRepository *repository.ReviewRepository
}

// GET /api/v1/repositories/:id/reviews
func (h *ReviewHandler) ListReviews(c *gin.Context) {
	repoID, _ := strconv.Atoi(c.Param("id"))

	reviews, err := h.ReviewRepository.GetReviewsByRepoID(c.Request.Context(), repoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	// If list is empty, return empty array instead of null
	if reviews == nil {
		c.JSON(http.StatusOK, []string{})
		return
	}

	c.JSON(http.StatusOK, reviews)
}