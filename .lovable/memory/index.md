# Project Memory

## Core
Go 1.22 CLI project (NOT web). Binary: `mahin`. Module: `mahin-cli-v2`. Ignore Lovable build errors.
One file per command, max ~200 lines. Shared helpers in movie_info.go and movie_resolve.go.
File naming: `01-name-of-file.md`. Keep folder file counts small.
Plans & suggestions tracked in single files, not per-item files.
Never modify `.release` folder. Any code change bumps at least minor version.
Malaysia timezone (UTC+8) for timestamps. Milestones in `readm.txt`.
Root spec files: lowercase (spec.md, ai-handoff.md, development-log.md). Keep README.md uppercase.
Spec resequenced: foundation 01-06, app at 08, app-issues at 09. Issues in spec/09-app-issues/.
Error spec flattened: spec/02-error-manage-spec/ (no nested subfolder).
cmd/ has 21 files, db/ has 6 files (including tags.go).

## Memories
- [Project overview](mem://01-project-overview) — Go CLI, command tree (21 cmds), architecture, file structure
- [Conventions](mem://02-conventions) — Code style, naming, build, deploy, config keys
- [Plan](mem://workflow/01-plan) — Done/pending task tracker, prioritized backlog
- [AI success plan](mem://workflow/01-ai-success-plan) — 7 rules for 98% AI success rate
- [Suggestions](mem://suggestions/01-suggestions) — Active suggestion tracker (S05 undo prompt is next P0)
- [Reliability report](mem://reports/01-reliability-risk-report) — Failure map, corrective actions, readiness decision
- [Timestamp bug](mem://issues/01-timestamp-bug) — Fixed: hardcoded "now" → RFC3339
- [Duplicate TMDb fetch](mem://issues/02-duplicate-tmdb-fetch) — Fixed: shared helpers
- [Large files](mem://issues/03-large-files) — Fixed: split to <200 lines
