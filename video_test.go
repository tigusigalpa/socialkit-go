package socialkit

import (
	"context"
	"net/http"
	"testing"
)

func TestVideoService_Transcript(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/video/transcript" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		writeResponse(t, w, `{"success":true,"data":{"url":"https://example.com/video.mp4","videoId":"","transcript":"hello video","transcriptSegments":[{"text":"hello","start":0,"duration":2.0}],"durationSeconds":120.0,"creditsUsed":10}}`)
	})
	defer srv.Close()

	transcript, _, err := c.Video.Transcript(context.Background(), &FetchRequest{
		URL: "https://example.com/video.mp4",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if transcript.Transcript != "hello video" {
		t.Errorf("unexpected transcript: %s", transcript.Transcript)
	}
	if transcript.DurationSeconds == nil || *transcript.DurationSeconds != 120.0 {
		t.Errorf("expected durationSeconds 120, got %v", transcript.DurationSeconds)
	}
	if transcript.CreditsUsed == nil || *transcript.CreditsUsed != 10 {
		t.Errorf("expected creditsUsed 10, got %v", transcript.CreditsUsed)
	}
}

func TestVideoService_Summarize(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/video/summarize" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"url":"https://example.com/video.mp4","summary":"A video summary","mainTopics":["topic"],"keyPoints":["point"],"tone":"neutral","targetAudience":"all","quotes":[],"durationSeconds":120.0,"creditsUsed":10}}`)
	})
	defer srv.Close()

	summary, _, err := c.Video.Summarize(context.Background(), &SummaryRequest{
		URL: "https://example.com/video.mp4",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Summary != "A video summary" {
		t.Errorf("unexpected summary: %s", summary.Summary)
	}
	if summary.CreditsUsed == nil || *summary.CreditsUsed != 10 {
		t.Errorf("expected creditsUsed 10, got %v", summary.CreditsUsed)
	}
}
