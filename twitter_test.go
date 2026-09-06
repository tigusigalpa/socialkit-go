package socialkit

import (
	"context"
	"net/http"
	"testing"
)

func TestTwitterService_Profile(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/twitter/profile" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		writeResponse(t, w, `{"success":true,"data":{"url":"https://twitter.com/user","id":"123","name":"User","username":"user","bio":"bio","followers":10000,"following":500,"tweets":2000,"verified":true,"profileImage":"https://img.com/p.jpg","bannerImage":"https://img.com/b.jpg","location":"Earth","website":"https://example.com","joinedAt":"2020-01-01"}}`)
	})
	defer srv.Close()

	profile, _, err := c.Twitter.Profile(context.Background(), &FetchRequest{
		URL: "https://twitter.com/user",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if profile.Username != "user" {
		t.Errorf("unexpected username: %s", profile.Username)
	}
	if profile.Followers == nil || *profile.Followers != 10000 {
		t.Errorf("expected followers 10000, got %v", profile.Followers)
	}
}

func TestTwitterService_Tweets(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/twitter/tweets" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"tweets":[{"id":"t1","text":"Hello world","likes":100,"retweets":10,"replies":5,"views":1000,"createdAt":"2026-01-01","author":{"username":"user"},"hashtags":["#tech"],"urls":["https://example.com"]}],"nextCursor":"cursor123"}}`)
	})
	defer srv.Close()

	tweets, _, err := c.Twitter.Tweets(context.Background(), &TweetsRequest{
		URL: "https://twitter.com/user",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tweets.Tweets) != 1 {
		t.Fatalf("expected 1 tweet, got %d", len(tweets.Tweets))
	}
	if tweets.NextCursor == nil || *tweets.NextCursor != "cursor123" {
		t.Errorf("expected nextCursor cursor123, got %v", tweets.NextCursor)
	}
}

func TestTwitterService_Tweets_NullNextCursor(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeResponse(t, w, `{"success":true,"data":{"tweets":[{"id":"t1","text":"Hello"}],"nextCursor":null}}`)
	})
	defer srv.Close()

	tweets, _, err := c.Twitter.Tweets(context.Background(), &TweetsRequest{
		URL: "https://twitter.com/user",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tweets.NextCursor != nil {
		t.Errorf("expected nil nextCursor, got %v", *tweets.NextCursor)
	}
}

func TestTwitterService_Tweet(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/twitter/tweet" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"id":"t1","text":"Hello world","likes":100,"retweets":10,"replies":5,"views":1000,"createdAt":"2026-01-01","author":{"username":"user"},"hashtags":["#tech"],"urls":["https://example.com"]}}`)
	})
	defer srv.Close()

	tweet, _, err := c.Twitter.Tweet(context.Background(), &FetchRequest{
		URL: "https://twitter.com/user/status/t1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tweet.ID != "t1" {
		t.Errorf("unexpected id: %s", tweet.ID)
	}
	if tweet.Likes == nil || *tweet.Likes != 100 {
		t.Errorf("expected likes 100, got %v", tweet.Likes)
	}
}

func TestTwitterService_Thread(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/twitter/thread" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"conversationId":"123","tweetCount":3,"isThread":true,"combinedText":"Tweet 1 Tweet 2 Tweet 3","author":{"username":"user"},"tweets":[{"id":"1","text":"Tweet 1"}]}}`)
	})
	defer srv.Close()

	thread, _, err := c.Twitter.Thread(context.Background(), &FetchRequest{
		URL: "https://twitter.com/user/status/123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if thread.ConversationID != "123" {
		t.Errorf("unexpected conversationId: %s", thread.ConversationID)
	}
	if thread.IsThread == nil || !*thread.IsThread {
		t.Error("expected isThread true")
	}
	if len(thread.Author) == 0 {
		t.Error("expected non-empty author")
	}
	if len(thread.Tweets) == 0 {
		t.Error("expected non-empty tweets")
	}
}

func TestTwitterService_Transcript(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/twitter/transcript" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeResponse(t, w, `{"success":true,"data":{"url":"https://twitter.com/user/status/123","videoId":"v1","transcript":"hello twitter","transcriptSegments":[{"text":"hello","start":0,"duration":1.0}],"creditsUsed":1}}`)
	})
	defer srv.Close()

	transcript, _, err := c.Twitter.Transcript(context.Background(), &FetchRequest{
		URL: "https://twitter.com/user/status/123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if transcript.Transcript != "hello twitter" {
		t.Errorf("unexpected transcript: %s", transcript.Transcript)
	}
}
