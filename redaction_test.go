package socialkit

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestRedaction_KeyNotInErrorOutput verifies that the access key never
// appears in any error output, including APIError messages, RawBody, or
// ResponseMeta.
func TestRedaction_KeyNotInErrorOutput(t *testing.T) {
	const secretKey = "super-secret-key-12345"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		// The response body includes the access_key field — it must be
		// redacted in the error output.
		writeResponse(t, w, `{"message":"forbidden","code":"access_denied","access_key":"super-secret-key-12345"}`)
	}))
	defer srv.Close()

	c := NewClient(secretKey, WithBaseURL(srv.URL))
	_, _, err := c.Credits.Get(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Check that the key does not appear in the error string.
	if strings.Contains(err.Error(), secretKey) {
		t.Errorf("access key found in error string: %s", err.Error())
	}

	// Check that the key does not appear in APIError.RawBody.
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		if strings.Contains(string(apiErr.RawBody), secretKey) {
			t.Errorf("access key found in RawBody: %s", string(apiErr.RawBody))
		}
		if apiErr.Meta != nil {
			if strings.Contains(string(apiErr.Meta.RawBody), secretKey) {
				t.Errorf("access key found in Meta.RawBody: %s", string(apiErr.Meta.RawBody))
			}
		}
	}
}

// TestRedaction_KeyNotInMetaHeaders verifies that auth headers are not
// included in ResponseMeta.Headers.
func TestRedaction_KeyNotInMetaHeaders(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set an auth-like header in the response — it should be excluded.
		w.Header().Set("X-Access-Key", "should-not-appear")
		w.Header().Set("Authorization", "should-not-appear")
		w.Header().Set("X-Custom", "should-appear")
		writeResponse(t, w, `{"success":true,"data":{"totalRemaining":100}}`)
	}))
	defer srv.Close()

	c := NewClient("test-key", WithBaseURL(srv.URL))
	_, meta, err := c.Credits.Get(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if meta == nil {
		t.Fatal("expected non-nil meta")
	}
	for k := range meta.Headers {
		if isAuthHeader(k) {
			t.Errorf("auth header %q should not be in meta headers", k)
		}
	}
	if _, ok := meta.Headers["X-Custom"]; !ok {
		t.Error("expected X-Custom header in meta")
	}
}

// TestRedaction_KeyInQueryMode verifies that when keyInQuery is enabled,
// the key is sent in the query/body but still does not appear in error
// output.
func TestRedaction_KeyInQueryMode(t *testing.T) {
	const secretKey = "my-secret-key"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the key is in the query string.
		if r.URL.Query().Get("access_key") != secretKey {
			t.Errorf("expected access_key in query, got %q", r.URL.Query().Get("access_key"))
		}
		w.WriteHeader(http.StatusUnauthorized)
		writeResponse(t, w, `{"message":"unauthorized","access_key":"my-secret-key"}`)
	}))
	defer srv.Close()

	c := NewClient(secretKey, WithBaseURL(srv.URL), WithKeyInQuery(true))
	_, _, err := c.Credits.Get(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if strings.Contains(err.Error(), secretKey) {
		t.Errorf("access key found in error string: %s", err.Error())
	}
}

// TestRedaction_KeyInQueryMode_POST verifies that when keyInQuery is
// enabled for POST requests, the key is injected into the JSON body.
func TestRedaction_KeyInQueryMode_POST(t *testing.T) {
	const secretKey = "post-secret-key"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body["access_key"] != secretKey {
			t.Errorf("expected access_key in body, got %v", body["access_key"])
		}
		if body["url"] != "https://youtube.com/watch?v=test" {
			t.Errorf("unexpected url: %v", body["url"])
		}
		writeResponse(t, w, `{"success":true,"data":{"url":"https://youtube.com/watch?v=test","transcript":"hello"}}`)
	}))
	defer srv.Close()

	c := NewClient(secretKey, WithBaseURL(srv.URL), WithKeyInQuery(true))
	_, _, err := c.YouTube.Transcript(context.Background(), &FetchRequest{
		URL: "https://youtube.com/watch?v=test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestRedactJSON verifies that the redactJSON function replaces known
// credential keys with [REDACTED].
func TestRedactJSON(t *testing.T) {
	input := []byte(`{"access_key":"secret123","data":{"api_key":"nested-secret","other":"value"}}`)
	output := redactJSON(input)
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if result["access_key"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %v", result["access_key"])
	}
	data := result["data"].(map[string]interface{})
	if data["api_key"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %v", data["api_key"])
	}
	if data["other"] != "value" {
		t.Errorf("expected value, got %v", data["other"])
	}
}

func TestRedactJSON_NilBody(t *testing.T) {
	output := redactJSON(nil)
	if output != nil {
		t.Errorf("expected nil, got %v", output)
	}
}

func TestRedactJSON_NonObjectJSON(t *testing.T) {
	input := []byte(`[1,2,3]`)
	output := redactJSON(input)
	if string(output) != `[1,2,3]` {
		t.Errorf("expected [1,2,3], got %s", string(output))
	}
}

func TestRedactJSON_RedactsNestedArraysAndCredentialSpelling(t *testing.T) {
	input := []byte(`{"items":[{"X-API-Key":"secret"},{"nested":{"Authorization":"Bearer secret"}}]}`)
	output := redactJSON(input)
	if strings.Contains(string(output), "secret") {
		t.Errorf("credential was not redacted: %s", output)
	}
}

func TestRedactResponseBody_RedactsPlainTextEcho(t *testing.T) {
	const accessKey = "secret-key-in-url"
	body := redactResponseBody([]byte("upstream request failed: ?access_key="+accessKey), accessKey)
	if strings.Contains(string(body), accessKey) {
		t.Errorf("configured key was not redacted: %s", body)
	}
}
