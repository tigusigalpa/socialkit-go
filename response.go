package socialkit

import (
	"encoding/json"
)

// apiResponse is the internal response envelope used by all SocialKit API
// endpoints. The Data field is kept as json.RawMessage so that each
// service method can decode it into its own typed struct.
type apiResponse struct {
	// Success indicates whether the API request succeeded.
	Success bool `json:"success"`
	// Data holds the response payload as raw JSON, decoded by the caller.
	Data json.RawMessage `json:"data"`
	// Message is an optional human-readable message from the API.
	Message string `json:"message"`
}

// Response is a generic response wrapper that pairs decoded data with
// response metadata. It is returned by service methods alongside the
// decoded data type.
type Response[T any] struct {
	// Data is the decoded response payload.
	Data T
	// Meta holds response metadata extracted from headers.
	Meta *ResponseMeta
}
