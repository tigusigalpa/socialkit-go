# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-09-06

### Added
- Initial release of the SocialKit Go SDK.
- `Client` with options pattern (`WithBaseURL`, `WithHTTPClient`, `WithTimeout`,
  `WithRetry`, `WithKeyInQuery`, `WithUserAgent`).
- Secure authentication via `x-access-key` header by default; optional
  `access_key` in query/body via `WithKeyInQuery`.
- `StatusService` for the public status endpoint.
- `CreditsService` for the credits balance endpoint.
- `YouTubeService` with transcript, summarize, stats, comments, channel-stats,
  search, videos, download, and bulk endpoints.
- `TikTokService` with transcript, summarize, stats, comments, channel-stats,
  channel-videos, search, hashtag search, download, and bulk endpoints.
- `InstagramService` with transcript, summarize, stats, comments, channel-stats,
  channel-posts, channel-reels, reels-search, download, and bulk endpoints.
- `FacebookService` with transcript, summarize, stats, and comments endpoints.
- `TwitterService` with profile, tweets, tweet, thread, and transcript endpoints.
- `LinkedInService` with profile, company, company-posts, post, and transcript
  endpoints.
- `VideoService` with transcript and summarize for direct video file URLs.
- `DownloadsService` for v2 async download jobs (start, get, wait).
- Typed error hierarchy with `errors.Is`/`errors.As` support.
- `ResponseMeta` with credits and rate-limit header metadata.
- Automatic retry with exponential backoff and jitter for 429/5xx responses.
- Comprehensive test suite using `httptest` with no network or real API key.
- Example programs for basic usage, async downloads, and error handling.
- GitHub Actions CI workflow.
