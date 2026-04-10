# Reliability & Failure-Chance Report — Movie CLI

**Updated:** 05-Apr-2026  
**Version:** 1.0.0  
**Status:** Active  
**Purpose:** Assess spec readiness for AI handoff and identify failure risks

---

## 1. Success Probability Estimates

### By Module Complexity Tier

| Tier | Modules | Success % | Assumptions |
|------|---------|-----------|-------------|
| **Simple** (single command, no external deps) | `hello`, `version`, `self-update`, `config`, `play` | **92–95%** | Cobra patterns are well-known; `play` is trivial OS exec; `config` is basic CRUD. Assumes AI reads project overview first. |
| **Medium** (DB + API + interactive I/O) | `ls`, `stats`, `info`, `search`, `scan` | **75–85%** | `ls` pagination is straightforward. `scan` + `search` depend on shared TMDb helpers being understood. `info` has a 3-level lookup priority chain that must be followed exactly. |
| **Complex Agentic** (multi-step workflows, file I/O) | `move`, `rename`, `undo`, `suggest` | **60–70%** | `move` has a 10-step flow with JSON logging, DB tracking, and interactive prompts. `suggest` mixes genre analysis, recommendations, and trending with dedup. `undo` must correctly reverse state. Cross-drive `os.Rename` failure is unhandled. |
| **End-to-End** (full system integration) | Full CLI with all 11 movie commands working together | **50–60%** | Shared helpers (`fetchMovieDetails`, `resolveMediaByQuery`) must be correctly imported. File splits across `db/` (5 files) and `cmd/` (15 files) increase context load. Schema migrations must match code expectations. |

### By AI Type

| AI Type | Estimated Success | Notes |
|---------|------------------|-------|
| Expert Go developer (human) | 95% | Full spec + HANDOFF.md is sufficient |
| Go-capable AI (Cursor, Copilot, Aider) | 70–80% | Must read ai-handoff.md first; may struggle with split DB files |
| Web-focused AI (Lovable, v0) | 20–30% | Fundamental platform mismatch — no `go build`, no terminal |
| AI without filesystem access | ~0% | File operations are core to 6/11 commands |

---

## 2. Failure Map

### 2.1 High-Risk Failures

| # | Module/Workflow | Why It Fails | How It Manifests |
|---|----------------|--------------|------------------|
| F1 | **`movie move`** — cross-drive | `os.Rename` silently fails across filesystems; no fallback exists | File appears "moved" in DB but stays in original location |
| F2 | **Shared helper imports** | AI doesn't realize `fetchMovieDetails` lives in `movie_info.go` and is unexported (lowercase) but accessible within `cmd` package | Compile error: `undefined: fetchMovieDetails` — AI creates duplicate function |
| F3 | **DB file split context** | 5 files in `db/` package; AI edits `media.go` but needs type from `db.go` | Import confusion, duplicate `DB` struct definitions, compile errors |
| F4 | **`movie undo`** — no confirmation | User accidentally undoes; no way to redo | Data loss if undo moves file to a cleaned-up directory |
| F5 | **TMDb API key flow** | Two sources (config DB + env var); scan continues without key but search exits | AI may unify behavior incorrectly, breaking graceful degradation |

### 2.2 Medium-Risk Failures

| # | Module/Workflow | Why It Fails | How It Manifests |
|---|----------------|--------------|------------------|
| F6 | **`movie scan`** — directory vs file handling | Directories containing video files use directory name for cleaning, not file name | Wrong metadata fetched; title mismatch |
| F7 | **`movie info`** — 3-level lookup | Numeric ID → local title search → TMDb fallback with auto-persist + dedup check | AI may skip dedup check (`GetMediaByTmdbID`), causing UNIQUE constraint violations |
| F8 | **`movie suggest`** — genre analysis | `TopGenres()` returns CSV-split genre counts; suggest picks random library items for recommendations | AI may not handle empty library edge case (should fall back to trending) |
| F9 | **Cleaner regex** — edge cases | Unicode filenames, multi-year strings, nested brackets | Silent incorrect cleaning; wrong movie matched on TMDb |
| F10 | **`movie ls`** — only scan-indexed items | Spec says ls shows only items from `scan`, not `search`/`info` | AI may show all DB records, mixing cataloged-only items with file-backed items |

### 2.3 Low-Risk Failures

| # | Module/Workflow | Why It Fails | How It Manifests |
|---|----------------|--------------|------------------|
| F11 | **Version injection** | `-ldflags` must be exact; wrong module path = no injection | Version shows `v0.0.1-dev` in production builds |
| F12 | **Poster download** | TMDb CDN path construction | Missing poster image; thumbnail_path set but file doesn't exist |
| F13 | **WAL mode** | Not set → concurrent access issues | DB locked errors on rapid successive commands |

---

## 3. Corrective Actions (Prioritized)

| Priority | Action | Where | Expected Reliability Gain |
|----------|--------|-------|--------------------------|
| **P0** | Add cross-drive move fallback (copy+delete) | `cmd/movie_move.go` | +5% on `move` command; eliminates F1 |
| **P0** | Add confirmation prompt to `movie undo` | `cmd/movie_undo.go` | Eliminates F4; prevents accidental data loss |
| **P1** | Add acceptance criteria (GIVEN/WHEN/THEN) to spec.md | `spec.md` §4 (each command) | +10% overall; AI can validate its own output |
| **P1** | Document shared helper locations explicitly | `ai-handoff.md` §3 | +5% overall; eliminates F2 |
| **P1** | Add `movie ls` filter clarification to code comments | `cmd/movie_ls.go` | Eliminates F10 |
| **P2** | Add error handling spec for TMDb rate limiting | `spec/08-app/01-project-spec.md` | +3% on scan/search; handles F5 edge |
| **P2** | Add cleaner unit test spec with edge cases | New: `spec/08-app/02-cleaner-test-spec.md` | Eliminates F9 ambiguity |
| **P3** | Add `movie tag` command spec | `spec.md` §4, new section | Table exists; commands missing; reduces tech debt |
| **P3** | Write JSON metadata per movie/TV on scan | `cmd/movie_scan.go` | Fulfills storage structure promise |

---

## 4. Readiness Decision

### ❌ NOT READY for full handoff — needs P0 and P1 fixes first

**What must be fixed before implementation handoff:**

1. **P0: Cross-drive move** — Without this, `movie move` silently fails on common setups (external drives, network mounts)
2. **P0: Undo confirmation** — Destructive operation without safety net
3. **P1: Acceptance criteria** — Without testable criteria, AI cannot self-validate
4. **P1: Shared helper documentation** — Without explicit "this function lives HERE", AI will duplicate code

**What CAN proceed now:**
- Simple commands (`hello`, `version`, `config`, `play`, `stats`) are fully specified
- DB schema is clear and well-documented
- TMDb API integration is well-documented
- Cleaner regex is well-documented (though edge case tests would help)

**After P0+P1 fixes, estimated success rate rises from 50–60% → 75–85% for end-to-end.**

---

## 5. Spec Gaps Summary

| Gap | Severity | Location |
|-----|----------|----------|
| No acceptance criteria for any command | High | spec.md §4 |
| Cross-drive move not handled | High | spec.md §4.10, code |
| No undo confirmation | Medium | spec.md §4.12, code |
| `tags` table unused | Low | spec.md §3, code |
| `DiscoverByGenre` unused | Low | spec.md §6, tmdb/client.go |
| JSON metadata files not written | Low | spec.md §3 |
| No error handling spec | Medium | spec/08-app/ |
| `movie ls` filter rule unclear in code | Medium | cmd/movie_ls.go |

---

*Report generated: 05-Apr-2026*
