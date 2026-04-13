// movie_scan_process.go — per-file processing and TMDb enrichment for movie scan
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

// processVideoFile handles a single video file: clean, check DB, fetch TMDb, insert, write JSON.
// Returns true if the file was processed (even if skipped), false on hard errors.
func processVideoFile(
	vf videoFile,
	database *db.DB,
	client *tmdb.Client,
	apiKey, outputDir string,
	totalFiles, movieCount, tvCount, skipped *int,
	scannedItems *[]db.Media,
	useTable bool,
) bool {
	*totalFiles++

	result := cleaner.Clean(vf.Name)
	if !useTable {
		fmt.Printf("  📄 %s\n", vf.Name)
		fmt.Printf("     → %s", result.CleanTitle)
		if result.Year > 0 {
			fmt.Printf(" (%d)", result.Year)
		}
		fmt.Printf(" [%s]\n", result.Type)
	}

	// Check if already in DB by path
	existing, searchErr := database.SearchMedia(result.CleanTitle)
	if searchErr != nil {
		fmt.Fprintf(os.Stderr, "     ⚠️  DB search error: %v\n", searchErr)
	}
	for i := range existing {
		if existing[i].OriginalFilePath == vf.FullPath {
			if useTable {
				printScanTableRow(buildMediaTableRow(*totalFiles, &db.Media{
					OriginalFileName: vf.Name,
					CleanTitle:       result.CleanTitle,
					Year:             result.Year,
					Type:             result.Type,
				}, "skipped"))
			} else {
				fmt.Println("     ⏩ Already in database, skipping")
			}
			*skipped++
			if result.Type == "movie" {
				*movieCount++
			} else {
				*tvCount++
			}
			return true
		}
	}

	fi, fiErr := os.Stat(vf.FullPath)
	if fiErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠️  Cannot stat file: %v\n", fiErr)
		return false
	}

	m := &db.Media{
		Title:            result.CleanTitle,
		CleanTitle:       result.CleanTitle,
		Year:             result.Year,
		Type:             result.Type,
		OriginalFileName: vf.Name,
		OriginalFilePath: vf.FullPath,
		CurrentFilePath:  vf.FullPath,
		FileExtension:    result.Extension,
	}
	if fi != nil {
		m.FileSize = fi.Size()
	}

	// Fetch metadata from TMDb
	if apiKey != "" {
		enrichFromTMDb(client, database, m, result)
	}

	// Insert into database
	_, insertErr := database.InsertMedia(m)
	if insertErr != nil {
		if m.TmdbID > 0 {
			if updateErr := database.UpdateMediaByTmdbID(m); updateErr != nil {
				fmt.Fprintf(os.Stderr, "     ⚠️  DB update error: %v\n", updateErr)
			}
		} else {
			fmt.Fprintf(os.Stderr, "     ❌ DB error: %v\n", insertErr)
		}
	}

	if jsonErr := writeMediaJSON(outputDir, m); jsonErr != nil {
		fmt.Fprintf(os.Stderr, "     ⚠️  JSON write error: %v\n", jsonErr)
	}

	*scannedItems = append(*scannedItems, *m)

	if useTable {
		printScanTableRow(buildMediaTableRow(*totalFiles, m, "new"))
	}

	if m.Type == "movie" {
		*movieCount++
	} else {
		*tvCount++
	}
	if !useTable {
		fmt.Println()
	}
	return true
}

// enrichFromTMDb fetches metadata, details, and thumbnail from TMDb.
func enrichFromTMDb(client *tmdb.Client, database *db.DB, m *db.Media, result cleaner.Result) {
	searchQuery := result.CleanTitle
	if result.Year > 0 {
		searchQuery += " " + strconv.Itoa(result.Year)
	}

	tmdbResults, tmdbErr := client.SearchMulti(searchQuery)
	if tmdbErr != nil || len(tmdbResults) == 0 {
		fmt.Println("     ⚠️  No TMDb match found")
		return
	}

	best := tmdbResults[0]
	m.TmdbID = best.ID
	m.TmdbRating = best.VoteAvg
	m.Popularity = best.Popularity
	m.Description = best.Overview
	m.Genre = tmdb.GenreNames(best.GenreIDs)

	if best.MediaType == "movie" || best.MediaType == "" {
		m.Type = "movie"
		fetchMovieDetails(client, best.ID, m)
	} else if best.MediaType == "tv" {
		m.Type = "tv"
		fetchTVDetails(client, best.ID, m)
	}

	// Download thumbnail
	if best.PosterPath != "" {
		slug := cleaner.ToSlug(m.CleanTitle)
		if m.Year > 0 {
			slug += "-" + strconv.Itoa(m.Year)
		}
		thumbDir := filepath.Join(database.BasePath, "thumbnails", slug)
		if mkdirErr := os.MkdirAll(thumbDir, 0755); mkdirErr != nil {
			fmt.Fprintf(os.Stderr, "     ⚠️  Cannot create thumbnail dir: %v\n", mkdirErr)
		}
		thumbPath := filepath.Join(thumbDir, slug+".jpg")
		if dlErr := client.DownloadPoster(best.PosterPath, thumbPath); dlErr != nil {
			fmt.Fprintf(os.Stderr, "     ⚠️  Thumbnail download failed: %v\n", dlErr)
		} else {
			m.ThumbnailPath = thumbPath
			fmt.Println("     🖼️  Thumbnail saved")
		}
	}

	fmt.Printf("     ✅ TMDb: %s (⭐ %.1f)\n", m.Title, m.TmdbRating)
}
