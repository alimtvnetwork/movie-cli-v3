// movie_rescan_helper.go — shared rescan logic used by both scan and rescan commands
package cmd

import (
	"regexp"
	"strconv"

	"github.com/alimtvnetwork/movie-cli-v3/db"
	"github.com/alimtvnetwork/movie-cli-v3/errlog"
	"github.com/alimtvnetwork/movie-cli-v3/tmdb"
)

// mediaNeedsRescan returns true if the entry is missing genre, rating, or description.
func mediaNeedsRescan(m *db.Media) bool {
	return m.Genre == "" || m.TmdbRating == 0 || m.Description == ""
}

// rescanMediaEntry re-fetches TMDb metadata for a single media entry.
// Returns true if the entry was updated successfully.
func rescanMediaEntry(database *db.DB, client *tmdb.Client, m *db.Media) bool {
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
		errlog.Warn("rescan TMDb search failed for '%s': %v", searchQuery, tmdbErr)
		return false
	}
	if len(tmdbResults) == 0 {
		return false
	}

	best := tmdbResults[0]
	m.TmdbID = best.ID
	m.TmdbRating = best.VoteAvg
	m.Popularity = best.Popularity
	m.Description = best.Overview
	m.Genre = tmdb.GenreNames(best.GenreIDs)

	if best.MediaType == "movie" || best.MediaType == "" {
		m.Type = "movie"
		fetchMovieDetails(client, best.ID, m)
	} else if best.MediaType == "tv" {
		m.Type = "tv"
		fetchTVDetails(client, best.ID, m)
	}

	// Update in DB
	if m.TmdbID > 0 {
		if updateErr := database.UpdateMediaByTmdbID(m); updateErr != nil {
			if updateErr2 := database.UpdateMediaByID(m); updateErr2 != nil {
				errlog.Error("rescan DB update failed for '%s': %v", m.Title, updateErr2)
				return false
			}
		}
	} else {
		if updateErr := database.UpdateMediaByID(m); updateErr != nil {
			errlog.Error("rescan DB update failed for '%s': %v", m.Title, updateErr)
			return false
		}
	}

	return true
}
