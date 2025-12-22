package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v57/github"
	"github.com/DHRUVV23/ai-code-review/backend/internal/repository"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github" 
	// "golang.org/x/oauth2/github"
)

type AuthHandler struct {
	UserRepo *repository.UserRepository
}

// oauthConf setup
var oauthConf = &oauth2.Config{
	ClientID:     "", // Set in Init() or use os.Getenv directly below
	ClientSecret: "",
	Scopes:       []string{"repo", "user:email"},
	Endpoint:     githuboauth.Endpoint, 
}

func init() {
	oauthConf.ClientID = os.Getenv("GITHUB_CLIENT_ID")
	oauthConf.ClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
}

// GitHubLogin redirects user to GitHub
func (h *AuthHandler) GitHubLogin(c *gin.Context) {
	url := oauthConf.AuthCodeURL("random_state_string", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GitHubCallback handles the return from GitHub
func (h *AuthHandler) GitHubCallback(c *gin.Context) {
	code := c.Query("code")
	
	// 1. Exchange code for token
	token, err := oauthConf.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// 2. Fetch User Info from GitHub
	oauthClient := oauthConf.Client(context.Background(), token)
	client := github.NewClient(oauthClient)

	githubUser, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	// 3. Extract Data (Safely handle pointers)
	githubID := githubUser.GetID()
	username := githubUser.GetLogin()
	email := githubUser.GetEmail()

	if email == "" {
		// Fallback: If email is private, we might need to fetch it explicitly.
		// For now, we'll create a dummy placeholder or error out.
		// A robust app would make a second call to /user/emails here.
		email = fmt.Sprintf("%s@no-email.github.com", username) 
	}

	// 4. Upsert User into Database
	// FIX: We now pass 4 arguments: Context, GithubID, Username, Email
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

	tokenString, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// 6. Success! Return the token
	c.JSON(http.StatusOK, gin.H{
		"token":    tokenString,
		"user_id":  userID,
		"username": username,
	})
}