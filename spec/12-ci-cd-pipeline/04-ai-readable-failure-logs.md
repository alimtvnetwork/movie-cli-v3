# 04 — AI-Readable Failure Logs

**Version:** 2.0.0  
**Updated:** 2026-04-10  

## Purpose

Define the pattern for writing structured CI/CD failure logs to the repository so that an AI agent can read, diagnose, and fix failures autonomously.

> **Problem**: When CI fails, the failure details are locked inside GitHub Actions logs — an AI working on the codebase cannot access them. This pattern makes failures **visible in the repo** as committed files inside `.github/logs/errors/`.

---

## Architecture

```
CI starts  → Step 1: delete .github/logs/errors/ folder (fresh slate)
           → Step 2: run all jobs (lint, test, build, vulncheck)

Job fails  → write structured error file to .github/logs/errors/<stage>.log
           → failure-report job commits .github/logs/errors/ with [skip ci]
           → AI reads files, applies fixes, pushes
           → CI re-runs from scratch (folder deleted → fresh logs)

All pass   → folder was already cleared at start → nothing to commit
```

### Key Principle: Fresh Logs Every Run

The `.github/logs/errors/` folder is **deleted at the start of every CI run**. This guarantees:
- No stale logs from previous runs
- Only the current run's failures are present
- If all jobs pass, the folder stays empty (deleted at start, nothing written)

### Infinite Loop Prevention

Three safeguards prevent CI from triggering itself endlessly:

1. **`[skip ci]` commit message** — GitHub Actions ignores commits with this tag
2. **`paths-ignore: [".github/logs/**"]`** — CI trigger excludes the log directory
3. **Folder deleted on every run start** — stale logs never accumulate

---

## Directory Structure

```
.github/
  logs/
    errors/
      lint.log          # go vet + golangci-lint failures
      test.log          # unit/integration test failures
      build.log         # compilation errors
      vulncheck.log     # govulncheck findings
      summary.log       # assembled overview with metadata + AI instructions
```

Each job writes **only its own file**. The `failure-report` job assembles `summary.log` from all individual logs.

---

## Log File Format

### Individual Stage Logs

**Path**: `.github/logs/errors/<stage>.log`  
**Format**: Plain text with structured markers

```
=== LINT ===
status: failure
timestamp: 2026-04-10T12:00:00Z

--- go vet ---
./cmd/movie_scan.go:42:6: unused variable 'x'

--- golangci-lint ---
./db/media.go:105:2: ineffectual assignment to err (ineffassign)
```

### Summary Log

**Path**: `.github/logs/errors/summary.log`  
**Format**: Markdown (human-readable AND AI-parseable)

```markdown
# CI/CD Failure Report

## Run Metadata

| Field         | Value |
|---------------|-------|
| Workflow      | CI |
| Run ID        | 12345678 |
| Commit        | abc1234567 |
| Branch        | main |
| Triggered By  | username |
| Timestamp     | 2026-04-10T12:00:00Z |
| Run URL       | https://github.com/.../actions/runs/12345678 |

## Job Results

| Job | Status | Log File |
|-----|--------|----------|
| Lint | failure | errors/lint.log |
| Vulnerability Scan | success | — |
| Test | failure | errors/test.log |
| Build | skipped | — |

## AI Fix Instructions

1. Read each `.log` file in `.github/logs/errors/`.
2. Identify the root cause from the error output.
3. Apply the fix to the relevant source file.
4. Push — CI will re-run, clear the folder, and only write new failures (if any).

### Common Fix Patterns

| Failure Type | Typical Fix |
|-------------|-------------|
| `go vet` | Fix the reported code issue |
| `golangci-lint` | Fix the violation or add `//nolint` with justification |
| `govulncheck` (3rd-party) | `go get <pkg>@latest && go mod tidy` |
| Test failure | Fix logic or update assertion |
| Build failure | Fix compilation errors |
```

---

## Workflow Implementation

### Step 1: Clear Logs (Every Run)

This step runs **first in every workflow**, before any jobs:

```yaml
jobs:
  clear-logs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6

      - name: Clear previous error logs
        run: |
          if [ -d ".github/logs/errors" ]; then
            rm -rf .github/logs/errors
            git add -A .github/logs/errors
            git diff --cached --quiet || {
              git commit -m "ci: clear error logs [skip ci]"
              git push
            }
          fi
```

All other jobs use `needs: [clear-logs]` to ensure the folder is clean before they run.

### Step 2: Each Job Writes Its Log

Every CI job captures its output to a structured log file and uploads it as an artifact:

```yaml
- name: Run lint
  id: lint
  run: |
    go vet ./... 2>&1 | tee /tmp/govet.out
  continue-on-error: true

- name: Write error log
  if: steps.lint.outcome == 'failure'
  run: |
    mkdir -p /tmp/ci-logs
    {
      echo "=== LINT ==="
      echo "status: failure"
      echo "timestamp: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
      echo ""
      echo "--- go vet ---"
      cat /tmp/govet.out
    } > /tmp/ci-logs/lint.log

- name: Upload log artifact
  if: always()
  uses: actions/upload-artifact@v4
  with:
    name: ci-log-lint
    path: /tmp/ci-logs/lint.log
    if-no-files-found: ignore

- name: Fail if lint failed
  if: steps.lint.outcome == 'failure'
  run: exit 1
```

**Key details:**
- `continue-on-error: true` — lets the log-writing step run even on failure
- Log file is **only created on failure** — successful stages produce no log
- Separate "Fail if..." step — ensures the job still reports failure status
- Artifact naming convention: `ci-log-<stage>` (e.g., `ci-log-lint`, `ci-log-test`)

### Step 3: Failure Report Job

The `failure-report` job:

1. Runs **only on failure** (`if: failure()`)
2. Downloads all `ci-log-*` artifacts
3. Copies them into `.github/logs/errors/`
4. Assembles `summary.log` from all individual logs
5. Commits and pushes with `[skip ci]`

```yaml
failure-report:
  needs: [clear-logs, lint, test, build, vulncheck]
  if: failure()
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v6

    - uses: actions/download-artifact@v4
      with:
        pattern: ci-log-*
        path: /tmp/all-logs
        merge-multiple: true

    - name: Assemble error logs
      run: |
        mkdir -p .github/logs/errors

        # Copy individual stage logs
        cp /tmp/all-logs/*.log .github/logs/errors/ 2>/dev/null || true

        # Build summary
        cat > .github/logs/errors/summary.log << 'HEADER'
        # CI/CD Failure Report

        ## Run Metadata

        | Field | Value |
        |-------|-------|
        | Workflow | ${{ github.workflow }} |
        | Run ID | ${{ github.run_id }} |
        | Commit | ${{ github.sha }} |
        | Branch | ${{ github.ref_name }} |
        | Triggered By | ${{ github.actor }} |
        | Timestamp | $(date -u +%Y-%m-%dT%H:%M:%SZ) |
        | Run URL | ${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }} |

        ## Failed Stages

        HEADER

        for f in .github/logs/errors/*.log; do
          [ "$(basename "$f")" = "summary.log" ] && continue
          echo "### $(basename "$f" .log)" >> .github/logs/errors/summary.log
          echo '```' >> .github/logs/errors/summary.log
          cat "$f" >> .github/logs/errors/summary.log
          echo '```' >> .github/logs/errors/summary.log
          echo "" >> .github/logs/errors/summary.log
        done

        cat >> .github/logs/errors/summary.log << 'FOOTER'

        ## AI Fix Instructions

        1. Read each `.log` file in `.github/logs/errors/`.
        2. Identify the root cause from the error output.
        3. Apply the fix to the relevant source file.
        4. Push — CI will re-run, clear the folder, and only write new failures.

        | Failure Type | Typical Fix |
        |-------------|-------------|
        | `go vet` | Fix the reported code issue |
        | `golangci-lint` | Fix the violation or add `//nolint` with justification |
        | `govulncheck` | `go get <pkg>@latest && go mod tidy` |
        | Test failure | Fix logic or update assertion |
        | Build failure | Fix compilation errors |
        FOOTER

    - name: Commit error logs
      run: |
        git config user.name "github-actions[bot]"
        git config user.email "github-actions[bot]@users.noreply.github.com"
        git add .github/logs/errors/
        git commit -m "ci: write failure logs [skip ci]"
        git push
```

---

## Permissions

The workflow requires `contents: write` to commit the log files back:

```yaml
permissions:
  contents: write
```

The commit uses `${{ secrets.GITHUB_TOKEN }}` (automatic, no setup needed).

---

## AI Agent Workflow

When an AI agent (Lovable, Cursor, Copilot, etc.) works on this repo:

1. **Check for failures**: Look for `.github/logs/errors/summary.log` — if it exists, CI is broken
2. **Read individual logs**: Each file in `errors/` has the raw error output for one stage
3. **Identify the fix**: Use the error messages and "Common Fix Patterns" table
4. **Apply the fix**: Edit the relevant source files
5. **Push**: CI re-runs, clears the folder first, runs all checks fresh

### Example AI Prompt

```
Read .github/logs/errors/summary.log and fix all CI failures.
The individual error details are in the other .log files in the same folder.
Follow the fix patterns in the AI Fix Instructions section.
```

---

## Acceptance Criteria

- GIVEN a new CI run starts WHEN the clear-logs job executes THEN `.github/logs/errors/` is deleted from the repo
- GIVEN a lint failure WHEN CI completes THEN `.github/logs/errors/lint.log` contains the lint error with file and line numbers
- GIVEN a test failure WHEN CI completes THEN `.github/logs/errors/test.log` contains the FAIL output with test names
- GIVEN a vulnerability WHEN CI completes THEN `.github/logs/errors/vulncheck.log` contains affected packages
- GIVEN a build failure WHEN CI completes THEN `.github/logs/errors/build.log` contains the compilation error
- GIVEN any failure WHEN the failure-report job runs THEN `summary.log` is assembled from all individual logs
- GIVEN all jobs pass WHEN CI completes THEN `.github/logs/errors/` does not exist (cleared at start, nothing written)
- GIVEN error logs are committed WHEN the commit message is checked THEN it contains `[skip ci]`
- GIVEN a push to `.github/logs/` WHEN CI trigger evaluates THEN the workflow does NOT run (paths-ignore)

---

*AI-readable failure logs v2 — updated: 2026-04-10*
