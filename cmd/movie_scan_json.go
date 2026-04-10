// movie_scan_json.go — JSON metadata file generation during scan
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mahin/mahin-cli-v2/cleaner"
	"github.com/mahin/mahin-cli-v2/db"
)

// scanMediaJSON is the JSON representation written to disk.
type scanMediaJSON struct {
	Title            string  `json:"title"`
	CleanTitle       string  `json:"clean_title"`
	Type             string  `json:"type"`
	ImdbID           string  `json:"imdb_id,omitempty"`
	Description      string  `json:"description,omitempty"`
	Genre            string  `json:"genre,omitempty"`
	Director         string  `json:"director,omitempty"`
	CastList         string  `json:"cast_list,omitempty"`
	ThumbnailPath    string  `json:"thumbnail_path,omitempty"`
	OriginalFileName string  `json:"original_file_name"`
	OriginalFilePath string  `json:"original_file_path"`
	CurrentFilePath  string  `json:"current_file_path"`
	FileExtension    string  `json:"file_extension"`
	GeneratedAt      string  `json:"generated_at"`
	ImdbRating       float64 `json:"imdb_rating,omitempty"`
	TmdbRating       float64 `json:"tmdb_rating,omitempty"`
	Popularity       float64 `json:"popularity,omitempty"`
	FileSize         int64   `json:"file_size"`
	Year             int     `json:"year,omitempty"`
	TmdbID           int     `json:"tmdb_id,omitempty"`
}

// writeMediaJSON writes a JSON metadata file for the given media record.
// Files are saved to <basePath>/json/movie/<slug>.json or json/tv/<slug>.json.
func writeMediaJSON(basePath string, m *db.Media) error {
	subDir := "movie"
	if m.Type == "tv" {
		subDir = "tv"
	}

	slug := cleaner.ToSlug(m.CleanTitle)
	if m.Year > 0 {
		slug += "-" + strconv.Itoa(m.Year)
	}

	dir := filepath.Join(basePath, "json", subDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("cannot create json dir: %w", err)
	}

	data := scanMediaJSON{
		Title:            m.Title,
		CleanTitle:       m.CleanTitle,
		Year:             m.Year,
		Type:             m.Type,
		TmdbID:           m.TmdbID,
		ImdbID:           m.ImdbID,
		Description:      m.Description,
		ImdbRating:       m.ImdbRating,
		TmdbRating:       m.TmdbRating,
		Popularity:       m.Popularity,
		Genre:            m.Genre,
		Director:         m.Director,
		CastList:         m.CastList,
		ThumbnailPath:    m.ThumbnailPath,
		OriginalFileName: m.OriginalFileName,
		OriginalFilePath: m.OriginalFilePath,
		CurrentFilePath:  m.CurrentFilePath,
		FileExtension:    m.FileExtension,
		FileSize:         m.FileSize,
		GeneratedAt:      time.Now().UTC().Format(time.RFC3339),
	}

	jsonPath := filepath.Join(dir, slug+".json")
	f, err := os.Create(jsonPath)
	if err != nil {
		return fmt.Errorf("cannot create json file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("cannot write json: %w", err)
	}

	fmt.Printf("     📝 JSON metadata saved: %s\n", jsonPath)
	return nil
}
