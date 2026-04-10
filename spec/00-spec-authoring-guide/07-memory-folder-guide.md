# Memory Folder Guide

**Version:** 2.0.0  
**Updated:** 2026-03-30

---

## Overview

The `.lovable/memories/` directory stores **institutional knowledge** — patterns, conventions, architectural decisions, and workflow rules that persist across AI sessions. This document explains its structure, naming rules, and how to create/update memory files.

---

## Memory Tree Structure

```
.lovable/memories/
├── 00-memory-index.md            # CANONICAL index of all memory files (source of truth)
├── readme.md                     # Simplified high-level overview
├── suggestions.md                # Legacy suggestion tracker
│
├── ai-integration/               # AI/LLM integration patterns
├── architecture/                 # System architecture decisions
├── audit/                        # Audit-related memories
├── completed-issues/             # ✅ Resolved issues archive
├── constraints/                  # Hard constraints and rules
├── done/                         # ✅ Completed tasks archive
├── features/                     # Feature-specific knowledge
├── frontend/                     # Frontend patterns
├── guidelines/                   # Development guidelines
├── logic/                        # Business logic rules
├── patterns/                     # Reusable patterns/templates
├── pending/                      # Work-in-progress items / active tasks
├── planned/                      # ✅ Queued tasks for future execution
├── project/                      # Project-level status/decisions
├── qa/                           # Quality assurance standards
├── reports/                      # Reliability reports, audit reports
├── spec-management/              # Spec management conventions
├── standards/                    # Technical standards
├── style/                        # Code style rules
├── suggestions/                  # Suggestion tracking
│   └── completed/               # Completed suggestions archive
├── technical/                    # Technical implementation details
├── training/                     # Training/learning materials
├── ui/                           # UI component patterns
├── workflow/                     # Workflow processes
└── wp-plugins/                   # WordPress plugin knowledge
```

---

## Naming Conventions

### Folders

Memory folders use **kebab-case WITHOUT numeric prefixes** (unlike spec folders):

```
✅ architecture/
✅ ai-integration/
✅ wp-plugins/

❌ 01-architecture/    # No numeric prefixes in memory folders
```

### Files

Memory files use **kebab-case** and MAY have numeric prefixes for ordering:

```
✅ 01-specification-structure.md
✅ database-standards.md
✅ 02-folder-conventions.md

❌ SpecificationStructure.md    # No PascalCase
❌ specification_structure.md   # No underscores
```

---

## File Format

Every memory file follows this template:

```markdown
# Memory: {category}/{topic}

**Updated:** YYYY-MM-DD  
**Version:** X.Y.Z  
**Status:** Active | Deprecated

---

## Overview

[Concise description of what this memory captures]

---

## [Content Sections]

[The actual knowledge, decisions, patterns, or rules]

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Related Spec | `spec/XX-module/00-overview.md` |
```

---

## Memory Categories Explained

| Category | Purpose | Example Files |
|----------|---------|---------------|
| `architecture/` | System design decisions, module structure, data pipelines | `01-specification-structure.md`, `dashboard-component-architecture.md` |
| `completed-issues/` | Resolved issues archive — bugs fixed, problems closed | `issue-broken-links-resolved.md` |
| `constraints/` | Non-negotiable rules and boundaries | `spec-only-repository.md`, `coding-guidelines.md` |
| `done/` | Completed tasks archive — finished work items | `batch-09-dashboard-update.md` |
| `features/` | Feature-specific implementation knowledge | `dashboard/spec-file-viewer.md` |
| `patterns/` | Reusable templates and patterns | `spec-template.md` |
| `pending/` | Active work-in-progress tasks | `fix-boolean-negatives.md` |
| `planned/` | Queued tasks for future execution | `update-all-overviews.md` |
| `project/` | Project-level status and migration tracking | `spec-migration-status.md` |
| `qa/` | Quality standards and compliance tracking | `consistency-standards.md`, `link-integrity.md` |
| `reports/` | Reliability reports, audit results | `01-reliability-risk-report.md` |
| `workflow/` | How work gets done — processes, conventions, status tracking | `02-folder-conventions.md`, `changelog-management.md` |

---

## Task & Issue Tracking

The memory folder includes dedicated folders for managing work items across their lifecycle:

| Folder | Purpose | When to Use |
|--------|---------|-------------|
| `planned/` | Tasks queued for future work | Upcoming batches, prioritized backlog items |
| `pending/` | Tasks currently in progress or awaiting action | Active work items, blocked tasks |
| `done/` | Completed tasks archive | Move here from `pending/` or `planned/` when finished |
| `completed-issues/` | Resolved issues archive | Bug fixes, closed issues, resolved problems |

### Workflow

```
planned/ → pending/ → done/
                    → completed-issues/ (for issue-type items)
```

1. **Create** a task file in `planned/` with a description
2. **Move** to `pending/` when work begins
3. **Move** to `done/` when complete (or `completed-issues/` for resolved issues)

---

## When to Create a Memory

Create a new memory when:

1. **A pattern is established** that should persist across sessions (e.g., "all database columns use PascalCase")
2. **An architectural decision** is made (e.g., "CLIs use 3-folder structure")
3. **A workflow convention** emerges (e.g., "consistency reports must be updated when files change")
4. **A constraint** is defined (e.g., "this repo is spec-only, no application code")

Do NOT create memories for:
- Information already in spec files (avoid duplication)
- Session-specific debugging notes

> **Note:** Temporary tasks and one-off fixes should go in `planned/` or `pending/`, NOT as standalone memory files.

### Consolidation Rule

> **There is only ONE memory folder: `.lovable/memories/`.** The legacy `.lovable/memory/` variant is prohibited and must be deleted if found.

---

## The Memory Index

`00-memory-index.md` is the **canonical source of truth** for the memory tree. It contains a complete inventory of all memory files with their paths and descriptions.

**Rules:**
- ALWAYS update `00-memory-index.md` when adding or removing memory files
- The readme.md provides a simplified view and does NOT need to list every file
- Cross-references within spec files should point to specific memory files, not the index

---

## Relationship to Spec Tree

| Aspect | `spec/` | `.lovable/memories/` |
|--------|---------|---------------------|
| Purpose | Formal specifications | Institutional knowledge |
| Naming | Numeric prefix required | Numeric prefix optional |
| Required files | `00-overview.md`, `99-consistency-report.md` | No mandatory files |
| Health scoring | ✅ Tracked by dashboard | ❌ Not scored |
| AI consumption | Referenced during implementation | Referenced during all interactions |
| Depth | Up to 3 levels | Up to 2 levels (category/file) |
