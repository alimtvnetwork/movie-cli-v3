// movie_scan_watch.go — file watcher for movie scan --watch
package cmd

import (
	"fmt"
	"time"

	"github.com/alimtvnetwork/movie-cli-v3/db"
	"github.com/alimtvnetwork/movie-cli-v3/errlog"
	"github.com/alimtvnetwork/movie-cli-v3/tmdb"
)

var scanWatch bool
var scanWatchInterval int

// runWatchLoop polls the scan directory for new video files at a fixed interval.
// It keeps a set of already-seen file paths and processes only new ones.
func runWatchLoop(scanDir, outputDir string, database *db.DB, creds tmdbCredentials) {
	seen := make(map[string]bool)

	// Seed with existing files so we don't re-process them
	initial := collectVideoFiles(scanDir, scanRecursive, scanDepth)
	for _, vf := range initial {
		seen[vf.FullPath] = true
	}

	interval := time.Duration(scanWatchInterval) * time.Second
	client := tmdb.NewClientWithToken(creds.APIKey, creds.Token)
	useTMDb := creds.HasAuth()

	fmt.Printf("\n  👁️  Watching for new files (every %ds) — press Ctrl+C to stop\n", scanWatchInterval)
	fmt.Println("  ──────────────────────────────────────────")

	cycle := 0
	for {
		time.Sleep(interval)
		cycle++

		current := collectVideoFiles(scanDir, scanRecursive, scanDepth)
		var newFiles []videoFile
		for _, vf := range current {
			if !seen[vf.FullPath] {
				newFiles = append(newFiles, vf)
				seen[vf.FullPath] = true
			}
		}

		if len(newFiles) == 0 {
			continue
		}

		fmt.Printf("\n  🔔 Detected %d new file(s) at %s\n",
			len(newFiles), time.Now().Format("15:04:05"))

		var totalFiles, movieCount, tvCount, skipped int
		var scannedItems []db.Media

		for _, vf := range newFiles {
			processVideoFile(vf, database, client, useTMDb, outputDir,
				&totalFiles, &movieCount, &tvCount, &skipped, &scannedItems, false, "")
		}

		// Log watch cycle
		if !scanDryRun {
			if histErr := database.InsertScanHistory(scanDir, totalFiles, movieCount, tvCount); histErr != nil {
				errlog.Warn("Could not log watch scan history: %v", histErr)
			}
		}

		fmt.Printf("  ✅ Processed: %d files (%d movies, %d TV)\n", totalFiles, movieCount, tvCount)
	}
}

