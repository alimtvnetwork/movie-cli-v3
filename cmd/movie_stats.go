// movie_stats.go — movie stats
package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/movie-cli-v3/db"
)

var movieStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show library statistics",
	Long:  `Display total counts, top genres, and average ratings.`,
	Run:   runMovieStats,
}

func runMovieStats(cmd *cobra.Command, args []string) {
	database, err := db.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Database error: %v\n", err)
		return
	}
	defer database.Close()

	totalMovies, err := database.CountMedia("movie")
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Count movies error: %v\n", err)
	}
	totalTV, err := database.CountMedia("tv")
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Count TV error: %v\n", err)
	}
	total, err := database.CountMedia("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Database error: %v\n", err)
		return
	}

	if total == 0 {
		fmt.Println("📭 No media in library. Run 'movie scan <folder>' first.")
		return
	}

	fmt.Println("📊 Library Statistics")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("  🎬 Total Movies:    %d\n", totalMovies)
	fmt.Printf("  📺 Total TV Shows:  %d\n", totalTV)
	fmt.Printf("  📁 Total:           %d\n", total)
	fmt.Println()

	// File size stats
	totalSize, largestSize, smallestSize, sizeErr := database.FileSizeStats()
	if sizeErr != nil {
		fmt.Fprintf(os.Stderr, "⚠️  File size stats error: %v\n", sizeErr)
	} else if totalSize > 0 {
		fmt.Println("  💾 Storage:")
		fmt.Printf("     Total Size:    %s\n", humanSize(totalSize))
		fmt.Printf("     Largest File:  %s\n", humanSize(largestSize))
		fmt.Printf("     Smallest File: %s\n", humanSize(smallestSize))
		if total > 0 {
			fmt.Printf("     Average Size:  %s\n", humanSize(totalSize/int64(total)))
		}
		fmt.Println()
	}

	// Top genres
	genres, genreErr := database.TopGenres(10)
	if genreErr != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Genre stats error: %v\n", genreErr)
	} else if len(genres) > 0 {
		type gc struct {
			name  string
			count int
		}
		var sorted []gc
		for n, c := range genres {
			sorted = append(sorted, gc{n, c})
		}
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].count > sorted[j].count
		})

		fmt.Println("  🎭 Top Genres:")
		for i, g := range sorted {
			if i >= 10 {
				break
			}
			bar := ""
			for j := 0; j < g.count && j < 30; j++ {
				bar += "█"
			}
			fmt.Printf("     %-20s %s %d\n", g.name, bar, g.count)
		}
		fmt.Println()
	}

	// Average ratings
	var avgImdb, avgTmdb float64
	var imdbCount, tmdbCount int

	allMedia, listErr := database.ListMedia(0, 10000)
	if listErr != nil {
		fmt.Fprintf(os.Stderr, "⚠️  List media error: %v\n", listErr)
	}
	for i := range allMedia {
		if allMedia[i].ImdbRating > 0 {
			avgImdb += allMedia[i].ImdbRating
			imdbCount++
		}
		if allMedia[i].TmdbRating > 0 {
			avgTmdb += allMedia[i].TmdbRating
			tmdbCount++
		}
	}

	if imdbCount > 0 {
		fmt.Printf("  ⭐ Avg IMDb Rating: %.1f\n", avgImdb/float64(imdbCount))
	}
	if tmdbCount > 0 {
		fmt.Printf("  ⭐ Avg TMDb Rating: %.1f\n", avgTmdb/float64(tmdbCount))
	}
}
