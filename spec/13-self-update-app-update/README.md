# Self-Update & App Update Specification

## Purpose

This folder defines the **self-update architecture** for the mahin CLI tool. It covers the full lifecycle: detecting where the binary is installed, building a new version, deploying it without file-lock errors, and cleaning up afterward.

Any AI or engineer reading these documents should be able to implement a complete self-update system from scratch without ambiguity.

> **Reference implementation**: gitmap-v2 ([generic-update spec](https://github.com/alimtvnetwork/gitmap-v2/tree/main/spec/generic-update), [generic-release spec](https://github.com/alimtvnetwork/gitmap-v2/tree/main/spec/generic-release))

---

## Documents

| File | Topic |
|------|-------|
| [01-self-update-overview.md](./01-self-update-overview.md) | The problem, approach, and platform differences |
| [02-deploy-path-resolution.md](./02-deploy-path-resolution.md) | 3-tier strategy for finding the installed binary |
| [03-rename-first-deploy.md](./03-rename-first-deploy.md) | Rename-first strategy to bypass file locks |
| [04-build-scripts.md](./04-build-scripts.md) | `run.ps1` and `build.ps1` patterns for build + deploy |
| [05-release-distribution.md](./05-release-distribution.md) | Cross-compilation, install scripts, checksums |
| [06-cleanup.md](./06-cleanup.md) | Post-update artifact removal |

---

## Core Principle

A running binary **cannot overwrite itself** on Windows. The entire update architecture exists to work around this constraint while maintaining a seamless user experience.

## Current Implementation

The mahin CLI uses **Strategy 1 (Source-Based Update)** via `git pull --ff-only`:
- `mahin self-update` → pulls latest source → user rebuilds via `run.ps1` or `build.ps1`
- See `updater/updater.go` and `cmd/update.go`

## Future Enhancement

**Strategy 2 (Binary-Based Update)** via GitHub Releases:
- Download pre-built binary from GitHub Releases
- Verify SHA256 checksum
- Replace installed binary using rename-first deploy
- No Go toolchain required on end-user machines

## Placeholders

| Placeholder | Meaning | Mahin Value |
|-------------|---------|-------------|
| `<binary>` | CLI binary name | `movie` (or `movie.exe`) |
| `<deploy-dir>` | Install directory | `$env:LOCALAPPDATA\movie` (Win) / `~/.local/bin` (Unix) |
| `<repo-root>` | Source repository root | Repository root containing `go.mod` |
| `<module>` | Go module path | `github.com/mahin/mahin-cli-v2` |

---

*Self-update specs — updated: 2026-04-10*
