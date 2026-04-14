// movie_rest_handlers.go — additional REST API handlers for tags, similar, and watched.
package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/movie-cli-v3/db"
	"github.com/alimtvnetwork/movie-cli-v3/tmdb"
)

// handleTags handles POST /api/tags (add/remove tags).
//
//	POST body: {"media_id": 1, "tag": "favorite"}
//	DELETE body: {"media_id": 1, "tag": "favorite"}
//	GET /api/tags?media_id=1 — list tags for a media item
//	GET /api/tags — list all tags with counts
func handleTags(w http.ResponseWriter, r *http.Request, database *db.DB) {
	switch r.Method {
	case http.MethodGet:
		idStr := r.URL.Query().Get("media_id")
		if idStr != "" {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "invalid media_id", http.StatusBadRequest)
				return
			}
			tags, tagErr := database.GetTagsByMediaID(id)
			if tagErr != nil {
				http.Error(w, tagErr.Error(), http.StatusInternalServerError)
				return
			}
			writeJSON(w, map[string]interface{}{"media_id": id, "tags": tags})
		} else {
			counts, cErr := database.GetAllTagCounts()
			if cErr != nil {
				http.Error(w, cErr.Error(), http.StatusInternalServerError)
				return
			}
			writeJSON(w, counts)
		}

	case http.MethodPost:
		var req struct {
			MediaID int    `json:"media_id"`
			Tag     string `json:"tag"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if req.MediaID == 0 || req.Tag == "" {
			http.Error(w, "media_id and tag are required", http.StatusBadRequest)
			return
		}
		if err := database.AddTag(req.MediaID, req.Tag); err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		writeJSON(w, map[string]string{"status": "added", "tag": req.Tag})

	case http.MethodDelete:
		var req struct {
			MediaID int    `json:"media_id"`
			Tag     string `json:"tag"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		removed, err := database.RemoveTag(req.MediaID, req.Tag)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !removed {
			http.Error(w, "tag not found", http.StatusNotFound)
			return
		}
		writeJSON(w, map[string]string{"status": "removed", "tag": req.Tag})

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleSimilar handles GET /api/media/{id}/similar — fetches TMDb recommendations.
func handleSimilar(w http.ResponseWriter, r *http.Request, database *db.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse ID from path: /api/media/{id}/similar
	path := strings.TrimPrefix(r.URL.Path, "/api/media/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "similar" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	m, getErr := database.GetMediaByID(id)
	if getErr != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if m.TmdbID == 0 {
		writeJSON(w, map[string]interface{}{"media_id": id, "similar": []interface{}{}, "message": "no TMDb ID available"})
		return
	}

	// Get TMDb credentials
	apiKey, _ := database.GetConfig("tmdb_api_key")
	token, _ := database.GetConfig("tmdb_token")
	client := tmdb.NewClientWithToken(apiKey, token)

	results, recErr := client.GetRecommendations(m.TmdbID, m.Type, 1)
	if recErr != nil {
		http.Error(w, fmt.Sprintf("TMDb error: %v", recErr), http.StatusBadGateway)
		return
	}

	writeJSON(w, map[string]interface{}{
		"media_id": id,
		"title":    m.Title,
		"similar":  results,
	})
}

// handleWatched handles PATCH /api/media/{id}/watched — marks a media item as watched.
func handleWatched(w http.ResponseWriter, r *http.Request, database *db.DB) {
	if r.Method != http.MethodPatch {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse ID from path: /api/media/{id}/watched
	path := strings.TrimPrefix(r.URL.Path, "/api/media/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "watched" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Check media exists
	m, getErr := database.GetMediaByID(id)
	if getErr != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	// Update watchlist status if in watchlist
	database.Exec("UPDATE watchlist SET status = 'watched', watched_at = CURRENT_TIMESTAMP WHERE tmdb_id = ?", m.TmdbID)

	// Also add a "watched" tag
	_ = database.AddTag(int(id), "watched")

	writeJSON(w, map[string]interface{}{
		"status":   "marked_watched",
		"media_id": id,
		"title":    m.Title,
	})
}
