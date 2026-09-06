package socialkit

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"
	"time"
)

func TestDownloadsService_Start(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/youtube/download" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		writeResponse(t, w, `{"success":true,"data":{"jobId":"job123","status":"queued","platform":"youtube","url":"https://youtube.com/watch?v=test","quality":"720p","format":"mp4"}}`)
	})
	defer srv.Close()

	job, _, err := c.Downloads.Start(context.Background(), "youtube", &V2DownloadRequest{
		URL:     "https://youtube.com/watch?v=test",
		Format:  "mp4",
		Quality: "720p",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if job.JobID != "job123" {
		t.Errorf("expected jobId job123, got %s", job.JobID)
	}
	if job.Status != "queued" {
		t.Errorf("expected status queued, got %s", job.Status)
	}
}

func TestDownloadsService_Start_InvalidPlatform(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not make a request")
	})
	defer srv.Close()

	_, _, err := c.Downloads.Start(context.Background(), "invalid", &V2DownloadRequest{
		URL: "https://example.com",
	})
	if err == nil {
		t.Fatal("expected error for invalid platform")
	}
}

func TestDownloadsService_Get(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/downloads/job123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		writeResponse(t, w, `{"success":true,"data":{"jobId":"job123","status":"ready","platform":"youtube","url":"https://youtube.com/watch?v=test","quality":"720p","format":"mp4","billed":true,"downloadUrl":"https://download.socialkit.dev/file.mp4","expiresIn":3600.0,"durationSeconds":600.0,"fileSize":10485760,"fileSizeMB":"10.0","creditsCost":2,"title":"Test Video","thumbnail":"https://img.com/t.jpg"}}`)
	})
	defer srv.Close()

	job, _, err := c.Downloads.Get(context.Background(), "job123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if job.Status != "ready" {
		t.Errorf("expected status ready, got %s", job.Status)
	}
	if job.DownloadURL != "https://download.socialkit.dev/file.mp4" {
		t.Errorf("unexpected downloadUrl: %s", job.DownloadURL)
	}
	if job.Billed == nil || !*job.Billed {
		t.Error("expected billed true")
	}
}

func TestDownloadsService_Get_EmptyJobID(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not make a request")
	})
	defer srv.Close()

	_, _, err := c.Downloads.Get(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty jobID")
	}
}

func TestDownloadsService_Wait_Ready(t *testing.T) {
	var callCount int32
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&callCount, 1)
		if n == 1 {
			writeResponse(t, w, `{"success":true,"data":{"jobId":"job123","status":"processing","platform":"youtube"}}`)
			return
		}
		writeResponse(t, w, `{"success":true,"data":{"jobId":"job123","status":"ready","platform":"youtube","downloadUrl":"https://download.socialkit.dev/file.mp4"}}`)
	})
	defer srv.Close()

	job, _, err := c.Downloads.Wait(context.Background(), "job123", &WaitOptions{
		Interval: 10 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if job.Status != "ready" {
		t.Errorf("expected status ready, got %s", job.Status)
	}
	if callCount != 2 {
		t.Errorf("expected 2 calls, got %d", callCount)
	}
}

func TestDownloadsService_Wait_Failed(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeResponse(t, w, `{"success":true,"data":{"jobId":"job123","status":"failed","platform":"youtube","error":"Video too long","errorCode":"MAX_DURATION_EXCEEDED","retryable":false}}`)
	})
	defer srv.Close()

	job, _, err := c.Downloads.Wait(context.Background(), "job123", &WaitOptions{
		Interval: 10 * time.Millisecond,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrAsyncJobFailed) {
		t.Errorf("expected ErrAsyncJobFailed, got %v", err)
	}
	var asyncErr *AsyncJobFailedError
	if !errors.As(err, &asyncErr) {
		t.Errorf("expected *AsyncJobFailedError, got %T", err)
	} else {
		if asyncErr.Message != "Video too long" {
			t.Errorf("expected message 'Video too long', got %s", asyncErr.Message)
		}
		if asyncErr.ErrorCode != "MAX_DURATION_EXCEEDED" {
			t.Errorf("expected errorCode MAX_DURATION_EXCEEDED, got %s", asyncErr.ErrorCode)
		}
		if asyncErr.Retryable {
			t.Error("expected retryable false")
		}
	}
	if job == nil || job.Status != "failed" {
		t.Errorf("expected job with status failed, got %v", job)
	}
}

func TestDownloadsService_Wait_ContextDeadline(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeResponse(t, w, `{"success":true,"data":{"jobId":"job123","status":"processing","platform":"youtube"}}`)
	})
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	job, _, err := c.Downloads.Wait(ctx, "job123", &WaitOptions{
		Interval: 30 * time.Millisecond,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected context.DeadlineExceeded, got %v", err)
	}
	if job == nil || job.Status != "processing" {
		t.Errorf("expected last job with status processing, got %v", job)
	}
}

func TestDownloadsService_Wait_MaxAttempts(t *testing.T) {
	var callCount int32
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		writeResponse(t, w, `{"success":true,"data":{"jobId":"job123","status":"processing","platform":"youtube"}}`)
	})
	defer srv.Close()

	_, _, err := c.Downloads.Wait(context.Background(), "job123", &WaitOptions{
		Interval:      5 * time.Millisecond,
		MaxAttempts:   3,
		BackoffFactor: 1.0,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if callCount != 3 {
		t.Errorf("expected 3 calls, got %d", callCount)
	}
}

func TestDownloadsService_Wait_DefaultOptions(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeResponse(t, w, `{"success":true,"data":{"jobId":"job123","status":"ready","platform":"youtube"}}`)
	})
	defer srv.Close()

	job, _, err := c.Downloads.Wait(context.Background(), "job123", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if job.Status != "ready" {
		t.Errorf("expected status ready, got %s", job.Status)
	}
}
