// movie_undo.go — movie undo
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/movie-cli-v3/db"
)

var movieUndoCmd = &cobra.Command{
	Use:   "undo",
	Short: "Undo the last move operation",
	Long:  `Reverts the most recent movie move operation.`,
	Run:   runMovieUndo,
}

func runMovieUndo(cmd *cobra.Command, args []string) {
	database, err := db.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Database error: %v\n", err)
		return
	}
	defer database.Close()

	lastMove, moveErr := database.GetLastMove()
	if moveErr != nil {
		fmt.Println("📭 No move operations to undo.")
		return
	}

	fmt.Println("⏪ Last move operation:")
	fmt.Println()
	fmt.Printf("  📁 %s\n", lastMove.ToPath)
	fmt.Printf("  → %s\n", lastMove.FromPath)
	fmt.Println()

	// Confirmation prompt
	fmt.Print("Undo this? [y/N]: ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return
	}
	confirm := strings.ToLower(strings.TrimSpace(scanner.Text()))
	if confirm != "y" && confirm != "yes" {
		fmt.Println("❌ Undo canceled.")
		return
	}

	// Check source exists
	if _, statErr := os.Stat(lastMove.ToPath); statErr != nil {
		if os.IsNotExist(statErr) {
			fmt.Fprintf(os.Stderr, "❌ File not found at: %s\n", lastMove.ToPath)
			fmt.Fprintln(os.Stderr, "   It may have been moved or deleted manually.")
		} else {
			fmt.Fprintf(os.Stderr, "❌ Cannot access file %s: %v\n", lastMove.ToPath, statErr)
		}
		return
	}

	// Move back
	if undoErr := MoveFile(lastMove.ToPath, lastMove.FromPath); undoErr != nil {
		fmt.Fprintf(os.Stderr, "❌ Undo failed: %v\n", undoErr)
		return
	}

	// Mark as undone in DB
	if markErr := database.MarkMoveUndone(lastMove.ID); markErr != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Could not mark move as undone: %v\n", markErr)
	}

	// Update media path back
	if pathErr := database.UpdateMediaPath(lastMove.MediaID, lastMove.FromPath); pathErr != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Could not update media path: %v\n", pathErr)
	}

	fmt.Println()
	fmt.Println("✅ Undo successful!")
	fmt.Printf("   %s\n", lastMove.ToPath)
	fmt.Printf("   → %s\n", lastMove.FromPath)
}
