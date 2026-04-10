# Issue #03: Large Files Exceeding 200 Lines

## Issue Summary

1. **What happened**: `cmd/movie_move.go` (348 lines) and `db/sqlite.go` (452 lines) grew too large, causing AI context loss and partial edits.
2. **Where**: `cmd/movie_move.go`, `db/sqlite.go`.
3. **Symptoms and impact**: AI made edits in wrong locations, missed related code, produced broken imports. Harder to navigate for contributors.
4. **How discovered**: AI success rate analysis — 15% of failures traced to large file context loss.

## Root Cause Analysis

1. **Direct cause**: Features added incrementally without splitting files at natural boundaries.
2. **Contributing factors**: No file size guideline existed. Single-file-per-package was the default Go pattern used.
3. **Triggering conditions**: Files exceeded ~250 lines, causing AI to lose track of function relationships.
4. **Why spec did not prevent it**: No maximum file size rule existed in conventions or spec.

## Fix Description

1. **Spec change**: Added "Max ~200 lines per file" rule to `spec/08-app/01-project-spec.md` and conventions.
2. **New rule**: Split files at natural boundaries when approaching 200 lines. One concern per file.
3. **Why it resolves root cause**: Smaller files keep full context visible to AI and human readers.
4. **Config changes**: None.
5. **Diagnostics**: None required.

## Iterations History

1. **Iteration 1** (17-Mar-2026): Split `cmd/movie_move.go` → `movie_move.go` (178) + `movie_move_helpers.go` (168). Split `db/sqlite.go` → 5 files (`db.go`, `media.go`, `config.go`, `history.go`, `helpers.go`). Fixed on first attempt.

## Prevention and Non-Regression

1. **Prevention rule**: Never let a file grow past ~200 lines without evaluating a split.
2. **Acceptance criteria**: No Go source file in the project exceeds 250 lines.
3. **Guardrails**: AI success plan Rule 2 (one file, one concern, <200 lines).
4. **Spec references**: `spec/08-app/01-project-spec.md` (File Structure Rules).

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
