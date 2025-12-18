package main

import (
	"log"
	"github.com/DHRUVV23/ai-code-review/backend/internal/worker"
)

func main() {
	// 1. Connect to Redis
	worker.InitClient()
	defer worker.CloseClient()

	// 2. Fake a PR Review Job
	repoID := 1
	prNumber := 99
	repoName := "fake/repo"

	log.Println("🚀 Sending fake job to queue...")
	task, err := worker.NewReviewPRTask(repoID, prNumber, repoName)
	if err != nil {
		log.Fatal(err)
	}

	// 3. Send it
	info, err := worker.Client.Enqueue(task)
	if err != nil {
		log.Fatal(err)
	}
	
	log.Printf("✅ Success! Job ID: %s", info.ID)
}