package socialkit

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// ResponseMeta holds metadata extracted from HTTP response headers and
// the raw response body. It is returned alongside the decoded data by
// every service method. No authentication information is ever included
// in the headers map or raw body.
type ResponseMeta struct {
	// HTTPStatus is the HTTP status code of the response.
	HTTPStatus int
	// Headers contains selected response headers (excluding any
	// authentication-related headers).
	Headers map[string][]string
	// CreditsUsed is the number of credits consumed by the request,
	// parsed from the X-Credits-Used header.
	CreditsUsed *float64
	// CreditsRemaining is the number of credits remaining on the
	// account, parsed from the X-Credits-Remaining header.
	CreditsRemaining *float64
	// RetryAfter is the number of seconds the client should wait before
	// retrying, parsed from the Retry-After header.
	RetryAfter *int
	// RateLimitLimit is the rate limit ceiling, parsed from the
	// X-RateLimit-Limit header.
	RateLimitLimit *int
	// RateLimitRemaining is the number of requests remaining in the
	// current rate limit window, parsed from the
	// X-RateLimit-Remaining header.
	RateLimitRemaining *int
	// RateLimitReset is the Unix timestamp at which the rate limit
	// window resets, parsed from the X-RateLimit-Reset header.
	RateLimitReset *int64
	// RawBody contains the raw, redacted response body for debugging.
	RawBody json.RawMessage
}

// buildMeta constructs a ResponseMeta from an http.Response's headers and
// status code, plus the raw response body. Authentication headers are
// excluded from the Headers map.
func buildMeta(statusCode int, header http.Header, rawBody []byte) *ResponseMeta {
	meta := &ResponseMeta{
		HTTPStatus: statusCode,
		Headers:    make(map[string][]string),
	}

	// Copy headers, excluding any auth-related ones.
	for k, vs := range header {
		if isAuthHeader(k) {
			continue
		}
		meta.Headers[k] = vs
	}

	// Parse credit headers.
	if v := header.Get("X-Credits-Used"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			meta.CreditsUsed = &f
		}
	}
	if v := header.Get("X-Credits-Remaining"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			meta.CreditsRemaining = &f
		}
	}

	// Parse rate limit headers.
	if v := header.Get("Retry-After"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			meta.RetryAfter = &n
		}
	}
	if v := header.Get("X-RateLimit-Limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			meta.RateLimitLimit = &n
		}
	}
	if v := header.Get("X-RateLimit-Remaining"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			meta.RateLimitRemaining = &n
		}
	}
	if v := header.Get("X-RateLimit-Reset"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			meta.RateLimitReset = &n
		}
	}

	// Store a redacted copy of the raw body.
	if len(rawBody) > 0 {
		meta.RawBody = redactJSON(rawBody)
	}

	return meta
}

// isAuthHeader reports whether the given canonical header name is
// authentication-related and should be excluded from ResponseMeta.
func isAuthHeader(canonical string) bool {
	switch canonical {
	case "X-Access-Key", "Authorization", "Access-Key":
		return true
	}
	return false
}
