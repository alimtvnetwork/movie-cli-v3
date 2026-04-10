# Issue: Hardcoded Timestamp in move-log.json

> **Status**: ✅ Resolved  
> **Severity**: Medium  
> **File**: `cmd/movie_move_helpers.go` (was `cmd/movie_move.go` line 345-346)  
> **Iteration**: 1 (fixed 17-Mar-2026)

## Root Cause

The `saveHistoryLog` function wrote `"timestamp":"now"` as a literal string instead of an actual timestamp.

## Solution Applied

Replaced `"now"` with `time.Now().Format(time.RFC3339)` and added `"time"` to imports.
Function moved to `movie_move_helpers.go` during refactoring.

## Impact (Before Fix)

- All move history JSON logs had useless timestamp data
- Could not reconstruct when moves happened from JSON logs
- DB `move_history.moved_at` was correct (uses `CURRENT_TIMESTAMP`), so DB was unaffected

## Learning

- Always use actual time functions, never placeholder strings
- Review all format strings for hardcoded test values before committing

## What Not to Repeat

- Don't use placeholder strings (`"now"`, `"TODO"`, `"test"`) in production code paths
- Always grep for placeholder strings before marking a feature complete
