package socialkit

import (
	"context"
	"encoding/json"
	"fmt"
)

// TikTokService provides access to the /tiktok/* endpoints of the
// SocialKit API, covering transcripts, summaries, stats, comments,
// channel stats, channel videos, search, hashtag search, downloads, and
// bulk operations.
type TikTokService struct {
	client *Client
}

// Transcript retrieves the transcript for a TikTok video.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/tiktok-transcript-api
func (s *TikTokService) Transcript(ctx context.Context, req *FetchRequest) (*Transcript, *ResponseMeta, error) {
	var out Transcript
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/tiktok/transcript", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Summarize generates an AI summary of a TikTok video.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/tiktok-summarize-api
func (s *TikTokService) Summarize(ctx context.Context, req *SummaryRequest) (*Summary, *ResponseMeta, error) {
	var out Summary
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/tiktok/summarize", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Stats retrieves statistics for a TikTok video.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/tiktok-stats-api
func (s *TikTokService) Stats(ctx context.Context, req *FetchRequest) (*TikTokStats, *ResponseMeta, error) {
	var out TikTokStats
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/tiktok/stats", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Comments retrieves comments for a TikTok video with cursor-based
// pagination. The response includes hasMore and cursor fields.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/tiktok-comments-api
func (s *TikTokService) Comments(ctx context.Context, req *CommentsRequest) (*TikTokComments, *ResponseMeta, error) {
	var out TikTokComments
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/tiktok/comments", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// ChannelStats retrieves statistics for a TikTok channel.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/tiktok-channel-stats-api
func (s *TikTokService) ChannelStats(ctx context.Context, req *FetchRequest) (*TikTokChannelStats, *ResponseMeta, error) {
	var out TikTokChannelStats
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/tiktok/channel-stats", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// ChannelVideos retrieves videos from a TikTok channel with cursor-based
// pagination. The default limit is 30 and the maximum is 100.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/tiktok-channel-videos-api
func (s *TikTokService) ChannelVideos(ctx context.Context, req *ListRequest) (*TikTokChannelVideos, *ResponseMeta, error) {
	if req != nil && req.Limit != nil && *req.Limit > 100 {
		return nil, nil, fmt.Errorf("socialkit: tiktok channel videos limit must be at most 100")
	}
	var out TikTokChannelVideos
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/tiktok/channel-videos", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Search searches for TikTok videos. SortBy accepts "relevance",
// "likes", or "date". DatePosted accepts "day", "week", "month",
// "3months", or "6months". Limit is capped at 100.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/tiktok-search-api
func (s *TikTokService) Search(ctx context.Context, req *SearchRequest) (*TikTokSearchResult, *ResponseMeta, error) {
	if req != nil && req.Limit != nil && *req.Limit > 100 {
		return nil, nil, fmt.Errorf("socialkit: tiktok search limit must be at most 100")
	}
	var out TikTokSearchResult
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/tiktok/search", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// HashtagSearch searches for TikTok videos by hashtag. The hashtag
// should be provided without the leading '#'. Limit is capped at 100.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/tiktok-hashtag-search-api
func (s *TikTokService) HashtagSearch(ctx context.Context, req *HashtagSearchRequest) (*TikTokHashtagResult, *ResponseMeta, error) {
	if req != nil && req.Limit != nil && *req.Limit > 100 {
		return nil, nil, fmt.Errorf("socialkit: tiktok hashtag search limit must be at most 100")
	}
	var out TikTokHashtagResult
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/tiktok/hashtag-search", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Download initiates a download for a TikTok video. The download URL
// expires after 1 hour. Maximum file size is 30MB.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/tiktok-download-api
func (s *TikTokService) Download(ctx context.Context, req *DownloadRequest) (*DownloadResult, *ResponseMeta, error) {
	var out DownloadResult
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/tiktok/download", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// BulkTranscript retrieves transcripts for multiple TikTok videos in a
// single request. This is an experimental endpoint; the request and
// response are passed as raw JSON.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/tiktok-transcript-api and https://api.socialkit.dev/openapi.json (experimental bulk)
func (s *TikTokService) BulkTranscript(ctx context.Context, body json.RawMessage) (json.RawMessage, *ResponseMeta, error) {
	var out json.RawMessage
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/tiktok/transcript/bulk", nil, body, &out)
	if err != nil {
		return nil, meta, err
	}
	return out, meta, nil
}

// BulkComments retrieves comments for multiple TikTok videos in a
// single request. This is an experimental endpoint; the request and
// response are passed as raw JSON.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/tiktok-comments-api and https://api.socialkit.dev/openapi.json (experimental bulk)
func (s *TikTokService) BulkComments(ctx context.Context, body json.RawMessage) (json.RawMessage, *ResponseMeta, error) {
	var out json.RawMessage
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/tiktok/comments/bulk", nil, body, &out)
	if err != nil {
		return nil, meta, err
	}
	return out, meta, nil
}

// BulkStats retrieves statistics for multiple TikTok videos in a single
// request. This is an experimental endpoint; the request and response
// are passed as raw JSON.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/tiktok-stats-api and https://api.socialkit.dev/openapi.json (experimental bulk)
func (s *TikTokService) BulkStats(ctx context.Context, body json.RawMessage) (json.RawMessage, *ResponseMeta, error) {
	var out json.RawMessage
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/tiktok/stats/bulk", nil, body, &out)
	if err != nil {
		return nil, meta, err
	}
	return out, meta, nil
}

// BulkChannelStats retrieves statistics for multiple TikTok channels in
// a single request. This is an experimental endpoint; the request and
// response are passed as raw JSON.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/tiktok-channel-stats-api and https://api.socialkit.dev/openapi.json (experimental bulk)
func (s *TikTokService) BulkChannelStats(ctx context.Context, body json.RawMessage) (json.RawMessage, *ResponseMeta, error) {
	var out json.RawMessage
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/tiktok/channel-stats/bulk", nil, body, &out)
	if err != nil {
		return nil, meta, err
	}
	return out, meta, nil
}

// BulkSummarize generates AI summaries for multiple TikTok videos in a
// single request. This is an experimental endpoint; the request and
// response are passed as raw JSON.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/tiktok-summarize-api and https://api.socialkit.dev/openapi.json (experimental bulk)
func (s *TikTokService) BulkSummarize(ctx context.Context, body json.RawMessage) (json.RawMessage, *ResponseMeta, error) {
	var out json.RawMessage
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/tiktok/summarize/bulk", nil, body, &out)
	if err != nil {
		return nil, meta, err
	}
	return out, meta, nil
}
