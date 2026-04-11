// Package tmdb provides a client for The Movie Database (TMDb) API.
package tmdb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const baseURL = "https://api.themoviedb.org/3"
const imageBaseURL = "https://image.tmdb.org/t/p/w500"

// Client interacts with the TMDb API.
type Client struct {
	HTTPClient *http.Client
	APIKey     string
}

// NewClient creates a new TMDb client. Reads API key from env or config.
func NewClient(apiKey string) *Client {
	if apiKey == "" {
		apiKey = os.Getenv("TMDB_API_KEY")
	}
	return &Client{
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
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
	Title       string  `json:"title"`
	Overview    string  `json:"overview"`
	ReleaseDate string  `json:"release_date"`
	PosterPath  string  `json:"poster_path"`
	ImdbID      string  `json:"imdb_id"`
	Genres      []Genre `json:"genres"`
	VoteAvg     float64 `json:"vote_average"`
	Popularity  float64 `json:"popularity"`
	ID          int     `json:"id"`
	Runtime     int     `json:"runtime"`
}

// TVDetails holds detailed TV show info.
type TVDetails struct {
	Name         string  `json:"name"`
	Overview     string  `json:"overview"`
	FirstAirDate string  `json:"first_air_date"`
	PosterPath   string  `json:"poster_path"`
	Genres       []Genre `json:"genres"`
	VoteAvg      float64 `json:"vote_average"`
	Popularity   float64 `json:"popularity"`
	ID           int     `json:"id"`
	Seasons      int     `json:"number_of_seasons"`
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
	u := fmt.Sprintf("%s/search/multi?api_key=%s&query=%s&page=1",
		baseURL, c.APIKey, url.QueryEscape(query))

	var resp searchResponse
	if err := c.get(u, &resp); err != nil {
		return nil, err
	}

	// Filter to only movie/tv
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
	u := fmt.Sprintf("%s/movie/%d?api_key=%s", baseURL, tmdbID, c.APIKey)
	var d MovieDetails
	if err := c.get(u, &d); err != nil {
		return nil, err
	}
	return &d, nil
}

// GetTVDetails returns detailed info for a TV show.
func (c *Client) GetTVDetails(tmdbID int) (*TVDetails, error) {
	u := fmt.Sprintf("%s/tv/%d?api_key=%s", baseURL, tmdbID, c.APIKey)
	var d TVDetails
	if err := c.get(u, &d); err != nil {
		return nil, err
	}
	return &d, nil
}

// GetMovieCredits returns cast and crew for a movie.
func (c *Client) GetMovieCredits(tmdbID int) (*Credits, error) {
	u := fmt.Sprintf("%s/movie/%d/credits?api_key=%s", baseURL, tmdbID, c.APIKey)
	var cr Credits
	if err := c.get(u, &cr); err != nil {
		return nil, err
	}
	return &cr, nil
}

// GetTVCredits returns cast and crew for a TV show.
func (c *Client) GetTVCredits(tmdbID int) (*Credits, error) {
	u := fmt.Sprintf("%s/tv/%d/credits?api_key=%s", baseURL, tmdbID, c.APIKey)
	var cr Credits
	if err := c.get(u, &cr); err != nil {
		return nil, err
	}
	return &cr, nil
}

// DownloadPoster downloads a poster image and saves it to dst.
func (c *Client) DownloadPoster(posterPath, dst string) error {
	if posterPath == "" {
		return fmt.Errorf("no poster available")
	}

	imgURL := imageBaseURL + posterPath
	resp, err := c.HTTPClient.Get(imgURL)
	if err != nil {
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
	u := fmt.Sprintf("%s/%s/%d/recommendations?api_key=%s&page=%d",
		baseURL, mediaType, tmdbID, c.APIKey, page)
	var resp searchResponse
	if err := c.get(u, &resp); err != nil {
		return nil, err
	}
	return resp.Results, nil
}

// DiscoverByGenre discovers content by genre ID.
func (c *Client) DiscoverByGenre(mediaType string, genreID int, page int) ([]SearchResult, error) {
	u := fmt.Sprintf("%s/discover/%s?api_key=%s&with_genres=%d&sort_by=popularity.desc&page=%d",
		baseURL, mediaType, c.APIKey, genreID, page)
	var resp searchResponse
	if err := c.get(u, &resp); err != nil {
		return nil, err
	}
	return resp.Results, nil
}

// Trending returns trending content.
func (c *Client) Trending(mediaType string) ([]SearchResult, error) {
	u := fmt.Sprintf("%s/trending/%s/week?api_key=%s", baseURL, mediaType, c.APIKey)
	var resp searchResponse
	if err := c.get(u, &resp); err != nil {
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

// MaxRetries is the number of retry attempts for rate-limited requests.
const MaxRetries = 3

func (c *Client) get(reqURL string, target interface{}) error {
	var lastErr error
	for attempt := 0; attempt <= MaxRetries; attempt++ {
		resp, err := c.HTTPClient.Get(reqURL)
		if err != nil {
			lastErr = fmt.Errorf("HTTP request failed: %w", err)
			backoff(attempt)
			continue
		}

		if resp.StatusCode == 429 {
			resp.Body.Close()
			lastErr = fmt.Errorf("TMDb rate limit (429)")
			backoff(attempt)
			continue
		}

		if resp.StatusCode != 200 {
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
	// TV-specific
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
