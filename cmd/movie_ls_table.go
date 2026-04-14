// movie_ls_table.go — table-formatted output for movie ls
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/movie-cli-v3/db"
)

// runMovieLsTable outputs all library items as a formatted table (no pager).
func runMovieLsTable(database *db.DB) {
	allMedia, err := database.ListMedia(0, 100000)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Database error: %v\n", err)
		return
	}

	if len(allMedia) == 0 {
		fmt.Println("📭 No media found. Run 'movie scan <folder>' first.")
		return
	}

	numW := 5
	titleW := 40
	yearW := 6
	typeW := 8
	ratingW := 6
	genreW := 25
	directorW := 20

	fmt.Println()
	fmt.Printf("  %-*s │ %-*s │ %-*s │ %-*s │ %-*s │ %-*s │ %-*s\n",
		numW, "#",
		titleW, "Title",
		yearW, "Year",
		typeW, "Type",
		ratingW, "Rating",
		genreW, "Genre",
		directorW, "Director")

	fmt.Printf("  %s─┼─%s─┼─%s─┼─%s─┼─%s─┼─%s─┼─%s\n",
		strings.Repeat("─", numW),
		strings.Repeat("─", titleW),
		strings.Repeat("─", yearW),
		strings.Repeat("─", typeW),
		strings.Repeat("─", ratingW),
		strings.Repeat("─", genreW),
		strings.Repeat("─", directorW))

	for i := range allMedia {
		m := &allMedia[i]

		title := truncate(m.CleanTitle, titleW)
		yearStr := ""
		if m.Year > 0 {
			yearStr = fmt.Sprintf("%d", m.Year)
		}

		typeLabel := capitalize(m.Type)

		rating := "N/A"
		if m.TmdbRating > 0 {
			rating = fmt.Sprintf("%.1f", m.TmdbRating)
		} else if m.ImdbRating > 0 {
			rating = fmt.Sprintf("%.1f", m.ImdbRating)
		}

		genre := truncate(m.Genre, genreW)
		director := truncate(m.Director, directorW)

		fmt.Printf("  %-*d │ %-*s │ %-*s │ %-*s │ %-*s │ %-*s │ %-*s\n",
			numW, i+1,
			titleW, title,
			yearW, yearStr,
			typeW, typeLabel,
			ratingW, rating,
			genreW, genre,
			directorW, director)
	}

	fmt.Printf("  %s─┴─%s─┴─%s─┴─%s─┴─%s─┴─%s─┴─%s\n",
		strings.Repeat("─", numW),
		strings.Repeat("─", titleW),
		strings.Repeat("─", yearW),
		strings.Repeat("─", typeW),
		strings.Repeat("─", ratingW),
		strings.Repeat("─", genreW),
		strings.Repeat("─", directorW))

	fmt.Printf("\n  Total: %d items\n\n", len(allMedia))
}

// truncate shortens a string to maxLen, adding "…" if needed.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 1 {
		return "…"
	}
	return s[:maxLen-1] + "…"
}
