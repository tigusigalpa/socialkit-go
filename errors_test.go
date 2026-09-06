package socialkit

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestErrBadRequest(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		writeResponse(t, w, `{"message":"bad request","code":"invalid_url"}`)
	})
	defer srv.Close()

	_, _, err := c.Credits.Get(context.Background())
	if !errors.Is(err, ErrBadRequest) {
		t.Errorf("expected ErrBadRequest, got %v", err)
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Errorf("expected *APIError, got %T", err)
	} else {
		if apiErr.Code != "invalid_url" {
			t.Errorf("expected code invalid_url, got %s", apiErr.Code)
		}
	}
}

func TestErrForbidden(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		writeResponse(t, w, `{"message":"forbidden","code":"access_denied"}`)
	})
	defer srv.Close()

	_, _, err := c.Credits.Get(context.Background())
	if !errors.Is(err, ErrForbidden) {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
}

func TestErrInsufficientCredits(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		writeResponse(t, w, `{"message":"insufficient credits","code":"insufficient_credits","data":{"requiredCredits":10,"remainingCredits":5,"shortfallCredits":5,"resetAt":"2026-02-01T00:00:00Z","upgradeUrl":"https://socialkit.dev/upgrade","topUpUrl":"https://socialkit.dev/topup","incidentId":"inc123"}}`)
	})
	defer srv.Close()

	_, _, err := c.Credits.Get(context.Background())
	if !errors.Is(err, ErrInsufficientCredits) {
		t.Errorf("expected ErrInsufficientCredits, got %v", err)
	}
	var icErr *InsufficientCreditsError
	if !errors.As(err, &icErr) {
		t.Errorf("expected *InsufficientCreditsError, got %T", err)
	} else {
		if icErr.RequiredCredits != 10 {
			t.Errorf("expected requiredCredits 10, got %f", icErr.RequiredCredits)
		}
		if icErr.RemainingCredits != 5 {
			t.Errorf("expected remainingCredits 5, got %f", icErr.RemainingCredits)
		}
		if icErr.ShortfallCredits != 5 {
			t.Errorf("expected shortfallCredits 5, got %f", icErr.ShortfallCredits)
		}
		if icErr.UpgradeURL != "https://socialkit.dev/upgrade" {
			t.Errorf("unexpected upgradeUrl: %s", icErr.UpgradeURL)
		}
		if icErr.TopUpURL != "https://socialkit.dev/topup" {
			t.Errorf("unexpected topUpUrl: %s", icErr.TopUpURL)
		}
		if icErr.IncidentID != "inc123" {
			t.Errorf("unexpected incidentId: %s", icErr.IncidentID)
		}
		if icErr.ResetAt == nil {
			t.Error("expected non-nil ResetAt")
		}
	}
	// Also verify *APIError is accessible.
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Errorf("expected *APIError, got %T", err)
	} else {
		if apiErr.StatusCode != http.StatusForbidden {
			t.Errorf("expected status 403, got %d", apiErr.StatusCode)
		}
		if apiErr.Code != "insufficient_credits" {
			t.Errorf("expected code insufficient_credits, got %s", apiErr.Code)
		}
	}
}

func TestErrInsufficientCredits_TopLevelFields(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		writeResponse(t, w, `{"message":"insufficient credits","code":"insufficient_credits","requiredCredits":20,"remainingCredits":0,"shortfallCredits":20}`)
	})
	defer srv.Close()

	_, _, err := c.Credits.Get(context.Background())
	if !errors.Is(err, ErrInsufficientCredits) {
		t.Errorf("expected ErrInsufficientCredits, got %v", err)
	}
	var icErr *InsufficientCreditsError
	if !errors.As(err, &icErr) {
		t.Fatalf("expected *InsufficientCreditsError, got %T", err)
	}
	if icErr.RequiredCredits != 20 {
		t.Errorf("expected requiredCredits 20, got %f", icErr.RequiredCredits)
	}
	if icErr.RemainingCredits != 0 {
		t.Errorf("expected remainingCredits 0, got %f", icErr.RemainingCredits)
	}
}

func TestErrServer_500(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		writeResponse(t, w, `{"message":"internal server error"}`)
	})
	defer srv.Close()

	_, _, err := c.Credits.Get(context.Background())
	if !errors.Is(err, ErrServer) {
		t.Errorf("expected ErrServer, got %v", err)
	}
}

func TestErrServer_503(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		writeResponse(t, w, `{"message":"service unavailable"}`)
	})
	defer srv.Close()

	_, _, err := c.Credits.Get(context.Background())
	if !errors.Is(err, ErrServer) {
		t.Errorf("expected ErrServer, got %v", err)
	}
}

func TestAsyncJobFailedError_Error(t *testing.T) {
	err := &AsyncJobFailedError{Message: "job failed", ErrorCode: "TIMEOUT"}
	want := "socialkit: async job failed: code=TIMEOUT error=job failed"
	if err.Error() != want {
		t.Errorf("expected %q, got %q", want, err.Error())
	}
}

func TestAsyncJobFailedError_Error_NoCode(t *testing.T) {
	err := &AsyncJobFailedError{Message: "job failed"}
	want := "socialkit: async job failed: job failed"
	if err.Error() != want {
		t.Errorf("expected %q, got %q", want, err.Error())
	}
}

func TestInsufficientCreditsError_Error(t *testing.T) {
	err := &InsufficientCreditsError{
		RequiredCredits:  10,
		RemainingCredits: 5,
		ShortfallCredits: 5,
	}
	if err.Error() == "" {
		t.Error("expected non-empty error message")
	}
}

func TestNewAPIError_UnknownStatus(t *testing.T) {
	err := newAPIError(418, "I'm a teapot", "", nil, nil)
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Errorf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 418 {
		t.Errorf("expected status 418, got %d", apiErr.StatusCode)
	}
}

func TestExtractMessageAndCode(t *testing.T) {
	msg, code := extractMessageAndCode([]byte(`{"message":"test error","code":"test_code"}`))
	if msg != "test error" {
		t.Errorf("expected message 'test error', got %s", msg)
	}
	if code != "test_code" {
		t.Errorf("expected code 'test_code', got %s", code)
	}
}

func TestExtractMessageAndCode_NestedData(t *testing.T) {
	msg, code := extractMessageAndCode([]byte(`{"data":{"message":"nested error","code":"nested_code"}}`))
	if msg != "nested error" {
		t.Errorf("expected message 'nested error', got %s", msg)
	}
	if code != "nested_code" {
		t.Errorf("expected code 'nested_code', got %s", code)
	}
}

func TestExtractMessageAndCode_EmptyBody(t *testing.T) {
	msg, _ := extractMessageAndCode(nil)
	if msg != "unknown error" {
		t.Errorf("expected 'unknown error', got %s", msg)
	}
}
