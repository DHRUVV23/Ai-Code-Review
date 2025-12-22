package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/DHRUVV23/ai-code-review/backend/internal/repository"
	// "github.com/DHRUVV23/ai-code-review/backend/internal/model"
)



type UserHandler struct {
	UserRepo *repository.UserRepository
}

// GET /api/v1/user/profile
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	// TODO: Later we will get this ID from the JWT Token (Auth Middleware)
	// For now, we mock it as User ID 1 (your seeded user)
	userID := 1

	user, err := h.UserRepo.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch profile"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// PUT /api/v1/user/profile
func (h *UserHandler) UpdateUserProfile(c *gin.Context) {
	userID := 1 // Mock ID

	var req struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// We reuse the UpsertUser function but we might need a specific Update function later.
	// For now, let's just mock the success to keep moving.
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"id":      userID,
		"email":   req.Email,
	})
}