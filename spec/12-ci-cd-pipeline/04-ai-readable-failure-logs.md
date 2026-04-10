# 04 — AI-Readable Failure Logs

## Purpose

Define the pattern for writing structured CI/CD failure logs to the repository so that an AI agent can read, diagnose, and fix failures autonomously.

> **Problem**: When CI fails, the failure details are locked inside GitHub Actions logs — an AI working on the codebase cannot access them. This pattern makes failures **visible in the repo** as a committed file.

---

## Architecture

```
CI fails → each job writes structured .log artifact
         → failure-report job collects all logs
         → assembles .github/logs/cicd.log (Markdown)
         → commits to repo with [skip ci]
         → AI reads file, applies fix, pushes
         → CI re-runs → on success, cicd.log is deleted

CI passes → success-cleanup job removes cicd.log if present
          → commits deletion with [skip ci]
```

### Infinite Loop Prevention

Three safeguards prevent CI from triggering itself endlessly:

1. **`[skip ci]` commit message** — GitHub Actions ignores commits with this tag
2. **`paths-ignore: [".github/logs/**"]`** — CI trigger excludes the log directory
3. **Log is deleted on success** — the file only exists while there's a failure to fix

---

## Log File Format

**Path**: `.github/logs/cicd.log`  
**Format**: Markdown (human-readable AND AI-parseable)

### Structure

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

| Job | Status |
|-----|--------|
| Lint | failure |
| Vulnerability Scan | success |
| Test | success |
| Build | skipped |

## Failure Details

### Stage: lint

\```
=== GO VET ===
status: failure
./cmd/movie_scan.go:42:6: unused variable 'x'

=== GOLANGCI-LINT ===
status: pass
\```

## AI Fix Instructions

1. Read each failure section above.
2. Identify the root cause.
3. Apply the fix to the relevant source file.
4. Push — CI will re-run and this file will be cleared on success.

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

## How Each Job Writes Logs

Every CI job captures its output to a structured log file and uploads it as an artifact:

### Pattern

```yaml
- name: Run step
  id: step-id
  run: |
    some-command 2>&1 | tee /tmp/step.out
  continue-on-error: true

- name: Write log
  if: always()
  run: |
    mkdir -p /tmp/ci-logs
    {
      echo "=== STAGE NAME ==="
      echo "status: ${{ steps.step-id.outcome }}"
      cat /tmp/step.out
    } > /tmp/ci-logs/stage.log

- name: Upload log
  if: always()
  uses: actions/upload-artifact@v4
  with:
    name: ci-log-stage
    path: /tmp/ci-logs/stage.log

- name: Fail if step failed
  if: steps.step-id.outcome == 'failure'
  run: exit 1
```

**Key details:**
- `continue-on-error: true` — lets the log-writing step run even on failure
- Separate "Fail if..." step — ensures the job still reports failure status
- Artifact naming convention: `ci-log-<stage>` (e.g., `ci-log-lint`, `ci-log-test-unit`)

---

## Failure Report Job

The `failure-report` job:

1. Runs **only on failure** (`if: failure()`)
2. Downloads all `ci-log-*` artifacts
3. Assembles them into a single `.github/logs/cicd.log` file
4. Commits and pushes with `[skip ci]`

```yaml
failure-report:
  needs: [all-jobs]
  if: failure()
  steps:
    - uses: actions/checkout@v6
    - uses: actions/download-artifact@v4
      with:
        pattern: ci-log-*
    - name: Assemble cicd.log
      run: # ... merge all logs into .github/logs/cicd.log
    - name: Commit
      run: |
        git add .github/logs/cicd.log
        git commit -m "ci: write failure log [skip ci]"
        git push
```

---

## Success Cleanup Job

The `success-cleanup` job:

1. Runs **only on success** (`if: success()`)
2. Checks if `.github/logs/cicd.log` exists
3. If yes, deletes it and commits with `[skip ci]`

This ensures the log file is **only present when there's an active failure**.

---

## Permissions

The workflow requires `contents: write` to commit the log file back:

```yaml
permissions:
  contents: write
```

The commit uses `${{ secrets.GITHUB_TOKEN }}` (automatic, no setup needed).

---

## AI Agent Workflow

When an AI agent (Lovable, Cursor, Copilot, etc.) works on this repo:

1. **Check for failures**: Read `.github/logs/cicd.log` — if it exists, CI is broken
2. **Parse the file**: Each `### Stage:` section contains the raw error output
3. **Identify the fix**: Use the error messages and "Common Fix Patterns" table
4. **Apply the fix**: Edit the relevant source files
5. **Push**: CI re-runs automatically; if it passes, `cicd.log` is deleted

### Example AI Prompt

```
Read .github/logs/cicd.log and fix all CI failures.
Follow the fix patterns in the AI Fix Instructions section.
```

---

## Acceptance Criteria

- GIVEN a lint failure WHEN CI completes THEN `.github/logs/cicd.log` contains the lint error output with file and line numbers
- GIVEN a test failure WHEN CI completes THEN `cicd.log` contains the FAIL output with test names and assertions
- GIVEN a vulnerability WHEN CI completes THEN `cicd.log` contains the govulncheck output listing affected packages
- GIVEN a build failure WHEN CI completes THEN `cicd.log` contains the compilation error
- GIVEN all jobs pass WHEN CI completes THEN `.github/logs/cicd.log` is deleted from the repo
- GIVEN `cicd.log` is committed WHEN the commit message is checked THEN it contains `[skip ci]`
- GIVEN a push to `.github/logs/` WHEN CI trigger evaluates THEN the workflow does NOT run (paths-ignore)

---

*AI-readable failure logs — updated: 2026-04-10*
