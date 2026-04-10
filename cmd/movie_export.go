// movie_export.go — movie export
// Dumps the media table as JSON to ./data/json/export/media.json.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/mahin/mahin-cli-v2/db"
)

var exportOutput string

var movieExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export media library as JSON",
	Long: `Dump the entire media table to a JSON file.

Default output: ./data/json/export/media.json

Examples:
  movie export                              # Export to default path
  movie export -o ~/Desktop/library.json    # Custom output path`,
	Run: runExport,
}

func init() {
	movieExportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "Output file path (default: ./data/json/export/media.json)")
}

// exportMediaJSON mirrors db.Media with JSON tags for clean output.
type exportMediaJSON struct {
	Title            string  `json:"title"`
	CleanTitle       string  `json:"clean_title"`
	Type             string  `json:"type"`
	ImdbID           string  `json:"imdb_id,omitempty"`
	Description      string  `json:"description,omitempty"`
	Genre            string  `json:"genre,omitempty"`
	Director         string  `json:"director,omitempty"`
	CastList         string  `json:"cast_list,omitempty"`
	ThumbnailPath    string  `json:"thumbnail_path,omitempty"`
	OriginalFileName string  `json:"original_file_name,omitempty"`
	OriginalFilePath string  `json:"original_file_path,omitempty"`
	CurrentFilePath  string  `json:"current_file_path,omitempty"`
	FileExtension    string  `json:"file_extension,omitempty"`
	ID               int64   `json:"id"`
	FileSize         int64   `json:"file_size,omitempty"`
	ImdbRating       float64 `json:"imdb_rating,omitempty"`
	TmdbRating       float64 `json:"tmdb_rating,omitempty"`
	Popularity       float64 `json:"popularity,omitempty"`
	Year             int     `json:"year"`
	TmdbID           int     `json:"tmdb_id"`
}

func toExportMediaJSON(m db.Media) exportMediaJSON {
	return exportMediaJSON{
		ID: m.ID, Title: m.Title, CleanTitle: m.CleanTitle,
		Year: m.Year, Type: m.Type, TmdbID: m.TmdbID, ImdbID: m.ImdbID,
		Description: m.Description, ImdbRating: m.ImdbRating, TmdbRating: m.TmdbRating,
		Popularity: m.Popularity, Genre: m.Genre, Director: m.Director,
		CastList: m.CastList, ThumbnailPath: m.ThumbnailPath,
		OriginalFileName: m.OriginalFileName, OriginalFilePath: m.OriginalFilePath,
		CurrentFilePath: m.CurrentFilePath, FileExtension: m.FileExtension,
		FileSize: m.FileSize,
	}
}

func runExport(cmd *cobra.Command, args []string) {
	database, err := db.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Database error: %v\n", err)
		return
	}
	defer database.Close()

	// Fetch all media (large limit)
	items, err := database.ListMedia(0, 100000)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to read media: %v\n", err)
		return
	}

	if len(items) == 0 {
		fmt.Println("📭 No media to export. Run 'movie scan <folder>' first.")
		return
	}

	// Convert to JSON-friendly structs
	out := make([]exportMediaJSON, len(items))
	for i := range items {
		out[i] = toExportMediaJSON(items[i])
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ JSON encoding error: %v\n", err)
		return
	}

	// Determine output path
	outPath := exportOutput
	if outPath == "" {
		outPath = filepath.Join(".", "data", "json", "export", "media.json")
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Cannot create directory: %v\n", err)
		return
	}

	if err := os.WriteFile(outPath, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to write file: %v\n", err)
		return
	}

	fmt.Printf("✅ Exported %d items → %s\n", len(items), outPath)
}
