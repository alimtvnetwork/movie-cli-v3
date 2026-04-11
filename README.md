<div align="center">

<img src="assets/icon.png" alt="Movie CLI icon" width="80" height="80">

# 🎬 Movie CLI

**Personal movie & TV show library manager — from the terminal**

[![CI](https://github.com/alimtvnetwork/movie-cli-v3/actions/workflows/ci.yml/badge.svg)](https://github.com/alimtvnetwork/movie-cli-v3/actions/workflows/ci.yml)
[![GitHub Release](https://img.shields.io/github/v/release/alimtvnetwork/movie-cli-v3?style=flat-square&label=version)](https://github.com/alimtvnetwork/movie-cli-v3/releases)
[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey?style=flat-square)](https://github.com/alimtvnetwork/movie-cli-v3)
[![SQLite](https://img.shields.io/badge/SQLite-WAL-003B57?style=flat-square&logo=sqlite&logoColor=white)](https://www.sqlite.org/)
[![License](https://img.shields.io/badge/license-Private-red?style=flat-square)](https://github.com/alimtvnetwork/movie-cli-v3)

_Scan folders, clean filenames, fetch TMDb metadata, organize files, and track your collection._

</div>

---

<div align="center">

## Quick Start

</div>

### Install

```
# Windows (PowerShell)
irm https://github.com/alimtvnetwork/movie-cli-v3/releases/latest/download/install.ps1 | iex
```

```
# Linux / macOS
curl -fsSL https://github.com/alimtvnetwork/movie-cli-v3/releases/latest/download/install.sh | bash
```

### Set up & scan

```bash
movie config set tmdb_api_key YOUR_KEY
movie scan ~/Downloads
movie ls
```

### Search & discover

```bash
movie search "Inception"
movie suggest 5
```

Every command supports `--help` or `-h` for detailed usage.

---

<div align="center">

## 🎥 Demo

</div>

### 📂 Scanning a Folder

<!-- Replace with actual GIF: docs/screenshots/demo-scan.gif -->
<!-- Record with: vhs docs/screenshots/scan.tape  or  asciinema rec -->

```
$ movie scan ~/Downloads

🔍 Scanning: /home/user/Downloads
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  Found 12 video files

  [1/12] Scream.2022.1080p.WEBRip.x264-RARBG.mkv
         → Title: Scream (2022)
         → TMDb: ★ 6.8 | Horror, Mystery, Thriller
         → Poster saved: thumbnails/scream-2022/scream-2022.jpg
         ✅ Saved to database

  [2/12] The.Batman.2022.2160p.BluRay.x265.mkv
         → Title: The Batman (2022)
         → TMDb: ★ 7.7 | Crime, Mystery, Thriller
         → Poster saved: thumbnails/the-batman-2022/the-batman-2022.jpg
         ✅ Saved to database

  ...

  ✅ Done — 12 items scanned, 11 new, 1 updated
```

<p align="center">
  <img src="docs/screenshots/demo-scan.gif" alt="movie scan demo" width="700">
  <br><em>↑ Replace with actual recording</em>
</p>

---

### 📋 Browsing Your Library

<!-- Replace with actual GIF: docs/screenshots/demo-ls.gif -->

```
$ movie ls

🎬 Library — Page 1 of 3 (20 per page)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  #   Title                          Year   Type    Rating
  ─── ────────────────────────────── ────── ─────── ──────
  1   Scream                         2022   🎬      ★ 6.8
  2   The Batman                     2022   🎬      ★ 7.7
  3   Everything Everywhere All...   2022   🎬      ★ 7.8
  4   Breaking Bad                   2008   📺      ★ 8.9
  5   Severance                      2022   📺      ★ 8.4
  ...

  [N]ext  [P]rev  [1-9] Detail  [Q]uit
```

<p align="center">
  <img src="docs/screenshots/demo-ls.gif" alt="movie ls demo" width="700">
  <br><em>↑ Replace with actual recording</em>
</p>

---

### 🎯 Getting Suggestions

<!-- Replace with actual GIF: docs/screenshots/demo-suggest.gif -->

```
$ movie suggest 5

🎯 Movie Suggest
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  Select category:
  1. 🎬 Movie
  2. 📺 TV
  3. 🎲 Random

  Choose: 1

  📽️  Recommendations based on your library:

  #   Title                          Year   Rating   Genre
  ─── ────────────────────────────── ────── ──────── ────────────────
  1   Nope                           2022   ★ 6.8    Horror, Sci-Fi
  2   X                              2022   ★ 6.6    Horror, Mystery
  3   Pearl                          2022   ★ 7.0    Drama, Horror
  4   Bodies Bodies Bodies            2022   ★ 6.5    Comedy, Horror
  5   Barbarian                      2022   ★ 7.0    Horror, Thriller

  🔥 Trending This Week:
  1   Oppenheimer                    2023   ★ 8.1    Drama, History
  2   Killers of the Flower Moon     2023   ★ 7.5    Crime, Drama
  3   Poor Things                    2023   ★ 7.9    Comedy, Drama
```

<p align="center">
  <img src="docs/screenshots/demo-suggest.gif" alt="movie suggest demo" width="700">
  <br><em>↑ Replace with actual recording</em>
</p>

> **📹 Recording your own demos:**
> Use [VHS](https://github.com/charmbracelet/vhs) or [asciinema](https://asciinema.org/) to record terminal sessions as GIFs.
> ```bash
> # VHS (recommended — deterministic, scriptable)
> vhs docs/screenshots/scan.tape
>
> # asciinema + agg (manual recording)
> asciinema rec demo.cast
> agg demo.cast docs/screenshots/demo-scan.gif
> ```

---

<div align="center">

## Installation

</div>

### One-Liner Install (recommended)

**Windows (PowerShell):**

```powershell
irm https://github.com/alimtvnetwork/movie-cli-v3/releases/latest/download/install.ps1 | iex
```

**Linux / macOS (Bash):**

```bash
curl -fsSL https://github.com/alimtvnetwork/movie-cli-v3/releases/latest/download/install.sh | bash
```

To pin a specific release, replace `latest/download` with `download/vX.Y.Z`:

```powershell
irm https://github.com/alimtvnetwork/movie-cli-v3/releases/download/v1.0.0/install.ps1 | iex
```

```bash
curl -fsSL https://github.com/alimtvnetwork/movie-cli-v3/releases/download/v1.0.0/install.sh | bash
```

> These one-liners use the installer assets attached to each GitHub Release. The repo-root `install.ps1` and `install.sh` are source bootstrap scripts for local builds.

### Installer Options

**Windows (PowerShell):**

| Flag | Description | Example |
|---|---|---|
| `-InstallDir` | Custom install directory | `-InstallDir C:\tools\movie` |
| `-Arch` | Force architecture (`amd64`, `arm64`) | `-Arch arm64` |
| `-NoPath` | Skip adding to user PATH | `-NoPath` |

**Linux / macOS (Bash):**

| Flag | Description | Example |
|---|---|---|
| `--dir` | Custom install directory | `--dir ~/bin` |
| `--arch` | Force architecture (`amd64`, `arm64`) | `--arch arm64` |
| `--no-path` | Skip adding to PATH | `--no-path` |

### Clone & Build (Development)

**Prerequisites:**

| Requirement | Minimum | Check |
|---|---|---|
| **Go** | 1.22+ | `go version` |
| **Git** | 2.x | `git --version` |
| **PowerShell** | 5.1+ (Win) / 7+ (Unix) | `$PSVersionTable.PSVersion` |

```bash
git clone https://github.com/alimtvnetwork/movie-cli-v3.git
cd movie-cli-v3
pwsh run.ps1
```

**Using the bootstrap installer:**

```powershell
pwsh install.ps1                      # Fresh install (clone + build + deploy)
pwsh install.ps1 -DeployPath ~/bin    # Custom deploy path
```

### Verify

```bash
movie version
# v1.0.0 (commit: abc1234, built: 2026-04-09)
#   Go:   go1.22.0
#   OS:   linux/amd64
```

> **Tip:** If `movie` is not found, add the deploy path to your `PATH`.
> Default: `E:\bin-run` (Windows) or `/usr/local/bin` (Unix) for source builds.

---

<div align="center">

## What It Does

</div>

A portable CLI that manages your personal movie and TV show library entirely from the terminal. Every scan produces:

- **Database** — structured metadata in SQLite (WAL mode)
- **Thumbnails** — poster images downloaded from TMDb
- **JSON** — per-file metadata written to `./data/json/`
- **Clean filenames** — `Scream.2022.1080p.WEBRip.x264.mkv` → `Scream (2022).mkv`

All data lives in `./data/` at the project root.

---

<div align="center">

## Command Reference

</div>

### Scanning & Library

| Command | Description |
|---|---|
| `movie scan [folder]` | Scan folder → DB + TMDb metadata |
| `movie ls` | Paginated interactive library browser |
| `movie search <name>` | Live TMDb search → save to DB |
| `movie info <id\|title>` | Detail view (local DB → TMDb fallback) |

```bash
movie scan ~/Downloads
movie ls
movie search "Inception"
movie info 1
```

---

### File Management

| Command | Description |
|---|---|
| `movie move [directory]` | Browse, select, move with clean name |
| `movie move --all` | Batch move all files (auto-route by type) |
| `movie rename` | Batch rename to clean format |
| `movie undo` | Revert last move/rename |
| `movie play <id>` | Open with default video player |

```bash
movie move ~/Downloads        # interactive single-file
movie move --all ~/Downloads  # batch move all
movie rename                  # clean all filenames
movie undo                    # revert last operation
movie play 1                  # play with system player
```

---

### Discovery & Organization

| Command | Description |
|---|---|
| `movie suggest [N]` | Genre-based recommendations + trending |
| `movie tag add <id> <tag>` | Add a tag to a media item |
| `movie tag remove <id> <tag>` | Remove a tag |
| `movie tag list [id]` | List tags (per item or all) |
| `movie stats` | Counts, storage, genre chart, avg ratings |

```bash
movie suggest 5
movie tag add 1 favorite
movie tag list
movie stats
```

---

### Configuration & System

| Command | Description |
|---|---|
| `movie config` | Show all configuration |
| `movie config set <key> <val>` | Set a config value |
| `movie export [-o path]` | Dump media table as JSON |
| `movie version` | Version, commit, build date, Go, OS/arch |
| `movie self-update` | Update via `git pull --ff-only` |
| `movie changelog [--latest]` | Show release notes |

```bash
movie config set movies_dir ~/Movies
movie config set tmdb_api_key YOUR_KEY
movie export -o ~/Desktop/library.json
movie self-update
```

**Config keys:**

| Key | Default | Purpose |
|---|---|---|
| `movies_dir` | `~/Movies` | Movie file destination |
| `tv_dir` | `~/TVShows` | TV show destination |
| `archive_dir` | `~/Archive` | Archive destination |
| `scan_dir` | `~/Downloads` | Default scan source |
| `tmdb_api_key` | *(none)* | TMDb API key |
| `page_size` | `20` | Items per page in `ls` |

---

<div align="center">

## Command Tree

</div>

```
movie
├── hello                         # Greeting with version
├── version                       # Version, commit, build date, Go, OS/arch
├── changelog [--latest]          # Show changelog (full or latest version)
├── self-update                   # git pull --ff-only
├── config [get|set] [key]        # View/set configuration
├── scan [folder]                 # Scan folder → DB + TMDb metadata
├── ls                            # Paginated library list (file-backed only)
├── search <name>                 # Live TMDb search → save to DB
├── info <id|title>               # Detail view (local DB → TMDb fallback)
├── suggest [N]                   # Recommendations + trending
├── move [directory]              # Browse, select, move with clean name
├── rename                        # Batch rename to clean format
├── undo                          # Revert last move/rename
├── play <id>                     # Open with default video player
├── stats                         # Counts, storage, genre chart, avg ratings
├── tag [add|remove|list]         # Manage user-defined tags
└── export [-o path]              # Dump media table as JSON
```

---

<div align="center">

## Build & Deploy

</div>

### Makefile Targets

| Target | Description |
|---|---|
| `make build` | Compile for current platform |
| `make build-all` | Cross-compile all 6 targets into `dist/` |
| `make build-windows` | Windows amd64 (with embedded icon) |
| `make build-windows-arm` | Windows arm64 (with embedded icon) |
| `make build-mac-arm` | macOS ARM64 |
| `make build-mac-intel` | macOS amd64 |
| `make build-linux` | Linux amd64 |
| `make build-linux-arm` | Linux arm64 |
| `make install` | Build + copy to `/usr/local/bin` |

### Build via run.ps1

```powershell
.\run.ps1                           # Full pipeline: pull, build, deploy
.\run.ps1 -NoPull                   # Skip git pull
.\run.ps1 -NoPull -NoDeploy        # Build only
.\run.ps1 -R movie scan D:\movies  # Build + run scan
.\run.ps1 -t                       # Run all unit tests
.\run.ps1 -ForcePull               # CI mode: discard changes + pull
```

| Flag | Description |
|---|---|
| `-NoPull` | Skip `git pull` |
| `-NoDeploy` | Skip deploy step |
| `-R` | Run movie after build (trailing args forwarded) |
| `-t` | Run all unit tests |
| `-ForcePull` | CI mode: discard changes + pull |

See [spec/03-general/04-run-guide.md](spec/03-general/04-run-guide.md) for the full usage guide.

---

<div align="center">

## Release Workflow

</div>

Releases are fully automated via GitHub Actions. Pushing to a `release/**` branch or a `v*` tag triggers:

1. **Cross-compilation** — 6 binaries (Windows/Linux/macOS × amd64/arm64)
2. **Packaging** — `.zip` (Windows) and `.tar.gz` (Unix)
3. **SHA256 checksums** — `checksums.txt` with all artifact hashes
4. **Install scripts** — version-pinned `install.ps1` and `install.sh`
5. **GitHub Release** — formatted page with changelog, checksums, and install instructions

### Creating a Release

```bash
# Option A: Push a release branch
git checkout -b release/v1.3.0
git push origin release/v1.3.0

# Option B: Tag directly
git tag v1.3.0
git push origin v1.3.0
```

Both trigger the same pipeline. Version is resolved from the ref name.

> **CI Pipeline:** Pushing a `release/*` branch or `v*` tag triggers GitHub Actions to cross-compile 6 targets, generate checksums, and create a GitHub release with changelog and install instructions.

See [spec/12-ci-cd-pipeline/02-release-pipeline.md](spec/12-ci-cd-pipeline/02-release-pipeline.md) for the full pipeline spec.

---

<div align="center">

## Project Structure

</div>

```
movie-cli-v3/
├── main.go                        # Entry point
├── cmd/                           # Cobra commands (one file per command)
│   ├── root.go                    # Root command, registers subcommands
│   ├── movie_config.go            # config get/set
│   ├── movie_scan.go              # scan folder
│   ├── movie_ls.go                # paginated list
│   ├── movie_search.go            # TMDb search
│   ├── movie_info.go              # detail view + shared fetch helpers
│   ├── movie_suggest.go           # recommendations
│   ├── movie_move.go              # interactive move
│   ├── movie_rename.go            # batch rename
│   ├── movie_undo.go              # undo last move/rename
│   ├── movie_stats.go             # library statistics
│   ├── movie_tag.go               # tag management
│   └── movie_resolve.go           # shared ID/title resolver
├── cleaner/cleaner.go             # Filename cleaning + slug generation
├── tmdb/client.go                 # TMDb API client
├── db/                            # SQLite database layer
│   ├── db.go                      # Connection + migrations
│   ├── media.go                   # Media CRUD operations
│   ├── config.go                  # Config get/set
│   └── history.go                 # Move + scan history
├── updater/updater.go             # Git-based self-update
├── version/version.go             # Build-time version variables
├── .github/
│   ├── workflows/
│   │   ├── ci.yml                 # Lint + test + vulncheck + failure log
│   │   ├── release.yml            # Cross-compile + GitHub Release
│   │   └── vulncheck.yml          # Weekly vulnerability scan
│   └── logs/
│       └── cicd.log               # AI-readable failure report (auto-generated)
├── run.ps1                        # PowerShell build + deploy pipeline
├── install.ps1                    # Bootstrap installer
├── CHANGELOG.md                   # Release notes
└── spec/                          # Detailed specifications
```

---

<div align="center">

## Data Storage

</div>

All data lives in `./data/`:

```
./data/
├── movie.db                  # SQLite database (WAL mode)
├── thumbnails/               # Downloaded poster images
└── json/
    ├── movie/                # Per-movie JSON metadata
    ├── tv/                   # Per-show JSON metadata
    └── history/              # Move operation logs (RFC3339)
```

---

<div align="center">

## Dependencies

</div>

| Package | Purpose |
|---|---|
| [`github.com/spf13/cobra`](https://github.com/spf13/cobra) | CLI framework |
| [`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite) | Pure-Go SQLite driver (no CGo) |

---

## 🤝 Contributing

Contributions are welcome! Here's how to get started:

1. **Fork** the repository
2. **Create a branch** for your feature or fix:
   ```bash
   git checkout -b feature/my-feature
   ```
3. **Follow the coding guidelines** in [`spec/01-coding-guidelines/`](spec/01-coding-guidelines/)
4. **Keep files small** — one file per command, max ~200 lines
5. **Run tests** before submitting:
   ```bash
   make tidy
   go test ./... -v
   ```
6. **Open a Pull Request** against `main` with a clear description

### Development Setup

```bash
git clone https://github.com/alimtvnetwork/movie-cli-v3.git
cd movie-cli-v3
make tidy
make build
```

See the [Install Guide](spec/03-general/01-install-guide.md) for full setup instructions.

### Code Style

- Go standard formatting (`gofmt`)
- Descriptive variable names, no abbreviations
- Error messages start lowercase, no trailing punctuation
- All new code must pass `go vet` and `golangci-lint`

---

## 📜 Code of Conduct

We are committed to providing a welcoming and inclusive experience for everyone.

**Our Standards:**

- Be respectful, constructive, and empathetic
- Welcome newcomers and help them get started
- Accept constructive criticism gracefully
- Focus on what's best for the project and community

**Unacceptable Behavior:**

- Harassment, trolling, or personal attacks
- Publishing others' private information
- Any conduct that would be considered inappropriate in a professional setting

**Enforcement:** Project maintainers may remove, edit, or reject contributions that violate this code. Repeated violations may result in a temporary or permanent ban.

This Code of Conduct is adapted from the [Contributor Covenant v2.1](https://www.contributor-covenant.org/version/2/1/code_of_conduct/).

---

<div align="center">

## License

</div>

Private project.
