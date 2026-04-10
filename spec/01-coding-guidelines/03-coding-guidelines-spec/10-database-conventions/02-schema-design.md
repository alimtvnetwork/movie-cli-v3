# Database Schema Design

**Version:** 1.0.0  
**Updated:** 2026-04-02

---

## Overview

Rules for designing database schemas that are efficient, maintainable, and correctly normalized. Covers key sizing, primary key strategy, normalization, and the Split DB pattern.

---

## 1. Primary Key Strategy

### 1.1 Always Use Integer Primary Keys

Primary keys MUST be integer-based. Choose the **smallest type** that fits the expected data volume:

| Expected Rows (10-year horizon) | Key Type | Range | Storage |
|--------------------------------|----------|-------|---------|
| < 32,000 | `SMALLINT` | ±32K | 2 bytes |
| < 2 billion | `INTEGER` (default) | ±2.1B | 4 bytes |
| > 2 billion | `BIGINT` | ±9.2 quintillion | 8 bytes |

### 1.2 Decision Flow

```
How many rows in 10 years?
├── < 32,000 → SMALLINT
├── < 2,000,000,000 → INTEGER (this is the default)
└── > 2,000,000,000 → BIGINT
```

### 1.3 UUID/GUID — Avoid

| Aspect | INTEGER | UUID |
|--------|---------|------|
| Storage | 4 bytes | 16 bytes (4x larger) |
| Index performance | Fast (sequential) | Slow (random distribution) |
| Readability | Easy to debug | Hard to read |
| Fragmentation | None | High (random inserts) |

> **Rule:** ❌ Do NOT use UUID/GUID as primary key unless there is an **explicit requirement** (e.g., distributed systems with no central authority, public-facing IDs that must not be guessable).

```sql
-- ❌ AVOID
CREATE TABLE Users (
    UserId TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16))))
);

-- ✅ PREFERRED
CREATE TABLE Users (
    UserId INTEGER PRIMARY KEY AUTOINCREMENT
);
```

### 1.4 When UUID Is Acceptable

Only use UUID when **all** of these are true:
1. Records are created across multiple disconnected systems
2. There is no central ID authority
3. IDs must be publicly exposed and non-guessable
4. The team has explicitly approved UUID for this table

If UUID is used, store it as `BLOB(16)` (not `TEXT(36)`) for storage efficiency.

---

## 2. Key Sizing — Smallest Possible Type

Apply the smallest-type principle to ALL columns, not just primary keys:

| Data | ❌ Oversized | ✅ Right-Sized |
|------|-------------|---------------|
| Status (5 values) | `TEXT` | `TINYINT` + lookup table |
| Age | `INTEGER` | `TINYINT` (0-255) |
| Year | `INTEGER` | `SMALLINT` (0-65535) |
| Boolean | `INTEGER` | `TINYINT(1)` or `BOOLEAN` |
| Country code | `TEXT` | `CHAR(2)` |
| Currency amount | `REAL` | `DECIMAL(10,2)` |

---

## 3. Normalization — Repeated Values Become Tables

### 3.1 The Rule

> **Any column that contains a repeated set of values MUST be extracted into a separate lookup table with a foreign key relationship.**

This applies to: status types, file types, category types, role types, priority levels, etc.

### 3.2 ❌ Wrong — Repeated Strings

```sql
CREATE TABLE Transactions (
    TransactionId INTEGER PRIMARY KEY AUTOINCREMENT,
    Status        TEXT,     -- 'Pending', 'Complete', 'Failed' repeated thousands of times
    FileType      TEXT      -- 'Plugin', 'Theme', 'MuPlugin' repeated thousands of times
);
```

### 3.3 ✅ Correct — Normalized with Lookup Tables

```sql
-- Lookup table for statuses
CREATE TABLE StatusTypes (
    StatusTypeId INTEGER PRIMARY KEY AUTOINCREMENT,
    Name         TEXT NOT NULL UNIQUE   -- 'Pending', 'Complete', 'Failed'
);

-- Lookup table for file types
CREATE TABLE FileTypes (
    FileTypeId INTEGER PRIMARY KEY AUTOINCREMENT,
    Name       TEXT NOT NULL UNIQUE   -- 'Plugin', 'Theme', 'MuPlugin'
);

-- Main table references lookup tables via FK
CREATE TABLE Transactions (
    TransactionId INTEGER PRIMARY KEY AUTOINCREMENT,
    StatusTypeId  INTEGER NOT NULL,
    FileTypeId    INTEGER NOT NULL,
    FOREIGN KEY (StatusTypeId) REFERENCES StatusTypes(StatusTypeId),
    FOREIGN KEY (FileTypeId)   REFERENCES FileTypes(FileTypeId)
);
```

### 3.4 Many-to-Many (N-to-M) Relationships

When an entity can have multiple values of a type (e.g., a user has multiple roles):

```sql
-- Entity tables
CREATE TABLE Users (
    UserId INTEGER PRIMARY KEY AUTOINCREMENT,
    Name   TEXT NOT NULL
);

CREATE TABLE Roles (
    RoleId INTEGER PRIMARY KEY AUTOINCREMENT,
    Name   TEXT NOT NULL UNIQUE   -- 'Admin', 'Editor', 'Viewer'
);

-- Junction table for N-to-M
CREATE TABLE UserRoles (
    UserRoleId INTEGER PRIMARY KEY AUTOINCREMENT,
    UserId     INTEGER NOT NULL,
    RoleId     INTEGER NOT NULL,
    UNIQUE (UserId, RoleId),
    FOREIGN KEY (UserId) REFERENCES Users(UserId),
    FOREIGN KEY (RoleId) REFERENCES Roles(RoleId)
);
```

---

## 4. Database Engine — SQLite First (Split DB)

### 4.1 Default: SQLite with Split DB Pattern

The project follows the **Split DB** pattern: multiple small SQLite databases per domain concern rather than one monolithic database.

| Advantage | Description |
|-----------|-------------|
| Zero config | No server process, no credentials |
| Portable | Single file per database, easy to backup/copy |
| Isolation | Domain failures don't cascade |
| Performance | Each DB has its own WAL, no lock contention across domains |
| Testable | In-memory mode for fast tests |

> See [07-split-db-pattern.md](./07-split-db-pattern.md) for the full Split DB specification including directory layout, DB registry, cross-domain rules, and migration strategy.

### 4.2 Fallback: MySQL

Use MySQL only when:
- High concurrent write volume exceeds SQLite's single-writer model
- Multi-server access to the same database is required
- The application requires server-level replication

All naming conventions (PascalCase) apply equally to MySQL.

---

## 5. Schema Documentation

Every database schema MUST be documented with:

1. **Table purpose** — one-line description
2. **Column definitions** — name, type, constraints, description
3. **Relationships** — FK references with cardinality
4. **Indexes** — which columns and why
5. **Expected volume** — estimated rows in 10 years (drives key sizing)

### Template

```markdown
### TableName

**Purpose:** [What this table stores]  
**Expected volume:** [N rows in 10 years]

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| TableNameId | INTEGER | PK, AUTOINCREMENT | Primary key |
| ForeignTableId | INTEGER | FK, NOT NULL | References ForeignTable |
| Name | TEXT | NOT NULL | Human-readable name |
| CreatedAt | TEXT | NOT NULL, DEFAULT CURRENT_TIMESTAMP | ISO 8601 timestamp |
```

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Naming conventions | [./01-naming-conventions.md](./01-naming-conventions.md) |
| ORM and views | [./03-orm-and-views.md](./03-orm-and-views.md) |
| Testing strategy | [./04-testing-strategy.md](./04-testing-strategy.md) |
| Cross-language DB naming | [../01-cross-language/07-database-naming.md](../01-cross-language/07-database-naming.md) |
