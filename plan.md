# Movie CLI — Future Work Roadmap

> **Last Updated**: 05-Apr-2026  
> **Purpose**: Handoff document for any AI or developer continuing work on this project

---

## Phase 1: Safety & Reliability (P0)

### 1.1 Cross-Drive Move Fallback ✅ DONE (2026-04-05)
- **Objective**: Make `movie move` work across filesystems (USB drives, network mounts)
- **Dependencies**: None
- **Expected outputs**: Modified `cmd/movie_move_helpers.go`
- **Implementation**: `MoveFile()` attempts `os.Rename`, detects `EXDEV`, falls back to `io.Copy` + `os.Remove`
- **Acceptance criteria**:
  - GIVEN a video file on drive C: WHEN moved to drive D: THEN file appears at destination and is removed from source ✅
  - GIVEN a same-drive move WHEN `os.Rename` succeeds THEN no copy fallback is used ✅

### 1.2 Undo Confirmation Prompt
- **Objective**: Prevent accidental undo of move/rename operations
- **Dependencies**: None
- **Expected outputs**: Modified `cmd/movie_undo.go`
- **Implementation**: Show from/to paths, ask `Are you sure you want to undo? [y/N]:`
- **Acceptance criteria**:
  - GIVEN a pending undo WHEN user types 'n' THEN no file is moved
  - GIVEN a pending undo WHEN user types 'y' THEN file is reverted and DB updated

---

## Phase 2: Spec Completeness (P1)

### 2.1 Acceptance Criteria for All Commands
- **Objective**: Add GIVEN/WHEN/THEN test cases to spec.md
- **Dependencies**: None
- **Expected outputs**: Updated `spec.md` §4, each command gets 2–4 acceptance criteria
- **Acceptance criteria**: Every command in §4 has at least 2 testable scenarios

### 2.2 Shared Helper Documentation
- **Objective**: Add comments marking shared helpers so AI doesn't duplicate them
- **Dependencies**: None
- **Expected outputs**: Updated comments in `cmd/movie_info.go`, `cmd/movie_resolve.go`
- **Acceptance criteria**: Each shared function has a comment like `// SHARED: used by scan, search, info`

### 2.3 Movie LS Filter Clarification
- **Objective**: Clarify that `movie ls` shows only file-backed (scanned) media
- **Dependencies**: None
- **Expected outputs**: Updated `cmd/movie_ls.go` comments + spec.md §4.6
- **Acceptance criteria**: Code comment + spec both state "only items with non-empty `original_file_path`"

---

## Phase 3: New Features (P2)

### 3.1 Movie Tag Command ✅ DONE (2026-04-06)
- **Objective**: Expose the existing `tags` table via CLI commands
- **Dependencies**: `tags` table (already exists in schema)
- **Expected outputs**: New `cmd/movie_tag.go`, updated DB methods in `db/media.go`
- **Subcommands**: `movie tag add <id> <tag>`, `movie tag remove <id> <tag>`, `movie tag list [id]`
- **Acceptance criteria**:
  - GIVEN a media ID WHEN `tag add 1 favorite` THEN tag is inserted
  - GIVEN a duplicate tag WHEN `tag add 1 favorite` again THEN error "tag already exists"
  - GIVEN tags exist WHEN `tag list 1` THEN all tags for media 1 are shown

### 3.2 File Size Stats
- **Objective**: Add file size information to `movie stats`
- **Dependencies**: `file_size` column in `media` table
- **Expected outputs**: Updated `cmd/movie_stats.go`
- **Acceptance criteria**: Stats output shows total library size, average file size, largest file

### 3.3 Error Handling Spec
- **Objective**: Document error handling for TMDb rate limits, DB locks, offline mode
- **Dependencies**: None
- **Expected outputs**: New `spec/01-app/02-error-handling-spec.md`
- **Acceptance criteria**: Spec covers TMDb 429 response, SQLite BUSY, network offline

### 3.4 README Update
- **Objective**: Document all movie management features in README.md
- **Dependencies**: None
- **Expected outputs**: Updated `README.md`
- **Acceptance criteria**: README lists all 11 movie commands with usage examples

---

## Phase 4: Enhancements (P3)

### 4.1 Batch Move
- **Objective**: Add `--all` flag to `movie move` to move all video files at once
- **Dependencies**: Phase 1.1 (cross-drive support)
- **Expected outputs**: Updated `cmd/movie_move.go`

### 4.2 JSON Metadata Files
- **Objective**: Write per-movie JSON metadata during scan
- **Dependencies**: None
- **Expected outputs**: Updated `cmd/movie_scan.go`, JSON files in `./data/json/movie/` and `json/tv/`

### 4.3 DiscoverByGenre Integration
- **Objective**: Use `DiscoverByGenre` TMDb method in `movie suggest`
- **Dependencies**: None
- **Expected outputs**: Updated `cmd/movie_suggest.go`

---

## Next Task Selection

**Ready to implement now (no dependencies):**
1. Cross-drive move fallback (Phase 1.1)
2. Undo confirmation prompt (Phase 1.2)
3. Movie tag command (Phase 3.1)
4. Acceptance criteria (Phase 2.1)
5. Shared helper documentation (Phase 2.2)
6. README update (Phase 3.4)

*Pick a task and I'll implement it.*
