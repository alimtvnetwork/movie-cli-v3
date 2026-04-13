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
package cmd

import (
	"strings"

	"github.com/alimtvnetwork/movie-cli-v3/db"
	"github.com/alimtvnetwork/movie-cli-v3/tmdb"
)

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
