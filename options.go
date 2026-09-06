package socialkit

import (
	"net/http"
	"strings"
	"time"
)

// Option configures a Client. Options are applied in the order they are
// passed to NewClient.
type Option func(*Client)

// WithBaseURL overrides the default SocialKit API base URL
// (https://api.socialkit.dev). This is primarily useful for testing
// against a mock server.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		if url != "" {
			c.baseURL = strings.TrimRight(url, "/")
		}
	}
}

// WithHTTPClient sets a custom *http.Client used to perform requests.
// This allows callers to configure transport-level behaviour such as
// proxies, TLS settings, or custom RoundTrippers.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		if client != nil {
			c.httpClient = client
		}
	}
}

// WithTimeout sets the timeout applied to the underlying HTTP client for
// every request made by the Client.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		if d > 0 {
			c.httpClient.Timeout = d
		}
	}
}

// WithRetry configures automatic retry behaviour for HTTP 429 (rate
// limited) and 5xx (server error) responses. maxAttempts is the total
// number of attempts (including the first one) and baseDelay is the
// initial backoff delay used for exponential backoff with jitter. The
// delay doubles on each subsequent retry unless a Retry-After header is
// present, in which case that value takes precedence.
//
// Calling WithRetry is explicit opt-in for retries of both GET and POST
// requests. A POST can be billable, so enable it only when your operation
// is safe to repeat or the upstream request is idempotent.
func WithRetry(maxAttempts int, baseDelay time.Duration) Option {
	return func(c *Client) {
		if maxAttempts > 0 {
			c.maxAttempts = maxAttempts
		}
		if baseDelay > 0 {
			c.baseDelay = baseDelay
		}
	}
}

// WithKeyInQuery enables the legacy/compatibility authentication mode
// where the access key is sent as an access_key query parameter (for GET
// requests) or in the JSON body (for POST requests) instead of (or in
// addition to) the x-access-key header. By default the key is only sent
// in the header and never appears in URLs or request bodies.
func WithKeyInQuery(enabled bool) Option {
	return func(c *Client) {
		c.keyInQuery = enabled
	}
}

// WithUserAgent overrides the default User-Agent header
// ("SocialKit-Go-SDK/1.0.0") sent with every request.
func WithUserAgent(ua string) Option {
	return func(c *Client) {
		if ua != "" {
			c.userAgent = ua
		}
	}
}
