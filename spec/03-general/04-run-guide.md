# 04 — run.ps1 Usage Guide

> Complete guide for using the `run.ps1` automation script.

## Purpose

`run.ps1` is the single command to pull, build, deploy, and optionally run the **movie** CLI. It replaces manual `go build` workflows with a repeatable pipeline that handles git safety, dependency resolution, cross-platform compilation, safe deployment with rollback, and argument forwarding.

## Syntax

```powershell
.\run.ps1 [-NoPull] [-NoDeploy] [-ForcePull] [-DeployPath <path>] [-R] [-t] [-- <args>]
```

## Flags

| Flag | Description |
|------|-------------|
| `-NoPull` | Skip git pull entirely |
| `-NoDeploy` | Build but don't deploy to the target directory |
| `-ForcePull` | Discard all local changes and untracked files before pulling (no prompt) |
| `-DeployPath` | Override the deploy directory from `powershell.json` |
| `-R` | After build+deploy, run the binary with forwarded arguments |
| `-t` | Run all Go unit tests and exit (skips build/deploy) |

## Default Execution Flow

```
.\run.ps1
```

Without any flags, the script runs the full pipeline:

1. **Banner** — displays the movie builder header
2. **Config** — loads `powershell.json` (falls back to defaults)
3. **Git Pull** — ensures `main` branch, pulls latest changes
4. **Dependencies** — verifies Go is installed, runs `go mod tidy`
5. **Build** — compiles with ldflags (version, commit, date), verifies output
6. **Deploy** — safe copy with rename-first strategy, rollback on failure, PATH check

## Argument Forwarding with `-R`

The `-R` switch activates run mode. All remaining arguments are forwarded to the binary:

```powershell
# Explicit arguments
.\run.ps1 -R movie scan D:\movies

# No arguments → defaults to: movie scan <parentDir>
.\run.ps1 -R
```

Relative paths in arguments are automatically resolved to absolute paths before execution.

## Test Mode with `-t`

```powershell
.\run.ps1 -t
```

- Runs `go test ./... -v`
- Writes a timestamped report to `data/unit-test-reports/test-report-YYYYMMDD-HHmmss.txt`
- Displays pass/fail/skip summary
- Exits immediately (does not build or deploy)

## Configuration

Config is loaded from `powershell.json`. See [03-config-reference.md](03-config-reference.md).

When no config file exists, these defaults apply:

| Field | Default |
|-------|---------|
| `deployPath` | `E:\bin-run` |
| `buildOutput` | `./bin` |
| `binaryName` | `movie.exe` |
| `copyData` | `false` |

## Git Safety

### Branch Protection

The script checks the current branch before pulling. If not on `main`, it automatically switches.

### Conflict Resolution

If `git pull` fails due to local changes, an interactive menu appears:

| Option | Action |
|--------|--------|
| **S** — Stash | Saves changes with `git stash`, then pulls. Restore later with `git stash pop`. |
| **D** — Discard | Runs `git checkout -- .` to discard tracked changes, then pulls. |
| **C** — Clean | Discards changes AND removes untracked files (`git clean -fd`), then pulls. |
| **Q** — Quit | Aborts without making any changes. |

Use `-ForcePull` to skip this menu and automatically run the **Clean** strategy (useful for CI).

## Build Details

- **Cross-platform**: auto-detects Windows, macOS (arm64/amd64), Linux
- **ldflags injection**: embeds `Version` (from git tag), `Commit` (short hash), `BuildDate` (RFC3339)
- **Post-build verification**: checks binary exists, reports file size, runs `version` command

## Deploy Safety

The deploy step uses a **rename-first** strategy:

1. Rename existing binary to `.bak`
2. Copy new binary to deploy path
3. If copy fails → rollback by renaming `.bak` back
4. If copy succeeds → delete `.bak`
5. On Unix: `chmod +x` the deployed binary
6. Check if deploy path is in `PATH` and warn if not

## Examples

```powershell
# Full pipeline: pull → build → deploy
.\run.ps1

# Skip pull, build and deploy only
.\run.ps1 -NoPull

# Build only, no deploy
.\run.ps1 -NoPull -NoDeploy

# Force-pull (CI mode), build, deploy
.\run.ps1 -ForcePull

# Build, deploy, then scan a folder
.\run.ps1 -R movie scan D:\movies

# Build, deploy, then show help
.\run.ps1 -R help

# Just run tests
.\run.ps1 -t

# Deploy to custom path
.\run.ps1 -DeployPath C:\tools
```

## Console Output

The script uses semantic color-coded output:

| Prefix | Color | Meaning |
|--------|-------|---------|
| `[N/4]` | Magenta | Phase step header |
| `OK` | Green | Success |
| `->` | Cyan | Informational |
| `!!` | Yellow | Warning |
| `XX` | Red | Error/failure |

## Troubleshooting

| Symptom | Cause | Fix |
|---------|-------|-----|
| Script exits at git pull | Local changes conflict | Use `-ForcePull` or choose from the interactive menu |
| `go: command not found` | Go not installed | Install from [go.dev/dl](https://go.dev/dl/) |
| Build fails | Syntax error or missing dependency | Check `go build` output, run `go mod tidy` manually |
| Deploy permission denied | No write access to deploy path | Run as admin or use `-DeployPath ~/bin` |
| Binary not found after deploy | Deploy path not in PATH | Add deploy path to system PATH |
| `-R` runs wrong command | Arguments not forwarded | Place all args after `-R`: `.\run.ps1 -R movie scan path` |

## Error Handling

All external commands (git, go) are wrapped with `$ErrorActionPreference = "Continue"` to capture output and exit codes without PowerShell terminating the script. Failures produce red `XX` messages with the captured output.
