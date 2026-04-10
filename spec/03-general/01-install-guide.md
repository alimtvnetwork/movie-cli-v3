# 01 — Install Guide

> How to set up the **movie** CLI development environment from scratch.

## Prerequisites

| Requirement | Minimum Version | Check Command |
|-------------|-----------------|---------------|
| **Go** | 1.22+ | `go version` |
| **Git** | 2.x | `git --version` |
| **PowerShell** | 5.1+ (Windows) or 7+ (cross-platform) | `$PSVersionTable.PSVersion` |

### Install Go

- **Windows**: Download from [go.dev/dl](https://go.dev/dl/) or `winget install GoLang.Go`
- **macOS**: `brew install go`
- **Linux**: `sudo apt install golang` or download from [go.dev/dl](https://go.dev/dl/)

### Install PowerShell (macOS / Linux only)

```bash
# macOS
brew install --cask powershell

# Ubuntu/Debian
sudo apt-get install -y powershell

# Then launch with:
pwsh
```

## Installation

### Step 1 — Clone the Repository

```powershell
git clone https://github.com/movie/movie-cli-v2.git
cd movie-cli-v2
```

### Step 2 — Set Execution Policy (Windows only)

PowerShell may block script execution by default. Run once as Administrator:

```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### Step 3 — Run the Build Pipeline

```powershell
.\run.ps1
```

This will:
1. Pull latest code from `main`
2. Run `go mod tidy`
3. Build the `movie` binary into `./bin/`
4. Deploy to the configured deploy path (default: `E:\bin-run` on Windows, `/usr/local/bin` on Unix)

### Step 4 — Verify Installation

```powershell
movie version
```

Expected output: version string with commit hash and build date.

If `movie` is not found, ensure the deploy path is in your `PATH`:

```powershell
# Windows — add to PATH for current session
$env:PATH += ";E:\bin-run"

# macOS / Linux — /usr/local/bin is usually already in PATH
```

## Quick Install (One-Liner)

### Windows (PowerShell)

```powershell
git clone https://github.com/movie/movie-cli-v2.git; cd movie-cli-v2; .\run.ps1
```

### macOS / Linux

```bash
git clone https://github.com/movie/movie-cli-v2.git && cd movie-cli-v2 && pwsh run.ps1
```

## Configuration

The build pipeline reads `powershell.json` from the repo root. See [03-config-reference.md](03-config-reference.md) for all fields.

Default `powershell.json`:

```json
{
  "deployPath": "E:\\bin-run",
  "buildOutput": "./bin",
  "binaryName": "movie.exe",
  "copyData": false
}
```

## Troubleshooting

| Problem | Solution |
|---------|----------|
| `go: command not found` | Install Go and ensure it's in PATH |
| `pwsh: command not found` | Install PowerShell 7+: `brew install --cask powershell` |
| Script blocked by execution policy | Run `Set-ExecutionPolicy RemoteSigned -Scope CurrentUser` |
| `movie: command not found` after build | Add deploy path to PATH (see Step 4) |
| Git pull fails with local changes | Use `-ForcePull` or follow the interactive menu |
| Build fails with missing modules | Delete `go.sum` and re-run `go mod tidy` |
| Permission denied on `/usr/local/bin` | Run with `sudo pwsh run.ps1` or set `-DeployPath ~/bin` |

## What Gets Installed

After a successful run, you should have:

```
<deployPath>/
  └── movie.exe          # (or 'movie' on Unix)

<repoRoot>/
  └── bin/
      └── movie.exe      # build artifact (intermediate)
```
