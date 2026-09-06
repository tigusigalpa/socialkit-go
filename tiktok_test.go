package socialkit

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestTikTokService_Transcript(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tiktok/transcript" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		writeResponse(t, w, `{"success":true,"data":{"url":"https://tiktok.com/@user/video/123","videoId":"123","transcript":"hello tiktok","transcriptSegments":[{"text":"hello","start":0,"duration":2.0}],"creditsUsed":1}}`)
	})
	defer srv.Close()

	transcript, _, err := c.TikTok.Transcript(context.Background(), &FetchRequest{
		URL: "https://tiktok.com/@user/video/123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if transcript.Transcript != "hello tiktok" {
		t.Errorf("unexpected transcript: %s", transcript.Transcript)
	}
}

func TestTikTokService_Summarize(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeResponse(t, w, `{"success":true,"data":{"url":"https://tiktok.com/@user/video/123","summary":"A fun tiktok","mainTopics":["fun"],"keyPoints":["point1"],"tone":"casual","targetAudience":"everyone","quotes":[]}}`)
	})
	defer srv.Close()

	summary, _, err := c.TikTok.Summarize(context.Background(), &SummaryRequest{
		URL: "https://tiktok.com/@user/video/123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Summary != "A fun tiktok" {
		t.Errorf("unexpected summary: %s", summary.Summary)
	}
}

func TestTikTokService_Stats(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tiktok/stats" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"url":"https://tiktok.com/@user/video/123","videoId":"123","title":"Cool video","channelName":"user","channelLink":"https://tiktok.com/@user","likes":1000,"comments":50,"collects":200,"shares":100,"views":50000,"description":"desc","duration":"00:30","thumbnailUrl":"https://img.com/thumb.jpg","publishedAt":"2026-01-01","contentType":"video","isShortForm":true}}`)
	})
	defer srv.Close()

	stats, _, err := c.TikTok.Stats(context.Background(), &FetchRequest{
		URL: "https://tiktok.com/@user/video/123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Likes == nil || *stats.Likes != 1000 {
		t.Errorf("expected likes 1000, got %v", stats.Likes)
	}
	if stats.Collects == nil || *stats.Collects != 200 {
		t.Errorf("expected collects 200, got %v", stats.Collects)
	}
	if stats.IsShortForm == nil || !*stats.IsShortForm {
		t.Error("expected isShortForm true")
	}
}

func TestTikTokService_Comments(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeResponse(t, w, `{"success":true,"data":{"videoId":"123","comments":[{"id":"c1","author":"User1","username":"user1","text":"Nice!","likes":10,"avatar":"https://img.com/a.jpg","createTime":1700000000,"replyCount":2}],"commentCount":100,"hasMore":true,"cursor":"abc123"}}`)
	})
	defer srv.Close()

	comments, _, err := c.TikTok.Comments(context.Background(), &CommentsRequest{
		URL: "https://tiktok.com/@user/video/123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments.Comments) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(comments.Comments))
	}
	if comments.HasMore == nil || !*comments.HasMore {
		t.Error("expected hasMore true")
	}
	if comments.Cursor == nil || *comments.Cursor != "abc123" {
		t.Errorf("expected cursor abc123, got %v", comments.Cursor)
	}
}

func TestTikTokService_ChannelStats(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tiktok/channel-stats" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"profileUrl":"https://tiktok.com/@user","username":"user","nickname":"User","signature":"bio","verified":false,"avatar":"https://img.com/a.jpg","followers":50000,"following":100,"totalLikes":1000000,"totalVideos":200,"bioLink":"https://example.com"}}`)
	})
	defer srv.Close()

	stats, _, err := c.TikTok.ChannelStats(context.Background(), &FetchRequest{
		URL: "https://tiktok.com/@user",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Followers == nil || *stats.Followers != 50000 {
		t.Errorf("expected followers 50000, got %v", stats.Followers)
	}
}

func TestTikTokService_ChannelVideos(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tiktok/channel-videos" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"profileUrl":"https://tiktok.com/@user","channelName":"user","results":[{"videoId":"v1","description":"video 1","url":"https://tiktok.com/@user/video/v1","thumbnail":"https://img.com/t.jpg","duration":"00:30","createTime":1700000000,"views":1000,"likes":50,"comments":5,"shares":10,"collects":3}],"hasMore":false,"cursor":null}}`)
	})
	defer srv.Close()

	videos, _, err := c.TikTok.ChannelVideos(context.Background(), &ListRequest{
		URL: "https://tiktok.com/@user",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(videos.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(videos.Results))
	}
	if videos.Results[0].VideoID != "v1" {
		t.Errorf("unexpected videoId: %s", videos.Results[0].VideoID)
	}
}

func TestTikTokService_ChannelVideos_LimitExceeded(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not make a request")
	})
	defer srv.Close()

	limit := 101
	_, _, err := c.TikTok.ChannelVideos(context.Background(), &ListRequest{
		URL:   "https://tiktok.com/@user",
		Limit: &limit,
	})
	if err == nil {
		t.Fatal("expected error for limit > 100")
	}
}

func TestTikTokService_Search(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tiktok/search" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"query":"dance","results":[{"id":"v1","desc":"dance video","url":"https://tiktok.com/@user/video/v1","createTime":"2026-01-01","author":{"username":"user"},"stats":{"likes":100}}],"hasMore":false}}`)
	})
	defer srv.Close()

	results, _, err := c.TikTok.Search(context.Background(), &SearchRequest{
		Query: "dance",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results.Results))
	}
}

func TestTikTokService_Search_LimitExceeded(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not make a request")
	})
	defer srv.Close()

	limit := 101
	_, _, err := c.TikTok.Search(context.Background(), &SearchRequest{
		Query: "dance",
		Limit: &limit,
	})
	if err == nil {
		t.Fatal("expected error for limit > 100")
	}
}

func TestTikTokService_HashtagSearch(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tiktok/hashtag-search" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body HashtagSearchRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body.Hashtag != "dance" {
			t.Errorf("expected hashtag dance, got %s", body.Hashtag)
		}
		writeResponse(t, w, `{"success":true,"data":{"hashtag":"dance","results":[{"id":"v1","desc":"dance","url":"https://tiktok.com/@user/video/v1","author":{},"stats":{}}],"hasMore":false}}`)
	})
	defer srv.Close()

	results, _, err := c.TikTok.HashtagSearch(context.Background(), &HashtagSearchRequest{
		Hashtag: "dance",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results.Hashtag != "dance" {
		t.Errorf("unexpected hashtag: %s", results.Hashtag)
	}
}

func TestTikTokService_HashtagSearch_LimitExceeded(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not make a request")
	})
	defer srv.Close()

	limit := 101
	_, _, err := c.TikTok.HashtagSearch(context.Background(), &HashtagSearchRequest{
		Hashtag: "dance",
		Limit:   &limit,
	})
	if err == nil {
		t.Fatal("expected error for limit > 100")
	}
}

func TestTikTokService_Download(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tiktok/download" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"url":"https://tiktok.com/@user/video/123","title":"Cool","duration":"00:30","downloadUrl":"https://download.socialkit.dev/file.mp4","fileSizeMB":"25.0","format":"mp4","quality":"no-watermark","expiresIn":"1h","creditsUsed":1}}`)
	})
	defer srv.Close()

	result, _, err := c.TikTok.Download(context.Background(), &DownloadRequest{
		URL: "https://tiktok.com/@user/video/123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.DownloadURL == "" {
		t.Error("expected non-empty downloadUrl")
	}
}

func TestTikTokService_BulkTranscript(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tiktok/transcript/bulk" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":[]}`)
	})
	defer srv.Close()

	_, _, err := c.TikTok.BulkTranscript(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTikTokService_BulkComments(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tiktok/comments/bulk" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":[]}`)
	})
	defer srv.Close()

	_, _, err := c.TikTok.BulkComments(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTikTokService_BulkStats(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tiktok/stats/bulk" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":[]}`)
	})
	defer srv.Close()

	_, _, err := c.TikTok.BulkStats(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTikTokService_BulkChannelStats(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tiktok/channel-stats/bulk" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":[]}`)
	})
	defer srv.Close()

	_, _, err := c.TikTok.BulkChannelStats(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTikTokService_BulkSummarize(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tiktok/summarize/bulk" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":[]}`)
	})
	defer srv.Close()

	_, _, err := c.TikTok.BulkSummarize(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
