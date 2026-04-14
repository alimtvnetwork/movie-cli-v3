// movie_ls_detail.go — detail view and helpers for movie ls
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/movie-cli-v3/db"
)

func showMediaDetail(database *db.DB, id int64) {
	m, err := database.GetMediaByID(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ❌ Not found: %v\n", err)
		return
	}

	fmt.Print("\033[H\033[2J")
	printMediaDetail(m)
}

func printMediaDetail(m *db.Media) {
	typeIcon := "🎬 Movie"
	if m.Type == "tv" {
		typeIcon = "📺 TV Show"
	}

	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Printf("║  %s\n", m.Title)
	if m.Tagline != "" {
		fmt.Printf("║  \"%s\"\n", m.Tagline)
	}
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()

	if m.Year > 0 {
		fmt.Printf("  📅 Year:        %d\n", m.Year)
	}
	fmt.Printf("  🏷️  Type:        %s\n", typeIcon)

	if m.Runtime > 0 {
		fmt.Printf("  ⏱️  Runtime:     %d min\n", m.Runtime)
	}
	if m.Language != "" {
		fmt.Printf("  🌐 Language:    %s\n", strings.ToUpper(m.Language))
	}

	if m.ImdbRating > 0 {
		fmt.Printf("  ⭐ IMDb:        %.1f\n", m.ImdbRating)
	}
	if m.TmdbRating > 0 {
		fmt.Printf("  ⭐ TMDb:        %.1f\n", m.TmdbRating)
	}
	if m.Popularity > 0 {
		fmt.Printf("  📈 Popularity:  %.0f\n", m.Popularity)
	}

	if m.Genre != "" {
		fmt.Printf("  🎭 Genre:       %s\n", m.Genre)
	}
	if m.Director != "" {
		fmt.Printf("  🎬 Director:    %s\n", m.Director)
	}
	if m.CastList != "" {
		fmt.Printf("  👥 Cast:        %s\n", m.CastList)
	}

	if m.Budget > 0 {
		fmt.Printf("  💰 Budget:      $%s\n", formatMoney(m.Budget))
	}
	if m.Revenue > 0 {
		fmt.Printf("  💵 Revenue:     $%s\n", formatMoney(m.Revenue))
	}

	if m.Description != "" {
		fmt.Println()
		fmt.Printf("  📝 %s\n", m.Description)
	}

	if m.TrailerURL != "" {
		fmt.Println()
		fmt.Printf("  🎥 Trailer:     %s\n", m.TrailerURL)
	}

	if m.ThumbnailPath != "" {
		fmt.Printf("  🖼️  Thumbnail:   %s\n", m.ThumbnailPath)
	}

	if m.CurrentFilePath != "" {
		fmt.Println()
		fmt.Printf("  📁 File:        %s\n", m.CurrentFilePath)
	}
}

// formatMoney formats an int64 as a human-readable money string (e.g. 1,234,567).
func formatMoney(n int64) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}
	var result []byte
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}
	return string(result)
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] >= 'a' && s[0] <= 'z' {
		return string(s[0]-32) + s[1:]
	}
	return s
}
