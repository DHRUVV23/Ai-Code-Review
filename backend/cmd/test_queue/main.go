package main

import (
	"log"
	"github.com/DHRUVV23/ai-code-review/backend/internal/worker"
)

func main() {
	
	worker.InitClient()
	defer worker.CloseClient()


	repoID := 1
	prNumber := 1
	repoName := "DHRUVV23/OmniScribe"

	log.Println("🚀 Sending fake job to queue...")
	task, err := worker.NewReviewPRTask(repoID, prNumber, repoName)
	if err != nil {
		log.Fatal(err)
	}

	
	info, err := worker.Client.Enqueue(task)
	if err != nil {
		log.Fatal(err)
	}
	
	log.Printf("Success! Job ID: %s", info.ID)
}