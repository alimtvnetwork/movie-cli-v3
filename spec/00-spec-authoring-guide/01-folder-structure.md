# Folder Structure

**Version:** 1.1.0  
**Updated:** 2026-04-08

---

## Overview

The `spec/` directory is the canonical location for all project specifications. It uses a **numbered folder hierarchy** organized into functional layers. This document explains the complete tree layout, how modules are grouped, and how the numbering scheme works.

---

## Root-Level Tree

```
spec/
├── 00-overview.md                          # Master index — links every module
├── 99-consistency-report.md               # Root-level consistency report
│
├── 00-spec-authoring-guide/                 # THIS GUIDE — how to write specs
├── 01-coding-guidelines/                   # Consolidated coding standards (5 sub-categories)
├── 02-error-manage-spec/                   # Error management, codes, architecture
├── 03-split-db-architecture/               # Split database pattern
├── 04-seedable-config-architecture/        # Seedable configuration with versioning
├── 05-design-spec/                         # Design specifications (Mermaid styling)
├── 06-diagrams/                            # Mermaid flow diagrams for all CLI commands
│
├── 08-app/                                  # Application project spec (movie-cli)
├── 09-app-issues/                           # Application issue tracking and resolution
│
├── 10-gsearch-cli/                         # CLI: Google/Bing search + BI suite
├── 11-brun-cli/                            # CLI: Build runner
├── 12-ai-bridge-cli/                       # CLI: LLM orchestration
├── 13-nexus-flow-cli/                      # CLI: Visual workflow editor
├── 14-wp-plugin/                           # WordPress plugins
├── 15-wp-plugin-builder/                   # CLI: WP plugin scaffolding
├── 16-spec-reverse-cli/                    # CLI: Code-to-spec reverse engineering
├── 17-ai-transcribe-cli/                   # CLI: Speech-to-text/text-to-speech
├── 18-ai-research/                         # AI research notes
├── 19-error-resolution/                    # Error handling patterns
├── 20-license-manager/                     # CLI: License key management
├── 21-shared-cli-frontend/                 # Shared React frontend for CLIs
├── 22-wp-seo-publish-cli/                  # CLI: WordPress SEO publishing
│
├── 24-how-app-issues-track/                # Issue tracking & prevention
│
├── 28-wp-plugin-development/               # WordPress plugin patterns
├── 29-upload-scripts/                      # Plugin upload/deployment scripts
├── 30-e2-activity-feed/                    # Activity feed spec
├── 31-generic-enforce/                     # Automated standard enforcement
├── 32-shared-preset-data/                  # Shared preset/seed data
├── 33-ai-bridge-non-vector-rag/            # Non-vector RAG system
├── 34-time-log-cli/                        # CLI: OS-level activity tracker (Rust)
├── 35-time-log-ui/                         # Web dashboard for Time Log
├── 36-time-log-combined/                   # Combined acceptance criteria
│
├── 99-archive/                             # Deprecated specifications
└── validation-reports/                     # Audit and validation artifacts
```

> **Note:** Numbers 07 is reserved for future foundation modules. The app project is at 08, app issues at 09. CLI tools start from 10+. Some numbers are non-contiguous for historical reasons.

---

## Functional Layers

Modules are organized into six conceptual layers:

### Layer 1: Foundation & Standards (00–06)

Foundational rules that all other modules depend on.

| # | Module | Purpose |
|---|--------|---------|
| 00 | spec-authoring-guide | This guide — how to write and maintain specs |
| 01 | coding-guidelines | Consolidated language standards (cross-language, TS, Go, PHP, Rust) |
| 02 | error-manage-spec | Error management, error codes, error architecture |
| 03 | split-db-architecture | Split database pattern (operational + config DBs) |
| 04 | seedable-config-architecture | Seedable configuration with changelog versioning |
| 05 | design-spec | Design specifications (Mermaid styling, visual guidelines) |
| 06 | diagrams | Mermaid flow diagrams for all CLI commands |

### Layer 2: Application (08–09)

Application project specifications and issue tracking.

| # | Module | Purpose |
|---|--------|---------|
| 08 | app | Application project spec (movie-cli) |
| 09 | app-issues | Issue tracking and resolution |

### Layer 3: CLI Tools (10–22)

Individual CLI tool specifications, each following the 3-folder pattern.

| # | Module | Language |
|---|--------|----------|
| 10 | gsearch-cli | Go |
| 11 | brun-cli | Go |
| 12 | ai-bridge-cli | Go |
| 13 | nexus-flow-cli | Go |
| 15 | wp-plugin-builder | Go |
| 16 | spec-reverse-cli | Go |
| 17 | ai-transcribe-cli | Go |
| 20 | license-manager | Go |
| 22 | wp-seo-publish-cli | Go |
| 31 | generic-enforce | Go |
| 34 | time-log-cli | Rust |

### Layer 4: WordPress (14, 28–29)

WordPress plugin specifications and deployment tooling.

| # | Module | Purpose |
|---|--------|---------|
| 14 | wp-plugin | Plugin specs (exam manager, link manager, etc.) |
| 28 | wp-plugin-development | Plugin development patterns |
| 29 | upload-scripts | Plugin deployment utilities |

### Layer 5: Shared Modules & Frontend (21, 32, 35–36)

Shared libraries, data, and UI dashboards.

| # | Module | Purpose |
|---|--------|---------|
| 21 | shared-cli-frontend | Shared React frontend for all CLIs |
| 32 | shared-preset-data | Shared preset/seed data files |
| 35 | time-log-ui | Web dashboard for Time Log data |
| 36 | time-log-combined | Combined acceptance criteria for Time Log |

### Layer 6: Archive & Governance (99, validation-reports)

Historical and governance artifacts.

| # | Module | Purpose |
|---|--------|---------|
| 99 | archive | Deprecated and superseded specs |
| — | validation-reports | Audit certificates and validation logs |

---

## Reserved Number Ranges

| Range | Purpose |
|-------|---------|
| 00 | Root files (overview) and `00-overview.md` within folders |
| 01–06 | Foundation & standards |
| 07 | Reserved for future foundation module |
| 08–09 | Application project and issues |
| 10–22 | CLI tools and core platform modules |
| 23 | Reserved / skipped |
| 24–36 | Research, WordPress, utilities, enforcement, shared data |
| 37–89 | Future modules |
| 90–95 | Meta documents (if needed) |
| 96 | AI context files |
| 97 | Acceptance criteria |
| 98 | Changelogs |
| 99 | Consistency reports and archive |

---

## Module Templates

Two templates exist depending on project type:

| Type | Template | Pattern |
|------|----------|---------|
| **CLI tool** | [CLI Module Template](./04-cli-module-template.md) | `01-backend/`, `02-frontend/`, `03-deploy/` |
| **App / WordPress** | [App Project Template](./05-app-project-template.md) | `01-fundamentals.md`, `02-features/`, `03-issues/` |

- **CLI tools** — Go/Rust command-line tools with backend, optional frontend, and deploy specs
- **App/WordPress projects** — Plugins, themes, or web apps with features and issue tracking. Features use `01-backend.md`, `02-frontend.md`, `03-wp-admin.md` sub-structure.

---

## Subfolder Depth Rules

- **Maximum depth:** 3 levels (e.g., `spec/14-wp-plugin/03-exam-manager/02-features/01-exam-builder/`)
- **Each level** follows the same `{NN}-{name}` convention
- **Every folder** at any depth must have `00-overview.md`
- **Only top-level** and major subfolder boundaries require `99-consistency-report.md`

---

## Non-Contiguous Numbering

Some numbers are intentionally skipped:
- **07** — Reserved for future foundation module
- **23** — Reserved (previously used, now skipped)

These gaps are preserved to avoid mass-renaming existing modules. New modules should use the **next available number** after the highest existing module.
