package socialkit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// newTestClient spins up an httptest.Server using handler and returns a
// *Client configured to talk to it, along with the server for cleanup.
func newTestClient(t *testing.T, handler http.HandlerFunc, opts ...Option) (*Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	allOpts := append([]Option{WithBaseURL(srv.URL)}, opts...)
	c := NewClient("test-access-key", allOpts...)
	return c, srv
}

// writeResponse writes an HTTP response body and reports an unexpected
// write failure to the active test.
func writeResponse(t *testing.T, w http.ResponseWriter, body string) {
	t.Helper()
	if _, err := w.Write([]byte(body)); err != nil {
		t.Errorf("failed to write test response: %v", err)
	}
}

// envelope wraps data in the standard API response envelope. The data
// argument must be a valid JSON string (e.g. `{"key":"value"}` or `[1,2,3]`).
func envelope(data string) string {
	return `{"success":true,"data":` + data + `}`
}

func TestNewClient_Defaults(t *testing.T) {
	c := NewClient("key")
	if c.baseURL != DefaultBaseURL {
		t.Errorf("expected baseURL %s, got %s", DefaultBaseURL, c.baseURL)
	}
	if c.httpClient.Timeout != DefaultTimeout {
		t.Errorf("expected timeout %s, got %s", DefaultTimeout, c.httpClient.Timeout)
	}
	if c.maxAttempts != DefaultMaxAttempts {
		t.Errorf("expected maxAttempts %d, got %d", DefaultMaxAttempts, c.maxAttempts)
	}
	if c.keyInQuery != false {
		t.Error("expected keyInQuery to be false by default")
	}
	if c.userAgent != defaultUserAgent {
		t.Errorf("expected userAgent %s, got %s", defaultUserAgent, c.userAgent)
	}
	if c.Status == nil || c.Credits == nil || c.YouTube == nil || c.TikTok == nil ||
		c.Instagram == nil || c.Facebook == nil || c.Twitter == nil ||
		c.LinkedIn == nil || c.Video == nil || c.Downloads == nil {
		t.Error("expected all services to be initialized")
	}
}

func TestNewClient_Options(t *testing.T) {
	customClient := &http.Client{}
	c := NewClient("key",
		WithBaseURL("https://example.com"),
		WithHTTPClient(customClient),
		WithTimeout(5*time.Second),
		WithRetry(3, 100*time.Millisecond),
		WithKeyInQuery(true),
		WithUserAgent("MyApp/1.0"),
	)
	if c.baseURL != "https://example.com" {
		t.Errorf("expected custom baseURL, got %s", c.baseURL)
	}
	if c.httpClient != customClient {
		t.Error("expected custom http client to be used")
	}
	if c.httpClient.Timeout != 5*time.Second {
		t.Errorf("expected timeout 5s, got %s", c.httpClient.Timeout)
	}
	if c.maxAttempts != 3 {
		t.Errorf("expected maxAttempts 3, got %d", c.maxAttempts)
	}
	if c.baseDelay != 100*time.Millisecond {
		t.Errorf("expected baseDelay 100ms, got %s", c.baseDelay)
	}
	if !c.keyInQuery {
		t.Error("expected keyInQuery to be true")
	}
	if c.userAgent != "MyApp/1.0" {
		t.Errorf("expected userAgent MyApp/1.0, got %s", c.userAgent)
	}
}

func TestWithBaseURL_TrailingSlash(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/status"; got != want {
			t.Errorf("expected path %q, got %q", want, got)
		}
		writeResponse(t, w, envelope(`{"overall":"operational"}`))
	}))
	defer srv.Close()

	c := NewClient("test-access-key", WithBaseURL(srv.URL+"/"))
	if _, _, err := c.Status.Get(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPathParametersAreEscaped(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.EscapedPath(), "/v2/downloads/a%2Fb%3Fc%23d%20e"; got != want {
			t.Errorf("expected escaped path %q, got %q", want, got)
		}
		writeResponse(t, w, envelope(`{"jobId":"test","status":"ready"}`))
	})
	defer srv.Close()

	if _, _, err := c.Downloads.Get(context.Background(), "a/b?c#d e"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDoRequest_Retry429(t *testing.T) {
	var attempts int
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			writeResponse(t, w, `{"message":"rate limited"}`)
			return
		}
		writeResponse(t, w, envelope(`{"overall":"operational"}`))
	}, WithRetry(5, 10*time.Millisecond))
	defer srv.Close()

	resp, _, err := c.Status.Get(context.Background())
	if err != nil {
		t.Fatalf("expected no error after retries, got %v", err)
	}
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
	if resp.Overall != "operational" {
		t.Errorf("expected overall operational, got %s", resp.Overall)
	}
}

func TestRetryDelay_TinyDelayDoesNotPanic(t *testing.T) {
	if got := retryDelay("", time.Nanosecond, 1); got != time.Nanosecond {
		t.Errorf("expected a 1ns delay, got %s", got)
	}
}

func TestDoRequest_RetriesAndRedactsTransportErrors(t *testing.T) {
	const secretKey = "retry-secret-key"
	attempts := 0
	client := NewClient(secretKey,
		WithHTTPClient(&http.Client{Transport: roundTripperFunc(func(*http.Request) (*http.Response, error) {
			attempts++
			return nil, fmt.Errorf("dial failed for access_key=%s", secretKey)
		})}),
		WithRetry(3, time.Nanosecond),
	)

	_, _, err := client.Credits.Get(context.Background())
	if !errors.Is(err, ErrTransport) {
		t.Fatalf("expected transport error, got %v", err)
	}
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
	if strings.Contains(err.Error(), secretKey) {
		t.Errorf("transport error leaked access key: %v", err)
	}
}

func TestDoRequest_RetryPreservesRequestBody(t *testing.T) {
	attempts := 0
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		attempts++
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body on attempt %d: %v", attempts, err)
		}
		if got, want := body["url"], "https://youtube.com/watch?v=test"; got != want {
			t.Errorf("expected body url %q, got %q", want, got)
		}
		if attempts == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		writeResponse(t, w, envelope(`{"url":"https://youtube.com/watch?v=test","transcript":"hello"}`))
	}, WithRetry(2, time.Millisecond))
	defer srv.Close()

	resp, _, err := c.YouTube.Transcript(context.Background(), &FetchRequest{
		URL: "https://youtube.com/watch?v=test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attempts != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts)
	}
	if resp.Transcript != "hello" {
		t.Errorf("expected transcript hello, got %s", resp.Transcript)
	}
}

func TestDoRequest_RateLimitedExhausted(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		writeResponse(t, w, `{"message":"rate limited"}`)
	}, WithRetry(2, 5*time.Millisecond))
	defer srv.Close()

	_, _, err := c.Status.Get(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrRateLimited) {
		t.Errorf("expected ErrRateLimited, got %v", err)
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Errorf("expected *APIError, got %T", err)
	} else if apiErr.StatusCode != http.StatusTooManyRequests {
		t.Errorf("expected status 429, got %d", apiErr.StatusCode)
	}
}

func TestDoRequest_Unauthorized(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		writeResponse(t, w, `{"message":"invalid access key"}`)
	})
	defer srv.Close()

	_, _, err := c.Credits.Get(context.Background())
	if !errors.Is(err, ErrUnauthorized) {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestDoRequest_NotFound(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		writeResponse(t, w, `{"message":"not found"}`)
	})
	defer srv.Close()

	_, _, err := c.Credits.Get(context.Background())
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestDoRequest_BadRequest(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		writeResponse(t, w, `{"message":"bad request"}`)
	})
	defer srv.Close()

	_, _, err := c.Credits.Get(context.Background())
	if !errors.Is(err, ErrBadRequest) {
		t.Errorf("expected ErrBadRequest, got %v", err)
	}
}

func TestDoRequest_ServerError(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		writeResponse(t, w, `{"message":"server error"}`)
	})
	defer srv.Close()

	_, _, err := c.Credits.Get(context.Background())
	if !errors.Is(err, ErrServer) {
		t.Errorf("expected ErrServer, got %v", err)
	}
}

func TestDoRequest_ContextCancellation(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		writeResponse(t, w, envelope(`{}`))
	})
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	_, _, err := c.Status.Get(ctx)
	if err == nil {
		t.Fatal("expected context deadline error, got nil")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected context.DeadlineExceeded, got %v", err)
	}
}

func TestDoRequest_AccessKeyHeader(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("x-access-key"); got != "test-access-key" {
			t.Errorf("expected x-access-key header 'test-access-key', got %q", got)
		}
		writeResponse(t, w, envelope(`{}`))
	})
	defer srv.Close()

	if _, _, err := c.Credits.Get(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDoRequest_UserAgentHeader(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("User-Agent"); got != defaultUserAgent {
			t.Errorf("expected User-Agent %q, got %q", defaultUserAgent, got)
		}
		writeResponse(t, w, envelope(`{}`))
	})
	defer srv.Close()

	if _, _, err := c.Status.Get(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAPIError_Error(t *testing.T) {
	err := &APIError{StatusCode: 500, Message: "server error"}
	want := "socialkit: api error: status=500 message=server error"
	if err.Error() != want {
		t.Errorf("expected %q, got %q", want, err.Error())
	}
}

func TestAPIError_ErrorWithCode(t *testing.T) {
	err := &APIError{StatusCode: 403, Code: "insufficient_credits", Message: "not enough credits"}
	want := "socialkit: api error: status=403 code=insufficient_credits message=not enough credits"
	if err.Error() != want {
		t.Errorf("expected %q, got %q", want, err.Error())
	}
}

func TestResponseMeta_CreditsHeaders(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Credits-Used", "2")
		w.Header().Set("X-Credits-Remaining", "998")
		w.Header().Set("X-RateLimit-Limit", "100")
		w.Header().Set("X-RateLimit-Remaining", "95")
		w.Header().Set("X-RateLimit-Reset", "1700000000")
		writeResponse(t, w, envelope(`{"totalRemaining":998}`))
	})
	defer srv.Close()

	_, meta, err := c.Credits.Get(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if meta == nil {
		t.Fatal("expected non-nil meta")
	}
	if meta.CreditsUsed == nil || *meta.CreditsUsed != 2 {
		t.Errorf("expected CreditsUsed 2, got %v", meta.CreditsUsed)
	}
	if meta.CreditsRemaining == nil || *meta.CreditsRemaining != 998 {
		t.Errorf("expected CreditsRemaining 998, got %v", meta.CreditsRemaining)
	}
	if meta.RateLimitLimit == nil || *meta.RateLimitLimit != 100 {
		t.Errorf("expected RateLimitLimit 100, got %v", meta.RateLimitLimit)
	}
	if meta.RateLimitRemaining == nil || *meta.RateLimitRemaining != 95 {
		t.Errorf("expected RateLimitRemaining 95, got %v", meta.RateLimitRemaining)
	}
	if meta.RateLimitReset == nil || *meta.RateLimitReset != 1700000000 {
		t.Errorf("expected RateLimitReset 1700000000, got %v", meta.RateLimitReset)
	}
}
