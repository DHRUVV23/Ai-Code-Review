package handler

import (
	"context"
	// "encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/go-resty/resty/v2" // You might need to go get this: github.com/go-resty/resty/v2
	"github.com/DHRUVV23/ai-code-review/backend/internal/repository"
)

type AuthHandler struct {
	UserRepo *repository.UserRepository
}

// 1. Login: Redirects user to GitHub
func (h *AuthHandler) GitHubLogin(c *gin.Context) {
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	// We redirect to GitHub's authorization page
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&scope=read:user",
		clientID,
	)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// 2. Callback: GitHub sends user back here
func (h *AuthHandler) GitHubCallback(c *gin.Context) {
	code := c.Query("code") // The temporary code from GitHub
	
	// A. Exchange code for Access Token
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")

	client := resty.New()
	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	
	// Call GitHub API to get token
	_, err := client.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]string{
			"client_id":     clientID,
			"client_secret": clientSecret,
			"code":          code,
		}).
		SetResult(&tokenResp).
		Post("https://github.com/login/oauth/access_token")

	if err != nil || tokenResp.AccessToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get access token"})
		return
	}

	// B. Fetch User Profile using the Token
	var githubUser struct {
		ID    int64  `json:"id"`
		Login string `json:"login"`
		Email string `json:"email"`
	}
	
	_, err = client.R().
		SetAuthToken(tokenResp.AccessToken).
		SetResult(&githubUser).
		Get("https://api.github.com/user")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
		return
	}

	// C. Save User to Database
	user, err := h.UserRepo.UpsertUser(context.Background(), githubUser.ID, githubUser.Login, githubUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// D. Generate JWT (Session Token)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hours
	})
	
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	// Success! Return the token
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user":  user,
	})
}