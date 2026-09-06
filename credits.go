package socialkit

import (
	"context"
)

// CreditsService provides access to the /credits endpoint of the
// SocialKit API, which reports the credit balance for the authenticated
// account.
type CreditsService struct {
	client *Client
}

// Get retrieves the current credit balance for the authenticated
// account, including monthly and purchased credit breakdowns.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/credits-api
func (s *CreditsService) Get(ctx context.Context) (*Credits, *ResponseMeta, error) {
	var out Credits
	meta, err := s.client.doRequestWithMeta(ctx, "GET", "/credits", nil, nil, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}
