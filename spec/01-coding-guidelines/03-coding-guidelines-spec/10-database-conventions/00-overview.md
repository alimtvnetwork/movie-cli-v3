# Database Conventions

**Version:** 1.0.0  
**Status:** Active  
**Updated:** 2026-04-02  
**AI Confidence:** High  
**Ambiguity:** None

---

## Keywords

`database` · `sqlite` · `split-db` · `orm` · `pascalcase` · `primary-key` · `foreign-key` · `views` · `testing` · `naming` · `schema-design`

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

Comprehensive database design and implementation conventions covering naming, schema design, key sizing, ORM usage, view patterns, relationship modeling, and testing strategies. This is the **single source of truth** for how databases are designed and used across all languages.

---

## Golden Rules

> 1. **PascalCase everything** — tables, columns, indexes, views
> 2. **SQLite first** (Split DB pattern) — MySQL as fallback
> 3. **Always use ORMs** — never write raw SQL in business logic
> 4. **Smallest possible key type** — `INTEGER` over `BIGINT`, never UUID unless required
> 5. **Repeated values → separate table** — normalize with foreign key relationships
> 6. **Views for joins** — define DB views instead of on-the-fly joins in code
> 7. **Test with in-memory DB** — unit test schemas, integration test with real queries

---

## Document Index

| # | File | Description |
|---|------|-------------|
| 01 | [01-naming-conventions.md](./01-naming-conventions.md) | PascalCase rules for tables, columns, indexes — references the cross-language spec |
| 02 | [02-schema-design.md](./02-schema-design.md) | Key sizing, primary keys, foreign keys, normalization rules |
| 03 | [03-orm-and-views.md](./03-orm-and-views.md) | ORM-first approach, view patterns, no raw SQL in business logic |
| 04 | [04-testing-strategy.md](./04-testing-strategy.md) | Unit tests for schemas, integration tests with in-memory DB |
| 05 | [05-relationship-diagrams.md](./05-relationship-diagrams.md) | Visual relationship patterns and AI-readable schema diagrams |
| 06 | [06-rest-api-format.md](./06-rest-api-format.md) | PascalCase REST API response format, full CRUD sample, response envelope |
| 07 | [07-split-db-pattern.md](./07-split-db-pattern.md) | Split DB pattern — one SQLite file per bounded context, registry, migrations, cross-domain rules |

---

## Quick Reference

| Topic | Rule |
|-------|------|
| Table names | PascalCase: `AgentSites`, `Transactions` |
| Column names | PascalCase: `PluginSlug`, `CreatedAt` |
| Primary key format | `{TableName}Id` (e.g., `TransactionId`) |
| Primary key type | `INTEGER` default — `BIGINT` only if >2B rows expected |
| UUID/GUID | ❌ Avoid unless explicitly required |
| Repeated values | Normalize into separate table with FK |
| Joins in code | ❌ Use DB views instead |
| Raw SQL in business logic | ❌ Use ORM |
| Default database | SQLite (Split DB pattern) |
| Fallback database | MySQL |
| Schema testing | Unit test + integration test with in-memory DB |

---

## Database Engine Priority

| Priority | Engine | When to Use |
|----------|--------|-------------|
| 1st | **SQLite** (Split DB) | Default for all projects — embedded, zero-config, portable |
| 2nd | **MySQL** | When concurrent write-heavy loads or multi-server access is needed |

> The **Split DB** pattern uses multiple small SQLite databases per domain concern rather than one monolithic database. See [07-split-db-pattern.md](./07-split-db-pattern.md) for the full specification.

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Parent Overview | [../00-overview.md](../00-overview.md) |
| Cross-Language DB Naming | [../01-cross-language/07-database-naming.md](../01-cross-language/07-database-naming.md) |
| File & Folder Naming | [../08-file-folder-naming/00-overview.md](../08-file-folder-naming/00-overview.md) |
| Security Guidelines | [../09-security/00-overview.md](../09-security/00-overview.md) |
| Test Naming Conventions | [../01-cross-language/14-test-naming-and-structure.md](../01-cross-language/14-test-naming-and-structure.md) |

---

*Single source of truth for database design and conventions across all languages.*
