// errors demonstrates error handling patterns for the SocialKit SDK.
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
		socialkit.WithTimeout(15*time.Second),
		socialkit.WithRetry(3, time.Second),
	)

	// Example 1: Check for specific error types using errors.Is.
	_, _, err := client.YouTube.Transcript(context.Background(), &socialkit.FetchRequest{
		URL: "https://www.youtube.com/watch?v=invalid",
	})
	if err != nil {
		switch {
		case errors.Is(err, socialkit.ErrUnauthorized):
			fmt.Println("Error: invalid access key (401)")
		case errors.Is(err, socialkit.ErrNotFound):
			fmt.Println("Error: video not found (404)")
		case errors.Is(err, socialkit.ErrRateLimited):
			fmt.Println("Error: rate limited after retries (429)")
		case errors.Is(err, socialkit.ErrBadRequest):
			fmt.Println("Error: bad request (400)")
		case errors.Is(err, socialkit.ErrServer):
			fmt.Println("Error: server error (500/503)")
		default:
			fmt.Printf("Error: %v\n", err)
		}

		// Extract the APIError for more details.
		var apiErr *socialkit.APIError
		if errors.As(err, &apiErr) {
			fmt.Printf("  Status: %d\n", apiErr.StatusCode)
			fmt.Printf("  Code: %s\n", apiErr.Code)
			fmt.Printf("  Message: %s\n", apiErr.Message)
		}
	}

	// Example 2: Handle insufficient credits.
	_, _, err = client.YouTube.Summarize(context.Background(), &socialkit.SummaryRequest{
		URL: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
	})
	if err != nil {
		if errors.Is(err, socialkit.ErrInsufficientCredits) {
			var icErr *socialkit.InsufficientCreditsError
			if errors.As(err, &icErr) {
				fmt.Printf("\nInsufficient credits:\n")
				fmt.Printf("  Required: %.0f\n", icErr.RequiredCredits)
				fmt.Printf("  Remaining: %.0f\n", icErr.RemainingCredits)
				fmt.Printf("  Shortfall: %.0f\n", icErr.ShortfallCredits)
				if icErr.ResetAt != nil {
					fmt.Printf("  Reset at: %s\n", icErr.ResetAt.Format(time.RFC3339))
				}
				if icErr.UpgradeURL != "" {
					fmt.Printf("  Upgrade: %s\n", icErr.UpgradeURL)
				}
				if icErr.TopUpURL != "" {
					fmt.Printf("  Top up: %s\n", icErr.TopUpURL)
				}
			}
		}
	}

	// Example 3: Handle async job failure.
	_, _, err = client.Downloads.Wait(context.Background(), "job-id", &socialkit.WaitOptions{
		Interval:    2 * time.Second,
		MaxAttempts: 10,
	})
	if err != nil {
		if errors.Is(err, socialkit.ErrAsyncJobFailed) {
			var asyncErr *socialkit.AsyncJobFailedError
			if errors.As(err, &asyncErr) {
				fmt.Printf("\nAsync job failed:\n")
				fmt.Printf("  Error: %s\n", asyncErr.Message)
				fmt.Printf("  Code: %s\n", asyncErr.ErrorCode)
				fmt.Printf("  Retryable: %v\n", asyncErr.Retryable)
			}
		}
	}
}
