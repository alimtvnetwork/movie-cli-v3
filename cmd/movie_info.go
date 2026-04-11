// movie_info.go — movie info <id-or-title>
//
// Accepts a numeric ID (from library) or a title string.
// Checks local DB first; if not found by title, falls back to TMDb API,
// fetches full details, stores in DB, then displays.
//
// -- Shared helpers exported from this file --
//
//	fetchMovieDetails(client, tmdbID, m)  — populate Media with TMDb movie details + credits
//	fetchTVDetails(client, tmdbID, m)     — populate Media with TMDb TV details + credits
//
// Consumers: movie_scan.go (scan + metadata fetch), movie_info.go (info fallback)
//
// These helpers centralize all TMDb detail+credit fetching so that scan
// and info share identical enrichment logic.  Any change to field mapping
// or credit extraction should happen here only.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mahin/mahin-cli-v2/cleaner"
	"github.com/mahin/mahin-cli-v2/db"
	"github.com/mahin/mahin-cli-v2/tmdb"
)

var movieInfoCmd = &cobra.Command{
	Use:   "info [id or title]",
	Short: "Show detailed info for a movie or TV show",
	Long: `Display full metadata for a media item.

If a numeric ID is given, it looks up the item from your local library.
If a title is given, it first searches the local database. If not found,
it queries the TMDb API, saves the result, and then displays it.`,
	Args: cobra.MinimumNArgs(1),
	Run:  runMovieInfo,
}

func runMovieInfo(cmd *cobra.Command, args []string) {
	database, err := db.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Database error: %v\n", err)
		return
	}
	defer database.Close()

	query := strings.Join(args, " ")

	// 1) Try local DB first (by ID or title)
	m, resolveErr := resolveMediaByQuery(database, query)
	if resolveErr == nil {
		fmt.Println("📚 Found in local library:")
		fmt.Println()
		printMediaDetail(m)
		return
	}

	// 3) Not in DB — fall back to TMDb API
	fmt.Printf("🔎 Not found locally. Searching TMDb for: %s\n\n", query)

	apiKey, cfgErr := database.GetConfig("tmdb_api_key")
	if cfgErr != nil {
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
	tmdbResults, searchErr := client.SearchMulti(query)
	if searchErr != nil {
		fmt.Fprintf(os.Stderr, "❌ TMDb search error: %v\n", searchErr)
		return
	}
	if len(tmdbResults) == 0 {
		fmt.Println("📭 No results found on TMDb either.")
		return
	}

	// Pick the first (most relevant) result
	selected := tmdbResults[0]
	title := selected.GetDisplayTitle()
	year := selected.GetYear()
	yearInt := 0
	if year != "" {
		yearInt, _ = strconv.Atoi(year)
	}

	fmt.Printf("⏳ Fetching details for: %s (%s)...\n", title, year)

	// Check if this TMDb ID already exists in DB (avoid duplicates)
	existing, existErr := database.GetMediaByTmdbID(selected.ID)
	if existErr != nil && existErr.Error() != "sql: no rows in result set" {
		fmt.Fprintf(os.Stderr, "⚠️  DB lookup error: %v\n", existErr)
	}
	if existing != nil {
		fmt.Println("📚 Already in your library:")
		fmt.Println()
		printMediaDetail(existing)
		return
	}

	// Build media record with full details
	m = &db.Media{
		Title:       title,
		CleanTitle:  title,
		Year:        yearInt,
		TmdbID:      selected.ID,
		TmdbRating:  selected.VoteAvg,
		Popularity:  selected.Popularity,
		Description: selected.Overview,
		Genre:       tmdb.GenreNames(selected.GenreIDs),
	}

	if selected.MediaType == "movie" || selected.MediaType == "" {
		m.Type = "movie"
		fetchMovieDetails(client, selected.ID, m)
	} else if selected.MediaType == "tv" {
		m.Type = "tv"
		fetchTVDetails(client, selected.ID, m)
	}

	// Download thumbnail
	if selected.PosterPath != "" {
		slug := cleaner.ToSlug(m.CleanTitle)
		if m.Year > 0 {
			slug += "-" + strconv.Itoa(m.Year)
		}
		thumbDir := filepath.Join(database.BasePath, "thumbnails", slug)
		if mkdirErr := os.MkdirAll(thumbDir, 0755); mkdirErr != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Cannot create thumbnail dir: %v\n", mkdirErr)
		}
		thumbPath := filepath.Join(thumbDir, slug+".jpg")
		if dlErr := client.DownloadPoster(selected.PosterPath, thumbPath); dlErr == nil {
			m.ThumbnailPath = thumbPath
		}
	}

	// Save to DB
	_, insertErr := database.InsertMedia(m)
	if insertErr != nil {
		if m.TmdbID > 0 {
			insertErr = database.UpdateMediaByTmdbID(m)
		}
		if insertErr != nil {
			fmt.Fprintf(os.Stderr, "❌ DB error: %v\n", insertErr)
			return
		}
	}

	fmt.Println()
	fmt.Println("✅ Saved to your library!")
	fmt.Println()
	printMediaDetail(m)
}

// fetchMovieDetails populates a Media record with TMDb movie details + credits + videos.
func fetchMovieDetails(client *tmdb.Client, tmdbID int, m *db.Media) {
	details, detailErr := client.GetMovieDetails(tmdbID)
	if detailErr == nil {
		m.ImdbID = details.ImdbID
		m.Title = details.Title
		m.Runtime = details.Runtime
		m.Language = details.OriginalLanguage
		m.Budget = details.Budget
		m.Revenue = details.Revenue
		m.Tagline = details.Tagline
		genres := make([]string, len(details.Genres))
		for i, g := range details.Genres {
			genres[i] = g.Name
		}
		m.Genre = strings.Join(genres, ", ")
	}

	credits, creditErr := client.GetMovieCredits(tmdbID)
	if creditErr == nil {
		var directors, castNames []string
		for _, c := range credits.Crew {
			if c.Job == "Director" {
				directors = append(directors, c.Name)
			}
		}
		m.Director = strings.Join(directors, ", ")

		for i, c := range credits.Cast {
			if i >= 10 {
				break
			}
			castNames = append(castNames, c.Name)
		}
		m.CastList = strings.Join(castNames, ", ")
	}

	videos, vidErr := client.GetMovieVideos(tmdbID)
	if vidErr == nil {
		m.TrailerURL = tmdb.TrailerURL(videos)
	}
}

// fetchTVDetails populates a Media record with TMDb TV details + credits + videos.
func fetchTVDetails(client *tmdb.Client, tmdbID int, m *db.Media) {
	details, detailErr := client.GetTVDetails(tmdbID)
	if detailErr == nil {
		m.Title = details.Name
		m.Language = details.OriginalLanguage
		m.Tagline = details.Tagline
		if len(details.EpisodeRunTime) > 0 {
			m.Runtime = details.EpisodeRunTime[0]
		}
		genres := make([]string, len(details.Genres))
		for i, g := range details.Genres {
			genres[i] = g.Name
		}
		m.Genre = strings.Join(genres, ", ")
	}

	credits, creditErr := client.GetTVCredits(tmdbID)
	if creditErr == nil {
		var directors, castNames []string
		for _, c := range credits.Crew {
			if c.Job == "Director" || c.Job == "Executive Producer" {
				directors = append(directors, c.Name)
			}
		}
		if len(directors) > 5 {
			directors = directors[:5]
		}
		m.Director = strings.Join(directors, ", ")

		for i, c := range credits.Cast {
			if i >= 10 {
				break
			}
			castNames = append(castNames, c.Name)
		}
		m.CastList = strings.Join(castNames, ", ")
	}

	videos, vidErr := client.GetTVVideos(tmdbID)
	if vidErr == nil {
		m.TrailerURL = tmdb.TrailerURL(videos)
	}
}
