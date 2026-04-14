// movie_db.go — movie db: print resolved database path for debugging
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/movie-cli-v3/db"
)

var movieDBCmd = &cobra.Command{
	Use:   "db",
	Short: "Show the resolved database path and status",
	Long:  `Prints the full resolved path to the SQLite database and data directory. Useful for debugging data location issues.`,
	Run:   runMovieDB,
}

func runMovieDB(cmd *cobra.Command, args []string) {
	database, err := db.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Database error: %v\n", err)
		return
	}
	defer database.Close()

	fmt.Println()
	fmt.Println("📂 Movie CLI — Data Location")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("  Data directory:  %s\n", database.BasePath)
	fmt.Printf("  Database file:   %s/movie.db\n", database.BasePath)
	fmt.Printf("  Thumbnails:      %s/thumbnails/\n", database.BasePath)
	fmt.Printf("  JSON metadata:   %s/json/\n", database.BasePath)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Show record counts
	total, countErr := database.CountMedia("")
	if countErr == nil {
		movies, _ := database.CountMedia("movie")
		tv, _ := database.CountMedia("tv")
		fmt.Printf("\n  Records: %d total (%d movies, %d TV shows)\n", total, movies, tv)
	}

	fmt.Println()
}
