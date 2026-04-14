// movie_scan_html.go — generates report.html from the embedded template after a scan.
package cmd

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/movie-cli-v3/db"
	"github.com/alimtvnetwork/movie-cli-v3/templates"
)

const defaultRESTPort = 8086

// htmlReportData is the data passed to the HTML template.
type htmlReportData struct {
	ScannedFolder string
	ScannedAt     string
	TotalFiles    int
	Movies        int
	TVShows       int
	Skipped       int
	Port          int
	Items         []htmlReportItem
}

// htmlReportItem represents a single media item in the HTML report.
type htmlReportItem struct {
	ID            int64
	Title         string
	Year          int
	Type          string
	Genre         string
	GenreList     []string
	Director      string
	CastList      string
	Description   string
	Tagline       string
	TmdbRating    float64
	ImdbRating    float64
	Runtime       int
	ThumbnailPath string
}

// writeHTMLReport generates report.html in the output directory.
func writeHTMLReport(outputDir, scanDir string, items []db.Media, total, movies, tv, skipped int) error {
	tmplBytes, err := templates.FS.ReadFile("report.html")
	if err != nil {
		return fmt.Errorf("read template: %w", err)
	}

	tmpl, err := template.New("report").Parse(string(tmplBytes))
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	reportItems := make([]htmlReportItem, 0, len(items))
	for _, m := range items {
		var genres []string
		if m.Genre != "" {
			for _, g := range strings.Split(m.Genre, ",") {
				g = strings.TrimSpace(g)
				if g != "" {
					genres = append(genres, g)
				}
			}
		}
		reportItems = append(reportItems, htmlReportItem{
			ID:            m.ID,
			Title:         m.Title,
			Year:          m.Year,
			Type:          m.Type,
			Genre:         m.Genre,
			GenreList:     genres,
			Director:      m.Director,
			CastList:      m.CastList,
			Description:   m.Description,
			Tagline:       m.Tagline,
			TmdbRating:    m.TmdbRating,
			ImdbRating:    m.ImdbRating,
			Runtime:       m.Runtime,
			ThumbnailPath: m.ThumbnailPath,
		})
	}

	data := htmlReportData{
		ScannedFolder: scanDir,
		ScannedAt:     time.Now().Format("2006-01-02 15:04:05"),
		TotalFiles:    total,
		Movies:        movies,
		TVShows:       tv,
		Skipped:       skipped,
		Port:          defaultRESTPort,
		Items:         reportItems,
	}

	outPath := filepath.Join(outputDir, "report.html")
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return nil
}
