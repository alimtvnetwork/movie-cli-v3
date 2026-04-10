# Movie CLI — Development Log

> **Project**: Movie CLI  
> **Language**: Go 1.22  
> **Module**: `github.com/mahin/movie-cli-v2`  
> **Purpose**: Cross-platform CLI tool for managing a personal movie and TV show library  
> **Development Period**: 15-Mar-2026 → 18-Mar-2026  
> **Last Updated**: 05-Apr-2026

---

## How to Read This File

This document records **every user prompt** and the **exact modifications** made in response. Each entry follows this format:

1. **User Prompt** — what was asked (verbatim or summarized)
2. **What Was Done** — files created, edited, or deleted
3. **Technical Details** — key code changes, logic, and architecture decisions

Any AI reading this file should be able to reconstruct the full project from scratch or continue development from where it left off.

---

## Table of Contents

1. [Session 1: Project Bootstrap](#session-1-project-bootstrap-15-mar-2026)
2. [Session 2: Movie CLI Feature — Specification](#session-2-movie-cli-feature--specification-16-mar-2026)
3. [Session 3: Movie CLI Feature — Full Implementation](#session-3-movie-cli-feature--full-implementation-16-mar-2026)
4. [Session 4: Command Refinements & Build Tooling](#session-4-command-refinements--build-tooling-17-mar-2026)
5. [Session 5: Specification, Memory & Quality Audit](#session-5-specification-memory--quality-audit-17-mar-2026)
6. [Session 6: Code Hygiene Fixes & Documentation](#session-6-code-hygiene-fixes--documentation-17-18-mar-2026)

---

## Session 1: Project Bootstrap (15-Mar-2026)

### Prompt 1.1 — Milestone Marker
> **User**: Write a readme.txt text file with 3 words with no context at all "let's start now date:dd-MMM-YYYY) (time:12 hr clock format exact time now for Malaysia}"

**What Was Done**:
- Created `readme.txt` with content: `let's start now date:15-Mar-2026 time:08:00 AM`

**Files Created**: `readme.txt`

---

### Prompt 1.2 — Milestone Marker (Time Correction)
> **User**: (Same prompt repeated — time correction)

**What Was Done**:
- Updated `readme.txt` to: `let's start now date:15-Mar-2026 time:04:00 PM`

**Files Modified**: `readme.txt`

---

### Prompt 1.3 — Milestone Marker (New Day)
> **User**: Write a readm.txt text file with 3 words...

**What Was Done**:
- Created `readm.txt` with content: `let's start now date:16-Mar-2026 time:11:00 PM`
- Note: From this point, the file is named `readm.txt` (not `readme.txt`)

**Files Created**: `readm.txt`

**Convention Established**: Project milestones are tracked via `readm.txt` containing `let's start now {date} {time}` in Malaysia time (UTC+8).

---

## Session 2: Movie CLI Feature — Specification (16-Mar-2026)

### Prompt 2.1 — Initial Feature Request
> **User**: I want to create another option in this CLI that is a movie CLI. This CLI can scan a folder where there can be a lot of movies. With this scan, it reads the folder name. If I write movie ls command, it reads all files and shows 20 of them from the top, with rename clean the file name...

**What Was Done**:
- Organized a structured plan for the Movie CLI feature
- Proposed commands: `scan`, `ls`, `search`, `info`, `suggest`, `move`, `undo`, `rename`, `play`, `stats`, `config`
- Proposed filename cleaning via regex (remove `1080p`, `BluRay`, `x264`, etc.)
- Proposed TMDb API integration for metadata
- Proposed SQLite database for storage

**Files Created**: None (planning only)

---

### Prompt 2.2 — Detailed Specification
> **User**: Corrected Description of Your Movie CLI. Your CLI will include a movie module that helps manage movies and TV shows stored in a folder. 1. movie scan — The CLI can scan a folder containing movies or TV shows. During scanning it will: Read all files and folders, Clean the file names, Identify whether it is a Movie or TV Show, Fetch metadata from an online database...

**What Was Done**:
- Refined the specification based on the user's corrected description
- Proposed interactive `movie move` design with fuzzy-search selection flow
- Decided on Go libraries (`survey` or `bubbletea`) for interactive prompts

**Files Created**: None (specification refinement only)

---

### Prompt 2.3 — File Structure & Database Design
> **User**: Show me what the file structure will be. And I want to use SQLiteDB

**What Was Done**:
- Proposed full project structure:
  ```
  cmd/movie.go           — Parent Cobra command
  cmd/movie_scan.go      — Scan command
  cmd/movie_ls.go        — List command
  cmd/movie_search.go    — Search command
  cmd/movie_info.go      — Info command
  cmd/movie_suggest.go   — Suggest command
  cmd/movie_move.go      — Move command
  cmd/movie_undo.go      — Undo command
  cmd/movie_rename.go    — Rename command
  cmd/movie_play.go      — Play command
  cmd/movie_stats.go     — Stats command
  cmd/movie_config.go    — Config command
  db/sqlite.go           — SQLite connection + CRUD
  cleaner/cleaner.go     — Filename cleaning
  tmdb/client.go         — TMDb API client
  ```
- Replaced previously discussed JSON-based storage with SQLite for scalability

**Files Created**: None (architecture planning only)

---

### Prompt 2.4 — SQL Schema
> **User**: Can you give me all the tables sql code for built db?

**What Was Done**:
- Provided full SQLite schema with 5 tables:
  - `media` — Core metadata (clean titles, TMDB/IMDb IDs, ratings, cast, genre, file paths)
  - `move_history` — Tracks `from_path`, `to_path`, `undone` status for undo feature
  - `config` — Key-value settings storage
  - `scan_history` — Logs scan operations (folder, file count, timestamp)
  - `tags` — User-defined tags per media item with `UNIQUE(media_id, tag)`

**Files Created**: None (schema documentation only)

---

## Session 3: Movie CLI Feature — Full Implementation (16-Mar-2026)

### Prompt 3.1 — Full Code Implementation
> **User**: Now create all the code for movie CLI

**What Was Done**:
- Implemented the complete Movie CLI module

**Files Created**:
| File | Purpose | Lines |
|------|---------|-------|
| `cmd/movie.go` | Parent Cobra command, registers all subcommands | ~30 |
| `cmd/movie_config.go` | `config get/set` — manages settings in SQLite | ~80 |
| `cmd/movie_scan.go` | `scan <folder>` — walks directory, cleans names, fetches TMDb, saves to DB | ~180 |
| `cmd/movie_ls.go` | `ls` — paginated library listing (20 items per page) | ~90 |
| `cmd/movie_search.go` | `search <name>` — TMDb API multi-search, interactive selection, saves to DB | ~150 |
| `cmd/movie_info.go` | `info <id-or-title>` — detail view with DB-first, TMDb fallback | ~170 |
| `cmd/movie_suggest.go` | `suggest [N]` — genre-based recommendations + trending | ~150 |
| `cmd/movie_move.go` | `move` — interactive file move with undo tracking | ~350 |
| `cmd/movie_rename.go` | `rename` — batch rename to clean format `Title (Year).ext` | ~80 |
| `cmd/movie_undo.go` | `undo` — reverts last move via `move_history` table | ~100 |
| `cmd/movie_play.go` | `play <id>` — opens file with system default player | ~60 |
| `cmd/movie_stats.go` | `stats` — genre chart, counts, average ratings | ~100 |
| `cmd/movie_resolve.go` | Shared ID/title resolver helper | ~60 |
| `db/sqlite.go` | SQLite connection, migrations, all CRUD operations | ~450 |
| `cleaner/cleaner.go` | Regex-based filename cleaning + slug generation | ~120 |
| `tmdb/client.go` | TMDb API client (search, details, credits, discover, trending) | ~200 |

**Technical Decisions**:
- Used `modernc.org/sqlite` (pure Go, no CGo dependency)
- Database stored at `./data/movie-cli.db` with WAL mode
- Thumbnails downloaded to `./data/thumbnails/`
- Move history logged to `./data/json/history/` as JSON

---

### Prompt 3.2 — Search Behavior Clarification
> **User**: I want it work like this, if I search some name use TMDb API for search the info in TMDb API and all info save in database and if it not find in TMDb then do nothing. And it store in movie or TV show folder base on movie type. It does not matter the movie or TV show does present in my library or not.

**What Was Done**:
- Rewrote `cmd/movie_search.go` to:
  1. Search TMDb API directly (not local database)
  2. Display interactive list of up to 15 results (with type, rating, year)
  3. User picks one → fetch full metadata (cast, director, genres, description)
  4. Save to local SQLite database
  5. If not found on TMDb → do nothing
  6. Media type (Movie/TV) determines categorization, not physical file existence

**Files Modified**: `cmd/movie_search.go`

---

### Prompt 3.3 — Interactive Selection Confirmation
> **User**: "If TMDb returns results → user picks one → fetch full details → save to local database" — what does "user picks one" mean?

**What Was Done**:
- Explained: CLI shows a numbered list, user types a number to select
- User confirmed this is the desired behavior ("ok its good")

**Files Modified**: None (clarification only)

---

### Prompt 3.4 — `movie ls` Scope Clarification
> **User**: movie-cli movie ls will only show my local libraries' movies and TV shows, not the databases. Read this and find out if these things are done or not.

**What Was Done**:
- Clarified: Currently `movie ls` reads from the database (entries created by `movie scan`)
- User confirmed: `movie ls` should show only items indexed via `movie scan`
- Current implementation already satisfies this (scan → DB → ls reads DB)

**Files Modified**: None

---

### Prompt 3.5 — `movie info` Full Implementation
> **User**: movie-cli movie info — Also, it will work on my Library when I use ID (firstly show how many movie and TV show present in my library, then base on ID show the result), and when I use the title, it first checks the database, if it cannot find it then it will use TMDb API and also store the data in the database.

**What Was Done**:
- Rewrote `cmd/movie_info.go` with the following priority logic:
  1. **Numeric ID** → direct lookup from local library via `GetMediaByID(id)`
  2. **Title (Local)** → search local database first using exact + prefix matching
  3. **Title (TMDb Fallback)** → query TMDb API, fetch full metadata, download poster, save to DB
  4. **Duplicate Prevention** → check `GetMediaByTmdbID` before saving
- Added `GetMediaByTmdbID(tmdbID int)` method to `db/sqlite.go`

**Files Modified**: `cmd/movie_info.go`, `db/sqlite.go`

---

## Session 4: Command Refinements & Build Tooling (17-Mar-2026)

### Prompt 4.1 — `movie move` Workflow Change
> **User**: With this "movie-cli movie move" command, I first want to choose a directory on my local machine. When I select a local directory, it shows all the movie or TV show files in it, and then I can perform a move operation.

**What Was Done**:
- Completely rewrote `cmd/movie_move.go` to a "browse-first" workflow:
  1. **Pick source directory** — choose Downloads, configured `scan_dir`, or enter custom path (or pass as argument: `movie-cli movie move ~/Videos`)
  2. **List video files** — shows all video files with cleaned titles, year, type icon, and file size
  3. **Select file** by number
  4. **Choose destination** — Movies dir, TV Shows dir, Archive dir, or custom path
  5. **Confirm & move** — shows from/to paths, asks `Are you sure? [y/N]:`
  6. **Track history** — logs to `move_history` table AND JSON history file for `movie undo` support
  7. If file isn't in DB yet, auto-inserts it

**Files Modified**: `cmd/movie_move.go` (complete rewrite, ~350 lines)

---

### Prompt 4.2 — Build Tooling
> **User**: We don't have a Makefile to run the Go code. Please fix these 2 things: add a Makefile, use gitignore to ignore binary files. Also, create a PowerShell file to do pull first, build the binary, and put the binary into e:\bin-run as the default in Windows.

**What Was Done**:
- Created `Makefile` with targets:
  - `build` — current OS
  - `build-windows` — Windows amd64
  - `build-mac-arm` — macOS ARM64
  - `build-mac-intel` — macOS amd64
  - `build-linux` — Linux amd64
  - `clean` — remove build artifacts
  - `install` — build + copy to `/usr/local/bin`
- Created `build.ps1` PowerShell script:
  1. `git pull` to fetch latest
  2. `go mod tidy`
  3. Build binary with `-ldflags="-s -w"` (optimized)
  4. Deploy to `E:\bin-run` (Windows) or `/usr/local/bin` (Mac/Linux)
  5. Clean up and verify
  - Supports custom path: `pwsh build.ps1 -BinDir "C:\custom\path"`
- Attempted `.gitignore` creation (read-only in Lovable environment, must be created manually)

**Files Created**: `Makefile`, `build.ps1`  
**Files Attempted**: `.gitignore` (read-only, user must create manually)

**Proposed `.gitignore` content**:
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

## Session 5: Specification, Memory & Quality Audit (17-Mar-2026)

### Prompt 5.1 — Full Specification Document
> **User**: Let's start implementing the specification. Before that, I want you to view the specification and tell me, if I give it to any AI, what are the failing chances? Create a detailed report first.

**What Was Done**:
- Read the entire codebase to understand all implementations
- Created comprehensive `spec.md` (350+ lines) covering:
  - Full command tree with behavior descriptions
  - SQLite database schema (all 5 tables with column definitions)
  - Cleaner regex pipeline documentation
  - TMDb API integration details (genre mapping, credit extraction)
  - Build pipeline documentation
  - Configuration keys and defaults
- Provided AI failure risk assessment:
  | Risk Area | Rating |
  |-----------|--------|
  | Web-focused AI (Lovable, v0) | 🔴 20-30% success |
  | General AI with Go support | 🟡 70-80% success |
  | Expert Go dev | 🟢 95% success |
- Identified 3 known issues:
  1. `timestamp` in move-log.json hardcoded as `"now"` (not actual time)
  2. `tags` table exists but no commands use it
  3. `DiscoverByGenre` API method exists but isn't used

**Files Created**: `spec.md`

---

### Prompt 5.2 — Memory & Tracking System
> **User**: Update Memory and Suggestions in Lovable. Include the suggestions you have, specifying what has been implemented and what is still pending. Update the issues with root cause, solution, iteration, learning.

**What Was Done**:
- Created comprehensive `.lovable/memory/` structure:

```
.lovable/memory/
├── 01-project-overview.md      # Project context, architecture, command tree
├── 02-conventions.md           # Naming, style, build, deploy preferences
├── workflow/
│   ├── 01-plan.md              # ✅ 20 completed items, 🔲 11 pending items
│   └── completed/.gitkeep
├── suggestions/
│   ├── 01-suggestions.md       # 10 suggestions ranked by priority
│   └── completed/.gitkeep
└── issues/
    ├── 01-timestamp-bug.md     # Hardcoded "now" in move-log.json
    ├── 02-duplicate-tmdb-fetch.md  # ~80 lines duplicated across 3 files
    └── 03-large-files.md       # movie_move.go + sqlite.go need splitting
```

Each issue includes: root cause, solution, impact, learning, and prevention rules.

**Files Created**: 7 new memory/tracking files

---

### Prompt 5.3 — Fix All Issues & Increase AI Success Rate
> **User**: Let's fix all the issues and increase the AI success rate to 98%. Create a plan if required.

**What Was Done**:

#### Fix 1: Timestamp Bug (Issue #01)
- **File**: `cmd/movie_move_helpers.go`
- **Change**: Replaced hardcoded `"now"` with `time.Now().Format(time.RFC3339)` in `saveHistoryLog`
- **Impact**: Move history JSON files now contain proper RFC3339 timestamps

#### Fix 2: TMDb Fetch Deduplication (Issue #02)
- **Files**: `cmd/movie_scan.go`, `cmd/movie_search.go`
- **Change**: Removed ~80 lines of duplicate TMDb fetch code; both files now call shared `fetchMovieDetails()` and `fetchTVDetails()` from `cmd/movie_info.go`
- **Impact**: Single source of truth for TMDb metadata fetching

#### Fix 3: File Refactoring (Issue #03)
- **`cmd/movie_move.go`**: Split from 348 lines into:
  - `cmd/movie_move.go` (178 lines) — main command flow
  - `cmd/movie_move_helpers.go` (168 lines) — utility functions (`promptSourceDirectory`, `promptDestination`, `listVideoFiles`, `humanSize`, `expandHome`, `saveHistoryLog`)
- **`db/sqlite.go`**: Split from 452 lines into 5 focused files:
  - `db/db.go` — connection, migrations, initialization
  - `db/media.go` — Media CRUD operations (insert, update, get, list, search)
  - `db/config.go` — Config get/set operations
  - `db/history.go` — Move history + scan history operations
  - `db/helpers.go` — String utility functions

**Files Deleted**: `db/sqlite.go` (replaced by 5 new files)

#### AI Success Plan
- Created `.lovable/memory/workflow/01-ai-success-plan.md` with 7 rules:
  1. Always read `01-project-overview.md` first
  2. Ignore `package.json` and `dev` server errors (Go project)
  3. Never create web/React/Node.js files
  4. Follow one-file-per-command pattern
  5. Use shared helpers (DRY principle)
  6. Keep files under 200 lines
  7. Document all fixes in the issue tracking system

---

### Prompt 5.4 — Requirements Reliability Analysis
> **User**: Read the full set of requirements and specifications. Assess whether the requirements are complete enough for another AI to successfully build the application.

**What Was Done**:
- Performed comprehensive reliability analysis
- **Success Probability Estimates**:
  - Simple isolated features: 85–90%
  - Medium features with DB + UI: 75–85%
  - Complex multi-step workflows: 60–70%
  - Full system integration: 55–65%
- **Key Gaps Identified**:
  1. `spec.md` still references deleted `db/sqlite.go` (stale)
  2. Appendix B shows hardcoded `"now"` timestamp (stale)
  3. No error handling specs (filesystem failures, TMDb rate limits, DB locks)
  4. No testable acceptance criteria (GIVEN/WHEN/THEN)
- **Recommendations**:
  - Consolidate cross-cutting rules into a central section
  - Add application-level spec in `spec/01-app/01-project-spec.md`
  - Update all stale references

**Files Created**: `spec/01-app/01-project-spec.md`, `spec/02-app/issues/04-wrong-project-context.md`

---

### Prompt 5.5 — Update spec.md
> **User**: Base on spec.md, if need to update spec files, then do.

**What Was Done**:
- Fixed 3 stale areas in `spec.md`:
  1. Project structure now shows the 5-file `db/` layout and `movie_move_helpers.go`
  2. Added `spec/` folder to the file tree
  3. Appendix B now shows RFC3339 timestamp instead of hardcoded `"now"`

**Files Modified**: `spec.md`

---

## Session 6: Code Hygiene Fixes & Documentation (17-18-Mar-2026)

### Prompt 6.1 — Full Codebase Audit
> **User**: Read full project, if anything makes you annoyed or irrelevant then inform me. I want to create Movie CLI, does it stay on its focus?

**What Was Done**:
- Read every file in the project
- **Confirmed**: Project is 100% focused on Movie CLI movie management — no scope creep
- **Found 8 hygiene issues**:
  1. `go.mod` missing `modernc.org/sqlite` dependency
  2. `db/helpers.go` reimplements `strings.Split`, `strings.Index`, `strings.TrimSpace` (46 unnecessary lines)
  3. `main.go` line 2 has stray `//Ab` comment
  4. `cmd/movie_scan.go` lines 175-181 have broken indentation
  5. `cmd/movie_info.go` creates redundant `Media` copy (loses `ImdbRating` field)
  6. `cmd/movie_resolve.go` — `resolveMediaByQuery` is never called by any command
  7. `readm.txt` and `readme.txt` both exist with different dates (redundant)
  8. Compiled `movie-cli` binary committed to repo

**Files Modified**: None (audit only)

---

### Prompt 6.2 — Fix All 8 Issues
> **User**: Fix all 8 issues found: remove db/helpers.go reimplementations, fix stray comment in main.go, fix indentation in movie_scan.go, remove redundant Media copy in movie_info.go, wire up resolveMediaByQuery or remove it, consolidate readm.txt/readme.txt, and update go.mod

**What Was Done**:

| # | Issue | Fix | File(s) |
|---|-------|-----|---------|
| 1 | `go.mod` missing SQLite dep | Added `modernc.org/sqlite v1.29.5` | `go.mod` |
| 2 | `db/helpers.go` reimplements stdlib | Replaced with `strings.Split`, `strings.TrimSpace` (46→13 lines) | `db/helpers.go` |
| 3 | Stray `//Ab` comment | Removed | `main.go` |
| 4 | Broken indentation | Fixed alignment | `cmd/movie_scan.go` |
| 5 | Redundant `Media` copy | Removed; passes `m` directly to `printMediaDetail` | `cmd/movie_info.go` |
| 6 | Unused `resolveMediaByQuery` | Wired into `movie_info.go` for title lookups; removed duplicate `pickBestMatch` | `cmd/movie_info.go`, `cmd/movie_resolve.go` |
| 7 | Duplicate readme files | Deleted `readme.txt`; `readm.txt` is the sole milestone file | `readme.txt` (deleted) |
| 8 | Binary committed to repo | Deleted `movie-cli` and legacy binaries binaries | `movie-cli` binaries (deleted) |

**Files Modified**: `go.mod`, `db/helpers.go`, `main.go`, `cmd/movie_scan.go`, `cmd/movie_info.go`  
**Files Deleted**: `readme.txt`, `movie-cli` binaries  
**Files Created**: `spec/02-app/issues/05-code-hygiene-audit.md`

---

### Prompt 6.3 — Verification
> **User**: Read the full codebase again to verify all 8 fixes are clean and no new issues were introduced.

**What Was Done**:
- Re-read all modified files
- Verified all 8 fixes are clean
- Found and deleted one additional binary (`movie-cli`) that was still in the repo
- **Result**: No new issues introduced

---

### Prompt 6.4 — README Update
> **User**: Base on the project, update README.md file.

**What Was Done**:
- Rewrote `README.md` with comprehensive documentation:
  - Feature list (scan, search, list, move, rename, undo, play, suggest, stats)
  - Full command tree with descriptions
  - Quick start guide (build, configure, scan, browse, search)
  - Configuration table (6 keys with defaults)
  - Project structure (modularized `cmd/` and `db/` directories)
  - Build targets (make build, cross-compile for Windows/Mac/Linux)
  - Dependencies table (Cobra, modernc.org/sqlite)
  - Data storage layout (`./data/`)

**Files Modified**: `README.md`

---

## Milestone Markers (readm.txt History)

The user uses `readm.txt` as session start markers. All entries below were requested:

| # | Date | Time (Malaysia UTC+8) | Prompt # |
|---|------|-----------------------|----------|
| 1 | 15-Mar-2026 | 08:00 AM | 1.1 |
| 2 | 15-Mar-2026 | 04:00 PM | 1.2 |
| 3 | 16-Mar-2026 | 11:00 PM | 1.3 |
| 4 | 17-Mar-2026 | 12:10 PM | — |
| 5 | 17-Mar-2026 | 08:40 PM | — |
| 6 | 17-Mar-2026 | 08:30 PM | — |
| 7 | 17-Mar-2026 | 09:45 PM | — |
| 8 | 18-Mar-2026 | 12:00 AM | — |
| 9 | 18-Mar-2026 | 12:45 AM | — |
| 10 | 18-Mar-2026 | 01:05 AM | — |

---

## Current Project State (as of 05-Apr-2026)

### ✅ Completed Features
1. CLI framework (Cobra) with `hello`, `version`, `self-update`
2. `movie config` — get/set settings in SQLite
3. `movie scan` — walk directory, clean names, fetch TMDb, save to DB
4. `movie ls` — paginated library listing (20 items per page)
5. `movie search` — TMDb API search, interactive selection, save to DB
6. `movie info` — ID/title lookup with DB-first, TMDb fallback, auto-persist
7. `movie suggest` — genre-based recommendations + trending
8. `movie move` — browse-first file move with undo tracking
9. `movie rename` — batch rename to `Title (Year).ext`
10. `movie undo` — revert last move via `move_history` table
11. `movie play` — open with system default player
12. `movie stats` — genre chart, counts, average ratings
13. `Makefile` — cross-platform build targets
14. `build.ps1` — PowerShell build + deploy script
15. `spec.md` — comprehensive project specification
16. Full `.lovable/memory/` tracking system
17. All 3 code issues fixed (timestamp, deduplication, file splitting)
18. All 8 hygiene issues fixed (helpers, comments, indentation, go.mod, etc.)

### 🔲 Pending Features
1. `movie tag` command (add/remove/list tags — `tags` table already exists)
2. `movie undo` confirmation prompt before reverting
3. Cross-drive move support (`os.Rename` fallback to copy+delete)
4. Batch move (`--all` flag for moving all video files at once)
5. File size stats in `movie stats`
6. `.gitignore` file (must be created manually — read-only in Lovable)

### Architecture Summary

```
movie-cli/
├── main.go                        # Entry point
├── cmd/
│   ├── root.go                    # Root command
│   ├── hello.go, version.go       # Utility commands
│   ├── update.go                  # Self-update (git pull)
│   ├── movie.go                   # Parent: movie-cli movie
│   ├── movie_config.go            # config get/set
│   ├── movie_scan.go              # scan folder → DB + TMDb
│   ├── movie_ls.go                # paginated list
│   ├── movie_search.go            # TMDb search → DB
│   ├── movie_info.go              # detail view + shared TMDb helpers
│   ├── movie_suggest.go           # recommendations + trending
│   ├── movie_move.go              # interactive move
│   ├── movie_move_helpers.go      # move utility functions
│   ├── movie_rename.go            # batch rename
│   ├── movie_undo.go              # undo last move
│   ├── movie_play.go              # open with default player
│   ├── movie_stats.go             # library statistics
│   └── movie_resolve.go           # shared ID/title resolver
├── cleaner/cleaner.go             # Filename cleaning + slug generation
├── tmdb/client.go                 # TMDb API client
├── db/
│   ├── db.go                      # SQLite connection + migrations
│   ├── media.go                   # Media CRUD
│   ├── config.go                  # Config get/set
│   ├── history.go                 # Move + scan history
│   └── helpers.go                 # String utilities
├── updater/updater.go             # Git-based self-update
├── version/version.go             # Build-time version variables
├── Makefile                       # Build targets
├── build.ps1                      # PowerShell build + deploy
├── spec.md                        # Full specification
├── README.md                      # Project documentation
├── readm.txt                      # Session milestone marker
└── spec/
    ├── 01-app/01-project-spec.md  # Application spec
    └── 02-app/issues/             # Issue write-ups (01-05)
```

### Database Schema

```sql
-- Core metadata
CREATE TABLE media (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT, clean_title TEXT, slug TEXT,
    year INTEGER, media_type TEXT,
    tmdb_id INTEGER UNIQUE, imdb_id TEXT,
    overview TEXT, genre TEXT, cast TEXT, director TEXT,
    tmdb_rating REAL, imdb_rating REAL, vote_count INTEGER,
    poster_path TEXT, thumbnail_local TEXT,
    file_path TEXT, file_size INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Move history for undo
CREATE TABLE move_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    media_id INTEGER, from_path TEXT, to_path TEXT,
    moved_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    undone INTEGER DEFAULT 0
);

-- Key-value config
CREATE TABLE config (key TEXT PRIMARY KEY, value TEXT);

-- Scan logs
CREATE TABLE scan_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    folder TEXT, file_count INTEGER,
    scanned_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- User tags
CREATE TABLE tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    media_id INTEGER, tag TEXT,
    UNIQUE(media_id, tag)
);
```

### Key Technical Decisions
- **Pure Go SQLite**: `modernc.org/sqlite` — no CGo dependency, cross-platform
- **WAL mode**: Enabled for better concurrent access
- **Cobra**: CLI framework for command structure
- **TMDb API**: Metadata source (requires API key via `movie config set tmdb_api_key`)
- **Data directory**: `./data/` (DB, thumbnails, JSON history)

---

## For Any AI Continuing This Project

1. **This is a Go CLI project** — NOT a web app. Ignore any `package.json`, `dev` server, or React-related errors.
2. **Read `spec/01-app/01-project-spec.md`** first for coding guidelines and conventions.
3. **One file per command**: `cmd/movie_<cmd>.go` — max ~200 lines.
4. **Use shared helpers**: `fetchMovieDetails()`/`fetchTVDetails()` from `movie_info.go`, `resolveMediaByQuery()` from `movie_resolve.go`.
5. **Never hardcode placeholders**: No `"now"`, `"TODO"`, `"test"` in production paths.
6. **Document every fix**: Create issue write-up in `spec/02-app/issues/`, update spec, update memory.
7. **Build locally**: `make build` or `pwsh build.ps1` — cannot build in Lovable's web environment.
