# Upstream Compatibility Notes

## Verification Date

**Checked:** 2026-09-06

## API Version

- **OpenAPI version:** 3.0.3
- **API version:** 1.1.0
- **OpenAPI URL:** https://api.socialkit.dev/openapi.json
- **Base URL:** https://api.socialkit.dev

## Authentication

- Primary method: `x-access-key` HTTP header (apiKey security scheme)
- Compatibility: `access_key` query parameter (GET) and `access_key` JSON field (POST)
- SDK default: header only; query/JSON key is opt-in for compatibility

## Endpoint Summary

Total: 104 operations across 52 unique paths.

### Service endpoints (GET-only)
- `GET /status` — no key required, cacheable 60s
- `GET /credits` — requires key, costs 0 credits

### Async downloads v2
- `GET, POST /v2/{platform}/download` — platform enum: youtube, tiktok, instagram, facebook
- `GET /v2/downloads/{jobId}` — poll job status

### YouTube (8 endpoints + 4 bulk)
- transcript, summarize, stats, comments, channel-stats, search, videos, download
- bulk: transcript/bulk, comments/bulk, stats/bulk, summarize/bulk

### TikTok (9 endpoints + 5 bulk)
- transcript, summarize, stats, comments, channel-stats, channel-videos, search, hashtag-search, download
- bulk: transcript/bulk, comments/bulk, stats/bulk, channel-stats/bulk, summarize/bulk

### Instagram (9 endpoints + 4 bulk)
- transcript, summarize, stats, comments, channel-stats, channel-posts, channel-reels, reels-search, download
- bulk: transcript/bulk, stats/bulk, channel-stats/bulk, summarize/bulk

### Facebook (4 endpoints, no bulk, no sync download)
- transcript, summarize, stats, comments, channel-stats

### Twitter/X (5 endpoints)
- profile, tweets, tweet, thread, transcript

### LinkedIn (5 endpoints)
- profile, company, company-posts, post, transcript

### Direct video (2 endpoints)
- transcript, summarize

## Discrepancies Between OpenAPI and Narrative Docs

### Parameters in narrative docs but NOT in OpenAPI

The following parameters are documented on endpoint pages but are absent from the
OpenAPI 3.0.3 contract for some operations:

1. **`cache`** (boolean, default false) — server-side response caching
2. **`cache_ttl`** (integer, default 2592000, range 3600–2592000) — cache TTL in seconds
3. **`cursor`** — pagination cursor for TikTok/Instagram comments, channel-videos, search, hashtag-search, channel-posts, channel-reels
4. **`sortBy`** — sort order for YouTube comments (top|new), YouTube search (relevance|date|views|rating), TikTok search (relevance|likes|date)
5. **`uploadDate`** — YouTube search filter (hour|today|week|month|year)
6. **`type`** — YouTube search filter (video|shorts)
7. **`datePosted`** — TikTok search filter (day|week|month|3months|6months)
8. **`custom_prompt`** — custom prompt for summarize endpoints
9. **`custom_response`** — custom response format for summarize endpoints

### Resolution

These parameters are supported via typed options in the SDK but documented as
narrative-only. They are NOT generated exclusively from OpenAPI. The SDK sends
them when explicitly set by the caller.

### Bulk endpoints

All 13 bulk POST paths have `requestBody` schemas in OpenAPI with `anyOf` containing:
- An object with `requests` array (minItems=1, maxItems=5) of `{url: string, additionalProperties: true}`
- A nullable empty object

The response 200 schema only declares `{success: boolean}` with no `data` field.
There is **no dedicated endpoint-doc page** for bulk operations.

**SDK approach:** Each bulk method is labelled experimental, accepts a raw
JSON-safe payload, and returns `ApiResponse<json.RawMessage|array>`. No invented
strongly typed request/result classes.

### YouTube comments

OpenAPI shows no `cursor` or `hasMore` fields — YouTube comments return a flat
array without pagination cursors. This differs from TikTok/Instagram/Facebook
comments which do have cursor-based pagination.

### LinkedIn transcript

LinkedIn transcript segments have `text`, `start`, `duration` but NO `timestamp`
field (unlike YouTube/TikTok/Instagram/Facebook/Video/Twitter which have `timestamp`).

### Twitter transcript

Twitter transcript returns the standard schema with `transcript` (string),
`transcriptSegments` (array with text/start/duration/timestamp), `videoId`, etc.
This is the same aggregate form as YouTube/TikTok — NOT the array-segments-plus-text
form described in the narrative docs. The SDK models it as the standard transcript
schema.

### LinkedIn post

The `data.post` field is an untyped object in OpenAPI (no properties defined).
The SDK uses `json.RawMessage`/`array` for this field.

### Twitter thread

The `data.author` and `data.tweets[]` fields are untyped objects in OpenAPI.
The SDK uses `json.RawMessage`/`array` for these fields.

### Twitter tweet

The `data.tweet` field is an untyped object in OpenAPI.
The SDK uses `json.RawMessage`/`array` for this field.

## Error Response Schema

All error responses follow `{success: false, message: string}` with optional fields:

### 403 insufficient_credits
- `code`: "insufficient_credits"
- `required_credits`: number
- `remaining_credits`: number
- `shortfall_credits`: number
- `reset_at`: string (date-time) | null
- `upgrade_url`: string (uri)
- `top_up_url`: string (uri)
- `incident_id`: string
- `upgradeUrl`: string (uri) — deprecated alias

### Status codes
- 400: Invalid or missing request input
- 401: Access key is missing
- 403: Access key invalid or insufficient credits
- 404: Content not found or unavailable
- 429: Project rate limit exceeded
- 500: Internal server error
- 503: Safe backing-store failure (status/credits)

## Response Headers

The SDK captures these headers in `ResponseMeta`:
- `X-Credits-Used`
- `X-Credits-Remaining`
- `Retry-After`
- `X-RateLimit-Limit`
- `X-RateLimit-Remaining`
- `X-RateLimit-Reset`

## Credits Information

- Summary endpoints: 2 credits per successful request
- Comments/search/videos: 1 credit per 50 results (YouTube), 1 credit per 20 results (Instagram)
- Video transcript/summary: 5 credits per started minute
- Status: 0 credits, cacheable 60s
- Credits: 0 credits
- V2 download polling: 0 credits
- V2 download max_duration failure: no credit charged
