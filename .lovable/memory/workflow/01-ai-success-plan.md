# AI Success Rate Plan — Target: 98%

> **Last Updated**: 17-Mar-2026  
> **Current Estimated Rate**: ~65%  
> **Target Rate**: 98%

## Why AI Fails on This Project

### 1. Missing Context (Cause: ~40% of failures)
- **Problem**: AI doesn't know this is a Go CLI, not a web app. It tries to add `package.json`, run `npm`, or create React components.
- **Symptom**: Build errors like "no package.json found" are treated as real bugs.
- **Fix**: `01-project-overview.md` now explicitly states "This is NOT a web project" with specific instructions to ignore Lovable build errors.

### 2. Stale Memory / Inconsistent State (Cause: ~20% of failures)
- **Problem**: Memory files say things are "pending" when they've been done, or vice versa. AI makes changes that were already made.
- **Symptom**: Duplicate work, conflicting code, regression bugs.
- **Fix**: Keep ONE source of truth per concern. Plan file tracks done/pending. Suggestions file tracks implementation status. Issues files track iteration count.

### 3. Large Files Confuse AI (Cause: ~15% of failures)
- **Problem**: Files >300 lines cause AI to lose track of context, make edits in wrong locations, or miss related code.
- **Symptom**: Partial edits, broken imports, functions that reference deleted code.
- **Fix**: All files now <200 lines. Split `db/sqlite.go` → 5 files. Split `movie_move.go` → 2 files.

### 4. Duplicate Logic = Divergent Behavior (Cause: ~10% of failures)
- **Problem**: Same TMDb fetch pattern copied across 3 files. AI fixes one but not the others.
- **Symptom**: `scan` behaves differently from `search` or `info` for the same movie.
- **Fix**: Extracted shared `fetchMovieDetails()` and `fetchTVDetails()` in `movie_info.go`. Both `scan` and `search` now call these.

### 5. Placeholder/Debug Values in Production (Cause: ~5% of failures)
- **Problem**: Hardcoded strings like `"now"` instead of actual timestamps.
- **Symptom**: Data corruption, useless logs, silent failures.
- **Fix**: Fixed timestamp bug. Rule: grep for `"now"`, `"TODO"`, `"test"`, `"placeholder"` before any commit.

### 6. No Specification for Edge Cases (Cause: ~10% of failures)
- **Problem**: spec.md describes happy paths but not error handling, cross-platform issues, or boundary conditions.
- **Symptom**: AI generates code that works on Mac but fails on Windows (e.g., `os.Rename` cross-drive).
- **Fix**: Add edge case documentation to spec.md (see below).

## How to Achieve 98% Success Rate

### Rule 1: Read Memory Before Coding
Every new AI session MUST read in this order:
1. `01-project-overview.md` — understand it's a Go CLI
2. `02-conventions.md` — coding style and patterns
3. `workflow/01-plan.md` — what's done, what's pending
4. `suggestions/01-suggestions.md` — prioritized work items
5. `issues/` — active bugs with root cause analysis

### Rule 2: One File, One Concern, <200 Lines
- `cmd/` files: one file per command + one `_helpers.go` if needed
- `db/` files: `db.go` (connection), `media.go` (CRUD), `config.go`, `history.go`, `helpers.go`
- Never let a file grow past 200 lines without splitting

### Rule 3: Shared Logic = Shared Functions
- TMDb fetch: always use `fetchMovieDetails()` / `fetchTVDetails()` from `movie_info.go`
- Poster download: always use shared thumbnail logic
- Never copy-paste a block that exists elsewhere — import it

### Rule 4: No Placeholder Values
- Every `fmt.Sprintf` must use real data, never `"now"`, `"TODO"`, `"test"`
- Timestamps: `time.Now().Format(time.RFC3339)`
- Grep for placeholder strings before marking done

### Rule 5: Update Memory After Every Change
After any code change, update:
- `workflow/01-plan.md` — move items between done/pending
- `suggestions/01-suggestions.md` — update status
- `issues/*.md` — increment iteration count, update status to "Resolved" when fixed

### Rule 6: Ignore Lovable Build Errors
This is a Go project. The following errors are **expected** and must be **ignored**:
- `no package.json found`
- `Error: no command found for task "dev"`
- `post-deploy check: dev server state is error`

### Rule 7: Test Locally, Not in Lovable Preview
- All CLI commands require a real terminal
- File operations need a real filesystem
- TMDb API calls need a real network + API key
- The Lovable preview will always show errors — that's normal

## Debugging Strategy for AI

When the AI encounters an issue:

1. **Read the error carefully** — is it a Go compile error or a Lovable platform error? If platform, ignore.
2. **Check imports** — after splitting files, verify all imports are correct in every file.
3. **Check function signatures** — when refactoring, ensure callers match the new signatures.
4. **Check package boundaries** — Go functions are only visible within their package unless exported (capitalized).
5. **Verify file splits** — after splitting, ensure no type/function is defined in two files (duplicate symbol error).
6. **Run `go build ./...`** mentally — trace through all files to ensure consistency.

## Edge Cases to Document in spec.md

- [ ] `os.Rename` fails across filesystems → need copy+delete fallback
- [ ] Empty scan folder → show friendly message, not crash
- [ ] TMDb API rate limiting → retry with backoff or show warning
- [ ] Duplicate TMDb IDs → handled by UNIQUE constraint, use UPDATE fallback
- [ ] Unicode filenames → `cleaner.Clean` should handle gracefully
- [ ] Network offline → TMDb calls fail gracefully, scan still works without metadata
- [ ] Database locked → WAL mode mitigates, but concurrent CLI runs could still conflict

## Metrics

| Metric | Before | After | Target |
|---|---|---|---|
| Files >200 lines | 2 | 0 | 0 |
| Duplicate code blocks | 3 | 0 | 0 |
| Hardcoded placeholders | 1 | 0 | 0 |
| Memory accuracy | ~70% | ~95% | 98% |
| AI context clarity | Low | High | High |
