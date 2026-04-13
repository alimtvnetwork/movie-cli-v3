// movie_ls.go — movie ls
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/movie-cli-v3/db"
)

var movieLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List scanned movies and TV shows from your library",
	Long: `Lists scan-indexed movies and TV shows (items with a known file path).
Only items added via 'movie scan' are shown.
Press N for next page, P for previous, Q to quit.`,
	Run: runMovieLs,
}

func runMovieLs(cmd *cobra.Command, args []string) {
	database, err := db.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Database error: %v\n", err)
		return
	}
	defer database.Close()

	pageSizeStr, cfgErr := database.GetConfig("page_size")
	if cfgErr != nil && cfgErr.Error() != "sql: no rows in result set" {
		fmt.Fprintf(os.Stderr, "⚠️  Config read error (page_size): %v\n", cfgErr)
	}
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize <= 0 {
		pageSize = 20
	}

	total, countErr := database.CountMedia("")
	if countErr != nil {
		fmt.Fprintf(os.Stderr, "❌ Database error: %v\n", countErr)
		return
	}
	if total == 0 {
		fmt.Println("📭 No media found. Run 'movie scan <folder>' first.")
		return
	}

	offset := 0
	scanner := bufio.NewScanner(os.Stdin)

	for {
		media, listErr := database.ListMedia(offset, pageSize)
		if listErr != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", listErr)
			return
		}

		// Clear screen
		fmt.Print("\033[H\033[2J")

		page := (offset / pageSize) + 1
		totalPages := (total + pageSize - 1) / pageSize

		fmt.Printf("🎬 Your Library — Page %d/%d (%d total)\n", page, totalPages, total)
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		for i := range media {
			num := offset + i + 1
			yearStr := ""
			if media[i].Year > 0 {
				yearStr = fmt.Sprintf("(%d)", media[i].Year)
			}

			rating := "N/A"
			if media[i].TmdbRating > 0 {
				rating = fmt.Sprintf("%.1f", media[i].TmdbRating)
			} else if media[i].ImdbRating > 0 {
				rating = fmt.Sprintf("%.1f", media[i].ImdbRating)
			}

			typeIcon := "🎬"
			if media[i].Type == "tv" {
				typeIcon = "📺"
			}

			fmt.Printf("  %3d. %-40s %-6s  ⭐ %-4s  %s %s\n",
				num, media[i].CleanTitle, yearStr, rating, typeIcon, capitalize(media[i].Type))
		}

		fmt.Println()
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Print("  [N] Next  [P] Previous  [Q] Quit  [1-9] View details → ")

		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		switch {
		case input == "n" || input == "N":
			if offset+pageSize < total {
				offset += pageSize
			} else {
				fmt.Println("  ⚠️  Already on last page")
			}
		case input == "p" || input == "P":
			if offset-pageSize >= 0 {
				offset -= pageSize
			} else {
				fmt.Println("  ⚠️  Already on first page")
			}
		case input == "q" || input == "Q":
			fmt.Println("👋 Bye!")
			return
		default:
			// Try to parse as number for detail view
			if num, parseErr := strconv.Atoi(input); parseErr == nil && num > 0 && num <= total {
				showMediaDetail(database, int64(num))
				fmt.Print("\nPress Enter to continue...")
				scanner.Scan()
			}
		}
	}
}

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
