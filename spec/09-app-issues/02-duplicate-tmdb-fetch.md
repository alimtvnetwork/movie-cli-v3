# Issue #02: Duplicate TMDb Fetch Logic

## Issue Summary

1. **What happened**: Three commands (`scan`, `search`, `info`) each contained nearly identical code (~80 lines each) for fetching movie/TV details + credits from TMDb.
2. **Where**: `cmd/movie_scan.go`, `cmd/movie_search.go`, `cmd/movie_info.go`.
3. **Symptoms and impact**: Bug fixes had to be applied in 3 places. Risk of behavior divergence between commands. ~80 lines of unnecessary duplicate code.
4. **How discovered**: Code review during refactoring pass.

## Root Cause Analysis

1. **Direct cause**: Features were built incrementally — each command copied the fetch pattern from the previous one.
2. **Contributing factors**: No shared helper existed when `scan` was written. `search` and `info` copied from `scan` instead of extracting.
3. **Triggering conditions**: Any change to TMDb fetch logic (e.g., adding a new field) required updating 3 files.
4. **Why spec did not prevent it**: `01-project-spec.md` described each command independently without noting shared patterns.

## Fix Description

1. **Spec change**: Added DRY principle rule and shared function references to `spec/08-app/01-project-spec.md`.
2. **New rule**: TMDb fetch logic MUST use `fetchMovieDetails()` / `fetchTVDetails()` from `movie_info.go`. Never copy-paste these blocks.
3. **Why it resolves root cause**: Single source of truth for TMDb fetch. Changes propagate automatically.
4. **Config changes**: None.
5. **Diagnostics**: None required.

## Iterations History

1. **Iteration 1** (17-Mar-2026): Extracted `fetchMovieDetails()` and `fetchTVDetails()` in `movie_info.go`. Updated `scan` and `search` to call shared helpers. Reduced ~80 lines per file. Fixed on first attempt.

## Prevention and Non-Regression

1. **Prevention rule**: Before writing new TMDb-related code, check `movie_info.go` for existing helpers. Extract to shared functions BEFORE duplicating into a second command.
2. **Acceptance criteria**: No file in `cmd/` should contain inline TMDb detail+credit fetch code except `movie_info.go`.
3. **Guardrails**: AI success plan Rule 3 (shared logic = shared functions).
4. **Spec references**: `spec/08-app/01-project-spec.md` (DRY Principle section).

## TODO and Follow-ups

- [x] Code fix applied
- [x] Spec updated
- [x] Memory updated

## Done Checklist

- [x] Spec updated under `/spec/08-app/`
- [x] Issue write-up created under `/spec/09-app-issues/`
- [x] Memory updated with summary and prevention rule
- [x] Acceptance criteria updated
- [x] Iterations recorded
