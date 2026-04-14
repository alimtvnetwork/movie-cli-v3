// movie_fetch_details.go — shared TMDb detail+credit fetching helpers
//
// -- Shared helpers --
//
//	fetchMovieDetails(client, tmdbID, m)  — populate Media with TMDb movie details + credits
//	fetchTVDetails(client, tmdbID, m)     — populate Media with TMDb TV details + credits
//
// Consumers: movie_scan_process.go, movie_info.go, movie_search.go
//
// These helpers centralize all TMDb detail+credit fetching so that scan,
// info, and search share identical enrichment logic. Any change to field
// mapping or credit extraction should happen here only.
//
// Error handling per spec/02-error-manage-spec/04-runtime-error-handling.md:
// - Network/timeout errors are logged as warnings, not fatal
// - Each sub-request (details, credits, videos) fails independently
// - The scan continues with whatever data was successfully fetched
package cmd

import (
	"errors"
	"strings"

	"github.com/alimtvnetwork/movie-cli-v3/db"
	"github.com/alimtvnetwork/movie-cli-v3/errlog"
	"github.com/alimtvnetwork/movie-cli-v3/tmdb"
)

// fetchMovieDetails populates a Media record with TMDb movie details + credits + videos.
func fetchMovieDetails(client *tmdb.Client, tmdbID int, m *db.Media) {
	details, detailErr := client.GetMovieDetails(tmdbID)
	if detailErr != nil {
		logTMDbSubError("movie details", tmdbID, detailErr)
	}
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
	if creditErr != nil {
		logTMDbSubError("movie credits", tmdbID, creditErr)
	}
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
	if vidErr != nil {
		logTMDbSubError("movie videos", tmdbID, vidErr)
	}
	if vidErr == nil {
		m.TrailerURL = tmdb.TrailerURL(videos)
	}
}

// fetchTVDetails populates a Media record with TMDb TV details + credits + videos.
func fetchTVDetails(client *tmdb.Client, tmdbID int, m *db.Media) {
	details, detailErr := client.GetTVDetails(tmdbID)
	if detailErr != nil {
		logTMDbSubError("TV details", tmdbID, detailErr)
	}
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
	if creditErr != nil {
		logTMDbSubError("TV credits", tmdbID, creditErr)
	}
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
	if vidErr != nil {
		logTMDbSubError("TV videos", tmdbID, vidErr)
	}
	if vidErr == nil {
		m.TrailerURL = tmdb.TrailerURL(videos)
	}
}

// logTMDbSubError logs a TMDb sub-request error with proper classification per spec.
func logTMDbSubError(what string, tmdbID int, err error) {
	switch {
	case errors.Is(err, tmdb.ErrAuthInvalid):
		errlog.Error("❌ TMDb API key is invalid. Run: movie config set tmdb_api_key YOUR_KEY")
	case errors.Is(err, tmdb.ErrTimeout):
		errlog.Warn("⚠️ TMDb %s request timed out for ID %d — skipping", what, tmdbID)
	case errors.Is(err, tmdb.ErrNetworkError):
		errlog.Warn("⚠️ Network unavailable — skipping %s for ID %d", what, tmdbID)
	case errors.Is(err, tmdb.ErrServerError):
		errlog.Warn("⚠️ TMDb temporarily unavailable — skipping %s for ID %d", what, tmdbID)
	case errors.Is(err, tmdb.ErrRateLimited):
		errlog.Warn("TMDb rate limited — skipping %s for ID %d", what, tmdbID)
	default:
		errlog.Warn("TMDb %s failed for ID %d: %v", what, tmdbID, err)
	}
}
