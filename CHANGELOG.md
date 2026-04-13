# Changelog

All notable changes to this project will be documented in this file.

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
