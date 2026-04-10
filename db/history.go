package db

// MoveRecord represents a row in move_history.
type MoveRecord struct {
	FromPath         string
	ToPath           string
	OriginalFileName string
	NewFileName      string
	ID               int64
	MediaID          int64
	Undone           bool
}

// InsertMoveHistory logs a move operation.
func (d *DB) InsertMoveHistory(mediaID int64, fromPath, toPath, origName, newName string) error {
	_, err := d.Exec(`
		INSERT INTO move_history (media_id, from_path, to_path, original_file_name, new_file_name)
		VALUES (?, ?, ?, ?, ?)`, mediaID, fromPath, toPath, origName, newName)
	return err
}

// GetLastMove returns the latest un-undone move.
func (d *DB) GetLastMove() (*MoveRecord, error) {
	row := d.QueryRow(`
		SELECT id, media_id, from_path, to_path, original_file_name, new_file_name, undone
		FROM move_history WHERE undone = 0 ORDER BY moved_at DESC LIMIT 1`)
	r := &MoveRecord{}
	err := row.Scan(&r.ID, &r.MediaID, &r.FromPath, &r.ToPath, &r.OriginalFileName, &r.NewFileName, &r.Undone)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// MarkMoveUndone marks a move_history record as undone.
func (d *DB) MarkMoveUndone(id int64) error {
	_, err := d.Exec("UPDATE move_history SET undone = 1 WHERE id = ?", id)
	return err
}

// InsertScanHistory logs a scan operation.
func (d *DB) InsertScanHistory(folder string, total, movies, tv int) error {
	_, err := d.Exec(`
		INSERT INTO scan_history (folder_path, total_files, movies_found, tv_found)
		VALUES (?, ?, ?, ?)`, folder, total, movies, tv)
	return err
}
