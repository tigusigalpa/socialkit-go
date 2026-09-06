# SocialKit Golang Client/SDK/Library

![SocialKit Golang client sdk library](https://i.postimg.cc/pdYjbmLX/socialkit-go-hero-banner.jpg)

[![CI](https://github.com/tigusigalpa/socialkit-go/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/tigusigalpa/socialkit-go/actions/workflows/ci.yml)
[![Tests](https://github.com/tigusigalpa/socialkit-go/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/tigusigalpa/socialkit-go/actions/workflows/test.yml)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green?style=flat-square)](LICENSE)
[![CodeQL](https://github.com/tigusigalpa/socialkit-go/actions/workflows/codeql.yml/badge.svg?branch=main)](https://github.com/tigusigalpa/socialkit-go/actions/workflows/codeql.yml)
[![Codecov](https://codecov.io/gh/tigusigalpa/socialkit-go/graph/badge.svg)](https://codecov.io/gh/tigusigalpa/socialkit-go)
[![GitHub Release](https://img.shields.io/github/v/release/tigusigalpa/socialkit-go?style=flat-square)](https://github.com/tigusigalpa/socialkit-go/releases)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue?style=flat-square&logo=go)](https://pkg.go.dev/github.com/tigusigalpa/socialkit-go)

A practical, typed Go SDK for the [SocialKit](https://socialkit.dev) REST API.
It keeps authentication, response envelopes, pagination, rate-limit metadata,
and typed API errors out of your application code. Every request accepts a
`context.Context`, so cancellation and deadlines work the Go way.

## Installation

```sh
go get github.com/tigusigalpa/socialkit-go
```

## Examples you can run

The examples use `SOCIALKIT_ACCESS_KEY` from the environment; they never put a
key in source code. Each is a small program you can copy from or run directly:

```sh
set SOCIALKIT_ACCESS_KEY=YOUR_SOCIALKIT_ACCESS_KEY
go run ./examples/basic       # status + a transcript
go run ./examples/pagination  # cursor-based TikTok comments
go run ./examples/download    # v2 async download and polling
go run ./examples/errors      # errors.Is / errors.As patterns
```

The URLs are public placeholders. Before making a billable request, substitute
a URL you are authorized to process.

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    socialkit "github.com/tigusigalpa/socialkit-go"
)

func main() {
    client := socialkit.NewClient(os.Getenv("SOCIALKIT_ACCESS_KEY"))

    transcript, meta, err := client.YouTube.Transcript(context.Background(), &socialkit.FetchRequest{
        URL: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(transcript.Transcript)
    if meta.CreditsUsed != nil {
        fmt.Printf("Credits used: %.0f\n", *meta.CreditsUsed)
    }
}
```

## Secure Configuration

The SDK sends your access key via the `x-access-key` HTTP header by default.
The key is **never** placed in URLs, query strings, or request bodies unless
you explicitly opt in with `WithKeyInQuery(true)`.

```go
client := socialkit.NewClient(
    os.Getenv("SOCIALKIT_ACCESS_KEY"),
    socialkit.WithTimeout(30*time.Second),
    socialkit.WithRetry(3, time.Second),
    socialkit.WithUserAgent("MyApp/1.0"),
)
```

### Key Redaction

All error output, response metadata, and raw bodies are automatically
redacted — any credential values are replaced with `[REDACTED]`. The key is
never logged, serialized, or exposed in error structs.

## Service Map

| Service | Field | Endpoints |
|---------|-------|-----------|
| Status | `client.Status` | `/status` |
| Credits | `client.Credits` | `/credits` |
| YouTube | `client.YouTube` | `/youtube/*` |
| TikTok | `client.TikTok` | `/tiktok/*` |
| Instagram | `client.Instagram` | `/instagram/*` |
| Facebook | `client.Facebook` | `/facebook/*` |
| Twitter | `client.Twitter` | `/twitter/*` |
| LinkedIn | `client.LinkedIn` | `/linkedin/*` |
| Video | `client.Video` | `/video/*` |
| Downloads | `client.Downloads` | `/v2/*` |

## Options

| Option | Description |
|--------|-------------|
| `WithBaseURL(url)` | Override the default base URL |
| `WithHTTPClient(c)` | Use a custom `*http.Client` |
| `WithTimeout(d)` | Set the HTTP client timeout |
| `WithRetry(n, d)` | Configure retry attempts and backoff |
| `WithKeyInQuery(b)` | Enable legacy `access_key` in query/body |
| `WithUserAgent(ua)` | Override the User-Agent header |

## Method Usage

### Status

```go
status, meta, err := client.Status.Get(ctx)
```
- Upstream docs: https://docs.socialkit.dev/api-reference/status-api

### Credits

```go
credits, meta, err := client.Credits.Get(ctx)
```
- Upstream docs: https://docs.socialkit.dev/api-reference/credits-api

### YouTube

```go
// Transcript
transcript, meta, err := client.YouTube.Transcript(ctx, &socialkit.FetchRequest{
    URL: "https://youtube.com/watch?v=...",
})

// Summary with custom prompt
summary, meta, err := client.YouTube.Summarize(ctx, &socialkit.SummaryRequest{
    URL:          "https://youtube.com/watch?v=...",
    CustomPrompt: socialkit.StringPtr("Focus on technical details"),
})

// Stats
stats, meta, err := client.YouTube.Stats(ctx, &socialkit.FetchRequest{URL: "..."})

// Comments (limit max 100, sortBy: top|new)
comments, meta, err := client.YouTube.Comments(ctx, &socialkit.CommentsRequest{
    URL:    "https://youtube.com/watch?v=...",
    Limit:  socialkit.IntPtr(50),
    SortBy: socialkit.StringPtr("top"),
})

// Channel stats
channelStats, meta, err := client.YouTube.ChannelStats(ctx, &socialkit.FetchRequest{URL: "..."})

// Search (sortBy: relevance|date|views|rating, uploadDate: hour|today|week|month|year, type: video|shorts)
results, meta, err := client.YouTube.Search(ctx, &socialkit.SearchRequest{
    Query:  "golang tutorial",
    SortBy: socialkit.StringPtr("relevance"),
})

// Videos (fullDetails, limit capped at 30)
videos, meta, err := client.YouTube.Videos(ctx, &socialkit.VideosRequest{
    URL:         "https://youtube.com/channel/...",
    FullDetails: socialkit.BoolPtr(true),
})

// Download (max 10MB, URL expires 1h)
result, meta, err := client.YouTube.Download(ctx, &socialkit.DownloadRequest{
    URL:     "https://youtube.com/watch?v=...",
    Format:  socialkit.StringPtr("mp4"),
    Quality: socialkit.StringPtr("720p"),
})

// Bulk (experimental — raw JSON in/out)
rawResult, meta, err := client.YouTube.BulkTranscript(ctx, json.RawMessage(`{"urls":[...]}`))
```

- Upstream docs: https://docs.socialkit.dev/api-reference/youtube-transcript-api, https://docs.socialkit.dev/api-reference/youtube-summarize-api, https://docs.socialkit.dev/api-reference/youtube-stats-api, https://docs.socialkit.dev/api-reference/youtube-comments-api, https://docs.socialkit.dev/api-reference/youtube-channel-stats-api, https://docs.socialkit.dev/api-reference/youtube-search-api, https://docs.socialkit.dev/api-reference/youtube-videos-api, https://docs.socialkit.dev/api-reference/youtube-download-api

### TikTok

```go
transcript, _, _ := client.TikTok.Transcript(ctx, &socialkit.FetchRequest{URL: "..."})
summary, _, _ := client.TikTok.Summarize(ctx, &socialkit.SummaryRequest{URL: "..."})
stats, _, _ := client.TikTok.Stats(ctx, &socialkit.FetchRequest{URL: "..."})
comments, _, _ := client.TikTok.Comments(ctx, &socialkit.CommentsRequest{URL: "..."})  // cursor+hasMore
channelStats, _, _ := client.TikTok.ChannelStats(ctx, &socialkit.FetchRequest{URL: "..."})
channelVideos, _, _ := client.TikTok.ChannelVideos(ctx, &socialkit.ListRequest{URL: "..."})  // default 30, max 100
searchResults, _, _ := client.TikTok.Search(ctx, &socialkit.SearchRequest{Query: "..."})  // max 100
hashtagResults, _, _ := client.TikTok.HashtagSearch(ctx, &socialkit.HashtagSearchRequest{Hashtag: "dance"})  // no #, max 100
downloadResult, _, _ := client.TikTok.Download(ctx, &socialkit.DownloadRequest{URL: "..."})  // max 30MB, expires 1h
```

- Upstream docs: https://docs.socialkit.dev/api-reference/tiktok-transcript-api, https://docs.socialkit.dev/api-reference/tiktok-summarize-api, https://docs.socialkit.dev/api-reference/tiktok-stats-api, https://docs.socialkit.dev/api-reference/tiktok-comments-api, https://docs.socialkit.dev/api-reference/tiktok-channel-stats-api, https://docs.socialkit.dev/api-reference/tiktok-channel-videos-api, https://docs.socialkit.dev/api-reference/tiktok-search-api, https://docs.socialkit.dev/api-reference/tiktok-hashtag-search-api, https://docs.socialkit.dev/api-reference/tiktok-download-api

### Instagram

```go
transcript, _, _ := client.Instagram.Transcript(ctx, &socialkit.FetchRequest{URL: "..."})
summary, _, _ := client.Instagram.Summarize(ctx, &socialkit.SummaryRequest{URL: "..."})
stats, _, _ := client.Instagram.Stats(ctx, &socialkit.FetchRequest{URL: "..."})
comments, _, _ := client.Instagram.Comments(ctx, &socialkit.CommentsRequest{URL: "..."})
channelStats, _, _ := client.Instagram.ChannelStats(ctx, &socialkit.FetchRequest{URL: "..."})
posts, _, _ := client.Instagram.ChannelPosts(ctx, &socialkit.ListRequest{URL: "..."})  // max 100
reels, _, _ := client.Instagram.ChannelReels(ctx, &socialkit.ListRequest{URL: "..."})  // can take 60s
reelsSearch, _, _ := client.Instagram.ReelsSearch(ctx, &socialkit.ReelsSearchRequest{Query: "..."})  // page must be 1
downloadResult, _, _ := client.Instagram.Download(ctx, &socialkit.DownloadRequest{URL: "..."})  // expires 1h
```

- Upstream docs: https://docs.socialkit.dev/api-reference/instagram-transcript-api, https://docs.socialkit.dev/api-reference/instagram-summarize-api, https://docs.socialkit.dev/api-reference/instagram-stats-api, https://docs.socialkit.dev/api-reference/instagram-comments-api, https://docs.socialkit.dev/api-reference/instagram-channel-stats-api, https://docs.socialkit.dev/api-reference/instagram-channel-posts-api, https://docs.socialkit.dev/api-reference/instagram-channel-reels-api, https://docs.socialkit.dev/api-reference/instagram-reels-search-api, https://docs.socialkit.dev/api-reference/instagram-download-api

### Facebook

```go
transcript, _, _ := client.Facebook.Transcript(ctx, &socialkit.FetchRequest{URL: "..."})
summary, _, _ := client.Facebook.Summarize(ctx, &socialkit.SummaryRequest{URL: "..."})
stats, _, _ := client.Facebook.Stats(ctx, &socialkit.FetchRequest{URL: "..."})  // includes reactions[]
comments, _, _ := client.Facebook.Comments(ctx, &socialkit.CommentsRequest{URL: "..."})  // commentCount nullable
channelStats, _, _ := client.Facebook.ChannelStats(ctx, &socialkit.FetchRequest{URL: "..."})
```

- Upstream docs: https://docs.socialkit.dev/api-reference/facebook-transcript-api, https://docs.socialkit.dev/api-reference/facebook-summarize-api, https://docs.socialkit.dev/api-reference/facebook-stats-api, https://docs.socialkit.dev/api-reference/facebook-comments-api, https://docs.socialkit.dev/api-reference/facebook-channel-stats-api

### Twitter

```go
profile, _, _ := client.Twitter.Profile(ctx, &socialkit.FetchRequest{URL: "..."})
tweets, _, _ := client.Twitter.Tweets(ctx, &socialkit.TweetsRequest{URL: "..."})  // nextCursor nullable
tweet, _, _ := client.Twitter.Tweet(ctx, &socialkit.FetchRequest{URL: "..."})
thread, _, _ := client.Twitter.Thread(ctx, &socialkit.FetchRequest{URL: "..."})
transcript, _, _ := client.Twitter.Transcript(ctx, &socialkit.FetchRequest{URL: "..."})
```

- Upstream docs: https://docs.socialkit.dev/api-reference/twitter-profile-api, https://docs.socialkit.dev/api-reference/twitter-tweets-api, https://docs.socialkit.dev/api-reference/twitter-tweet-api, https://docs.socialkit.dev/api-reference/twitter-thread-api, https://docs.socialkit.dev/api-reference/twitter-transcript-api

### LinkedIn

```go
profile, _, _ := client.LinkedIn.Profile(ctx, &socialkit.FetchRequest{URL: "..."})  // recentArticles[]
company, _, _ := client.LinkedIn.Company(ctx, &socialkit.FetchRequest{URL: "..."})
companyPosts, _, _ := client.LinkedIn.CompanyPosts(ctx, &socialkit.CompanyPostsRequest{URL: "..."})  // max 50
post, _, _ := client.LinkedIn.Post(ctx, &socialkit.FetchRequest{URL: "..."})  // untyped json.RawMessage
transcript, _, _ := client.LinkedIn.Transcript(ctx, &socialkit.FetchRequest{URL: "..."})  // segments without timestamp
```

- Upstream docs: https://docs.socialkit.dev/api-reference/linkedin-profile-api, https://docs.socialkit.dev/api-reference/linkedin-company-api, https://docs.socialkit.dev/api-reference/linkedin-company-posts-api, https://docs.socialkit.dev/api-reference/linkedin-post-api, https://docs.socialkit.dev/api-reference/linkedin-transcript-api

### Video

```go
// Direct video file URL, 5 credits/started minute
transcript, _, _ := client.Video.Transcript(ctx, &socialkit.FetchRequest{URL: "https://example.com/video.mp4"})
summary, _, _ := client.Video.Summarize(ctx, &socialkit.SummaryRequest{URL: "https://example.com/video.mp4"})
```

- Upstream docs: https://docs.socialkit.dev/api-reference/video-transcript-api, https://docs.socialkit.dev/api-reference/video-summarize-api

### Downloads (V2 Async)

```go
// Start a download job (platform: youtube|tiktok|instagram|facebook)
job, _, err := client.Downloads.Start(ctx, "youtube", &socialkit.V2DownloadRequest{
    URL:     "https://youtube.com/watch?v=...",
    Format:  "mp4",
    Quality: "720p",
})

// Get job status (always GET)
job, _, err = client.Downloads.Get(ctx, job.JobID)

// Wait for completion with polling
finalJob, _, err := client.Downloads.Wait(ctx, job.JobID, &socialkit.WaitOptions{
    Interval:      3 * time.Second,
    MaxAttempts:   20,
    BackoffFactor: 1.5,
})
```

- Upstream docs: https://docs.socialkit.dev/api-reference/async-download-api

## Custom Summary Response

The `Summary` type captures unknown fields in an `Extra` map for forward
compatibility with `custom_response` fields:

```go
summary, _, _ := client.YouTube.Summarize(ctx, &socialkit.SummaryRequest{
    URL:            "https://youtube.com/watch?v=...",
    CustomResponse: map[string]interface{}{
        "sentiment":  "",
        "keyTakeaway": "",
    },
})

// Access custom fields via Extra
if raw, ok := summary.Extra["sentiment"]; ok {
    fmt.Println(string(raw))
}
```

## Pagination

Cursor-based pagination is supported via `cursor` and `hasMore` fields:

```go
// First page
page1, _, _ := client.TikTok.Comments(ctx, &socialkit.CommentsRequest{
    URL: "https://tiktok.com/@user/video/123",
})

// Next page using cursor
if page1.HasMore != nil && *page1.HasMore && page1.Cursor != nil {
    page2, _, _ := client.TikTok.Comments(ctx, &socialkit.CommentsRequest{
        URL:    "https://tiktok.com/@user/video/123",
        Cursor: page1.Cursor,
    })
}
```

## Transient Errors and Retries

By default, the SDK does not retry. Enable retries with `WithRetry`:

```go
client := socialkit.NewClient(key,
    socialkit.WithRetry(3, time.Second),
)
```

- Retries are attempted for HTTP 429 and 5xx responses.
- Exponential backoff with jitter is used, honoring `Retry-After` headers.
- HTTP 400, 401, 403, and 404 are never retried.
- `WithRetry` is explicit opt-in, including for POST requests. Because a POST
  can be billable, enable it only for operations you are comfortable repeating.

## Credits Metadata

Every response includes a `*ResponseMeta` with credit and rate-limit headers:

```go
_, meta, _ := client.YouTube.Transcript(ctx, req)

if meta.CreditsUsed != nil {
    fmt.Printf("Credits used: %.0f\n", *meta.CreditsUsed)
}
if meta.CreditsRemaining != nil {
    fmt.Printf("Credits remaining: %.0f\n", *meta.CreditsRemaining)
}
if meta.RateLimitRemaining != nil {
    fmt.Printf("Rate limit remaining: %d\n", *meta.RateLimitRemaining)
}
```

## V2 Async Jobs

V2 download jobs have four states: `queued`, `processing`, `ready`, `failed`.

The `Wait` method polls with caller-controlled interval and backoff,
respecting context deadlines. It returns `AsyncJobFailedError` when the job
reaches the `failed` state.

```go
finalJob, _, err := client.Downloads.Wait(ctx, jobID, &socialkit.WaitOptions{
    Interval:    2 * time.Second,
    MaxAttempts: 10,
})
if errors.Is(err, socialkit.ErrAsyncJobFailed) {
    var asyncErr *socialkit.AsyncJobFailedError
    if errors.As(err, &asyncErr) {
        log.Printf("job failed: %s (retryable: %v)", asyncErr.Message, asyncErr.Retryable)
    }
}
```

## Error Handling

The SDK provides a typed error hierarchy with `errors.Is` and `errors.As`
support:

```go
_, _, err := client.YouTube.Transcript(ctx, req)
if err != nil {
    switch {
    case errors.Is(err, socialkit.ErrUnauthorized):
        // 401 — invalid key
    case errors.Is(err, socialkit.ErrInsufficientCredits):
        var ic *socialkit.InsufficientCreditsError
        if errors.As(err, &ic) {
            fmt.Printf("Need %.0f more credits\n", ic.ShortfallCredits)
        }
    case errors.Is(err, socialkit.ErrRateLimited):
        // 429 — retries exhausted
    case errors.Is(err, socialkit.ErrNotFound):
        // 404
    }

    // Access the raw API error
    var apiErr *socialkit.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("HTTP %d: %s (code: %s)\n", apiErr.StatusCode, apiErr.Message, apiErr.Code)
    }
}
```

### Sentinel Errors

| Error | HTTP Status |
|-------|-------------|
| `ErrBadRequest` | 400 |
| `ErrUnauthorized` | 401 |
| `ErrForbidden` | 403 |
| `ErrInsufficientCredits` | 403 (code: insufficient_credits) |
| `ErrNotFound` | 404 |
| `ErrRateLimited` | 429 |
| `ErrServer` | 500, 503 |
| `ErrTransport` | Network error |
| `ErrTimeout` | Timeout |
| `ErrDecode` | JSON decode error |
| `ErrAsyncJobFailed` | V2 job failed |

## Testing

Tests use `httptest` and require no network or real API key:

```sh
go test ./... -v
```

## Compatibility

- Go 1.21+
- The SDK uses only the Go standard library (no external dependencies).
- `json.RawMessage` is used for forward-compatible and untyped fields.

## Security

- The access key is sent via the `x-access-key` header by default.
- All error output is redacted — credentials are replaced with `[REDACTED]`.
- See [SECURITY.md](SECURITY.md) for responsible disclosure.

## Endpoint Documentation

### Service & Async Downloads

- [Status API](https://docs.socialkit.dev/api-reference/status-api) — `GET /status`
- [Credits API](https://docs.socialkit.dev/api-reference/credits-api) — `GET /credits`
- [Async Download API v2](https://docs.socialkit.dev/api-reference/async-download-api) — `GET, POST /v2/{platform}/download`, `GET /v2/downloads/{jobId}`

### YouTube

- [Transcript](https://docs.socialkit.dev/api-reference/youtube-transcript-api) — `GET, POST /youtube/transcript`
- [Summary](https://docs.socialkit.dev/api-reference/youtube-summarize-api) — `GET, POST /youtube/summarize`
- [Stats](https://docs.socialkit.dev/api-reference/youtube-stats-api) — `GET, POST /youtube/stats`
- [Comments](https://docs.socialkit.dev/api-reference/youtube-comments-api) — `GET, POST /youtube/comments`
- [Channel Stats](https://docs.socialkit.dev/api-reference/youtube-channel-stats-api) — `GET, POST /youtube/channel-stats`
- [Search](https://docs.socialkit.dev/api-reference/youtube-search-api) — `GET, POST /youtube/search`
- [Videos](https://docs.socialkit.dev/api-reference/youtube-videos-api) — `GET, POST /youtube/videos`
- [Download](https://docs.socialkit.dev/api-reference/youtube-download-api) — `GET, POST /youtube/download`
- Bulk: `POST /youtube/transcript/bulk`, `POST /youtube/comments/bulk`, `POST /youtube/stats/bulk`, `POST /youtube/summarize/bulk` — [OpenAPI](https://api.socialkit.dev/openapi.json), [API Reference](https://docs.socialkit.dev/api-reference)

### TikTok

- [Transcript](https://docs.socialkit.dev/api-reference/tiktok-transcript-api) — `GET, POST /tiktok/transcript`
- [Summary](https://docs.socialkit.dev/api-reference/tiktok-summarize-api) — `GET, POST /tiktok/summarize`
- [Stats](https://docs.socialkit.dev/api-reference/tiktok-stats-api) — `GET, POST /tiktok/stats`
- [Comments](https://docs.socialkit.dev/api-reference/tiktok-comments-api) — `GET, POST /tiktok/comments`
- [Channel Stats](https://docs.socialkit.dev/api-reference/tiktok-channel-stats-api) — `GET, POST /tiktok/channel-stats`
- [Channel Videos](https://docs.socialkit.dev/api-reference/tiktok-channel-videos-api) — `GET, POST /tiktok/channel-videos`
- [Search](https://docs.socialkit.dev/api-reference/tiktok-search-api) — `GET, POST /tiktok/search`
- [Hashtag Search](https://docs.socialkit.dev/api-reference/tiktok-hashtag-search-api) — `GET, POST /tiktok/hashtag-search`
- [Download](https://docs.socialkit.dev/api-reference/tiktok-download-api) — `GET, POST /tiktok/download`
- Bulk: `POST /tiktok/transcript/bulk`, `POST /tiktok/comments/bulk`, `POST /tiktok/stats/bulk`, `POST /tiktok/channel-stats/bulk`, `POST /tiktok/summarize/bulk` — [OpenAPI](https://api.socialkit.dev/openapi.json), [API Reference](https://docs.socialkit.dev/api-reference)

### Instagram

- [Transcript](https://docs.socialkit.dev/api-reference/instagram-transcript-api) — `GET, POST /instagram/transcript`
- [Summary](https://docs.socialkit.dev/api-reference/instagram-summarize-api) — `GET, POST /instagram/summarize`
- [Stats](https://docs.socialkit.dev/api-reference/instagram-stats-api) — `GET, POST /instagram/stats`
- [Comments](https://docs.socialkit.dev/api-reference/instagram-comments-api) — `GET, POST /instagram/comments`
- [Channel Stats](https://docs.socialkit.dev/api-reference/instagram-channel-stats-api) — `GET, POST /instagram/channel-stats`
- [Channel Posts](https://docs.socialkit.dev/api-reference/instagram-channel-posts-api) — `GET, POST /instagram/channel-posts`
- [Channel Reels](https://docs.socialkit.dev/api-reference/instagram-channel-reels-api) — `GET, POST /instagram/channel-reels`
- [Reels Search](https://docs.socialkit.dev/api-reference/instagram-reels-search-api) — `GET, POST /instagram/reels-search`
- [Download](https://docs.socialkit.dev/api-reference/instagram-download-api) — `GET, POST /instagram/download`
- Bulk: `POST /instagram/transcript/bulk`, `POST /instagram/stats/bulk`, `POST /instagram/channel-stats/bulk`, `POST /instagram/summarize/bulk` — [OpenAPI](https://api.socialkit.dev/openapi.json), [API Reference](https://docs.socialkit.dev/api-reference)

### Facebook

- [Transcript](https://docs.socialkit.dev/api-reference/facebook-transcript-api) — `GET, POST /facebook/transcript`
- [Summary](https://docs.socialkit.dev/api-reference/facebook-summarize-api) — `GET, POST /facebook/summarize`
- [Stats](https://docs.socialkit.dev/api-reference/facebook-stats-api) — `GET, POST /facebook/stats`
- [Comments](https://docs.socialkit.dev/api-reference/facebook-comments-api) — `GET, POST /facebook/comments`
- [Channel Stats](https://docs.socialkit.dev/api-reference/facebook-channel-stats-api) — `GET, POST /facebook/channel-stats`

### Twitter/X

- [Profile](https://docs.socialkit.dev/api-reference/twitter-profile-api) — `GET, POST /twitter/profile`
- [Tweets](https://docs.socialkit.dev/api-reference/twitter-tweets-api) — `GET, POST /twitter/tweets`
- [Tweet](https://docs.socialkit.dev/api-reference/twitter-tweet-api) — `GET, POST /twitter/tweet`
- [Thread](https://docs.socialkit.dev/api-reference/twitter-thread-api) — `GET, POST /twitter/thread`
- [Transcript](https://docs.socialkit.dev/api-reference/twitter-transcript-api) — `GET, POST /twitter/transcript`

### LinkedIn

- [Profile](https://docs.socialkit.dev/api-reference/linkedin-profile-api) — `GET, POST /linkedin/profile`
- [Company](https://docs.socialkit.dev/api-reference/linkedin-company-api) — `GET, POST /linkedin/company`
- [Company Posts](https://docs.socialkit.dev/api-reference/linkedin-company-posts-api) — `GET, POST /linkedin/company-posts`
- [Post](https://docs.socialkit.dev/api-reference/linkedin-post-api) — `GET, POST /linkedin/post`
- [Transcript](https://docs.socialkit.dev/api-reference/linkedin-transcript-api) — `GET, POST /linkedin/transcript`

### Direct Video

- [Transcript](https://docs.socialkit.dev/api-reference/video-transcript-api) — `GET, POST /video/transcript`
- [Summary](https://docs.socialkit.dev/api-reference/video-summarize-api) — `GET, POST /video/summarize`

### Guides

- [Getting Started](https://docs.socialkit.dev/getting-started)
- [Authentication](https://docs.socialkit.dev/authentication)
- [API Reference](https://docs.socialkit.dev/api-reference)
- [Errors](https://docs.socialkit.dev/api-reference/errors)
- [OpenAPI 3.0 Contract](https://api.socialkit.dev/openapi.json)

## License

[MIT](LICENSE) — Copyright (c) 2026 Igor Sazonov

## Author

Igor Sazonov — [GitHub](https://github.com/tigusigalpa) — sovletig@gmail.com
