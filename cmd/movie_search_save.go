// movie_search_save.go — save selected search result to database and print summary
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/alimtvnetwork/movie-cli-v3/cleaner"
	"github.com/alimtvnetwork/movie-cli-v3/db"
	"github.com/alimtvnetwork/movie-cli-v3/tmdb"
)

// saveSearchResult builds a Media record from the selected TMDb result,
// fetches full details, downloads the thumbnail, and persists to the database.
func saveSearchResult(client *tmdb.Client, database *db.DB, selected tmdb.SearchResult) {
	title := selected.GetDisplayTitle()
	year := selected.GetYear()
	yearInt := 0
	if year != "" {
		yearInt, _ = strconv.Atoi(year)
	}

	fmt.Printf("\n⏳ Fetching full details for: %s...\n", title)

	m := &db.Media{
		Title:       title,
		CleanTitle:  title,
		Year:        yearInt,
		TmdbID:      selected.ID,
		TmdbRating:  selected.VoteAvg,
		Popularity:  selected.Popularity,
		Description: selected.Overview,
		Genre:       tmdb.GenreNames(selected.GenreIDs),
	}

	if selected.MediaType == "movie" || selected.MediaType == "" {
		m.Type = "movie"
		fetchMovieDetails(client, selected.ID, m)
	} else if selected.MediaType == "tv" {
		m.Type = "tv"
		fetchTVDetails(client, selected.ID, m)
	}

	downloadSearchThumbnail(client, database, selected, m)
	persistMedia(database, m)
	printSavedSummary(m)
}

// downloadSearchThumbnail downloads the poster image for a search result.
func downloadSearchThumbnail(client *tmdb.Client, database *db.DB, selected tmdb.SearchResult, m *db.Media) {
	if selected.PosterPath == "" {
		return
	}

	slug := cleaner.ToSlug(m.CleanTitle)
	if m.Year > 0 {
		slug += "-" + strconv.Itoa(m.Year)
	}

	thumbDir := filepath.Join(database.BasePath, "thumbnails", slug)
	if mkdirErr := os.MkdirAll(thumbDir, 0755); mkdirErr != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Cannot create thumbnail dir: %v\n", mkdirErr)
	}

	thumbPath := filepath.Join(thumbDir, slug+".jpg")
	if dlErr := client.DownloadPoster(selected.PosterPath, thumbPath); dlErr != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Thumbnail download failed: %v\n", dlErr)
	} else {
		m.ThumbnailPath = thumbPath
		fmt.Println("🖼️  Thumbnail saved")
	}
}

// persistMedia inserts (or updates) the media record in the database.
func persistMedia(database *db.DB, m *db.Media) {
	jsonDir := filepath.Join(database.BasePath, "json", m.Type)
	if mkdirErr := os.MkdirAll(jsonDir, 0755); mkdirErr != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Cannot create JSON dir: %v\n", mkdirErr)
	}

	_, insertErr := database.InsertMedia(m)
	if insertErr != nil {
		if m.TmdbID > 0 {
			updateErr := database.UpdateMediaByTmdbID(m)
			if updateErr == nil {
				fmt.Printf("🔄 Updated existing record for: %s\n", m.Title)
			} else {
				fmt.Fprintf(os.Stderr, "❌ DB error: %v\n", updateErr)
			}
		} else {
			fmt.Fprintf(os.Stderr, "❌ DB error: %v\n", insertErr)
		}
	}
}

// printSavedSummary prints the saved media summary to stdout.
func printSavedSummary(m *db.Media) {
	typeIcon := "🎬"
	typeLabel := "Movie"
	folder := "movie"
	if m.Type == "tv" {
		typeIcon = "📺"
		typeLabel = "TV Show"
		folder = "tv"
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("✅ Saved to database!\n\n")
	fmt.Printf("  %s  %s (%s)\n", typeIcon, m.Title, typeLabel)
	fmt.Printf("  📅  Year: %d\n", m.Year)
	fmt.Printf("  ⭐  Rating: %.1f\n", m.TmdbRating)
	fmt.Printf("  🎭  Genre: %s\n", m.Genre)

	if m.Director != "" {
		fmt.Printf("  🎬  Director: %s\n", m.Director)
	}
	if m.CastList != "" {
		fmt.Printf("  👥  Cast: %s\n", m.CastList)
	}
	if m.Description != "" {
		desc := m.Description
		if len(desc) > 150 {
			desc = desc[:147] + "..."
		}
		fmt.Printf("  📝  %s\n", desc)
	}

	fmt.Printf("  📁  Stored in: %s/ folder\n", folder)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
