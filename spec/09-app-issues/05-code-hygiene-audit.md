# Issue #05: Code Hygiene — 8 Quality Issues

## Issue Summary

1. **What happened**: Audit of the full codebase found 8 quality/hygiene issues.
2. **Where**: `db/helpers.go`, `main.go`, `cmd/movie_scan.go`, `cmd/movie_info.go`, `cmd/movie_resolve.go`, `readm.txt`/`readme.txt`, `go.mod`, root `movie-cli` binary.
3. **Symptoms and impact**: Redundant code, inconsistent formatting, unused functions, duplicate files, incomplete dependency list, committed binary.
4. **How discovered**: Full project read and audit.

## Root Cause Analysis

1. **Direct cause**: Incremental development without periodic cleanup passes.
2. **Contributing factors**: No linting or CI pipeline to catch issues automatically.
3. **Triggering conditions**: Accumulated over multiple development sessions.
4. **Why spec did not prevent it**: No automated code quality checks specified.

## Fix Description

1. **`db/helpers.go`**: Replaced hand-rolled `split()`, `indexOf()`, `trim()` (46 lines) with `strings.Split`, `strings.TrimSpace` (13 lines).
2. **`main.go`**: Removed stray `//Ab` comment on line 2.
3. **`cmd/movie_scan.go`**: Fixed indentation of `fetchMovieDetails`/`fetchTVDetails` block (lines 175-181) to align with surrounding code.
4. **`cmd/movie_info.go`**: Removed redundant field-by-field `Media` copy at line 159-173; now passes `m` directly to `printMediaDetail`.
5. **`cmd/movie_info.go`**: Removed unused `pickBestMatch` function; `runMovieInfo` now uses shared `resolveMediaByQuery` from `movie_resolve.go`.
6. **`readme.txt`**: Deleted duplicate file; `readm.txt` is the single milestone marker.
7. **`go.mod`**: Added missing `modernc.org/sqlite v1.29.5` to require block.
8. **`movie-cli` binary**: Deleted committed binary from repo root.

## Prevention and Non-Regression

1. **Prevention rule**: Run `go vet ./...` mentally before marking done. Use stdlib functions instead of reimplementing.
2. **Acceptance criteria**: No unused functions, no duplicate files, no committed binaries, `go.mod` lists all imports.
3. **Spec references**: `spec/08-app/01-project-spec.md` (DRY, file structure rules).

## Done Checklist

- [x] Spec updated under `/spec/08-app/`
- [x] Issue write-up created under `/spec/09-app-issues/`
- [x] Memory updated
- [x] All 8 fixes applied
