package db

import (
	"database/sql"
	"fmt"
	"path/filepath"
)

// mediaColumns is the standard SELECT column list for media queries.
const mediaColumns = `id, title, clean_title, year, type,
	COALESCE(tmdb_id, 0), COALESCE(imdb_id, ''), COALESCE(description, ''),
	COALESCE(imdb_rating, 0), COALESCE(tmdb_rating, 0), COALESCE(popularity, 0),
	COALESCE(genre, ''), COALESCE(director, ''), COALESCE(cast_list, ''),
	COALESCE(thumbnail_path, ''), COALESCE(original_file_name, ''),
	COALESCE(original_file_path, ''), COALESCE(current_file_path, ''),
	COALESCE(file_extension, ''), COALESCE(file_size, 0),
	COALESCE(runtime, 0), COALESCE(language, ''), COALESCE(budget, 0),
	COALESCE(revenue, 0), COALESCE(trailer_url, ''), COALESCE(tagline, '')`

// Media represents a row in the media table.
type Media struct {
	Title            string
	CleanTitle       string
	Type             string // "movie" or "tv"
	ImdbID           string
	Description      string
	Genre            string
	Director         string
	CastList         string
	ThumbnailPath    string
	OriginalFileName string
	OriginalFilePath string
	CurrentFilePath  string
	FileExtension    string
	Language         string
	TrailerURL       string
	Tagline          string
	ID               int64
	FileSize         int64
	Budget           int64
	Revenue          int64
	ImdbRating       float64
	TmdbRating       float64
	Popularity       float64
	Year             int
	TmdbID           int
	Runtime          int
}

// UpdateMediaByID updates an existing record by its primary key ID.
func (d *DB) UpdateMediaByID(m *Media) error {
	var tmdbID interface{}
	if m.TmdbID > 0 {
		tmdbID = m.TmdbID
	}
	_, err := d.Exec(`
		UPDATE media SET title=?, clean_title=?, year=?, type=?, tmdb_id=?, imdb_id=?,
			description=?, imdb_rating=?, tmdb_rating=?, popularity=?, genre=?,
			director=?, cast_list=?, thumbnail_path=?, current_file_path=?,
			file_extension=?, file_size=?,
			runtime=?, language=?, budget=?, revenue=?, trailer_url=?, tagline=?,
			updated_at=CURRENT_TIMESTAMP
		WHERE id=?`,
		m.Title, m.CleanTitle, m.Year, m.Type, tmdbID, m.ImdbID,
		m.Description, m.ImdbRating, m.TmdbRating, m.Popularity, m.Genre,
		m.Director, m.CastList, m.ThumbnailPath, m.CurrentFilePath,
		m.FileExtension, m.FileSize,
		m.Runtime, m.Language, m.Budget, m.Revenue, m.TrailerURL, m.Tagline,
		m.ID,
	)
	return err
}

// InsertMedia inserts a new media record and returns the ID.
func (d *DB) InsertMedia(m *Media) (int64, error) {
	var tmdbID interface{}
	if m.TmdbID > 0 {
		tmdbID = m.TmdbID
	}

	res, err := d.Exec(`
		INSERT INTO media (title, clean_title, year, type, tmdb_id, imdb_id,
			description, imdb_rating, tmdb_rating, popularity, genre, director,
			cast_list, thumbnail_path, original_file_name, original_file_path,
			current_file_path, file_extension, file_size,
			runtime, language, budget, revenue, trailer_url, tagline)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		m.Title, m.CleanTitle, m.Year, m.Type, tmdbID, m.ImdbID,
		m.Description, m.ImdbRating, m.TmdbRating, m.Popularity, m.Genre, m.Director,
		m.CastList, m.ThumbnailPath, m.OriginalFileName, m.OriginalFilePath,
		m.CurrentFilePath, m.FileExtension, m.FileSize,
		m.Runtime, m.Language, m.Budget, m.Revenue, m.TrailerURL, m.Tagline,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateMediaByTmdbID updates an existing record matched by tmdb_id.
func (d *DB) UpdateMediaByTmdbID(m *Media) error {
	_, err := d.Exec(`
		UPDATE media SET title=?, clean_title=?, year=?, type=?, imdb_id=?,
			description=?, imdb_rating=?, tmdb_rating=?, popularity=?, genre=?,
			director=?, cast_list=?, thumbnail_path=?, current_file_path=?,
			file_extension=?, file_size=?,
			runtime=?, language=?, budget=?, revenue=?, trailer_url=?, tagline=?,
			updated_at=CURRENT_TIMESTAMP
		WHERE tmdb_id=?`,
		m.Title, m.CleanTitle, m.Year, m.Type, m.ImdbID,
		m.Description, m.ImdbRating, m.TmdbRating, m.Popularity, m.Genre,
		m.Director, m.CastList, m.ThumbnailPath, m.CurrentFilePath,
		m.FileExtension, m.FileSize,
		m.Runtime, m.Language, m.Budget, m.Revenue, m.TrailerURL, m.Tagline,
		m.TmdbID,
	)
	return err
}

// UpdateMediaPath updates the current file path.
func (d *DB) UpdateMediaPath(mediaID int64, newPath string) error {
	_, err := d.Exec("UPDATE media SET current_file_path = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", newPath, mediaID)
	return err
}

// ListMedia returns paginated media records (only scan-indexed items with a file path).
func (d *DB) ListMedia(offset, limit int) ([]Media, error) {
	rows, err := d.Query(`SELECT `+mediaColumns+`
		FROM media WHERE original_file_path != ''
		ORDER BY clean_title ASC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMediaRows(rows)
}

// SearchMedia searches by title (fuzzy via LIKE).
func (d *DB) SearchMedia(query string) ([]Media, error) {
	rows, err := d.Query(`SELECT `+mediaColumns+`
		FROM media WHERE clean_title LIKE ? OR title LIKE ?
		ORDER BY popularity DESC LIMIT 20`, "%"+query+"%", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMediaRows(rows)
}

// GetMediaByID returns a single media record.
func (d *DB) GetMediaByID(id int64) (*Media, error) {
	row := d.QueryRow(`SELECT `+mediaColumns+` FROM media WHERE id = ?`, id)
	return scanMediaRow(row)
}

// GetMediaByTmdbID returns a media record by its TMDb ID.
func (d *DB) GetMediaByTmdbID(tmdbID int) (*Media, error) {
	row := d.QueryRow(`SELECT `+mediaColumns+` FROM media WHERE tmdb_id = ?`, tmdbID)
	return scanMediaRow(row)
}

// CountMedia returns total count of scan-indexed items, optionally filtered by type.
func (d *DB) CountMedia(mediaType string) (int, error) {
	var count int
	var err error
	if mediaType == "" {
		err = d.QueryRow("SELECT COUNT(*) FROM media WHERE original_file_path != ''").Scan(&count)
	} else {
		err = d.QueryRow("SELECT COUNT(*) FROM media WHERE type = ? AND original_file_path != ''", mediaType).Scan(&count)
	}
	return count, err
}

// ListAllMedia returns all media records that have a file path.
func (d *DB) ListAllMedia() ([]Media, error) {
	rows, err := d.Query(`SELECT `+mediaColumns+`
		FROM media WHERE original_file_path != ''
		ORDER BY clean_title ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMediaRows(rows)
}

// GetMediaWithMissingData returns entries that have empty genre, zero rating, or empty description.
func (d *DB) GetMediaWithMissingData() ([]Media, error) {
	rows, err := d.Query(`SELECT `+mediaColumns+`
		FROM media WHERE original_file_path != ''
		AND (COALESCE(genre, '') = '' OR COALESCE(tmdb_rating, 0) = 0 OR COALESCE(description, '') = '')
		ORDER BY clean_title ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMediaRows(rows)
}

// GetMediaByScanDir returns all media whose original_file_path starts with the given directory prefix.
func (d *DB) GetMediaByScanDir(scanDir string) ([]Media, error) {
	prefix := scanDir
	// Ensure prefix ends with path separator for accurate matching
	if prefix != "" && prefix[len(prefix)-1] != '/' && prefix[len(prefix)-1] != '\\' {
		prefix += string([]byte{filepath.Separator})
	}
	rows, err := d.Query(`SELECT `+mediaColumns+`
		FROM media WHERE original_file_path LIKE ?
		ORDER BY clean_title ASC`, prefix+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMediaRows(rows)
}

// DeleteMediaByIDs deletes multiple media records by their IDs.
func (d *DB) DeleteMediaByIDs(ids []int64) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	tx, err := d.Begin()
	if err != nil {
		return 0, err
	}
	deleted := 0
	for _, id := range ids {
		if _, err := tx.Exec("DELETE FROM media WHERE id = ?", id); err != nil {
			tx.Rollback()
			return deleted, err
		}
		deleted++
	}
	return deleted, tx.Commit()
}

// FileSizeStats returns total file size, largest file size, and smallest file size (non-zero) across all media.
func (d *DB) FileSizeStats() (total int64, largest int64, smallest int64, err error) {
	err = d.QueryRow(`
		SELECT COALESCE(SUM(file_size), 0),
		       COALESCE(MAX(file_size), 0),
		       COALESCE(MIN(NULLIF(file_size, 0)), 0)
		FROM media WHERE file_size > 0`).Scan(&total, &largest, &smallest)
	return
}

// MediaByType returns media filtered by type.
func (d *DB) MediaByType(mediaType string, limit int) ([]Media, error) {
	rows, err := d.Query(`SELECT `+mediaColumns+`
		FROM media WHERE type = ? ORDER BY popularity DESC LIMIT ?`, mediaType, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMediaRows(rows)
}

// TopGenres returns genres sorted by frequency.
func (d *DB) TopGenres(limit int) (map[string]int, error) {
	rows, err := d.Query("SELECT genre FROM media WHERE genre != ''")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var genre string
		if err := rows.Scan(&genre); err != nil {
			return nil, fmt.Errorf("scanning genre: %w", err)
		}
		for _, g := range splitCSV(genre) {
			counts[g]++
		}
	}
	return counts, nil
}

// scanMediaRow scans a single media row.
func scanMediaRow(row *sql.Row) (*Media, error) {
	m := &Media{}
	err := row.Scan(&m.ID, &m.Title, &m.CleanTitle, &m.Year, &m.Type,
		&m.TmdbID, &m.ImdbID, &m.Description, &m.ImdbRating, &m.TmdbRating,
		&m.Popularity, &m.Genre, &m.Director, &m.CastList, &m.ThumbnailPath,
		&m.OriginalFileName, &m.OriginalFilePath, &m.CurrentFilePath,
		&m.FileExtension, &m.FileSize,
		&m.Runtime, &m.Language, &m.Budget, &m.Revenue, &m.TrailerURL, &m.Tagline)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// scanMediaRows scans multiple media rows from a query result.
func scanMediaRows(rows *sql.Rows) ([]Media, error) {
	var list []Media
	for rows.Next() {
		var m Media
		if err := rows.Scan(&m.ID, &m.Title, &m.CleanTitle, &m.Year, &m.Type,
			&m.TmdbID, &m.ImdbID, &m.Description, &m.ImdbRating, &m.TmdbRating,
			&m.Popularity, &m.Genre, &m.Director, &m.CastList, &m.ThumbnailPath,
			&m.OriginalFileName, &m.OriginalFilePath, &m.CurrentFilePath,
			&m.FileExtension, &m.FileSize,
			&m.Runtime, &m.Language, &m.Budget, &m.Revenue, &m.TrailerURL, &m.Tagline); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, rows.Err()
}
