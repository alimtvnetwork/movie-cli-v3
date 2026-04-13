// movie_scan_json_output.go — JSON stdout output for movie scan --format json
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/alimtvnetwork/movie-cli-v3/cleaner"
	"github.com/alimtvnetwork/movie-cli-v3/db"
)

// scanJSONOutput is the top-level JSON written to stdout with --format json.
type scanJSONOutput struct {
	ScannedFolder string           `json:"scanned_folder"`
	ScannedAt     string           `json:"scanned_at"`
	DryRun        bool             `json:"dry_run"`
	TotalFiles    int              `json:"total_files"`
	Movies        int              `json:"movies"`
	TVShows       int              `json:"tv_shows"`
	Skipped       int              `json:"skipped"`
	Items         []scanJSONItem   `json:"items"`
}

// scanJSONItem is one item in the JSON output.
type scanJSONItem struct {
	FileName   string  `json:"file_name"`
	FilePath   string  `json:"file_path"`
	CleanTitle string  `json:"clean_title"`
	Year       int     `json:"year,omitempty"`
	Type       string  `json:"type"`
	TmdbID     int     `json:"tmdb_id,omitempty"`
	TmdbRating float64 `json:"tmdb_rating,omitempty"`
	Genre      string  `json:"genre,omitempty"`
	Status     string  `json:"status"` // "new", "skipped", "updated"
}

// buildDryRunJSONItems creates JSON items from video files in dry-run mode.
func buildDryRunJSONItems(videoFiles []videoFile) (items []scanJSONItem, movies, tvShows int) {
	for _, vf := range videoFiles {
		result := cleaner.Clean(vf.Name)
		items = append(items, scanJSONItem{
			FileName:   vf.Name,
			FilePath:   vf.FullPath,
			CleanTitle: result.CleanTitle,
			Year:       result.Year,
			Type:       result.Type,
			Status:     "new",
		})
		if result.Type == "movie" {
			movies++
		} else {
			tvShows++
		}
	}
	return
}

// buildMediaJSONItem creates a JSON item from a processed Media record.
func buildMediaJSONItem(m *db.Media, status string) scanJSONItem {
	return scanJSONItem{
		FileName:   m.OriginalFileName,
		FilePath:   m.CurrentFilePath,
		CleanTitle: m.CleanTitle,
		Year:       m.Year,
		Type:       m.Type,
		TmdbID:     m.TmdbID,
		TmdbRating: m.TmdbRating,
		Genre:      m.Genre,
		Status:     status,
	}
}

// printScanJSON writes the full scan result as JSON to stdout.
func printScanJSON(scanDir string, items []scanJSONItem, totalFiles, movies, tvShows, skipped int) {
	output := scanJSONOutput{
		ScannedFolder: scanDir,
		ScannedAt:     time.Now().UTC().Format(time.RFC3339),
		DryRun:        scanDryRun,
		TotalFiles:    totalFiles,
		Movies:        movies,
		TVShows:       tvShows,
		Skipped:       skipped,
		Items:         items,
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(output); err != nil {
		fmt.Fprintf(os.Stderr, "❌ JSON encode error: %v\n", err)
	}
}
