# AI Handoff — Movie CLI

> **Last Updated**: 05-Apr-2026  
> **Purpose**: Single file containing ALL project context, memory, conventions, issues, and recent changes. Share this with any AI to get them fully up to speed.

---

## Table of Contents

1. [Project Identity](#1-project-identity)
2. [Critical Rules for AI](#2-critical-rules-for-ai)
3. [Architecture & File Structure](#3-architecture--file-structure)
4. [Command Tree & Behavior](#4-command-tree--behavior)
5. [Database Schema](#5-database-schema)
6. [Key Technical Decisions](#6-key-technical-decisions)
7. [Conventions & Code Style](#7-conventions--code-style)
8. [Build & Deploy](#8-build--deploy)
9. [Configuration Keys](#9-configuration-keys)
10. [Issues Found & Fixed](#10-issues-found--fixed)
11. [Pending Work](#11-pending-work)
12. [Full Conversation History](#12-full-conversation-history)
13. [Memory & Documentation System](#13-memory--documentation-system)
14. [AI Success Rate Plan](#14-ai-success-rate-plan)
15. [Edge Cases to Handle](#15-edge-cases-to-handle)

---

## 1. Project Identity

| Field | Value |
|-------|-------|
| **Name** | Movie CLI |
| **Type** | Go CLI application — **NOT a web app** |
| **Binary** | `movie-cli` |
| **Language** | Go 1.22 |
| **Module** | `github.com/mahin/movie-cli-v2` |
| **Framework** | Cobra (CLI), SQLite (storage), TMDb API (metadata) |
| **Purpose** | Cross-platform CLI tool for managing a personal movie and TV show library |
| **Development Period** | 15-Mar-2026 → 18-Mar-2026 (active), maintained since |

### What It Does

Scans local folders for video files, cleans messy filenames (removes `1080p`, `BluRay`, `x264`, etc.), fetches metadata from TMDb (The Movie Database), stores everything in SQLite, and organizes files into configured directories. Supports interactive browsing, searching, moving, renaming, undoing, and playing media files.

---

## 2. Critical Rules for AI

### ⛔ MUST NOT Do
1. **Never create web files** — no `package.json`, no `index.html`, no React/Vue/Angular components, no `npm install`
2. **Never treat Lovable build errors as real bugs** — errors like `no package.json found` and `Error: no command found for task "dev"` are **expected** and must be **ignored**
3. **Never use placeholder values** — no `"now"`, `"TODO"`, `"test"` in production code paths
4. **Never copy-paste code that exists elsewhere** — use shared helpers
5. **Never let a file exceed ~200 lines** — split at natural boundaries
6. **Never create files for wrong technologies** — if a request mentions PHP, WordPress, TypeScript, React, Chrome extensions, etc., STOP and ask the user to verify they're in the right project

### ✅ MUST Do
1. **Read this file first** before making any changes
2. **Follow one-file-per-command pattern** — `cmd/movie_<cmd>.go`
3. **Use shared helpers** — `fetchMovieDetails()` / `fetchTVDetails()` from `movie_info.go`, `resolveMediaByQuery()` from `movie_resolve.go`
4. **Use real timestamps** — `time.Now().Format(time.RFC3339)`
5. **Document every fix** — create issue write-up, update spec, update memory (3-step mandatory process)
6. **Test locally** — all CLI commands need a real terminal, not a browser preview

---

## 3. Architecture & File Structure

```
movie-cli/
├── main.go                          # Entry point — calls cmd.Execute()
├── cmd/
│   ├── root.go                      # Root Cobra command
│   ├── hello.go                     # Greeting with version
│   ├── version.go                   # Version/commit/build date (ldflags)
│   ├── update.go                    # Self-update via git pull --ff-only
│   ├── movie.go                     # Parent: movie-cli movie (registers all subcommands)
│   ├── movie_config.go              # config get/set — manages settings in SQLite
│   ├── movie_scan.go                # scan <folder> — walk dir, clean names, fetch TMDb, save to DB
│   ├── movie_ls.go                  # ls — paginated library listing (20 items/page)
│   ├── movie_search.go              # search <name> — TMDb API search, interactive select, save to DB
│   ├── movie_info.go                # info <id|title> — DB-first lookup, TMDb fallback, auto-persist
│   │                                # ⭐ ALSO contains shared fetchMovieDetails() / fetchTVDetails()
│   ├── movie_suggest.go             # suggest [N] — genre-based recommendations + trending
│   ├── movie_move.go                # move — interactive browse-first file move with undo tracking
│   ├── movie_move_helpers.go        # Move utility functions (promptSourceDirectory, listVideoFiles, etc.)
│   ├── movie_rename.go              # rename — batch clean rename to "Title (Year).ext"
│   ├── movie_undo.go                # undo — revert last move via move_history table
│   ├── movie_play.go                # play <id> — open file with system default player
│   ├── movie_stats.go               # stats — genre chart, counts, average ratings
│   └── movie_resolve.go             # ⭐ Shared ID/title resolver helper (resolveMediaByQuery)
├── cleaner/
│   └── cleaner.go                   # Regex-based filename cleaning + slug generation (~120 lines)
├── tmdb/
│   └── client.go                    # TMDb API client: search, details, credits, discover, trending, posters (~200 lines)
├── db/
│   ├── db.go                        # DB struct, Open(), migrate() — connection + schema (~140 lines)
│   ├── media.go                     # Media struct, all CRUD: Insert, Update, Get, List, Search, TopGenres
│   ├── config.go                    # GetConfig, SetConfig
│   ├── history.go                   # MoveRecord, InsertMoveHistory, GetLastMove, MarkMoveUndone, InsertScanHistory
│   └── helpers.go                   # splitCSV — uses strings.Split + strings.TrimSpace (13 lines)
├── updater/
│   └── updater.go                   # Git-based self-update (git pull --ff-only)
├── version/
│   └── version.go                   # Build-time variables (Version, Commit, BuildDate) injected via ldflags
├── Makefile                         # Build targets: build, build-windows, build-mac-arm, clean, install
├── build.ps1                        # PowerShell build + deploy (Windows: E:\bin-run, Mac: /usr/local/bin)
├── spec.md                          # Full project specification (350+ lines)
├── README.md                        # Project documentation
├── development-log.md               # Every prompt + modification recorded
├── readm.txt                        # Session milestone marker
├── go.mod / go.sum                  # Go module files
├── spec/
│   ├── 01-app/
│   │   └── 01-project-spec.md       # Application-level spec with coding guidelines
│   └── 02-app/
│       └── issues/
│           ├── 00-issue-template.md  # Template for new issue write-ups
│           ├── 01-hardcoded-timestamp.md
│           ├── 02-duplicate-tmdb-fetch.md
│           ├── 03-large-file-refactor.md
│           ├── 04-wrong-project-context.md
│           └── 05-code-hygiene-audit.md
└── .lovable/memory/
    ├── 01-project-overview.md       # Project context, architecture, command tree
    ├── 02-conventions.md            # Naming, style, build, deploy preferences
    ├── workflow/
    │   ├── 01-plan.md               # ✅ 20+ completed, 🔲 8 pending items
    │   └── 01-ai-success-plan.md    # Rules for 98% AI success rate
    ├── suggestions/
    │   └── 01-suggestions.md        # 10 suggestions, 3 completed, 7 pending
    └── issues/
        ├── 01-timestamp-bug.md
        ├── 02-duplicate-tmdb-fetch.md
        └── 03-large-files.md
```

### File Count Summary
- `cmd/` — 15 Go files (root + hello + version + update + movie parent + 10 subcommands + move_helpers + resolve)
- `cleaner/` — 1 file
- `tmdb/` — 1 file
- `db/` — 5 files (db.go, media.go, config.go, history.go, helpers.go)
- `updater/` — 1 file
- `version/` — 1 file

---

## 4. Command Tree & Behavior

```
movie-cli
├── hello                      # Prints greeting with version number
├── version                    # Shows Version, Commit, BuildDate (injected via ldflags)
├── self-update                # Runs git pull --ff-only in the binary's source directory
└── movie
    ├── config                 # View/set configuration
    │   ├── config get <key>   # Read a config value
    │   └── config set <k> <v> # Write a config value (API key display is masked)
    ├── scan <folder>          # Scan folder for video files:
    │                          #   1. Walk directory tree
    │                          #   2. Clean filenames (regex: remove 1080p, BluRay, x264, etc.)
    │                          #   3. Detect Movie vs TV Show (SxxExx pattern)
    │                          #   4. Fetch metadata from TMDb API
    │                          #   5. Download poster/thumbnail
    │                          #   6. Save to SQLite database
    │                          #   7. Log scan to scan_history table
    ├── ls                     # List library (from DB, populated by scan):
    │                          #   - 20 items per page (configurable via page_size)
    │                          #   - Shows: #, 🎬/📺 icon, clean title, year, ⭐ rating
    │                          #   - Interactive: type number for detail, 'n' for next page, 'q' to quit
    ├── search <name>          # Search TMDb API directly (NOT local DB):
    │                          #   1. Query TMDb multi-search
    │                          #   2. Show up to 15 results with type, rating, year
    │                          #   3. User picks by number
    │                          #   4. Fetch full details (cast, director, genres, description)
    │                          #   5. Save to local database
    │                          #   6. Media doesn't need to exist as a local file
    ├── info <id|title>        # Detail view with priority logic:
    │                          #   - Numeric ID → direct DB lookup
    │                          #   - Title → search local DB first (exact + prefix match)
    │                          #   - Title not found locally → TMDb API fallback → auto-save to DB
    │                          #   - Uses shared resolveMediaByQuery() from movie_resolve.go
    ├── suggest [N]            # Recommendations:
    │                          #   - Analyzes top genres from library
    │                          #   - Fetches TMDb discover results for those genres
    │                          #   - Falls back to TMDb trending if library is empty
    │                          #   - N = number of suggestions (default 5)
    ├── move                   # Interactive browse-first file move:
    │                          #   1. Choose source: Downloads / scan_dir / custom / CLI arg
    │                          #   2. List all video files in source with cleaned titles + size
    │                          #   3. Select file by number
    │                          #   4. Choose destination: Movies dir / TV dir / Archive / custom
    │                          #   5. Confirm with [y/N] prompt
    │                          #   6. Move file + track in move_history table + JSON log
    │                          #   7. Auto-insert to DB if not already tracked
    ├── rename                 # Batch rename to "Title (Year).ext":
    │                          #   - Uses cleaner to generate clean names
    │                          #   - Tracks original names for undo
    ├── undo                   # Revert last move/rename:
    │                          #   - Reads latest non-undone entry from move_history
    │                          #   - Moves file back to original location
    │                          #   - Marks entry as undone in DB
    │                          #   ⚠️ No confirmation prompt (pending enhancement)
    ├── play <id>              # Open file with system default player:
    │                          #   - macOS: open
    │                          #   - Windows: start
    │                          #   - Linux: xdg-open
    └── stats                  # Library statistics:
                               #   - Total movies vs TV shows
                               #   - Average TMDb/IMDb ratings
                               #   - Top genres with bar chart
                               #   - Recently added items
```

---

## 5. Database Schema

Database location: `./data/movie-cli.db` (relative to executable)  
Engine: SQLite via `modernc.org/sqlite` (pure Go, no CGo)  
Journal mode: WAL (Write-Ahead Logging for concurrency)

### Table: `media` — Core metadata
```sql
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
    genre            TEXT,               -- comma-separated genre names
    director         TEXT,
    cast_list        TEXT,               -- comma-separated actor names
    thumbnail_path   TEXT,               -- local path to downloaded poster
    original_file_name TEXT,
    original_file_path TEXT,
    current_file_path  TEXT,
    file_extension   TEXT,
    file_size        INTEGER,
    scanned_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at       DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Table: `move_history` — Undo support
```sql
CREATE TABLE IF NOT EXISTS move_history (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    media_id         INTEGER NOT NULL,
    from_path        TEXT NOT NULL,
    to_path          TEXT NOT NULL,
    original_file_name TEXT,
    new_file_name    TEXT,
    moved_at         DATETIME DEFAULT CURRENT_TIMESTAMP,
    undone           INTEGER DEFAULT 0,  -- 0=active, 1=undone
    FOREIGN KEY (media_id) REFERENCES media(id) ON DELETE CASCADE
);
```

### Table: `config` — Key-value settings
```sql
CREATE TABLE IF NOT EXISTS config (
    key   TEXT PRIMARY KEY,
    value TEXT NOT NULL
);
-- Default values inserted on migration:
-- movies_dir=~/Movies, tv_dir=~/TVShows, archive_dir=~/Archive,
-- scan_dir=~/Downloads, page_size=20
```

### Table: `scan_history` — Scan operation logs
```sql
CREATE TABLE IF NOT EXISTS scan_history (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    folder_path   TEXT NOT NULL,
    total_files   INTEGER DEFAULT 0,
    movies_found  INTEGER DEFAULT 0,
    tv_found      INTEGER DEFAULT 0,
    scanned_at    DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Table: `tags` — User-defined tags (⚠️ table exists, commands NOT YET implemented)
```sql
CREATE TABLE IF NOT EXISTS tags (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    media_id   INTEGER NOT NULL,
    tag        TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (media_id) REFERENCES media(id) ON DELETE CASCADE,
    UNIQUE(media_id, tag)
);
```

### Indexes
```sql
CREATE INDEX IF NOT EXISTS idx_media_type         ON media(type);
CREATE INDEX IF NOT EXISTS idx_media_title        ON media(clean_title);
CREATE INDEX IF NOT EXISTS idx_media_year         ON media(year);
CREATE INDEX IF NOT EXISTS idx_media_tmdb         ON media(tmdb_id);
CREATE INDEX IF NOT EXISTS idx_move_history_media  ON move_history(media_id);
CREATE INDEX IF NOT EXISTS idx_move_history_undone ON move_history(undone);
CREATE INDEX IF NOT EXISTS idx_tags_media         ON tags(media_id);
```

### Data Storage Layout
```
./data/
├── movie-cli.db                    # SQLite database
├── thumbnails/                 # Downloaded poster images
└── json/
    ├── movie/                  # (Directories exist, JSON metadata not yet written)
    ├── tv/
    └── history/                # Move operation JSON logs
```

---

## 6. Key Technical Decisions

| Decision | Rationale |
|----------|-----------|
| `modernc.org/sqlite` (pure Go) | No CGo dependency → easy cross-compilation for Windows/Mac/Linux |
| WAL journal mode | Better concurrent read/write performance |
| Cobra framework | Industry-standard Go CLI framework, supports subcommands |
| TMDb API (not IMDb) | Free API with rich metadata, poster images, recommendations |
| `./data/` data dir | Single location for all data, easy to backup/migrate |
| Git-based self-update | Simple `git pull --ff-only` — assumes binary is built from source |
| One file per command | Keeps files small, easy for AI and humans to navigate |
| Shared TMDb helpers | `fetchMovieDetails()` / `fetchTVDetails()` in `movie_info.go` — single source of truth |
| Shared resolver | `resolveMediaByQuery()` in `movie_resolve.go` — consistent ID/title lookup |

---

## 7. Conventions & Code Style

### Go Code
- Standard Go formatting (`gofmt`)
- Max ~200 lines per file — split at natural boundaries
- One file per Cobra command: `cmd/movie_<command>.go`
- Helper functions go in `_helpers.go` suffix files
- Explicit methods over boolean flags (single responsibility)
- DRY: use shared helpers, never copy-paste

### CLI Output Emojis
- ✅ Success messages
- ❌ Error messages
- 🎬 Movie type indicator
- 📺 TV show type indicator
- ⭐ Rating display
- 📭 Empty/no results

### Documentation
- Project milestones: `readm.txt` — `let's start now {dd-MMM-YYYY} {hh:mm AM/PM}` (Malaysia time UTC+8)
- Specifications: `spec.md` (root) + `spec/01-app/` (app spec) + `spec/01-app/03-issues/` (issue write-ups)
- Memory: `.lovable/memory/` — project overview, conventions, workflow, suggestions, issues

### Issue Tracking (Mandatory 3-Step Process)
Every fix MUST include:
1. Create `spec/01-app/03-issues/XX-{slug}.md` using the template
2. Update the relevant spec in `spec/01-app/`
3. Update `.lovable/memory/` with summary and prevention rule

---

## 8. Build & Deploy

### Makefile Targets
```makefile
make build           # Build for current OS → ./movie-cli
make build-windows   # Cross-compile → movie-cli.exe (amd64)
make build-mac-arm   # Cross-compile → movie-cli-darwin-arm64
make build-mac-intel # Cross-compile → movie-cli-darwin-amd64
make build-linux     # Cross-compile → movie-cli-linux-amd64
make clean           # Remove build artifacts
make install         # Build + copy to /usr/local/bin
```

### PowerShell Deploy (`build.ps1`)
```powershell
pwsh build.ps1                        # Default: E:\bin-run (Windows) or /usr/local/bin (Mac/Linux)
pwsh build.ps1 -BinDir "C:\custom"    # Custom deploy path
```
Steps: `git pull` → `go mod tidy` → build with `-ldflags="-s -w"` → deploy → verify

### Version Injection
```bash
go build -ldflags="-X version.Version=1.0.0 -X version.Commit=$(git rev-parse --short HEAD) -X version.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o movie-cli .
```

### Dependencies (`go.mod`)
| Package | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI framework |
| `modernc.org/sqlite` | Pure-Go SQLite driver |

---

## 9. Configuration Keys

| Key | Default | Purpose |
|-----|---------|---------|
| `movies_dir` | `~/Movies` | Destination for movie files |
| `tv_dir` | `~/TVShows` | Destination for TV show files |
| `archive_dir` | `~/Archive` | Archive destination |
| `scan_dir` | `~/Downloads` | Default scan source directory |
| `tmdb_api_key` | (none) | TMDb API key — required for metadata features |
| `page_size` | `20` | Items per page in `movie ls` |

Set via: `movie-cli movie config set <key> <value>`  
Get via: `movie-cli movie config get <key>`

---

## 10. Issues Found & Fixed

### Issue #01: Hardcoded Timestamp in move-log.json ✅ RESOLVED
- **Severity**: Medium
- **File**: `cmd/movie_move_helpers.go` (was `cmd/movie_move.go` line 345-346)
- **Root Cause**: `saveHistoryLog` wrote `"timestamp":"now"` as literal string instead of actual time
- **Fix**: Replaced with `time.Now().Format(time.RFC3339)`
- **Impact**: All move history JSON logs now have proper timestamps
- **Prevention**: Never use placeholder strings in production code; grep for `"now"`, `"TODO"`, `"test"` before committing

### Issue #02: Duplicate TMDb Fetch Logic ✅ RESOLVED
- **Severity**: Low (code quality)
- **Files**: `cmd/movie_scan.go`, `cmd/movie_search.go`, `cmd/movie_info.go`
- **Root Cause**: Three commands had ~80 lines of nearly identical TMDb fetch code
- **Fix**: Refactored `scan` and `search` to call shared `fetchMovieDetails()` / `fetchTVDetails()` from `movie_info.go`
- **Prevention**: Always use shared helpers; never copy-paste TMDb fetch blocks

### Issue #03: Large Files Need Refactoring ✅ RESOLVED
- **Severity**: Low (maintainability)
- **Files**: `cmd/movie_move.go` (348 lines), `db/sqlite.go` (452 lines)
- **Fix**:
  - `movie_move.go` → `movie_move.go` (178 lines) + `movie_move_helpers.go` (168 lines)
  - `db/sqlite.go` → 5 files: `db.go`, `media.go`, `config.go`, `history.go`, `helpers.go`
- **Prevention**: Split files at natural boundaries; never exceed ~200 lines

### Issue #04: Wrong Project Context ✅ RESOLVED
- **Severity**: Critical (process)
- **Root Cause**: User had multiple Lovable projects open and sent a Chrome extension request to this Go CLI project
- **Fix**: Added scope validation rule — AI must verify incoming requests match Go CLI before proceeding
- **Prevention**: If request mentions PHP, WordPress, TypeScript, React, Chrome → stop and ask

### Issue #05: Code Hygiene — 8 Quality Issues ✅ RESOLVED
All 8 issues found during full codebase audit and fixed:

| # | Issue | Fix |
|---|-------|-----|
| 1 | `go.mod` missing SQLite dep | Added `modernc.org/sqlite v1.29.5` |
| 2 | `db/helpers.go` reimplements stdlib (46 lines) | Replaced with `strings.Split`/`strings.TrimSpace` (13 lines) |
| 3 | Stray `//Ab` comment in `main.go` | Removed |
| 4 | Broken indentation in `cmd/movie_scan.go` | Fixed alignment |
| 5 | Redundant `Media` field copy in `cmd/movie_info.go` | Removed; passes `m` directly |
| 6 | Unused `resolveMediaByQuery` function | Wired into `movie_info.go` for title lookups |
| 7 | Duplicate `readm.txt` + `readme.txt` | Deleted `readme.txt`; `readm.txt` is sole milestone file |
| 8 | Compiled `movie-cli` binary committed to repo | Deleted binaries from repo |

---

## 11. Pending Work

### 🔴 Known Bugs
- No confirmation prompt on `movie undo` before reverting (currently undoes immediately)

### 🟡 Missing Features
| Feature | Notes |
|---------|-------|
| `movie tag` command | `tags` table exists in DB but no commands use it. Needs: `tag add <id> <tag>`, `tag remove <id> <tag>`, `tag list [id]` |
| `DiscoverByGenre` TMDb method | Method exists in `tmdb/client.go` but unused by any command |
| Per-media JSON metadata files | Directories exist (`json/movie/`, `json/tv/`) but files are never written |
| `.gitignore` file | Must be created manually (Lovable environment limitation) |

### 🟢 Enhancements
| Enhancement | Notes |
|-------------|-------|
| Cross-drive move support | `os.Rename` fails across filesystems → need copy+delete fallback |
| File size statistics | Add total/average/largest file size to `movie stats` |
| Batch move (`--all` flag) | Move all video files from source at once |
| TMDb rate limit handling | Retry with backoff or show warning |
| Network offline graceful handling | TMDb calls fail gracefully, scan still works without metadata |

### Proposed `.gitignore` Content
```gitignore
movie-cli
movie-cli.exe
movie-cli-darwin-*
movie-cli
*.exe
*.o
*.a
*.so
*.dylib
data/
.DS_Store
Thumbs.db
.idea/
.vscode/
```

---

## 12. Full Conversation History

### Session 1: Project Bootstrap (15-Mar-2026)
- Created `readme.txt` with milestone marker
- Corrected time to Malaysia timezone (UTC+8)
- Established `readm.txt` convention for session milestones

### Session 2: Movie CLI Specification (16-Mar-2026)
- User described the full movie CLI vision: scan folders, clean names, fetch metadata, manage library
- Planned command tree: `scan`, `ls`, `search`, `info`, `suggest`, `move`, `rename`, `undo`, `play`, `stats`, `config`
- Decided on SQLite (not JSON) for storage
- Designed interactive `movie move` with browse-first workflow
- Designed 5-table database schema

### Session 3: Full Implementation (16-Mar-2026)
- Implemented ALL 16 Go files for the movie CLI
- Key design: `movie search` queries TMDb directly (not local DB), saves results to DB
- Key design: `movie info` uses DB-first → TMDb fallback → auto-persist
- Key design: `movie ls` shows only items indexed via `movie scan`
- Rewrote `movie_search.go` based on user clarification (search TMDb, not local)
- Rewrote `movie_info.go` with ID vs title priority logic

### Session 4: Command Refinements & Build Tooling (17-Mar-2026)
- Rewrote `movie move` to browse-first workflow (user picks directory → sees files → selects → confirms)
- Created `Makefile` with cross-compilation targets
- Created `build.ps1` PowerShell deploy script
- Attempted `.gitignore` (blocked by Lovable environment)

### Session 5: Specification & Quality Audit (17-Mar-2026)
- Created comprehensive `spec.md` (350+ lines)
- Set up `.lovable/memory/` tracking system (7 files)
- Fixed 3 code issues: timestamp bug, TMDb deduplication, file splitting
- Created AI success rate plan targeting 98%
- Performed requirements reliability analysis
- Updated stale references in `spec.md`

### Session 6: Code Hygiene (17-18-Mar-2026)
- Full codebase audit — confirmed project is focused, no scope creep
- Found and fixed 8 hygiene issues (see Issue #05 above)
- Updated `README.md` with comprehensive documentation
- Created `development-log.md` recording every prompt and modification

### Session 7: Documentation & Handoff (18-Mar-2026 → 05-Apr-2026)
- Updated `readm.txt` with new milestone markers
- Updated `README.md` based on project state
- Created this `ai-handoff.md` file

---

## 13. Memory & Documentation System

The project uses a structured memory system in `.lovable/memory/`:

| File | Purpose |
|------|---------|
| `01-project-overview.md` | Project identity, architecture, command tree, important notes for AI |
| `02-conventions.md` | Documentation format, file naming, code style, build/deploy, config keys |
| `workflow/01-plan.md` | Checklist of completed (20+) and pending (8) items |
| `workflow/01-ai-success-plan.md` | 7 rules for achieving 98% AI success rate |
| `suggestions/01-suggestions.md` | 10 suggestions: 3 completed, 7 pending with priority levels |
| `issues/01-timestamp-bug.md` | Root cause analysis + fix for hardcoded timestamp |
| `issues/02-duplicate-tmdb-fetch.md` | Root cause analysis + fix for duplicate TMDb logic |
| `issues/03-large-files.md` | Root cause analysis + fix for oversized files |

Additional documentation:
| File | Purpose |
|------|---------|
| `spec.md` | Master technical specification (350+ lines) |
| `README.md` | User-facing project documentation |
| `development-log.md` | Every prompt and modification recorded |
| `spec/01-app/01-project-spec.md` | Application-level spec with coding guidelines |
| `spec/01-app/03-issues/00-issue-template.md` | Template for new issue write-ups |
| `spec/01-app/03-issues/01-05` | Individual issue write-ups |

---

## 14. AI Success Rate Plan

### Why AI Fails on This Project (and How We Fixed It)

| Cause | % of Failures | Fix Applied |
|-------|---------------|-------------|
| Missing context (treats as web app) | ~40% | Project overview says "NOT a web app"; ignore Lovable build errors |
| Stale memory / inconsistent state | ~20% | Single source of truth per concern; mandatory updates after changes |
| Large files confuse AI | ~15% | All files now <200 lines; split db/sqlite.go (452→5 files) |
| Duplicate logic → divergent behavior | ~10% | Extracted shared TMDb fetch helpers |
| Placeholder/debug values in production | ~5% | Fixed timestamp bug; rule: grep for placeholders |
| No edge case specs | ~10% | Added edge case documentation |

### 7 Rules for 98% Success

1. **Read memory before coding** — `01-project-overview.md` → `02-conventions.md` → `workflow/01-plan.md`
2. **One file, one concern, <200 lines**
3. **Shared logic = shared functions** — never duplicate TMDb fetch or resolver code
4. **No placeholder values** — always use real data and `time.Now().Format(time.RFC3339)`
5. **Update memory after every change** — plan, suggestions, and issues
6. **Ignore Lovable build errors** — `no package.json found` is expected
7. **Test locally, not in Lovable preview** — CLI needs real terminal + filesystem

### Debugging Strategy
1. Is it a Go compile error or Lovable platform error? If platform → ignore
2. Check imports after splitting files
3. Check function signatures after refactoring
4. Check package boundaries (exported vs unexported)
5. Verify no duplicate symbols after file splits
6. Mentally trace `go build ./...` through all files

---

## 15. Edge Cases to Handle

These are known edge cases that are **not yet implemented** but should be considered:

| Edge Case | Current State | Recommended Fix |
|-----------|---------------|-----------------|
| `os.Rename` across filesystems | Fails silently | Detect error, fallback to io.Copy + os.Remove |
| Empty scan folder | May show confusing output | Show "📭 No video files found in <path>" |
| TMDb API rate limiting | No handling | Retry with exponential backoff or show warning |
| Duplicate TMDb IDs | UNIQUE constraint throws error | Use `INSERT OR REPLACE` or `UPDATE` fallback |
| Unicode filenames | Untested | Verify `cleaner.Clean` handles non-ASCII gracefully |
| Network offline | TMDb calls fail | Fail gracefully; scan still saves cleaned filenames without metadata |
| Database locked | WAL mode helps | Handle `SQLITE_BUSY` with retry |
| Very large libraries (10k+) | Pagination exists | May need query optimization / indexing review |

---

## Appendix: Quick Start for New AI

```bash
# 1. Understand the project
# Read this file (ai-handoff.md) — you're already here ✅

# 2. Build the binary
make build
# OR
go build -o movie-cli .

# 3. Configure
./movie-cli movie config set tmdb_api_key YOUR_API_KEY
./movie-cli movie config set movies_dir ~/Movies
./movie-cli movie config set scan_dir ~/Downloads

# 4. Use it
./movie-cli movie scan ~/Downloads      # Scan for videos
./movie-cli movie ls                    # List library
./movie-cli movie search "Inception"    # Search TMDb
./movie-cli movie info 1                # View details by ID
./movie-cli movie info "The Matrix"     # View details by title
./movie-cli movie move                  # Interactive file move
./movie-cli movie suggest 5             # Get 5 recommendations
./movie-cli movie stats                 # Library statistics
```

---

*This file is the single source of truth for any AI continuing work on Movie CLI. If you've read this, you're up to speed.*
