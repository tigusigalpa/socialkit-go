package socialkit

import (
	"context"
	"net/http"
	"testing"
)

func TestLinkedInService_Profile(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/linkedin/profile" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		writeResponse(t, w, `{"success":true,"data":{"url":"https://linkedin.com/in/user","name":"User Name","headline":"Software Engineer","followers":5000,"connections":"500+","profileImage":"https://img.com/p.jpg","recentArticles":[{"title":"Article 1","url":"https://linkedin.com/article/1","date":"2026-01-01"}]}}`)
	})
	defer srv.Close()

	profile, _, err := c.LinkedIn.Profile(context.Background(), &FetchRequest{
		URL: "https://linkedin.com/in/user",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if profile.Name != "User Name" {
		t.Errorf("unexpected name: %s", profile.Name)
	}
	if len(profile.RecentArticles) != 1 {
		t.Fatalf("expected 1 article, got %d", len(profile.RecentArticles))
	}
	if profile.RecentArticles[0].Title != "Article 1" {
		t.Errorf("unexpected article title: %s", profile.RecentArticles[0].Title)
	}
}

func TestLinkedInService_Company(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/linkedin/company" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"url":"https://linkedin.com/company/test","name":"Test Co","description":"A company","followers":50000,"employees":"100-200","logo":"https://img.com/logo.jpg"}}`)
	})
	defer srv.Close()

	company, _, err := c.LinkedIn.Company(context.Background(), &FetchRequest{
		URL: "https://linkedin.com/company/test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if company.Name != "Test Co" {
		t.Errorf("unexpected name: %s", company.Name)
	}
	if company.Followers == nil || *company.Followers != 50000 {
		t.Errorf("expected followers 50000, got %v", company.Followers)
	}
}

func TestLinkedInService_CompanyPosts(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/linkedin/company-posts" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"posts":[{"id":"p1","text":"Post text","likes":100,"comments":10,"shares":5,"imageUrl":"https://img.com/img.jpg","publishedAt":"2026-01-01","author":{"name":"Test Co"}}]}}`)
	})
	defer srv.Close()

	posts, _, err := c.LinkedIn.CompanyPosts(context.Background(), &CompanyPostsRequest{
		URL: "https://linkedin.com/company/test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(posts.Posts) != 1 {
		t.Fatalf("expected 1 post, got %d", len(posts.Posts))
	}
	if posts.Posts[0].ID != "p1" {
		t.Errorf("unexpected post id: %s", posts.Posts[0].ID)
	}
}

func TestLinkedInService_CompanyPosts_LimitExceeded(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not make a request")
	})
	defer srv.Close()

	limit := 51
	_, _, err := c.LinkedIn.CompanyPosts(context.Background(), &CompanyPostsRequest{
		URL:   "https://linkedin.com/company/test",
		Limit: &limit,
	})
	if err == nil {
		t.Fatal("expected error for limit > 50")
	}
}

func TestLinkedInService_Post(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/linkedin/post" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"post":{"id":"p1","text":"Post text","author":"User","likes":50}}}`)
	})
	defer srv.Close()

	post, _, err := c.LinkedIn.Post(context.Background(), &FetchRequest{
		URL: "https://linkedin.com/posts/p1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(post.Post) == 0 {
		t.Error("expected non-empty post")
	}
}

func TestLinkedInService_Transcript(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/linkedin/transcript" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"transcript":"hello linkedin","transcriptSegments":[{"text":"hello","start":0,"duration":1.0}],"wordCount":1,"durationSeconds":5.0,"creditsUsed":1}}`)
	})
	defer srv.Close()

	transcript, _, err := c.LinkedIn.Transcript(context.Background(), &FetchRequest{
		URL: "https://linkedin.com/posts/p1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if transcript.Transcript != "hello linkedin" {
		t.Errorf("unexpected transcript: %s", transcript.Transcript)
	}
	// LinkedIn transcript segments should NOT have a timestamp field.
	if len(transcript.TranscriptSegments) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(transcript.TranscriptSegments))
	}
	if transcript.TranscriptSegments[0].Text != "hello" {
		t.Errorf("unexpected segment text: %s", transcript.TranscriptSegments[0].Text)
	}
}
