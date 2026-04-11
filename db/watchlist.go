// watchlist.go — CRUD for the watchlist table (to-watch / watched tracking).
// SHARED: used by cmd/movie_watch.go
package db

import "database/sql"

// WatchlistEntry represents a row in the watchlist table.
type WatchlistEntry struct {
	ID        int64
	MediaID   sql.NullInt64
	TmdbID    int
	Title     string
	Year      int
	Type      string // "movie" or "tv"
	Status    string // "to-watch" or "watched"
	AddedAt   string
	WatchedAt sql.NullString
}

// AddToWatchlist inserts or updates a watchlist entry as "to-watch".
func (d *DB) AddToWatchlist(tmdbID int, title string, year int, mediaType string, mediaID int64) error {
	var mid sql.NullInt64
	if mediaID > 0 {
		mid = sql.NullInt64{Int64: mediaID, Valid: true}
	}
	_, err := d.Exec(`
		INSERT INTO watchlist (tmdb_id, title, year, type, status, media_id)
		VALUES (?, ?, ?, ?, 'to-watch', ?)
		ON CONFLICT(tmdb_id) DO UPDATE SET
			title = excluded.title,
			year  = excluded.year,
			type  = excluded.type,
			media_id = COALESCE(excluded.media_id, watchlist.media_id)`,
		tmdbID, title, year, mediaType, mid)
	return err
}

// MarkWatched updates a watchlist entry to "watched".
func (d *DB) MarkWatched(tmdbID int) error {
	_, err := d.Exec(`
		UPDATE watchlist SET status = 'watched', watched_at = CURRENT_TIMESTAMP
		WHERE tmdb_id = ?`, tmdbID)
	return err
}

// MarkToWatch updates a watchlist entry back to "to-watch".
func (d *DB) MarkToWatch(tmdbID int) error {
	_, err := d.Exec(`
		UPDATE watchlist SET status = 'to-watch', watched_at = NULL
		WHERE tmdb_id = ?`, tmdbID)
	return err
}

// RemoveFromWatchlist deletes a watchlist entry.
func (d *DB) RemoveFromWatchlist(tmdbID int) error {
	_, err := d.Exec("DELETE FROM watchlist WHERE tmdb_id = ?", tmdbID)
	return err
}

// ListWatchlist returns entries filtered by status ("to-watch", "watched", or "" for all).
func (d *DB) ListWatchlist(status string) ([]WatchlistEntry, error) {
	var rows *sql.Rows
	var err error
	if status == "" {
		rows, err = d.Query(`
			SELECT id, media_id, tmdb_id, title, year, type, status, added_at, watched_at
			FROM watchlist ORDER BY added_at DESC`)
	} else {
		rows, err = d.Query(`
			SELECT id, media_id, tmdb_id, title, year, type, status, added_at, watched_at
			FROM watchlist WHERE status = ? ORDER BY added_at DESC`, status)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []WatchlistEntry
	for rows.Next() {
		var e WatchlistEntry
		if err := rows.Scan(&e.ID, &e.MediaID, &e.TmdbID, &e.Title, &e.Year,
			&e.Type, &e.Status, &e.AddedAt, &e.WatchedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

// GetWatchlistByTmdbID returns a single watchlist entry.
func (d *DB) GetWatchlistByTmdbID(tmdbID int) (*WatchlistEntry, error) {
	row := d.QueryRow(`
		SELECT id, media_id, tmdb_id, title, year, type, status, added_at, watched_at
		FROM watchlist WHERE tmdb_id = ?`, tmdbID)
	var e WatchlistEntry
	err := row.Scan(&e.ID, &e.MediaID, &e.TmdbID, &e.Title, &e.Year,
		&e.Type, &e.Status, &e.AddedAt, &e.WatchedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}
