# Issue #01: Hardcoded Timestamp in move-log.json

## Issue Summary

1. **What happened**: The `saveHistoryLog` function wrote `"timestamp":"now"` as a literal string instead of an actual RFC3339 timestamp.
2. **Where**: `cmd/movie_move_helpers.go` (originally `cmd/movie_move.go` line 345-346), `saveHistoryLog` function.
3. **Symptoms and impact**: All move history JSON logs had useless timestamp data. Could not reconstruct when moves happened from JSON logs. DB `move_history.moved_at` was unaffected (uses `CURRENT_TIMESTAMP`).
4. **How discovered**: Manual review of `move-log.json` output during spec creation.

## Root Cause Analysis

1. **Direct cause**: Developer used `"now"` as a placeholder string and never replaced it with a real timestamp call.
2. **Contributing factors**: No code review step, no test that validates JSON log content.
3. **Triggering conditions**: Every `movie move` operation wrote the placeholder.
4. **Why spec did not prevent it**: Original `01-project-spec.md` documented the `"now"` value in Appendix B without flagging it as a bug. Spec described current behavior, not desired behavior.

## Fix Description

1. **Spec change**: Updated Appendix B in `spec.md` to show `time.RFC3339` format. Added "Known Pitfalls" section referencing this issue.
2. **New rule**: Never use placeholder strings (`"now"`, `"TODO"`, `"test"`) in production code paths. Grep for these before marking any feature complete.
3. **Why it resolves root cause**: The production code now calls `time.Now().Format(time.RFC3339)`, producing real timestamps.
4. **Config changes**: None.
5. **Diagnostics**: JSON logs are now parseable with standard timestamp tools.

## Iterations History

1. **Iteration 1** (17-Mar-2026): Replaced `"now"` with `time.Now().Format(time.RFC3339)` in `cmd/movie_move_helpers.go`. Added `"time"` import. Fixed on first attempt.

## Prevention and Non-Regression

1. **Prevention rule**: Grep for `"now"`, `"TODO"`, `"test"`, `"placeholder"` in all `fmt.Sprintf` and string literals before marking a feature complete.
2. **Acceptance criteria**: `move-log.json` entries must contain valid RFC3339 timestamps parseable by `time.Parse(time.RFC3339, ...)`.
3. **Guardrails**: Add to AI success plan Rule 4 (no placeholder values).
4. **Spec references**: `spec/08-app/01-project-spec.md` (Known Pitfalls table, Appendix B).

## TODO and Follow-ups

- [x] Code fix applied
- [x] Spec updated
- [x] Memory updated

## Done Checklist

- [x] Spec updated under `/spec/08-app/`
- [x] Issue write-up created under `/spec/09-app-issues/`
- [x] Memory updated with summary and prevention rule
- [x] Acceptance criteria updated
- [x] Iterations recorded
