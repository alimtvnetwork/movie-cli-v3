# Issue: Duplicate TMDb Fetch Logic

> **Status**: ✅ Resolved  
> **Severity**: Low (code quality)  
> **Files**: `cmd/movie_scan.go`, `cmd/movie_search.go`, `cmd/movie_info.go`  
> **Iteration**: 1 (fixed 17-Mar-2026)

## Root Cause

Three commands (`scan`, `search`, `info`) each contained nearly identical code for fetching movie/TV details + credits from TMDb.

## Solution Applied

Refactored `scan` and `search` to call the shared `fetchMovieDetails()` and `fetchTVDetails()` functions from `movie_info.go`. Reduced ~80 lines of duplicate code to 6 lines of function calls.

## Impact (Before Fix)

- ~80 lines of duplicate code across 3 files
- Bug fixes had to be applied in 3 places
- Risk of behavior divergence between commands

## Learning

- Extract shared logic BEFORE duplicating into a second command
- When adding a third copy of the same pattern, stop and refactor immediately

## What Not to Repeat

- Don't copy-paste TMDb fetch blocks — always use the shared helpers
- When adding new metadata fields, update the shared functions, not individual commands
