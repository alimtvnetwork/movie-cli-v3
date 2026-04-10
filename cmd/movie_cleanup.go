// movie_cleanup.go — find and remove stale DB entries where files no longer exist.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mahin/mahin-cli-v2/db"
)

var cleanupDryRun bool

var movieCleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Find stale entries where files no longer exist",
	Long: `Scan the library for media entries whose file path no longer exists
on disk. By default, shows a preview (dry run). Use --remove to delete
stale entries from the database.

Examples:
  movie cleanup              # preview stale entries (dry run)
  movie cleanup --remove     # delete stale entries from DB`,
	Run: runMovieCleanup,
}

func init() {
	movieCleanupCmd.Flags().BoolVar(&cleanupDryRun, "remove", false,
		"actually delete stale entries (default is dry-run preview)")
}

func runMovieCleanup(cmd *cobra.Command, args []string) {
	database, err := db.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Database error: %v\n", err)
		return
	}
	defer database.Close()

	stale, err := database.FindStaleEntries(10000)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error scanning for stale entries: %v\n", err)
		return
	}

	if len(stale) == 0 {
		fmt.Println("✅ No stale entries found — all files exist on disk.")
		return
	}

	fmt.Printf("🔍 Found %d stale entries (file missing from disk):\n", len(stale))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for i, s := range stale {
		fmt.Printf("  %d. [ID %d] %s (%d)\n", i+1, s.Media.ID, s.Media.Title, s.Media.Year)
		fmt.Printf("     Missing: %s\n", s.FilePath)
	}

	if !cleanupDryRun {
		fmt.Printf("\n📋 Dry run — no changes made. Use --remove to delete these entries.\n")
		return
	}

	// Confirm before deleting
	fmt.Printf("\n⚠️  Delete %d stale entries from the database? [y/N] ", len(stale))
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer != "y" && answer != "yes" {
		fmt.Println("❌ Aborted.")
		return
	}

	deleted := 0
	for _, s := range stale {
		if err := database.DeleteMedia(s.Media.ID); err != nil {
			fmt.Fprintf(os.Stderr, "  ⚠️  Failed to delete ID %d: %v\n", s.Media.ID, err)
			continue
		}
		deleted++
	}

	fmt.Printf("✅ Deleted %d stale entries from the database.\n", deleted)
}
