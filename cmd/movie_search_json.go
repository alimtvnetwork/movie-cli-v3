// movie_search_json.go — JSON output for movie search
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alimtvnetwork/movie-cli-v3/tmdb"
)

// searchJSONItem represents a single TMDb search result in JSON output.
type searchJSONItem struct {
	Index      int     `json:"index"`
	Title      string  `json:"title"`
	Year       string  `json:"year,omitempty"`
	Type       string  `json:"type"`
	TmdbID     int     `json:"tmdb_id"`
	Rating     float64 `json:"rating"`
	Popularity float64 `json:"popularity"`
	Overview   string  `json:"overview,omitempty"`
	PosterPath string  `json:"poster_path,omitempty"`
	GenreIDs   []int   `json:"genre_ids,omitempty"`
}

// printSearchResultsJSON outputs search results as a JSON array to stdout.
func printSearchResultsJSON(results []tmdb.SearchResult) {
	items := make([]searchJSONItem, 0, len(results))
	for i, r := range results {
		if i >= 15 {
			break
		}
		mediaType := r.MediaType
		if mediaType == "" {
			mediaType = "movie"
		}
		items = append(items, searchJSONItem{
			Index:      i + 1,
			Title:      r.GetDisplayTitle(),
			Year:       r.GetYear(),
			Type:       mediaType,
			TmdbID:     r.ID,
			Rating:     r.VoteAvg,
			Popularity: r.Popularity,
			Overview:   r.Overview,
			PosterPath: r.PosterPath,
			GenreIDs:   r.GenreIDs,
		})
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(items); err != nil {
		fmt.Fprintf(os.Stderr, "❌ JSON encode error: %v\n", err)
	}
}
