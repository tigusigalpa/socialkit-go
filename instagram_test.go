package socialkit

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestInstagramService_Transcript(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/instagram/transcript" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		writeResponse(t, w, `{"success":true,"data":{"url":"https://instagram.com/reel/abc","videoId":"abc","transcript":"hello ig","transcriptSegments":[{"text":"hello","start":0,"duration":1.0}],"creditsUsed":1}}`)
	})
	defer srv.Close()

	transcript, _, err := c.Instagram.Transcript(context.Background(), &FetchRequest{
		URL: "https://instagram.com/reel/abc",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if transcript.Transcript != "hello ig" {
		t.Errorf("unexpected transcript: %s", transcript.Transcript)
	}
}

func TestInstagramService_Summarize(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeResponse(t, w, `{"success":true,"data":{"url":"https://instagram.com/reel/abc","summary":"An IG reel","mainTopics":["topic"],"keyPoints":["point"],"tone":"fun","targetAudience":"all","quotes":[]}}`)
	})
	defer srv.Close()

	summary, _, err := c.Instagram.Summarize(context.Background(), &SummaryRequest{
		URL: "https://instagram.com/reel/abc",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Summary != "An IG reel" {
		t.Errorf("unexpected summary: %s", summary.Summary)
	}
}

func TestInstagramService_Stats(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/instagram/stats" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"postUrl":"https://instagram.com/p/abc","id":"123","shortcode":"abc","title":"Test","description":"desc","views":500,"likes":100,"comments":10,"publishedAt":"2026-01-01","author":"user","authorLink":"https://instagram.com/user","duration":"00:30","thumbnail":"https://img.com/t.jpg","isVideo":true,"contentType":"reel"}}`)
	})
	defer srv.Close()

	stats, _, err := c.Instagram.Stats(context.Background(), &FetchRequest{
		URL: "https://instagram.com/p/abc",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Shortcode != "abc" {
		t.Errorf("unexpected shortcode: %s", stats.Shortcode)
	}
	if stats.IsVideo == nil || !*stats.IsVideo {
		t.Error("expected isVideo true")
	}
}

func TestInstagramService_Comments(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeResponse(t, w, `{"success":true,"data":{"postUrl":"https://instagram.com/p/abc","comments":[{"id":"c1","text":"Nice!","likeCount":5,"replyCount":1,"user":{"username":"user1"}}],"commentCount":50,"hasMore":true,"cursor":"cur123"}}`)
	})
	defer srv.Close()

	comments, _, err := c.Instagram.Comments(context.Background(), &CommentsRequest{
		URL: "https://instagram.com/p/abc",
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
}

func TestInstagramService_ChannelStats(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/instagram/channel-stats" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"profileUrl":"https://instagram.com/user","userId":"123","username":"user","fullName":"User Name","verified":true,"followers":10000,"following":500,"totalPosts":200,"bio":"bio","avatar":"https://img.com/a.jpg"}}`)
	})
	defer srv.Close()

	stats, _, err := c.Instagram.ChannelStats(context.Background(), &FetchRequest{
		URL: "https://instagram.com/user",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.FullName != "User Name" {
		t.Errorf("unexpected fullName: %s", stats.FullName)
	}
	if stats.Followers == nil || *stats.Followers != 10000 {
		t.Errorf("expected followers 10000, got %v", stats.Followers)
	}
}

func TestInstagramService_ChannelPosts(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/instagram/channel-posts" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"profileUrl":"https://instagram.com/user","username":"user","items":[{"id":"p1","shortcode":"sc1","url":"https://instagram.com/p/sc1","type":"image","caption":"caption","likes":50,"comments":5,"timestamp":"2026-01-01","thumbnailUrl":"https://img.com/t.jpg","author":"user"}],"count":1,"hasMore":false}}`)
	})
	defer srv.Close()

	posts, _, err := c.Instagram.ChannelPosts(context.Background(), &ListRequest{
		URL: "https://instagram.com/user",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(posts.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(posts.Items))
	}
	if posts.Items[0].Shortcode != "sc1" {
		t.Errorf("unexpected shortcode: %s", posts.Items[0].Shortcode)
	}
}

func TestInstagramService_ChannelPosts_LimitExceeded(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not make a request")
	})
	defer srv.Close()

	limit := 101
	_, _, err := c.Instagram.ChannelPosts(context.Background(), &ListRequest{
		URL:   "https://instagram.com/user",
		Limit: &limit,
	})
	if err == nil {
		t.Fatal("expected error for limit > 100")
	}
}

func TestInstagramService_ChannelReels(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/instagram/channel-reels" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"profileUrl":"https://instagram.com/user","username":"user","items":[{"id":"r1","shortcode":"rs1","url":"https://instagram.com/reel/rs1","caption":"reel","likes":30,"comments":3,"views":1000,"duration":"00:15","thumbnailUrl":"https://img.com/t.jpg","author":"user"}],"count":1,"hasMore":false}}`)
	})
	defer srv.Close()

	reels, _, err := c.Instagram.ChannelReels(context.Background(), &ListRequest{
		URL: "https://instagram.com/user",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reels.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(reels.Items))
	}
}

func TestInstagramService_ReelsSearch(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/instagram/reels-search" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"query":"dance","items":[{"id":"r1","url":"https://instagram.com/reel/r1","caption":"dance","likes":100,"comments":5,"views":5000,"author":"user"}],"count":1,"hasMore":false}}`)
	})
	defer srv.Close()

	results, _, err := c.Instagram.ReelsSearch(context.Background(), &ReelsSearchRequest{
		Query: "dance",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(results.Items))
	}
}

func TestInstagramService_ReelsSearch_PageValidation(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not make a request for page > 1")
	})
	defer srv.Close()

	page := 2
	_, _, err := c.Instagram.ReelsSearch(context.Background(), &ReelsSearchRequest{
		Query: "dance",
		Page:  &page,
	})
	if err == nil {
		t.Fatal("expected error for page > 1")
	}
}

func TestInstagramService_Download(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/instagram/download" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"url":"https://instagram.com/reel/abc","title":"Reel","duration":"00:30","downloadUrl":"https://download.socialkit.dev/file.mp4","fileSizeMB":"5.0","format":"mp4","quality":"720p","expiresIn":"1h","creditsUsed":1}}`)
	})
	defer srv.Close()

	result, _, err := c.Instagram.Download(context.Background(), &DownloadRequest{
		URL: "https://instagram.com/reel/abc",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.DownloadURL == "" {
		t.Error("expected non-empty downloadUrl")
	}
}

func TestInstagramService_BulkTranscript(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/instagram/transcript/bulk" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":[]}`)
	})
	defer srv.Close()

	_, _, err := c.Instagram.BulkTranscript(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInstagramService_BulkStats(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/instagram/stats/bulk" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":[]}`)
	})
	defer srv.Close()

	_, _, err := c.Instagram.BulkStats(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInstagramService_BulkChannelStats(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/instagram/channel-stats/bulk" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":[]}`)
	})
	defer srv.Close()

	_, _, err := c.Instagram.BulkChannelStats(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInstagramService_BulkSummarize(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/instagram/summarize/bulk" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":[]}`)
	})
	defer srv.Close()

	_, _, err := c.Instagram.BulkSummarize(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
