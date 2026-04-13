# Movie CLI — Full Specification

> **Version**: 1.0  
> **Date**: 17-Mar-2026  
> **Binary Name**: `movie-cli`  
> **Language**: Go 1.22  
> **Module**: `github.com/movie/movie-cli`

---

## Table of Contents

1. [Overview](#1-overview)
2. [Architecture](#2-architecture)
3. [Data Storage](#3-data-storage)
4. [Commands](#4-commands)
5. [Filename Cleaner](#5-filename-cleaner)
6. [TMDb Integration](#6-tmdb-integration)
7. [Build & Deploy](#7-build--deploy)
8. [Configuration Keys](#8-configuration-keys)
9. [Dependencies](#9-dependencies)
10. [AI Implementation Risk Assessment](#10-ai-implementation-risk-assessment)

---

## 1. Overview

`movie-cli` is a cross-platform CLI tool for managing a personal movie and TV show library. It:

- **Scans** local folders for video files
- **Cleans** messy filenames (removes quality tags, release groups, etc.)
- **Fetches metadata** from The Movie Database (TMDb) API
- **Stores** everything in a local SQLite database
- **Moves**, **renames**, and **organizes** files into configured directories
- **Suggests** new content based on library genre frequency or trending
- **Plays** media files using the system's default player
- Supports **undo** for all move/rename operations
- **Self-updates** via `git pull --ff-only`

---

## 2. Architecture

### Project Structure

```
movie-cli/
├── main.go                    # Entry point → cmd.Execute()
├── cmd/
│   ├── root.go                # Root cobra command, registers all subcommands
│   ├── hello.go               # movie-cli hello
│   ├── version.go             # movie-cli version
│   ├── update.go              # movie-cli self-update
│   ├── movie.go               # movie-cli movie (parent, registers subcommands)
│   ├── movie_config.go        # movie-cli movie config [get|set]
│   ├── movie_scan.go          # movie-cli movie scan [folder]
│   ├── movie_ls.go            # movie-cli movie ls
│   ├── movie_search.go        # movie-cli movie search <name>
│   ├── movie_info.go          # movie-cli movie info <id|title> + shared fetchMovieDetails/fetchTVDetails
│   ├── movie_suggest.go       # movie-cli movie suggest [N]
│   ├── movie_move.go          # movie-cli movie move [directory] (main flow)
│   ├── movie_move_helpers.go  # Move helpers: promptSourceDirectory, promptDestination, listVideoFiles, humanSize, expandHome, saveHistoryLog
│   ├── movie_rename.go        # movie-cli movie rename
│   ├── movie_undo.go          # movie-cli movie undo
│   ├── movie_play.go          # movie-cli movie play <id>
│   ├── movie_stats.go         # movie-cli movie stats
│   └── movie_resolve.go       # Shared helper: resolve media by ID or title
├── cleaner/
│   └── cleaner.go             # Filename cleaning, slug generation, video detection
├── tmdb/
│   └── client.go              # TMDb API client (search, details, credits, posters, trending)
├── db/
│   ├── db.go                  # DB struct, Open(), migrate() (schema + defaults)
│   ├── media.go               # Media struct, all CRUD methods, scanMediaRows, TopGenres
│   ├── config.go              # GetConfig, SetConfig
│   ├── history.go             # MoveRecord, InsertMoveHistory, GetLastMove, MarkMoveUndone, InsertScanHistory
│   └── helpers.go             # splitCSV, split, indexOf, trim
├── updater/
│   └── updater.go             # git-based self-update logic
├── version/
│   └── version.go             # Build-time version variables (ldflags)
├── spec/
│   ├── 08-app/                # Application specs and coding guidelines
│   └── 02-app/issues/         # Issue write-ups (root cause, fix, prevention)
├── Makefile                   # Build targets
├── build.ps1                  # PowerShell build + deploy script
├── go.mod
└── go.sum
```

### Command Tree

```
movie-cli
├── hello                      # Print greeting with version
├── version                    # Show version, commit, build date
├── self-update (alias: update) # git pull --ff-only
└── movie                     # Parent command
    ├── config [get|set] [key] [value]
    ├── scan [folder]
    ├── ls
    ├── search <name>
    ├── info <id|title>
    ├── suggest [N]
    ├── move [directory]
    ├── rename
    ├── undo
    ├── play <id>
    ├── stats
    └── tag [add|remove|list] [id] [tag]
```

---

## 3. Data Storage

### Base Directory

All data resides in `./data/`:

```
./data/
├── movie.db                   # SQLite database (WAL mode)
├── thumbnails/
│   └── <slug>/
│       └── <slug>.jpg         # Downloaded poster images
└── json/
    ├── movie/                 # JSON metadata per movie (future use)
    ├── tv/                    # JSON metadata per TV show (future use)
    └── history/
        └── <slug>/
            └── move-log.json  # Append-only move operation log
```

### Database Schema

**SQLite with WAL mode enabled.**

#### `media` table
| Column | Type | Notes |
|---|---|---|
| id | INTEGER PK AUTO | |
| title | TEXT NOT NULL | Display title (from TMDb or cleaned) |
| clean_title | TEXT NOT NULL | Cleaned filename title |
| year | INTEGER | Release year |
| type | TEXT | `'movie'` or `'tv'` (CHECK constraint) |
| tmdb_id | INTEGER UNIQUE | TMDb ID (prevents duplicates) |
| imdb_id | TEXT | IMDb ID |
| description | TEXT | Overview/synopsis |
| imdb_rating | REAL | |
| tmdb_rating | REAL | |
| popularity | REAL | TMDb popularity score |
| genre | TEXT | Comma-separated genre names |
| director | TEXT | Comma-separated director names |
| cast_list | TEXT | Comma-separated cast names (top 10) |
| thumbnail_path | TEXT | Local path to downloaded poster |
| original_file_name | TEXT | Original filename before cleaning |
| original_file_path | TEXT | Original full path |
| current_file_path | TEXT | Current full path (updated on move) |
| file_extension | TEXT | e.g., `.mkv`, `.mp4` |
| file_size | INTEGER | Size in bytes |
| scanned_at | DATETIME | Auto: CURRENT_TIMESTAMP |
| updated_at | DATETIME | Auto: CURRENT_TIMESTAMP |

#### `move_history` table
| Column | Type | Notes |
|---|---|---|
| id | INTEGER PK AUTO | |
| media_id | INTEGER FK→media | CASCADE delete |
| from_path | TEXT NOT NULL | Source path |
| to_path | TEXT NOT NULL | Destination path |
| original_file_name | TEXT | Name before move |
| new_file_name | TEXT | Name after move |
| moved_at | DATETIME | Auto |
| undone | INTEGER DEFAULT 0 | 0=active, 1=undone |

#### `config` table
| Column | Type | Notes |
|---|---|---|
| key | TEXT PK | Config key |
| value | TEXT NOT NULL | Config value |

#### `scan_history` table
| Column | Type | Notes |
|---|---|---|
| id | INTEGER PK AUTO | |
| folder_path | TEXT NOT NULL | Scanned directory |
| total_files | INTEGER | |
| movies_found | INTEGER | |
| tv_found | INTEGER | |
| scanned_at | DATETIME | Auto |

#### `tags` table
| Column | Type | Notes |
|---|---|---|
| id | INTEGER PK AUTO | |
| media_id | INTEGER FK→media | CASCADE delete |
| tag | TEXT NOT NULL | |
| created_at | DATETIME | Auto |
| UNIQUE(media_id, tag) | | |

#### Indexes
- `idx_media_type` on media(type)
- `idx_media_title` on media(clean_title)
- `idx_media_year` on media(year)
- `idx_media_tmdb` on media(tmdb_id)
- `idx_move_history_media` on move_history(media_id)
- `idx_move_history_undone` on move_history(undone)
- `idx_tags_media` on tags(media_id)

#### Default Config Values
| Key | Default |
|---|---|
| movies_dir | `~/Movies` |
| tv_dir | `~/TVShows` |
| archive_dir | `~/Archive` |
| scan_dir | `~/Downloads` |
| page_size | `20` |

---

## 4. Commands

### 4.1 `movie-cli hello`

**Purpose**: Print a greeting with version info.  
**Args**: None  
**Output**: `👋 Hello from Movie CLI!` + version string  

### 4.2 `movie-cli version`

**Purpose**: Show version, commit hash, and build date.  
**Args**: None  
**Output**: `v1.2.0 (commit: abc1234, built: 2024-06-01)`  
**Note**: Values injected via `-ldflags` at build time. Defaults: `v0.0.1-dev`, `none`, `unknown`.

### 4.3 `movie-cli self-update`

**Aliases**: `update`  
**Purpose**: Sync the local source repo to the latest code, or bootstrap a fresh local repo if none exists.  
**Args**: None  
**Behavior**:
1. Verify `git` is in PATH
2. Resolve repo path by checking the binary directory, current working directory, and a sibling `movie-cli-v3/` folder
3. If no local repo exists, clone a fresh copy next to the binary
4. If an existing repo is found, run `git status --porcelain` — repo must be clean (no uncommitted changes)
5. For an existing repo, record current commit (`git rev-parse --short HEAD`)
6. For an existing repo, run `git pull --ff-only`
7. Record resulting commit
8. Report one of: bootstrap success, already up-to-date, or old→new commit
9. Instruct the user to rebuild with `pwsh run.ps1`

**Error cases**: git not found, clone failed, dirty working tree, merge conflicts.

### 4.4 `movie-cli movie config`

**Usage**: `config [get|set] [key] [value]`  
**No args**: Show all config keys with values (API key is masked: first 4 + `...` + last 4 chars).  
**`get <key>`**: Show single config value.  
**`set <key> <value>`**: Update config value.  
**Valid keys**: `movies_dir`, `tv_dir`, `archive_dir`, `scan_dir`, `tmdb_api_key`, `page_size`.

### 4.5 `movie-cli movie scan [folder]`

**Purpose**: Scan a directory for video files, clean names, fetch TMDb metadata, save to DB.  
**Args**: Optional folder path. Falls back to `scan_dir` config.  
**Behavior**:
1. Resolve folder (arg → config → error)
2. Expand `~` to home directory
3. Validate folder exists and is a directory
4. Get TMDb API key (config → env `TMDB_API_KEY`). Warns if missing but continues without metadata.
5. Iterate entries:
   - **Files**: Must pass `IsVideoFile()` check
   - **Directories**: Look inside for the first video file; use directory name for cleaning
6. For each video file:
   - Clean filename → extract title, year, type (movie/tv)
   - Skip if `OriginalFilePath` already exists in DB (dedup by path match via `SearchMedia`)
   - Build `Media` record with file metadata
   - If API key available:
     - `SearchMulti(cleanTitle + year)`
     - Take first result
     - Fetch full details: `GetMovieDetails`/`GetTVDetails` + `GetMovieCredits`/`GetTVCredits`
     - Extract: IMDb ID, genres, directors, top 10 cast
     - For TV: include `Executive Producer` in directors (max 5)
     - Download poster → `./data/thumbnails/<slug>/<slug>.jpg`
   - Insert into DB (or update if `tmdb_id` conflict)
7. Log to `scan_history` table
8. Print summary: total files, movies, TV shows, skipped

### 4.6 `movie-cli movie ls`

**Purpose**: Paginated list of locally indexed media that have actual files on disk.  
**Args**: None  

**Filter Rule — File-Backed Only**:  
`ls` displays ONLY media records that have a non-empty `current_file_path` (i.e., items added via `scan` or `move`). Records created via `search` or `info` that only have metadata (no local file) are excluded. This prevents the list from being cluttered with catalog-only entries the user cannot play, move, or rename.

**Behavior**:
1. Read `page_size` from config (default 20)
2. Query media WHERE `current_file_path IS NOT NULL AND current_file_path != ''`
3. Count total matching media
4. Display page with: number, clean title, year, rating (TMDb→IMDb fallback), type icon
5. Interactive navigation: `N`=next, `P`=prev, `Q`=quit, number=view detail
6. Detail view: full metadata card with title, year, type, ratings, genres, director, cast, description, thumbnail path, file path
7. Clears terminal screen (`\033[H\033[2J`) between pages

### 4.7 `movie-cli movie search <name>`

**Purpose**: Search TMDb API live, select a result, fetch full details, save to DB.  
**Args**: One or more words joined as query  
**CRITICAL**: Does NOT require the file to exist locally. Catalogs metadata regardless.  
**Behavior**:
1. Require TMDb API key (config → env). Exits if missing.
2. `SearchMulti(query)` → filter to movie/tv only
3. Display up to 15 results with: number, icon, title, year, rating, type label
4. User selects a number (0 to cancel)
5. Fetch full details + credits (same as scan)
6. Download poster
7. Insert or update in DB
8. Print saved details summary

### 4.8 `movie-cli movie info <id|title>`

**Purpose**: Show detailed info for a media item.  
**Args**: Numeric ID or title string  
**Lookup Priority**:
1. **Numeric ID** → `GetMediaByID(id)` from local DB
2. **Title string** → `SearchMedia(title)` from local DB with match priority:
   - Exact match (case-insensitive)
   - Prefix match
   - First result
3. **Not found locally** → TMDb API search:
   - Check if TMDb ID already in DB (avoid duplicates)
   - Fetch full details + credits + poster
   - Auto-save to DB
   - Display result

### 4.9 `movie-cli movie suggest [N]`

**Purpose**: Recommend movies/TV shows based on library patterns.  
**Args**: Optional count (default 10)  
**Behavior**:
1. Require TMDb API key
2. Interactive category selection: Movie / TV / Random
3. **Movie or TV**:
   - Analyze library genre frequency via `TopGenres()`
   - Show user's top 3 genres
   - Get existing media IDs to avoid duplicates
   - Pick random library items → `GetRecommendations()` from TMDb
   - Fill remaining slots with `Trending(mediaType)`
   - Falls back to trending if not enough library data
4. **Random**:
   - Fetch trending movies + trending TV
   - Merge, shuffle, deduplicate
5. Display results: title, year, rating, genre names

### 4.10 `movie-cli movie move [directory]`

**Purpose**: Browse a directory, select a video file, move it to a configured destination with clean naming.  
**Args**: Optional source directory  
**Behavior**:
1. **Source directory resolution**:
   - From argument (expand `~`)
   - Or interactive prompt: Downloads / Scan Dir (if different) / Custom path
2. List video files in directory with clean titles, type icons, file sizes
3. User selects a file by number
4. **Destination prompt**: Movies dir / TV Shows dir / Archive dir / Custom path (from config)
5. Generate clean filename: `Title (Year).ext`
6. Show confirmation with from/to paths
7. Create destination directory if needed
8. **Move file** (with cross-drive fallback):
   - Attempt `os.Rename()` first (atomic, same-filesystem only)
   - If `os.Rename()` returns `EXDEV` (cross-device link error): fallback to `io.Copy` + `os.Remove`
   - Fallback copies file content, preserves original permissions, then removes source only after successful copy + close
   - Any copy/remove error is reported and source file is NOT deleted
9. **Track in database**:
   - Search for existing media record by title+path
   - Insert new record if not found, or update path if found
   - Log to `move_history` table (for undo)
   - Append to JSON history file: `./data/json/history/<slug>/move-log.json`
10. Print success message

### 4.11 `movie-cli movie rename`

**Purpose**: Batch rename files in the library to clean format.  
**Args**: None  
**Behavior**:
1. Load all media from DB (up to 10,000)
2. For each with a `current_file_path`:
   - Compare current filename vs `ToCleanFileName(cleanTitle, year, ext)`
   - Collect items where names differ
3. Show preview of all renames
4. Confirm with user (`y/N`)
5. For each rename:
    - Move file using `MoveFile()` (same cross-drive fallback as `movie move`)
    - Update `current_file_path` in DB
   - Log to `move_history` (for undo support)
6. Print summary: `X/Y files renamed`

### 4.12 `movie-cli movie undo`

**Purpose**: Revert the most recent move/rename operation.  
**Args**: None  
**Behavior**:
1. `GetLastMove()` — latest `move_history` record where `undone=0`
2. Print what will be undone (from → to paths)
3. Prompt user: `Undo this? [y/N]: `
4. If declined, print cancellation and exit
5. Verify file exists at `to_path`
6. Move file back using `MoveFile()` (same cross-drive fallback as `movie move`)
7. Mark record as `undone=1` in DB
8. Update media `current_file_path` back to `from_path`
9. Print success message

### 4.13 `movie-cli movie play <id>`

**Purpose**: Open a media file with the system's default video player.  
**Args**: Numeric media ID (required)  
**Behavior**:
1. Look up media by ID
2. Verify `current_file_path` exists on disk
3. Open with platform-specific command:
   - **macOS**: `open <path>`
   - **Linux**: `xdg-open <path>`
   - **Windows**: `cmd /c start "" <path>`

### 4.14 `movie-cli movie stats`

**Purpose**: Display library statistics.  
**Args**: None  
**Output**:
1. Counts: total movies, total TV shows, total
2. Storage: total file size, largest file, smallest file, average file size (human-readable)
3. Top 10 genres with visual bar chart (`█` characters, max 30 width)
4. Average IMDb rating (if available)
5. Average TMDb rating (if available)
**Note**: Loads all media (up to 10,000) to compute averages. File size stats use a dedicated SQL aggregate query.

### 4.15 `movie-cli movie tag`

**Purpose**: Manage user-defined tags on media items.  
**Subcommands**:

#### `movie tag add <id> <tag>`
- Look up media by ID → error if not found
- Check if tag already exists → error `"tag already exists"`
- Insert into `tags` table
- Print confirmation: `Tag "favorite" added to "Title (Year)"`

#### `movie tag remove <id> <tag>`
- Look up media by ID → error if not found
- Delete tag from `tags` table → error `"tag not found"` if no rows affected
- Print confirmation: `Tag "favorite" removed from "Title (Year)"`

#### `movie tag list [id]`
- **With ID**: Show all tags for that media item, or `"No tags"` if none
- **Without ID**: Show all unique tags with media count, e.g. `favorite (3)`, `watchlist (7)`

**Acceptance Criteria**:
- GIVEN a media ID WHEN `tag add 1 favorite` THEN tag is inserted and confirmation printed
- GIVEN a duplicate tag WHEN `tag add 1 favorite` again THEN error "tag already exists"
- GIVEN a tag exists WHEN `tag remove 1 favorite` THEN tag is deleted and confirmation printed
- GIVEN a non-existent tag WHEN `tag remove 1 unknown` THEN error "tag not found"
- GIVEN tags exist WHEN `tag list 1` THEN all tags for media 1 are shown
- GIVEN tags exist WHEN `tag list` (no ID) THEN all unique tags with counts are shown

---

## 5. Filename Cleaner

**Package**: `cleaner`

### Video Extensions Supported
`.mkv`, `.mp4`, `.avi`, `.mov`, `.wmv`, `.flv`, `.webm`, `.m4v`, `.ts`, `.vob`, `.ogv`, `.mpg`, `.mpeg`, `.3gp`

### Cleaning Pipeline (`Clean(filename)`)
1. Extract file extension
2. **Detect TV show**: Regex `S\d{1,2}E\d{1,2}|Season\s*\d+|Episode\s*\d+` → type = `"tv"`, else `"movie"`
3. **Extract year**: Regex `\b(19|20)\d{2}\b` (first match)
4. Replace `.` and `_` with spaces
5. Remove junk patterns (case-insensitive):
   - Quality: `1080p`, `720p`, `480p`, `2160p`, `4k`
   - Source: `bluray`, `bdrip`, `webrip`, `web-dl`, `hdtv`, `dvdrip`, etc.
   - Codec: `x264`, `x265`, `h264`, `h265`, `hevc`, `aac`, `dts`, etc.
   - Groups: `RARBG`, `YTS`, `YIFY`, `EZTV`, etc.
   - Editions: `extended`, `unrated`, `directors cut`, `remastered`, etc.
   - Language: `multi`, `dual`, `eng`, `hindi`, `subs`, etc.
   - Bracketed content: `[...]`, `(...)`, `{...}`
6. Collapse multiple spaces
7. If year found, truncate everything after the year
8. Return: `{ CleanTitle, Year, Type, Extension }`

### Helper Functions
- `IsVideoFile(name)` — check extension against known video types
- `ToSlug(title)` — lowercase, strip non-alphanumeric, spaces→hyphens
- `ToCleanFileName(title, year, ext)` — `"Title (Year).ext"` or `"Title.ext"` if no year

---

## 6. TMDb Integration

**Package**: `tmdb`  
**Base URL**: `https://api.themoviedb.org/3`  
**Image Base**: `https://image.tmdb.org/t/p/w500`  
**HTTP Timeout**: 15 seconds

### API Key Resolution
1. Config table: `tmdb_api_key`
2. Environment variable: `TMDB_API_KEY`

### Endpoints Used
| Method | Endpoint | Purpose |
|---|---|---|
| `SearchMulti` | `/search/multi` | Search movies + TV, filters to movie/tv only |
| `GetMovieDetails` | `/movie/{id}` | Full movie details (IMDb ID, genres, runtime) |
| `GetTVDetails` | `/tv/{id}` | Full TV details (genres, seasons) |
| `GetMovieCredits` | `/movie/{id}/credits` | Cast + crew (directors) |
| `GetTVCredits` | `/tv/{id}/credits` | Cast + crew (directors + exec producers) |
| `GetRecommendations` | `/{type}/{id}/recommendations` | Similar content |
| `DiscoverByGenre` | `/discover/{type}` | Content by genre (not currently used in commands) |
| `Trending` | `/trending/{type}/week` | Weekly trending content |
| `DownloadPoster` | Image CDN | Download poster JPG to local storage |

### Genre ID Mapping
Hardcoded map of 27 genre IDs covering both movie and TV genres (Action through War & Politics).

---

## 7. Build & Deploy

### Makefile Targets
| Target | Description |
|---|---|
| `make build` | Build for current OS → `./movie-cli` |
| `make build-windows` | Cross-compile → `movie-cli.exe` |
| `make build-mac-arm` | macOS ARM64 → `movie-cli-darwin-arm64` |
| `make build-mac-intel` | macOS AMD64 → `movie-cli-darwin-amd64` |
| `make clean` | Remove all binaries |
| `make install` | Build + copy to `/usr/local/bin` |

### Build Flags
```
-ldflags "-s -w
  -X github.com/movie/movie-cli/version.Version=<git tag>
  -X github.com/movie/movie-cli/version.Commit=<short SHA>
  -X github.com/movie/movie-cli/version.BuildDate=<YYYY-MM-DD>"
```

### PowerShell Deploy (`build.ps1`)
1. `git pull` (update source)
2. `go mod tidy`
3. Build binary with version ldflags
4. Copy to deploy directory:
   - **Windows default**: `E:\bin-run`
   - **Mac/Linux default**: `/usr/local/bin`
   - **Custom**: `pwsh build.ps1 -BinDir "C:\custom"`
5. Clean local build artifacts
6. Verify binary runs

---

## 8. Configuration Keys

| Key | Purpose | Default |
|---|---|---|
| `movies_dir` | Destination for movie files | `~/Movies` |
| `tv_dir` | Destination for TV show files | `~/TVShows` |
| `archive_dir` | Archive destination | `~/Archive` |
| `scan_dir` | Default scan source directory | `~/Downloads` |
| `tmdb_api_key` | TMDb API key | (none) |
| `page_size` | Items per page in `ls` | `20` |

---

## 9. Dependencies

| Package | Version | Purpose |
|---|---|---|
| `github.com/spf13/cobra` | v1.8.0 | CLI framework |
| `modernc.org/sqlite` | (in go.sum) | Pure-Go SQLite driver (no CGo) |

---

## 10. AI Implementation Risk Assessment

### Risk Matrix

| Area | Risk | Reason | Mitigation |
|---|---|---|---|
| **Go CLI in web IDE** | 🔴 Critical | No `go build`, no terminal, no file system access | Must develop locally or in Go-compatible environment |
| **File system ops** (move, rename, scan) | 🔴 High | Requires real OS paths, permissions, cross-drive moves | Test on actual OS; `os.Rename` fails across drives |
| **TMDb API** | 🟡 Medium | API key management, rate limits, response shape changes | Handle errors gracefully; cache responses |
| **SQLite migrations** | 🟡 Medium | Schema evolution, unique constraint handling | Test insert/update/upsert flows |
| **Cross-platform paths** | 🟡 Medium | `~/` expansion, path separators, `E:\bin-run` | Use `filepath.Join`, test on Windows + Mac |
| **Undo correctness** | 🟡 Medium | File may be manually deleted; multiple undo states | Check file exists before undo; handle errors |
| **Cleaner regex accuracy** | 🟢 Low | Edge cases in filenames | Expand test cases for unusual naming patterns |
| **Cobra CLI structure** | 🟢 Low | Well-established pattern | Follow existing patterns |
| **Version injection** | 🟢 Low | ldflags is standard Go practice | Already implemented correctly |

### Key Pitfalls for AI Implementation
1. **Don't use `package.json`** — this is Go, not Node.js
2. **Don't add a web server** — this is a CLI tool
3. **`os.Rename` doesn't work across filesystems** — may need copy+delete for cross-drive moves
4. **TMDb API key must never be committed** — use config or env var only
5. **SQLite `UNIQUE` on `tmdb_id`** — insert will fail on duplicates; must upsert
6. **`move_history.undone` is an INTEGER** (0/1), not BOOLEAN — SQLite limitation
7. **The `tags` table exists in schema but NO commands use it yet** — future feature
8. **`DiscoverByGenre` exists in tmdb client but is NOT used by any command** — future feature
9. **JSON metadata files** are referenced in storage structure but only `move-log.json` is actually written

### Estimated Success Rate
- **Expert Go developer**: ~95%
- **General-purpose AI (Go-capable)**: ~70-80%
- **Web-focused AI (Lovable, v0, etc.)**: ~20-30% — fundamental platform mismatch
- **AI without file system access**: ~0% for file operations

---

## Appendix A: File Size Formatting

```
bytes >= 1 GB → "X.X GB"
bytes >= 1 MB → "X.X MB"  
bytes >= 1 KB → "X.X KB"
bytes < 1 KB  → "X B"
```

## Appendix B: Move History JSON Format

Appended to `./data/json/history/<slug>/move-log.json`:

```json
{"from":"/path/from","to":"/path/to","timestamp":"2026-03-17T21:45:00+08:00"}
```

Timestamps use `time.Now().Format(time.RFC3339)`. See `spec/02-app/issues/01-hardcoded-timestamp.md` for the fix history.

## Appendix C: TV Director Handling

For TV shows, both `"Director"` and `"Executive Producer"` crew jobs are included in the director field, capped at 5 entries. Movies only include `"Director"`.
