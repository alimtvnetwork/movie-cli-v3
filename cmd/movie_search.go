// movie_search.go — movie search <name>
// Searches TMDb API, fetches full details, and saves to local database.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/movie-cli-v3/cleaner"
	"github.com/alimtvnetwork/movie-cli-v3/db"
	"github.com/alimtvnetwork/movie-cli-v3/tmdb"
)

var searchFormat string

var movieSearchCmd = &cobra.Command{
	Use:   "search [name]",
	Short: "Search TMDb for a movie or TV show and save to database",
	Long: `Searches the TMDb API for movies/TV shows matching the query.
Fetches full metadata (rating, genres, cast, crew, poster) and saves
to the local database. Categorizes as Movie or TV Show automatically.
Does NOT require the file to exist in your library.

Use --format json to output search results as JSON (no interactive prompt).
Use --format table to output search results as a formatted table (no interactive prompt).`,
	Args: cobra.MinimumNArgs(1),
	Run:  runMovieSearch,
}

func init() {
	movieSearchCmd.Flags().StringVar(&searchFormat, "format", "", "Output format: json, table")
}



func runMovieSearch(cmd *cobra.Command, args []string) {
	database, err := db.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Database error: %v\n", err)
		return
	}
	defer database.Close()

	// Get TMDb API key (GetConfig returns "" with nil error when key is absent)
	apiKey, cfgErr := database.GetConfig("tmdb_api_key")
	if cfgErr != nil && cfgErr.Error() != "sql: no rows in result set" {
		fmt.Fprintf(os.Stderr, "⚠️  Config read error: %v\n", cfgErr)
	}
	if apiKey == "" {
		apiKey = os.Getenv("TMDB_API_KEY")
	}
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "❌ No TMDb API key configured.")
		fmt.Fprintln(os.Stderr, "   Set it with: movie config set tmdb_api_key YOUR_KEY")
		return
	}

	client := tmdb.NewClient(apiKey)
	query := strings.Join(args, " ")
	if searchFormat != "json" && searchFormat != "table" {
		fmt.Printf("🔎 Searching TMDb for: %s\n\n", query)
	}

	// Search TMDb API
	results, searchErr := client.SearchMulti(query)
	if searchErr != nil {
		fmt.Fprintf(os.Stderr, "❌ TMDb search error: %v\n", searchErr)
		return
	}

	if len(results) == 0 {
		if searchFormat == "json" {
			fmt.Println("[]")
		} else {
			fmt.Println("📭 No results found on TMDb.")
		}
		return
	}

	// JSON mode: output results and exit (no interactive prompt)
	if searchFormat == "json" {
		printSearchResultsJSON(results)
		return
	}

	// Table mode: output results and exit (no interactive prompt)
	if searchFormat == "table" {
		printSearchResultsTable(results)
		return
	}

	// Show results and let user pick
	fmt.Printf("Found %d results:\n\n", len(results))
	for i := range results {
		if i >= 15 {
			break
		}
		title := results[i].GetDisplayTitle()
		year := results[i].GetYear()
		typeIcon := "🎬"
		typeLabel := "Movie"
		if results[i].MediaType == "tv" {
			typeIcon = "📺"
			typeLabel = "TV Show"
		}

		rating := "N/A"
		if results[i].VoteAvg > 0 {
			rating = fmt.Sprintf("%.1f", results[i].VoteAvg)
		}

		yearStr := ""
		if year != "" {
			yearStr = fmt.Sprintf("(%s)", year)
		}

		fmt.Printf("  %d. %s %-35s %-6s  ⭐ %-4s  [%s]\n",
			i+1, typeIcon, title, yearStr, rating, typeLabel)
	}

	fmt.Println()
	fmt.Print("Enter number to save (0 to cancel): ")

	var choice int
	_, scanErr := fmt.Scan(&choice)
	if scanErr != nil || choice < 1 || choice > len(results) || choice > 15 {
		fmt.Println("❌ Canceled.")
		return
	}

	selected := results[choice-1]
	saveSearchResult(client, database, selected)
}
