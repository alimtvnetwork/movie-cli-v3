// movie_move_helpers.go — shared helpers for move/rename/undo operations
//
// -- Shared helpers exported from this file --
//
//	expandHome(path, home)                     — resolve ~ in paths
//	listVideoFiles(dir) ([]FileInfo, error)    — list video files in a directory
//	humanSize(bytes) string                    — format bytes as human-readable size
//	promptSourceDirectory(scanner, db, home)   — interactive source dir picker
//	promptDestination(scanner, db, home)       — interactive destination picker
//	MoveFile(src, dst) error                   — move with cross-device fallback
//	crossDeviceMove(src, dst) error            — copy+delete for cross-filesystem moves
//	saveHistoryLog(basePath, title, year, from, to) — write move-log.json
//
// Consumers: movie_move.go, movie_rename.go, movie_undo.go, movie_stats.go
//
// Do NOT duplicate move/size/path logic elsewhere — use these helpers.
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/alimtvnetwork/movie-cli-v3/cleaner"
	"github.com/alimtvnetwork/movie-cli-v3/db"
)

// expandHome replaces ~ with actual home directory.
func expandHome(path, home string) string {
	if strings.HasPrefix(path, "~") {
		return filepath.Join(home, path[1:])
	}
	return path
}

// listVideoFiles returns all video files in a directory.
// Returns an error if the directory cannot be read.
func listVideoFiles(dir string) ([]os.FileInfo, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot read directory %s: %w", dir, err)
	}

	var files []os.FileInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !cleaner.IsVideoFile(entry.Name()) {
			continue
		}
		info, infoErr := entry.Info()
		if infoErr != nil {
			fmt.Fprintf(os.Stderr, "  ⚠️  Cannot stat %s: %v\n", entry.Name(), infoErr)
			continue
		}
		files = append(files, info)
	}
	return files, nil
}

// humanSize formats bytes into human-readable form.
func humanSize(bytes int64) string {
	const (
		kb = 1024
		mb = kb * 1024
		gb = mb * 1024
	)
	switch {
	case bytes >= gb:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(gb))
	case bytes >= mb:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(mb))
	case bytes >= kb:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(kb))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// promptSourceDirectory asks the user to pick a directory.
func promptSourceDirectory(scanner interface {
	Scan() bool
	Text() string
}, database *db.DB, home string) string {
	scanDir, cfgErr := database.GetConfig("scan_dir")
	if cfgErr != nil && cfgErr.Error() != "sql: no rows in result set" {
		fmt.Fprintf(os.Stderr, "⚠️  Config read error (scan_dir): %v\n", cfgErr)
	}
	scanDir = expandHome(scanDir, home)

	fmt.Println("📂 Where are your video files?")
	fmt.Println()

	options := []string{}
	labels := []string{}

	if scanDir != "" {
		options = append(options, scanDir)
		labels = append(labels, fmt.Sprintf("Scan folder (%s)", scanDir))
	}

	downloads := filepath.Join(home, "Downloads")
	if info, err := os.Stat(downloads); err == nil && info.IsDir() {
		options = append(options, downloads)
		labels = append(labels, fmt.Sprintf("Downloads (%s)", downloads))
	}

	desktop := filepath.Join(home, "Desktop")
	if info, err := os.Stat(desktop); err == nil && info.IsDir() {
		options = append(options, desktop)
		labels = append(labels, fmt.Sprintf("Desktop (%s)", desktop))
	}

	for i, label := range labels {
		fmt.Printf("  %d. %s\n", i+1, label)
	}
	fmt.Printf("  %d. Enter custom path\n", len(labels)+1)
	fmt.Println()
	fmt.Print("  Choose: ")

	if !scanner.Scan() {
		return ""
	}
	choice, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil || choice < 1 || choice > len(options)+1 {
		fmt.Fprintln(os.Stderr, "❌ Invalid choice")
		return ""
	}

	if choice <= len(options) {
		return options[choice-1]
	}

	fmt.Print("  Enter path: ")
	if !scanner.Scan() {
		return ""
	}
	return expandHome(strings.TrimSpace(scanner.Text()), home)
}

// promptDestination asks the user to choose a move destination.
func promptDestination(scanner interface {
	Scan() bool
	Text() string
}, database *db.DB, home string) string {
	moviesDir, cfgErr := database.GetConfig("movies_dir")
	if cfgErr != nil && cfgErr.Error() != "sql: no rows in result set" {
		fmt.Fprintf(os.Stderr, "⚠️  Config read error (movies_dir): %v\n", cfgErr)
	}
	tvDir, cfgErr := database.GetConfig("tv_dir")
	if cfgErr != nil && cfgErr.Error() != "sql: no rows in result set" {
		fmt.Fprintf(os.Stderr, "⚠️  Config read error (tv_dir): %v\n", cfgErr)
	}
	archiveDir, cfgErr := database.GetConfig("archive_dir")
	if cfgErr != nil && cfgErr.Error() != "sql: no rows in result set" {
		fmt.Fprintf(os.Stderr, "⚠️  Config read error (archive_dir): %v\n", cfgErr)
	}
	moviesDir = expandHome(moviesDir, home)
	tvDir = expandHome(tvDir, home)
	archiveDir = expandHome(archiveDir, home)

	if moviesDir == "" {
		moviesDir = expandHome("~/Movies", home)
	}
	if tvDir == "" {
		tvDir = expandHome("~/TVShows", home)
	}
	if archiveDir == "" {
		archiveDir = expandHome("~/Archive", home)
	}

	fmt.Println()
	fmt.Println("  📁 Move to:")
	fmt.Printf("  1. 🎬 Movies (%s)\n", moviesDir)
	fmt.Printf("  2. 📺 TV Shows (%s)\n", tvDir)
	fmt.Printf("  3. 📦 Archive (%s)\n", archiveDir)
	fmt.Println("  4. 📂 Custom path")
	fmt.Println()
	fmt.Print("  Choose [1-4]: ")

	if !scanner.Scan() {
		return ""
	}
	choice := strings.TrimSpace(scanner.Text())

	switch choice {
	case "1":
		return moviesDir
	case "2":
		return tvDir
	case "3":
		return archiveDir
	case "4":
		fmt.Print("  Enter path: ")
		if !scanner.Scan() {
			return ""
		}
		return expandHome(strings.TrimSpace(scanner.Text()), home)
	default:
		fmt.Fprintln(os.Stderr, "❌ Invalid choice")
		return ""
	}
}

// MoveFile moves a file from src to dst using os.Rename with cross-device fallback.
func MoveFile(src, dst string) error {
	err := os.Rename(src, dst)
	if err != nil && isCrossDeviceError(err) {
		return crossDeviceMove(src, dst)
	}
	return err
}

// isCrossDeviceError checks whether the error is an EXDEV (cross-device link)
// error, which occurs when os.Rename is called across different filesystems
// (e.g., USB drives, network mounts, different partitions).
func isCrossDeviceError(err error) bool {
	var linkErr *os.LinkError
	if errors.As(err, &linkErr) {
		var errno syscall.Errno
		if errors.As(linkErr.Err, &errno) {
			return errno == syscall.EXDEV
		}
	}
	return false
}

// crossDeviceMove copies the file from src to dst, preserves the original file
// permissions, and removes the source only after the destination is fully
// written and synced. This is the fallback when os.Rename fails with EXDEV.
func crossDeviceMove(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("stat source: %w", err)
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("create destination: %w", err)
	}

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		dstFile.Close()
		os.Remove(dst)
		return fmt.Errorf("copy data: %w", err)
	}

	if err := dstFile.Sync(); err != nil {
		dstFile.Close()
		os.Remove(dst)
		return fmt.Errorf("sync destination: %w", err)
	}
	dstFile.Close()

	return os.Remove(src)
}

// saveHistoryLog writes a JSON move record to the history log.
// All errors are logged to stderr — never swallowed.
func saveHistoryLog(basePath, title string, year int, fromPath, toPath string) {
	historyDir := filepath.Join(basePath, "json", "history")
	if err := os.MkdirAll(historyDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠️  Cannot create history dir: %v\n", err)
		return
	}

	record := map[string]interface{}{
		"title":     title,
		"year":      year,
		"from_path": fromPath,
		"to_path":   toPath,
		"moved_at":  time.Now().UTC().Format(time.RFC3339),
	}

	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠️  Cannot marshal history JSON: %v\n", err)
		return
	}

	filename := fmt.Sprintf("move-%s.json", time.Now().UTC().Format("20060102-150405"))
	historyPath := filepath.Join(historyDir, filename)
	if writeErr := os.WriteFile(historyPath, data, 0644); writeErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠️  Cannot write history file: %v\n", writeErr)
	}
}
