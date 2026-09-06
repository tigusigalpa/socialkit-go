package socialkit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusService_Get(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/status" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		writeResponse(t, w, `{"success":true,"data":{"overall":"operational","generatedAt":"2026-01-01T00:00:00Z","tools":[{"id":"youtube","channel":"youtube","name":"YouTube","status":"operational","message":"","updatedAt":"2026-01-01T00:00:00Z"}]}}`)
	})
	defer srv.Close()

	status, meta, err := c.Status.Get(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.Overall != "operational" {
		t.Errorf("expected overall operational, got %s", status.Overall)
	}
	if status.GeneratedAt != "2026-01-01T00:00:00Z" {
		t.Errorf("unexpected generatedAt: %s", status.GeneratedAt)
	}
	if len(status.Tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(status.Tools))
	}
	if status.Tools[0].ID != "youtube" {
		t.Errorf("expected tool id youtube, got %s", status.Tools[0].ID)
	}
	if status.Tools[0].Name != "YouTube" {
		t.Errorf("expected tool name YouTube, got %s", status.Tools[0].Name)
	}
	if meta == nil {
		t.Fatal("expected non-nil meta")
	}
	if meta.HTTPStatus != http.StatusOK {
		t.Errorf("expected HTTP status 200, got %d", meta.HTTPStatus)
	}
}

func TestStatusService_Get_NoKeyRequired(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Status endpoint should still send the header, but it works
		// even without a valid key.
		writeResponse(t, w, `{"success":true,"data":{"overall":"operational","generatedAt":"","tools":[]}}`)
	}))
	defer srv.Close()

	// Create a client with an empty key.
	c2 := NewClient("", WithBaseURL(srv.URL))
	status, _, err := c2.Status.Get(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.Overall != "operational" {
		t.Errorf("expected operational, got %s", status.Overall)
	}
}
