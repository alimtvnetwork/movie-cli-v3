# Rust Naming Conventions

**Version:** 1.0.0  
**Updated:** 2026-03-27

---

## Overview

Naming rules for Rust code, aligned with the Rust API Guidelines (RFC 430) and project-wide PascalCase serialization mandate.

---

## Identifier Casing

| Item | Convention | Example |
|------|-----------|---------|
| Types, Traits, Enums | `PascalCase` | `BrowserCollector`, `ActivityEvent` |
| Functions, Methods | `snake_case` | `get_active_window`, `start_daemon` |
| Local variables | `snake_case` | `session_id`, `dwell_seconds` |
| Constants | `SCREAMING_SNAKE_CASE` | `MAX_BUFFER_SIZE`, `DEFAULT_PORT` |
| Static variables | `SCREAMING_SNAKE_CASE` | `GLOBAL_CONFIG` |
| Modules, Crates | `snake_case` | `browser_tracking`, `os_integration` |
| Enum variants | `PascalCase` | `ScreenshotTrigger::TabChange` |
| Trait methods | `snake_case` | `fn is_supported(&self) -> bool` |
| Type parameters | Single uppercase or `PascalCase` | `T`, `EventType` |
| Lifetimes | Short lowercase | `'a`, `'ctx` |
| Feature flags | `kebab-case` | `wayland-support`, `browser-extension` |

---

## Meaningful Identifiers

> **Inherited rule:** All identifiers must use full, descriptive names. The abbreviation `ctx` is prohibited — use `context` instead.

| ❌ Forbidden | ✅ Required |
|-------------|------------|
| `ctx` | `context` |
| `cfg` | `config` |
| `mgr` | `manager` |
| `btn` | `button` |
| `evt` | `event` |
| `msg` | `message` |
| `req` | `request` |
| `res` | `response` |
| `cb` | `callback` |
| `idx` | `index` |

**Exception:** Single-letter variables in closures and iterators are acceptable when the scope is ≤ 3 lines:

```rust
// ✅ Acceptable — short closure
let total: f64 = activities.iter().map(|a| a.dwell_seconds).sum();

// ❌ Forbidden — longer closure needs descriptive name
let results: Vec<_> = activities.iter().filter(|activity| {
    activity.dwell_seconds > min_threshold
        && activity.category == UrlCategory::Work
}).collect();
```

---

## Abbreviation Casing

Abbreviations are treated as regular words in PascalCase — only capitalize the first letter:

| ❌ Forbidden | ✅ Required |
|-------------|------------|
| `URLParser` | `UrlParser` |
| `HTTPClient` | `HttpClient` |
| `getJSON` | `get_json` (function) |
| `SQLiteDB` | `SqliteDb` |
| `APIServer` | `ApiServer` |
| `FFIBridge` | `FfiBridge` |
| `WALMode` | `WalMode` |
| `PIDFile` | `PidFile` |

---

## Serialization: PascalCase Wire Format

All JSON serialization uses PascalCase to match the project-wide standard:

```rust
// ✅ Correct — derive with rename_all
#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "PascalCase")]
pub struct BrowserActivity {
    pub id: Uuid,
    pub session_id: Uuid,
    pub url: Option<String>,
    pub title: String,
    pub dwell_seconds: f64,
    pub started_at: DateTime<Utc>,
}
// Serializes to: { "Id": "...", "SessionId": "...", "Url": "...", ... }
```

```rust
// ❌ Forbidden — default serde (produces snake_case JSON)
#[derive(Serialize)]
pub struct BrowserActivity {
    pub id: Uuid,           // Would serialize as "id" — wrong
    pub session_id: Uuid,   // Would serialize as "session_id" — wrong
}
```

### Enum Serialization

```rust
// ✅ PascalCase variants serialize correctly by default
#[derive(Debug, Serialize, Deserialize)]
pub enum ScreenshotTrigger {
    Periodic,
    TabChange,
    AppSwitch,
    Idle,
    Manual,
}
// Serializes to: "TabChange" — correct PascalCase
```

---

## Module Structure

```
src/
├── main.rs                  # Entry point, CLI parsing
├── daemon.rs                # Daemon lifecycle
├── config.rs                # Configuration loading
├── event.rs                 # ActivityEvent enum, EventSender type
├── storage/
│   ├── mod.rs               # Storage engine trait
│   ├── sqlite.rs            # SQLite implementation
│   └── migrations/          # Embedded SQL migrations
├── collectors/
│   ├── mod.rs               # Collector trait, registry
│   ├── browser.rs           # BrowserCollector
│   ├── app_focus.rs         # AppFocusCollector
│   ├── click.rs             # ClickCollector
│   ├── screenshot.rs        # ScreenshotCollector
│   └── idle.rs              # IdleCollector
├── platform/
│   ├── mod.rs               # Platform abstraction layer
│   ├── windows.rs           # #[cfg(target_os = "windows")]
│   ├── linux.rs             # #[cfg(target_os = "linux")]
│   └── macos.rs             # #[cfg(target_os = "macos")]
├── api/
│   ├── mod.rs               # HTTP server setup
│   ├── routes.rs            # Route definitions
│   └── handlers.rs          # Request handlers
└── models/
    ├── mod.rs               # Re-exports
    ├── activity.rs          # BrowserActivity, AppActivity, etc.
    ├── session.rs           # Session model
    └── screenshot.rs        # Screenshot model
```

### Naming Rules

- One type per file when the type is complex (> 50 lines)
- Module files use `snake_case.rs`
- Re-export public items from `mod.rs` for clean import paths
- Group related types in a single file when each is < 20 lines

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Cross-Language Function Naming | `../01-cross-language/10-function-naming.md` |
| Abbreviation Casing | `../01-cross-language/04-code-style/00-overview.md` |
| Database Naming | `../01-cross-language/07-database-naming.md` |
| Cross-Language Guidelines | `../01-cross-language/00-overview.md` |
| Boolean Flag Method Splitting | `../01-cross-language/24-boolean-flag-methods.md` |
