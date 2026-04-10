# Coding Guidelines

**Version:** 2.4.0  
**Updated:** 2026-04-02  
**AI Confidence:** Production-Ready  
**Ambiguity:** None

---

## Purpose

Consolidated coding standards and conventions organized by category. This folder is the **single canonical location** for all language-specific and cross-language coding guidelines, including file naming, security policies, and database design conventions.

---

## Keywords

`coding-standards` · `cross-language` · `typescript` · `golang` · `php` · `rust` · `csharp` · `naming-conventions` · `boolean-patterns` · `dry` · `strict-typing` · `file-naming` · `folder-naming` · `security` · `dependency-pinning` · `database` · `orm` · `schema-design` · `slug-conventions`

---

## Scoring

| Metric | Value |
|--------|-------|
| AI Confidence | Production-Ready |
| Ambiguity | None |
| Health Score | 100/100 (A+) |

---

## Categories

### Language & Cross-Language Standards

| # | Category | Description | Files |
|---|----------|-------------|-------|
| 01 | [Cross-Language](./01-cross-language/00-overview.md) | Language-agnostic rules: DRY, naming, booleans, typing, complexity, lazy eval, regex, mutation, null safety, nesting, slugs | 29 |
| 02 | [TypeScript](./02-typescript/00-overview.md) | TypeScript enum patterns, type safety, promise/await patterns | 13 |
| 03 | [Golang](./03-golang/00-overview.md) | Go coding standards, enum specification, boolean rules, defer, internals, severity | 16 |
| 04 | [PHP](./04-php/00-overview.md) | PHP coding standards, enums, forbidden patterns, naming, spacing/imports, ResponseKeyType | 12 |
| 05 | [Rust](./05-rust/00-overview.md) | Rust standards: naming, error handling, async, memory safety, FFI | 10 |
| 06 | [AI Optimization](./06-ai-optimization/00-overview.md) | Anti-hallucination rules, AI quick-reference checklist, common AI mistakes, enum naming reference | 8 |
| 07 | [C#](./07-csharp/00-overview.md) | C# standards: naming, method design, error handling, type safety | 5 |

### Infrastructure & Convention Standards

| # | Category | Description | Files |
|---|----------|-------------|-------|
| 08 | [File & Folder Naming](./08-file-folder-naming/00-overview.md) | Per-language file and folder naming conventions (PHP/WordPress, Go, TS/JS, Rust, C#) | 7 |
| 09 | [Security](./09-security/00-overview.md) | Security policies, dependency pinning (Axios), vulnerability tracking | 6 |
| 10 | [Database Conventions](./10-database-conventions/00-overview.md) | Schema design, PascalCase naming, ORM usage, views, key sizing, REST API format, testing | 8 |

---

## Migration History

| Date | Change |
|------|--------|
| 2026-04-02 | Added `10-database-conventions/` (8 files: schema design, ORM, views, testing, REST API format) |
| 2026-04-02 | Added `09-security/` and moved Axios version control from `spec/08-app/` |
| 2026-04-02 | Added `08-file-folder-naming/` (per-language conventions) |
| 2026-04-02 | Added `28-slug-conventions.md` to cross-language |
| 2026-03-31 | Consolidated 5 guideline sources into this canonical location |

---

## Document Inventory

| File | Description |
|------|-------------|
| [01-consolidated-review-guide.md](../01-consolidated-review-guide.md) | Full code review guide with examples (all languages) |
| [02-consolidated-review-guide-condensed.md](../02-consolidated-review-guide-condensed.md) | One-liner bullet-point checklist for quick scanning |
| 97-acceptance-criteria.md | |
| 99-consistency-report.md | |

---

## Cross-References

- [Spec Authoring Guide](../../00-spec-authoring-guide/00-overview.md)
- [Error Management Spec](../../02-error-manage-spec/00-overview.md)
- Consolidation Plan (archived)
- [Consolidated Review Guide](../01-consolidated-review-guide.md)
- [Condensed Review Checklist](../02-consolidated-review-guide-condensed.md)
