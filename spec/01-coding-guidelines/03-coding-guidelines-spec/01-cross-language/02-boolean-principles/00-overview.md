# Boolean Principles

**Version:** 2.6.0  
**Updated:** 2026-03-31  
**AI Confidence:** Production-Ready  
**Ambiguity:** None

---

## Keywords

`02-boolean-principles` · `coding-standards`

---

## Scoring

| Criterion | Status |
|-----------|--------|
| `00-overview.md` present | ✅ |
| AI Confidence assigned | ✅ |
| Ambiguity assigned | ✅ |
| Keywords present | ✅ |
| Scoring table present | ✅ |

---

## Purpose

Previously a single 858-line file, now split into focused modules under 300 lines each.

---

## Document Inventory

| # | File | Purpose | Lines |
|---|------|---------|-------|
| — | [01-naming-prefixes.md](./01-naming-prefixes.md) | P1: is/has prefixes, P2: no negative words | 134 |
| — | [02-guards-and-extraction.md](./02-guards-and-extraction.md) | P3: named guards, P4: extract complex expressions | 205 |
| — | [03-parameters-and-conditions.md](./03-parameters-and-conditions.md) | P5: explicit params, P6: no mixed booleans, P7: no inline statements, P8: no raw system calls | 262 |
| — | [04-quick-reference.md](./04-quick-reference.md) | Quick reference table, common mistakes | 155 |
| — | [05-exemptions-and-api.md](./05-exemptions-and-api.md) | Static factory exemption, Result wrapper API | 139 |
| — | 99-consistency-report.md | — | — |

| — | 99-consistency-report.md | — | — |
---

## Cross-References

- [No Raw Negations](../12-no-negatives.md) — Full guard function inventory
- [Code Style — Rule 3](../04-code-style/00-overview.md) — Complex condition extraction
- [Function Naming](../10-function-naming.md) — No boolean flag parameters
- [PHP Boolean Guard Inventory](../../04-php/07-php-standards-reference/00-overview.md) — PHP-specific helpers
- [Go Boolean Standards](../../03-golang/02-boolean-standards.md) — Go-specific rules and exemptions (P7, P8)
- [Master Coding Guidelines](../15-master-coding-guidelines/00-overview.md) — Consolidated reference
- [Issues & Fixes Log](../01-issues-and-fixes-log.md) — Historical fixes
- [apperror Package — Result Guard Rule](../../../../02-error-manage-spec/02-error-architecture/06-apperror-package/01-apperror-reference/00-overview.md#12-result-guard-rule-mandatory-error-check-before-value-access)

