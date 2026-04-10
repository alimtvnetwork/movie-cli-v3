# Database Relationship Diagrams

**Version:** 1.1.0  
**Updated:** 2026-04-02

---

## Overview

Visual relationship patterns for AI-readable schema design. These diagrams show how tables relate to each other following the project's database conventions.

> **AI note:** When implementing a schema, use Section 7 (AI Implementation Checklist) as a step-by-step guide. Section 6 provides a copy-paste-ready complete SQL example.

---

## 1. Naming Convention Quick Reference

```
RULE                    CONVENTION                 EXAMPLE
─────────────────────────────────────────────────────────────
Table name              PascalCase                 AgentSites
Column name             PascalCase                 PluginSlug
Primary key             {TableName}Id              TransactionId
Foreign key             Same name as referenced PK AgentSiteId
Boolean column          Is{Positive} / Has{Noun}   IsActive, HasLicense
Index                   Idx{Table}_{Column}        IdxTransactions_CreatedAt
View                    Vw{DescriptiveName}        VwTransactionDetails
Abbreviations           First letter only caps     Id, Url, Api (NOT ID, URL)
```

---

## 2. Primary Key and Foreign Key Pattern

Shows how `{TableName}Id` creates self-documenting FK relationships.

**Reading the diagram:** Arrows show FK → PK direction (the FK column *references* the PK).

```
┌───────────────────────────┐
│       AgentSites          │
├───────────────────────────┤
│ AgentSiteId   INTEGER (PK)│
│ SiteName      TEXT        │
│ SiteUrl       TEXT        │
│ IsActive      BOOLEAN     │
│ CreatedAt     TEXT        │
└─────────────┬─────────────┘
              │ referenced by
              │
              ▼
┌──────────────────────────────────────────┐
│              Transactions                │
├──────────────────────────────────────────┤
│ TransactionId  INTEGER (PK)              │
│ AgentSiteId    INTEGER (FK) ─────────────│──→ AgentSites.AgentSiteId
│ StatusTypeId   SMALLINT (FK) ────────────│──→ StatusTypes.StatusTypeId
│ FileTypeId     SMALLINT (FK) ────────────│──→ FileTypes.FileTypeId
│ PluginSlug     TEXT                      │
│ Amount         REAL                      │
│ IsActive       BOOLEAN                   │
│ CreatedAt      TEXT                       │
└──────────────────────────────────────────┘
              │                │
   ┌──────────┘                └──────────┐
   ▼                                      ▼
┌─────────────────────┐    ┌─────────────────────┐
│    StatusTypes       │    │     FileTypes        │
├─────────────────────┤    ├─────────────────────┤
│ StatusTypeId (PK)   │    │ FileTypeId  (PK)    │
│ Name         TEXT   │    │ Name        TEXT     │
├─────────────────────┤    ├─────────────────────┤
│ 1 = Pending         │    │ 1 = Plugin           │
│ 2 = Complete        │    │ 2 = Theme            │
│ 3 = Failed          │    │ 3 = MuPlugin         │
└─────────────────────┘    └─────────────────────┘
```

**Key rules demonstrated:**
- FK column `AgentSiteId` in `Transactions` has the **exact same name** as PK `AgentSiteId` in `AgentSites`
- Lookup tables (`StatusTypes`, `FileTypes`) use SMALLINT PK — they will never exceed 32K rows
- Main table stores only the integer FK — the human-readable `Name` lives in the lookup table
- Boolean columns use `Is` prefix (`IsActive`) — never negative names

---

## 3. Many-to-Many (N-to-M) Pattern

The junction table holds FKs to both sides. FK columns use the **exact PK names** from their source tables.

```
┌─────────────────────┐                          ┌─────────────────────┐
│       Users         │                          │       Roles         │
├─────────────────────┤                          ├─────────────────────┤
│ UserId     (PK)     │                          │ RoleId     (PK)     │
│ Name       TEXT     │    ┌──────────────────┐  │ Name       TEXT     │
│ Email      TEXT     │    │    UserRoles     │  ├─────────────────────┤
│ IsActive   BOOLEAN  │    ├──────────────────┤  │ 1 = Admin           │
│ IsVerified BOOLEAN  │    │ UserRoleId (PK)  │  │ 2 = Editor          │
│ HasLicense BOOLEAN  │    │ UserId     (FK)──│──│ 3 = Viewer          │
└─────────────────────┘    │ RoleId     (FK)──│──┘                     
         │                 │ UNIQUE(UserId,   │                        
         └────────────────>│        RoleId)   │                        
            UserId FK      └──────────────────┘                        
```

**Junction table rules:**
- PK: `{JunctionTable}Id` → `UserRoleId`
- FK columns: exact same names as source PKs → `UserId`, `RoleId`
- `UNIQUE(UserId, RoleId)` prevents duplicate assignments
- No extra columns unless the relationship itself has attributes (e.g., `AssignedAt`)

---

## 4. View Pattern — Flattening Joins

Views pre-define JOINs so the business layer queries a flat result via ORM.

```
SOURCE TABLES                              VIEW (flat output)
─────────────────────                      ──────────────────────────────────
Transactions.TransactionId ──────────────→ TransactionId
Transactions.PluginSlug    ──────────────→ PluginSlug
Transactions.Amount        ──────────────→ Amount
Transactions.IsActive      ──────────────→ IsActive
Transactions.CreatedAt     ──────────────→ CreatedAt
                                           
StatusTypes.Name           ─── JOIN on ──→ StatusName
                           StatusTypeId    
                                           
FileTypes.Name             ─── JOIN on ──→ FileTypeName
                           FileTypeId      
                                           
AgentSites.SiteName        ─── JOIN on ──→ AgentSiteName
AgentSites.SiteUrl         AgentSiteId  ─→ AgentSiteUrl

═══════════════════════════════════════════════════════════════
RESULT: VwTransactionDetails

Business layer code:
  orm.FindAll("VwTransactionDetails", filter)

API response (PascalCase JSON):
  {"TransactionId": 42, "StatusName": "Pending", "IsActive": true, ...}

No JOINs in application code. No key transformation.
═══════════════════════════════════════════════════════════════
```

> See [06-rest-api-format.md](./06-rest-api-format.md) for the full REST API response format.

---

## 5. Key Sizing Decision Tree

```
How many rows will this table have in 10 years?
│
├─ < 32,000 rows
│  └─ SMALLINT (2 bytes)
│     Use for: lookup tables (StatusTypes, FileTypes, Roles)
│
├─ < 2,000,000,000 rows
│  └─ INTEGER (4 bytes) ← THIS IS THE DEFAULT
│     Use for: most entity tables (Users, Transactions, Logs)
│
├─ > 2,000,000,000 rows
│  └─ BIGINT (8 bytes)
│     Use for: event streams, analytics, high-volume logs
│
└─ Need UUID? Check ALL 3 criteria:
   [1] Created across multiple disconnected systems?
   [2] No central ID authority?
   [3] Must be publicly exposed and non-guessable?
   │
   ├─ All 3 = YES → BLOB(16) (not TEXT(36))
   └─ Any = NO   → Use INTEGER instead
```

---

## 6. Complete Schema Example

Copy-paste-ready SQL showing all conventions together:

```sql
-- ============================================================
-- LOOKUP TABLES (SMALLINT PKs — under 32K rows expected)
-- ============================================================

CREATE TABLE StatusTypes (
    StatusTypeId SMALLINT PRIMARY KEY,
    Name         TEXT NOT NULL UNIQUE
);

CREATE TABLE FileTypes (
    FileTypeId SMALLINT PRIMARY KEY,
    Name       TEXT NOT NULL UNIQUE
);

CREATE TABLE Roles (
    RoleId SMALLINT PRIMARY KEY,
    Name   TEXT NOT NULL UNIQUE
);

-- ============================================================
-- ENTITY TABLES (INTEGER PKs — under 2B rows expected)
-- ============================================================

CREATE TABLE AgentSites (
    AgentSiteId INTEGER PRIMARY KEY AUTOINCREMENT,
    SiteName    TEXT NOT NULL,
    SiteUrl     TEXT NOT NULL,
    IsActive    BOOLEAN NOT NULL DEFAULT 1,
    CreatedAt   TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE Users (
    UserId     INTEGER PRIMARY KEY AUTOINCREMENT,
    Name       TEXT NOT NULL,
    Email      TEXT NOT NULL UNIQUE,
    IsActive   BOOLEAN NOT NULL DEFAULT 1,
    IsVerified BOOLEAN NOT NULL DEFAULT 0,
    HasLicense BOOLEAN NOT NULL DEFAULT 0,
    CreatedAt  TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE Transactions (
    TransactionId INTEGER PRIMARY KEY AUTOINCREMENT,
    AgentSiteId   INTEGER NOT NULL,
    StatusTypeId  SMALLINT NOT NULL,
    FileTypeId    SMALLINT NOT NULL,
    PluginSlug    TEXT NOT NULL,
    Amount        REAL NOT NULL DEFAULT 0,
    IsActive      BOOLEAN NOT NULL DEFAULT 1,
    CreatedAt     TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (AgentSiteId)  REFERENCES AgentSites(AgentSiteId),
    FOREIGN KEY (StatusTypeId) REFERENCES StatusTypes(StatusTypeId),
    FOREIGN KEY (FileTypeId)   REFERENCES FileTypes(FileTypeId)
);

-- ============================================================
-- JUNCTION TABLE (N-to-M relationship)
-- ============================================================

CREATE TABLE UserRoles (
    UserRoleId INTEGER PRIMARY KEY AUTOINCREMENT,
    UserId     INTEGER NOT NULL,
    RoleId     SMALLINT NOT NULL,
    UNIQUE (UserId, RoleId),
    FOREIGN KEY (UserId) REFERENCES Users(UserId),
    FOREIGN KEY (RoleId) REFERENCES Roles(RoleId)
);

-- ============================================================
-- INDEXES (PascalCase: Idx{Table}_{Column})
-- ============================================================

CREATE INDEX IdxTransactions_CreatedAt    ON Transactions(CreatedAt);
CREATE INDEX IdxTransactions_PluginSlug   ON Transactions(PluginSlug);
CREATE INDEX IdxTransactions_StatusTypeId ON Transactions(StatusTypeId);
CREATE INDEX IdxUserRoles_UserId          ON UserRoles(UserId);

-- ============================================================
-- VIEWS (Vw prefix — flatten joins for ORM queries)
-- ============================================================

CREATE VIEW VwTransactionDetails AS
SELECT
    t.TransactionId,
    t.PluginSlug,
    t.Amount,
    t.IsActive,
    t.CreatedAt,
    st.Name       AS StatusName,
    ft.Name       AS FileTypeName,
    a.SiteName    AS AgentSiteName,
    a.SiteUrl     AS AgentSiteUrl
FROM Transactions t
INNER JOIN StatusTypes st ON t.StatusTypeId = st.StatusTypeId
INNER JOIN FileTypes ft   ON t.FileTypeId = ft.FileTypeId
LEFT JOIN AgentSites a    ON t.AgentSiteId = a.AgentSiteId;

CREATE VIEW VwUserRoleSummary AS
SELECT
    u.UserId,
    u.Name       AS UserName,
    u.Email,
    u.IsActive,
    u.IsVerified,
    u.HasLicense,
    r.Name       AS RoleName
FROM Users u
INNER JOIN UserRoles ur ON u.UserId = ur.UserId
INNER JOIN Roles r      ON ur.RoleId = r.RoleId;

-- ============================================================
-- SEED LOOKUP DATA
-- ============================================================

INSERT INTO StatusTypes (StatusTypeId, Name) VALUES
    (1, 'Pending'),
    (2, 'Complete'),
    (3, 'Failed');

INSERT INTO FileTypes (FileTypeId, Name) VALUES
    (1, 'Plugin'),
    (2, 'Theme'),
    (3, 'MuPlugin');

INSERT INTO Roles (RoleId, Name) VALUES
    (1, 'Admin'),
    (2, 'Editor'),
    (3, 'Viewer');
```

---

## 7. AI Implementation Checklist

When an AI is asked to create or modify a database schema, follow this checklist in order:

| # | Step | Rule |
|---|------|------|
| 1 | **Name everything PascalCase** | Tables, columns, indexes, views |
| 2 | **Primary key = `{TableName}Id`** | Never bare `Id` |
| 3 | **FK column = exact PK name** | `AgentSiteId` in both tables |
| 4 | **Boolean = `Is`/`Has` + positive** | `IsActive`, `HasLicense` — never `IsDisabled`, `IsNotActive` |
| 5 | **Size the key type** | Estimate 10-year rows → SMALLINT / INTEGER / BIGINT |
| 6 | **No UUID** | Unless distributed + public + non-guessable (all 3) |
| 7 | **Extract repeated values** | Create lookup table + FK relationship |
| 8 | **Create views for joins** | `Vw` prefix — business layer queries views only |
| 9 | **Use ORM** | No raw SQL in service/business layer |
| 10 | **Unit test schemas** | In-memory SQLite — verify tables, columns, constraints |
| 11 | **Integration test CRUD** | In-memory SQLite — verify insert, select, update, delete |
| 12 | **REST API = PascalCase JSON** | Response keys match DB column names exactly |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Naming conventions | [./01-naming-conventions.md](./01-naming-conventions.md) |
| Schema design | [./02-schema-design.md](./02-schema-design.md) |
| ORM and views | [./03-orm-and-views.md](./03-orm-and-views.md) |
| Testing strategy | [./04-testing-strategy.md](./04-testing-strategy.md) |
| REST API format | [./06-rest-api-format.md](./06-rest-api-format.md) |
| Boolean principles | [../01-cross-language/02-boolean-principles/00-overview.md](../01-cross-language/02-boolean-principles/00-overview.md) |
| No-negatives rule | [../01-cross-language/12-no-negatives.md](../01-cross-language/12-no-negatives.md) |
