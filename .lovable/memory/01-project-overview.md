# Project Overview

> **Last Updated**: 17-Mar-2026

## Project

- **Name**: Movie CLI
- **Type**: Go CLI application (NOT a web app)
- **Binary**: `movie-cli`
- **Language**: Go 1.22
- **Module**: `github.com/mahin/mahin-cli-v1`
- **Framework**: Cobra (CLI), SQLite (storage), TMDb API (metadata)

## Purpose

A cross-platform CLI tool for managing a personal movie and TV show library. It scans local folders for video files, cleans messy filenames, fetches metadata from TMDb, stores everything in SQLite, and organizes files into configured directories.

## Key Architecture Decisions

1. **Pure-Go SQLite** (`modernc.org/sqlite`) — no CGo dependency
2. **WAL mode** for SQLite concurrency
3. **TMDb API** for metadata (requires user-provided API key)
4. **git-based self-update** (`git pull --ff-only`)
5. **All data** stored in `./data/` (DB, thumbnails, JSON logs)

## Command Tree

```
movie-cli
├── hello                      # Greeting with version
├── version                    # Version/commit/build date
├── self-update                # git pull --ff-only
└── movie
    ├── config                 # View/set configuration
    ├── scan                   # Scan folder → DB + TMDb
    ├── ls                     # Paginated library list
    ├── search                 # Live TMDb search → save
    ├── info                   # Local DB → TMDb fallback
    ├── suggest                # Recommendations/trending
    ├── move                   # Browse + move + track
    ├── rename                 # Batch clean rename
    ├── undo                   # Revert last move/rename
    ├── play                   # Open with default player
    └── stats                  # Library statistics
```

## Important Notes for AI

- **This is NOT a web project** — no `package.json`, no dev server, no preview
- Build errors in Lovable (`no package.json found`, `no command found for task "dev"`) are **expected and MUST be ignored**
- All file operations require a real OS/terminal to test
- Full specification lives in `spec.md` at project root
- Milestone markers use `readm.txt` format: `let's start now {date} {time Malaysia}`
- `.gitignore` cannot be created in Lovable environment — must be added manually
- **Always read memory files before making changes** (see `workflow/01-ai-success-plan.md`)

## File Structure (as of 17-Mar-2026)

- `cmd/` — 15 Go files (root + hello + version + update + movie parent + 10 subcommands + move_helpers)
- `cleaner/` — 1 file (filename cleaning)
- `tmdb/` — 1 file (API client)
- `db/` — 5 files (db.go, media.go, config.go, history.go, helpers.go)
- `updater/` — 1 file (git self-update)
- `version/` — 1 file (build-time vars)

## Session History (17-Mar-2026)

1. Updated `readm.txt` with milestone marker
2. Attempted `.gitignore` creation (blocked by Lovable environment)
3. Fixed timestamp bug in `saveHistoryLog` — `"now"` → `time.Now().Format(time.RFC3339)`
4. Deduplicated TMDb fetch logic — `scan` and `search` now use shared helpers from `movie_info.go`
5. Split `cmd/movie_move.go` (348 lines) → `movie_move.go` + `movie_move_helpers.go`
6. Split `db/sqlite.go` (452 lines) → 5 focused files
7. Updated all memory files
8. Created AI success rate plan
