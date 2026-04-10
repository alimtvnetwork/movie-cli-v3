# Project Overview

> **Last Updated**: 10-Apr-2026

## Project

- **Name**: Mahin CLI (formerly Movie CLI)
- **Type**: Go CLI application (NOT a web app)
- **Binary**: `mahin`
- **Language**: Go 1.22
- **Module**: `github.com/mahin/mahin-cli-v2`
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
mahin
├── hello                      # Greeting with version
├── version                    # Version/commit/build date + Go/OS info
├── self-update                # git pull --ff-only
├── changelog                  # Show changelog
└── movie
    ├── config                 # View/set configuration
    ├── scan                   # Scan folder → DB + TMDb + JSON metadata
    ├── ls                     # Paginated library list (file-backed only)
    ├── search                 # Live TMDb search → save
    ├── info                   # Local DB → TMDb fallback
    ├── suggest                # Recommendations/trending + genre discover
    ├── move                   # Browse + move + track (--all batch support)
    ├── rename                 # Batch clean rename
    ├── undo                   # Revert last move/rename
    ├── play                   # Open with default player
    ├── stats                  # Library statistics
    ├── tag                    # Add/remove/list tags
    └── export                 # Export library data
```

## Important Notes for AI

- **This is NOT a web project** — no dev server, no preview
- Build errors in Lovable (`no package.json found`, `no command found for task "dev"`) are **expected and MUST be ignored**
- All file operations require a real OS/terminal to test
- Full specification lives in `spec/` folder
- Milestone markers use `readm.txt` format: `let's start now {date} {time Malaysia}`
- **Always read memory files before making changes**

## File Structure (as of 10-Apr-2026)

- `cmd/` — 21 Go files (root, hello, version, update, changelog + movie parent + 14 subcommands + move_helpers)
- `cleaner/` — 1 file (filename cleaning)
- `tmdb/` — 1 file (API client)
- `db/` — 6 files (db.go, media.go, config.go, history.go, helpers.go, tags.go)
- `updater/` — 1 file (git self-update)
- `version/` — 1 file (build-time vars)
- `.github/` — Release pipeline (release.yml)
- `spec/` — Structured specification docs
- `docs/` — Additional documentation
