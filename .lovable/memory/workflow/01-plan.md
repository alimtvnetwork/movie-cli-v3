# Project Plan & Status

> **Last Updated**: 09-Apr-2026

## ✅ Completed

### Core CLI Structure
- [x] Root command with Cobra (`movie-cli`)
- [x] `hello` command with version display
- [x] `version` command with ldflags injection
- [x] `self-update` command via git pull --ff-only

### Movie Management Commands
- [x] `movie config` — get/set configuration with masked API key display
- [x] `movie scan` — folder scanning with TMDb metadata + poster download
- [x] `movie ls` — paginated list with interactive navigation + detail view
- [x] `movie search` — live TMDb search, select, save to DB
- [x] `movie info` — local DB lookup → TMDb fallback → auto-persist
- [x] `movie suggest` — genre-based recommendations + trending fallback
- [x] `movie move` — interactive browse, move, track history
- [x] `movie rename` — batch clean rename with undo tracking
- [x] `movie undo` — revert last move/rename operation
- [x] `movie play` — open file with system default player (cross-platform)
- [x] `movie stats` — counts, genre chart, average ratings

### Infrastructure
- [x] SQLite database with migrations (5 tables, 7 indexes)
- [x] TMDb API client (search, details, credits, recommendations, trending, posters)
- [x] Filename cleaner (junk removal, year extraction, TV detection, slugs)
- [x] Makefile with build + cross-compile targets
- [x] build.ps1 PowerShell deploy script
- [x] spec.md — full project specification
- [x] Shared resolver helper (`movie_resolve.go`)

### Bug Fixes
- [x] Fixed timestamp bug — `saveHistoryLog` now uses `time.Now().Format(time.RFC3339)`
- [x] Deduplicated TMDb fetch logic — shared `fetchMovieDetails()`/`fetchTVDetails()`

### Refactoring
- [x] Split `cmd/movie_move.go` → `movie_move.go` + `movie_move_helpers.go`
- [x] Split `db/sqlite.go` → 5 focused files

### Documentation
- [x] README.md (basic), spec.md, ai-handoff.md, development-log.md
- [x] .lovable/memory structure with suggestions, issues, workflow
- [x] AI success rate plan
- [x] Reliability risk report (05-Apr-2026)

### Spec Restructuring (Phase 1-5)
- [x] Phase 1: Spec authoring guideline review
- [x] Phase 2: Spec folder audit
- [x] Phase 3: Naming/placement normalization (root lowercase, merge 02-app, flatten error spec)
- [x] Phase 4: Ignore rule verification (.gitignore audit)
- [x] Phase 5: Final consistency pass (N1-N4 renames, C1-C5 missing files created)

### PowerShell Automation (Phase 1-8) ✅
- [x] Phase 1: Core parameters & environment detection
- [x] Phase 2: Git operations (pull, conflict resolution, force-pull)
- [x] Phase 3: Go build pipeline integration
- [x] Phase 4: Deployment with backup & rollback
- [x] Phase 5: Logging & colored output helpers
- [x] Phase 6: Error handling audit (no swallowed errors)
- [x] Phase 7: install.ps1 bootstrap script
- [x] Phase 8: README.md automation docs update

### Release Pipeline & Install Scripts (09-Apr-2026) ✅
- [x] GitHub Actions release.yml — triggers on `release/**` branches and `v*` tags
- [x] Cross-compiled binaries for 6 targets (windows/linux/darwin × amd64/arm64)
- [x] Version-specific install.ps1 (Windows) with SHA256 verification
- [x] Version-specific install.sh (Linux/macOS) with SHA256 verification
- [x] Release page with changelog, checksums, install instructions, asset table
- [x] CHANGELOG.md created for release note extraction
- [x] Pipeline spec documentation (`spec/pipeline/`)

### CLI UX Improvements (09-Apr-2026) ✅
- [x] Root command shows version + comprehensive help with examples
- [x] `mahin --version` flag support
- [x] `mahin version` shows Go version and OS/arch
- [x] `mahin movie` shows categorized subcommand help with examples

---

## 🔲 Pending — Prioritized Backlog

### Phase 1: Safety & Reliability (P0)
- [ ] `movie undo` confirmation prompt before reverting

### Phase 2: Spec Completeness (P1)
- [x] Clarify `movie ls` filter rule (scan-indexed items only) ✅ 09-Apr-2026

### Phase 3: Enhancements (P3)
- [x] Batch move (`--all` flag for `movie move`) ✅ 09-Apr-2026
- [x] JSON metadata files per movie/TV show on scan ✅ 09-Apr-2026
- [x] Use `DiscoverByGenre` in `movie suggest` ✅ 09-Apr-2026
- [ ] CI pipeline (lint, test, vuln scan) — follow gitmap-v2 pattern

---

## Next Task Selection

Pick one of these to implement next:

1. **Undo confirmation prompt** — Add `[y/N]` prompt before reverting. Affects `cmd/movie_undo.go`.
2. **Movie ls filter clarification** — Document that only file-backed items show.
3. **CI pipeline** — Add `.github/workflows/ci.yml` with lint, test, vuln scan.
4. **Batch move** — Add `--all` flag to `movie move`.

*Tell me which task to implement.*
