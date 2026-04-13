// movie_info_json.go — JSON output for movie info
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alimtvnetwork/movie-cli-v3/db"
)

// infoJSONOutput represents a media item in JSON output for movie info.
type infoJSONOutput struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	CleanTitle  string  `json:"clean_title"`
	Year        int     `json:"year"`
	Type        string  `json:"type"`
	TmdbID      int     `json:"tmdb_id,omitempty"`
	ImdbID      string  `json:"imdb_id,omitempty"`
	TmdbRating  float64 `json:"tmdb_rating"`
	ImdbRating  float64 `json:"imdb_rating,omitempty"`
	Popularity  float64 `json:"popularity,omitempty"`
	Genre       string  `json:"genre,omitempty"`
	Director    string  `json:"director,omitempty"`
	Cast        string  `json:"cast,omitempty"`
	Runtime     int     `json:"runtime,omitempty"`
	Language    string  `json:"language,omitempty"`
	Tagline     string  `json:"tagline,omitempty"`
	Description string  `json:"description,omitempty"`
	TrailerURL  string  `json:"trailer_url,omitempty"`
	Budget      int64   `json:"budget,omitempty"`
	Revenue     int64   `json:"revenue,omitempty"`
	FilePath    string  `json:"file_path,omitempty"`
	FileSize    int64   `json:"file_size,omitempty"`
	Thumbnail   string  `json:"thumbnail,omitempty"`
	Source      string  `json:"source"` // "local" or "tmdb"
}

// printMediaDetailJSON outputs a media item as JSON to stdout.
func printMediaDetailJSON(m *db.Media, source string) {
	out := infoJSONOutput{
		ID:          m.ID,
		Title:       m.Title,
		CleanTitle:  m.CleanTitle,
		Year:        m.Year,
		Type:        m.Type,
		TmdbID:      m.TmdbID,
		ImdbID:      m.ImdbID,
		TmdbRating:  m.TmdbRating,
		ImdbRating:  m.ImdbRating,
		Popularity:  m.Popularity,
		Genre:       m.Genre,
		Director:    m.Director,
		Cast:        m.CastList,
		Runtime:     m.Runtime,
		Language:    m.Language,
		Tagline:     m.Tagline,
		Description: m.Description,
		TrailerURL:  m.TrailerURL,
		Budget:      m.Budget,
		Revenue:     m.Revenue,
		FilePath:    m.CurrentFilePath,
		FileSize:    m.FileSize,
		Thumbnail:   m.ThumbnailPath,
		Source:      source,
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		fmt.Fprintf(os.Stderr, "❌ JSON encode error: %v\n", err)
	}
}
