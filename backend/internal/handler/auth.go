package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/DHRUVV23/ai-code-review/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v50/github" // Ensure this matches your go.mod version
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

type AuthHandler struct {
	UserRepo *repository.UserRepository
}

// Helper to get config (Corrects the init() bug by loading when needed)
func getOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes:       []string{"repo", "user:email"},
		Endpoint:     githuboauth.Endpoint,
		RedirectURL:  "http://localhost:8080/auth/github/callback",
	}
}

// GitHubLogin redirects user to GitHub
func (h *AuthHandler) GitHubLogin(c *gin.Context) {
	conf := getOAuthConfig()
	// Redirect user to GitHub's consent page
	url := conf.AuthCodeURL("random_state_string", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GitHubCallback handles the return from GitHub
func (h *AuthHandler) GitHubCallback(c *gin.Context) {
	code := c.Query("code")
	conf := getOAuthConfig()

	// 1. Exchange code for token
	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// 2. Fetch User Info from GitHub
	oauthClient := conf.Client(context.Background(), token)
	client := github.NewClient(oauthClient)

	githubUser, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	// 3. Extract Data
	githubID := githubUser.GetID()
	username := githubUser.GetLogin()
	email := githubUser.GetEmail()
	if email == "" {
		email = fmt.Sprintf("%s@no-email.github.com", username)
	}

	// 4. Upsert User into Database (Your smart logic)
	userID, err := h.UserRepo.UpsertUser(c.Request.Context(), githubID, username, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		return
	}

	// 5. Generate JWT Token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	// Ensure JWT_SECRET is in your .env
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_secret_please_change" // Fallback to prevent crash
	}
	
	tokenString, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// 6. REDIRECT to Frontend (The Critical Fix)
	// Instead of showing JSON, we send the user back to the React App with the token
	frontendURL := fmt.Sprintf("http://localhost:3000/dashboard?token=%s&user=%s", tokenString, username)
	c.Redirect(http.StatusFound, frontendURL)
}