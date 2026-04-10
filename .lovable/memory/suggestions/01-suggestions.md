# Suggestions Tracker

> **Last Updated**: 05-Apr-2026

## Status Legend
- ✅ Done — implemented and verified
- 🔲 Open — not started
- 🔄 In Progress — actively being worked on

---

## ✅ Completed

| # | Suggestion | Completed | Notes |
|---|-----------|-----------|-------|
| S01 | Fix timestamp bug in move-log.json | 17-Mar-2026 | Replaced `"now"` with `time.Now().Format(time.RFC3339)` |
| S02 | Refactor large files (>200 lines) | 17-Mar-2026 | Split `movie_move.go` and `db/sqlite.go` |
| S03 | Extract shared TMDb fetch logic | 17-Mar-2026 | `fetchMovieDetails()`/`fetchTVDetails()` in `movie_info.go` |

---

## 🔲 Open — Priority Order

### P0 — Must Fix Before Handoff

| # | Suggestion | Affected | Rationale |
|---|-----------|----------|-----------|
| S04 | Add cross-drive move fallback (copy+delete) | `cmd/movie_move.go` | `os.Rename` fails across filesystems; silent data inconsistency |
| S05 | Add confirmation prompt to `movie undo` | `cmd/movie_undo.go` | Destructive operation without safety net |

### P1 — High Priority

| # | Suggestion | Affected | Rationale |
|---|-----------|----------|-----------|
| S06 | Add GIVEN/WHEN/THEN acceptance criteria to spec.md | `spec.md` §4 | AI cannot self-validate without testable criteria |
| S07 | Document shared helper locations in code comments | `cmd/movie_info.go`, `cmd/movie_resolve.go` | Prevent duplicate code creation by AI |
| S08 | Clarify `movie ls` filter rule (scan-only items) | `cmd/movie_ls.go` | Ambiguous whether all DB records or only file-backed records shown |

### P2 — Medium Priority

| # | Suggestion | Affected | Rationale |
|---|-----------|----------|-----------|
| S09 | Implement `movie tag` command | New: `cmd/movie_tag.go` | `tags` table exists but no commands use it |
| S10 | Add file size stats to `movie stats` | `cmd/movie_stats.go` | Total size, average size, largest file |
| S11 | Add error handling spec (TMDb rate limits, DB locks, offline) | `spec/08-app/` | No error handling documentation |
| S12 | Update README.md with full feature documentation | `README.md` | Currently only documents hello/version/self-update |

### P3 — Low Priority

| # | Suggestion | Affected | Rationale |
|---|-----------|----------|-----------|
| S13 | Add batch move (`--all` flag) | `cmd/movie_move.go` | Move all video files from source at once |
| S14 | Write JSON metadata per movie/TV on scan | `cmd/movie_scan.go` | Storage structure promises this but it's not implemented |
| S15 | Use `DiscoverByGenre` TMDb method in suggest | `cmd/movie_suggest.go` | Method exists in client but is unused |

---

*Tracker updated: 05-Apr-2026 — Remove entries when done, add completion date to Completed table.*
