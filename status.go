package socialkit

import (
	"context"
)

// StatusService provides access to the /status endpoint of the SocialKit
// API, which reports the operational status of all tools and channels.
type StatusService struct {
	client *Client
}

// Get retrieves the current operational status of all SocialKit tools.
// This endpoint does not require an access key.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/status-api
func (s *StatusService) Get(ctx context.Context) (*Status, *ResponseMeta, error) {
	var out Status
	meta, err := s.client.doRequestWithMeta(ctx, "GET", "/status", nil, nil, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}
