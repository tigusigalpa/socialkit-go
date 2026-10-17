package socialkit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// maxErrorBodySize limits how much of a non-2xx response body is read
// into memory before decoding, to prevent unbounded memory usage from
// misbehaving servers.
const maxErrorBodySize = 1 << 20 // 1 MB

// escapePathSegment escapes a caller-provided value that occupies one URL
// path segment. This prevents characters such as '/', '?', and '#' from
// changing the endpoint being requested.
func escapePathSegment(value string) string {
	return url.PathEscape(value)
}

// doRequestWithMeta performs an HTTP request against the SocialKit API,
// decodes the JSON response body into out (if non-nil), and returns the
// response metadata alongside any error. It handles authentication,
// query string construction, JSON error mapping, and automatic retry
// with exponential backoff for HTTP 429 and 5xx responses.
//
// For POST requests, the body is marshalled to JSON and sent in the
// request body. For GET requests, the body struct is ignored (query
// params are used instead).
//
// Authentication is sent via the x-access-key header by default. When
// keyInQuery is enabled, the key is also added to the query string (GET)
// or JSON body (POST).
func (c *Client) doRequestWithMeta(ctx context.Context, method, path string, query url.Values, body interface{}, out interface{}) (*ResponseMeta, error) {
	method = strings.ToUpper(method)

	var bodyBytes []byte
	if body != nil && method != http.MethodGet {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("socialkit: failed to marshal request body: %w", err)
		}
	}

	// When keyInQuery is enabled, add the access key to the query or body.
	if c.keyInQuery {
		if method == http.MethodGet {
			if query == nil {
				query = url.Values{}
			}
			query.Set("access_key", c.accessKey)
		} else if body != nil {
			// Inject access_key into the JSON body.
			var raw map[string]json.RawMessage
			if err := json.Unmarshal(bodyBytes, &raw); err == nil {
				raw["access_key"], _ = json.Marshal(c.accessKey)
				bodyBytes, _ = json.Marshal(raw)
			}
		}
	}

	fullURL := c.baseURL + path
	if method == http.MethodGet && len(query) > 0 {
		fullURL += "?" + query.Encode()
	}

	maxAttempts := c.maxAttempts
	if maxAttempts < 1 {
		maxAttempts = 1
	}

	// POST requests are only retried if the caller explicitly opts in
	// by setting maxAttempts > 1. This is the default behaviour since
	// WithRetry is opt-in.
	canRetryPOST := maxAttempts > 1

	var lastErr error
	var lastMeta *ResponseMeta
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return lastMeta, err
		}

		var bodyReader io.Reader
		if body != nil && method != http.MethodGet {
			// Recreate the reader for every retry attempt.
			bodyReader = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
		if err != nil {
			return nil, fmt.Errorf("socialkit: failed to build request: %w", err)
		}

		// Authentication via header (default, always set unless keyInQuery
		// is the only auth mode — but we set it regardless for safety).
		req.Header.Set("x-access-key", c.accessKey)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", c.userAgent)
		if body != nil && method != http.MethodGet {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			if ctx.Err() != nil {
				return lastMeta, ctx.Err()
			}
			lastErr = fmt.Errorf("%w: %s", ErrTransport, redactString(err.Error(), c.accessKey))
			if (method != http.MethodPost || canRetryPOST) && attempt < maxAttempts {
				if err := waitForRetry(ctx, retryDelay("", c.baseDelay, attempt)); err != nil {
					return lastMeta, err
				}
				continue
			}
			return lastMeta, lastErr
		}

		// Read the full response body, bounding non-2xx bodies.
		var respBody []byte
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			respBody, err = io.ReadAll(resp.Body)
		} else {
			respBody, err = io.ReadAll(io.LimitReader(resp.Body, maxErrorBodySize))
		}
		// The read result determines success; a cleanup error must not
		// cause an already completed, potentially billable call to repeat.
		_ = resp.Body.Close()
		if err != nil {
			meta := buildMeta(resp.StatusCode, resp.Header, nil)
			return meta, fmt.Errorf("%w: %v", ErrDecode, err)
		}

		// Build metadata from headers.
		// Redact before the body is used anywhere caller-visible. Besides
		// known JSON credential fields, replace the configured key itself:
		// some proxies return plain-text/HTML errors that echo request URLs.
		safeBody := redactResponseBody(respBody, c.accessKey)
		meta := buildMeta(resp.StatusCode, resp.Header, safeBody)

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if out != nil && len(respBody) > 0 {
				if err := decodeAPIResponse(respBody, out); err != nil {
					return meta, fmt.Errorf("%w: %v", ErrDecode, err)
				}
			}
			return meta, nil
		}

		// Non-2xx: build typed error.
		msg, code := extractMessageAndCode(safeBody)
		lastErr = newAPIError(resp.StatusCode, msg, code, safeBody, meta)
		lastMeta = meta

		// Determine if this status is retryable.
		retryable := isRetryableStatus(resp.StatusCode)
		// POST requests are only retried when explicitly opted in.
		if method == http.MethodPost && !canRetryPOST {
			retryable = false
		}

		if retryable && attempt < maxAttempts {
			if err := waitForRetry(ctx, retryDelay(resp.Header.Get("Retry-After"), c.baseDelay, attempt)); err != nil {
				return lastMeta, err
			}
			continue
		}

		return lastMeta, lastErr
	}

	if lastErr != nil {
		return lastMeta, lastErr
	}
	return lastMeta, newAPIError(http.StatusInternalServerError, "unknown error", "", nil, nil)
}

// decodeAPIResponse decodes the SocialKit response envelope. If the body
// contains the standard {success, data, message} envelope, the Data
// field is extracted and decoded into out. Otherwise, the body is
// decoded directly into out (for endpoints that return raw payloads).
func decodeAPIResponse(body []byte, out interface{}) error {
	var envelope apiResponse
	if err := json.Unmarshal(body, &envelope); err != nil {
		// Not valid JSON envelope — try decoding directly.
		if err2 := json.Unmarshal(body, out); err2 != nil {
			return err
		}
		return nil
	}

	// If the envelope has a Data field, decode it into out.
	if len(envelope.Data) > 0 {
		if err := json.Unmarshal(envelope.Data, out); err != nil {
			return err
		}
		return nil
	}

	// If there's no Data field but the body is valid JSON, try decoding
	// the whole body directly (some endpoints return data without an
	// envelope).
	if err := json.Unmarshal(body, out); err != nil {
		// If out is a *json.RawMessage, assign the body directly.
		if raw, ok := out.(*json.RawMessage); ok {
			*raw = body
			return nil
		}
		return err
	}
	return nil
}

// extractMessageAndCode attempts to pull a human-readable message and
// machine-readable code from a JSON error response body.
func extractMessageAndCode(body []byte) (string, string) {
	var parsed struct {
		Message string `json:"message"`
		Error   string `json:"error"`
		Code    string `json:"code"`
		Data    struct {
			Message string `json:"message"`
			Error   string `json:"error"`
			Code    string `json:"code"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &parsed); err == nil {
		msg := parsed.Message
		if msg == "" {
			msg = parsed.Error
		}
		if msg == "" {
			msg = parsed.Data.Message
		}
		if msg == "" {
			msg = parsed.Data.Error
		}
		code := parsed.Code
		if code == "" {
			code = parsed.Data.Code
		}
		if msg != "" || code != "" {
			return msg, code
		}
	}
	if len(body) == 0 {
		return "unknown error", ""
	}
	snippet := string(body)
	if len(snippet) > 200 {
		snippet = snippet[:200]
	}
	return snippet, ""
}

// isRetryableStatus reports whether the given HTTP status code should
// trigger a retry attempt.
func isRetryableStatus(statusCode int) bool {
	switch statusCode {
	case http.StatusTooManyRequests, http.StatusInternalServerError, http.StatusServiceUnavailable, http.StatusBadGateway, http.StatusGatewayTimeout:
		return true
	}
	return false
}

// retryDelay computes the delay to wait before the next retry attempt.
// It honors a Retry-After header (in seconds) if present and valid,
// otherwise it falls back to exponential backoff based on baseDelay and
// the current attempt number, with up to 25% jitter.
func retryDelay(retryAfterHeader string, baseDelay time.Duration, attempt int) time.Duration {
	if retryAfterHeader != "" {
		if secs, err := strconv.Atoi(strings.TrimSpace(retryAfterHeader)); err == nil && secs >= 0 {
			return time.Duration(secs) * time.Second
		}
		if retryAt, err := http.ParseTime(strings.TrimSpace(retryAfterHeader)); err == nil {
			return max(time.Until(retryAt), 0)
		}
	}
	if baseDelay <= 0 {
		baseDelay = DefaultBaseDelay
	}
	multiplier := 1 << (attempt - 1)
	delay := baseDelay * time.Duration(multiplier)
	// Add up to 25% jitter to avoid thundering herd.
	// rand.Int63n(0) panics, which is possible with deliberately tiny
	// delays in tests or latency-sensitive callers.
	if delay < 4 {
		return delay
	}
	jitter := time.Duration(rand.Int63n(int64(delay) / 4))
	return delay + jitter
}

// waitForRetry waits without masking cancellation or leaking timers.
func waitForRetry(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// redactJSON returns a copy of the given JSON body with any access key
// values replaced by [REDACTED]. This ensures no credential ever appears
// in error output, metadata, or logs.
func redactJSON(body []byte) json.RawMessage {
	if len(body) == 0 {
		return nil
	}
	var raw interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		// Not a JSON object — return as-is (no string scanning to avoid
		// false positives).
		return json.RawMessage(body)
	}
	redactValue(raw)
	out, err := json.Marshal(raw)
	if err != nil {
		return json.RawMessage(body)
	}
	return out
}

// redactResponseBody redacts structured credential fields and the
// configured credential itself. The latter covers text or HTML responses
// that contain an echoed request URL.
func redactResponseBody(body []byte, accessKey string) json.RawMessage {
	redacted := redactJSON(body)
	if accessKey == "" {
		return redacted
	}
	return json.RawMessage(bytes.ReplaceAll(redacted, []byte(accessKey), []byte("[REDACTED]")))
}

func redactString(value, accessKey string) string {
	if accessKey == "" {
		return value
	}
	return strings.ReplaceAll(value, accessKey, "[REDACTED]")
}

// redactValue recursively replaces values for credential-like keys. It
// handles nested objects and arrays and accepts common spelling variants.
func redactValue(value interface{}) {
	switch v := value.(type) {
	case map[string]interface{}:
		for key, child := range v {
			if isCredentialKey(key) {
				v[key] = "[REDACTED]"
				continue
			}
			redactValue(child)
		}
	case []interface{}:
		for _, child := range v {
			redactValue(child)
		}
	}
}

func isCredentialKey(key string) bool {
	normalized := strings.NewReplacer("_", "", "-", "", ".", "").Replace(strings.ToLower(key))
	switch normalized {
	case "key", "accesskey", "apikey", "xaccesskey", "xapikey", "authorization", "token", "accesstoken", "bearertoken":
		return true
	default:
		return false
	}
}

// jsonUnmarshal is a thin wrapper around json.Unmarshal used by other
// files to avoid duplicate imports.
func jsonUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
