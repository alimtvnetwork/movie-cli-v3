// action_history.go — action_history table: migration, types, and helpers.
//
// Tracks all reversible state changes (scan adds/removes, deletions, popouts,
// rescans) so that movie undo / redo / history can operate on them.
package db

import (
	"database/sql"
	"fmt"
)

// ActionType enumerates the allowed action_type values.
type ActionType string

const (
	ActionScanAdd      ActionType = "scan_add"
	ActionScanRemove   ActionType = "scan_remove"
	ActionDelete       ActionType = "delete"
	ActionPopout       ActionType = "popout"
	ActionRestore      ActionType = "restore"
	ActionRescanUpdate ActionType = "rescan_update"
)

// ActionRecord represents a row in action_history.
type ActionRecord struct {
	ID            int64
	ActionType    ActionType
	MediaID       sql.NullInt64
	MediaSnapshot string // JSON snapshot of the media row before change
	Detail        string // human-readable description
	BatchID       string // groups related actions (UUID)
	Undone        bool
	CreatedAt     string
}

// migrateActionHistory creates the action_history table and indexes.
func (d *DB) migrateActionHistory() error {
	_, err := d.Exec(`
		CREATE TABLE IF NOT EXISTS action_history (
			id             INTEGER PRIMARY KEY AUTOINCREMENT,
			action_type    TEXT NOT NULL CHECK(action_type IN (
				'scan_add', 'scan_remove', 'delete', 'popout', 'restore', 'rescan_update'
			)),
			media_id       INTEGER,
			media_snapshot TEXT,
			detail         TEXT,
			batch_id       TEXT,
			undone         INTEGER DEFAULT 0,
			created_at     DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (media_id) REFERENCES media(id) ON DELETE SET NULL
		);
		CREATE INDEX IF NOT EXISTS idx_action_history_type   ON action_history(action_type);
		CREATE INDEX IF NOT EXISTS idx_action_history_batch  ON action_history(batch_id);
		CREATE INDEX IF NOT EXISTS idx_action_history_undone ON action_history(undone);
	`)
	if err != nil {
		return fmt.Errorf("migrate action_history: %w", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Insert helpers
// ---------------------------------------------------------------------------

// InsertAction logs a state-changing action to action_history.
func (d *DB) InsertAction(actionType ActionType, mediaID sql.NullInt64, snapshot, detail, batchID string) (int64, error) {
	res, err := d.Exec(`
		INSERT INTO action_history (action_type, media_id, media_snapshot, detail, batch_id)
		VALUES (?, ?, ?, ?, ?)`,
		string(actionType), mediaID, snapshot, detail, batchID,
	)
	if err != nil {
		return 0, fmt.Errorf("insert action (%s): %w", actionType, err)
	}
	return res.LastInsertId()
}

// InsertActionSimple is a convenience wrapper when media_id is a plain int64.
func (d *DB) InsertActionSimple(actionType ActionType, mediaID int64, snapshot, detail, batchID string) (int64, error) {
	mid := sql.NullInt64{Int64: mediaID, Valid: mediaID > 0}
	return d.InsertAction(actionType, mid, snapshot, detail, batchID)
}

// ---------------------------------------------------------------------------
// Query helpers
// ---------------------------------------------------------------------------

// GetLastUndoableAction returns the most recent un-undone action.
func (d *DB) GetLastUndoableAction() (*ActionRecord, error) {
	row := d.QueryRow(`
		SELECT id, action_type, media_id, media_snapshot, detail, batch_id, undone, created_at
		FROM action_history
		WHERE undone = 0
		ORDER BY id DESC LIMIT 1`)
	return scanActionRow(row)
}

// GetActionByID returns a single action by primary key.
func (d *DB) GetActionByID(id int64) (*ActionRecord, error) {
	row := d.QueryRow(`
		SELECT id, action_type, media_id, media_snapshot, detail, batch_id, undone, created_at
		FROM action_history WHERE id = ?`, id)
	return scanActionRow(row)
}

// GetLastRedoableAction returns the most recent undone action (for redo).
func (d *DB) GetLastRedoableAction() (*ActionRecord, error) {
	row := d.QueryRow(`
		SELECT id, action_type, media_id, media_snapshot, detail, batch_id, undone, created_at
		FROM action_history
		WHERE undone = 1
		ORDER BY id DESC LIMIT 1`)
	return scanActionRow(row)
}

// ListActions returns recent action_history records, newest first.
func (d *DB) ListActions(limit int) ([]ActionRecord, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := d.Query(`
		SELECT id, action_type, media_id, media_snapshot, detail, batch_id, undone, created_at
		FROM action_history
		ORDER BY id DESC LIMIT ?`, limit)
	if err != nil {
		return nil, fmt.Errorf("list actions: %w", err)
	}
	defer rows.Close()
	return scanActionRows(rows)
}

// ListActionsByType filters by action_type.
func (d *DB) ListActionsByType(actionType ActionType, limit int) ([]ActionRecord, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := d.Query(`
		SELECT id, action_type, media_id, media_snapshot, detail, batch_id, undone, created_at
		FROM action_history
		WHERE action_type = ?
		ORDER BY id DESC LIMIT ?`, string(actionType), limit)
	if err != nil {
		return nil, fmt.Errorf("list actions by type: %w", err)
	}
	defer rows.Close()
	return scanActionRows(rows)
}

// ListActionsByBatch returns all actions sharing a batch_id.
func (d *DB) ListActionsByBatch(batchID string) ([]ActionRecord, error) {
	rows, err := d.Query(`
		SELECT id, action_type, media_id, media_snapshot, detail, batch_id, undone, created_at
		FROM action_history
		WHERE batch_id = ?
		ORDER BY id ASC`, batchID)
	if err != nil {
		return nil, fmt.Errorf("list actions by batch: %w", err)
	}
	defer rows.Close()
	return scanActionRows(rows)
}

// ---------------------------------------------------------------------------
// Undo / Redo state helpers
// ---------------------------------------------------------------------------

// MarkActionUndone sets undone = 1 for the given action.
func (d *DB) MarkActionUndone(id int64) error {
	_, err := d.Exec("UPDATE action_history SET undone = 1 WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("mark action undone %d: %w", id, err)
	}
	return nil
}

// MarkActionRedone sets undone = 0 for the given action (redo).
func (d *DB) MarkActionRedone(id int64) error {
	_, err := d.Exec("UPDATE action_history SET undone = 0 WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("mark action redone %d: %w", id, err)
	}
	return nil
}

// MarkBatchUndone sets undone = 1 for all actions in a batch.
func (d *DB) MarkBatchUndone(batchID string) error {
	_, err := d.Exec("UPDATE action_history SET undone = 1 WHERE batch_id = ?", batchID)
	if err != nil {
		return fmt.Errorf("mark batch undone %s: %w", batchID, err)
	}
	return nil
}

// MarkBatchRedone sets undone = 0 for all actions in a batch.
func (d *DB) MarkBatchRedone(batchID string) error {
	_, err := d.Exec("UPDATE action_history SET undone = 0 WHERE batch_id = ?", batchID)
	if err != nil {
		return fmt.Errorf("mark batch redone %s: %w", batchID, err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Internal scan helpers
// ---------------------------------------------------------------------------

func scanActionRow(row *sql.Row) (*ActionRecord, error) {
	r := &ActionRecord{}
	err := row.Scan(&r.ID, &r.ActionType, &r.MediaID, &r.MediaSnapshot,
		&r.Detail, &r.BatchID, &r.Undone, &r.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("scan action row: %w", err)
	}
	return r, nil
}

func scanActionRows(rows *sql.Rows) ([]ActionRecord, error) {
	var records []ActionRecord
	for rows.Next() {
		var r ActionRecord
		if err := rows.Scan(&r.ID, &r.ActionType, &r.MediaID, &r.MediaSnapshot,
			&r.Detail, &r.BatchID, &r.Undone, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan action rows: %w", err)
		}
		records = append(records, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("action rows iteration: %w", err)
	}
	return records, nil
}
