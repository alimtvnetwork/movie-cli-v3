// movie_move.go — movie move
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/movie-cli-v3/cleaner"
	"github.com/alimtvnetwork/movie-cli-v3/db"
)

var moveAllFlag bool

var movieMoveCmd = &cobra.Command{
	Use:   "move [directory]",
	Short: "Browse a local directory and move a movie/TV show file",
	Long: `Browse a local directory for video files, select one, and move it
to a configured destination (Movies, TV Shows, Archive, or custom path).
The move is logged for undo support.

Use --all to move all video files at once. Movies go to movies_dir,
TV shows go to tv_dir (auto-detected from filename).

If no directory is given, you'll be prompted to choose one.`,
	Args: cobra.MaximumNArgs(1),
	Run:  runMovieMove,
}

func init() {
	movieMoveCmd.Flags().BoolVar(&moveAllFlag, "all", false, "Move all video files in the directory at once")
}

func runMovieMove(cmd *cobra.Command, args []string) {
	database, err := db.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Database error: %v\n", err)
		return
	}
	defer database.Close()

	scanner := bufio.NewScanner(os.Stdin)
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		fmt.Fprintf(os.Stderr, "❌ Cannot determine home directory: %v\n", homeErr)
		return
	}

	// Step 1: Determine the source directory
	sourceDir := ""
	if len(args) > 0 {
		sourceDir = expandHome(args[0], home)
	} else {
		sourceDir = promptSourceDirectory(scanner, database, home)
		if sourceDir == "" {
			return
		}
	}

	// Validate directory
	info, statErr := os.Stat(sourceDir)
	if statErr != nil {
		fmt.Fprintf(os.Stderr, "❌ Cannot access directory: %v\n", statErr)
		return
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "❌ Path is not a directory: %s\n", sourceDir)
		return
	}

	// Step 2: List video files in the directory
	files, listErr := listVideoFiles(sourceDir)
	if listErr != nil {
		fmt.Fprintf(os.Stderr, "❌ %v\n", listErr)
		return
	}
	if len(files) == 0 {
		fmt.Printf("📭 No video files found in: %s\n", sourceDir)
		return
	}

	if moveAllFlag {
		runBatchMove(database, scanner, sourceDir, files, home)
		return
	}

	runInteractiveMove(database, scanner, sourceDir, files, home)
}

// runBatchMove moves all video files at once, auto-routing by type.
func runBatchMove(database *db.DB, scanner *bufio.Scanner, sourceDir string, files []os.FileInfo, home string) {
	moviesDir, cfgErr := database.GetConfig("movies_dir")
	if cfgErr != nil && cfgErr.Error() != "sql: no rows in result set" {
		fmt.Fprintf(os.Stderr, "⚠️  Config read error (movies_dir): %v\n", cfgErr)
	}
	tvDir, cfgErr := database.GetConfig("tv_dir")
	if cfgErr != nil && cfgErr.Error() != "sql: no rows in result set" {
		fmt.Fprintf(os.Stderr, "⚠️  Config read error (tv_dir): %v\n", cfgErr)
	}
	moviesDir = expandHome(moviesDir, home)
	tvDir = expandHome(tvDir, home)

	if moviesDir == "" {
		moviesDir = expandHome("~/Movies", home)
	}
	if tvDir == "" {
		tvDir = expandHome("~/TVShows", home)
	}

	// Preview all moves
	type moveItem struct {
		srcPath   string
		destPath  string
		destDir   string
		cleanName string
		result    cleaner.Result
		fileInfo  os.FileInfo
	}

	var moves []moveItem

	fmt.Printf("\n🎬 Batch move — %d video files in: %s\n\n", len(files), sourceDir)
	fmt.Println("  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for _, f := range files {
		result := cleaner.Clean(f.Name())
		cleanName := cleaner.ToCleanFileName(result.CleanTitle, result.Year, result.Extension)

		destDir := moviesDir
		typeIcon := "🎬"
		if result.Type == "tv" {
			destDir = tvDir
			typeIcon = "📺"
		}

		srcPath := filepath.Join(sourceDir, f.Name())
		destPath := filepath.Join(destDir, cleanName)

		yearStr := ""
		if result.Year > 0 {
			yearStr = fmt.Sprintf(" (%d)", result.Year)
		}

		fmt.Printf("  %s %s%s  [%s]\n", typeIcon, result.CleanTitle, yearStr, humanSize(f.Size()))
		fmt.Printf("     → %s\n", destPath)

		moves = append(moves, moveItem{
			srcPath:   srcPath,
			destPath:  destPath,
			destDir:   destDir,
			cleanName: cleanName,
			result:    result,
			fileInfo:  f,
		})
	}

	fmt.Println("  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("\n  Move all %d files? [y/N]: ", len(moves))

	if !scanner.Scan() {
		return
	}
	confirm := strings.ToLower(strings.TrimSpace(scanner.Text()))
	if confirm != "y" && confirm != "yes" {
		fmt.Println("  ❌ Batch move canceled.")
		return
	}

	// Execute all moves
	success := 0
	failed := 0

	for i := range moves {
		// Create destination directory
		if mkdirErr := os.MkdirAll(moves[i].destDir, 0755); mkdirErr != nil {
			fmt.Fprintf(os.Stderr, "  ❌ Cannot create dir %s: %v\n", moves[i].destDir, mkdirErr)
			failed++
			continue
		}

		// Move file
		if moveErr := MoveFile(moves[i].srcPath, moves[i].destPath); moveErr != nil {
			fmt.Fprintf(os.Stderr, "  ❌ Failed: %s — %v\n", moves[i].fileInfo.Name(), moveErr)
			failed++
			continue
		}

		// Track in DB
		trackMove(database, moves[i].result, moves[i].fileInfo, moves[i].srcPath, moves[i].destPath, moves[i].cleanName)
		success++
	}

	fmt.Println()
	if failed == 0 {
		fmt.Printf("  ✅ All %d files moved successfully!\n", success)
	} else {
		fmt.Printf("  ⚠️  %d moved, %d failed\n", success, failed)
	}
}

// runInteractiveMove is the original single-file interactive flow.
func runInteractiveMove(database *db.DB, scanner *bufio.Scanner, sourceDir string, files []os.FileInfo, home string) {
	fmt.Printf("\n🎬 Video files in: %s\n\n", sourceDir)
	for i, f := range files {
		result := cleaner.Clean(f.Name())
		typeIcon := "🎬"
		if result.Type == "tv" {
			typeIcon = "📺"
		}
		yearStr := ""
		if result.Year > 0 {
			yearStr = fmt.Sprintf("(%d)", result.Year)
		}
		fmt.Printf("  %2d. %s %s %s  [%s]\n", i+1, typeIcon, result.CleanTitle, yearStr, humanSize(f.Size()))
	}

	// Select a file
	fmt.Println()
	fmt.Print("  Select file [number]: ")
	if !scanner.Scan() {
		return
	}
	choice, parseErr := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if parseErr != nil || choice < 1 || choice > len(files) {
		fmt.Fprintln(os.Stderr, "❌ Invalid selection")
		return
	}

	selectedFile := files[choice-1]
	selectedPath := filepath.Join(sourceDir, selectedFile.Name())
	result := cleaner.Clean(selectedFile.Name())

	fmt.Printf("\n  Selected: %s\n", result.CleanTitle)
	if result.Year > 0 {
		fmt.Printf("  Year:     %d\n", result.Year)
	}
	fmt.Printf("  Type:     %s\n", result.Type)

	// Choose destination
	destDir := promptDestination(scanner, database, home)
	if destDir == "" {
		return
	}

	// Build clean filename and confirm
	cleanName := cleaner.ToCleanFileName(result.CleanTitle, result.Year, result.Extension)
	destPath := filepath.Join(destDir, cleanName)

	fmt.Println()
	fmt.Println("  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("  📄 From: %s\n", selectedPath)
	fmt.Printf("  📁 To:   %s\n", destPath)
	fmt.Println("  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Print("  Are you sure? [y/N]: ")

	if !scanner.Scan() {
		return
	}
	confirm := strings.ToLower(strings.TrimSpace(scanner.Text()))
	if confirm != "y" && confirm != "yes" {
		fmt.Println("  ❌ Move canceled.")
		return
	}

	// Create destination directory
	if mkdirErr := os.MkdirAll(destDir, 0755); mkdirErr != nil {
		fmt.Fprintf(os.Stderr, "  ❌ Cannot create directory: %v\n", mkdirErr)
		return
	}

	// Move the file
	if moveErr := MoveFile(selectedPath, destPath); moveErr != nil {
		fmt.Fprintf(os.Stderr, "  ❌ Move failed: %v\n", moveErr)
		return
	}

	// Track history
	trackMove(database, result, selectedFile, selectedPath, destPath, cleanName)

	fmt.Println()
	fmt.Println("  ✅ Moved successfully!")
	fmt.Printf("     %s\n", selectedPath)
	fmt.Printf("     → %s\n", destPath)
}

// trackMove records a move in the database and JSON history log.
func trackMove(database *db.DB, result cleaner.Result, fileInfo os.FileInfo, srcPath, destPath, cleanName string) {
	var mediaID int64
	existing, searchErr := database.SearchMedia(result.CleanTitle)
	if searchErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠️  DB search error: %v\n", searchErr)
	}
	for i := range existing {
		if existing[i].CurrentFilePath == srcPath || existing[i].OriginalFilePath == srcPath {
			mediaID = existing[i].ID
			break
		}
	}

	if mediaID == 0 {
		m := &db.Media{
			Title:            result.CleanTitle,
			CleanTitle:       result.CleanTitle,
			Year:             result.Year,
			Type:             result.Type,
			OriginalFileName: fileInfo.Name(),
			OriginalFilePath: srcPath,
			CurrentFilePath:  destPath,
			FileExtension:    result.Extension,
			FileSize:         fileInfo.Size(),
		}
		var insertErr error
		mediaID, insertErr = database.InsertMedia(m)
		if insertErr != nil {
			fmt.Fprintf(os.Stderr, "  ❌ DB insert error: %v\n", insertErr)
		}
	} else {
		if updateErr := database.UpdateMediaPath(mediaID, destPath); updateErr != nil {
			fmt.Fprintf(os.Stderr, "  ❌ DB update path error: %v\n", updateErr)
		}
	}

	if mediaID > 0 {
		if histErr := database.InsertMoveHistory(mediaID, srcPath, destPath,
			fileInfo.Name(), cleanName); histErr != nil {
			fmt.Fprintf(os.Stderr, "  ⚠️  DB history error: %v\n", histErr)
		}
	}

	saveHistoryLog(database.BasePath, result.CleanTitle, result.Year, srcPath, destPath)
}
