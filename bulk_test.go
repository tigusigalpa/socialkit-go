package socialkit

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

// TestBulk_AllEndpoints verifies that all bulk experimental endpoints
// accept raw JSON and return raw JSON responses.
func TestBulk_AllEndpoints(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		callFunc func(c *Client, body json.RawMessage) (json.RawMessage, *ResponseMeta, error)
	}{
		{
			name: "YouTubeBulkTranscript",
			path: "/youtube/transcript/bulk",
			callFunc: func(c *Client, body json.RawMessage) (json.RawMessage, *ResponseMeta, error) {
				return c.YouTube.BulkTranscript(context.Background(), body)
			},
		},
		{
			name: "YouTubeBulkComments",
			path: "/youtube/comments/bulk",
			callFunc: func(c *Client, body json.RawMessage) (json.RawMessage, *ResponseMeta, error) {
				return c.YouTube.BulkComments(context.Background(), body)
			},
		},
		{
			name: "YouTubeBulkStats",
			path: "/youtube/stats/bulk",
			callFunc: func(c *Client, body json.RawMessage) (json.RawMessage, *ResponseMeta, error) {
				return c.YouTube.BulkStats(context.Background(), body)
			},
		},
		{
			name: "YouTubeBulkSummarize",
			path: "/youtube/summarize/bulk",
			callFunc: func(c *Client, body json.RawMessage) (json.RawMessage, *ResponseMeta, error) {
				return c.YouTube.BulkSummarize(context.Background(), body)
			},
		},
		{
			name: "TikTokBulkTranscript",
			path: "/tiktok/transcript/bulk",
			callFunc: func(c *Client, body json.RawMessage) (json.RawMessage, *ResponseMeta, error) {
				return c.TikTok.BulkTranscript(context.Background(), body)
			},
		},
		{
			name: "TikTokBulkComments",
			path: "/tiktok/comments/bulk",
			callFunc: func(c *Client, body json.RawMessage) (json.RawMessage, *ResponseMeta, error) {
				return c.TikTok.BulkComments(context.Background(), body)
			},
		},
		{
			name: "TikTokBulkStats",
			path: "/tiktok/stats/bulk",
			callFunc: func(c *Client, body json.RawMessage) (json.RawMessage, *ResponseMeta, error) {
				return c.TikTok.BulkStats(context.Background(), body)
			},
		},
		{
			name: "TikTokBulkChannelStats",
			path: "/tiktok/channel-stats/bulk",
			callFunc: func(c *Client, body json.RawMessage) (json.RawMessage, *ResponseMeta, error) {
				return c.TikTok.BulkChannelStats(context.Background(), body)
			},
		},
		{
			name: "TikTokBulkSummarize",
			path: "/tiktok/summarize/bulk",
			callFunc: func(c *Client, body json.RawMessage) (json.RawMessage, *ResponseMeta, error) {
				return c.TikTok.BulkSummarize(context.Background(), body)
			},
		},
		{
			name: "InstagramBulkTranscript",
			path: "/instagram/transcript/bulk",
			callFunc: func(c *Client, body json.RawMessage) (json.RawMessage, *ResponseMeta, error) {
				return c.Instagram.BulkTranscript(context.Background(), body)
			},
		},
		{
			name: "InstagramBulkStats",
			path: "/instagram/stats/bulk",
			callFunc: func(c *Client, body json.RawMessage) (json.RawMessage, *ResponseMeta, error) {
				return c.Instagram.BulkStats(context.Background(), body)
			},
		},
		{
			name: "InstagramBulkChannelStats",
			path: "/instagram/channel-stats/bulk",
			callFunc: func(c *Client, body json.RawMessage) (json.RawMessage, *ResponseMeta, error) {
				return c.Instagram.BulkChannelStats(context.Background(), body)
			},
		},
		{
			name: "InstagramBulkSummarize",
			path: "/instagram/summarize/bulk",
			callFunc: func(c *Client, body json.RawMessage) (json.RawMessage, *ResponseMeta, error) {
				return c.Instagram.BulkSummarize(context.Background(), body)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != tt.path {
					t.Errorf("unexpected path: got %s, want %s", r.URL.Path, tt.path)
				}
				if r.Method != http.MethodPost {
					t.Errorf("expected POST, got %s", r.Method)
				}
				// Verify the body is passed through as raw JSON.
				var body map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatalf("failed to decode body: %v", err)
				}
				if body["key"] != "value" {
					t.Errorf("expected key=value in body, got %v", body["key"])
				}
				writeResponse(t, w, `{"success":true,"data":{"results":["item1","item2"]}}`)
			})
			defer srv.Close()

			result, meta, err := tt.callFunc(c, json.RawMessage(`{"key":"value"}`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) == 0 {
				t.Error("expected non-empty result")
			}
			if meta == nil {
				t.Error("expected non-nil meta")
			}
		})
	}
}

// TestBulk_ErrorHandling verifies that bulk endpoints properly propagate
// API errors.
func TestBulk_ErrorHandling(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		writeResponse(t, w, `{"message":"invalid bulk request","code":"invalid_request"}`)
	})
	defer srv.Close()

	_, _, err := c.YouTube.BulkTranscript(context.Background(), json.RawMessage(`{}`))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
