# Suggestions Tracker

> **Last Updated**: 10-Apr-2026

## Status Legend
- ✅ Done — implemented and verified
- 🔲 Open — not started

---

## ✅ Completed

| # | Suggestion | Completed | Notes |
|---|-----------|-----------|-------|
| S01 | Fix timestamp bug in move-log.json | 17-Mar-2026 | Replaced `"now"` with `time.Now().Format(time.RFC3339)` |
| S02 | Refactor large files (>200 lines) | 17-Mar-2026 | Split `movie_move.go` and `db/sqlite.go` |
| S03 | Extract shared TMDb fetch logic | 17-Mar-2026 | `fetchMovieDetails()`/`fetchTVDetails()` in `movie_info.go` |
| S04 | Cross-drive move fallback (copy+delete) | 05-Apr-2026 | `MoveFile()` detects EXDEV, falls back to copy+remove |
| S08 | Clarify `movie ls` filter rule | 09-Apr-2026 | Only file-backed (scanned) items shown |
| S09 | Implement `movie tag` command | 06-Apr-2026 | `cmd/movie_tag.go` + `db/tags.go` |
| S13 | Batch move (`--all` flag) | 09-Apr-2026 | Move all video files from source at once |
| S14 | JSON metadata per movie/TV on scan | 09-Apr-2026 | `cmd/movie_scan_json.go` |
| S15 | Use `DiscoverByGenre` in suggest | 09-Apr-2026 | Genre-based discovery integrated |
| S05 | Add confirmation prompt to `movie undo` | 10-Apr-2026 | Already implemented with `[y/N]` prompt |
| S16 | CI pipeline (lint, test, vuln scan) | 10-Apr-2026 | ci.yml + vulncheck.yml + spec/12-ci-cd-pipeline/ |
| S06 | Add GIVEN/WHEN/THEN acceptance criteria | 10-Apr-2026 | 16 ACs covering all commands + export + batch move |
| S07 | Document shared helper locations | 10-Apr-2026 | Annotated movie_info.go, movie_resolve.go, movie_move_helpers.go, movie_scan_json.go |
| S12 | Update README.md with full docs | 10-Apr-2026 | 620+ lines, all commands, install, build, project structure |
| S11 | Add error handling spec | 10-Apr-2026 | TMDb rate limits, DB locks, offline mode, filesystem errors |
| S10 | Add file size stats to `movie stats` | 10-Apr-2026 | Already implemented: total, largest, smallest, average in CLI + total in web UI |
| S18 | Add `movie duplicates` command | 10-Apr-2026 | Detect by TMDb ID, filename, or file size — db/duplicates.go + cmd/movie_duplicates.go |
| S21 | Apply error log spec v2 to ci.yml | 10-Apr-2026 | clear-logs job, per-stage .github/logs/errors/*.log, summary.log assembly |
| S19 | Add `movie cleanup` command | 10-Apr-2026 | Find/remove stale DB entries — db/cleanup.go + cmd/movie_cleanup.go |
| S22 | Add `movie watch` / watchlist | 11-Apr-2026 | to-watch/watched tracking — db/watchlist.go + cmd/movie_watch.go |
| S17 | Retry logic with exponential backoff | 11-Apr-2026 | 429 rate-limit handling in tmdb/client.go get() — 3 retries, 1s/2s/4s backoff |

---

## 🔲 Open — Priority Order

### P3 — Low Priority

| # | Suggestion | Affected | Rationale |
|---|-----------|----------|-----------|
| S20 | Integration tests for CI | `tests/`, `ci.yml` | No test job with SQLite fixtures |

---

*Tracker updated: 11-Apr-2026*
