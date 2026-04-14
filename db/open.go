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

// exeDir returns the directory where the running binary is located.
// Symlinks are resolved so the real physical location is used.
func exeDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("cannot locate executable: %w", err)
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return "", fmt.Errorf("cannot resolve symlinks for executable: %w", err)
	}
	return filepath.Dir(exe), nil
}

// Open opens (or creates) the SQLite database and runs migrations.
// The database is stored in <binary-dir>/data/movie.db — always relative
// to where the CLI binary physically resides, not the working directory.
func Open() (*DB, error) {
	binDir, dirErr := exeDir()
	if dirErr != nil {
		return nil, dirErr
	}

	base := filepath.Join(binDir, "data")
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
		runtime          INTEGER DEFAULT 0,
		language         TEXT DEFAULT '',
		budget           INTEGER DEFAULT 0,
		revenue          INTEGER DEFAULT 0,
		trailer_url      TEXT DEFAULT '',
		tagline          TEXT DEFAULT '',
		scanned_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at       DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Migration: add new columns if upgrading from older schema
	-- SQLite ignores ALTER TABLE ADD COLUMN if column already exists via IF NOT EXISTS workaround
	`

	// Run the main schema
	if _, err := d.Exec(schema); err != nil {
		return err
	}

	// Add columns if they don't exist (for existing databases)
	newCols := []struct{ name, def string }{
		{"runtime", "INTEGER DEFAULT 0"},
		{"language", "TEXT DEFAULT ''"},
		{"budget", "INTEGER DEFAULT 0"},
		{"revenue", "INTEGER DEFAULT 0"},
		{"trailer_url", "TEXT DEFAULT ''"},
		{"tagline", "TEXT DEFAULT ''"},
	}
	for _, col := range newCols {
		q := fmt.Sprintf("ALTER TABLE media ADD COLUMN %s %s", col.name, col.def)
		d.Exec(q) // ignore error = column already exists
	}

	rest := `

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

	CREATE TABLE IF NOT EXISTS watchlist (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		media_id   INTEGER,
		tmdb_id    INTEGER NOT NULL,
		title      TEXT NOT NULL,
		year       INTEGER,
		type       TEXT CHECK(type IN ('movie', 'tv')) NOT NULL,
		status     TEXT CHECK(status IN ('to-watch', 'watched')) NOT NULL DEFAULT 'to-watch',
		added_at   DATETIME DEFAULT CURRENT_TIMESTAMP,
		watched_at DATETIME,
		FOREIGN KEY (media_id) REFERENCES media(id) ON DELETE SET NULL,
		UNIQUE(tmdb_id)
	);

	CREATE INDEX IF NOT EXISTS idx_watchlist_status ON watchlist(status);
	CREATE INDEX IF NOT EXISTS idx_watchlist_tmdb   ON watchlist(tmdb_id);

	INSERT OR IGNORE INTO config (key, value) VALUES ('movies_dir',  '~/Movies');
	INSERT OR IGNORE INTO config (key, value) VALUES ('tv_dir',      '~/TVShows');
	INSERT OR IGNORE INTO config (key, value) VALUES ('archive_dir', '~/Archive');
	INSERT OR IGNORE INTO config (key, value) VALUES ('scan_dir',    '~/Downloads');
	INSERT OR IGNORE INTO config (key, value) VALUES ('page_size',   '20');
	`
	_, err := d.Exec(rest)
	return err
}
