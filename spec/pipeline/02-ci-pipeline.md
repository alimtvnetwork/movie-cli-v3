# CI Pipeline

## Overview

The CI pipeline validates every push and pull request to the `main` branch. It runs linting, vulnerability scanning, parallel test suites, and cross-compiled builds — then caches the result so identical commits are never re-validated.

---

## Trigger and Concurrency

### Trigger

```yaml
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
```

### Concurrency Control

Scope concurrent runs by branch reference. Cancel superseded runs on feature branches, but **never cancel release branches** (they must always run to completion):

```yaml
concurrency:
  group: ci-${{ github.ref }}
  cancel-in-progress: ${{ !startsWith(github.ref, 'refs/heads/release/') }}
```

---

## Job Graph

```
sha-check ──┬── lint ──────┬── test (matrix) ──┬── test-summary ── build (matrix: 6 targets) ── build-summary
             │              │                    │
             └── vulncheck ─┘                    └── (cache SHA on success)
```

All jobs depend on `sha-check`. If the SHA is already cached, every job prints a skip message and exits successfully — no work is repeated.

---

## Pattern: SHA-Based Build Deduplication

A gate job probes the GitHub Actions cache for a key tied to the commit SHA. Downstream jobs always run (never `if: false` at job level) but use **step-level conditionals** to skip expensive work when the SHA is cached.

**Why not job-level `if`?** GitHub treats skipped jobs as neither success nor failure. Required status checks see "skipped" as not passing — blocking merges. The passthrough pattern ensures every job shows a green checkmark.

---

## Job: Lint

Runs `go vet ./...` and `golangci-lint` (pinned to `v1.64.8`, 5-minute timeout).

## Job: Vulnerability Scan

Runs `govulncheck@v1.1.4`. Differentiates:
- **Third-party vulnerabilities** → fail the build
- **Stdlib vulnerabilities** → warn only (unfixable until Go upgrade)

## Job: Test (Matrix)

Runs `go test` with coverage. Uses `fail-fast: false` so all suites complete even if one fails. Uploads test output and coverage as artifacts.

## Job: Test Summary

Aggregates test results, prints pass/fail summary, and writes the SHA cache on success.

## Job: Build (Matrix)

Cross-compiles for 6 targets with `CGO_ENABLED=0`. For Windows targets, an icon embedding step runs first using `go-winres` (v0.3.3) to generate a `.syso` resource from `assets/icon.png`. CI builds use `dev-{sha}` versioning. Uploads binaries as artifacts (14-day retention).

## Job: Build Summary

Downloads all build artifacts and prints a formatted table of binary names and sizes.

---

## Constraints

- No `cd` in CI steps — use `working-directory`
- All tool installs use exact version tags — `@latest` is prohibited
- Never use job-level `if` for SHA deduplication — use step-level conditionals
- Inline cache writes into the last validation job — never a separate job

---

*CI pipeline spec — updated: 2026-04-09*
