# Pipeline Specifications

Generic, portable documentation for the project's CI/CD pipeline architecture. These specs describe **what** each pipeline does, **why** each pattern exists, and **how** to implement it — in enough detail for any AI or engineer to reproduce the workflows from scratch.

---

## Documents

| Document | Purpose |
|----------|---------|
| [01-release-pipeline.md](./01-release-pipeline.md) | Release automation: version resolution, binary packaging, install scripts, GitHub releases |
| [02-ci-pipeline.md](./02-ci-pipeline.md) | CI: lint, vulnerability scan, parallel tests, cross-compiled builds, SHA deduplication |

---

## Quick Reference

### Pipeline Triggers

| Workflow | Trigger | Branch/Tag |
|----------|---------|------------|
| CI | Push, Pull Request | `main` |
| Release | Push | `release/**`, `v*` tags |

### Shared Conventions

- **Platform**: GitHub Actions
- **Runner**: `ubuntu-latest`
- **Language toolchain**: Go (version from `go.mod`)
- **Node.js compatibility**: `FORCE_JAVASCRIPT_ACTIONS_TO_NODE24: true` environment variable
- **Action versions**: Pinned to exact tags (e.g., `@v6`), never `@latest` or `@main`
- **Build mode**: Static linking (`CGO_ENABLED=0`) for all binaries
- **Cross-compilation targets**: `windows/amd64`, `windows/arm64`, `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`

### Pinned Tool Versions

| Tool | Version | Used In |
|------|---------|---------|
| `actions/checkout` | `@v6` | Release |
| `actions/setup-go` | `@v6` | Release |
| `softprops/action-gh-release` | `@v2` | Release |

---

*Pipeline specs — updated: 2026-04-09*
