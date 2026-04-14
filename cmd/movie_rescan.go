// movie_rescan.go — movie rescan — re-fetches TMDb data for entries with missing metadata
package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/movie-cli-v3/cleaner"
	"github.com/alimtvnetwork/movie-cli-v3/db"
	"github.com/alimtvnetwork/movie-cli-v3/errlog"
	"github.com/alimtvnetwork/movie-cli-v3/tmdb"
)

var rescanAll bool
var rescanLimit int

var movieRescanCmd = &cobra.Command{
	Use:   "rescan",
	Short: "Re-fetch TMDb metadata for entries with missing data",
	Long: `Scans the database for media entries that have missing genre, rating,
or description, and re-fetches their metadata from TMDb.

This is useful after fixing API keys or when earlier scans failed to
retrieve complete metadata. No folder scan is needed.

Examples:
  movie rescan              Re-fetch only entries with missing data
  movie rescan --all        Re-fetch TMDb data for ALL entries
  movie rescan --limit 50   Process at most 50 entries`,
	Run: runMovieRescan,
}

func init() {
	movieRescanCmd.Flags().BoolVar(&rescanAll, "all", false,
		"re-fetch TMDb data for all entries, not just those with missing data")
	movieRescanCmd.Flags().IntVar(&rescanLimit, "limit", 0,
		"max number of entries to process (0 = unlimited)")
}

func runMovieRescan(cmd *cobra.Command, args []string) {
	database, err := db.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Database error: %v\n", err)
		return
	}
	defer database.Close()

	creds := resolveScanTMDbCredentials(database)
	if !creds.HasAuth() {
		fmt.Fprintln(os.Stderr, "❌ No TMDb credentials configured. Run: movie config set tmdb_api_key YOUR_KEY")
		return
	}

	// Initialize error logger
	if initErr := errlog.Init("", "rescan"); initErr != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Could not init error logger: %v\n", initErr)
	} else {
		defer errlog.Close()
	}

	// Fetch entries to rescan
	var entries []db.Media
	if rescanAll {
		entries, err = database.ListAllMedia()
	} else {
		entries, err = database.GetMediaWithMissingData()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Database query error: %v\n", err)
		return
	}

	if rescanLimit > 0 && len(entries) > rescanLimit {
		entries = entries[:rescanLimit]
	}

	if len(entries) == 0 {
		fmt.Println("✅ All entries have complete metadata. Nothing to rescan.")
		return
	}

	client := tmdb.NewClientWithToken(creds.APIKey, creds.Token)

	fmt.Printf("\n🔄 Rescanning %d entries for TMDb metadata...\n\n", len(entries))

	updated, failed := 0, 0
	for i, m := range entries {
		fmt.Printf("  %d/%d  %s", i+1, len(entries), m.CleanTitle)
		if m.Year > 0 {
			fmt.Printf(" (%d)", m.Year)
		}

		// Build search query
		searchTitle := m.CleanTitle
		if m.Year > 0 {
			yearStr := strconv.Itoa(m.Year)
			re := regexp.MustCompile(`\s+` + regexp.QuoteMeta(yearStr) + `$`)
			searchTitle = re.ReplaceAllString(searchTitle, "")
		}
		searchQuery := searchTitle
		if m.Year > 0 {
			searchQuery += " " + strconv.Itoa(m.Year)
		}

		tmdbResults, tmdbErr := client.SearchMulti(searchQuery)
		if tmdbErr != nil {
			fmt.Printf("  ❌ %v\n", tmdbErr)
			errlog.Warn("rescan TMDb search failed for '%s': %v", searchQuery, tmdbErr)
			failed++
			continue
		}
		if len(tmdbResults) == 0 {
			fmt.Println("  ⚠️  No match")
			failed++
			continue
		}

		best := tmdbResults[0]
		m.TmdbID = best.ID
		m.TmdbRating = best.VoteAvg
		m.Popularity = best.Popularity
		m.Description = best.Overview
		m.Genre = tmdb.GenreNames(best.GenreIDs)

		// Fetch full details
		result := cleaner.Result{
			CleanTitle: m.CleanTitle,
			Year:       m.Year,
			Type:       m.Type,
		}
		_ = result // used for type context

		if best.MediaType == "movie" || best.MediaType == "" {
			m.Type = "movie"
			fetchMovieDetails(client, best.ID, &m)
		} else if best.MediaType == "tv" {
			m.Type = "tv"
			fetchTVDetails(client, best.ID, &m)
		}

		// Update in DB
		if m.TmdbID > 0 {
			if updateErr := database.UpdateMediaByTmdbID(&m); updateErr != nil {
				// Try update by ID instead
				if updateErr2 := database.UpdateMediaByID(&m); updateErr2 != nil {
					fmt.Printf("  ❌ DB update failed\n")
					errlog.Error("rescan DB update failed for '%s': %v", m.Title, updateErr2)
					failed++
					continue
				}
			}
		} else {
			if updateErr := database.UpdateMediaByID(&m); updateErr != nil {
				fmt.Printf("  ❌ DB update failed\n")
				errlog.Error("rescan DB update failed for '%s': %v", m.Title, updateErr)
				failed++
				continue
			}
		}

		fmt.Printf("  ✅ ⭐%.1f %s\n", m.TmdbRating, m.Genre)
		updated++
	}

	fmt.Printf("\n📊 Rescan Complete!\n")
	fmt.Printf("   Updated:  %d\n", updated)
	fmt.Printf("   Failed:   %d\n", failed)
	fmt.Printf("   Total:    %d\n\n", len(entries))
}
