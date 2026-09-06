package socialkit

import (
	"context"
)

// TwitterService provides access to the /twitter/* endpoints of the
// SocialKit API, covering profile, tweets, individual tweet, thread, and
// transcript operations.
type TwitterService struct {
	client *Client
}

// Profile retrieves a Twitter/X user profile.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/twitter-profile-api
func (s *TwitterService) Profile(ctx context.Context, req *FetchRequest) (*TwitterProfile, *ResponseMeta, error) {
	var out TwitterProfile
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/twitter/profile", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Tweets retrieves tweets from a Twitter/X user timeline. The
// nextCursor field may be null when there are no more pages.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/twitter-tweets-api
func (s *TwitterService) Tweets(ctx context.Context, req *TweetsRequest) (*TwitterTweets, *ResponseMeta, error) {
	var out TwitterTweets
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/twitter/tweets", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Tweet retrieves a single tweet by URL.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/twitter-tweet-api
func (s *TwitterService) Tweet(ctx context.Context, req *FetchRequest) (*TwitterTweet, *ResponseMeta, error) {
	var out TwitterTweet
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/twitter/tweet", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Thread retrieves a Twitter/X thread by URL.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/twitter-thread-api
func (s *TwitterService) Thread(ctx context.Context, req *FetchRequest) (*TwitterThread, *ResponseMeta, error) {
	var out TwitterThread
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/twitter/thread", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Transcript retrieves the transcript for a Twitter/X video.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/twitter-transcript-api
func (s *TwitterService) Transcript(ctx context.Context, req *FetchRequest) (*TwitterTranscript, *ResponseMeta, error) {
	var out TwitterTranscript
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/twitter/transcript", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}
