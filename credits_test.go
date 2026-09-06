package socialkit

import (
	"context"
	"net/http"
	"testing"
)

func TestCreditsService_Get(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/credits" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.Header.Get("x-access-key") != "test-access-key" {
			t.Errorf("expected x-access-key header, got %q", r.Header.Get("x-access-key"))
		}
		w.Header().Set("X-Credits-Remaining", "500")
		writeResponse(t, w, `{"success":true,"data":{"totalRemaining":500,"monthly":{"remaining":300,"used":700,"total":1000},"purchased":{"remaining":200,"used":0,"total":200}}}`)
	})
	defer srv.Close()

	credits, meta, err := c.Credits.Get(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if credits.TotalRemaining != 500 {
		t.Errorf("expected totalRemaining 500, got %f", credits.TotalRemaining)
	}
	if credits.Monthly == nil {
		t.Fatal("expected monthly breakdown")
	}
	if *credits.Monthly.Remaining != 300 {
		t.Errorf("expected monthly remaining 300, got %v", *credits.Monthly.Remaining)
	}
	if *credits.Monthly.Used != 700 {
		t.Errorf("expected monthly used 700, got %v", *credits.Monthly.Used)
	}
	if credits.Purchased == nil {
		t.Fatal("expected purchased breakdown")
	}
	if *credits.Purchased.Remaining != 200 {
		t.Errorf("expected purchased remaining 200, got %v", *credits.Purchased.Remaining)
	}
	if meta == nil || meta.CreditsRemaining == nil || *meta.CreditsRemaining != 500 {
		t.Errorf("expected CreditsRemaining 500 in meta, got %v", meta)
	}
}

func TestCreditsService_Get_Unauthorized(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		writeResponse(t, w, `{"message":"invalid access key"}`)
	})
	defer srv.Close()

	_, _, err := c.Credits.Get(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
