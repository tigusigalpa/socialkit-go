// pagination demonstrates cursor pagination without parsing or altering
// the opaque cursor returned by SocialKit.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	socialkit "github.com/tigusigalpa/socialkit-go"
)

func main() {
	accessKey := os.Getenv("SOCIALKIT_ACCESS_KEY")
	if accessKey == "" {
		log.Fatal("set SOCIALKIT_ACCESS_KEY before running this example")
	}

	client := socialkit.NewClient(accessKey)
	ctx := context.Background()
	request := &socialkit.CommentsRequest{
		URL: "https://www.tiktok.com/@scout2015/video/6718335390845095173",
	}

	for page := 1; ; page++ {
		result, meta, err := client.TikTok.Comments(ctx, request)
		if err != nil {
			log.Fatalf("get page %d: %v", page, err)
		}
		fmt.Printf("page %d: %d comments\n", page, len(result.Comments))
		if meta != nil && meta.CreditsUsed != nil {
			fmt.Printf("  credits used: %.0f\n", *meta.CreditsUsed)
		}

		if result.HasMore == nil || !*result.HasMore || result.Cursor == nil {
			break
		}
		// Cursors are opaque: pass them back verbatim and never construct them.
		request.Cursor = result.Cursor
	}
}
