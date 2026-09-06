package socialkit

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestYouTubeService_Transcript(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/youtube/transcript" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body FetchRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body.URL != "https://youtube.com/watch?v=test" {
			t.Errorf("unexpected url: %s", body.URL)
		}
		writeResponse(t, w, `{"success":true,"data":{"url":"https://youtube.com/watch?v=test","videoId":"test","transcript":"hello world","transcriptSegments":[{"text":"hello","start":0,"duration":1.5,"timestamp":"00:00:00"}],"wordCount":2,"segments":1,"language":"en","durationSeconds":10.5,"creditsUsed":1}}`)
	})
	defer srv.Close()

	transcript, _, err := c.YouTube.Transcript(context.Background(), &FetchRequest{
		URL: "https://youtube.com/watch?v=test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if transcript.Transcript != "hello world" {
		t.Errorf("expected transcript 'hello world', got %s", transcript.Transcript)
	}
	if transcript.VideoID != "test" {
		t.Errorf("expected videoId test, got %s", transcript.VideoID)
	}
	if len(transcript.TranscriptSegments) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(transcript.TranscriptSegments))
	}
	if transcript.TranscriptSegments[0].Text != "hello" {
		t.Errorf("expected segment text hello, got %s", transcript.TranscriptSegments[0].Text)
	}
	if transcript.TranscriptSegments[0].Timestamp == nil || *transcript.TranscriptSegments[0].Timestamp != "00:00:00" {
		t.Errorf("expected timestamp 00:00:00, got %v", transcript.TranscriptSegments[0].Timestamp)
	}
	if transcript.WordCount == nil || *transcript.WordCount != 2 {
		t.Errorf("expected wordCount 2, got %v", transcript.WordCount)
	}
	if transcript.DurationSeconds == nil || *transcript.DurationSeconds != 10.5 {
		t.Errorf("expected durationSeconds 10.5, got %v", transcript.DurationSeconds)
	}
}

func TestYouTubeService_Summarize(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/youtube/summarize" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"url":"https://youtube.com/watch?v=test","summary":"A great video","mainTopics":["tech","AI"],"keyPoints":["point1","point2"],"tone":"informative","targetAudience":"developers","quotes":["quote1"],"durationSeconds":120.0,"creditsUsed":3}}`)
	})
	defer srv.Close()

	summary, _, err := c.YouTube.Summarize(context.Background(), &SummaryRequest{
		URL: "https://youtube.com/watch?v=test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Summary != "A great video" {
		t.Errorf("unexpected summary: %s", summary.Summary)
	}
	if len(summary.MainTopics) != 2 {
		t.Errorf("expected 2 main topics, got %d", len(summary.MainTopics))
	}
	if len(summary.KeyPoints) != 2 {
		t.Errorf("expected 2 key points, got %d", len(summary.KeyPoints))
	}
	if summary.Tone != "informative" {
		t.Errorf("unexpected tone: %s", summary.Tone)
	}
	if summary.TargetAudience != "developers" {
		t.Errorf("unexpected targetAudience: %s", summary.TargetAudience)
	}
}

func TestYouTubeService_Summarize_CustomResponse(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeResponse(t, w, `{"success":true,"data":{"url":"https://youtube.com/watch?v=test","summary":"A great video","customField":"customValue","anotherField":42}}`)
	})
	defer srv.Close()

	summary, _, err := c.YouTube.Summarize(context.Background(), &SummaryRequest{
		URL: "https://youtube.com/watch?v=test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Extra == nil {
		t.Fatal("expected non-nil Extra map")
	}
	if _, ok := summary.Extra["customField"]; !ok {
		t.Error("expected customField in Extra map")
	}
	if _, ok := summary.Extra["anotherField"]; !ok {
		t.Error("expected anotherField in Extra map")
	}
	if _, ok := summary.Extra["url"]; ok {
		t.Error("url should not be in Extra map")
	}
}

func TestYouTubeService_Stats(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/youtube/stats" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"url":"https://youtube.com/watch?v=test","videoId":"test","title":"Test Video","description":"A test","channelName":"TestChannel","channelLink":"https://youtube.com/c/test","channelHandle":"@test","views":1000,"likes":50,"comments":10,"duration":"10:30","publishedAt":"2026-01-01","thumbnailUrl":"https://img.youtube.com/vi/test/maxresdefault.jpg","contentType":"video","isShortForm":false}}`)
	})
	defer srv.Close()

	stats, _, err := c.YouTube.Stats(context.Background(), &FetchRequest{
		URL: "https://youtube.com/watch?v=test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Title != "Test Video" {
		t.Errorf("unexpected title: %s", stats.Title)
	}
	if stats.Views == nil || *stats.Views != 1000 {
		t.Errorf("expected views 1000, got %v", stats.Views)
	}
	if stats.ChannelHandle != "@test" {
		t.Errorf("expected channelHandle @test, got %s", stats.ChannelHandle)
	}
	if stats.IsShortForm == nil || *stats.IsShortForm != false {
		t.Errorf("expected isShortForm false, got %v", stats.IsShortForm)
	}
}

func TestYouTubeService_Comments(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/youtube/comments" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"url":"https://youtube.com/watch?v=test","comments":[{"author":"User1","text":"Great video","likes":5,"date":"2026-01-01","avatar":"https://img.com/avatar1.jpg","replyCount":2,"position":0}]}}`)
	})
	defer srv.Close()

	sortBy := "top"
	limit := 50
	comments, _, err := c.YouTube.Comments(context.Background(), &CommentsRequest{
		URL:    "https://youtube.com/watch?v=test",
		Limit:  &limit,
		SortBy: &sortBy,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments.Comments) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(comments.Comments))
	}
	if comments.Comments[0].Author != "User1" {
		t.Errorf("unexpected author: %s", comments.Comments[0].Author)
	}
	if comments.Comments[0].Likes == nil || *comments.Comments[0].Likes != 5 {
		t.Errorf("expected likes 5, got %v", comments.Comments[0].Likes)
	}
}

func TestYouTubeService_Comments_LimitExceeded(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not make a request")
	})
	defer srv.Close()

	limit := 101
	_, _, err := c.YouTube.Comments(context.Background(), &CommentsRequest{
		URL:   "https://youtube.com/watch?v=test",
		Limit: &limit,
	})
	if err == nil {
		t.Fatal("expected error for limit > 100")
	}
}

func TestYouTubeService_ChannelStats(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/youtube/channel-stats" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"profileUrl":"https://youtube.com/c/test","username":"test","nickname":"Test Channel","bio":"A test channel","verified":true,"avatar":"https://img.com/avatar.jpg","subscribers":10000,"totalVideos":50,"bioLink":"https://example.com","banner":"https://img.com/banner.jpg"}}`)
	})
	defer srv.Close()

	stats, _, err := c.YouTube.ChannelStats(context.Background(), &FetchRequest{
		URL: "https://youtube.com/c/test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Nickname != "Test Channel" {
		t.Errorf("unexpected nickname: %s", stats.Nickname)
	}
	if stats.Verified == nil || !*stats.Verified {
		t.Error("expected verified true")
	}
	if stats.Subscribers == nil || *stats.Subscribers != 10000 {
		t.Errorf("expected subscribers 10000, got %v", stats.Subscribers)
	}
}

func TestYouTubeService_Search(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/youtube/search" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"query":"golang tutorial","results":[{"videoId":"abc","title":"Go Tutorial","thumbnail":"https://img.com/thumb.jpg","channelName":"Go Channel","channelId":"UC123","channelUrl":"https://youtube.com/c/go","publishedTime":"1 week ago","duration":"10:00","views":5000,"viewsFormatted":"5K","description":"Learn Go","url":"https://youtube.com/watch?v=abc","verified":true}]}}`)
	})
	defer srv.Close()

	sortBy := "relevance"
	results, _, err := c.YouTube.Search(context.Background(), &SearchRequest{
		Query:  "golang tutorial",
		SortBy: &sortBy,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results.Query != "golang tutorial" {
		t.Errorf("unexpected query: %s", results.Query)
	}
	if len(results.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results.Results))
	}
	if results.Results[0].VideoID != "abc" {
		t.Errorf("unexpected videoId: %s", results.Results[0].VideoID)
	}
}

func TestYouTubeService_Videos(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/youtube/videos" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"url":"https://youtube.com/c/test","type":"channel","results":[{"videoId":"v1","title":"Video 1","thumbnail":"https://img.com/t1.jpg","channelName":"Test","duration":"5:00","views":100,"viewsFormatted":"100","publishedTime":"1 day ago","publishedAt":"2026-01-01","description":"desc","url":"https://youtube.com/watch?v=v1"}]}}`)
	})
	defer srv.Close()

	fullDetails := true
	videos, _, err := c.YouTube.Videos(context.Background(), &VideosRequest{
		URL:         "https://youtube.com/c/test",
		FullDetails: &fullDetails,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if videos.Type != "channel" {
		t.Errorf("unexpected type: %s", videos.Type)
	}
	if len(videos.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(videos.Results))
	}
}

func TestYouTubeService_Videos_LimitExceeded(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not make a request")
	})
	defer srv.Close()

	limit := 31
	_, _, err := c.YouTube.Videos(context.Background(), &VideosRequest{
		URL:   "https://youtube.com/c/test",
		Limit: &limit,
	})
	if err == nil {
		t.Fatal("expected error for limit > 30")
	}
}

func TestYouTubeService_Download(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/youtube/download" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"url":"https://youtube.com/watch?v=test","title":"Test Video","duration":"10:30","durationSeconds":630.0,"downloadUrl":"https://download.socialkit.dev/file.mp4","fileSize":10485760,"fileSizeMB":"10.0","format":"mp4","quality":"720p","expiresIn":"1h","creditsUsed":2}}`)
	})
	defer srv.Close()

	format := "mp4"
	quality := "720p"
	result, _, err := c.YouTube.Download(context.Background(), &DownloadRequest{
		URL:     "https://youtube.com/watch?v=test",
		Format:  &format,
		Quality: &quality,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.DownloadURL != "https://download.socialkit.dev/file.mp4" {
		t.Errorf("unexpected downloadUrl: %s", result.DownloadURL)
	}
	if result.Format != "mp4" {
		t.Errorf("unexpected format: %s", result.Format)
	}
}

func TestYouTubeService_BulkTranscript(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/youtube/transcript/bulk" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":[{"url":"https://youtube.com/watch?v=1","transcript":"first"},{"url":"https://youtube.com/watch?v=2","transcript":"second"}]}`)
	})
	defer srv.Close()

	body := json.RawMessage(`{"urls":["https://youtube.com/watch?v=1","https://youtube.com/watch?v=2"]}`)
	result, _, err := c.YouTube.BulkTranscript(context.Background(), body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) == 0 {
		t.Error("expected non-empty result")
	}
}

func TestYouTubeService_BulkComments(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/youtube/comments/bulk" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":[]}`)
	})
	defer srv.Close()

	_, _, err := c.YouTube.BulkComments(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestYouTubeService_BulkStats(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/youtube/stats/bulk" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":[]}`)
	})
	defer srv.Close()

	_, _, err := c.YouTube.BulkStats(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestYouTubeService_BulkSummarize(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/youtube/summarize/bulk" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":[]}`)
	})
	defer srv.Close()

	_, _, err := c.YouTube.BulkSummarize(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
