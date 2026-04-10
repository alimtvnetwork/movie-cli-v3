// Package db manages the SQLite database for the movie CLI.
package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const dbFile = "movie.db"

// DB wraps the sql.DB connection.
type DB struct {
	*sql.DB
	BasePath string // path to data directory
}

// Open opens (or creates) the SQLite database and runs migrations.
// The database is stored in ./data/movie.db relative to the working directory.
func Open() (*DB, error) {
	base := filepath.Join(".", "data")
	dirs := []string{
		base,
		filepath.Join(base, "json", "movie"),
		filepath.Join(base, "json", "tv"),
		filepath.Join(base, "json", "history"),
		filepath.Join(base, "thumbnails"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return nil, fmt.Errorf("cannot create directory %s: %w", d, err)
		}
	}

	dbPath := filepath.Join(base, dbFile)
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open database: %w", err)
	}

	// Enable WAL mode for better concurrency
	if _, err := conn.Exec("PRAGMA journal_mode=WAL"); err != nil {
		conn.Close()
		return nil, fmt.Errorf("cannot set WAL mode: %w", err)
	}

	d := &DB{DB: conn, BasePath: base}
	if err := d.migrate(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return d, nil
}

func (d *DB) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS media (
		id               INTEGER PRIMARY KEY AUTOINCREMENT,
		title            TEXT NOT NULL,
		clean_title      TEXT NOT NULL,
		year             INTEGER,
		type             TEXT CHECK(type IN ('movie', 'tv')) NOT NULL,
		tmdb_id          INTEGER UNIQUE,
		imdb_id          TEXT,
		description      TEXT,
		imdb_rating      REAL,
		tmdb_rating      REAL,
		popularity       REAL,
		genre            TEXT,
		director         TEXT,
		cast_list        TEXT,
		thumbnail_path   TEXT,
		original_file_name TEXT,
		original_file_path TEXT,
		current_file_path  TEXT,
		file_extension   TEXT,
		file_size        INTEGER,
		scanned_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at       DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS move_history (
		id               INTEGER PRIMARY KEY AUTOINCREMENT,
		media_id         INTEGER NOT NULL,
		from_path        TEXT NOT NULL,
		to_path          TEXT NOT NULL,
		original_file_name TEXT,
		new_file_name    TEXT,
		moved_at         DATETIME DEFAULT CURRENT_TIMESTAMP,
		undone           INTEGER DEFAULT 0,
		FOREIGN KEY (media_id) REFERENCES media(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS config (
		key   TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS scan_history (
		id            INTEGER PRIMARY KEY AUTOINCREMENT,
		folder_path   TEXT NOT NULL,
		total_files   INTEGER DEFAULT 0,
		movies_found  INTEGER DEFAULT 0,
		tv_found      INTEGER DEFAULT 0,
		scanned_at    DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS tags (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		media_id   INTEGER NOT NULL,
		tag        TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (media_id) REFERENCES media(id) ON DELETE CASCADE,
		UNIQUE(media_id, tag)
	);

	CREATE INDEX IF NOT EXISTS idx_media_type       ON media(type);
	CREATE INDEX IF NOT EXISTS idx_media_title      ON media(clean_title);
	CREATE INDEX IF NOT EXISTS idx_media_year       ON media(year);
	CREATE INDEX IF NOT EXISTS idx_media_tmdb       ON media(tmdb_id);
	CREATE INDEX IF NOT EXISTS idx_move_history_media ON move_history(media_id);
	CREATE INDEX IF NOT EXISTS idx_move_history_undone ON move_history(undone);
	CREATE INDEX IF NOT EXISTS idx_tags_media       ON tags(media_id);

	INSERT OR IGNORE INTO config (key, value) VALUES ('movies_dir',  '~/Movies');
	INSERT OR IGNORE INTO config (key, value) VALUES ('tv_dir',      '~/TVShows');
	INSERT OR IGNORE INTO config (key, value) VALUES ('archive_dir', '~/Archive');
	INSERT OR IGNORE INTO config (key, value) VALUES ('scan_dir',    '~/Downloads');
	INSERT OR IGNORE INTO config (key, value) VALUES ('page_size',   '20');
	`
	_, err := d.Exec(schema)
	return err
}
