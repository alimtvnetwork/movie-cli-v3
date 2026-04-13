// movie_scan.go — movie scan [folder] — command definition and orchestrator
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
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
  movie scan                      Scan current directory (top-level)
  movie scan ~/Movies             Scan specific folder
  movie scan -r                   Scan current directory recursively
  movie scan ~/Movies --recursive
  movie scan -r --depth 2         Scan only 2 levels deep
  movie scan --dry-run            Preview files without writing to DB
  movie scan --format table       Show results as a formatted table
  movie scan --format json        Output results as JSON to stdout`,
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

	if !useJSON {
		printScanHeader(scanDir, outputDir)
	}

	var totalFiles, movieCount, tvCount, skipped int
	var scannedItems []db.Media
	var jsonItems []scanJSONItem

	videoFiles := collectVideoFiles(scanDir, scanRecursive, scanDepth)
	useTable := scanFormat == "table"
	useTMDb := creds.HasAuth()

	if useTable {
		printScanTableHeader()
	}

	if scanDryRun {
		if useJSON {
			items, mc, tc := buildDryRunJSONItems(videoFiles)
			jsonItems = items
			totalFiles = len(items)
			movieCount = mc
			tvCount = tc
		} else if useTable {
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
		client := tmdb.NewClientWithToken(creds.APIKey, creds.Token)
		for _, vf := range videoFiles {
			processVideoFile(vf, database, client, useTMDb, outputDir,
				&totalFiles, &movieCount, &tvCount, &skipped, &scannedItems,
				useTable || useJSON)
		}
		if useJSON {
			for i := range scannedItems {
				jsonItems = append(jsonItems, buildMediaJSONItem(&scannedItems[i], "new"))
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
		printScanFooter(scanDir, outputDir, scannedItems, totalFiles, movieCount, tvCount, skipped)
	}
}

// resolveScanDir determines and validates the scan directory from args.
func resolveScanDir(args []string, quiet bool) (string, error) {
	var scanDir string
	var err error

	if len(args) > 0 {
		scanDir = args[0]
	} else {
		scanDir, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("cannot determine current directory: %v", err)
		}
		if !quiet {
			fmt.Printf("📂 No folder specified — scanning current directory\n\n")
		}
	}

	if strings.HasPrefix(scanDir, "~") {
		home, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return "", fmt.Errorf("cannot determine home directory: %v", homeErr)
		}
		scanDir = filepath.Join(home, scanDir[1:])
	}

	info, statErr := os.Stat(scanDir)
	if statErr != nil || !info.IsDir() {
		return "", fmt.Errorf("folder not found: %s", scanDir)
	}

	return scanDir, nil
}

// createOutputDirs creates the .movie-output directory structure.
func createOutputDirs(outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("cannot create output directory: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(outputDir, "json", "movie"), 0755); err != nil {
		return fmt.Errorf("cannot create json/movie dir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(outputDir, "json", "tv"), 0755); err != nil {
		return fmt.Errorf("cannot create json/tv dir: %v", err)
	}
	return nil
}

// printScanHeader prints the scan mode banner.
func printScanHeader(scanDir, outputDir string) {
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
}

// printScanFooter prints the summary after scanning completes.
func printScanFooter(scanDir, outputDir string, scannedItems []db.Media,
	totalFiles, movieCount, tvCount, skipped int) {
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
		if summaryErr := writeScanSummary(outputDir, scanDir, scannedItems,
			totalFiles, movieCount, tvCount, skipped); summaryErr != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Could not write summary.json: %v\n", summaryErr)
		} else {
			fmt.Printf("\n📋 Summary saved: %s\n", filepath.Join(outputDir, "summary.json"))
		}
	}
}
