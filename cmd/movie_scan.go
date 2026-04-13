// movie_scan.go — movie scan [folder] [--recursive]
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/movie-cli-v3/cleaner"
	"github.com/alimtvnetwork/movie-cli-v3/db"
	"github.com/alimtvnetwork/movie-cli-v3/tmdb"
)

var scanRecursive bool

var movieScanCmd = &cobra.Command{
	Use:   "scan [folder]",
	Short: "Scan a folder for movies and TV shows",
	Long: `Scans a folder for video files, cleans filenames, fetches metadata
from TMDb, downloads thumbnails, and stores everything in the database.

If no folder is specified, scans the current working directory.
Use --recursive (-r) to scan all subdirectories recursively.

Results are saved to .movie-output/ inside the scanned folder, including:
  - summary.json   — full scan report with categories, counts, and per-item metadata
  - json/movie/    — individual JSON files per movie
  - json/tv/       — individual JSON files per TV show

Examples:
  movie scan                     Scan current directory (top-level)
  movie scan ~/Movies            Scan specific folder
  movie scan -r                  Scan current directory recursively
  movie scan ~/Movies --recursive`,
	Args: cobra.MaximumNArgs(1),
	Run:  runMovieScan,
}

func init() {
	movieScanCmd.Flags().BoolVarP(&scanRecursive, "recursive", "r", false,
		"scan all subdirectories recursively")
}

func runMovieScan(cmd *cobra.Command, args []string) {
	database, err := db.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Database error: %v\n", err)
		return
	}
	defer database.Close()

	// Determine scan folder
	scanDir := ""
	if len(args) > 0 {
		scanDir = args[0]
	} else {
		// Default to current working directory
		scanDir, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Cannot determine current directory: %v\n", err)
			return
		}
		fmt.Printf("📂 No folder specified — scanning current directory\n\n")
	}

	// Expand ~ to home
	if strings.HasPrefix(scanDir, "~") {
		home, homeErr := os.UserHomeDir()
		if homeErr != nil {
			fmt.Fprintf(os.Stderr, "❌ Cannot determine home directory: %v\n", homeErr)
			return
		}
		scanDir = filepath.Join(home, scanDir[1:])
	}

	// Check folder exists
	info, statErr := os.Stat(scanDir)
	if statErr != nil || !info.IsDir() {
		fmt.Fprintf(os.Stderr, "❌ Folder not found: %s\n", scanDir)
		return
	}

	// Get TMDb API key
	apiKey, cfgErr := database.GetConfig("tmdb_api_key")
	if cfgErr != nil && cfgErr.Error() != "sql: no rows in result set" {
		fmt.Fprintf(os.Stderr, "⚠️  Config read error: %v\n", cfgErr)
	}
	if apiKey == "" {
		apiKey = os.Getenv("TMDB_API_KEY")
	}
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "⚠️  No TMDb API key configured.")
		fmt.Fprintln(os.Stderr, "   Set it with: movie config set tmdb_api_key YOUR_KEY")
		fmt.Fprintln(os.Stderr, "   Or set TMDB_API_KEY environment variable.")
		fmt.Fprintln(os.Stderr, "   Scanning will proceed without metadata fetching.")
		fmt.Println()
	}

	client := tmdb.NewClient(apiKey)

	// Set up .movie-output directory inside the scanned folder
	outputDir := filepath.Join(scanDir, ".movie-output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Cannot create output directory: %v\n", err)
		return
	}
	jsonMovieDir := filepath.Join(outputDir, "json", "movie")
	jsonTVDir := filepath.Join(outputDir, "json", "tv")
	if err := os.MkdirAll(jsonMovieDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Cannot create json/movie dir: %v\n", err)
		return
	}
	if err := os.MkdirAll(jsonTVDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Cannot create json/tv dir: %v\n", err)
		return
	}

	fmt.Printf("🔍 Scanning: %s\n", scanDir)
	if scanRecursive {
		fmt.Println("🔄 Mode: recursive (all subdirectories)")
	}
	fmt.Printf("📁 Output:   %s\n\n", outputDir)

	var totalFiles, movieCount, tvCount, skipped int
	var scannedItems []db.Media

	// Collect video files based on scan mode
	videoFiles := collectVideoFiles(scanDir, scanRecursive)

	for _, vf := range videoFiles {
		result := processVideoFile(vf, database, client, apiKey, outputDir,
			&totalFiles, &movieCount, &tvCount, &skipped, &scannedItems)
		if !result {
			continue
		}
	}

	// Log scan history
	if histErr := database.InsertScanHistory(scanDir, totalFiles, movieCount, tvCount); histErr != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Could not log scan history: %v\n", histErr)
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("📊 Scan Complete!\n")
	fmt.Printf("   Total files: %d\n", totalFiles)
	fmt.Printf("   Movies:      %d\n", movieCount)
	fmt.Printf("   TV Shows:    %d\n", tvCount)
	if skipped > 0 {
		fmt.Printf("   Skipped:     %d (already in DB)\n", skipped)
	}
	fmt.Printf("   Output:      %s\n", outputDir)

	// Write summary.json to .movie-output/
	if summaryErr := writeScanSummary(outputDir, scanDir, scannedItems,
		totalFiles, movieCount, tvCount, skipped); summaryErr != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Could not write summary.json: %v\n", summaryErr)
	} else {
		fmt.Printf("\n📋 Summary saved: %s\n", filepath.Join(outputDir, "summary.json"))
	}
}

// videoFile holds a discovered video file's display name and full path.
type videoFile struct {
	Name     string // display name used for cleaning (dir name or filename)
	FullPath string // absolute path to the actual video file
}

// collectVideoFiles finds video files in the given directory.
// When recursive is true, it walks all subdirectories.
func collectVideoFiles(scanDir string, recursive bool) []videoFile {
	var files []videoFile

	if recursive {
		_ = filepath.WalkDir(scanDir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				fmt.Fprintf(os.Stderr, "  ⚠️  Cannot access %s: %v\n", path, err)
				return nil // continue walking
			}
			// Skip .movie-output and hidden directories
			if d.IsDir() {
				base := d.Name()
				if base == ".movie-output" || (strings.HasPrefix(base, ".") && base != ".") {
					return filepath.SkipDir
				}
				return nil
			}
			if cleaner.IsVideoFile(d.Name()) {
				// Use parent directory name if it differs from scanDir, else use filename
				parentDir := filepath.Dir(path)
				name := d.Name()
				if parentDir != scanDir {
					name = filepath.Base(parentDir)
				}
				files = append(files, videoFile{Name: name, FullPath: path})
			}
			return nil
		})
	} else {
		entries, readErr := os.ReadDir(scanDir)
		if readErr != nil {
			fmt.Fprintf(os.Stderr, "❌ Cannot read folder: %v\n", readErr)
			return nil
		}
		for _, entry := range entries {
			name := entry.Name()
			fullPath := filepath.Join(scanDir, name)

			if entry.IsDir() {
				// Look for first video file inside the subdirectory
				subEntries, subErr := os.ReadDir(fullPath)
				if subErr != nil {
					fmt.Fprintf(os.Stderr, "  ⚠️  Cannot read subdirectory %s: %v\n", name, subErr)
					continue
				}
				for _, sub := range subEntries {
					if !sub.IsDir() && cleaner.IsVideoFile(sub.Name()) {
						files = append(files, videoFile{
							Name:     entry.Name(),
							FullPath: filepath.Join(fullPath, sub.Name()),
						})
						break
					}
				}
			} else if cleaner.IsVideoFile(name) {
				files = append(files, videoFile{Name: name, FullPath: fullPath})
			}
		}
	}

	return files
}

// processVideoFile handles a single video file: clean, check DB, fetch TMDb, insert, write JSON.
// Returns true if the file was processed (even if skipped), false on hard errors.
func processVideoFile(
	vf videoFile,
	database *db.DB,
	client *tmdb.Client,
	apiKey, outputDir string,
	totalFiles, movieCount, tvCount, skipped *int,
	scannedItems *[]db.Media,
) bool {
	*totalFiles++

	result := cleaner.Clean(vf.Name)
	fmt.Printf("  📄 %s\n", vf.Name)
	fmt.Printf("     → %s", result.CleanTitle)
	if result.Year > 0 {
		fmt.Printf(" (%d)", result.Year)
	}
	fmt.Printf(" [%s]\n", result.Type)

	// Check if already in DB by path
	existing, searchErr := database.SearchMedia(result.CleanTitle)
	if searchErr != nil {
		fmt.Fprintf(os.Stderr, "     ⚠️  DB search error: %v\n", searchErr)
	}
	for i := range existing {
		if existing[i].OriginalFilePath == vf.FullPath {
			fmt.Println("     ⏩ Already in database, skipping")
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

	if m.Type == "movie" {
		*movieCount++
	} else {
		*tvCount++
	}
	fmt.Println()
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
