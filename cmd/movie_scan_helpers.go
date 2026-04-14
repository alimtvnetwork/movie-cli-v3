// movie_scan_helpers.go — shared helpers for movie scan (dir resolution, output dirs, print)
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/movie-cli-v3/db"
	"github.com/alimtvnetwork/movie-cli-v3/errlog"
)

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

// printScanHeader prints the scan mode banner (gitmap-style box).
func printScanHeader(scanDir, outputDir string) {
	fmt.Println()
	fmt.Println("  ╔══════════════════════════════════════╗")
	fmt.Println("  ║         🎬  Movie CLI Scanner        ║")
	fmt.Println("  ╚══════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("  📂 Scanning: %s\n", scanDir)
	if scanDryRun {
		fmt.Println("  🧪 Mode: dry run (no writes)")
	}
	if scanRecursive {
		if scanDepth > 0 {
			fmt.Printf("  🔄 Mode: recursive (max depth: %d)\n", scanDepth)
		} else {
			fmt.Println("  🔄 Mode: recursive (all subdirectories)")
		}
	}
	if !scanDryRun {
		fmt.Printf("  📁 Output: %s\n", outputDir)
	}
	fmt.Println()
	fmt.Println("  ■ Scanned Items")
	fmt.Println("  ──────────────────────────────────────────")
}

// printScanFooter prints the summary after scanning completes (gitmap-style).
func printScanFooter(scanDir, outputDir string, scannedItems []db.Media,
	totalFiles, movieCount, tvCount, skipped int) {
	fmt.Println()
	fmt.Println("  ■ Summary")
	fmt.Println("  ──────────────────────────────────────────")
	if scanDryRun {
		fmt.Println("  📊 Dry Run Complete!")
	} else {
		fmt.Println("  📊 Scan Complete!")
	}
	fmt.Printf("     Total files: %d\n", totalFiles)
	fmt.Printf("     Movies:      %d\n", movieCount)
	fmt.Printf("     TV Shows:    %d\n", tvCount)
	if skipped > 0 {
		fmt.Printf("     Skipped:     %d (already in DB)\n", skipped)
	}
	if scanDryRun {
		fmt.Println("\n  💡 Run without --dry-run to actually scan and save.")
	} else {
		fmt.Println()
		fmt.Println("  ■ Output Files")
		fmt.Println("  ──────────────────────────────────────────")
		fmt.Printf("  📁 %s/\n", outputDir)
		if summaryErr := writeScanSummary(outputDir, scanDir, scannedItems,
			totalFiles, movieCount, tvCount, skipped); summaryErr != nil {
			errlog.Warn("Could not write summary.json: %v", summaryErr)
		} else {
			fmt.Printf("  ├── 📄 summary.json      Scan report with metadata\n")
		}
		if htmlErr := writeHTMLReport(outputDir, scanDir, scannedItems,
			totalFiles, movieCount, tvCount, skipped); htmlErr != nil {
			errlog.Warn("Could not write report.html: %v", htmlErr)
		} else {
			fmt.Printf("  ├── 🌐 report.html       Interactive HTML report\n")
		}
		fmt.Printf("  ├── 📁 json/movie/       Per-movie JSON metadata\n")
		fmt.Printf("  ├── 📁 json/tv/          Per-show JSON metadata\n")
		fmt.Printf("  └── 📁 thumbnails/       Movie poster thumbnails\n")
	}
	fmt.Println()
}
