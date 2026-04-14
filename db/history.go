package db

// MoveRecord represents a row in move_history.
type MoveRecord struct {
	FromPath         string
	ToPath           string
	OriginalFileName string
	NewFileName      string
	MovedAt          string
	ID               int64
	MediaID          int64
	Undone           bool
}

// ListMoveHistory returns all move records ordered by most recent first.
func (d *DB) ListMoveHistory(limit int) ([]MoveRecord, error) {
	if limit <= 0 {
		limit = 1000
	}
	rows, err := d.Query(`
		SELECT id, media_id, from_path, to_path, original_file_name, new_file_name, moved_at, undone
		FROM move_history ORDER BY moved_at DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []MoveRecord
	for rows.Next() {
		var r MoveRecord
		if scanErr := rows.Scan(&r.ID, &r.MediaID, &r.FromPath, &r.ToPath,
			&r.OriginalFileName, &r.NewFileName, &r.MovedAt, &r.Undone); scanErr != nil {
			return nil, scanErr
		}
		records = append(records, r)
	}
	return records, rows.Err()
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

// ScanRecord represents a row in scan_history.
type ScanRecord struct {
	ID         int64
	FolderPath string
	TotalFiles int
	Movies     int
	TV         int
	ScannedAt  string
}

// ListScanFolders returns distinct scanned folder paths, most recent first.
func (d *DB) ListScanFolders(limit int) ([]ScanRecord, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := d.Query(`
		SELECT id, folder_path, total_files, movies_found, tv_found, scanned_at
		FROM scan_history
		ORDER BY scanned_at DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []ScanRecord
	for rows.Next() {
		var r ScanRecord
		if scanErr := rows.Scan(&r.ID, &r.FolderPath, &r.TotalFiles, &r.Movies, &r.TV, &r.ScannedAt); scanErr != nil {
			return nil, scanErr
		}
		records = append(records, r)
	}
	return records, rows.Err()
}

// ListDistinctScanFolders returns unique folder paths from scan history.
func (d *DB) ListDistinctScanFolders() ([]string, error) {
	rows, err := d.Query(`
		SELECT DISTINCT folder_path FROM scan_history ORDER BY folder_path`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folders []string
	for rows.Next() {
		var f string
		if scanErr := rows.Scan(&f); scanErr != nil {
			return nil, scanErr
		}
		folders = append(folders, f)
	}
	return folders, rows.Err()
}

// InsertScanHistory logs a scan operation.
func (d *DB) InsertScanHistory(folder string, total, movies, tv int) error {
	_, err := d.Exec(`
		INSERT INTO scan_history (folder_path, total_files, movies_found, tv_found)
		VALUES (?, ?, ?, ?)`, folder, total, movies, tv)
	return err
}
