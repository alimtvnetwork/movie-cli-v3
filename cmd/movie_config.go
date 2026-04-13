// movie_config.go — movie config
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/movie-cli-v3/db"
)

var movieConfigCmd = &cobra.Command{
	Use:   "config [get|set] [key] [value]",
	Short: "Manage movie CLI configuration",
	Long: `View or update configuration settings.

Keys:
  movies_dir     - Default movies directory
  tv_dir         - Default TV shows directory
  archive_dir    - Default archive directory
  scan_dir       - Default scan directory
  tmdb_api_key   - TMDb API key
  tmdb_token     - TMDb access token
  page_size      - Items per page in list view

Examples:
  movie config                           # Show all
  movie config get movies_dir            # Get one
  movie config set movies_dir ~/Movies   # Set one
  movie config set tmdb_api_key abc123   # Set API key
  movie config set tmdb_token eyJ...     # Set access token`,
	Run: runMovieConfig,
}

func runMovieConfig(cmd *cobra.Command, args []string) {
	database, err := db.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Database error: %v\n", err)
		return
	}
	defer database.Close()

	if len(args) == 0 {
		showAllConfig(database)
		return
	}

	action := args[0]

	switch action {
	case "get":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "❌ Usage: movie config get <key>")
			return
		}
		val, getErr := database.GetConfig(args[1])
		if getErr != nil {
			fmt.Printf("  %s = (not set)\n", args[1])
		} else {
			fmt.Printf("  %s = %s\n", args[1], val)
		}

	case "set":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "❌ Usage: movie config set <key> <value>")
			return
		}
		key, value := args[1], args[2]
		if setErr := database.SetConfig(key, value); setErr != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", setErr)
			return
		}
		fmt.Printf("  ✅ %s = %s\n", key, value)

	default:
		fmt.Fprintf(os.Stderr, "❌ Unknown action: %s. Use 'get' or 'set'.\n", action)
	}
}

func showAllConfig(database *db.DB) {
	fmt.Println("⚙️  Configuration:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	keys := []string{"movies_dir", "tv_dir", "archive_dir", "scan_dir", "tmdb_api_key", "tmdb_token", "page_size"}
	for _, key := range keys {
		val, err := database.GetConfig(key)
		if err != nil {
			val = "(not set)"
		}
		if (key == "tmdb_api_key" || key == "tmdb_token") && len(val) > 8 {
			val = val[:4] + "..." + val[len(val)-4:]
		}
		fmt.Printf("  %-15s = %s\n", key, val)
	}
}
