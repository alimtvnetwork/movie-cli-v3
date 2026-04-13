# 01 — Self-Update Overview

## Purpose

Explain why CLI self-update is non-trivial, what platform-specific constraints exist, and how the overall architecture addresses them.

> **Reference**: Adapted from gitmap-v2 ([01-self-update-overview.md](https://github.com/alimtvnetwork/gitmap-v2/blob/main/spec/generic-update/01-self-update-overview.md))

---

## The Problem

When a user runs `mahin self-update`, the tool must:

1. Fetch the latest source or download a pre-built binary.
2. Build (if from source) or extract the new binary.
3. Replace the currently running binary with the new one.
4. Verify the update succeeded.
5. Clean up temporary artifacts.

**Step 3 is the hard part.** On Windows, a running `.exe` file is locked by the OS — it cannot be overwritten or deleted while the process is alive. On Linux/macOS, the file can be replaced in-place (the OS uses inode references, not file paths, for running processes).

---

## Platform Behavior

| Operation | Windows | Linux / macOS |
|-----------|---------|---------------|
| Overwrite running binary | ❌ Blocked (file lock) | ✅ Works |
| Rename running binary | ✅ Allowed | ✅ Works |
| Delete running binary | ❌ Blocked | ✅ Works |
| Replace after rename | ✅ Works | ✅ Works |

**Key insight**: Windows allows **renaming** a running executable but not **overwriting** or **deleting** it. This is the foundation of the rename-first deploy strategy.

---

## Two Update Strategies

### Strategy 1: Source-Based Update (Build from Repo) — Current

Used when the binary was installed from a source repository:

```
1. Resolve the source repo location
2. If the repo is missing, clone a fresh copy next to the binary
3. If an existing repo is found, verify it is clean
4. Pull latest code (git pull --ff-only) for existing repos only
5. Resolve dependencies (go mod tidy)
6. Build new binary (go build with ldflags)
7. Deploy to the installed location (rename-first)
8. Verify version
9. Clean up
```

**Current mahin implementation** (`updater/updater.go`):
- Checks git is installed and in PATH
- Resolves the repo from the binary directory, current working directory, or sibling clone path
- Clones a fresh repo next to the binary when no local repo exists
- Treats a fresh clone as **bootstrap success**, not as an "already up to date" result
- Verifies no local changes only for the resolved local repo
- Runs `git pull --ff-only` (safe, no merge conflicts) for existing repos
- Reports either **bootstrap success**, **already up to date**, or **old→new commit SHAs**
- User runs `run.ps1` or `build.ps1` for rebuild + deploy

### Strategy 2: Binary-Based Update (Download from Releases) — Future

Used when the user installed from a GitHub Release:

```
1. Query GitHub API for latest release
2. Download the correct archive for current OS/arch
3. Download checksums.txt
4. Verify SHA256
5. Extract binary
6. Deploy using rename-first
7. Verify version
8. Clean up
```

This strategy requires no Go toolchain on the end-user's machine.

---

## Decision: Which Strategy to Use

| Signal | Strategy |
|--------|----------|
| `.git` directory exists in binary's parent path | Source-based |
| Sibling source clone exists next to binary | Source-based |
| No local source repo found but git is available | Source-based with bootstrap clone |
| User passes `--from-source` flag | Source-based |
| User passes `--from-release` flag | Binary-based |

---

## Acceptance Criteria

- GIVEN a clean existing git repo WHEN `mahin self-update` runs THEN the latest code is pulled
- GIVEN no local repo exists WHEN `mahin self-update` runs THEN it clones a fresh repo next to the binary and reports bootstrap success
- GIVEN local changes WHEN `mahin self-update` runs THEN it refuses with an error message
- GIVEN git is not installed WHEN `mahin self-update` runs THEN a clear error is shown
- GIVEN the repo is already at latest WHEN `mahin self-update` runs THEN it reports "already up to date" only for an existing repo with no new commits

---

*Self-update overview — updated: 2026-04-13*
