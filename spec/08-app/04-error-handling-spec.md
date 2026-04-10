# Application Error Handling Specification

**Version:** 1.0.0  
**Created:** 09-Apr-2026  
**Scope:** Runtime error handling for `movie` CLI — TMDb API, SQLite, filesystem, and network failures

---

## 1. Overview

This spec defines how the `movie` CLI handles three categories of runtime errors:

1. **TMDb API errors** — rate limits, auth failures, server errors
2. **SQLite database errors** — locks, corruption, migration failures
3. **Offline / network errors** — no connectivity, DNS resolution, timeouts

### Design Principles

| Principle | Rule |
|-----------|------|
| Never silently fail | Every error is logged to stderr or returned |
| Degrade gracefully | Missing metadata ≠ command failure |
| Preserve user data | Never delete source files on partial operations |
| Be actionable | Error messages include what to do next |

---

## 2. TMDb API Errors

### 2.1 Error Classification

| HTTP Status | Meaning | Behavior |
|-------------|---------|----------|
| `200` | Success | Process response normally |
| `401` | Invalid API key | Print error + hint: `movie config set tmdb_api_key YOUR_KEY` |
| `404` | Resource not found | Treat as "no result" — do not retry |
| `422` | Invalid parameters | Log warning, skip item |
| `429` | Rate limit exceeded | Retry with backoff (see §2.2) |
| `500-503` | Server error | Retry once after 2s, then skip with warning |

### 2.2 Rate Limiting Strategy

TMDb enforces ~40 requests per 10 seconds. The CLI must handle `429 Too Many Requests`:

```
GIVEN the TMDb API returns HTTP 429
WHEN any command is making API calls
THEN the client waits for the duration specified in the `Retry-After` header
  (or 10 seconds if header is absent)
AND retries the request up to 3 times
AND if all retries fail, logs a warning and continues to the next item
```

**Implementation location:** `tmdb/client.go` → `get()` method

```go
// Pseudocode for retry logic
func (c *Client) getWithRetry(url string, target interface{}) error {
    for attempt := 0; attempt < 3; attempt++ {
        err := c.get(url, target)
        if err == nil {
            return nil
        }
        if isRateLimited(err) {
            wait := parseRetryAfter(err) // default 10s
            fmt.Fprintf(os.Stderr, "⚠️  Rate limited, waiting %ds...\n", wait)
            time.Sleep(time.Duration(wait) * time.Second)
            continue
        }
        return err // non-retryable error
    }
    return fmt.Errorf("TMDb rate limit: max retries exceeded")
}
```

### 2.3 Commands Affected

| Command | TMDb Calls | Failure Mode |
|---------|-----------|--------------|
| `scan` | SearchMulti + Details + Credits + Poster per file | Skip metadata for failed items, continue scan |
| `search` | SearchMulti + Details + Credits + Poster | Fatal — cannot proceed without results |
| `info` | SearchMulti + Details + Credits + Poster (fallback only) | Show local data if available, else error |
| `suggest` | Recommendations + Trending | Show partial results or "no suggestions" |

### 2.4 Acceptance Criteria

```
GIVEN a valid TMDb API key
WHEN the API returns 429 on a scan with 50 files
THEN the scan pauses, retries, and completes without crashing
AND skipped items are reported in the summary

GIVEN an invalid TMDb API key
WHEN any TMDb-dependent command runs
THEN the error message includes the exact config command to fix it

GIVEN the TMDb API returns 500
WHEN fetching movie details
THEN one retry is attempted after 2 seconds
AND if retry fails, the item is skipped with a warning
```

---

## 3. SQLite Database Errors

### 3.1 Error Classification

| Error | Cause | Behavior |
|-------|-------|----------|
| `SQLITE_BUSY` (5) | Another process holds the lock | Retry up to 5 times with 200ms intervals |
| `SQLITE_LOCKED` (6) | Table-level lock within same connection | Should not occur with WAL mode — log and abort |
| `SQLITE_CORRUPT` (11) | Database file damaged | Fatal error: print path, suggest backup restore |
| `SQLITE_CANTOPEN` (14) | File permissions or missing directory | Print path + hint: check permissions |
| `SQLITE_CONSTRAINT` (19) | Unique violation (e.g., duplicate tmdb_id) | Handled: fall through to `UpdateMediaByTmdbID` |
| `SQLITE_FULL` (13) | Disk full | Fatal error: print disk usage hint |

### 3.2 BUSY Retry Strategy

SQLite in WAL mode rarely encounters BUSY, but it can happen with concurrent CLI instances:

```
GIVEN the database returns SQLITE_BUSY
WHEN any write operation is attempted
THEN the operation is retried up to 5 times with 200ms delay between attempts
AND if all retries fail, an error is printed: "Database is locked. Another movie process may be running."
```

**Implementation location:** `db/db.go` — add `PRAGMA busy_timeout = 5000` after WAL mode

```go
// Add after WAL mode pragma
if _, err := conn.Exec("PRAGMA busy_timeout = 5000"); err != nil {
    conn.Close()
    return nil, fmt.Errorf("cannot set busy timeout: %w", err)
}
```

### 3.3 Database Recovery

```
GIVEN the database file is corrupted (SQLITE_CORRUPT)
WHEN any command attempts to open the database
THEN the error message includes:
  - The full path to the database file
  - Suggestion: "Delete ~/data/movie.db and re-scan your library"
  - Note: "Thumbnails in ~/data/thumbnails/ are preserved"
```

### 3.4 Migration Failure

```
GIVEN the migration SQL fails (e.g., schema change on existing DB)
WHEN db.Open() is called
THEN the error is wrapped with "migration failed" context
AND the full SQLite error is included for debugging
AND the command exits with code 1
```

### 3.5 Acceptance Criteria

```
GIVEN two movie processes run simultaneously
WHEN both attempt to write to the database
THEN neither crashes — one waits via busy_timeout and succeeds

GIVEN the database file does not exist
WHEN any command runs
THEN the database is created with full schema and default config

GIVEN a duplicate tmdb_id insert
WHEN scan or search saves a record
THEN the existing record is updated instead of failing
```

---

## 4. Offline / Network Errors

### 4.1 Error Classification

| Scenario | Detection | Behavior |
|----------|-----------|----------|
| No network | `net.OpError` / DNS failure | Warn + continue without TMDb |
| DNS resolution failure | `no such host` in error | Same as no network |
| Connection timeout | HTTP client 15s timeout | Warn per item, continue |
| TLS handshake failure | `tls: handshake failure` | Warn + suggest checking time/date |
| Proxy / firewall block | Connection refused | Warn + suggest checking network config |

### 4.2 Offline Mode Behavior

Commands should degrade gracefully when offline:

| Command | Online Behavior | Offline Behavior |
|---------|----------------|-----------------|
| `scan` | Fetch TMDb metadata + poster | Store filename-only records (title, year, type, path, size) |
| `search` | Query TMDb API | Fatal: "No network connection. Search requires TMDb." |
| `info` (local hit) | Show from DB | Works fully — no network needed |
| `info` (TMDb fallback) | Search + save | Show "not found locally" + "cannot reach TMDb" |
| `suggest` | Recommendations + trending | Fatal: "No network connection. Suggest requires TMDb." |
| `ls` | N/A (local only) | Works fully |
| `move` | N/A (local only) | Works fully |
| `rename` | N/A (local only) | Works fully |
| `undo` | N/A (local only) | Works fully |
| `play` | N/A (local only) | Works fully |
| `stats` | N/A (local only) | Works fully |
| `tag` | N/A (local only) | Works fully |
| `config` | N/A (local only) | Works fully |

### 4.3 Network Error Detection

**Implementation location:** `tmdb/client.go` → `get()` method

```go
func isNetworkError(err error) bool {
    var netErr *net.OpError
    if errors.As(err, &netErr) {
        return true
    }
    var dnsErr *net.DNSError
    if errors.As(err, &dnsErr) {
        return true
    }
    if os.IsTimeout(err) {
        return true
    }
    return strings.Contains(err.Error(), "no such host") ||
           strings.Contains(err.Error(), "connection refused")
}
```

### 4.4 Acceptance Criteria

```
GIVEN the device has no internet connection
WHEN movie scan ~/Downloads is run
THEN files are scanned and stored with filename-only metadata
AND a warning is printed: "Offline — metadata not fetched. Re-run scan when online to enrich."
AND the summary shows how many items lack metadata

GIVEN the device has no internet connection
WHEN movie search "Inception" is run
THEN the command exits with: "Cannot reach TMDb. Check your network connection."

GIVEN the device has no internet connection
WHEN movie ls is run
THEN the list displays normally — no network dependency

GIVEN TMDb is reachable but responds with timeout on 3 of 20 files
WHEN movie scan runs
THEN 17 files get full metadata, 3 get filename-only records
AND the summary reports: "3 items skipped due to network errors"
```

---

## 5. Error Message Format

All error messages must follow this format:

```
❌ <What failed>: <Technical detail>
   <What to do about it>
```

Examples:

```
❌ TMDb API error: 401 Unauthorized
   Set your API key: movie config set tmdb_api_key YOUR_KEY

❌ Database error: database is locked
   Another movie process may be running. Wait and retry.

❌ Cannot reach TMDb: no such host api.themoviedb.org
   Check your internet connection and try again.

❌ File not found: /home/user/Movies/Inception (2010).mkv
   The file may have been moved or deleted manually.
```

### Severity Levels

| Icon | Level | Behavior |
|------|-------|----------|
| `❌` | Fatal | Command cannot continue — exits with code 1 |
| `⚠️` | Warning | Operation degraded but command continues |
| `ℹ️` | Info | Non-critical note (e.g., "already in database") |

---

## 6. Summary Table

| Error Category | Detection | Retry | Degrade | Fatal |
|---------------|-----------|-------|---------|-------|
| TMDb 429 (rate limit) | HTTP status | 3x with backoff | Skip item | No |
| TMDb 401 (auth) | HTTP status | No | N/A | Yes (with hint) |
| TMDb 500 (server) | HTTP status | 1x after 2s | Skip item | No |
| SQLite BUSY | Error code 5 | busy_timeout 5s | N/A | After timeout |
| SQLite CORRUPT | Error code 11 | No | N/A | Yes (with path) |
| No network | net.OpError | No | Filename-only | For search/suggest |
| DNS failure | net.DNSError | No | Filename-only | For search/suggest |
| Timeout | 15s HTTP timeout | No | Skip item | No |

---

## 7. Implementation Checklist

- [ ] Add `PRAGMA busy_timeout = 5000` to `db/db.go` after WAL pragma
- [ ] Add retry logic to `tmdb/client.go` `get()` for HTTP 429 responses
- [ ] Add `isNetworkError()` helper to `tmdb/client.go`
- [ ] Update `scan` to report metadata-skip count in summary
- [ ] Ensure all error messages follow the format in §5

---

## Cross-References

- [Error Management Architecture](../02-error-manage-spec/00-overview.md)
- [Acceptance Criteria](./97-acceptance-criteria.md)
- [Dashboard Spec](./03-dashboard-spec.md)
