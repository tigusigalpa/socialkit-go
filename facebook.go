package socialkit

import (
	"context"
)

// FacebookService provides access to the /facebook/* endpoints of the
// SocialKit API, covering transcripts, summaries, stats, and comments.
// Facebook does not support sync download or bulk operations.
type FacebookService struct {
	client *Client
}

// Transcript retrieves the transcript for a Facebook video.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/facebook-transcript-api
func (s *FacebookService) Transcript(ctx context.Context, req *FetchRequest) (*Transcript, *ResponseMeta, error) {
	var out Transcript
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/facebook/transcript", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Summarize generates an AI summary of a Facebook video.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/facebook-summarize-api
func (s *FacebookService) Summarize(ctx context.Context, req *SummaryRequest) (*Summary, *ResponseMeta, error) {
	var out Summary
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/facebook/summarize", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Stats retrieves statistics for a Facebook post, including reactions.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/facebook-stats-api
func (s *FacebookService) Stats(ctx context.Context, req *FetchRequest) (*FacebookStats, *ResponseMeta, error) {
	var out FacebookStats
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/facebook/stats", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Comments retrieves comments for a Facebook post with cursor-based
// pagination. The commentCount field may be null.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/facebook-comments-api
func (s *FacebookService) Comments(ctx context.Context, req *CommentsRequest) (*FacebookComments, *ResponseMeta, error) {
	var out FacebookComments
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/facebook/comments", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// ChannelStats retrieves statistics for a Facebook page/profile.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/facebook-channel-stats-api
func (s *FacebookService) ChannelStats(ctx context.Context, req *FetchRequest) (*FacebookChannelStats, *ResponseMeta, error) {
	var out FacebookChannelStats
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/facebook/channel-stats", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}
