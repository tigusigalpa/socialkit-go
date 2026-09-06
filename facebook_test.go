package socialkit

import (
	"context"
	"net/http"
	"testing"
)

func TestFacebookService_Transcript(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/facebook/transcript" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		writeResponse(t, w, `{"success":true,"data":{"url":"https://facebook.com/watch?v=123","videoId":"123","transcript":"hello fb","transcriptSegments":[{"text":"hello","start":0,"duration":1.0}],"creditsUsed":1}}`)
	})
	defer srv.Close()

	transcript, _, err := c.Facebook.Transcript(context.Background(), &FetchRequest{
		URL: "https://facebook.com/watch?v=123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if transcript.Transcript != "hello fb" {
		t.Errorf("unexpected transcript: %s", transcript.Transcript)
	}
}

func TestFacebookService_Summarize(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeResponse(t, w, `{"success":true,"data":{"url":"https://facebook.com/watch?v=123","summary":"A FB video","mainTopics":["topic"],"keyPoints":["point"],"tone":"neutral","targetAudience":"all","quotes":[]}}`)
	})
	defer srv.Close()

	summary, _, err := c.Facebook.Summarize(context.Background(), &SummaryRequest{
		URL: "https://facebook.com/watch?v=123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Summary != "A FB video" {
		t.Errorf("unexpected summary: %s", summary.Summary)
	}
}

func TestFacebookService_Stats(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/facebook/stats" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"postUrl":"https://facebook.com/post/123","id":"123","description":"desc","views":1000,"likes":100,"comments":50,"shares":25,"reactions":[{"name":"like","count":80,"formatted":"80"},{"name":"love","count":20,"formatted":"20"}],"author":"Page","authorLink":"https://facebook.com/page","isVideo":true}}`)
	})
	defer srv.Close()

	stats, _, err := c.Facebook.Stats(context.Background(), &FetchRequest{
		URL: "https://facebook.com/post/123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stats.Reactions) != 2 {
		t.Fatalf("expected 2 reactions, got %d", len(stats.Reactions))
	}
	if stats.Reactions[0].Name != "like" {
		t.Errorf("expected reaction like, got %s", stats.Reactions[0].Name)
	}
	if stats.Reactions[0].Count == nil || *stats.Reactions[0].Count != 80 {
		t.Errorf("expected count 80, got %v", stats.Reactions[0].Count)
	}
	if stats.IsVideo == nil || !*stats.IsVideo {
		t.Error("expected isVideo true")
	}
}

func TestFacebookService_Comments(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeResponse(t, w, `{"success":true,"data":{"postUrl":"https://facebook.com/post/123","comments":[{"id":"c1","text":"Nice!","likeCount":3,"replyCount":0,"user":{"name":"User1"}}],"commentCount":null,"hasMore":false}}`)
	})
	defer srv.Close()

	comments, _, err := c.Facebook.Comments(context.Background(), &CommentsRequest{
		URL: "https://facebook.com/post/123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments.Comments) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(comments.Comments))
	}
	if comments.CommentCount != nil {
		t.Errorf("expected nil commentCount, got %v", *comments.CommentCount)
	}
}

func TestFacebookService_ChannelStats(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/facebook/channel-stats" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"profileUrl":"https://facebook.com/page","id":"123","fullName":"Test Page","bio":"page bio","avatar":"https://img.com/a.jpg","coverPhoto":"https://img.com/c.jpg","verified":true,"followers":5000,"category":"Brand"}}`)
	})
	defer srv.Close()

	stats, _, err := c.Facebook.ChannelStats(context.Background(), &FetchRequest{
		URL: "https://facebook.com/page",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.FullName != "Test Page" {
		t.Errorf("unexpected fullName: %s", stats.FullName)
	}
	if stats.Category != "Brand" {
		t.Errorf("unexpected category: %s", stats.Category)
	}
}
