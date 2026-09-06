// basic demonstrates a status check and a YouTube transcript request.
package main

import (
	"context"
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

	// Check API status (no key required).
	status, _, err := client.Status.Get(context.Background())
	if err != nil {
		log.Printf("status check failed: %v", err)
	} else {
		fmt.Printf("API status: %s (generated at %s)\n", status.Overall, status.GeneratedAt)
		for _, tool := range status.Tools {
			fmt.Printf("  %s: %s\n", tool.Name, tool.Status)
		}
	}

	// Fetch a YouTube transcript.
	transcript, meta, err := client.YouTube.Transcript(context.Background(), &socialkit.FetchRequest{
		URL: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
	})
	if err != nil {
		log.Fatalf("transcript failed: %v", err)
	}

	fmt.Println("\nTranscript:")
	if transcript.WordCount != nil {
		fmt.Printf("Transcript (%d words):\n", *transcript.WordCount)
	}
	if transcript.DurationSeconds != nil {
		fmt.Printf("Duration: %.1f seconds\n", *transcript.DurationSeconds)
	}
	fmt.Printf("First 200 chars: %s\n", truncate(transcript.Transcript, 200))

	if meta != nil && meta.CreditsUsed != nil {
		fmt.Printf("Credits used: %.0f\n", *meta.CreditsUsed)
	}
	if meta != nil && meta.CreditsRemaining != nil {
		fmt.Printf("Credits remaining: %.0f\n", *meta.CreditsRemaining)
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
