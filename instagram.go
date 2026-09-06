package socialkit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// InstagramService provides access to the /instagram/* endpoints of the
// SocialKit API, covering transcripts, summaries, stats, comments,
// channel stats, channel posts, channel reels, reels search, downloads,
// and bulk operations.
type InstagramService struct {
	client *Client
}

// Transcript retrieves the transcript for an Instagram video/reel.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/instagram-transcript-api
func (s *InstagramService) Transcript(ctx context.Context, req *FetchRequest) (*Transcript, *ResponseMeta, error) {
	var out Transcript
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/instagram/transcript", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Summarize generates an AI summary of an Instagram video/reel.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/instagram-summarize-api
func (s *InstagramService) Summarize(ctx context.Context, req *SummaryRequest) (*Summary, *ResponseMeta, error) {
	var out Summary
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/instagram/summarize", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Stats retrieves statistics for an Instagram post.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/instagram-stats-api
func (s *InstagramService) Stats(ctx context.Context, req *FetchRequest) (*InstagramStats, *ResponseMeta, error) {
	var out InstagramStats
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/instagram/stats", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Comments retrieves comments for an Instagram post with cursor-based
// pagination. The response includes hasMore and cursor fields.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/instagram-comments-api
func (s *InstagramService) Comments(ctx context.Context, req *CommentsRequest) (*InstagramComments, *ResponseMeta, error) {
	var out InstagramComments
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/instagram/comments", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// ChannelStats retrieves statistics for an Instagram channel.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/instagram-channel-stats-api
func (s *InstagramService) ChannelStats(ctx context.Context, req *FetchRequest) (*InstagramChannelStats, *ResponseMeta, error) {
	var out InstagramChannelStats
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/instagram/channel-stats", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// ChannelPosts retrieves posts from an Instagram channel with cursor-based
// pagination. Limit is capped at 100.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/instagram-channel-posts-api
func (s *InstagramService) ChannelPosts(ctx context.Context, req *ListRequest) (*InstagramChannelPosts, *ResponseMeta, error) {
	if req != nil && req.Limit != nil && *req.Limit > 100 {
		return nil, nil, fmt.Errorf("socialkit: instagram channel posts limit must be at most 100")
	}
	var out InstagramChannelPosts
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/instagram/channel-posts", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// ChannelReels retrieves reels from an Instagram channel with cursor-based
// pagination. This endpoint can take up to 60 seconds to respond.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/instagram-channel-reels-api
func (s *InstagramService) ChannelReels(ctx context.Context, req *ListRequest) (*InstagramChannelReels, *ResponseMeta, error) {
	var out InstagramChannelReels
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/instagram/channel-reels", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// ReelsSearch searches for Instagram reels. The page parameter must be 1;
// pages greater than 1 are not supported and will return an error before
// making any request.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/instagram-reels-search-api
func (s *InstagramService) ReelsSearch(ctx context.Context, req *ReelsSearchRequest) (*InstagramReelsSearch, *ResponseMeta, error) {
	if req != nil && req.Page != nil && *req.Page > 1 {
		return nil, nil, errors.New("socialkit: instagram reels search only supports page 1")
	}
	var out InstagramReelsSearch
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/instagram/reels-search", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Download initiates a download for an Instagram video/reel. The
// download URL expires after 1 hour.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/instagram-download-api
func (s *InstagramService) Download(ctx context.Context, req *DownloadRequest) (*DownloadResult, *ResponseMeta, error) {
	var out DownloadResult
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/instagram/download", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// BulkTranscript retrieves transcripts for multiple Instagram videos in
// a single request. This is an experimental endpoint; the request and
// response are passed as raw JSON.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/instagram-transcript-api and https://api.socialkit.dev/openapi.json (experimental bulk)
func (s *InstagramService) BulkTranscript(ctx context.Context, body json.RawMessage) (json.RawMessage, *ResponseMeta, error) {
	var out json.RawMessage
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/instagram/transcript/bulk", nil, body, &out)
	if err != nil {
		return nil, meta, err
	}
	return out, meta, nil
}

// BulkStats retrieves statistics for multiple Instagram posts in a
// single request. This is an experimental endpoint; the request and
// response are passed as raw JSON.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/instagram-stats-api and https://api.socialkit.dev/openapi.json (experimental bulk)
func (s *InstagramService) BulkStats(ctx context.Context, body json.RawMessage) (json.RawMessage, *ResponseMeta, error) {
	var out json.RawMessage
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/instagram/stats/bulk", nil, body, &out)
	if err != nil {
		return nil, meta, err
	}
	return out, meta, nil
}

// BulkChannelStats retrieves statistics for multiple Instagram channels
// in a single request. This is an experimental endpoint; the request and
// response are passed as raw JSON.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/instagram-channel-stats-api and https://api.socialkit.dev/openapi.json (experimental bulk)
func (s *InstagramService) BulkChannelStats(ctx context.Context, body json.RawMessage) (json.RawMessage, *ResponseMeta, error) {
	var out json.RawMessage
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/instagram/channel-stats/bulk", nil, body, &out)
	if err != nil {
		return nil, meta, err
	}
	return out, meta, nil
}

// BulkSummarize generates AI summaries for multiple Instagram videos in
// a single request. This is an experimental endpoint; the request and
// response are passed as raw JSON.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/instagram-summarize-api and https://api.socialkit.dev/openapi.json (experimental bulk)
func (s *InstagramService) BulkSummarize(ctx context.Context, body json.RawMessage) (json.RawMessage, *ResponseMeta, error) {
	var out json.RawMessage
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/instagram/summarize/bulk", nil, body, &out)
	if err != nil {
		return nil, meta, err
	}
	return out, meta, nil
}
