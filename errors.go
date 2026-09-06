package socialkit

import (
	"errors"
	"fmt"
	"time"
)

// APIError represents an error response returned by the SocialKit API.
// It is returned by every service method when the API responds with a
// non-2xx HTTP status code.
type APIError struct {
	// StatusCode is the HTTP status code returned by the API.
	StatusCode int
	// Message is a human readable error message, extracted from the
	// response body when possible, or a generic description otherwise.
	Message string
	// Code is the machine-readable error code from the response body, if
	// present (e.g. "insufficient_credits").
	Code string
	// RawBody contains the raw response body (redacted of any auth info)
	// for debugging purposes.
	RawBody []byte
	// Meta holds response metadata extracted from headers, if available.
	Meta *ResponseMeta
}

// Error implements the error interface for APIError.
func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("socialkit: api error: status=%d code=%s message=%s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("socialkit: api error: status=%d message=%s", e.StatusCode, e.Message)
}

// Sentinel errors that can be checked with errors.Is against errors
// returned from service methods. The underlying *APIError is always
// available via errors.As for accessing StatusCode, Code, RawBody, and
// Meta.
var (
	// ErrBadRequest is returned when the API responds with HTTP 400.
	ErrBadRequest = errors.New("socialkit: bad request")
	// ErrUnauthorized is returned when the API responds with HTTP 401.
	ErrUnauthorized = errors.New("socialkit: unauthorized")
	// ErrForbidden is returned when the API responds with HTTP 403
	// (without a specific insufficient_credits code).
	ErrForbidden = errors.New("socialkit: forbidden")
	// ErrInsufficientCredits is returned when the API responds with HTTP
	// 403 and the error code "insufficient_credits". The wrapped error
	// also implements the InsufficientCreditsError interface.
	ErrInsufficientCredits = errors.New("socialkit: insufficient credits")
	// ErrNotFound is returned when the API responds with HTTP 404.
	ErrNotFound = errors.New("socialkit: not found")
	// ErrRateLimited is returned when the API responds with HTTP 429 and
	// all retry attempts have been exhausted.
	ErrRateLimited = errors.New("socialkit: rate limited")
	// ErrServer is returned when the API responds with HTTP 500 or 503.
	ErrServer = errors.New("socialkit: server error")
	// ErrTransport is returned when the HTTP request fails at the
	// transport level (e.g. network error).
	ErrTransport = errors.New("socialkit: transport error")
	// ErrTimeout is returned when the request times out.
	ErrTimeout = errors.New("socialkit: timeout")
	// ErrDecode is returned when the response body cannot be decoded.
	ErrDecode = errors.New("socialkit: decode error")
	// ErrAsyncJobFailed is returned when a v2 async download job ends in
	// the "failed" state.
	ErrAsyncJobFailed = errors.New("socialkit: async job failed")
)

// InsufficientCreditsError provides detailed information about a credit
// shortage. It is attached to errors returned when the API responds with
// HTTP 403 and the code "insufficient_credits".
type InsufficientCreditsError struct {
	// RequiredCredits is the number of credits the operation needed.
	RequiredCredits float64
	// RemainingCredits is the number of credits the account has left.
	RemainingCredits float64
	// ShortfallCredits is the difference between required and remaining.
	ShortfallCredits float64
	// ResetAt is the time when monthly credits reset, if applicable.
	ResetAt *time.Time
	// UpgradeURL is a URL to upgrade the plan, if provided by the API.
	UpgradeURL string
	// TopUpURL is a URL to purchase additional credits, if provided.
	TopUpURL string
	// IncidentID is a reference ID for the failed request, if provided.
	IncidentID string
}

// Error implements the error interface for InsufficientCreditsError.
func (e *InsufficientCreditsError) Error() string {
	return fmt.Sprintf("socialkit: insufficient credits: required=%.0f remaining=%.0f shortfall=%.0f",
		e.RequiredCredits, e.RemainingCredits, e.ShortfallCredits)
}

// AsyncJobFailedError provides details about a failed v2 async download
// job. It is returned by DownloadsService.Wait when the job reaches the
// "failed" state.
type AsyncJobFailedError struct {
	// Message is the human-readable error message from the job.
	Message string
	// ErrorCode is the machine-readable error code, if any.
	ErrorCode string
	// Retryable indicates whether the job can be retried.
	Retryable bool
}

// Error implements the error interface for AsyncJobFailedError.
func (e *AsyncJobFailedError) Error() string {
	if e.ErrorCode != "" {
		return fmt.Sprintf("socialkit: async job failed: code=%s error=%s", e.ErrorCode, e.Message)
	}
	return fmt.Sprintf("socialkit: async job failed: %s", e.Message)
}

// Is reports whether this error matches target, allowing errors.Is to
// match against ErrAsyncJobFailed.
func (e *AsyncJobFailedError) Is(target error) bool {
	return target == ErrAsyncJobFailed
}

// wrappedAPIError wraps an *APIError together with a sentinel error so
// that both errors.As(&APIError{}) and errors.Is(sentinel) work on the
// value returned to the caller. It may also carry an
// InsufficientCreditsError for 403 insufficient_credits responses.
type wrappedAPIError struct {
	*APIError
	sentinel     error
	insufficient *InsufficientCreditsError
}

// Unwrap allows errors.Is/errors.As to traverse to the sentinel error.
func (w *wrappedAPIError) Unwrap() error {
	if w.insufficient != nil {
		return w.sentinel
	}
	return w.sentinel
}

// Is reports whether this error matches target, supporting comparison
// against the wrapped sentinel error.
func (w *wrappedAPIError) Is(target error) bool {
	return errors.Is(w.sentinel, target)
}

// As allows errors.As(err, &apiErr) to retrieve the underlying
// *APIError, and errors.As(err, &insufficientErr) to retrieve the
// *InsufficientCreditsError when present.
func (w *wrappedAPIError) As(target interface{}) bool {
	if ptr, ok := target.(**APIError); ok {
		*ptr = w.APIError
		return true
	}
	if w.insufficient != nil {
		if ptr, ok := target.(**InsufficientCreditsError); ok {
			*ptr = w.insufficient
			return true
		}
	}
	return false
}

// newAPIError builds an error value for the given status code, message,
// code, raw body, and metadata, attaching the appropriate sentinel
// error when the status code is a well-known one. When the code is
// "insufficient_credits" and the status is 403, an
// InsufficientCreditsError is also attached.
func newAPIError(statusCode int, message, code string, rawBody []byte, meta *ResponseMeta) error {
	// Redact any credential values from the raw body before storing.
	redactedBody := rawBody
	if len(rawBody) > 0 {
		redactedBody = []byte(redactJSON(rawBody))
	}

	apiErr := &APIError{
		StatusCode: statusCode,
		Message:    message,
		Code:       code,
		RawBody:    redactedBody,
		Meta:       meta,
	}

	switch {
	case statusCode == 400:
		return &wrappedAPIError{APIError: apiErr, sentinel: ErrBadRequest}
	case statusCode == 401:
		return &wrappedAPIError{APIError: apiErr, sentinel: ErrUnauthorized}
	case statusCode == 403 && code == "insufficient_credits":
		ic := parseInsufficientCredits(rawBody)
		return &wrappedAPIError{APIError: apiErr, sentinel: ErrInsufficientCredits, insufficient: ic}
	case statusCode == 403:
		return &wrappedAPIError{APIError: apiErr, sentinel: ErrForbidden}
	case statusCode == 404:
		return &wrappedAPIError{APIError: apiErr, sentinel: ErrNotFound}
	case statusCode == 429:
		return &wrappedAPIError{APIError: apiErr, sentinel: ErrRateLimited}
	case statusCode == 500 || statusCode == 503:
		return &wrappedAPIError{APIError: apiErr, sentinel: ErrServer}
	default:
		return apiErr
	}
}

// parseInsufficientCredits attempts to extract credit shortage details
// from a 403 insufficient_credits response body.
func parseInsufficientCredits(body []byte) *InsufficientCreditsError {
	var parsed struct {
		Data struct {
			RequiredCredits  *float64 `json:"requiredCredits"`
			RemainingCredits *float64 `json:"remainingCredits"`
			ShortfallCredits *float64 `json:"shortfallCredits"`
			ResetAt          *string  `json:"resetAt"`
			UpgradeURL       *string  `json:"upgradeUrl"`
			TopUpURL         *string  `json:"topUpUrl"`
			IncidentID       *string  `json:"incidentId"`
		} `json:"data"`
		RequiredCredits  *float64 `json:"requiredCredits"`
		RemainingCredits *float64 `json:"remainingCredits"`
		ShortfallCredits *float64 `json:"shortfallCredits"`
		ResetAt          *string  `json:"resetAt"`
		UpgradeURL       *string  `json:"upgradeUrl"`
		TopUpURL         *string  `json:"topUpUrl"`
		IncidentID       *string  `json:"incidentId"`
	}
	if err := jsonUnmarshal(body, &parsed); err != nil {
		return &InsufficientCreditsError{}
	}

	ic := &InsufficientCreditsError{}

	// Prefer nested data fields, fall back to top-level.
	if parsed.Data.RequiredCredits != nil {
		ic.RequiredCredits = *parsed.Data.RequiredCredits
	} else if parsed.RequiredCredits != nil {
		ic.RequiredCredits = *parsed.RequiredCredits
	}
	if parsed.Data.RemainingCredits != nil {
		ic.RemainingCredits = *parsed.Data.RemainingCredits
	} else if parsed.RemainingCredits != nil {
		ic.RemainingCredits = *parsed.RemainingCredits
	}
	if parsed.Data.ShortfallCredits != nil {
		ic.ShortfallCredits = *parsed.Data.ShortfallCredits
	} else if parsed.ShortfallCredits != nil {
		ic.ShortfallCredits = *parsed.ShortfallCredits
	}
	if parsed.Data.ResetAt != nil && *parsed.Data.ResetAt != "" {
		if t, err := time.Parse(time.RFC3339, *parsed.Data.ResetAt); err == nil {
			ic.ResetAt = &t
		}
	} else if parsed.ResetAt != nil && *parsed.ResetAt != "" {
		if t, err := time.Parse(time.RFC3339, *parsed.ResetAt); err == nil {
			ic.ResetAt = &t
		}
	}
	if parsed.Data.UpgradeURL != nil {
		ic.UpgradeURL = *parsed.Data.UpgradeURL
	} else if parsed.UpgradeURL != nil {
		ic.UpgradeURL = *parsed.UpgradeURL
	}
	if parsed.Data.TopUpURL != nil {
		ic.TopUpURL = *parsed.Data.TopUpURL
	} else if parsed.TopUpURL != nil {
		ic.TopUpURL = *parsed.TopUpURL
	}
	if parsed.Data.IncidentID != nil {
		ic.IncidentID = *parsed.Data.IncidentID
	} else if parsed.IncidentID != nil {
		ic.IncidentID = *parsed.IncidentID
	}

	return ic
}
