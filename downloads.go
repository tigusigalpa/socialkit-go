package socialkit

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// DownloadsService provides access to the /v2/* download endpoints of
// the SocialKit API, covering async download job start, get, and wait
// operations.
type DownloadsService struct {
	client *Client
}

// Start initiates a v2 async download job for the given platform.
// Platform must be one of: "youtube", "tiktok", "instagram", "facebook".
// POST is the recommended high-level path; GET parity is preserved in
// the advanced API.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/async-download-api
func (s *DownloadsService) Start(ctx context.Context, platform string, req *V2DownloadRequest) (*V2Job, *ResponseMeta, error) {
	switch platform {
	case "youtube", "tiktok", "instagram", "facebook":
	default:
		return nil, nil, fmt.Errorf("socialkit: unsupported platform %q (must be youtube, tiktok, instagram, or facebook)", platform)
	}
	path := fmt.Sprintf("/v2/%s/download", escapePathSegment(platform))
	var out V2Job
	meta, err := s.client.doRequestWithMeta(ctx, "POST", path, nil, req, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Get retrieves the current state of a v2 async download job by its
// job ID. This endpoint always uses GET.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/async-download-api
func (s *DownloadsService) Get(ctx context.Context, jobID string) (*V2Job, *ResponseMeta, error) {
	if jobID == "" {
		return nil, nil, errors.New("socialkit: jobID is required")
	}
	path := fmt.Sprintf("/v2/downloads/%s", escapePathSegment(jobID))
	var out V2Job
	meta, err := s.client.doRequestWithMeta(ctx, "GET", path, nil, nil, &out)
	if err != nil {
		return nil, meta, err
	}
	return &out, meta, nil
}

// Wait polls a v2 async download job until it reaches a terminal state
// ("ready" or "failed") or the context deadline is reached. The polling
// interval and backoff are controlled by opts. If opts is nil, sensible
// defaults are used (2-second interval, no backoff, no max attempts).
//
// When the job reaches the "failed" state, an *AsyncJobFailedError is
// returned (also matching ErrAsyncJobFailed via errors.Is). When the
// context deadline is reached, the last known job state is returned
// alongside the context error.
//
// Upstream docs: https://docs.socialkit.dev/api-reference/async-download-api
func (s *DownloadsService) Wait(ctx context.Context, jobID string, opts *WaitOptions) (*V2Job, *ResponseMeta, error) {
	interval := 2 * time.Second
	maxAttempts := 0
	backoffFactor := 1.0
	if opts != nil {
		if opts.Interval > 0 {
			interval = opts.Interval
		}
		if opts.MaxAttempts > 0 {
			maxAttempts = opts.MaxAttempts
		}
		if opts.BackoffFactor > 0 {
			backoffFactor = opts.BackoffFactor
		}
	}

	currentInterval := interval
	var lastJob *V2Job
	var lastMeta *ResponseMeta

	for attempt := 1; maxAttempts == 0 || attempt <= maxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			if lastJob != nil {
				return lastJob, lastMeta, err
			}
			return nil, lastMeta, err
		}

		job, meta, err := s.Get(ctx, jobID)
		if err != nil {
			// If the context was cancelled, return what we have.
			if ctx.Err() != nil {
				return lastJob, lastMeta, ctx.Err()
			}
			return lastJob, meta, err
		}
		lastJob = job
		lastMeta = meta

		switch job.Status {
		case "ready":
			return job, meta, nil
		case "failed":
			retryable := false
			if job.Retryable != nil {
				retryable = *job.Retryable
			}
			return job, meta, &AsyncJobFailedError{
				Message:   job.Error,
				ErrorCode: job.ErrorCode,
				Retryable: retryable,
			}
		}

		// Wait before the next poll, respecting context cancellation.
		timer := time.NewTimer(currentInterval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return job, meta, ctx.Err()
		case <-timer.C:
		}

		// Apply backoff for the next interval.
		if backoffFactor > 1 {
			currentInterval = time.Duration(float64(currentInterval) * backoffFactor)
		}
	}

	// Exhausted maxAttempts without reaching a terminal state.
	if lastJob != nil {
		return lastJob, lastMeta, fmt.Errorf("socialkit: wait exhausted after %d attempts, job status is %q", maxAttempts, lastJob.Status)
	}
	return nil, lastMeta, fmt.Errorf("socialkit: wait exhausted after %d attempts", maxAttempts)
}
