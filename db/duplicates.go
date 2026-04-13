// duplicates.go — duplicate detection queries for the media table.
// SHARED: used by cmd/movie_duplicates.go
package db

import "fmt"

// DuplicateGroup represents a set of media records that share a duplicate key.
type DuplicateGroup struct {
	Key   string  // the shared value (e.g. TMDb ID, filename, size)
	Items []Media // media records in this group
}

// FindDuplicatesByTmdbID returns groups of media that share the same tmdb_id.
// Only includes tmdb_id values with 2+ entries.
func (d *DB) FindDuplicatesByTmdbID() ([]DuplicateGroup, error) {
	rows, err := d.Query(`
		SELECT tmdb_id FROM media
		WHERE tmdb_id > 0
		GROUP BY tmdb_id
		HAVING COUNT(*) > 1
		ORDER BY COUNT(*) DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			continue
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var groups []DuplicateGroup
	for _, id := range ids {
		items, err := d.mediaByTmdbIDAll(id)
		if err != nil {
			continue
		}
		groups = append(groups, DuplicateGroup{
			Key:   formatInt(id),
			Items: items,
		})
	}
	return groups, nil
}

// FindDuplicatesByFileName returns groups of media that share the same original_file_name.
func (d *DB) FindDuplicatesByFileName() ([]DuplicateGroup, error) {
	rows, err := d.Query(`
		SELECT original_file_name FROM media
		WHERE original_file_name != ''
		GROUP BY original_file_name
		HAVING COUNT(*) > 1
		ORDER BY COUNT(*) DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			continue
		}
		names = append(names, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var groups []DuplicateGroup
	for _, name := range names {
		items, err := d.mediaByFileName(name)
		if err != nil {
			continue
		}
		groups = append(groups, DuplicateGroup{
			Key:   name,
			Items: items,
		})
	}
	return groups, nil
}

// FindDuplicatesByFileSize returns groups of media that share the same file_size.
// Only considers files with size > 0.
func (d *DB) FindDuplicatesByFileSize() ([]DuplicateGroup, error) {
	rows, err := d.Query(`
		SELECT file_size FROM media
		WHERE file_size > 0
		GROUP BY file_size
		HAVING COUNT(*) > 1
		ORDER BY file_size DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sizes []int64
	for rows.Next() {
		var size int64
		if err := rows.Scan(&size); err != nil {
			continue
		}
		sizes = append(sizes, size)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var groups []DuplicateGroup
	for _, size := range sizes {
		items, err := d.mediaByFileSize(size)
		if err != nil {
			continue
		}
		groups = append(groups, DuplicateGroup{
			Key:   HumanSize(size),
			Items: items,
		})
	}
	return groups, nil
}

// --- internal helpers ---

func (d *DB) mediaByTmdbIDAll(tmdbID int) ([]Media, error) {
	rows, err := d.Query(`SELECT `+mediaColumns+`
		FROM media WHERE tmdb_id = ? ORDER BY id`, tmdbID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMediaRows(rows)
}

func (d *DB) mediaByFileName(name string) ([]Media, error) {
	rows, err := d.Query(`SELECT `+mediaColumns+`
		FROM media WHERE original_file_name = ? ORDER BY id`, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMediaRows(rows)
}

func (d *DB) mediaByFileSize(size int64) ([]Media, error) {
	rows, err := d.Query(`SELECT `+mediaColumns+`
		FROM media WHERE file_size = ? ORDER BY id`, size)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMediaRows(rows)
}

func formatInt(n int) string {
	return fmt.Sprintf("%d", n)
}

// humanSize formats bytes into human-readable form (shared with cmd package).
func HumanSize(b int64) string {
	const (
		gb = 1024 * 1024 * 1024
		mb = 1024 * 1024
		kb = 1024
	)
	switch {
	case b >= gb:
		return fmt.Sprintf("%.1f GB", float64(b)/float64(gb))
	case b >= mb:
		return fmt.Sprintf("%.1f MB", float64(b)/float64(mb))
	case b >= kb:
		return fmt.Sprintf("%.1f KB", float64(b)/float64(kb))
	default:
		return fmt.Sprintf("%d B", b)
	}
}
