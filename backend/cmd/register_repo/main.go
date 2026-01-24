package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

const (
	MyGithubUsername = "DHRUVV23"       
	MyRepoName       = "ai-code-review" 
)

func main() {
	
	if err := godotenv.Load(".env"); err != nil {
		godotenv.Load("../.env")
	}


	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"),
	)
	
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf(" Connect error: %v", err)
	}
	defer conn.Close(context.Background())
	ctx := context.Background()

	
	var userID int
	err = conn.QueryRow(ctx, `SELECT id FROM users WHERE username=$1`, MyGithubUsername).Scan(&userID)
	if err != nil {
		err = conn.QueryRow(ctx, `
			INSERT INTO users (github_id, username, email) 
			VALUES (12345, $1, 'test@example.com') 
			ON CONFLICT (email) DO UPDATE SET username = $1
			RETURNING id`, MyGithubUsername).Scan(&userID)
		if err != nil {
			log.Fatalf(" Failed to get User ID: %v", err)
		}
	}


	var repoID int
	

	fullName := fmt.Sprintf("%s/%s", MyGithubUsername, MyRepoName)
	dummyGithubID := time.Now().Unix() 

	
	err = conn.QueryRow(ctx, `SELECT id FROM repositories WHERE name=$1`, MyRepoName).Scan(&repoID)

	if err == nil {
		fmt.Printf(" Found existing repo (ID: %d). Updating info...\n", repoID)
	
		_, err = conn.Exec(ctx, `
			UPDATE repositories 
			SET owner=$1, user_id=$2, full_name=$3 
			WHERE id=$4`, 
			MyGithubUsername, userID, fullName, repoID)
		if err != nil {
			log.Fatalf(" Failed to update repo: %v", err)
		}
	} else {
		fmt.Println("Repo not found. Creating new one...")
		
		
		err = conn.QueryRow(ctx, `
			INSERT INTO repositories (user_id, name, owner, full_name, github_repo_id) 
			VALUES ($1, $2, $3, $4, $5) 
			RETURNING id`, 
			userID, MyRepoName, MyGithubUsername, fullName, dummyGithubID).Scan(&repoID)
		
		if err != nil {
			log.Fatalf(" Failed to create repo: %v", err)
		}
	}

	fmt.Println("------------------------------------------------")
	fmt.Printf("REPO REGISTERED!\n")
	fmt.Printf("Repo Name: %s\n", fullName)
	fmt.Printf(" Repo ID:   %d  <--- USE THIS ID IN YOUR TEST!\n", repoID)
	fmt.Println("------------------------------------------------")
}