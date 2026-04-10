# CI/CD Pipeline Specifications

Generic, portable documentation for the project's CI/CD pipeline architecture. These specs describe **what** each pipeline does, **why** each pattern exists, and **how** to implement it — in enough detail for any AI or engineer to reproduce the workflows from scratch.

> **Reference implementation**: `gitmap-v2` ([spec/pipeline](https://github.com/alimtvnetwork/gitmap-v2/tree/main/spec/pipeline))
> **Workflow files**: `.github/workflows/ci.yml`, `.github/workflows/release.yml`, `.github/workflows/vulncheck.yml`

---

## Documents

| Document | Purpose |
|----------|---------|
| [01-ci-pipeline.md](./01-ci-pipeline.md) | CI: lint, vulnerability scan, parallel tests, cross-compiled builds, SHA deduplication |
| [02-release-pipeline.md](./02-release-pipeline.md) | Release automation: version resolution, binary packaging, install scripts, GitHub releases |
| [03-vulnerability-scanning.md](./03-vulnerability-scanning.md) | Standalone vulnerability scanning: scheduled and manual |
| [04-ai-readable-failure-logs.md](./04-ai-readable-failure-logs.md) | AI-readable failure log pattern: structured log → commit → AI fixes → auto-clear |

---

## Quick Reference

### Pipeline Triggers

| Workflow | Trigger | Branch/Tag |
|----------|---------|------------|
| CI | Push, Pull Request | `main` |
| Release | Push | `release/**`, `v*` tags |
| Vulnerability Scan | Weekly schedule, Manual | Any (default branch) |

### Shared Conventions

- **Platform**: GitHub Actions
- **Runner**: `ubuntu-latest`
- **Language toolchain**: Go (version from `go.mod`)
- **Node.js compatibility**: `FORCE_JAVASCRIPT_ACTIONS_TO_NODE24: true` environment variable
- **Action versions**: Pinned to exact tags (e.g., `@v6`), never `@latest` or `@main`
- **Tool versions**: Pinned to exact versions (e.g., `golangci-lint@v1.64.8`, `govulncheck@v1.1.4`)
- **Build mode**: Static linking (`CGO_ENABLED=0`) for all binaries
- **Cross-compilation targets**: `windows/amd64`, `windows/arm64`, `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`

### Pinned Tool Versions

| Tool | Version | Used In |
|------|---------|---------|
| `golangci-lint` | `v1.64.8` | CI pipeline |
| `govulncheck` | `v1.1.4` | CI pipeline, Vulnerability scan |
| `go-winres` | `v0.3.3` | CI pipeline (Windows icon), Release pipeline |
| `actions/checkout` | `@v6` | All workflows |
| `actions/setup-go` | `@v6` | All workflows |
| `actions/cache` | `@v4` | CI pipeline |
| `actions/upload-artifact` | `@v4` | CI, Release |
| `actions/download-artifact` | `@v4` | CI pipeline |
| `softprops/action-gh-release` | `@v2` | Release pipeline |
| `golangci/golangci-lint-action` | `@v6` | CI pipeline |

---

*Pipeline specs — updated: 2026-04-10*
