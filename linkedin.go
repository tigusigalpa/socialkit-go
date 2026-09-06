package socialkit

import (
	"context"
	"fmt"
)

// LinkedInService provides access to the /linkedin/* endpoints of the
// SocialKit API, covering profile, company, company posts, individual
// post, and transcript operations.
type LinkedInService struct {
	client *Client
}

// Profile retrieves a LinkedIn user profile, including recent articles.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/linkedin-profile-api
func (s *LinkedInService) Profile(ctx context.Context, req *FetchRequest) (*LinkedInProfile, *ResponseMeta, error) {
	var out LinkedInProfile
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/linkedin/profile", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Company retrieves a LinkedIn company page.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/linkedin-company-api
func (s *LinkedInService) Company(ctx context.Context, req *FetchRequest) (*LinkedInCompany, *ResponseMeta, error) {
	var out LinkedInCompany
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/linkedin/company", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// CompanyPosts retrieves posts from a LinkedIn company page. Limit is
// capped at 50.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/linkedin-company-api-posts
func (s *LinkedInService) CompanyPosts(ctx context.Context, req *CompanyPostsRequest) (*LinkedInCompanyPosts, *ResponseMeta, error) {
	if req != nil && req.Limit != nil && *req.Limit > 50 {
		return nil, nil, fmt.Errorf("socialkit: linkedin company posts limit must be at most 50")
	}
	var out LinkedInCompanyPosts
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/linkedin/company-posts", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Post retrieves a single LinkedIn post. The post content is returned as
// raw JSON because the upstream API does not guarantee a fixed schema.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/linkedin-post-api
func (s *LinkedInService) Post(ctx context.Context, req *FetchRequest) (*LinkedInPost, *ResponseMeta, error) {
	var out LinkedInPost
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/linkedin/post", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Transcript retrieves the transcript for a LinkedIn video. Segments do
// not include a timestamp field.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/linkedin-transcript-api
func (s *LinkedInService) Transcript(ctx context.Context, req *FetchRequest) (*LinkedInTranscript, *ResponseMeta, error) {
	var out LinkedInTranscript
	meta, err := s.client.doRequestWithMeta(ctx, "POST", "/linkedin/transcript", nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}
