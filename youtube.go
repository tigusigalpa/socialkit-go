package socialkit

import (
	"context"
	"encoding/json"
	"fmt"
)

// YouTubeService provides access to the /youtube/* endpoints of the
// SocialKit API, covering transcripts, summaries, stats, comments,
// channel stats, search, videos, downloads, and bulk operations.
type YouTubeService struct {
	client *Client
}

// Transcript retrieves the transcript for a YouTube video.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/youtube-transcript-api
func (s *YouTubeService) Transcript(ctx context.Context, req *FetchRequest) (*Transcript, *ResponseMeta, error) {
	var out Transcript
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/youtube/transcript", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Summarize generates an AI summary of a YouTube video.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/youtube-summarize-api
func (s *YouTubeService) Summarize(ctx context.Context, req *SummaryRequest) (*Summary, *ResponseMeta, error) {
	var out Summary
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/youtube/summarize", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Stats retrieves statistics for a YouTube video.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/youtube-stats-api
func (s *YouTubeService) Stats(ctx context.Context, req *FetchRequest) (*YouTubeStats, *ResponseMeta, error) {
	var out YouTubeStats
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/youtube/stats", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Comments retrieves comments for a YouTube video. Limit is capped at
// 100. SortBy accepts "top" or "new".
//
// Upstream docs: https://docs.socialkit.dev/api-reference/youtube-comments-api
func (s *YouTubeService) Comments(ctx context.Context, req *CommentsRequest) (*YouTubeComments, *ResponseMeta, error) {
	if req != nil && req.Limit != nil && *req.Limit > 100 {
		return nil, nil, fmt.Errorf("socialkit: youtube comments limit must be at most 100")
	}
	var out YouTubeComments
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/youtube/comments", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// ChannelStats retrieves statistics for a YouTube channel.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/youtube-channel-stats-api
func (s *YouTubeService) ChannelStats(ctx context.Context, req *FetchRequest) (*YouTubeChannelStats, *ResponseMeta, error) {
	var out YouTubeChannelStats
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/youtube/channel-stats", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Search searches for YouTube videos. SortBy accepts "relevance",
// "date", "views", or "rating". UploadDate accepts "hour", "today",
// "week", "month", or "year". Type accepts "video" or "shorts".
//
// Upstream docs: https://docs.socialkit.dev/api-reference/youtube-search-api
func (s *YouTubeService) Search(ctx context.Context, req *SearchRequest) (*YouTubeSearchResult, *ResponseMeta, error) {
	var out YouTubeSearchResult
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/youtube/search", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Videos retrieves a list of videos from a YouTube channel or playlist.
// FullDetails controls whether extended metadata is returned. Limit is
// capped at 30.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/youtube-videos-api
func (s *YouTubeService) Videos(ctx context.Context, req *VideosRequest) (*YouTubeVideos, *ResponseMeta, error) {
	if req != nil && req.Limit != nil && *req.Limit > 30 {
		return nil, nil, fmt.Errorf("socialkit: youtube videos limit must be at most 30")
	}
	var out YouTubeVideos
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/youtube/videos", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Download initiates a download for a YouTube video. The download URL
// expires after 1 hour. Maximum file size is 10MB.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/youtube-download-api
func (s *YouTubeService) Download(ctx context.Context, req *DownloadRequest) (*DownloadResult, *ResponseMeta, error) {
	var out DownloadResult
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/youtube/download", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// BulkTranscript retrieves transcripts for multiple YouTube videos in a
// single request. This is an experimental endpoint; the request and
// response are passed as raw JSON.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/youtube-transcript-api and https://api.socialkit.dev/openapi.json (experimental bulk)
func (s *YouTubeService) BulkTranscript(ctx context.Context, body json.RawMessage) (json.RawMessage, *ResponseMeta, error) {
	var out json.RawMessage
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/youtube/transcript/bulk", nil, body, &out)
	if err != nil {
		return nil, meta, err
	}
	return out, meta, nil
}

// BulkComments retrieves comments for multiple YouTube videos in a
// single request. This is an experimental endpoint; the request and
// response are passed as raw JSON.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/youtube-comments-api and https://api.socialkit.dev/openapi.json (experimental bulk)
func (s *YouTubeService) BulkComments(ctx context.Context, body json.RawMessage) (json.RawMessage, *ResponseMeta, error) {
	var out json.RawMessage
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/youtube/comments/bulk", nil, body, &out)
	if err != nil {
		return nil, meta, err
	}
	return out, meta, nil
}

// BulkStats retrieves statistics for multiple YouTube videos in a
// single request. This is an experimental endpoint; the request and
// response are passed as raw JSON.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/youtube-stats-api and https://api.socialkit.dev/openapi.json (experimental bulk)
func (s *YouTubeService) BulkStats(ctx context.Context, body json.RawMessage) (json.RawMessage, *ResponseMeta, error) {
	var out json.RawMessage
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/youtube/stats/bulk", nil, body, &out)
	if err != nil {
		return nil, meta, err
	}
	return out, meta, nil
}

// BulkSummarize generates AI summaries for multiple YouTube videos in a
// single request. This is an experimental endpoint; the request and
// response are passed as raw JSON.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/youtube-summarize-api and https://api.socialkit.dev/openapi.json (experimental bulk)
func (s *YouTubeService) BulkSummarize(ctx context.Context, body json.RawMessage) (json.RawMessage, *ResponseMeta, error) {
	var out json.RawMessage
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/youtube/summarize/bulk", nil, body, &out)
	if err != nil {
		return nil, meta, err
	}
	return out, meta, nil
}
