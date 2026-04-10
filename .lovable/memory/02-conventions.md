# Conventions & Preferences

> **Last Updated**: 17-Mar-2026

## Documentation

- Project milestones documented via `readm.txt` (or `readme.txt`) with format:
  `let's start now {dd-MMM-YYYY} {hh:mm AM/PM}` (Malaysia time, UTC+8)
- Full specification maintained in `spec.md` at project root
- App spec in `spec/08-app/`
- Issue write-ups in `spec/09-app-issues/` (format: `01-{issue-slug}.md`)
- Memory and tracking in `.lovable/memory/`

## Issue Tracking Workflow

Every time a fix is made:

1. Create `spec/09-app-issues/XX-{slug}.md` using the template in `00-issue-template.md`
2. Update the relevant spec in `spec/08-app/` with prevention rules
3. Update `.lovable/memory/` with summary and prevention rule
4. All three steps are mandatory — a fix is incomplete without them

## File Naming

- Memory/plan files: `01-name-of-the-file.md` (numbered prefix)
- Issue files: `01-{issue-slug-name}.md` (lowercase, hyphen-separated)
- Keep folder file count small — consolidate where possible
- Completed plans/suggestions move to `completed/` subfolder

## Code Style

- Go standard formatting
- Emoji usage in CLI output (✅ ❌ 📭 🎬 📺 ⭐ etc.)
- Error messages prefixed with `❌`
- Success messages prefixed with `✅`
- Cobra command pattern: one file per command (`movie_<cmd>.go`)
- Max ~200 lines per file
- Explicit methods over boolean flags (single responsibility)
- DRY: use shared helpers, never copy-paste

## Build & Deploy

- **Makefile** for build targets (build, cross-compile, clean, install)
- **build.ps1** for PowerShell automated deploy:
  - Windows default: `E:\bin-run`
  - Mac/Linux default: `/usr/local/bin`
- Version injected via `-ldflags` at build time

## Config Keys

| Key | Default | Purpose |
|---|---|---|
| movies_dir | ~/Movies | Movie file destination |
| tv_dir | ~/TVShows | TV show destination |
| archive_dir | ~/Archive | Archive destination |
| scan_dir | ~/Downloads | Default scan source |
| tmdb_api_key | (none) | TMDb API key |
| page_size | 20 | Items per page in ls |
