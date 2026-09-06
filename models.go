package socialkit

import (
	"encoding/json"
	"time"
)

// ---------------------------------------------------------------------------
// Status models
// ---------------------------------------------------------------------------

// Status represents the overall API status response.
type Status struct {
	Overall     string       `json:"overall"`
	GeneratedAt string       `json:"generatedAt"`
	Tools       []StatusTool `json:"tools"`
}

// StatusTool represents the status of a single SocialKit tool/channel.
type StatusTool struct {
	ID        string `json:"id"`
	Channel   string `json:"channel"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	UpdatedAt string `json:"updatedAt"`
}

// ---------------------------------------------------------------------------
// Credits models
// ---------------------------------------------------------------------------

// Credits represents the credit balance for the authenticated account.
type Credits struct {
	TotalRemaining float64           `json:"totalRemaining"`
	Monthly        *CreditsBreakdown `json:"monthly,omitempty"`
	Purchased      *CreditsBreakdown `json:"purchased,omitempty"`
}

// CreditsBreakdown represents a breakdown of credits by source.
type CreditsBreakdown struct {
	Remaining *float64 `json:"remaining,omitempty"`
	Used      *float64 `json:"used,omitempty"`
	Total     *float64 `json:"total,omitempty"`
}

// ---------------------------------------------------------------------------
// Transcript models
// ---------------------------------------------------------------------------

// Transcript represents a video transcript with optional segment-level
// timing data.
type Transcript struct {
	URL                string              `json:"url"`
	VideoID            string              `json:"videoId"`
	Transcript         string              `json:"transcript"`
	TranscriptSegments []TranscriptSegment `json:"transcriptSegments"`
	WordCount          *int                `json:"wordCount,omitempty"`
	Segments           *int                `json:"segments,omitempty"`
	Language           *string             `json:"language,omitempty"`
	DurationSeconds    *float64            `json:"durationSeconds,omitempty"`
	CreditsUsed        *float64            `json:"creditsUsed,omitempty"`
}

// TranscriptSegment represents a single segment of a transcript with
// timing information.
type TranscriptSegment struct {
	Text      string  `json:"text"`
	Start     float64 `json:"start"`
	Duration  float64 `json:"duration"`
	Timestamp *string `json:"timestamp,omitempty"`
}

// LinkedInTranscript represents a LinkedIn video transcript. Segments do
// not include a timestamp field.
type LinkedInTranscript struct {
	Transcript         string                      `json:"transcript"`
	TranscriptSegments []LinkedInTranscriptSegment `json:"transcriptSegments"`
	WordCount          *int                        `json:"wordCount,omitempty"`
	DurationSeconds    *float64                    `json:"durationSeconds,omitempty"`
	CreditsUsed        *float64                    `json:"creditsUsed,omitempty"`
}

// LinkedInTranscriptSegment represents a single segment of a LinkedIn
// transcript. Unlike TranscriptSegment, it does not include a timestamp.
type LinkedInTranscriptSegment struct {
	Text     string  `json:"text"`
	Start    float64 `json:"start"`
	Duration float64 `json:"duration"`
}

// TwitterTranscript represents a Twitter/X video transcript.
type TwitterTranscript struct {
	URL                string              `json:"url"`
	VideoID            string              `json:"videoId"`
	Transcript         string              `json:"transcript"`
	TranscriptSegments []TranscriptSegment `json:"transcriptSegments"`
	WordCount          *int                `json:"wordCount,omitempty"`
	Segments           *int                `json:"segments,omitempty"`
	Language           *string             `json:"language,omitempty"`
	DurationSeconds    *float64            `json:"durationSeconds,omitempty"`
	CreditsUsed        *float64            `json:"creditsUsed,omitempty"`
}

// ---------------------------------------------------------------------------
// Summary model
// ---------------------------------------------------------------------------

// Summary represents a video or post summary with key points, tone, and
// other AI-generated insights. The Extra map captures any custom_response
// fields not explicitly modeled.
type Summary struct {
	URL             string                     `json:"url"`
	Summary         string                     `json:"summary"`
	MainTopics      []string                   `json:"mainTopics"`
	KeyPoints       []string                   `json:"keyPoints"`
	Tone            string                     `json:"tone"`
	TargetAudience  string                     `json:"targetAudience"`
	Quotes          []string                   `json:"quotes"`
	DurationSeconds *float64                   `json:"durationSeconds,omitempty"`
	CreditsUsed     *float64                   `json:"creditsUsed,omitempty"`
	Extra           map[string]json.RawMessage `json:"-"`
}

// UnmarshalJSON implements custom JSON unmarshalling for Summary. It
// decodes known fields normally and captures any unknown top-level fields
// into the Extra map for forward compatibility with custom_response
// fields.
func (s *Summary) UnmarshalJSON(data []byte) error {
	type alias Summary
	var raw alias
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*s = Summary(raw)

	// Capture unknown fields into Extra.
	var generic map[string]json.RawMessage
	if err := json.Unmarshal(data, &generic); err != nil {
		return nil
	}
	known := map[string]bool{
		"url": true, "summary": true, "mainTopics": true, "keyPoints": true,
		"tone": true, "targetAudience": true, "quotes": true,
		"durationSeconds": true, "creditsUsed": true,
	}
	s.Extra = make(map[string]json.RawMessage)
	for k, v := range generic {
		if !known[k] {
			s.Extra[k] = v
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// YouTube stats model
// ---------------------------------------------------------------------------

// YouTubeStats represents statistics for a YouTube video.
type YouTubeStats struct {
	URL           string `json:"url"`
	VideoID       string `json:"videoId"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	ChannelName   string `json:"channelName"`
	ChannelLink   string `json:"channelLink"`
	ChannelHandle string `json:"channelHandle"`
	Views         *int   `json:"views,omitempty"`
	Likes         *int   `json:"likes,omitempty"`
	Comments      *int   `json:"comments,omitempty"`
	Duration      string `json:"duration"`
	PublishedAt   string `json:"publishedAt"`
	ThumbnailURL  string `json:"thumbnailUrl"`
	ContentType   string `json:"contentType"`
	IsShortForm   *bool  `json:"isShortForm,omitempty"`
}

// ---------------------------------------------------------------------------
// TikTok stats model
// ---------------------------------------------------------------------------

// TikTokStats represents statistics for a TikTok video.
type TikTokStats struct {
	URL          string `json:"url"`
	VideoID      string `json:"videoId"`
	Title        string `json:"title"`
	ChannelName  string `json:"channelName"`
	ChannelLink  string `json:"channelLink"`
	Likes        *int   `json:"likes,omitempty"`
	Comments     *int   `json:"comments,omitempty"`
	Collects     *int   `json:"collects,omitempty"`
	Shares       *int   `json:"shares,omitempty"`
	Views        *int   `json:"views,omitempty"`
	Description  string `json:"description"`
	Duration     string `json:"duration"`
	ThumbnailURL string `json:"thumbnailUrl"`
	PublishedAt  string `json:"publishedAt"`
	ContentType  string `json:"contentType"`
	IsShortForm  *bool  `json:"isShortForm,omitempty"`
}

// ---------------------------------------------------------------------------
// Instagram stats model
// ---------------------------------------------------------------------------

// InstagramStats represents statistics for an Instagram post.
type InstagramStats struct {
	PostURL     string `json:"postUrl"`
	ID          string `json:"id"`
	Shortcode   string `json:"shortcode"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Views       *int   `json:"views,omitempty"`
	Likes       *int   `json:"likes,omitempty"`
	Comments    *int   `json:"comments,omitempty"`
	PublishedAt string `json:"publishedAt"`
	Author      string `json:"author"`
	AuthorLink  string `json:"authorLink"`
	Duration    string `json:"duration"`
	Thumbnail   string `json:"thumbnail"`
	IsVideo     *bool  `json:"isVideo,omitempty"`
	ContentType string `json:"contentType"`
}

// ---------------------------------------------------------------------------
// Facebook stats model
// ---------------------------------------------------------------------------

// FacebookStats represents statistics for a Facebook post.
type FacebookStats struct {
	PostURL     string     `json:"postUrl"`
	ID          string     `json:"id"`
	Description string     `json:"description"`
	Views       *int       `json:"views,omitempty"`
	Likes       *int       `json:"likes,omitempty"`
	Comments    *int       `json:"comments,omitempty"`
	Shares      *int       `json:"shares,omitempty"`
	Reactions   []Reaction `json:"reactions,omitempty"`
	Author      string     `json:"author"`
	AuthorLink  string     `json:"authorLink"`
	IsVideo     *bool      `json:"isVideo,omitempty"`
}

// Reaction represents a single Facebook reaction type and its count.
type Reaction struct {
	Name      string `json:"name"`
	Count     *int   `json:"count,omitempty"`
	Formatted string `json:"formatted"`
}

// ---------------------------------------------------------------------------
// Comments models
// ---------------------------------------------------------------------------

// YouTubeComments represents comments for a YouTube video.
type YouTubeComments struct {
	URL      string           `json:"url"`
	Comments []YouTubeComment `json:"comments"`
}

// YouTubeComment represents a single YouTube comment.
type YouTubeComment struct {
	Author     string `json:"author"`
	Text       string `json:"text"`
	Likes      *int   `json:"likes,omitempty"`
	Date       string `json:"date"`
	Avatar     string `json:"avatar"`
	ReplyCount *int   `json:"replyCount,omitempty"`
	Position   *int   `json:"position,omitempty"`
}

// TikTokComments represents comments for a TikTok video with pagination.
type TikTokComments struct {
	VideoID      string          `json:"videoId"`
	Comments     []TikTokComment `json:"comments"`
	CommentCount *int            `json:"commentCount,omitempty"`
	HasMore      *bool           `json:"hasMore,omitempty"`
	Cursor       *string         `json:"cursor,omitempty"`
}

// TikTokComment represents a single TikTok comment.
type TikTokComment struct {
	ID         string `json:"id"`
	Author     string `json:"author"`
	Username   string `json:"username"`
	Text       string `json:"text"`
	Likes      *int   `json:"likes,omitempty"`
	Avatar     string `json:"avatar"`
	CreateTime *int   `json:"createTime,omitempty"`
	ReplyCount *int   `json:"replyCount,omitempty"`
}

// InstagramComments represents comments for an Instagram post with
// pagination.
type InstagramComments struct {
	PostURL      string             `json:"postUrl"`
	Comments     []InstagramComment `json:"comments"`
	CommentCount *int               `json:"commentCount,omitempty"`
	HasMore      *bool              `json:"hasMore,omitempty"`
	Cursor       *string            `json:"cursor,omitempty"`
}

// FacebookComments represents comments for a Facebook post with
// pagination.
type FacebookComments struct {
	PostURL      string            `json:"postUrl"`
	Comments     []FacebookComment `json:"comments"`
	CommentCount *int              `json:"commentCount,omitempty"`
	HasMore      *bool             `json:"hasMore,omitempty"`
	Cursor       *string           `json:"cursor,omitempty"`
}

// InstagramComment represents a single Instagram comment.
type InstagramComment struct {
	ID         string                 `json:"id"`
	Text       string                 `json:"text"`
	LikeCount  *int                   `json:"likeCount,omitempty"`
	ReplyCount *int                   `json:"replyCount,omitempty"`
	User       map[string]interface{} `json:"user,omitempty"`
}

// FacebookComment represents a single Facebook comment. It has the same
// structure as InstagramComment.
type FacebookComment struct {
	ID         string                 `json:"id"`
	Text       string                 `json:"text"`
	LikeCount  *int                   `json:"likeCount,omitempty"`
	ReplyCount *int                   `json:"replyCount,omitempty"`
	User       map[string]interface{} `json:"user,omitempty"`
}

// ---------------------------------------------------------------------------
// Channel stats models
// ---------------------------------------------------------------------------

// YouTubeChannelStats represents statistics for a YouTube channel.
type YouTubeChannelStats struct {
	ProfileURL  string `json:"profileUrl"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	Bio         string `json:"bio"`
	Verified    *bool  `json:"verified,omitempty"`
	Avatar      string `json:"avatar"`
	Subscribers *int   `json:"subscribers,omitempty"`
	TotalVideos *int   `json:"totalVideos,omitempty"`
	BioLink     string `json:"bioLink"`
	Banner      string `json:"banner"`
}

// TikTokChannelStats represents statistics for a TikTok channel.
type TikTokChannelStats struct {
	ProfileURL  string `json:"profileUrl"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	Signature   string `json:"signature"`
	Verified    *bool  `json:"verified,omitempty"`
	Avatar      string `json:"avatar"`
	Followers   *int   `json:"followers,omitempty"`
	Following   *int   `json:"following,omitempty"`
	TotalLikes  *int   `json:"totalLikes,omitempty"`
	TotalVideos *int   `json:"totalVideos,omitempty"`
	BioLink     string `json:"bioLink"`
}

// InstagramChannelStats represents statistics for an Instagram channel.
type InstagramChannelStats struct {
	ProfileURL string `json:"profileUrl"`
	UserID     string `json:"userId"`
	Username   string `json:"username"`
	FullName   string `json:"fullName"`
	Verified   *bool  `json:"verified,omitempty"`
	Followers  *int   `json:"followers,omitempty"`
	Following  *int   `json:"following,omitempty"`
	TotalPosts *int   `json:"totalPosts,omitempty"`
	Bio        string `json:"bio"`
	Avatar     string `json:"avatar"`
}

// FacebookChannelStats represents statistics for a Facebook page/profile.
type FacebookChannelStats struct {
	ProfileURL string `json:"profileUrl"`
	ID         string `json:"id"`
	FullName   string `json:"fullName"`
	Bio        string `json:"bio"`
	Avatar     string `json:"avatar"`
	CoverPhoto string `json:"coverPhoto"`
	Verified   *bool  `json:"verified,omitempty"`
	Followers  *int   `json:"followers,omitempty"`
	Category   string `json:"category"`
}

// ---------------------------------------------------------------------------
// Search/list models
// ---------------------------------------------------------------------------

// YouTubeSearchResult represents search results for YouTube videos.
type YouTubeSearchResult struct {
	Query   string              `json:"query"`
	Results []YouTubeSearchItem `json:"results"`
}

// YouTubeSearchItem represents a single YouTube search result.
type YouTubeSearchItem struct {
	VideoID        string `json:"videoId"`
	Title          string `json:"title"`
	Thumbnail      string `json:"thumbnail"`
	ChannelName    string `json:"channelName"`
	ChannelID      string `json:"channelId"`
	ChannelURL     string `json:"channelUrl"`
	PublishedTime  string `json:"publishedTime"`
	Duration       string `json:"duration"`
	Views          *int   `json:"views,omitempty"`
	ViewsFormatted string `json:"viewsFormatted"`
	Description    string `json:"description"`
	URL            string `json:"url"`
	Verified       *bool  `json:"verified,omitempty"`
}

// TikTokSearchResult represents search results for TikTok videos with
// pagination.
type TikTokSearchResult struct {
	Query   string             `json:"query"`
	Results []TikTokSearchItem `json:"results"`
	HasMore *bool              `json:"hasMore,omitempty"`
	Cursor  *string            `json:"cursor,omitempty"`
}

// TikTokSearchItem represents a single TikTok search result.
type TikTokSearchItem struct {
	ID         string                 `json:"id"`
	Desc       string                 `json:"desc"`
	URL        string                 `json:"url"`
	CreateTime string                 `json:"createTime"`
	Author     map[string]interface{} `json:"author,omitempty"`
	Stats      map[string]interface{} `json:"stats,omitempty"`
}

// TikTokHashtagResult represents videos for a TikTok hashtag with
// pagination.
type TikTokHashtagResult struct {
	Hashtag string              `json:"hashtag"`
	Results []TikTokHashtagItem `json:"results"`
	HasMore *bool               `json:"hasMore,omitempty"`
	Cursor  *string             `json:"cursor,omitempty"`
}

// TikTokHashtagItem represents a single TikTok hashtag video.
type TikTokHashtagItem struct {
	ID     string                 `json:"id"`
	Desc   string                 `json:"desc"`
	URL    string                 `json:"url"`
	Author map[string]interface{} `json:"author,omitempty"`
	Stats  map[string]interface{} `json:"stats,omitempty"`
}

// TikTokChannelVideos represents videos from a TikTok channel with
// pagination.
type TikTokChannelVideos struct {
	ProfileURL  string            `json:"profileUrl"`
	ChannelName string            `json:"channelName"`
	Results     []TikTokVideoItem `json:"results"`
	HasMore     *bool             `json:"hasMore,omitempty"`
	Cursor      *string           `json:"cursor,omitempty"`
}

// TikTokVideoItem represents a single TikTok video from a channel.
type TikTokVideoItem struct {
	VideoID     string `json:"videoId"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Thumbnail   string `json:"thumbnail"`
	Duration    string `json:"duration"`
	CreateTime  *int   `json:"createTime,omitempty"`
	Views       *int   `json:"views,omitempty"`
	Likes       *int   `json:"likes,omitempty"`
	Comments    *int   `json:"comments,omitempty"`
	Shares      *int   `json:"shares,omitempty"`
	Collects    *int   `json:"collects,omitempty"`
}

// ---------------------------------------------------------------------------
// Instagram feed models
// ---------------------------------------------------------------------------

// InstagramChannelPosts represents posts from an Instagram channel with
// pagination.
type InstagramChannelPosts struct {
	ProfileURL string              `json:"profileUrl"`
	Username   string              `json:"username"`
	Items      []InstagramPostItem `json:"items"`
	Count      *int                `json:"count,omitempty"`
	HasMore    *bool               `json:"hasMore,omitempty"`
	Cursor     *string             `json:"cursor,omitempty"`
}

// InstagramPostItem represents a single Instagram post from a channel.
type InstagramPostItem struct {
	ID           string `json:"id"`
	Shortcode    string `json:"shortcode"`
	URL          string `json:"url"`
	Type         string `json:"type"`
	Caption      string `json:"caption"`
	Likes        *int   `json:"likes,omitempty"`
	Comments     *int   `json:"comments,omitempty"`
	Views        *int   `json:"views,omitempty"`
	Timestamp    string `json:"timestamp"`
	ThumbnailURL string `json:"thumbnailUrl"`
	Author       string `json:"author"`
}

// InstagramChannelReels represents reels from an Instagram channel with
// pagination.
type InstagramChannelReels struct {
	ProfileURL string              `json:"profileUrl"`
	Username   string              `json:"username"`
	Items      []InstagramReelItem `json:"items"`
	Count      *int                `json:"count,omitempty"`
	HasMore    *bool               `json:"hasMore,omitempty"`
	Cursor     *string             `json:"cursor,omitempty"`
}

// InstagramReelItem represents a single Instagram reel from a channel.
type InstagramReelItem struct {
	ID           string `json:"id"`
	Shortcode    string `json:"shortcode"`
	URL          string `json:"url"`
	Caption      string `json:"caption"`
	Likes        *int   `json:"likes,omitempty"`
	Comments     *int   `json:"comments,omitempty"`
	Views        *int   `json:"views,omitempty"`
	Duration     string `json:"duration"`
	ThumbnailURL string `json:"thumbnailUrl"`
	Author       string `json:"author"`
}

// InstagramReelsSearch represents search results for Instagram reels.
type InstagramReelsSearch struct {
	Query   string                    `json:"query"`
	Items   []InstagramReelSearchItem `json:"items"`
	Count   *int                      `json:"count,omitempty"`
	HasMore *bool                     `json:"hasMore,omitempty"`
}

// InstagramReelSearchItem represents a single Instagram reel search result.
type InstagramReelSearchItem struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Caption  string `json:"caption"`
	Likes    *int   `json:"likes,omitempty"`
	Comments *int   `json:"comments,omitempty"`
	Views    *int   `json:"views,omitempty"`
	Author   string `json:"author"`
}

// ---------------------------------------------------------------------------
// YouTube videos model
// ---------------------------------------------------------------------------

// YouTubeVideos represents a list of YouTube videos from a channel or
// playlist.
type YouTubeVideos struct {
	URL     string             `json:"url"`
	Type    string             `json:"type"`
	Results []YouTubeVideoItem `json:"results"`
}

// YouTubeVideoItem represents a single YouTube video in a list.
type YouTubeVideoItem struct {
	VideoID        string `json:"videoId"`
	Title          string `json:"title"`
	Thumbnail      string `json:"thumbnail"`
	ChannelName    string `json:"channelName"`
	Duration       string `json:"duration"`
	Views          *int   `json:"views,omitempty"`
	ViewsFormatted string `json:"viewsFormatted"`
	PublishedTime  string `json:"publishedTime"`
	PublishedAt    string `json:"publishedAt"`
	Description    string `json:"description"`
	URL            string `json:"url"`
}

// ---------------------------------------------------------------------------
// Download model
// ---------------------------------------------------------------------------

// DownloadResult represents a download result for a video.
type DownloadResult struct {
	URL             string   `json:"url"`
	Title           string   `json:"title"`
	Duration        string   `json:"duration"`
	DurationSeconds *float64 `json:"durationSeconds,omitempty"`
	DownloadURL     string   `json:"downloadUrl"`
	FileSize        *float64 `json:"fileSize,omitempty"`
	FileSizeMB      string   `json:"fileSizeMB"`
	Format          string   `json:"format"`
	Quality         string   `json:"quality"`
	ExpiresIn       string   `json:"expiresIn"`
	CreditsUsed     *float64 `json:"creditsUsed,omitempty"`
}

// ---------------------------------------------------------------------------
// Twitter models
// ---------------------------------------------------------------------------

// TwitterProfile represents a Twitter/X user profile.
type TwitterProfile struct {
	URL          string `json:"url"`
	ID           string `json:"id"`
	Name         string `json:"name"`
	Username     string `json:"username"`
	Bio          string `json:"bio"`
	Followers    *int   `json:"followers,omitempty"`
	Following    *int   `json:"following,omitempty"`
	Tweets       *int   `json:"tweets,omitempty"`
	Verified     *bool  `json:"verified,omitempty"`
	ProfileImage string `json:"profileImage"`
	BannerImage  string `json:"bannerImage"`
	Location     string `json:"location"`
	Website      string `json:"website"`
	JoinedAt     string `json:"joinedAt"`
}

// TwitterTweets represents a list of tweets with an optional next cursor.
type TwitterTweets struct {
	Tweets     []TwitterTweet `json:"tweets"`
	NextCursor *string        `json:"nextCursor,omitempty"`
}

// TwitterTweet represents a single tweet.
type TwitterTweet struct {
	ID        string                 `json:"id"`
	Text      string                 `json:"text"`
	Likes     *int                   `json:"likes,omitempty"`
	Retweets  *int                   `json:"retweets,omitempty"`
	Replies   *int                   `json:"replies,omitempty"`
	Views     *int                   `json:"views,omitempty"`
	CreatedAt string                 `json:"createdAt"`
	Author    map[string]interface{} `json:"author,omitempty"`
	Hashtags  []string               `json:"hashtags,omitempty"`
	URLs      []string               `json:"urls,omitempty"`
}

// TwitterThread represents a Twitter/X thread.
type TwitterThread struct {
	ConversationID string          `json:"conversationId"`
	TweetCount     *int            `json:"tweetCount,omitempty"`
	IsThread       *bool           `json:"isThread,omitempty"`
	CombinedText   string          `json:"combinedText"`
	Author         json.RawMessage `json:"author,omitempty"`
	Tweets         json.RawMessage `json:"tweets,omitempty"`
}

// ---------------------------------------------------------------------------
// LinkedIn models
// ---------------------------------------------------------------------------

// LinkedInProfile represents a LinkedIn user profile.
type LinkedInProfile struct {
	URL            string            `json:"url"`
	Name           string            `json:"name"`
	Headline       string            `json:"headline"`
	Followers      *int              `json:"followers,omitempty"`
	Connections    string            `json:"connections"`
	ProfileImage   string            `json:"profileImage"`
	RecentArticles []LinkedInArticle `json:"recentArticles,omitempty"`
}

// LinkedInArticle represents a recent article on a LinkedIn profile.
type LinkedInArticle struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Date  string `json:"date"`
}

// LinkedInCompany represents a LinkedIn company page.
type LinkedInCompany struct {
	URL         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Followers   *int   `json:"followers,omitempty"`
	Employees   string `json:"employees"`
	Logo        string `json:"logo"`
}

// LinkedInCompanyPosts represents posts from a LinkedIn company page.
type LinkedInCompanyPosts struct {
	Posts []LinkedInPostItem `json:"posts"`
}

// LinkedInPostItem represents a single LinkedIn company post.
type LinkedInPostItem struct {
	ID          string                 `json:"id"`
	Text        string                 `json:"text"`
	Likes       *int                   `json:"likes,omitempty"`
	Comments    *int                   `json:"comments,omitempty"`
	Shares      *int                   `json:"shares,omitempty"`
	ImageURL    *string                `json:"imageUrl,omitempty"`
	VideoURL    *string                `json:"videoUrl,omitempty"`
	PublishedAt string                 `json:"publishedAt"`
	Author      map[string]interface{} `json:"author,omitempty"`
}

// LinkedInPost represents a single LinkedIn post. The post content is
// returned as raw JSON because the upstream API does not guarantee a
// fixed schema.
type LinkedInPost struct {
	Post json.RawMessage `json:"post,omitempty"`
}

// ---------------------------------------------------------------------------
// V2 download models
// ---------------------------------------------------------------------------

// V2DownloadRequest represents the request body for starting a v2 async
// download job.
type V2DownloadRequest struct {
	URL         string `json:"url"`
	Format      string `json:"format,omitempty"`
	Quality     string `json:"quality,omitempty"`
	Country     string `json:"country,omitempty"`
	MaxDuration *int   `json:"max_duration,omitempty"`
}

// V2Job represents a v2 async download job.
type V2Job struct {
	JobID           string   `json:"jobId"`
	Status          string   `json:"status"`
	Platform        string   `json:"platform"`
	URL             string   `json:"url"`
	Quality         string   `json:"quality"`
	Format          string   `json:"format"`
	Billed          *bool    `json:"billed,omitempty"`
	DownloadURL     string   `json:"downloadUrl"`
	ExpiresIn       *float64 `json:"expiresIn,omitempty"`
	DurationSeconds *float64 `json:"durationSeconds,omitempty"`
	FileSize        *float64 `json:"fileSize,omitempty"`
	FileSizeMB      string   `json:"fileSizeMB"`
	CreditsCost     *float64 `json:"creditsCost,omitempty"`
	Title           string   `json:"title"`
	Thumbnail       string   `json:"thumbnail"`
	Error           string   `json:"error"`
	ErrorCode       string   `json:"errorCode"`
	Retryable       *bool    `json:"retryable,omitempty"`
}

// ---------------------------------------------------------------------------
// Request param structs
// ---------------------------------------------------------------------------

// FetchRequest holds parameters for simple fetch operations (transcript,
// stats, etc.).
type FetchRequest struct {
	URL      string `json:"url"`
	Cache    *bool  `json:"cache,omitempty"`
	CacheTTL *int   `json:"cacheTtl,omitempty"`
}

// SummaryRequest holds parameters for summary operations.
type SummaryRequest struct {
	URL            string      `json:"url"`
	CustomPrompt   *string     `json:"customPrompt,omitempty"`
	CustomResponse interface{} `json:"customResponse,omitempty"`
	Cache          *bool       `json:"cache,omitempty"`
	CacheTTL       *int        `json:"cacheTtl,omitempty"`
}

// ListRequest holds parameters for list/feed operations with cursor
// pagination.
type ListRequest struct {
	URL      string  `json:"url"`
	Limit    *int    `json:"limit,omitempty"`
	Cursor   *string `json:"cursor,omitempty"`
	Cache    *bool   `json:"cache,omitempty"`
	CacheTTL *int    `json:"cacheTtl,omitempty"`
}

// CommentsRequest holds parameters for comments operations.
type CommentsRequest struct {
	URL      string  `json:"url"`
	Limit    *int    `json:"limit,omitempty"`
	Cursor   *string `json:"cursor,omitempty"`
	SortBy   *string `json:"sortBy,omitempty"`
	Cache    *bool   `json:"cache,omitempty"`
	CacheTTL *int    `json:"cacheTtl,omitempty"`
}

// SearchRequest holds parameters for search operations.
type SearchRequest struct {
	Query      string  `json:"query"`
	Limit      *int    `json:"limit,omitempty"`
	SortBy     *string `json:"sortBy,omitempty"`
	UploadDate *string `json:"uploadDate,omitempty"`
	DatePosted *string `json:"datePosted,omitempty"`
	Type       *string `json:"type,omitempty"`
	Cache      *bool   `json:"cache,omitempty"`
	CacheTTL   *int    `json:"cacheTtl,omitempty"`
}

// DownloadRequest holds parameters for download operations.
type DownloadRequest struct {
	URL     string  `json:"url"`
	Format  *string `json:"format,omitempty"`
	Quality *string `json:"quality,omitempty"`
	Country *string `json:"country,omitempty"`
}

// HashtagSearchRequest holds parameters for TikTok hashtag search.
type HashtagSearchRequest struct {
	Hashtag  string  `json:"hashtag"`
	Limit    *int    `json:"limit,omitempty"`
	Cursor   *string `json:"cursor,omitempty"`
	Cache    *bool   `json:"cache,omitempty"`
	CacheTTL *int    `json:"cacheTtl,omitempty"`
}

// VideosRequest holds parameters for YouTube videos list operations.
type VideosRequest struct {
	URL         string `json:"url"`
	Limit       *int   `json:"limit,omitempty"`
	FullDetails *bool  `json:"fullDetails,omitempty"`
	Cache       *bool  `json:"cache,omitempty"`
	CacheTTL    *int   `json:"cacheTtl,omitempty"`
}

// ReelsSearchRequest holds parameters for Instagram reels search.
type ReelsSearchRequest struct {
	Query    string `json:"query"`
	Limit    *int   `json:"limit,omitempty"`
	Page     *int   `json:"page,omitempty"`
	Cache    *bool  `json:"cache,omitempty"`
	CacheTTL *int   `json:"cacheTtl,omitempty"`
}

// TweetsRequest holds parameters for Twitter tweets list operations.
type TweetsRequest struct {
	URL      string  `json:"url"`
	Limit    *int    `json:"limit,omitempty"`
	Cursor   *string `json:"cursor,omitempty"`
	Cache    *bool   `json:"cache,omitempty"`
	CacheTTL *int    `json:"cacheTtl,omitempty"`
}

// CompanyPostsRequest holds parameters for LinkedIn company posts.
type CompanyPostsRequest struct {
	URL      string `json:"url"`
	Limit    *int   `json:"limit,omitempty"`
	Cache    *bool  `json:"cache,omitempty"`
	CacheTTL *int   `json:"cacheTtl,omitempty"`
}

// WaitOptions controls the polling behaviour of DownloadsService.Wait.
type WaitOptions struct {
	// Interval is the delay between poll attempts. Defaults to 2 seconds
	// if zero.
	Interval time.Duration
	// MaxAttempts is the maximum number of poll attempts. If zero, the
	// poll continues until the context deadline is reached.
	MaxAttempts int
	// BackoffFactor multiplies the interval after each attempt. A value
	// of 0 or 1 disables backoff.
	BackoffFactor float64
}

// ensure time package is used (for WaitOptions.Interval).
var _ = time.Second
