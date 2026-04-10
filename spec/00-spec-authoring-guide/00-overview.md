# Spec Authoring Guide

**Version:** 2.0.0  
**Updated:** 2026-03-30  
**Status:** Active  
**AI Confidence:** Production-Ready  
**Ambiguity:** None

---

## Overview

This is the **definitive guide** for any AI agent or human contributor to understand, navigate, create, and maintain specifications in this repository. It covers folder structure, naming conventions, required files, templates, cross-referencing, scoring metrics, reliability validation, and the `.lovable/` institutional memory system — everything needed to produce spec-compliant documentation from scratch.

### How to Write an Overview (`00-overview.md`)

Every module's `00-overview.md` must follow this structure:

1. **Title** — H1 heading matching the module name
2. **Metadata block** — Version, Updated date, Status, AI Confidence Score (0–100%), Ambiguity Score (0–100%)
3. **Overview paragraph** — 1–2 paragraphs explaining what the module covers, its purpose, and who benefits
4. **Scoring section** — AI Confidence Score, Ambiguity Score, and Health Score with thresholds
5. **Keywords section** — Searchable tags for AI discovery
6. **Files table** — Numbered inventory of all files with links and descriptions
7. **Cross-References table** — Links to related modules, memories, and external docs

> **Rule:** The overview is the **entry point**. An AI agent reading only the overview should understand the module's scope, readiness, and how to find any file within it.

---

## Scoring Metrics

Every `00-overview.md` MUST include these three scores:

### AI Confidence

Measures how ready the specification is for an AI agent to implement the described feature.

| Tier | Icon | Meaning | When to Use |
|------|------|---------|-------------|
| **Production-Ready** | ✅ | Specs are complete, unambiguous, and fully implementable | All interfaces defined, acceptance criteria explicit, error codes mapped |
| **High** | 🟢 | Minor gaps exist but AI can proceed with reasonable assumptions | Most sections complete; a few edge cases undefined |
| **Medium** | 🟡 | Significant gaps; AI will need clarification or make risky assumptions | Missing types, partial acceptance criteria, unclear validation rules |
| **Low** | 🔴 | Major sections missing; do NOT attempt implementation | No interfaces, no acceptance criteria, vague requirements |

**Factors that increase confidence:** Complete interfaces, explicit acceptance criteria, error code mappings, clear data models, defined API contracts.

**Factors that decrease confidence:** Missing types, vague requirements ("should handle errors"), undefined edge cases, no acceptance criteria.

### Ambiguity

Measures how much interpretation is required. **Lower tiers are better.**

| Tier | Icon | Meaning | When to Use |
|------|------|---------|-------------|
| **None** | ✅ | No interpretation needed; every detail is explicit | All fields typed, all flows documented, all errors handled |
| **Low** | 🟢 | A few areas need assumptions; acceptable for implementation | Minor gaps in edge cases or optional features |
| **Medium** | 🟡 | Multiple areas require interpretation; review recommended | Several undefined behaviors, partial validation rules |
| **High** | 🟠 | Many areas open to interpretation; rewrite recommended | Missing data models, unclear permissions, vague UI specs |
| **Critical** | 🔴 | Spec is too vague to implement; MUST be rewritten | No clear structure, contradictory requirements, missing core definitions |

**Common ambiguity sources:** Undefined field types, unclear validation rules, missing error handling paths, unspecified permissions, vague UI requirements.

### Health Score (0–100)

Structural compliance score calculated by the dashboard scanner:

| Criterion | Weight |
|-----------|--------|
| `00-overview.md` present | 25% |
| `99-consistency-report.md` present | 25% |
| Lowercase kebab-case naming | 25% |
| Unique numeric sequence prefixes | 25% |

**100/100 = A+** — All four criteria met.

---

## Keywords

`spec-authoring` · `ai-guide` · `folder-structure` · `naming-conventions` · `required-files` · `module-template` · `cli-template` · `cross-references` · `health-score` · `ai-confidence` · `ambiguity-score` · `consistency-report` · `reliability-report` · `memory-folder` · `lovable-folder` · `kebab-case` · `numeric-prefix`

---

## File Categories

Every spec file falls into one of these categories. When creating a new file, assign it the correct category to help AI agents and contributors navigate:

| Category | Purpose | Examples |
|----------|---------|----------|
| **Overview** | Module entry point with metadata and file index | `00-overview.md` |
| **Architecture** | System design, data flow, component structure | `01-architecture.md`, `02-data-model.md` |
| **API / Interface** | Endpoints, contracts, request/response schemas | `03-api-design.md`, `04-rest-endpoints.md` |
| **Logic / Rules** | Business logic, validation rules, algorithms | `05-validation-rules.md`, `06-scoring-logic.md` |
| **UI / Frontend** | Component specs, layouts, user flows | `07-ui-components.md`, `08-user-flows.md` |
| **Backend** | Server-side implementation, database, services | `09-database-schema.md`, `10-service-layer.md` |
| **Diagrams** | Mermaid diagrams, architecture visuals, flow charts | `11-diagrams.md`, `12-sequence-diagrams.md` |
| **Testing** | Test plans, acceptance criteria, QA standards | `97-acceptance-criteria.md` |
| **Meta / Reports** | Consistency reports, changelogs, migration notes | `99-consistency-report.md`, `98-changelog.md` |

> **Rule:** Assign a category in the file's metadata block so AI agents can filter and prioritize what to read.

---

## Folder Structure Examples

### Standard Module (Flat)

```
spec/17-research-queries/
├── 00-overview.md
├── 01-query-types.md
├── 02-search-algorithms.md
├── 97-acceptance-criteria.md
└── 99-consistency-report.md
```

### CLI Tool Module (3-Folder Pattern)

```
spec/09-gsearch-cli/
├── 00-overview.md
├── 01-backend/
│   ├── 01-architecture.md
│   ├── 02-commands.md
│   └── 03-api-design.md
├── 02-frontend/
│   ├── 01-ui-components.md
│   └── 02-user-flows.md
├── 03-diagrams/
│   └── 01-architecture-diagram.md
├── 97-acceptance-criteria.md
└── 99-consistency-report.md
```

### WordPress / App Module (Features + Issues Pattern)

App and WordPress projects use `01-fundamentals.md` as the first content file, then `02-features/` and `03-issues/` folders:

```
spec/13-wp-plugin/03-exam-manager/
├── 00-overview.md
├── 01-fundamentals.md                    # Core architecture, schema, lifecycle
│
├── 02-features/                          # Feature specifications
│   ├── 00-overview.md                   # Feature index with status table
│   ├── 01-exam-builder/
│   │   ├── 00-overview.md
│   │   ├── 01-backend.md
│   │   ├── 02-frontend.md
│   │   └── 03-wp-admin.md
│   └── 02-question-bank/
│       ├── 00-overview.md
│       ├── 01-backend.md
│       └── 02-frontend.md
│
├── 03-issues/                            # Tracked issues and investigations
│   ├── 00-overview.md                   # Issue index with status/severity
│   ├── 01-score-rounding/
│   │   ├── 00-overview.md
│   │   ├── 01-investigation.md
│   │   └── 02-resolution.md
│   └── 02-wp-65-compat.md              # Simple issues can be single files
│
├── 97-acceptance-criteria.md
└── 99-consistency-report.md
```

> **Key insight:** App/WP projects split features into `01-backend.md`, `02-frontend.md`, `03-wp-admin.md` inside each feature folder. Issues follow the same `{NN}-{kebab-name}` convention with `00-overview.md` required for multi-file issues. See [05-app-project-template.md](./05-app-project-template.md) for the full template.

---

## Files

| # | File | Category | Description |
|---|------|----------|-------------|
| 01 | [01-folder-structure.md](./01-folder-structure.md) | Architecture | Complete spec tree layout, layer grouping, and numbering ranges |
| 02 | [02-naming-conventions.md](./02-naming-conventions.md) | Rules | File and folder naming rules (kebab-case, numeric prefixes, reserved ranges) |
| 03 | [03-required-files.md](./03-required-files.md) | Rules | Mandatory files every module must contain (overview, consistency report, etc.) |
| 04 | [04-cli-module-template.md](./04-cli-module-template.md) | Template | Step-by-step template for CLI tool spec modules (3-folder pattern) |
| 05 | [05-app-project-template.md](./05-app-project-template.md) | Template | Template for app/WordPress projects (fundamentals + features + issues) |
| 06 | [06-non-cli-module-template.md](./06-non-cli-module-template.md) | Template | Template for flat/non-CLI modules (research, utilities, standards) |
| 07 | [07-memory-folder-guide.md](./07-memory-folder-guide.md) | Guide | Structure and conventions for the `.lovable/memories/` tree |
| 08 | [08-cross-references.md](./08-cross-references.md) | Rules | How to write cross-references, relative paths, and link integrity rules |
| 09 | [09-exceptions.md](./09-exceptions.md) | Rules | All known exception cases with folder structure examples |

---

## File Naming Convention (Quick Reference)

All files and folders in `spec/` and `.lovable/` MUST use **lowercase kebab-case**:

```
✅ 01-backend/                   ✅ 00-overview.md
✅ 09-gsearch-cli/               ✅ 03-api-design.md
✅ .lovable/memories/workflow/   ✅ file-naming-conventions.md

❌ 01-Backend/                   ❌ ApiDesign.md
❌ 09_gsearch_cli/               ❌ file_naming.md
```

**Numeric prefixes** are mandatory for spec files/folders and optional for memory files. See [02-naming-conventions.md](./02-naming-conventions.md) for full rules.

---

## Cross-Reference Validation Checklist

Every spec file must pass these cross-reference checks:

| # | Check | Rule |
|---|-------|------|
| 1 | **Relative paths only** | Never use root-relative (`/spec/...`) or absolute filesystem paths |
| 2 | **File extension included** | Always end with `.md` — never bare paths |
| 3 | **Target exists** | Every linked file must exist on disk |
| 4 | **Lowercase paths** | Path segments must be entirely lowercase kebab-case |
| 5 | **Depth correctness** | Count `../` levels carefully from source to target |
| 6 | **Bidirectional linking** | If module A references module B, module B should reference A |
| 7 | **Post-rename audit** | After any module renumbering, grep and update ALL references |

**Automated validation:** Run `node scripts/generate-dashboard-data.cjs` — the output JSON reports all broken links. Zero broken links = passing.

See [08-cross-references.md](./08-cross-references.md) for full syntax and examples.

---

## Reliability Check Report

Every module SHOULD include a **reliability risk assessment** to evaluate implementation feasibility before coding begins. These reports are stored in `spec/validation-reports/` or inline within the module.

### What It Covers

| Section | Content |
|---------|---------|
| **Complexity Tier** | Simple / Medium / Complex Agentic / End-to-End |
| **Success Probability** | Estimated % chance of first-pass implementation success |
| **Failure Modes** | Where, why, and how failures can manifest |
| **Risk Mitigations** | Specific actions to reduce failure likelihood |
| **Dependency Risks** | External APIs, shared modules, or integrations that add risk |

### When to Create

- Before implementing any **Complex Agentic** or **End-to-End** module
- When a module's AI Confidence Score is below 70%
- When multiple modules have interdependencies
- After a major spec rewrite or architectural change

---

## The `.lovable/` Folder

The `.lovable/` directory is the **institutional knowledge hub** for the project. It persists AI-learned patterns, decisions, and workflows across sessions.

### Canonical Structure

```
.lovable/
├── memories/                    # ← CANONICAL memory folder (single source of truth)
│   ├── 00-memory-index.md       # Complete inventory of all memory files
│   ├── readme.md                # Simplified high-level overview
│   ├── architecture/            # System design decisions
│   ├── constraints/             # Hard rules (e.g., no-code policy, coding standards)
│   ├── features/                # Feature-specific knowledge
│   ├── guidelines/              # Development guidelines
│   ├── logic/                   # Business logic, formulas, algorithms
│   ├── patterns/                # Reusable templates
│   ├── pending/                 # Work-in-progress / pending tasks
│   ├── planned/                 # Planned tasks (queued for future work)
│   ├── done/                    # Completed tasks archive
│   ├── completed-issues/        # Resolved issues archive
│   ├── project/                 # Project status and tracking
│   ├── qa/                      # Quality standards
│   ├── reports/                 # Reliability reports, audit reports
│   ├── spec-management/         # Spec management conventions
│   ├── suggestions/             # Suggestion tracking
│   │   └── completed/           # Archived completed suggestions
│   ├── training/                # AI training materials
│   ├── ui/                      # UI component patterns
│   ├── workflow/                # Process conventions
│   └── wp-plugins/              # WordPress plugin knowledge
├── plan.md                      # Current execution plan
├── reliability-risk-report.md   # Project-level reliability assessment
└── [other root files]           # Standards archive, audit history, etc.
```

### Task & Issue Tracking Folders

The memory folder includes dedicated folders for tracking work items:

| Folder | Purpose | When to Use |
|--------|---------|-------------|
| `pending/` | Tasks currently in progress or awaiting action | Active work items, blocked tasks |
| `planned/` | Tasks queued for future execution | Upcoming batches, prioritized backlog |
| `done/` | Completed tasks archive | Finished work items (move from pending/planned) |
| `completed-issues/` | Resolved issues archive | Bug fixes, resolved problems, closed issues |

> **Workflow:** Create task files in `planned/` → move to `pending/` when work starts → move to `done/` when complete. Issues follow the same flow but end in `completed-issues/`.

### Consolidation Rule

> **There is only ONE memory folder: `.lovable/memories/`.** The legacy `.lovable/memory/` variant is prohibited. If found during audits, migrate all contents to `.lovable/memories/` and delete the legacy folder.

### Memory File Conventions

| Rule | Detail |
|------|--------|
| **Naming** | Lowercase kebab-case; numeric prefixes optional (unlike spec files) |
| **Format** | Every file: H1 title → metadata (Updated, Version, Status) → Overview → Content → Cross-References |
| **Index** | Always update `00-memory-index.md` when adding/removing files |
| **Depth** | Maximum 2 levels: `memories/{category}/{file}.md` |
| **No duplication** | Don't duplicate information already in spec files |

### What Goes in Memories vs. Specs

| Content | Location |
|---------|----------|
| Formal specifications, APIs, data models | `spec/` |
| Architectural decisions, conventions, patterns | `.lovable/memories/` |
| Execution plans and batch tracking | `.lovable/plan.md` |
| Suggestion tracking | `.lovable/memories/suggestions/` |
| Pending / planned / done tasks | `.lovable/memories/pending/`, `planned/`, `done/` |
| Completed issues | `.lovable/memories/completed-issues/` |
| Reliability assessments | `.lovable/memories/reports/` or `spec/validation-reports/` |

### Where AI Should Write Updates

When an AI agent learns something new or the user provides instructions:

| User Says | AI Writes To |
|-----------|-------------|
| "Remember this pattern" | `.lovable/memories/patterns/` or relevant category |
| "Add this to the plan" | `.lovable/plan.md` |
| "Track this task" | `.lovable/memories/planned/` or `pending/` |
| "This issue is resolved" | Move to `.lovable/memories/completed-issues/` |
| "Update coding guidelines" | `.lovable/memories/constraints/` |
| "Add WP plugin spec" | `spec/XX-wp-plugin-name/` (spec tree) |
| "Remember WP plugin convention" | `.lovable/memories/wp-plugins/` |

See [07-memory-folder-guide.md](./07-memory-folder-guide.md) for the complete memory folder guide.

---

## Quick Start for AI Agents

1. **Read this overview** to understand the file inventory, scoring, and conventions
2. **Read [01-folder-structure.md](./01-folder-structure.md)** for the tree layout
3. **Read [02-naming-conventions.md](./02-naming-conventions.md)** for naming rules
4. **Read [03-required-files.md](./03-required-files.md)** for mandatory file checklist
5. **Choose a template**: [04-cli-module-template.md](./04-cli-module-template.md), [05-app-project-template.md](./05-app-project-template.md), or [06-non-cli-module-template.md](./06-non-cli-module-template.md)
6. **Check [09-exceptions.md](./09-exceptions.md)** for edge cases before creating files
7. **Score the module** — set AI Confidence and Ambiguity percentages in `00-overview.md`
8. **Validate cross-references** — run the link scanner and fix any broken links
9. **Check `.lovable/memories/`** — read relevant memories before writing new specs

---

## Document Inventory

| File |
|------|
| 97-acceptance-criteria.md |
| 98-changelog.md |
| 99-consistency-report.md |


## Cross-References

| Reference | Location |
|-----------|----------|
| Master Index | `../00-overview.md` |
| Folder Structure Guideline | `../00-folder-structure-guideline.md` |
| Coding Guidelines | `../03-coding-guidelines/00-overview.md` |
| Memory Index | `../../.lovable/memories/00-memory-index.md` |
| Reliability Reports | `../validation-reports/` |
| Required Files | `./03-required-files.md` |
| Cross-Reference Rules | `./08-cross-references.md` |
