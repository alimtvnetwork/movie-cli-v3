# Changelog

All notable changes to this project will be documented in this file.

## v1.28.0

### Added
- **Centralized error logging system** (`errlog/logger.go`) — all errors are now logged to:
  - `.movie-output/logs/error.txt` (file-based, append-only, with timestamp/source/stack trace)
  - `error_logs` DB table (queryable, with level/source/function/command/workdir/stack trace)
- **`error_logs` table** (`db/errorlog.go`) — new table with columns: timestamp, level (ERROR/WARN/INFO), source, function, command, work_dir, message, stack_trace; includes `RecentErrorLogs()` query
- **`errlog` package** — `Error()`, `Warn()`, `Info()` functions with automatic caller detection, stack trace capture (errors only), and dual output (file + DB)
- **DB writer injection** — `errlog.SetDBWriter()` allows wiring DB logging without circular imports

### Changed
- **`movie scan` errors** — DB search, stat, insert, update, JSON write, TMDb, and thumbnail errors now use `errlog` instead of raw `fmt.Fprintf(os.Stderr)`
- **`movie rest` errors** — JSON encode, template render, watchlist update, tag add, config read errors now use `errlog`
- **Error entries include**: timestamp, severity, source file:line, function name, CLI command, working directory, message, and full Go stack trace

## v1.27.0

### Changed
- **Modernized HTML report** — complete UI overhaul: sticky toolbar with inline search, genre/rating/sort dropdowns, type filter pills, dark zinc theme, result count, empty state, keyboard shortcut (`/` to search, `Esc` to close modal), responsive layout
- **Search now searches titles, directors, and cast** — not just titles
- **Genre filter dropdown** — auto-populated from scan data
- **Rating filter dropdown** — filter by minimum rating (5+ through 9+)
- **Sort options** — sort by title, rating, or year (ascending/descending)
- **Connected REST indicator** — banner shows green dot when REST server is detected

### Fixed
- **`writeJSON` error swallowed** — `json.Encoder.Encode` error now logged to stderr
- **`tmpl.Execute` error swallowed** — template render error now logged to stderr
- **`GetConfig` errors swallowed** — `tmdb_api_key` and `tmdb_token` config read errors now logged
- **`database.Exec` watchlist update error swallowed** — now logged to stderr
- **`database.AddTag` watched tag error swallowed** — now logged to stderr
- **JS error handling** — all `catch(e)` blocks now show specific error messages; `fetch` non-ok responses show HTTP status/body

## v1.26.0

### Added
- **`GET /` on REST server** — serves a live HTML library report rendered from the database; always up-to-date, no need to open a static file

## v1.25.0

### Added
- **HTML report: tag management** — add/remove tags per card with inline input; tags shown as purple pills with ✕ to remove
- **HTML report: mark watched** — 👁 button marks a movie as watched via REST API; card gets green border and "watched" tag
- **HTML report: similar movies** — 🔍 button opens a modal with TMDb recommendations (poster, title, year, rating, description)
- **HTML report: watched filter** — new "✅ Watched" filter button in the toolbar
- **HTML report: tags auto-load** — when REST server is detected, all tags load automatically on page open

## v1.24.0

### Added
- **`GET/POST/DELETE /api/tags`** — full tag management via REST: list all tags with counts, list tags per media, add tag, remove tag
- **`GET /api/media/{id}/similar`** — fetches TMDb recommendations for a media item
- **`PATCH /api/media/{id}/watched`** — marks a media item as watched (updates watchlist + adds "watched" tag)
- **Refactored REST handlers** — new endpoints in `cmd/movie_rest_handlers.go` to keep files under 200 lines

## v1.23.0

### Added
- **`movie rest --open`** — automatically opens the HTML report in the default browser when the REST server starts; supports macOS (`open`), Windows (`rundll32`), and Linux (`xdg-open`)

## v1.22.0

### Added
- **`movie rest`** — starts a local REST API server (default port 8086, `--port` to override) exposing library endpoints: `GET /api/media`, `GET/DELETE/PATCH /api/media/{id}`, `GET /api/stats`; enables interactive features in the HTML report
- **HTML report** — `movie scan` now generates `report.html` in `.movie-output/` with responsive card layout showing thumbnail, title, year, rating, genre, director, cast, description, and tagline; includes search, filter, and delete via REST API
- **`templates/report.html`** — external HTML template file (not embedded in Go code); bundled via Go `embed` at compile time through `templates/embed.go`

## v1.21.0

### Added
- **`movie db`** — prints the resolved database path, data directory, and record counts for debugging

## v1.20.0

### Changed
- **Renamed `<package>/<package>.go` files** — `db/db.go` → `db/open.go`, `cleaner/cleaner.go` → `cleaner/parse.go`, `updater/updater.go` → `updater/run.go`, `version/version.go` → `version/info.go`; enforced as a permanent naming convention

## v1.19.0

### Added
- **`movie history --format table`** — output move history as a formatted table with columns: #, Title, From, To, Date, Status

## v1.18.0

### Added
- **Binary-relative data storage** — all data (database, thumbnails, JSON metadata) is now stored in `data/` next to the CLI binary, not the working directory
- **`run.ps1` deploys data folder** — build script copies data directory alongside the deployed binary

## v1.17.0

### Added
- **`movie ls --format table`** — output library listing as a formatted table with columns: #, Title, Year, Type, Rating, Genre, Director (no interactive pager)

### Changed
- **Refactored `movie_ls.go`** — split 313-line file into `movie_ls.go` (196), `movie_ls_table.go` (99), and `movie_ls_detail.go` (120)

## v1.16.0

### Changed
- **Refactored `movie_search.go`** — extracted save-and-print logic into `cmd/movie_search_save.go` (135 lines); `movie_search.go` reduced from 240 to 135 lines

## v1.15.0

### Added
- **`movie stats --format table`** — output library statistics as a formatted key-value table with sections for counts, storage, genres, and ratings

## v1.14.0

### Changed
- **Refactored `movie_info.go`** — extracted `fetchMovieDetails` and `fetchTVDetails` into `cmd/movie_fetch_details.go`

## v1.13.0

### Fixed
- **`movie update` fresh-clone flow** — when no local repo exists, a new clone is now reported as bootstrap success instead of incorrectly saying "Already up to date"
- **Self-update specs** — documented repo bootstrap vs existing-repo pull behavior using the GitMap-aligned update flow

## v1.12.0

### Added
- **`movie search --format table`** — output TMDb search results as a formatted table (no interactive prompt); columns: #, Title, Year, Type, Rating, TMDb ID
- **`movie info --format table`** — output media detail as a key-value formatted table; shows all metadata fields dynamically

## v1.11.0

### Added
- **`movie search --format json`** — output TMDb search results as a JSON array to stdout (no interactive prompt); pipeable to `jq` and scripts
- **`movie info --format json`** — output media detail as a JSON object to stdout; includes source field ("local" or "tmdb")

## v1.10.0

### Added
- **`movie ls --format json`** — output entire library as a JSON array to stdout; includes id, title, year, type, ratings, genre, file path, and file size per item
- **`movie stats --format json`** — output library statistics as a JSON object to stdout; includes counts, storage, top genres, and average ratings

## v1.9.0

### Added
- **`movie scan --format json`** — output scan results as structured JSON to stdout for piping to `jq`, scripts, or other tools; includes metadata, counts, and per-item details; works with `--dry-run` too

## v1.8.0

### Fixed
- **`movie scan` no longer fails when TMDb is unset** — media with no TMDb match/key now store `tmdb_id` as `NULL` instead of `0`, so bulk scans no longer hit `UNIQUE constraint failed: media.tmdb_id`
- **Interactive TMDb setup before scan** — when TMDb is not configured, `movie scan` now prompts for a TMDb API key and TMDb access token before scanning starts; leaving both blank continues without metadata
- **TMDb bearer token support** — scan can now authenticate with either `tmdb_api_key` or `tmdb_token`

## v1.7.1

### Changed
- **Refactored `movie_scan.go`** — split from ~500 lines into 4 focused files:
  - `movie_scan.go` (~120 lines) — command definition, orchestrator, helpers
  - `movie_scan_collect.go` (~110 lines) — video file discovery and path utilities
  - `movie_scan_process.go` (~170 lines) — per-file processing and TMDb enrichment
  - `movie_scan_table.go`, `movie_scan_json.go`, `movie_scan_summary.go` — unchanged

## v1.7.0

### Added
- **`movie scan --format table`** — display scan results as a formatted table with columns for #, filename, clean title, year, type, rating, and status; works with `--dry-run` too

## v1.6.0

### Added
- **`movie scan --dry-run`** — preview what would be scanned (files found, cleaned titles, types) without writing to DB or creating `.movie-output/`

## v1.5.0

### Added
- **`movie scan --depth N` (`-d`)** — limit recursive scan to N subdirectory levels (0 = unlimited); e.g. `movie scan -r -d 2`

## v1.4.0

### Fixed
- **`movie update` works from anywhere** — no longer requires CWD to be inside the git repo; finds the repo next to the binary, clones fresh if needed

## v1.3.0

### Added
- **`movie scan --recursive` (`-r`)** — scan all subdirectories recursively instead of just top-level entries; skips `.movie-output` and hidden directories automatically

### Changed
- **Refactored scan internals** — extracted `collectVideoFiles`, `processVideoFile`, and `enrichFromTMDb` helpers for cleaner architecture and reuse

## v1.2.0
### Changed
- **`movie scan` defaults to current directory** — running `movie scan` without arguments now scans the CWD instead of a config-stored `scan_dir` path
- **Scan output to `.movie-output/`** — all scan results (per-item JSON, summary.json with categories/descriptions/metadata) are now written to `.movie-output/` inside the scanned folder

### Added
- **`summary.json`** — comprehensive scan report with total counts, genre-based categories, and full TMDb metadata per item

## v1.1.0

### Fixed
- **`run.ps1` version stamping** — now reads the version from `version/version.go` and injects commit/build date into the correct `version` package variables
- **`run.ps1` version summary** — now reports the binary that was just built/deployed instead of accidentally showing an older `movie` found earlier in `PATH`
- **Deployed changelog visibility** — `run.ps1` now copies `CHANGELOG.md` beside the deployed binary and verifies `movie changelog --latest`

## v0.2.4

### Fixed
- **`GetConfig` false warnings** — `movie_info.go` and `movie_scan.go` now explicitly ignore `sql: no rows in result set` from `GetConfig`, preventing false-positive error messages when config keys are unset
- **Indentation fix** — corrected misleading indentation in `movie_scan.go` error block

### Changed
- **JSON export completeness** — `movie_export.go` now includes all 6 previously missing metadata fields: `Runtime`, `Language`, `Budget`, `Revenue`, `TrailerURL`, `Tagline`

## v0.2.3

### Fixed
- **`db/media.go` silent scan error** — `TopGenres` now returns a wrapped error on `rows.Scan` failure instead of silently using `continue`
- **`movie_info.go` poster error swallowed** — `DownloadPoster` failures now logged to stderr
- **`movie_scan.go` poster error swallowed** — `DownloadPoster` failures now logged to stderr
- **`movie_scan.go` subdirectory read error** — `os.ReadDir` failures in subdirectory scanning now logged instead of silently skipped
- **`movie_undo.go` permission error masked** — `os.Stat` now distinguishes permission errors from file-not-found and logs them separately

## v0.2.2

### Fixed
- **`movie_search.go` unchecked `GetConfig`** — API key lookup now checks for errors before proceeding
- **`movie_suggest.go` unchecked `GetConfig`** — API key lookup now checks for errors and handles `sql: no rows` correctly
- **`movie_resolve.go` unbounded query** — `resolveByTitle` now uses `LIMIT 1` to prevent scanning full table
- **`db/media.go` missing `rows.Err()` check** — `TopGenres` now checks `rows.Err()` after iteration loop

### Changed
- **`movie_search.go` duplicate detail fetch removed** — eliminated redundant `GetMovieDetails`/`GetTVDetails` calls that were already handled by shared `fetchMovieDetails`/`fetchTVDetails` helpers

## v0.2.1

### Fixed
- **`movie_move.go` unchecked error** — `database.GetConfig("movies_dir")` error now handled instead of silently ignored
- **`movie_move.go` unchecked error** — `database.GetConfig("tv_dir")` error now handled instead of silently ignored
- **`movie_move_helpers.go` cross-drive cleanup** — copy+delete fallback now removes the source file after successful copy
- **`movie_rename.go` unchecked `InsertMoveHistory`** — rename history logging error now reported to stderr
- **`movie_play.go` unchecked `exec.Command` error** — player launch error now reported to stderr
- **`movie_stats.go` unchecked `CountMedia`** — movie/TV count errors now handled instead of silently returning zero
- **`movie_watch.go` unchecked `GetConfig`** — API key lookup now checks for errors before proceeding
- **`tmdb/client.go` unchecked `json.NewDecoder` error** — HTTP response body decoding errors now properly returned
- **`updater/updater.go` unchecked exec errors** — `git pull` and `go build` errors now returned instead of silently ignored

## v1.0.0

### Added
- **Batch move** (`movie move --all`) — move all video files at once with auto-routing to movies/TV directories, preview table, and `[y/N]` confirmation
- **JSON metadata export** — `movie scan` now writes per-file JSON metadata to `./data/json/movie/` and `json/tv/`
- **Genre-based discovery** — `movie suggest` uses `DiscoverByGenre` for TMDb genre-based recommendations (3-phase: genre discovery → recommendations → trending fallback)
- **`GenreNameToID()` helper** — reverse genre map in tmdb package for name-to-ID lookups
- **CI pipeline** (`.github/workflows/ci.yml`) — lint (`go vet` + `golangci-lint`), vulnerability scanning (`govulncheck`), parallel test matrix, cross-compiled builds (6 targets), SHA deduplication
- **Release pipeline** (`.github/workflows/release.yml`) — triggers on `release/**` branches and `v*` tags, cross-compiled binaries, SHA256 checksums, version-pinned install scripts, changelog extraction
- **Cross-platform install scripts** — `install.sh` (Linux/macOS) and `install.ps1` (Windows) with checksum verification and PATH setup
- **`.golangci.yml`** — sensible linter defaults (errcheck, govet, staticcheck, gocritic, misspell, errorlint, etc.)
- **Undo confirmation prompt** — `movie undo` shows from/to paths and asks `[y/N]` before reverting
- **Tag command** (`movie tag`) — add, remove, and list tags on media entries
- **Comprehensive CLI help** — root command shows version + categorized help with examples; `movie --version` flag; `movie version` shows Go/OS/arch

### Changed
- **`movie ls`** now only shows scan-indexed items (filters by non-empty `original_file_path`)
- **`movie suggest`** upgraded from recommendations-only to 3-phase strategy (DiscoverByGenre → Recommendations → Trending)
- **Repository migrated** from `movie-cli-v1` to `movie-cli-v2` to `movie-cli-v3` across all imports, workflows, and docs

### Fixed
- Timestamp bug — `saveHistoryLog` now uses `time.Now().Format(time.RFC3339)` instead of hardcoded "now"
- Deduplicated TMDb fetch logic — shared `fetchMovieDetails()`/`fetchTVDetails()` helpers
- Cross-drive move fallback — copy+delete when `os.Rename` fails with `EXDEV`

## v0.1.0

### Added
- Core CLI with Cobra: `hello`, `version`, `self-update` commands
- Movie management: `scan`, `ls`, `search`, `info`, `suggest`, `move`, `rename`, `undo`, `play`, `stats`, `config`
- SQLite database with WAL mode, 5 tables, 7 indexes
- TMDb API client (search, details, credits, recommendations, trending, posters)
- Filename cleaner (junk removal, year extraction, TV detection)
- PowerShell build & deploy pipeline (`run.ps1`)
- Full project specification in `spec/`
