# App / WordPress Project Template

**Version:** 1.0.0  
**Updated:** 2026-04-03

---

## Overview

App and WordPress projects follow a different folder structure than CLI tools. Instead of the CLI 3-folder pattern (`01-backend`, `02-frontend`, `03-deploy`), app projects use a **fundamentals-first** layout with dedicated `features` and `issues` folders. This document defines the canonical structure for all app-type specifications.

---

## When to Use This Template

Use this template when the spec describes:

- A **WordPress plugin or theme** (PHP + JS/React)
- A **web application** with a UI
- A **platform or product** (not a CLI tool)
- Any module that has **user-facing features and tracked issues**

For CLI tools, use the [CLI Module Template](./04-cli-module-template.md) instead.

---

## Standard App/WordPress Folder Structure

```
{NN}-{app-name}/
├── 00-overview.md                        # Project overview, purpose, and folder index
├── 01-fundamentals.md                    # Core concepts, architecture, data model
│
├── 02-features/                          # All feature specifications
│   ├── 00-overview.md                   # Feature index — lists all features with status
│   │
│   ├── 01-{feature-name}/              # Each feature in its own numbered folder
│   │   ├── 00-overview.md              # Feature purpose, scope, dependencies
│   │   ├── 01-backend.md               # Server-side logic, API endpoints, DB queries
│   │   ├── 02-frontend.md              # UI components, pages, user interactions
│   │   ├── 03-wp-admin.md              # WordPress admin screens, settings, hooks
│   │   ├── 97-acceptance-criteria.md   # Testable criteria for this feature
│   │   └── 99-consistency-report.md
│   │
│   ├── 02-{feature-name}/
│   │   ├── 00-overview.md
│   │   ├── 01-backend.md
│   │   ├── 02-frontend.md
│   │   └── ...
│   │
│   └── ...                              # More features, numbered sequentially
│
├── 03-issues/                            # Tracked issues, bugs, and investigations
│   ├── 00-overview.md                   # Issue index — lists all issues with status
│   │
│   ├── 01-{issue-name}/                # Multi-file issues get their own folder
│   │   ├── 00-overview.md              # Issue description, impact, root cause
│   │   ├── 01-investigation.md         # Analysis, logs, reproduction steps
│   │   ├── 02-resolution.md            # Fix details, code changes, verification
│   │   └── 99-consistency-report.md
│   │
│   ├── 02-{issue-name}.md              # Simple issues can be a single file
│   │
│   └── ...                              # More issues, numbered sequentially
│
├── 97-acceptance-criteria.md             # Consolidated acceptance criteria
├── 98-changelog.md                       # Project changelog
└── 99-consistency-report.md              # Project consistency report
```

---

## Folder Details

### `00-overview.md` (Root)

The first file anyone reads. Must contain:

- Project name, version, status, language/stack
- One-paragraph purpose statement
- Full folder structure tree
- Links to coding guidelines and related specs

### `01-fundamentals.md`

Core architectural decisions that apply across all features:

- Data model and database schema overview
- Authentication and authorization approach
- Plugin/theme lifecycle (for WordPress)
- Key design patterns used
- Configuration and environment setup

### `02-features/`

Each feature lives in its own numbered subfolder. Inside each feature folder:

| File | Purpose | Required? |
|------|---------|-----------|
| `00-overview.md` | Feature scope, dependencies, user stories | ✅ Yes |
| `01-backend.md` | API endpoints, DB queries, server logic | ✅ Yes |
| `02-frontend.md` | UI components, pages, interactions | If applicable |
| `03-wp-admin.md` | Admin screens, settings pages, hooks | WordPress only |
| `97-acceptance-criteria.md` | Testable criteria for this feature | Recommended |
| `99-consistency-report.md` | Consistency check | ✅ Yes |

**Feature overview index** (`02-features/00-overview.md`) must list all features with their status:

```markdown
## Feature Index

| # | Feature | Status | Files |
|---|---------|--------|-------|
| 01 | User Registration | Complete | 4 |
| 02 | Exam Builder | In Progress | 3 |
| 03 | Results Dashboard | Planned | 1 |
```

### `03-issues/`

Tracks bugs, investigations, and resolutions. Each issue follows the same `{NN}-{kebab-name}` convention.

| Format | When to Use |
|--------|------------|
| **Folder** (`{NN}-{issue-name}/`) | Multi-file issues requiring investigation + resolution docs. Must include `00-overview.md`. |
| **Single file** (`{NN}-{issue-name}.md`) | Simple issues with a straightforward description and fix |

**Issue overview index** (`03-issues/00-overview.md`) must list all issues with their status:

```markdown
## Issue Index

| # | Issue | Status | Severity | Files |
|---|-------|--------|----------|-------|
| 01 | Login timeout on slow networks | Resolved | High | 3 |
| 02 | Admin menu not rendering in WP 6.5 | Open | Medium | 1 |
```

---

## Comparison: CLI vs App Templates

| Aspect | CLI Template | App Template |
|--------|-------------|--------------|
| First content file | `01-backend/01-architecture.md` | `01-fundamentals.md` |
| Feature organization | Flat files in `01-backend/` | `02-features/{NN}-{name}/` subfolders |
| Sub-structure per feature | N/A | `01-backend.md`, `02-frontend.md`, `03-wp-admin.md` |
| Issue tracking | External (GitHub Issues) | `03-issues/` folder with indexed issues |
| Deploy specs | `03-deploy/` folder | Not included (handled by deploy tooling specs) |
| WordPress admin | Not applicable | `03-wp-admin.md` inside each feature |

---

## Real-World Example: Exam Manager Plugin

```
13-wp-plugin/03-exam-manager/
├── 00-overview.md                        # Plugin overview, tech stack, folder index
├── 01-fundamentals.md                    # Schema, auth model, plugin lifecycle
│
├── 02-features/
│   ├── 00-overview.md                   # Feature index with 8 features listed
│   ├── 01-exam-builder/
│   │   ├── 00-overview.md              # Exam creation workflow
│   │   ├── 01-backend.md               # REST API, DB tables, validation
│   │   ├── 02-frontend.md              # React exam editor component
│   │   ├── 03-wp-admin.md              # Admin menu registration, settings
│   │   └── 99-consistency-report.md
│   ├── 02-question-bank/
│   │   ├── 00-overview.md
│   │   ├── 01-backend.md
│   │   ├── 02-frontend.md
│   │   └── 99-consistency-report.md
│   ├── 03-results-dashboard/
│   │   ├── 00-overview.md
│   │   ├── 01-backend.md
│   │   └── 02-frontend.md
│   └── 04-certificate-generator/
│       ├── 00-overview.md
│       ├── 01-backend.md
│       └── 02-frontend.md
│
├── 03-issues/
│   ├── 00-overview.md                   # Issue index — 3 tracked issues
│   ├── 01-score-calculation-rounding/
│   │   ├── 00-overview.md              # Bug: scores rounded incorrectly
│   │   ├── 01-investigation.md         # Root cause: float precision
│   │   └── 02-resolution.md            # Fix: decimal(10,4) + round-half-up
│   └── 02-wp-65-compat.md              # Simple issue: single file
│
├── 97-acceptance-criteria.md
├── 98-changelog.md
└── 99-consistency-report.md
```

---

## Checklist: New App/WordPress Module

- [ ] Module number selected (next available after highest)
- [ ] Root `00-overview.md` created with full metadata and folder tree
- [ ] `01-fundamentals.md` written with core architecture
- [ ] `02-features/00-overview.md` created with feature index table
- [ ] At least one feature folder with `00-overview.md` + `01-backend.md`
- [ ] `03-issues/00-overview.md` created (even if empty initially)
- [ ] `97-acceptance-criteria.md` with numbered, testable criteria
- [ ] `99-consistency-report.md` at root and each major subfolder
- [ ] `spec/00-overview.md` master index updated
- [ ] Cross-references added to related modules
- [ ] Error codes registered in error code registry (if applicable)

---

## Cross-References

| Reference | Location |
|-----------|----------|
| CLI Module Template | [04-cli-module-template.md](./04-cli-module-template.md) |
| Folder Structure | [01-folder-structure.md](./01-folder-structure.md) |
| Naming Conventions | [02-naming-conventions.md](./02-naming-conventions.md) |
| Required Files | [03-required-files.md](./03-required-files.md) |
