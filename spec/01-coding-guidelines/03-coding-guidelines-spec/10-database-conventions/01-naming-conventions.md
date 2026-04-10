# Database Naming Conventions

**Version:** 1.1.0  
**Updated:** 2026-04-02

---

## Overview

All database objects use **PascalCase**. This document summarizes the rules and references the full cross-language database naming spec.

> **Full specification:** [../01-cross-language/07-database-naming.md](../01-cross-language/07-database-naming.md)

---

## Summary of Rules

| Object | Convention | Example |
|--------|-----------|---------|
| Table names | PascalCase, no prefix | `AgentSites`, `Transactions` |
| Column names | PascalCase | `PluginSlug`, `CreatedAt` |
| Primary key | `{TableName}Id` | `TransactionId`, `AgentSiteId` |
| Foreign key column | Same name as referenced PK | `AgentSiteId` references `AgentSites.AgentSiteId` |
| Boolean columns | `Is` or `Has` prefix, **positive only** | `IsActive`, `HasLicense` |
| Index names | `Idx{Table}_{Column}` | `IdxTransactions_CreatedAt` |
| View names | PascalCase with `Vw` prefix | `VwTransactionDetails`, `VwActiveAgentSites` |
| Abbreviations | First letter only capitalized | `Id`, `Url`, `Api` — never `ID`, `URL`, `API` |

---

## Boolean Column Rules

Boolean columns follow the same principles as the [cross-language boolean conventions](../01-cross-language/02-boolean-principles/00-overview.md) and [no-negatives rule](../01-cross-language/12-no-negatives.md).

### Rule 1: Always Use `Is` or `Has` Prefix

Every boolean column MUST start with `Is` or `Has`:

```sql
-- ✅ CORRECT
IsActive      BOOLEAN NOT NULL DEFAULT 1
IsVerified    BOOLEAN NOT NULL DEFAULT 0
IsPublished   BOOLEAN NOT NULL DEFAULT 0
HasLicense    BOOLEAN NOT NULL DEFAULT 0
HasChildren   BOOLEAN NOT NULL DEFAULT 0
IsEnabled     BOOLEAN NOT NULL DEFAULT 1

-- ❌ WRONG — no prefix
Active        BOOLEAN
Verified      BOOLEAN
Published     BOOLEAN
Licensed      BOOLEAN
```

### Rule 2: Never Use Negative Boolean Column Names

Boolean columns MUST express the **positive** state. Never name a column as a negation — it causes double-negative confusion in queries and code.

```sql
-- ❌ FORBIDDEN — negative names
IsNotActive       BOOLEAN    -- double negative: WHERE IsNotActive = 0
IsDisabled        BOOLEAN    -- negative: WHERE IsDisabled = 0 means "enabled"
IsInvalid         BOOLEAN    -- negative: WHERE IsInvalid = 0 means "valid"
HasNoLicense      BOOLEAN    -- negative: WHERE HasNoLicense = 0 means "has license"
IsUnverified      BOOLEAN    -- negative prefix "Un"
IsNotPublished    BOOLEAN    -- explicit negation

-- ✅ CORRECT — always positive
IsActive          BOOLEAN    -- WHERE IsActive = 1 (clear intent)
IsEnabled         BOOLEAN    -- WHERE IsEnabled = 1 (clear intent)
IsValid           BOOLEAN    -- WHERE IsValid = 1 (clear intent)
HasLicense        BOOLEAN    -- WHERE HasLicense = 1 (clear intent)
IsVerified        BOOLEAN    -- WHERE IsVerified = 1 (clear intent)
IsPublished       BOOLEAN    -- WHERE IsPublished = 1 (clear intent)
```

### Rule 3: Query Readability Test

A well-named boolean column reads naturally in both true and false checks:

| Column | True Check | False Check | Reads Naturally? |
|--------|-----------|-------------|-----------------|
| `IsActive` | `WHERE IsActive = 1` → "is active" | `WHERE IsActive = 0` → "is not active" | ✅ Yes |
| `HasLicense` | `WHERE HasLicense = 1` → "has license" | `WHERE HasLicense = 0` → "has no license" | ✅ Yes |
| `IsDisabled` | `WHERE IsDisabled = 1` → "is disabled" | `WHERE IsDisabled = 0` → "is not disabled" (??) | ❌ Confusing |
| `IsNotActive` | `WHERE IsNotActive = 1` → "is not active" | `WHERE IsNotActive = 0` → "is not not active" (??) | ❌ Double negative |

### Rule 4: ORM Mapping

Boolean columns with `Is`/`Has` prefix map cleanly to code:

```go
// Go struct — boolean fields match column names
type User struct {
    UserId     int64 `db:"UserId"`
    IsActive   bool  `db:"IsActive"`
    IsVerified bool  `db:"IsVerified"`
    HasLicense bool  `db:"HasLicense"`
}

// Clean, readable business logic
if user.IsActive && user.HasLicense {
    // grant access
}
```

```php
// PHP — clean property access
if ($user->IsActive && $user->HasLicense) {
    // grant access
}
```

```typescript
// TypeScript — clean conditionals
if (user.IsActive && user.HasLicense) {
    // grant access
}
```

### Rule 5: Only `Is` and `Has` — No Other Prefixes

Only two prefixes are allowed. Do NOT use `Can`, `Should`, `Was`, `Will`, `Did`, or bare adjectives:

```sql
-- ❌ WRONG — other prefixes
CanEdit       BOOLEAN    -- use IsEditable
ShouldSync    BOOLEAN    -- use IsSyncRequired
WasProcessed  BOOLEAN    -- use IsProcessed
WillExpire    BOOLEAN    -- use IsExpiring or use ExpiresAt timestamp
DidComplete   BOOLEAN    -- use IsComplete

-- ✅ CORRECT
IsEditable       BOOLEAN NOT NULL DEFAULT 0
IsSyncRequired   BOOLEAN NOT NULL DEFAULT 0
IsProcessed      BOOLEAN NOT NULL DEFAULT 0
IsComplete       BOOLEAN NOT NULL DEFAULT 0
```

### Rule 6: Always `NOT NULL DEFAULT`

Boolean columns MUST never be nullable. A three-state boolean (`true`/`false`/`NULL`) is a logic bug waiting to happen.

```sql
-- ❌ WRONG — nullable boolean
IsActive BOOLEAN    -- NULL = unknown state, breaks WHERE IsActive = 0

-- ✅ CORRECT — always NOT NULL with explicit default
IsActive BOOLEAN NOT NULL DEFAULT 1
IsVerified BOOLEAN NOT NULL DEFAULT 0
```

### Rule 7: Prefer Timestamp Over Boolean When Applicable

If you need to know **when** something happened (not just whether it happened), use a nullable timestamp column instead of a boolean:

| Boolean | Timestamp Alternative | When to Use Timestamp |
|---------|----------------------|----------------------|
| `IsDeleted` | `DeletedAt TEXT NULL` | Soft deletes — need to know when |
| `IsExpired` | `ExpiresAt TEXT NULL` | Expiration tracking |
| `IsBanned` | `BannedAt TEXT NULL` | Audit trail needed |
| `IsCompleted` | `CompletedAt TEXT NULL` | Duration tracking |

```sql
-- Boolean is fine when you only need true/false
IsActive   BOOLEAN NOT NULL DEFAULT 1

-- Timestamp is better when "when" matters
DeletedAt  TEXT NULL    -- NULL = not deleted, non-NULL = deleted at this time
ExpiresAt  TEXT NULL    -- NULL = never expires, non-NULL = expires at this time
```

> **Rule of thumb:** If the business logic ever asks "when did this happen?", use a timestamp. If it only asks "is this the case?", use a boolean.

### Rule 8: Common Negative-to-Positive Conversion

Reference table for converting common negative column names:

| ❌ Negative Name | ✅ Positive Replacement | Notes |
|-----------------|----------------------|-------|
| `IsNotActive` | `IsActive` | Flip the default |
| `IsDisabled` | `IsEnabled` | Positive state |
| `IsInvalid` | `IsValid` | Positive state |
| `IsIncomplete` | `IsComplete` | Positive state |
| `IsUnavailable` | `IsAvailable` | Positive state |
| `IsUnread` | `IsRead` | Positive state, default 0 |
| `IsHidden` | `IsVisible` | Positive state |
| `IsBroken` | `IsWorking` | Positive state |
| `HasNoAccess` | `HasAccess` | Positive state |
| `HasNoChildren` | `HasChildren` | Positive state |
| `IsUnpublished` | `IsPublished` | Positive state, default 0 |
| `IsLocked` | `IsEditable` | Reframe positively |

### Complete Example

```sql
CREATE TABLE Users (
    UserId        INTEGER PRIMARY KEY AUTOINCREMENT,
    Name          TEXT NOT NULL,
    Email         TEXT NOT NULL UNIQUE,
    IsActive      BOOLEAN NOT NULL DEFAULT 1,
    IsVerified    BOOLEAN NOT NULL DEFAULT 0,
    HasLicense    BOOLEAN NOT NULL DEFAULT 0,
    IsAdmin       BOOLEAN NOT NULL DEFAULT 0,
    CreatedAt     TEXT NOT NULL DEFAULT (datetime('now')),
    DeletedAt     TEXT NULL                                -- soft delete (timestamp > boolean)
);

-- Clean queries
SELECT * FROM Users WHERE IsActive = 1 AND IsVerified = 1;
SELECT * FROM Users WHERE HasLicense = 1 AND IsAdmin = 0;
SELECT * FROM Users WHERE DeletedAt IS NULL;              -- not deleted
```

---

## Primary Key Naming Pattern

The primary key MUST be named `{TableName}Id` — not just `Id`:

```sql
-- ❌ WRONG — generic Id
CREATE TABLE Transactions (
    Id INTEGER PRIMARY KEY AUTOINCREMENT
);

-- ✅ CORRECT — TableNameId
CREATE TABLE Transactions (
    TransactionId INTEGER PRIMARY KEY AUTOINCREMENT
);
```

**Why:** When this column appears as a foreign key in another table, the name is self-documenting:

```sql
CREATE TABLE TransactionLogs (
    TransactionLogId INTEGER PRIMARY KEY AUTOINCREMENT,
    TransactionId    INTEGER NOT NULL,  -- clearly references Transactions
    LogMessage       TEXT,
    FOREIGN KEY (TransactionId) REFERENCES Transactions(TransactionId)
);
```

---

## Foreign Key Column Naming

Foreign key columns MUST use the **exact same name** as the primary key they reference:

```sql
-- Source table
CREATE TABLE AgentSites (
    AgentSiteId INTEGER PRIMARY KEY AUTOINCREMENT,
    SiteName    TEXT NOT NULL
);

-- Referencing table — FK column matches PK name exactly
CREATE TABLE Transactions (
    TransactionId INTEGER PRIMARY KEY AUTOINCREMENT,
    AgentSiteId   INTEGER NOT NULL,  -- same name as AgentSites.AgentSiteId
    Amount        REAL,
    FOREIGN KEY (AgentSiteId) REFERENCES AgentSites(AgentSiteId)
);
```

---

## WordPress Exception

WordPress core tables (`wp_posts`, `wp_options`) retain their native `snake_case` naming. Only custom tables follow PascalCase.

See full details: [../01-cross-language/07-database-naming.md](../01-cross-language/07-database-naming.md)

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Full naming spec | [../01-cross-language/07-database-naming.md](../01-cross-language/07-database-naming.md) |
| Key naming PascalCase | [../01-cross-language/11-key-naming-pascalcase.md](../01-cross-language/11-key-naming-pascalcase.md) |
| Boolean principles | [../01-cross-language/02-boolean-principles/00-overview.md](../01-cross-language/02-boolean-principles/00-overview.md) |
| No-negatives rule | [../01-cross-language/12-no-negatives.md](../01-cross-language/12-no-negatives.md) |
