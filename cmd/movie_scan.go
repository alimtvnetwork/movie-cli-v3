// movie_scan.go — movie scan [folder] — command definition and orchestrator
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/movie-cli-v3/cleaner"
	"github.com/alimtvnetwork/movie-cli-v3/db"
	"github.com/alimtvnetwork/movie-cli-v3/errlog"
	"github.com/alimtvnetwork/movie-cli-v3/tmdb"
)

var scanRecursive bool
var scanDepth int
var scanDryRun bool
var scanFormat string
var scanRest bool
var scanRestPort int

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
  movie scan                      Scan current directory (top-level)
  movie scan ~/Movies             Scan specific folder
  movie scan -r                   Scan current directory recursively
  movie scan ~/Movies --recursive
  movie scan -r --depth 2         Scan only 2 levels deep
  movie scan --dry-run            Preview files without writing to DB
  movie scan --format table       Show results as a formatted table
  movie scan --format json        Output results as JSON to stdout
  movie scan --rest               Scan and start REST server + open browser
  movie scan --rest --port 9000   Scan and start REST on custom port
  movie scan --watch              Scan then watch for new files
  movie scan --watch --interval 5 Watch with 5-second polling`,
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
		"output format: default, table, or json")
	movieScanCmd.Flags().BoolVar(&scanRest, "rest", false,
		"start REST server and open HTML report in browser after scan")
	movieScanCmd.Flags().IntVar(&scanRestPort, "port", 8086,
		"port for REST server when using --rest")
	movieScanCmd.Flags().BoolVarP(&scanWatch, "watch", "w", false,
		"watch for new files after initial scan")
	movieScanCmd.Flags().IntVar(&scanWatchInterval, "interval", 10,
		"polling interval in seconds for --watch mode")
}

func runMovieScan(cmd *cobra.Command, args []string) {
	useJSON := scanFormat == "json"

	database, err := db.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Database error: %v\n", err)
		return
	}
	defer database.Close()

	scanDir, err := resolveScanDir(args, useJSON)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ %v\n", err)
		return
	}

	creds := resolveScanTMDbCredentials(database)
	outputDir := filepath.Join(scanDir, ".movie-output")

	if !scanDryRun {
		if err := createOutputDirs(outputDir); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			return
		}
	}

	// Initialize error logger — writes to .movie-output/logs/error.txt + DB
	if initErr := errlog.Init(outputDir, "scan"); initErr != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Could not init error logger: %v\n", initErr)
	} else {
		defer errlog.Close()
		// Wire DB writer
		errlog.SetDBWriter(func(e errlog.Entry) {
			if dbErr := database.InsertErrorLog(
				e.Timestamp, string(e.Level), e.Source, e.Function,
				e.Command, e.WorkDir, e.Message, e.StackTrace,
			); dbErr != nil {
				fmt.Fprintf(os.Stderr, "⚠️  Could not write error to DB: %v\n", dbErr)
			}
		})
	}

	if !useJSON {
		printScanHeader(scanDir, outputDir)
	}

	var totalFiles, movieCount, tvCount, skipped, removed int
	var scannedItems []db.Media
	var jsonItems []scanJSONItem

	videoFiles := collectVideoFiles(scanDir, scanRecursive, scanDepth)
	useTable := scanFormat == "table"
	useTMDb := creds.HasAuth()

	// Generate batch ID for action_history tracking
	scanBatchID := generateBatchID()

	if useTable {
		printScanTableHeader()
	}

	if scanDryRun {
		runDryRunScan(videoFiles, useJSON, useTable,
			&jsonItems, &totalFiles, &movieCount, &tvCount)
	} else {
		// ── Re-scan: detect removed files and clean up ──
		existingMedia, _ := database.GetMediaByScanDir(scanDir)
		diskPaths := make(map[string]bool, len(videoFiles))
		for _, vf := range videoFiles {
			diskPaths[vf.FullPath] = true
		}
		var removeIDs []int64
		var removeMedia []*db.Media
		for i := range existingMedia {
			if !diskPaths[existingMedia[i].OriginalFilePath] {
				removeIDs = append(removeIDs, existingMedia[i].ID)
				removeMedia = append(removeMedia, &existingMedia[i])
			}
		}
		if len(removeIDs) > 0 {
			// Snapshot each removed entry before deletion for undo support
			for _, rm := range removeMedia {
				snapshot, snapErr := db.MediaToJSON(rm)
				if snapErr != nil {
					errlog.Warn("Could not snapshot media %d for undo: %v", rm.ID, snapErr)
					continue
				}
				detail := fmt.Sprintf("Scan removed: %s (%s)", rm.CleanTitle, rm.OriginalFilePath)
				database.InsertActionSimple(db.ActionScanRemove, rm.ID, snapshot, detail, scanBatchID)
			}

			delCount, delErr := database.DeleteMediaByIDs(removeIDs)
			if delErr != nil {
				errlog.Warn("Could not remove %d stale entries: %v", len(removeIDs), delErr)
			} else {
				removed = delCount
				if !useJSON && !useTable {
					fmt.Printf("  🗑️  Removed %d entries (files no longer on disk)\n\n", removed)
				}
			}
		}

		// ── Build set of existing DB paths for skip detection ──
		existingPaths := make(map[string]*db.Media, len(existingMedia))
		for i := range existingMedia {
			existingPaths[existingMedia[i].OriginalFilePath] = &existingMedia[i]
		}

		// ── Process files: skip existing, enrich new ──
		client := tmdb.NewClientWithToken(creds.APIKey, creds.Token)
		for _, vf := range videoFiles {
		if em, found := existingPaths[vf.FullPath]; found {
				// Already in DB — auto-rescan if missing metadata
				totalFiles++
				status := "existing"
				if useTMDb && mediaNeedsRescan(em) {
					// Snapshot before rescan for undo
					preSnapshot, _ := db.MediaToJSON(em)
					if rescanMediaEntry(database, client, em) {
						status = "rescanned"
						detail := fmt.Sprintf("Rescan updated: %s", em.CleanTitle)
						database.InsertActionSimple(db.ActionRescanUpdate, em.ID, preSnapshot, detail, scanBatchID)
						if !useTable && !useJSON {
							typeIcon := "🎬"
							if em.Type == "tv" {
								typeIcon = "📺"
							}
							fmt.Printf("\n  %d. %s %s", totalFiles, typeIcon, em.CleanTitle)
							if em.Year > 0 {
								fmt.Printf(" (%d)", em.Year)
							}
							fmt.Printf(" [%s]\n", em.Type)
							fmt.Printf("     🔄 Rescanned — ⭐%.1f %s\n", em.TmdbRating, em.Genre)
						}
					} else {
						skipped++
						if !useTable && !useJSON {
							fmt.Printf("\n  %d. %s", totalFiles, em.CleanTitle)
							fmt.Printf(" [%s]\n", em.Type)
							fmt.Println("     ⚠️  Rescan failed — kept existing data")
						}
					}
				} else {
					skipped++
					if useTable {
						printScanTableRow(buildMediaTableRow(totalFiles, em, "existing"))
					} else if !useJSON {
						typeIcon := "🎬"
						if em.Type == "tv" {
							typeIcon = "📺"
						}
						fmt.Printf("\n  %d. %s %s", totalFiles, typeIcon, em.CleanTitle)
						if em.Year > 0 {
							fmt.Printf(" (%d)", em.Year)
						}
						fmt.Printf(" [%s]\n", em.Type)
						fmt.Println("     ⏩ Already in database")
					}
				}
				scannedItems = append(scannedItems, *em)
				if em.Type == "movie" {
					movieCount++
				} else {
					tvCount++
				}
				if useTable && status != "existing" {
					printScanTableRow(buildMediaTableRow(totalFiles, em, status))
				}
				continue
			}
			processVideoFile(vf, database, client, useTMDb, outputDir,
				&totalFiles, &movieCount, &tvCount, &skipped, &scannedItems,
				useTable || useJSON, scanBatchID)
		}
		if useJSON {
			for i := range scannedItems {
				status := "existing"
				if existingPaths[scannedItems[i].OriginalFilePath] == nil {
					status = "new"
				}
				jsonItems = append(jsonItems, buildMediaJSONItem(&scannedItems[i], status))
			}
		}
	}

	if useTable {
		printScanTableFooter()
	}

	if !scanDryRun {
		if histErr := database.InsertScanHistory(scanDir, totalFiles, movieCount, tvCount); histErr != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Could not log scan history: %v\n", histErr)
		}
	}

	if useJSON {
		printScanJSON(scanDir, jsonItems, totalFiles, movieCount, tvCount, skipped)
	} else {
		printScanFooter(scanDir, outputDir, scannedItems, totalFiles, movieCount, tvCount, skipped, removed)
	}

	// Start REST server if --rest was specified
	if scanRest && !scanDryRun {
		restPort = scanRestPort
		fmt.Printf("\n🚀 Starting REST server on http://localhost:%d ...\n", restPort)
		go openBrowser(fmt.Sprintf("http://localhost:%d", restPort))
		if scanWatch {
			// Run watch in background, REST in foreground
			go runWatchLoop(scanDir, outputDir, database, creds)
		}
		runMovieRest(cmd, []string{})
		return
	}

	// Start watch mode if --watch was specified (without --rest)
	if scanWatch && !scanDryRun {
		runWatchLoop(scanDir, outputDir, database, creds)
	}
}

// runDryRunScan handles the dry-run scanning loop for all output formats.
func runDryRunScan(videoFiles []videoFile, useJSON, useTable bool,
	jsonItems *[]scanJSONItem, totalFiles, movieCount, tvCount *int) {
	if useJSON {
		items, mc, tc := buildDryRunJSONItems(videoFiles)
		*jsonItems = items
		*totalFiles = len(items)
		*movieCount = mc
		*tvCount = tc
	} else if useTable {
		rows, mc, tc := buildDryRunTableRows(videoFiles)
		for _, row := range rows {
			printScanTableRow(row)
		}
		*totalFiles = len(rows)
		*movieCount = mc
		*tvCount = tc
	} else {
		for _, vf := range videoFiles {
			*totalFiles++
			result := cleaner.Clean(vf.Name)
			typeIcon := "🎬"
			if result.Type == "tv" {
				typeIcon = "📺"
			}
			fmt.Printf("\n  %d. %s %s", *totalFiles, typeIcon, result.CleanTitle)
			if result.Year > 0 {
				fmt.Printf(" (%d)", result.Year)
			}
			fmt.Printf(" [%s]\n", result.Type)
			fmt.Printf("     └─ %s\n", vf.Name)
			if result.Type == "movie" {
				*movieCount++
			} else {
				*tvCount++
			}
		}
	}
}
