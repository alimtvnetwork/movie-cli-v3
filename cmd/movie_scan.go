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
var scanDepth int
var scanDryRun bool
var scanFormat string

var movieScanCmd = &cobra.Command{
	Use:   "scan [folder]",
	Short: "Scan a folder for movies and TV shows",
	Long: `Scans a folder for video files, cleans filenames, fetches metadata
from TMDb, downloads thumbnails, and stores everything in the database.

If no folder is specified, scans the current working directory.
Use --recursive (-r) to scan all subdirectories recursively.
Use --depth to limit how many levels deep the recursive scan goes.
Use --dry-run to preview what would be scanned without writing anything.

Results are saved to .movie-output/ inside the scanned folder, including:
  - summary.json   — full scan report with categories, counts, and per-item metadata
  - json/movie/    — individual JSON files per movie
  - json/tv/       — individual JSON files per TV show

Examples:
  movie scan                     Scan current directory (top-level)
  movie scan ~/Movies            Scan specific folder
  movie scan -r                  Scan current directory recursively
  movie scan ~/Movies --recursive
  movie scan -r --depth 2        Scan only 2 levels deep
  movie scan --dry-run            Preview files without writing to DB
  movie scan --format table       Show results as a formatted table`,
	Args: cobra.MaximumNArgs(1),
	Run:  runMovieScan,
}

func init() {
	movieScanCmd.Flags().BoolVarP(&scanRecursive, "recursive", "r", false,
		"scan all subdirectories recursively")
	movieScanCmd.Flags().IntVarP(&scanDepth, "depth", "d", 0,
		"max subdirectory depth for recursive scan (0 = unlimited)")
	movieScanCmd.Flags().BoolVar(&scanDryRun, "dry-run", false,
		"preview what would be scanned without writing to DB or .movie-output")
	movieScanCmd.Flags().StringVar(&scanFormat, "format", "default",
		"output format: default or table")
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

	// Set up .movie-output directory inside the scanned folder
	outputDir := filepath.Join(scanDir, ".movie-output")

	if !scanDryRun {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Cannot create output directory: %v\n", err)
			return
		}
		if err := os.MkdirAll(filepath.Join(outputDir, "json", "movie"), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Cannot create json/movie dir: %v\n", err)
			return
		}
		if err := os.MkdirAll(filepath.Join(outputDir, "json", "tv"), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Cannot create json/tv dir: %v\n", err)
			return
		}
	}

	fmt.Printf("🔍 Scanning: %s\n", scanDir)
	if scanDryRun {
		fmt.Println("🧪 Mode: dry run (no writes)")
	}
	if scanRecursive {
		if scanDepth > 0 {
			fmt.Printf("🔄 Mode: recursive (max depth: %d)\n", scanDepth)
		} else {
			fmt.Println("🔄 Mode: recursive (all subdirectories)")
		}
	}
	if !scanDryRun {
		fmt.Printf("📁 Output:   %s\n", outputDir)
	}
	fmt.Println()

	var totalFiles, movieCount, tvCount, skipped int
	var scannedItems []db.Media

	// Collect video files based on scan mode
	videoFiles := collectVideoFiles(scanDir, scanRecursive, scanDepth)
	useTable := scanFormat == "table"

	if useTable {
		printScanTableHeader()
	}

	if scanDryRun {
		if useTable {
			rows, mc, tc := buildDryRunTableRows(videoFiles)
			for _, row := range rows {
				printScanTableRow(row)
			}
			totalFiles = len(rows)
			movieCount = mc
			tvCount = tc
		} else {
			for _, vf := range videoFiles {
				totalFiles++
				result := cleaner.Clean(vf.Name)
				fmt.Printf("  📄 %s\n", vf.Name)
				fmt.Printf("     → %s", result.CleanTitle)
				if result.Year > 0 {
					fmt.Printf(" (%d)", result.Year)
				}
				fmt.Printf(" [%s]\n", result.Type)
				fmt.Printf("     📂 %s\n\n", vf.FullPath)
				if result.Type == "movie" {
					movieCount++
				} else {
					tvCount++
				}
			}
		}
	} else {
		client := tmdb.NewClient(apiKey)
		for _, vf := range videoFiles {
			processVideoFile(vf, database, client, apiKey, outputDir,
				&totalFiles, &movieCount, &tvCount, &skipped, &scannedItems,
				useTable)
		}
	}

	if useTable {
		printScanTableFooter()
	}

	// Log scan history
	if !scanDryRun {
		if histErr := database.InsertScanHistory(scanDir, totalFiles, movieCount, tvCount); histErr != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Could not log scan history: %v\n", histErr)
		}
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	if scanDryRun {
		fmt.Printf("📊 Dry Run Complete!\n")
	} else {
		fmt.Printf("📊 Scan Complete!\n")
	}
	fmt.Printf("   Total files: %d\n", totalFiles)
	fmt.Printf("   Movies:      %d\n", movieCount)
	fmt.Printf("   TV Shows:    %d\n", tvCount)
	if skipped > 0 {
		fmt.Printf("   Skipped:     %d (already in DB)\n", skipped)
	}
	if scanDryRun {
		fmt.Println("\n💡 Run without --dry-run to actually scan and save.")
	} else {
		fmt.Printf("   Output:      %s\n", outputDir)

		// Write summary.json to .movie-output/
		if summaryErr := writeScanSummary(outputDir, scanDir, scannedItems,
			totalFiles, movieCount, tvCount, skipped); summaryErr != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Could not write summary.json: %v\n", summaryErr)
		} else {
			fmt.Printf("\n📋 Summary saved: %s\n", filepath.Join(outputDir, "summary.json"))
		}
	}
}

// videoFile holds a discovered video file's display name and full path.
type videoFile struct {
	Name     string // display name used for cleaning (dir name or filename)
	FullPath string // absolute path to the actual video file
}

// collectVideoFiles finds video files in the given directory.
// When recursive is true, it walks subdirectories up to maxDepth levels (0 = unlimited).
func collectVideoFiles(scanDir string, recursive bool, maxDepth int) []videoFile {
	var files []videoFile

	if recursive {
		scanDir = filepath.Clean(scanDir)
		baseParts := len(splitPath(scanDir))

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
				// Enforce depth limit
				if maxDepth > 0 {
					dirParts := len(splitPath(filepath.Clean(path)))
					if dirParts-baseParts > maxDepth {
						return filepath.SkipDir
					}
				}
				return nil
			}
			// Check depth for files too
			if maxDepth > 0 {
				fileParts := len(splitPath(filepath.Clean(filepath.Dir(path))))
				if fileParts-baseParts > maxDepth {
					return nil
				}
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

// splitPath splits a filepath into its components.
func splitPath(p string) []string {
	var parts []string
	for p != "" && p != "." && p != "/" && p != string(filepath.Separator) {
		dir, file := filepath.Split(p)
		if file != "" {
			parts = append(parts, file)
		}
		p = filepath.Clean(dir)
		if p == dir {
			break
		}
	}
	return parts
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
