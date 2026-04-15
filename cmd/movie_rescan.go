// movie_rescan.go — movie rescan — re-fetches TMDb data for entries with missing metadata
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/movie-cli-v3/db"
	"github.com/alimtvnetwork/movie-cli-v3/errlog"
	"github.com/alimtvnetwork/movie-cli-v3/tmdb"
)

// regenerateReports rebuilds HTML report and summary.json for every scan
// directory that has media in the DB.
func regenerateReports(database *db.DB) {
	allMedia, err := database.ListAllMedia()
	if err != nil {
		errlog.Warn("Could not list media for report regeneration: %v", err)
		return
	}
	if len(allMedia) == 0 {
		return
	}

	// Group media by scan directory (parent of .movie-output)
	dirMap := make(map[string][]db.Media)
	for _, m := range allMedia {
		if m.OriginalFilePath == "" {
			continue
		}
		scanDir := filepath.Dir(m.OriginalFilePath)
		dirMap[scanDir] = append(dirMap[scanDir], m)
	}

	for scanDir, items := range dirMap {
		outputDir := filepath.Join(scanDir, ".movie-output")
		if _, statErr := os.Stat(outputDir); os.IsNotExist(statErr) {
			continue // no .movie-output folder — skip
		}

		movieCount, tvCount := 0, 0
		for _, m := range items {
			if m.Type == "movie" {
				movieCount++
			} else {
				tvCount++
			}
		}

		if summaryErr := writeScanSummary(outputDir, scanDir, items,
			len(items), movieCount, tvCount, 0); summaryErr != nil {
			errlog.Warn("Could not regenerate summary.json for %s: %v", scanDir, summaryErr)
		}
		if htmlErr := writeHTMLReport(outputDir, scanDir, items,
			len(items), movieCount, tvCount, 0); htmlErr != nil {
			errlog.Warn("Could not regenerate report.html for %s: %v", scanDir, htmlErr)
		} else {
			fmt.Printf("🌐 Regenerated report.html → %s\n", filepath.Join(outputDir, "report.html"))
		}
	}
}

var rescanAll bool
var rescanLimit int

var movieRescanCmd = &cobra.Command{
	Use:   "rescan",
	Short: "Re-fetch TMDb metadata for entries with missing data",
	Long: `Scans the database for media entries that have missing genre, rating,
or description, and re-fetches their metadata from TMDb.

This is useful after fixing API keys or when earlier scans failed to
retrieve complete metadata. No folder scan is needed.

Examples:
  movie rescan              Re-fetch only entries with missing data
  movie rescan --all        Re-fetch TMDb data for ALL entries
  movie rescan --limit 50   Process at most 50 entries`,
	Run: runMovieRescan,
}

func init() {
	movieRescanCmd.Flags().BoolVar(&rescanAll, "all", false,
		"re-fetch TMDb data for all entries, not just those with missing data")
	movieRescanCmd.Flags().IntVar(&rescanLimit, "limit", 0,
		"max number of entries to process (0 = unlimited)")
}

func runMovieRescan(cmd *cobra.Command, args []string) {
	database, err := db.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Database error: %v\n", err)
		return
	}
	defer database.Close()

	creds := resolveScanTMDbCredentials(database)
	if !creds.HasAuth() {
		fmt.Fprintln(os.Stderr, "❌ No TMDb credentials configured. Run: movie config set tmdb_api_key YOUR_KEY")
		return
	}

	// Initialize error logger
	if initErr := errlog.Init("", "rescan"); initErr != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Could not init error logger: %v\n", initErr)
	} else {
		defer errlog.Close()
	}

	// Fetch entries to rescan
	var entries []db.Media
	if rescanAll {
		entries, err = database.ListAllMedia()
	} else {
		entries, err = database.GetMediaWithMissingData()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Database query error: %v\n", err)
		return
	}

	if rescanLimit > 0 && len(entries) > rescanLimit {
		entries = entries[:rescanLimit]
	}

	if len(entries) == 0 {
		fmt.Println("✅ All entries have complete metadata. Nothing to rescan.")
		return
	}

	client := tmdb.NewClientWithToken(creds.APIKey, creds.Token)

	fmt.Printf("\n🔄 Rescanning %d entries for TMDb metadata...\n\n", len(entries))

	updated, failed := 0, 0
	for i, m := range entries {
		fmt.Printf("  %d/%d  %s", i+1, len(entries), m.CleanTitle)
		if m.Year > 0 {
			fmt.Printf(" (%d)", m.Year)
		}

		if rescanMediaEntry(database, client, &m) {
			fmt.Printf("  ✅ ⭐%.1f %s\n", m.TmdbRating, m.Genre)
			updated++
		} else {
			fmt.Printf("  ❌ failed\n")
			failed++
		}
	}

	fmt.Printf("\n📊 Rescan Complete!\n")
	fmt.Printf("   Updated:  %d\n", updated)
	fmt.Printf("   Failed:   %d\n", failed)
	fmt.Printf("   Total:    %d\n\n", len(entries))

	// Regenerate HTML reports for affected scan directories
	if updated > 0 {
		regenerateReports(database)
	}
}
