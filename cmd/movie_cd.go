// movie_cd.go — movie cd [folder] — print scanned folder path for shell cd
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/movie-cli-v3/db"
	"github.com/alimtvnetwork/movie-cli-v3/errlog"
)

var movieCdCmd = &cobra.Command{
	Use:   "cd [folder-name]",
	Short: "Print the path of a scanned folder for quick navigation",
	Long: `Prints the full path of a previously scanned folder so you can navigate
to it in your terminal.

Usage with shell:
  cd $(movie cd Movies)          # Jump to the folder matching "Movies"
  cd $(movie cd)                 # Jump to the most recently scanned folder

Without arguments, lists all known scan folders with numbers for selection.

Examples:
  movie cd                       List all scanned folders
  movie cd Movies                Print path matching "Movies"
  movie cd 1                     Print path of folder #1 from list`,
	Args: cobra.MaximumNArgs(1),
	Run:  runMovieCd,
}

func init() {}

func runMovieCd(cmd *cobra.Command, args []string) {
	database, err := db.Open()
	if err != nil {
		errlog.Error("Database error: %v", err)
		return
	}
	defer database.Close()

	folders, err := database.ListDistinctScanFolders()
	if err != nil {
		errlog.Error("Database error: %v", err)
		return
	}

	if len(folders) == 0 {
		fmt.Fprintln(os.Stderr, "📭 No scanned folders found. Run 'movie scan <folder>' first.")
		return
	}

	// No argument: list all folders
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "📂 Scanned folders:")
		fmt.Fprintln(os.Stderr, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		for i, f := range folders {
			fmt.Fprintf(os.Stderr, "  %d. %s\n", i+1, f)
		}
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "💡 Usage: cd $(movie cd <name-or-number>)")
		return
	}

	query := args[0]

	// Try as a number first
	if num := 0; true {
		_, scanErr := fmt.Sscanf(query, "%d", &num)
		if scanErr == nil && num > 0 && num <= len(folders) {
			fmt.Print(folders[num-1])
			return
		}
	}

	// Try as a substring match (case-insensitive)
	queryLower := strings.ToLower(query)
	var matches []string
	for _, f := range folders {
		if strings.Contains(strings.ToLower(f), queryLower) {
			matches = append(matches, f)
		}
	}

	switch len(matches) {
	case 0:
		fmt.Fprintf(os.Stderr, "❌ No scanned folder matches '%s'\n", query)
		fmt.Fprintln(os.Stderr, "Run 'movie cd' to see all scanned folders.")
		os.Exit(1)
	case 1:
		fmt.Print(matches[0])
	default:
		fmt.Fprintln(os.Stderr, "⚠️  Multiple matches:")
		for i, m := range matches {
			fmt.Fprintf(os.Stderr, "  %d. %s\n", i+1, m)
		}
		fmt.Fprintln(os.Stderr, "\n💡 Be more specific or use: cd $(movie cd <number>)")
		os.Exit(1)
	}
}
