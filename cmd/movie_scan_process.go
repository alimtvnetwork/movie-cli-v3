// movie_scan_process.go — per-file processing and TMDb enrichment for movie scan
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/alimtvnetwork/movie-cli-v3/cleaner"
	"github.com/alimtvnetwork/movie-cli-v3/db"
	"github.com/alimtvnetwork/movie-cli-v3/errlog"
	"github.com/alimtvnetwork/movie-cli-v3/tmdb"
)

// processVideoFile handles a single video file: clean, check DB, fetch TMDb, insert, write JSON.
// Returns true if the file was processed (even if skipped), false on hard errors.
func processVideoFile(
	vf videoFile,
	database *db.DB,
	client *tmdb.Client,
	hasTMDb bool,
	outputDir string,
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
		errlog.Warn("DB search error for '%s': %v", result.CleanTitle, searchErr)
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
		errlog.Error("cannot stat file %s: %v", vf.FullPath, fiErr)
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
	if hasTMDb {
		enrichFromTMDb(client, database, m, result, outputDir)
	}

	// Insert into database
	_, insertErr := database.InsertMedia(m)
	if insertErr != nil {
		if m.TmdbID > 0 {
			if updateErr := database.UpdateMediaByTmdbID(m); updateErr != nil {
				errlog.Error("DB update error for '%s': %v", m.Title, updateErr)
			}
		} else {
			errlog.Error("DB insert error for '%s': %v", m.Title, insertErr)
		}
	}

	if jsonErr := writeMediaJSON(outputDir, m); jsonErr != nil {
		errlog.Warn("JSON write error for '%s': %v", m.Title, jsonErr)
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
func enrichFromTMDb(client *tmdb.Client, database *db.DB, m *db.Media, result cleaner.Result, outputDir string) {
	searchQuery := result.CleanTitle
	if result.Year > 0 {
		searchQuery += " " + strconv.Itoa(result.Year)
	}

	tmdbResults, tmdbErr := client.SearchMulti(searchQuery)
	if tmdbErr != nil || len(tmdbResults) == 0 {
		errlog.Warn("no TMDb match for '%s': %v", searchQuery, tmdbErr)
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

	// Download thumbnail into the .movie-output/thumbnails folder
	if best.PosterPath != "" {
		slug := cleaner.ToSlug(m.CleanTitle)
		if m.Year > 0 {
			slug += "-" + strconv.Itoa(m.Year)
		}
		thumbFileName := slug + "-" + strconv.Itoa(m.TmdbID) + ".jpg"
		thumbDir := filepath.Join(outputDir, "thumbnails")
		if mkdirErr := os.MkdirAll(thumbDir, 0755); mkdirErr != nil {
			errlog.Error("cannot create thumbnail dir %s: %v", thumbDir, mkdirErr)
		}
		thumbPath := filepath.Join(thumbDir, thumbFileName)
		if dlErr := client.DownloadPoster(best.PosterPath, thumbPath); dlErr != nil {
			errlog.Warn("thumbnail download failed for '%s': %v", m.CleanTitle, dlErr)
		} else {
			// Store relative path for HTML (thumbnails/name-id.jpg)
			m.ThumbnailPath = "thumbnails/" + thumbFileName
			fmt.Println("     🖼️  Thumbnail saved")
		}
	}

	fmt.Printf("     ✅ TMDb: %s (⭐ %.1f)\n", m.Title, m.TmdbRating)
}
