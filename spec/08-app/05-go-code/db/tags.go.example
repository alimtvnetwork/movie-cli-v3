package db

import "fmt"

// TagCount holds a tag name and its usage count.
type TagCount struct {
	Tag   string
	Count int
}

// AddTag inserts a tag for a media item.
// Returns UNIQUE constraint error if tag already exists.
func (d *DB) AddTag(mediaID int, tag string) error {
	_, err := d.Exec(
		`INSERT INTO tags (media_id, tag) VALUES (?, ?)`,
		mediaID, tag,
	)
	return err
}

// RemoveTag deletes a tag from a media item.
// Returns (true, nil) if deleted, (false, nil) if tag didn't exist.
func (d *DB) RemoveTag(mediaID int, tag string) (bool, error) {
	result, err := d.Exec(
		`DELETE FROM tags WHERE media_id = ? AND tag = ?`,
		mediaID, tag,
	)
	if err != nil {
		return false, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}

// GetTagsByMediaID returns all tags for a specific media item.
func (d *DB) GetTagsByMediaID(mediaID int) ([]string, error) {
	rows, err := d.Query(
		`SELECT tag FROM tags WHERE media_id = ? ORDER BY tag`,
		mediaID,
	)
	if err != nil {
		return nil, fmt.Errorf("query tags: %w", err)
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, fmt.Errorf("scan tag: %w", err)
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}

// GetAllTagCounts returns all unique tags with their usage count,
// ordered by count descending.
func (d *DB) GetAllTagCounts() ([]TagCount, error) {
	rows, err := d.Query(
		`SELECT tag, COUNT(*) as cnt FROM tags GROUP BY tag ORDER BY cnt DESC, tag ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("query tag counts: %w", err)
	}
	defer rows.Close()

	var counts []TagCount
	for rows.Next() {
		var tc TagCount
		if err := rows.Scan(&tc.Tag, &tc.Count); err != nil {
			return nil, fmt.Errorf("scan tag count: %w", err)
		}
		counts = append(counts, tc)
	}
	return counts, rows.Err()
}
