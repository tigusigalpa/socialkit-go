package socialkit

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

// TestPagination_CursorRoundtrip verifies that a cursor returned from
// one request is sent back in the next request.
func TestPagination_CursorRoundtrip(t *testing.T) {
	var callCount int
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if r.URL.Path != "/tiktok/comments" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		// Read the request body to check for cursor.
		var body CommentsRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if callCount == 1 {
			if body.Cursor != nil {
				t.Errorf("expected nil cursor on first call, got %v", *body.Cursor)
			}
			writeResponse(t, w, `{"success":true,"data":{"videoId":"123","comments":[{"id":"c1","text":"first"}],"hasMore":true,"cursor":"cursor-page-2"}}`)
			return
		}
		if body.Cursor == nil || *body.Cursor != "cursor-page-2" {
			t.Errorf("expected cursor cursor-page-2, got %v", body.Cursor)
		}
		writeResponse(t, w, `{"success":true,"data":{"videoId":"123","comments":[{"id":"c2","text":"second"}],"hasMore":false}}`)
	})
	defer srv.Close()

	// First page.
	page1, _, err := c.TikTok.Comments(context.Background(), &CommentsRequest{
		URL: "https://tiktok.com/@user/video/123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(page1.Comments) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(page1.Comments))
	}
	if page1.HasMore == nil || !*page1.HasMore {
		t.Error("expected hasMore true")
	}
	if page1.Cursor == nil {
		t.Fatal("expected non-nil cursor")
	}

	// Second page using cursor from first page.
	page2, _, err := c.TikTok.Comments(context.Background(), &CommentsRequest{
		URL:    "https://tiktok.com/@user/video/123",
		Cursor: page1.Cursor,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(page2.Comments) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(page2.Comments))
	}
	if page2.HasMore != nil && *page2.HasMore {
		t.Error("expected hasMore false")
	}
	if callCount != 2 {
		t.Errorf("expected 2 calls, got %d", callCount)
	}
}

// TestPagination_TwitterNextCursor verifies that the Twitter nextCursor
// field roundtrips correctly.
func TestPagination_TwitterNextCursor(t *testing.T) {
	var callCount int
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		var body TweetsRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if callCount == 1 {
			writeResponse(t, w, `{"success":true,"data":{"tweets":[{"id":"t1","text":"first"}],"nextCursor":"twitter-cursor-2"}}`)
			return
		}
		if body.Cursor == nil || *body.Cursor != "twitter-cursor-2" {
			t.Errorf("expected cursor twitter-cursor-2, got %v", body.Cursor)
		}
		writeResponse(t, w, `{"success":true,"data":{"tweets":[{"id":"t2","text":"second"}],"nextCursor":null}}`)
	})
	defer srv.Close()

	page1, _, err := c.Twitter.Tweets(context.Background(), &TweetsRequest{
		URL: "https://twitter.com/user",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page1.NextCursor == nil {
		t.Fatal("expected non-nil nextCursor")
	}

	page2, _, err := c.Twitter.Tweets(context.Background(), &TweetsRequest{
		URL:    "https://twitter.com/user",
		Cursor: page1.NextCursor,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page2.NextCursor != nil {
		t.Errorf("expected nil nextCursor, got %v", *page2.NextCursor)
	}
}
