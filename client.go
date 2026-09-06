// Package socialkit provides a Go SDK for the SocialKit REST API
// (https://api.socialkit.dev). It exposes service structs grouped by
// resource (Status, Credits, YouTube, TikTok, Instagram, Facebook,
// Twitter, LinkedIn, Video, Downloads) accessible as fields on Client.
//
// Authentication uses the x-access-key HTTP header by default. The key is
// never placed in the URL or query string unless WithKeyInQuery(true) is
// explicitly set for backward compatibility.
package socialkit

import (
	"net/http"
	"time"
)

// DefaultBaseURL is the default SocialKit API base URL used by NewClient
// unless overridden with WithBaseURL.
const DefaultBaseURL = "https://api.socialkit.dev"

// DefaultTimeout is the default HTTP client timeout used unless
// overridden with WithTimeout.
const DefaultTimeout = 30 * time.Second

// DefaultMaxAttempts is the default number of attempts (including the
// initial request) made for a request before giving up when receiving
// retryable HTTP responses, unless overridden with WithRetry.
const DefaultMaxAttempts = 1

// DefaultBaseDelay is the default initial backoff delay used for retry
// attempts, unless overridden with WithRetry.
const DefaultBaseDelay = 500 * time.Millisecond

// defaultUserAgent is the default User-Agent header sent with every
// request unless overridden with WithUserAgent.
const defaultUserAgent = "SocialKit-Go-SDK/1.0.0"

// Client is the entry point of the SocialKit Go SDK. It holds the HTTP
// configuration and exposes one service struct per API resource group.
// A Client is safe for concurrent use by multiple goroutines.
type Client struct {
	accessKey   string
	baseURL     string
	httpClient  *http.Client
	maxAttempts int
	baseDelay   time.Duration
	keyInQuery  bool
	userAgent   string

	// Status provides access to the /status endpoint.
	Status *StatusService
	// Credits provides access to the /credits endpoint.
	Credits *CreditsService
	// YouTube provides access to the /youtube/* endpoints.
	YouTube *YouTubeService
	// TikTok provides access to the /tiktok/* endpoints.
	TikTok *TikTokService
	// Instagram provides access to the /instagram/* endpoints.
	Instagram *InstagramService
	// Facebook provides access to the /facebook/* endpoints.
	Facebook *FacebookService
	// Twitter provides access to the /twitter/* endpoints.
	Twitter *TwitterService
	// LinkedIn provides access to the /linkedin/* endpoints.
	LinkedIn *LinkedInService
	// Video provides access to the /video/* endpoints.
	Video *VideoService
	// Downloads provides access to the /v2/* download endpoints.
	Downloads *DownloadsService
}

// NewClient creates a new SocialKit API client authenticated with the
// given accessKey. The accessKey is sent in the x-access-key header on
// every request by default. Behaviour can be customized via Option values
// such as WithBaseURL, WithHTTPClient, WithTimeout, WithRetry,
// WithKeyInQuery, and WithUserAgent.
//
// Example:
//
//	client := socialkit.NewClient("YOUR_ACCESS_KEY",
//	    socialkit.WithTimeout(15*time.Second),
//	    socialkit.WithRetry(3, time.Second),
//	)
func NewClient(accessKey string, opts ...Option) *Client {
	c := &Client{
		accessKey:   accessKey,
		baseURL:     DefaultBaseURL,
		httpClient:  &http.Client{Timeout: DefaultTimeout},
		maxAttempts: DefaultMaxAttempts,
		baseDelay:   DefaultBaseDelay,
		keyInQuery:  false,
		userAgent:   defaultUserAgent,
	}

	for _, opt := range opts {
		opt(c)
	}

	c.Status = &StatusService{client: c}
	c.Credits = &CreditsService{client: c}
	c.YouTube = &YouTubeService{client: c}
	c.TikTok = &TikTokService{client: c}
	c.Instagram = &InstagramService{client: c}
	c.Facebook = &FacebookService{client: c}
	c.Twitter = &TwitterService{client: c}
	c.LinkedIn = &LinkedInService{client: c}
	c.Video = &VideoService{client: c}
	c.Downloads = &DownloadsService{client: c}

	return c
}
