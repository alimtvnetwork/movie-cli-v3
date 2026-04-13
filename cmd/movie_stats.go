// movie_stats.go — movie stats
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/movie-cli-v3/db"
)

var statsFormat string

var movieStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show library statistics",
	Long: `Display total counts, top genres, and average ratings.

Use --format json to output stats as JSON to stdout for piping.`,
	Run: runMovieStats,
}

func init() {
	movieStatsCmd.Flags().StringVar(&statsFormat, "format", "default",
		"output format: default or json")
}

// statsJSONOutput is the JSON structure for --format json.
type statsJSONOutput struct {
	TotalMovies  int              `json:"total_movies"`
	TotalTV      int              `json:"total_tv_shows"`
	Total        int              `json:"total"`
	Storage      *statsStorage    `json:"storage,omitempty"`
	TopGenres    []statsGenre     `json:"top_genres,omitempty"`
	AvgImdb      float64          `json:"avg_imdb_rating,omitempty"`
	AvgTmdb      float64          `json:"avg_tmdb_rating,omitempty"`
}

type statsStorage struct {
	TotalSize    int64  `json:"total_bytes"`
	TotalHuman   string `json:"total_human"`
	LargestFile  int64  `json:"largest_file_bytes"`
	SmallestFile int64  `json:"smallest_file_bytes"`
	AverageSize  int64  `json:"average_bytes"`
}

type statsGenre struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func runMovieStats(cmd *cobra.Command, args []string) {
	database, err := db.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Database error: %v\n", err)
		return
	}
	defer database.Close()

	totalMovies, _ := database.CountMedia("movie")
	totalTV, _ := database.CountMedia("tv")
	total, err := database.CountMedia("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Database error: %v\n", err)
		return
	}

	if total == 0 {
		if statsFormat == "json" {
			fmt.Println("{}")
		} else {
			fmt.Println("📭 No media in library. Run 'movie scan <folder>' first.")
		}
		return
	}

	if statsFormat == "json" {
		printStatsJSON(database, totalMovies, totalTV, total)
	} else {
		printStatsDefault(database, totalMovies, totalTV, total)
	}
}

func printStatsJSON(database *db.DB, totalMovies, totalTV, total int) {
	out := statsJSONOutput{
		TotalMovies: totalMovies,
		TotalTV:     totalTV,
		Total:       total,
	}

	// Storage
	totalSize, largestSize, smallestSize, sizeErr := database.FileSizeStats()
	if sizeErr == nil && totalSize > 0 {
		out.Storage = &statsStorage{
			TotalSize:    totalSize,
			TotalHuman:   humanSize(totalSize),
			LargestFile:  largestSize,
			SmallestFile: smallestSize,
			AverageSize:  totalSize / int64(total),
		}
	}

	// Genres
	genres, genreErr := database.TopGenres(10)
	if genreErr == nil && len(genres) > 0 {
		for n, c := range genres {
			out.TopGenres = append(out.TopGenres, statsGenre{Name: n, Count: c})
		}
		sort.Slice(out.TopGenres, func(i, j int) bool {
			return out.TopGenres[i].Count > out.TopGenres[j].Count
		})
	}

	// Ratings
	allMedia, listErr := database.ListMedia(0, 100000)
	if listErr == nil {
		var sumImdb, sumTmdb float64
		var cntImdb, cntTmdb int
		for i := range allMedia {
			if allMedia[i].ImdbRating > 0 {
				sumImdb += allMedia[i].ImdbRating
				cntImdb++
			}
			if allMedia[i].TmdbRating > 0 {
				sumTmdb += allMedia[i].TmdbRating
				cntTmdb++
			}
		}
		if cntImdb > 0 {
			out.AvgImdb = sumImdb / float64(cntImdb)
		}
		if cntTmdb > 0 {
			out.AvgTmdb = sumTmdb / float64(cntTmdb)
		}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if encErr := enc.Encode(out); encErr != nil {
		fmt.Fprintf(os.Stderr, "❌ JSON encode error: %v\n", encErr)
	}
}

func printStatsDefault(database *db.DB, totalMovies, totalTV, total int) {
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
