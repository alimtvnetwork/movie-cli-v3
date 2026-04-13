// movie_rename.go — movie rename
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/movie-cli-v3/cleaner"
	"github.com/alimtvnetwork/movie-cli-v3/db"
)

var movieRenameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename files to clean names",
	Long: `Automatically renames messy filenames to clean format.
Example: Scream.2022.1080p.WEBRip.x264-RARBG.mkv → Scream (2022).mkv`,
	Run: runMovieRename,
}

func runMovieRename(cmd *cobra.Command, args []string) {
	database, err := db.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Database error: %v\n", err)
		return
	}
	defer database.Close()

	media, listErr := database.ListMedia(0, 10000)
	if listErr != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to read media: %v\n", listErr)
		return
	}
	if len(media) == 0 {
		fmt.Println("📭 No media found.")
		return
	}

	// Find files that need renaming
	type renameItem struct {
		oldPath string
		newPath string
		oldName string
		newName string
		media   db.Media
	}

	var items []renameItem
	for i := range media {
		if media[i].CurrentFilePath == "" {
			continue
		}
		dir := filepath.Dir(media[i].CurrentFilePath)
		oldName := filepath.Base(media[i].CurrentFilePath)
		newName := cleaner.ToCleanFileName(media[i].CleanTitle, media[i].Year, media[i].FileExtension)

		if oldName != newName {
			items = append(items, renameItem{
				media:   media[i],
				oldPath: media[i].CurrentFilePath,
				newPath: filepath.Join(dir, newName),
				oldName: oldName,
				newName: newName,
			})
		}
	}

	if len(items) == 0 {
		fmt.Println("✅ All files already have clean names!")
		return
	}

	fmt.Printf("📝 Found %d files to rename:\n\n", len(items))
	for i := range items {
		fmt.Printf("  %d. %s\n", i+1, items[i].oldName)
		fmt.Printf("     → %s\n\n", items[i].newName)
	}

	fmt.Print("Rename all? [y/N]: ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return
	}
	confirm := strings.ToLower(strings.TrimSpace(scanner.Text()))
	if confirm != "y" && confirm != "yes" {
		fmt.Println("❌ Canceled.")
		return
	}

	success := 0
	for i := range items {
		if moveErr := MoveFile(items[i].oldPath, items[i].newPath); moveErr != nil {
			fmt.Fprintf(os.Stderr, "  ❌ Failed: %s → %v\n", items[i].oldName, moveErr)
			continue
		}
		if updateErr := database.UpdateMediaPath(items[i].media.ID, items[i].newPath); updateErr != nil {
			fmt.Fprintf(os.Stderr, "  ⚠️  DB update path error: %v\n", updateErr)
		}
		if histErr := database.InsertMoveHistory(items[i].media.ID, items[i].oldPath, items[i].newPath,
			items[i].oldName, items[i].newName); histErr != nil {
			fmt.Fprintf(os.Stderr, "  ⚠️  DB history error: %v\n", histErr)
		}
		fmt.Printf("  ✅ %s → %s\n", items[i].oldName, items[i].newName)
		success++
	}

	fmt.Printf("\n✅ Renamed %d/%d files.\n", success, len(items))
}
