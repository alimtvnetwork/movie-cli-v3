// cleanup.go — find stale media entries where the file no longer exists on disk.
// SHARED: used by cmd/movie_cleanup.go
package db

import "os"

// StaleEntry represents a media record whose file is missing from disk.
type StaleEntry struct {
	Media    Media
	FilePath string // the path that was checked
}

// FindStaleEntries returns media records where current_file_path or
// original_file_path no longer exists on disk.
func (d *DB) FindStaleEntries(limit int) ([]StaleEntry, error) {
	rows, err := d.Query(`
		SELECT id, title, clean_title, year, type, tmdb_id, imdb_id,
			description, imdb_rating, tmdb_rating, popularity, genre,
			director, cast_list, thumbnail_path, original_file_name,
			original_file_path, current_file_path, file_extension, file_size
		FROM media
		WHERE original_file_path != ''
		ORDER BY clean_title ASC
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	all, err := scanMediaRows(rows)
	if err != nil {
		return nil, err
	}

	var stale []StaleEntry
	for _, m := range all {
		path := m.CurrentFilePath
		if path == "" {
			path = m.OriginalFilePath
		}
		if path == "" {
			continue
		}
		if _, err := os.Stat(path); os.IsNotExist(err) {
			stale = append(stale, StaleEntry{Media: m, FilePath: path})
		}
	}
	return stale, nil
}

// DeleteMedia removes a media record by ID.
func (d *DB) DeleteMedia(id int64) error {
	_, err := d.Exec("DELETE FROM media WHERE id = ?", id)
	return err
}
