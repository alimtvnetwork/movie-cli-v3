# CI Pipeline

## Overview

The CI pipeline validates every push and pull request to the `main` branch. It runs linting, vulnerability scanning, parallel test suites, and cross-compiled builds — then caches the result so identical commits are never re-validated.

> **Reference**: Adapted from gitmap-v2 CI patterns ([source](https://github.com/alimtvnetwork/gitmap-v2/blob/main/spec/pipeline/01-ci-pipeline.md))

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

**Why**: If two pushes land on the same branch in quick succession, the first run is cancelled to save resources. Release branches are exempt because every release commit must produce complete artifacts.

---

## Job Graph

```
sha-check ──┬── lint ──────┬── test (matrix) ──┬── test-summary ── build (matrix: 6 targets) ── build-summary
             │              │                    │
             └── vulncheck ─┘                    └── (cache SHA on success)
```

All jobs depend on `sha-check`. If the SHA is already cached, every job prints a skip message and exits successfully — no work is repeated.

---

## Pattern: SHA-Based Build Deduplication (Passthrough Gate)

### Problem

Re-running the pipeline on an already-validated commit wastes compute and blocks PRs with slow re-checks.

### Solution

A gate job (`sha-check`) probes the GitHub Actions cache for a key tied to the commit SHA. Downstream jobs always run (never `if: false` at job level) but use **step-level conditionals** to skip expensive work when the SHA is cached.

**Why not job-level `if`?** GitHub treats skipped jobs as neither success nor failure. Required status checks see "skipped" as not passing — blocking merges. The passthrough pattern ensures every job shows a green checkmark.

### Implementation

```yaml
sha-check:
  name: SHA Check
  runs-on: ubuntu-latest
  outputs:
    already-built: ${{ steps.cache-check.outputs.cache-hit }}
  steps:
    - name: Check SHA cache
      id: cache-check
      uses: actions/cache@v4
      with:
        path: /tmp/ci-passed
        key: ci-passed-${{ github.sha }}
        lookup-only: true
```

Every downstream step uses:

```yaml
- name: Some expensive step
  if: needs.sha-check.outputs.already-built != 'true'
  run: ...
```

### Cache Write

The cache is written by the `test-summary` job on success:

```yaml
- name: Mark SHA as built
  if: success() && needs.sha-check.outputs.already-built != 'true'
  run: mkdir -p /tmp/ci-passed && echo "${{ github.sha }}" > /tmp/ci-passed/sha.txt

- name: Save SHA to cache
  if: success() && needs.sha-check.outputs.already-built != 'true'
  uses: actions/cache/save@v4
  with:
    path: /tmp/ci-passed
    key: ci-passed-${{ github.sha }}
```

**Key rule**: Inline cache writes into the last validation job — never a separate job.

---

## Job: Lint

Runs `go vet ./...` and `golangci-lint` (pinned to `v1.64.8`, 5-minute timeout).

```yaml
- name: Go vet
  run: go vet ./...

- name: golangci-lint
  uses: golangci/golangci-lint-action@v6
  with:
    version: v1.64.8
    args: --timeout=5m
```

---

## Job: Vulnerability Scan

Runs `govulncheck@v1.1.4`. Differentiates:
- **Third-party vulnerabilities** → fail the build
- **Stdlib vulnerabilities** → warn only (unfixable until Go upgrade)

```bash
set +e
govulncheck ./... 2>&1 | tee /tmp/vulncheck.out
rc=$?
if [ $rc -ne 0 ]; then
  if grep "packages you import" /tmp/vulncheck.out | grep -qv "0 vulnerabilities in packages"; then
    echo "::error::Third-party vulnerabilities detected"
    exit 1
  fi
  echo "::warning::Only stdlib vulnerabilities found (no fix available in current Go version)"
fi
```

---

## Job: Test (Matrix)

Runs `go test` with coverage. Uses `fail-fast: false` so all suites complete even if one fails. Uploads test output and coverage as artifacts (7-day retention).

```yaml
strategy:
  fail-fast: false
  matrix:
    include:
      - name: unit
        packages: ./...
```

As the project grows, add separate matrix entries:
```yaml
- name: db
  packages: ./db/...
- name: integration
  packages: ./tests/...
```

---

## Job: Test Summary

Aggregates test results from all matrix entries, prints pass/fail summary, and writes the SHA cache on success.

---

## Job: Build (Matrix)

Cross-compiles for 6 targets with `CGO_ENABLED=0`. For Windows targets, an icon embedding step runs first using `go-winres` (v0.3.3) to generate a `.syso` resource from `assets/icon.ico`. CI builds use `dev-{sha}` versioning. Uploads binaries as artifacts (14-day retention).

```yaml
strategy:
  matrix:
    include:
      - { os: windows, arch: amd64, ext: .exe }
      - { os: windows, arch: arm64, ext: .exe }
      - { os: linux,   arch: amd64, ext: "" }
      - { os: linux,   arch: arm64, ext: "" }
      - { os: darwin,  arch: amd64, ext: "" }
      - { os: darwin,  arch: arm64, ext: "" }
```

---

## Job: Build Summary

Downloads all build artifacts and prints a formatted table of binary names and sizes.

---

## Constraints

- No `cd` in CI steps — use `working-directory` (project root is default for this repo)
- All tool installs use exact version tags — `@latest` is prohibited
- Never use job-level `if` for SHA deduplication — use step-level conditionals
- Inline cache writes into the last validation job — never a separate job
- Module download (`go mod download`) runs as a separate step for caching clarity

---

## Acceptance Criteria

- GIVEN a push to `main` WHEN CI runs THEN lint, vulncheck, test, and build all pass
- GIVEN the same SHA pushed twice WHEN CI runs the second time THEN all jobs skip expensive work and show ✅
- GIVEN a third-party vulnerability WHEN vulncheck runs THEN the build fails with `::error`
- GIVEN a stdlib-only vulnerability WHEN vulncheck runs THEN only a `::warning` is emitted
- GIVEN a test failure WHEN the test matrix runs THEN other suites still complete (fail-fast: false)

---

*CI pipeline spec — updated: 2026-04-10*
