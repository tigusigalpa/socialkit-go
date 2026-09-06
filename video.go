package socialkit

import (
	"context"
)

// VideoService provides access to the /video/* endpoints of the
// SocialKit API, covering transcript and summarize operations for direct
// video file URLs. Both operations cost 5 credits per started minute.
type VideoService struct {
	client *Client
}

// Transcript retrieves the transcript for a direct video file URL.
// Costs 5 credits per started minute.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/video-transcript-api
func (s *VideoService) Transcript(ctx context.Context, req *FetchRequest) (*Transcript, *ResponseMeta, error) {
	var out Transcript
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/video/transcript", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Summarize generates an AI summary for a direct video file URL.
// Costs 5 credits per started minute.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/video-summarize-api
func (s *VideoService) Summarize(ctx context.Context, req *SummaryRequest) (*Summary, *ResponseMeta, error) {
	var out Summary
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/video/summarize", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}
