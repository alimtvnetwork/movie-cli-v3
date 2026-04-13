// movie_search_table.go — table-formatted output for movie search
package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/movie-cli-v3/tmdb"
)

// printSearchResultsTable outputs TMDb search results as a formatted table.
func printSearchResultsTable(results []tmdb.MultiResult) {
	fmt.Println()
	fmt.Printf("  %-3s │ %-35s │ %-6s │ %-8s │ %-6s │ %-6s\n",
		"#", "Title", "Year", "Type", "Rating", "TMDb ID")
	fmt.Printf("  %s─┼─%s─┼─%s─┼─%s─┼─%s─┼─%s\n",
		strings.Repeat("─", 3),
		strings.Repeat("─", 35),
		strings.Repeat("─", 6),
		strings.Repeat("─", 8),
		strings.Repeat("─", 6),
		strings.Repeat("─", 6))

	for i, r := range results {
		if i >= 15 {
			break
		}
		title := truncate(r.GetDisplayTitle(), 35)
		year := r.GetYear()
		if year == "" {
			year = "  -   "
		}
		mediaType := "Movie"
		if r.MediaType == "tv" {
			mediaType = "TV Show"
		}
		rating := "   -  "
		if r.VoteAvg > 0 {
			rating = fmt.Sprintf("%5.1f ", r.VoteAvg)
		}
		fmt.Printf("  %-3d │ %-35s │ %-6s │ %-8s │ %s│ %6d\n",
			i+1, title, year, mediaType, rating, r.ID)
	}

	fmt.Printf("  %s─┴─%s─┴─%s─┴─%s─┴─%s─┴─%s\n",
		strings.Repeat("─", 3),
		strings.Repeat("─", 35),
		strings.Repeat("─", 6),
		strings.Repeat("─", 8),
		strings.Repeat("─", 6),
		strings.Repeat("─", 6))
	fmt.Println()
}
