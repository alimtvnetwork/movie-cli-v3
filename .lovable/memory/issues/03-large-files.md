# Issue: Large Files Need Refactoring

> **Status**: ✅ Resolved  
> **Severity**: Low (maintainability)  
> **Files**: `cmd/movie_move.go` (was 348 lines), `db/sqlite.go` (was 452 lines)  
> **Iteration**: 1 (fixed 17-Mar-2026)

## Root Cause

Features were added incrementally without splitting files at natural boundaries.

## Solution Applied

### `cmd/movie_move.go` → split into:
- `movie_move.go` (178 lines) — command definition + `runMovieMove` main flow
- `movie_move_helpers.go` (168 lines) — `promptSourceDirectory`, `promptDestination`, `promptCustomPath`, `listVideoFiles`, `humanSize`, `expandHome`, `saveHistoryLog`

### `db/sqlite.go` → split into:
- `db/db.go` — `DB` struct, `Open()`, `migrate()` (schema + defaults)
- `db/media.go` — `Media` struct, all CRUD methods, `scanMediaRows`, `TopGenres`
- `db/config.go` — `GetConfig`, `SetConfig`
- `db/history.go` — `MoveRecord`, `InsertMoveHistory`, `GetLastMove`, `MarkMoveUndone`, `InsertScanHistory`
- `db/helpers.go` — `splitCSV`, `split`, `indexOf`, `trim`

## Impact (Before Fix)

- Harder to navigate and understand for new contributors/AI
- Higher risk of merge conflicts
- AI lost context in large files, causing partial/broken edits

## Learning

- Split files at natural module boundaries early
- One "concern" per file: schema, queries per entity, helpers
- Target ~200 lines max per file

## What Not to Repeat

- Don't let files grow past ~200 lines without evaluating a split
