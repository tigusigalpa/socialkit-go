// download demonstrates a v2 async download with wait polling.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	socialkit "github.com/tigusigalpa/socialkit-go"
)

func main() {
	accessKey := os.Getenv("SOCIALKIT_ACCESS_KEY")
	if accessKey == "" {
		log.Fatal("set SOCIALKIT_ACCESS_KEY before running this example")
	}

	client := socialkit.NewClient(accessKey,
		socialkit.WithTimeout(60*time.Second),
	)

	// Start a v2 async download job.
	fmt.Println("Starting download job...")
	job, _, err := client.Downloads.Start(context.Background(), "youtube", &socialkit.V2DownloadRequest{
		URL:     "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		Format:  "mp4",
		Quality: "720p",
	})
	if err != nil {
		log.Fatalf("failed to start download: %v", err)
	}
	fmt.Printf("Job started: %s (status: %s)\n", job.JobID, job.Status)

	// Wait for the job to complete.
	fmt.Println("Waiting for job to complete...")
	finalJob, _, err := client.Downloads.Wait(context.Background(), job.JobID, &socialkit.WaitOptions{
		Interval:      3 * time.Second,
		MaxAttempts:   20,
		BackoffFactor: 1.5,
	})
	if err != nil {
		if errors.Is(err, socialkit.ErrAsyncJobFailed) {
			var asyncErr *socialkit.AsyncJobFailedError
			if errors.As(err, &asyncErr) {
				log.Fatalf("job failed: %s (code: %s, retryable: %v)",
					asyncErr.Message, asyncErr.ErrorCode, asyncErr.Retryable)
			}
		}
		log.Fatalf("wait failed: %v", err)
	}

	fmt.Printf("Job completed: %s\n", finalJob.Status)
	if finalJob.DownloadURL != "" {
		fmt.Printf("Download URL: %s\n", finalJob.DownloadURL)
	}
	if finalJob.Title != "" {
		fmt.Printf("Title: %s\n", finalJob.Title)
	}
	if finalJob.CreditsCost != nil {
		fmt.Printf("Credits cost: %.0f\n", *finalJob.CreditsCost)
	}
	if finalJob.FileSizeMB != "" {
		fmt.Printf("File size: %s MB\n", finalJob.FileSizeMB)
	}
}
