// Package tmdb provides a client for The Movie Database (TMDb) API.
package tmdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const baseURL = "https://api.themoviedb.org/3"
const imageBaseURL = "https://image.tmdb.org/t/p/w500"

// Sentinel errors for callers to classify TMDb failures.
var (
	ErrAuthInvalid   = errors.New("TMDb API key is invalid")
	ErrAuthMissing   = errors.New("no TMDb API key configured")
	ErrRateLimited   = errors.New("TMDb rate limit exceeded")
	ErrServerError   = errors.New("TMDb is temporarily unavailable")
	ErrNetworkError  = errors.New("network error reaching TMDb")
	ErrTimeout       = errors.New("TMDb request timed out")
)

// Client interacts with the TMDb API.
type Client struct {
	HTTPClient  *http.Client
	APIKey      string
	AccessToken string
}

// NewClient creates a new TMDb client from an API key or env vars.
func NewClient(apiKey string) *Client {
	return NewClientWithToken(apiKey, "")
}

// NewClientWithToken creates a TMDb client using either an API key or bearer token.
func NewClientWithToken(apiKey, accessToken string) *Client {
	if apiKey == "" {
		apiKey = os.Getenv("TMDB_API_KEY")
	}
	if accessToken == "" {
		accessToken = os.Getenv("TMDB_TOKEN")
	}
	return &Client{
		APIKey:      apiKey,
		AccessToken: accessToken,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// HasAuth returns true if the client has either an API key or access token.
func (c *Client) HasAuth() bool {
	return c.APIKey != "" || c.AccessToken != ""
}

// SearchResult holds a search result from TMDb.
type SearchResult struct {
	Overview    string  `json:"overview"`
	Title       string  `json:"title"`
	Name        string  `json:"name"`
	ReleaseDate string  `json:"release_date"`
	FirstAir    string  `json:"first_air_date"`
	PosterPath  string  `json:"poster_path"`
	MediaType   string  `json:"media_type"`
	GenreIDs    []int   `json:"genre_ids"`
	VoteAvg     float64 `json:"vote_average"`
	Popularity  float64 `json:"popularity"`
	ID          int     `json:"id"`
}

type searchResponse struct {
	Results []SearchResult `json:"results"`
}

// MovieDetails holds detailed movie info.
type MovieDetails struct {
	Title            string  `json:"title"`
	Overview         string  `json:"overview"`
	ReleaseDate      string  `json:"release_date"`
	PosterPath       string  `json:"poster_path"`
	ImdbID           string  `json:"imdb_id"`
	OriginalLanguage string  `json:"original_language"`
	Tagline          string  `json:"tagline"`
	Genres           []Genre `json:"genres"`
	VoteAvg          float64 `json:"vote_average"`
	Popularity       float64 `json:"popularity"`
	ID               int     `json:"id"`
	Runtime          int     `json:"runtime"`
	Budget           int64   `json:"budget"`
	Revenue          int64   `json:"revenue"`
}

// VideoResult holds a single video from TMDb.
type VideoResult struct {
	Key  string `json:"key"`
	Site string `json:"site"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type videosResponse struct {
	Results []VideoResult `json:"results"`
}

// TVDetails holds detailed TV show info.
type TVDetails struct {
	Name             string   `json:"name"`
	Overview         string   `json:"overview"`
	FirstAirDate     string   `json:"first_air_date"`
	PosterPath       string   `json:"poster_path"`
	OriginalLanguage string   `json:"original_language"`
	Tagline          string   `json:"tagline"`
	Genres           []Genre  `json:"genres"`
	Languages        []string `json:"languages"`
	VoteAvg          float64  `json:"vote_average"`
	Popularity       float64  `json:"popularity"`
	ID               int      `json:"id"`
	Seasons          int      `json:"number_of_seasons"`
	EpisodeRunTime   []int    `json:"episode_run_time"`
}

// Genre is a TMDb genre.
type Genre struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

// Credits holds cast and crew.
type Credits struct {
	Cast []CastMember `json:"cast"`
	Crew []CrewMember `json:"crew"`
}

// CastMember is a cast member.
type CastMember struct {
	Name      string `json:"name"`
	Character string `json:"character"`
	Order     int    `json:"order"`
}

// CrewMember is a crew member.
type CrewMember struct {
	Name string `json:"name"`
	Job  string `json:"job"`
}

// SearchMulti searches for movies and TV shows.
func (c *Client) SearchMulti(query string) ([]SearchResult, error) {
	params := url.Values{}
	params.Set("query", query)
	params.Set("page", "1")

	var resp searchResponse
	if err := c.get(c.buildURL("/search/multi", params), &resp); err != nil {
		return nil, err
	}

	var filtered []SearchResult
	for i := range resp.Results {
		if resp.Results[i].MediaType == "movie" || resp.Results[i].MediaType == "tv" {
			filtered = append(filtered, resp.Results[i])
		}
	}
	return filtered, nil
}

// GetMovieDetails returns detailed info for a movie.
func (c *Client) GetMovieDetails(tmdbID int) (*MovieDetails, error) {
	var d MovieDetails
	if err := c.get(c.buildURL(fmt.Sprintf("/movie/%d", tmdbID), nil), &d); err != nil {
		return nil, err
	}
	return &d, nil
}

// GetTVDetails returns detailed info for a TV show.
func (c *Client) GetTVDetails(tmdbID int) (*TVDetails, error) {
	var d TVDetails
	if err := c.get(c.buildURL(fmt.Sprintf("/tv/%d", tmdbID), nil), &d); err != nil {
		return nil, err
	}
	return &d, nil
}

// GetMovieCredits returns cast and crew for a movie.
func (c *Client) GetMovieCredits(tmdbID int) (*Credits, error) {
	var cr Credits
	if err := c.get(c.buildURL(fmt.Sprintf("/movie/%d/credits", tmdbID), nil), &cr); err != nil {
		return nil, err
	}
	return &cr, nil
}

// GetTVCredits returns cast and crew for a TV show.
func (c *Client) GetTVCredits(tmdbID int) (*Credits, error) {
	var cr Credits
	if err := c.get(c.buildURL(fmt.Sprintf("/tv/%d/credits", tmdbID), nil), &cr); err != nil {
		return nil, err
	}
	return &cr, nil
}

// GetMovieVideos returns videos (trailers, teasers) for a movie.
func (c *Client) GetMovieVideos(tmdbID int) ([]VideoResult, error) {
	var resp videosResponse
	if err := c.get(c.buildURL(fmt.Sprintf("/movie/%d/videos", tmdbID), nil), &resp); err != nil {
		return nil, err
	}
	return resp.Results, nil
}

// GetTVVideos returns videos (trailers, teasers) for a TV show.
func (c *Client) GetTVVideos(tmdbID int) ([]VideoResult, error) {
	var resp videosResponse
	if err := c.get(c.buildURL(fmt.Sprintf("/tv/%d/videos", tmdbID), nil), &resp); err != nil {
		return nil, err
	}
	return resp.Results, nil
}

// TrailerURL finds the best YouTube trailer URL from a list of videos.
func TrailerURL(videos []VideoResult) string {
	for _, v := range videos {
		if v.Site == "YouTube" && v.Type == "Trailer" {
			return "https://www.youtube.com/watch?v=" + v.Key
		}
	}
	for _, v := range videos {
		if v.Site == "YouTube" {
			return "https://www.youtube.com/watch?v=" + v.Key
		}
	}
	return ""
}

// DownloadPoster downloads a poster image and saves it to dst.
func (c *Client) DownloadPoster(posterPath, dst string) error {
	if posterPath == "" {
		return fmt.Errorf("no poster available")
	}

	imgURL := imageBaseURL + posterPath
	resp, err := c.HTTPClient.Get(imgURL)
	if err != nil {
		if IsNetworkError(err) {
			return fmt.Errorf("%w: %v", ErrNetworkError, err)
		}
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

// GetRecommendations returns recommended movies or TV shows.
func (c *Client) GetRecommendations(tmdbID int, mediaType string, page int) ([]SearchResult, error) {
	params := url.Values{}
	params.Set("page", fmt.Sprintf("%d", page))

	var resp searchResponse
	if err := c.get(c.buildURL(fmt.Sprintf("/%s/%d/recommendations", mediaType, tmdbID), params), &resp); err != nil {
		return nil, err
	}
	return resp.Results, nil
}

// DiscoverByGenre discovers content by genre ID.
func (c *Client) DiscoverByGenre(mediaType string, genreID int, page int) ([]SearchResult, error) {
	params := url.Values{}
	params.Set("with_genres", fmt.Sprintf("%d", genreID))
	params.Set("sort_by", "popularity.desc")
	params.Set("page", fmt.Sprintf("%d", page))

	var resp searchResponse
	if err := c.get(c.buildURL(fmt.Sprintf("/discover/%s", mediaType), params), &resp); err != nil {
		return nil, err
	}
	return resp.Results, nil
}

// Trending returns trending content.
func (c *Client) Trending(mediaType string) ([]SearchResult, error) {
	var resp searchResponse
	if err := c.get(c.buildURL(fmt.Sprintf("/trending/%s/week", mediaType), nil), &resp); err != nil {
		return nil, err
	}
	return resp.Results, nil
}

// GetDisplayTitle returns the correct display title (title for movies, name for TV).
func (r *SearchResult) GetDisplayTitle() string {
	if r.Title != "" {
		return r.Title
	}
	return r.Name
}

// GetYear extracts year from release_date or first_air_date.
func (r *SearchResult) GetYear() string {
	date := r.ReleaseDate
	if date == "" {
		date = r.FirstAir
	}
	if len(date) >= 4 {
		return date[:4]
	}
	return ""
}

// GenreNames converts genre IDs to names.
func GenreNames(ids []int) string {
	names := make([]string, 0, len(ids))
	for _, id := range ids {
		if n, ok := genreMap[id]; ok {
			names = append(names, n)
		}
	}
	return strings.Join(names, ", ")
}

// PosterURL returns the full poster URL.
func PosterURL(path string) string {
	if path == "" {
		return ""
	}
	return imageBaseURL + path
}

func (c *Client) buildURL(path string, params url.Values) string {
	if params == nil {
		params = url.Values{}
	}
	if c.AccessToken == "" && c.APIKey != "" {
		params.Set("api_key", c.APIKey)
	}
	encoded := params.Encode()
	if encoded == "" {
		return baseURL + path
	}
	return baseURL + path + "?" + encoded
}

// MaxRetries is the number of retry attempts for rate-limited requests.
const MaxRetries = 3

// IsNetworkError returns true if the error is a network-level failure (DNS, connection, timeout).
func IsNetworkError(err error) bool {
	if err == nil {
		return false
	}
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return true
	}
	// Check for common network error strings as fallback
	msg := err.Error()
	return strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "no such host") ||
		strings.Contains(msg, "network is unreachable")
}

// IsTimeoutError returns true if the error is specifically a timeout.
func IsTimeoutError(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	return false
}

func (c *Client) get(reqURL string, target interface{}) error {
	var lastErr error
	for attempt := 0; attempt <= MaxRetries; attempt++ {
		req, reqErr := http.NewRequest(http.MethodGet, reqURL, nil)
		if reqErr != nil {
			lastErr = fmt.Errorf("build request failed: %w", reqErr)
			backoff(attempt)
			continue
		}
		if c.AccessToken != "" {
			req.Header.Set("Authorization", "Bearer "+c.AccessToken)
		}
		req.Header.Set("Accept", "application/json")

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			// Classify network errors per spec §1.4 and §4
			if IsTimeoutError(err) {
				return fmt.Errorf("%w: check your internet connection", ErrTimeout)
			}
			if IsNetworkError(err) {
				return fmt.Errorf("%w: %v", ErrNetworkError, err)
			}
			lastErr = fmt.Errorf("HTTP request failed: %w", err)
			backoff(attempt)
			continue
		}

		// Handle specific HTTP status codes per spec
		switch {
		case resp.StatusCode == 401:
			resp.Body.Close()
			return fmt.Errorf("%w. Run: movie config set tmdb_api_key YOUR_KEY", ErrAuthInvalid)

		case resp.StatusCode == 429:
			// Rate limit — retry per spec §1.1
			resp.Body.Close()
			lastErr = ErrRateLimited
			retryAfter := resp.Header.Get("Retry-After")
			delay := 2 * time.Second
			if secs, parseErr := time.ParseDuration(retryAfter + "s"); parseErr == nil && secs > 0 {
				delay = secs
			}
			time.Sleep(delay)
			continue

		case resp.StatusCode >= 500:
			// Server errors — retry once per spec §1.3
			resp.Body.Close()
			lastErr = fmt.Errorf("%w (HTTP %d)", ErrServerError, resp.StatusCode)
			if attempt == 0 {
				delay := 3 * time.Second
				if resp.StatusCode == 502 || resp.StatusCode == 503 || resp.StatusCode == 504 {
					delay = 5 * time.Second
				}
				time.Sleep(delay)
				continue
			}
			return lastErr

		case resp.StatusCode != 200:
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return fmt.Errorf("TMDb API error %d: %s", resp.StatusCode, string(body))
		}

		err = json.NewDecoder(resp.Body).Decode(target)
		resp.Body.Close()
		return err
	}
	return fmt.Errorf("TMDb request failed after %d retries: %w", MaxRetries, lastErr)
}

// backoff sleeps for exponential duration: 1s, 2s, 4s, ...
func backoff(attempt int) {
	if attempt >= MaxRetries {
		return
	}
	d := time.Duration(1<<uint(attempt)) * time.Second
	time.Sleep(d)
}

// genreMap maps TMDb genre IDs to names (combined movie + TV).
var genreMap = map[int]string{
	28: "Action", 12: "Adventure", 16: "Animation", 35: "Comedy",
	80: "Crime", 99: "Documentary", 18: "Drama", 10751: "Family",
	14: "Fantasy", 36: "History", 27: "Horror", 10402: "Music",
	9648: "Mystery", 10749: "Romance", 878: "Science Fiction",
	10770: "TV Movie", 53: "Thriller", 10752: "War", 37: "Western",
	10759: "Action & Adventure", 10762: "Kids", 10763: "News",
	10764: "Reality", 10765: "Sci-Fi & Fantasy", 10766: "Soap",
	10767: "Talk", 10768: "War & Politics",
}

// GenreNameToID returns a reverse map of genre name → TMDb genre ID.
func GenreNameToID() map[string]int {
	m := make(map[string]int, len(genreMap))
	for id, name := range genreMap {
		m[name] = id
	}
	return m
}
