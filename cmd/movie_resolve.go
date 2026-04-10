// movie_resolve.go — shared media resolver
//
// -- Shared helper exported from this file --
//
//	resolveMediaByQuery(db, query) (*Media, error)
//	    Resolves a media item from the local DB by numeric ID or fuzzy
//	    title match (exact → prefix → first result).
//
// Consumers: movie_info.go, movie_play.go, movie_ls.go (detail view)
//
// All commands that accept an <id-or-title> argument should use this
// helper to keep resolution logic consistent.  Do NOT duplicate the
// ID-parse → exact-match → prefix-match → fallback chain elsewhere.
package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mahin/mahin-cli-v2/db"
)

// resolveMediaByQuery resolves a media item by numeric ID or fuzzy title query.
func resolveMediaByQuery(database *db.DB, query string) (*db.Media, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, fmt.Errorf("empty media identifier")
	}

	if id, parseErr := strconv.ParseInt(query, 10, 64); parseErr == nil {
		m, getErr := database.GetMediaByID(id)
		if getErr != nil {
			return nil, fmt.Errorf("media not found for ID %d", id)
		}
		return m, nil
	}

	results, searchErr := database.SearchMedia(query)
	if searchErr != nil {
		return nil, searchErr
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("media not found for %q", query)
	}

	for i := range results {
		if strings.EqualFold(results[i].CleanTitle, query) || strings.EqualFold(results[i].Title, query) {
			picked := results[i]
			return &picked, nil
		}
	}

	queryLower := strings.ToLower(query)
	for i := range results {
		if strings.HasPrefix(strings.ToLower(results[i].CleanTitle), queryLower) ||
			strings.HasPrefix(strings.ToLower(results[i].Title), queryLower) {
			picked := results[i]
			return &picked, nil
		}
	}

	picked := results[0]
	return &picked, nil
}
