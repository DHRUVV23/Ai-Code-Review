package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"fmt"

	// Using your preferred import structure
	"github.com/DHRUVV23/ai-code-review/backend/internal/config"
	"github.com/DHRUVV23/ai-code-review/backend/internal/repository"
	
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type AuthHandler struct {
	UserRepo *repository.UserRepository
	Config   *config.Config // <--- Added this missing field
}

// GitHubLogin redirects the user to GitHub
func (h *AuthHandler) GitHubLogin(c *gin.Context) {
	// Force the consent screen to ensure we get permissions
	conf := &oauth2.Config{
		ClientID:     h.Config.GithubClientID,
		ClientSecret: h.Config.GithubClientSecret,
		RedirectURL:  "http://localhost:8080/auth/github/callback",
		Endpoint:     github.Endpoint,
		Scopes:       []string{"user:email", "read:user", "repo", "admin:repo_hook"},
	}
	// AccessTypeOffline asks for a refresh token (optional), ApprovalForce forces the screen
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GitHubCallback handles the code -> token exchange
func (h *AuthHandler) GitHubCallback(c *gin.Context) {
	code := c.Query("code")

	conf := &oauth2.Config{
		ClientID:     h.Config.GithubClientID,
		ClientSecret: h.Config.GithubClientSecret,
		RedirectURL:  "http://localhost:8080/auth/github/callback",
		Endpoint:     github.Endpoint,
	}

	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	client := conf.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var githubUser struct {
		ID       int64  `json:"id"`
		Login    string `json:"login"`
		Email    string `json:"email"`
		Name     string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
		return
	}

	// Save User + Token to DB
	// Note: UpsertUser now takes 4 args: context, githubID, username, email, token
	userID, err := h.UserRepo.UpsertUser(c.Request.Context(), githubUser.ID, githubUser.Login, githubUser.Email, token.AccessToken)
	if err != nil {
		fmt.Println("❌ DATABASE ERROR:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		return
	}

	// Generate JWT
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     token.Expiry.Unix(),
	})
	tokenString, _ := jwtToken.SignedString([]byte("your_jwt_secret")) 

	// Redirect to Frontend with Token
	c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/dashboard?token="+tokenString)
}

// GetUserProfile - Handles GET /api/v1/user/profile
func (h *AuthHandler) GetUserProfile(c *gin.Context) {
	userID := getUserIDFromToken(c)
	if userID == 0 {
		return
	}

	user, err := h.UserRepo.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUserProfile - Handles PUT /api/v1/user/profile
func (h *AuthHandler) UpdateUserProfile(c *gin.Context) {
	// Just a placeholder for now to satisfy the interface
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated"})
}

// Helper function
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