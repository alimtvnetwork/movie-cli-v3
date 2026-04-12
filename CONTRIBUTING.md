# Contributing to Movie CLI

Thank you for your interest in contributing! This guide covers everything you need to get started.

---

## Table of Contents

- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Code Guidelines](#code-guidelines)
- [Commit Messages](#commit-messages)
- [Pull Requests](#pull-requests)
- [Issue Reporting](#issue-reporting)
- [Architecture Overview](#architecture-overview)
- [Testing](#testing)
- [Release Process](#release-process)

---

## Getting Started

### Prerequisites

| Tool | Minimum Version | Install |
|------|----------------|---------|
| **Go** | 1.22+ | [go.dev/dl](https://go.dev/dl/) |
| **Git** | 2.x | [git-scm.com](https://git-scm.com/) |
| **PowerShell** | 7+ (cross-platform) | `brew install --cask powershell` |

### Setup

```bash
# Fork the repo on GitHub, then:
git clone https://github.com/<your-username>/movie-cli-v3.git
cd movie-cli-v3
go mod tidy
make build
```

Verify:

```bash
./mahin version
```

See the [Install Guide](spec/03-general/01-install-guide.md) for detailed setup instructions.

---

## Development Workflow

1. **Create a branch** from `main`:
   ```bash
   git checkout -b feature/my-feature
   ```

2. **Make your changes** — keep each commit focused on a single concern.

3. **Run checks locally** before pushing:
   ```bash
   go vet ./...
   go test ./... -v -count=1
   ```

4. **Push and open a Pull Request** against `main`.

### Branch Naming

| Prefix | Purpose | Example |
|--------|---------|---------|
| `feature/` | New functionality | `feature/export-csv` |
| `fix/` | Bug fixes | `fix/scan-empty-dir` |
| `refactor/` | Code restructuring | `refactor/split-move-cmd` |
| `docs/` | Documentation only | `docs/update-readme` |
| `release/` | Release branches (maintainers) | `release/v1.4.0` |

---

## Code Guidelines

### General Rules

- **One file per command**, max ~200 lines
- Shared helpers go in `movie_info.go` and `movie_resolve.go`
- Use descriptive variable names — no abbreviations
- All exported functions must have doc comments

### Go Style

- Format with `gofmt` (enforced by CI)
- Pass `go vet ./...` and `golangci-lint` with zero warnings
- Error messages: lowercase, no trailing punctuation
  ```go
  // ✅ Good
  return fmt.Errorf("file not found: %s", path)

  // ❌ Bad
  return fmt.Errorf("File not found: %s.", path)
  ```

### Error Handling

- Always handle errors explicitly — never use `_` for error returns
- Wrap errors with context using `fmt.Errorf("operation: %w", err)`
- Use early returns to reduce nesting

### File & Folder Naming

- Spec files: `01-name-of-file.md` (numbered prefix)
- Go files: `snake_case.go`
- Root spec files: lowercase (`spec.md`, `ai-handoff.md`)

### Dependencies

- Prefer stdlib over third-party packages
- Any new dependency requires justification in the PR description
- Current deps: Cobra (CLI), modernc.org/sqlite (no CGo)

---

## Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <short description>

[optional body]
```

### Types

| Type | When to use |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation changes |
| `refactor` | Code change that neither fixes a bug nor adds a feature |
| `test` | Adding or updating tests |
| `chore` | Build, CI, tooling changes |

### Examples

```
feat(scan): add recursive directory scanning
fix(move): handle cross-device rename with fallback copy
docs(readme): add build-all target to make table
chore(ci): pin golangci-lint to v1.64.8
```

---

## Pull Requests

### Before Submitting

- [ ] Code compiles: `make build`
- [ ] Tests pass: `go test ./... -v -count=1`
- [ ] Linter clean: `go vet ./...`
- [ ] No new files exceed 200 lines
- [ ] Commit messages follow conventions
- [ ] Any new dependency is justified

### PR Description Template

```markdown
## What

Brief description of the change.

## Why

Context and motivation.

## How

Implementation approach and key decisions.

## Testing

How this was tested (commands run, edge cases checked).
```

### Review Process

1. All PRs require at least one approval
2. CI must pass (lint, vulncheck, tests, cross-compile)
3. Maintainers may request changes — address each comment
4. Squash-merge is preferred for clean history

---

## Issue Reporting

### Bug Reports

Include:
- **Movie CLI version**: output of `movie version`
- **OS and architecture**: e.g., Windows 11 amd64, macOS 15 arm64
- **Steps to reproduce**: minimal, specific commands
- **Expected vs actual behavior**
- **Relevant logs or error output**

### Feature Requests

Include:
- **Use case**: what problem does this solve?
- **Proposed solution**: how should it work?
- **Alternatives considered**: what else could work?

---

## Architecture Overview

```
movie-cli-v3/
├── main.go              # Entry point
├── cmd/                  # One file per CLI command
│   ├── scan.go
│   ├── ls.go
│   ├── move.go
│   └── ...
├── movie_info.go         # Shared TMDb/metadata helpers
├── movie_resolve.go      # Shared resolution helpers
├── version/              # Version info (injected at build)
├── data/                 # Runtime data (SQLite DB, posters)
├── spec/                 # Project specifications
├── assets/               # Icons and static assets
└── .github/workflows/    # CI and release pipelines
```

### Key Design Decisions

- **Pure-Go SQLite** (`modernc.org/sqlite`) — no CGo, single binary
- **WAL mode** — concurrent reads during scans
- **TMDb API** — user provides their own API key
- **Self-contained binary** — no external runtime dependencies

---

## Testing

### Running Tests

```bash
# All tests
go test ./... -v -count=1

# With coverage
go test ./... -v -coverprofile=coverage.out -covermode=atomic

# Specific package
go test ./cmd/... -v

# Via Makefile
make tidy
```

### Writing Tests

- Place tests in `_test.go` files alongside the code
- Use table-driven tests for multiple cases
- Mock external APIs (TMDb) — never call real APIs in tests
- Test error paths, not just happy paths

---

## Release Process

Releases are automated via GitHub Actions. Only maintainers create releases.

### Creating a Release

```bash
# Option A: release branch
git checkout -b release/v1.4.0
git push origin release/v1.4.0

# Option B: tag
git tag v1.4.0
git push origin v1.4.0
```

### What Happens

1. CI validates (lint, vuln scan, tests)
2. Cross-compiles 6 binaries (Windows/Linux/macOS × amd64/arm64)
3. Embeds icon into Windows binaries
4. Generates platform-specific install scripts
5. Creates GitHub Release with checksums and changelog

### Version Bumping

- Any code change bumps at least the **minor** version
- Bug-only fixes bump the **patch** version
- Breaking changes bump the **major** version

See [Release Pipeline Spec](spec/pipeline/01-release-pipeline.md) for details.

---

## Maintainer

### [Md. Alim Ul Karim](https://www.google.com/search?q=alim+ul+karim)

**[Creator & Lead Architect](https://alimkarim.com)** | [Chief Software Engineer](https://www.google.com/search?q=alim+ul+karim), [Riseup Asia LLC](https://riseup-asia.com)

|  |  |
|---|---|
| **Website** | [alimkarim.com](https://alimkarim.com/) · [my.alimkarim.com](https://my.alimkarim.com/) |
| **LinkedIn** | [linkedin.com/in/alimkarim](https://linkedin.com/in/alimkarim) |
| **Stack Overflow** | [stackoverflow.com/users/361646/alim-ul-karim](https://stackoverflow.com/users/361646/alim-ul-karim) |
| **Google** | [Alim Ul Karim](https://www.google.com/search?q=Alim+Ul+Karim) |
| **Role** | Chief Software Engineer, [Riseup Asia LLC](https://riseup-asia.com) |

---

## Questions?

Open an issue or start a discussion. We're happy to help!
