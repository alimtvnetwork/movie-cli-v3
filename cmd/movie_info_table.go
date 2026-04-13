// movie_info_table.go — table-formatted output for movie info
package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/movie-cli-v3/db"
)

// printMediaDetailTable outputs a media item as a formatted key-value table.
func printMediaDetailTable(m *db.Media) {
	labelWidth := 12
	valueWidth := 55

	fmt.Println()
	fmt.Printf("  %-*s │ %-*s\n", labelWidth, "Field", valueWidth, "Value")
	fmt.Printf("  %s─┼─%s\n",
		strings.Repeat("─", labelWidth),
		strings.Repeat("─", valueWidth))

	rows := []struct {
		label string
		value string
	}{
		{"Title", m.Title},
		{"Year", fmt.Sprintf("%d", m.Year)},
		{"Type", m.Type},
	}

	if m.TmdbID > 0 {
		rows = append(rows, struct{ label, value string }{"TMDb ID", fmt.Sprintf("%d", m.TmdbID)})
	}
	if m.ImdbID != "" {
		rows = append(rows, struct{ label, value string }{"IMDb ID", m.ImdbID})
	}

	rating := "N/A"
	if m.TmdbRating > 0 {
		rating = fmt.Sprintf("%.1f", m.TmdbRating)
	}
	rows = append(rows, struct{ label, value string }{"Rating", rating})

	if m.Genre != "" {
		rows = append(rows, struct{ label, value string }{"Genre", m.Genre})
	}
	if m.Director != "" {
		rows = append(rows, struct{ label, value string }{"Director", truncate(m.Director, valueWidth)})
	}
	if m.CastList != "" {
		rows = append(rows, struct{ label, value string }{"Cast", truncate(m.CastList, valueWidth)})
	}
	if m.Runtime > 0 {
		rows = append(rows, struct{ label, value string }{"Runtime", fmt.Sprintf("%d min", m.Runtime)})
	}
	if m.Language != "" {
		rows = append(rows, struct{ label, value string }{"Language", m.Language})
	}
	if m.Tagline != "" {
		rows = append(rows, struct{ label, value string }{"Tagline", truncate(m.Tagline, valueWidth)})
	}
	if m.TrailerURL != "" {
		rows = append(rows, struct{ label, value string }{"Trailer", m.TrailerURL})
	}
	if m.Budget > 0 {
		rows = append(rows, struct{ label, value string }{"Budget", fmt.Sprintf("$%d", m.Budget)})
	}
	if m.Revenue > 0 {
		rows = append(rows, struct{ label, value string }{"Revenue", fmt.Sprintf("$%d", m.Revenue)})
	}
	if m.CurrentFilePath != "" {
		rows = append(rows, struct{ label, value string }{"File", truncate(m.CurrentFilePath, valueWidth)})
	}
	if m.Description != "" {
		rows = append(rows, struct{ label, value string }{"Description", truncate(m.Description, valueWidth)})
	}

	for _, r := range rows {
		fmt.Printf("  %-*s │ %-*s\n", labelWidth, r.label, valueWidth, r.value)
	}

	fmt.Printf("  %s─┴─%s\n",
		strings.Repeat("─", labelWidth),
		strings.Repeat("─", valueWidth))
	fmt.Println()
}
