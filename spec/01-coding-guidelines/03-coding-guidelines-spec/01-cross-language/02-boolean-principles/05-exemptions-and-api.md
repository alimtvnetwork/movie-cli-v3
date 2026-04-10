# Boolean Principles — Static factory exemption, Result wrapper API

> **Parent:** [Boolean Principles](./00-overview.md)  
> **Version:** 2.6.0  
> **Updated:** 2026-03-31

---

## Static Factory Constructor Exemption

Methods like `DbResult::empty()`, `DbResultSet::empty()`, and `ResultSlice::empty()` are **static factory constructors** — they create a new empty instance, not query boolean state. These are **exempt** from the `is`/`has` prefix requirement (P1).

Boolean query methods on the **same classes** — such as `isEmpty()`, `isDefined()`, `hasError()`, `isSafe()`, `hasItems()` — **do** follow P1 correctly and must retain their prefixes.

| Method | Type | P1 Applies? |
|--------|------|-------------|
| `DbResult::empty()` | Static factory constructor | ❌ Exempt |
| `DbResultSet::empty()` | Static factory constructor | ❌ Exempt |
| `$result->isEmpty()` | Boolean query | ✅ Yes |
| `$result->hasError()` | Boolean query | ✅ Yes |
| `$result->isDefined()` | Boolean query | ✅ Yes |
| `result.IsSafe()` | Boolean query | ✅ Yes |

---


---

## Result Wrapper — Full Public API Reference

> **Cross-language invariant:** The `.AppError()` (Go) / `.error()` (PHP) method on every result wrapper returns the **framework's structured error type**, never a raw string or generic exception. In Go this is `*apperror.AppError` (carrying stack trace, error code, and contextual values) — named `.AppError()` (not `.Error()`) to avoid confusion with Go's native `error` interface. In PHP this is `Throwable` (typically a framework exception with trace). This guarantees that propagated errors always preserve diagnostic context — callers can safely pass `.AppError()` output to `Fail()`, `FailSlice()`, `FailMap()`, or log it with full traceability.

### Go — `apperror.Result[T]`

| Method | Returns | Description |
|--------|---------|-------------|
| `Ok[T](value)` | `Result[T]` | Static: successful result with value |
| `Fail[T](err)` | `Result[T]` | Static: failed result from `*AppError` |
| `FailWrap[T](cause, code, msg)` | `Result[T]` | Static: failed result wrapping raw error |
| `FailNew[T](code, msg)` | `Result[T]` | Static: failed result from new error |
| `HasError()` | `bool` | True when the operation failed |
| `IsSafe()` | `bool` | True when a value exists and no error |
| `IsDefined()` | `bool` | True when a value was set (regardless of error) |
| `IsEmpty()` | `bool` | True when no value was set (absent, not an error) |
| `Value()` | `T` | Returns value; **panics** if `HasError()` is true |
| `ValueOr(fallback)` | `T` | Returns value if defined, otherwise fallback |
| `AppError()` | `*AppError` | Returns underlying error, or nil. Named `AppError()` to avoid confusion with Go's `error` interface |
| `Unwrap()` | `(T, error)` | Bridges to standard Go `(T, error)` pattern |

### Go — `apperror.ResultSlice[T]`

| Method | Returns | Description |
|--------|---------|-------------|
| `OkSlice[T](items)` | `ResultSlice[T]` | Static: successful slice result |
| `FailSlice[T](err)` | `ResultSlice[T]` | Static: failed slice from `*AppError` |
| `FailSliceWrap[T](cause, code, msg)` | `ResultSlice[T]` | Static: failed slice wrapping raw error |
| `FailSliceNew[T](code, msg)` | `ResultSlice[T]` | Static: failed slice from new error |
| `HasError()` | `bool` | True when the operation failed |
| `IsSafe()` | `bool` | True when no error (items may be empty) |
| `HasItems()` | `bool` | True when slice has ≥1 item |
| `IsEmpty()` | `bool` | True when slice has zero items |
| `Count()` | `int` | Number of items |
| `Items()` | `[]T` | Returns underlying slice (nil if error) |
| `First()` | `Result[T]` | First item as `Result[T]`, or empty |
| `Last()` | `Result[T]` | Last item as `Result[T]`, or empty |
| `GetAt(index)` | `Result[T]` | Item at index as `Result[T]`, or empty |
| `Append(items...)` | — | Adds items; no-op if in error state |
| `AppError()` | `*AppError` | Returns underlying error, or nil |

### Go — `apperror.ResultMap[K, V]`

| Method | Returns | Description |
|--------|---------|-------------|
| `OkMap[K,V](items)` | `ResultMap[K,V]` | Static: successful map result |
| `FailMap[K,V](err)` | `ResultMap[K,V]` | Static: failed map from `*AppError` |
| `FailMapWrap[K,V](cause, code, msg)` | `ResultMap[K,V]` | Static: failed map wrapping raw error |
| `FailMapNew[K,V](code, msg)` | `ResultMap[K,V]` | Static: failed map from new error |
| `HasError()` | `bool` | True when the operation failed |
| `IsSafe()` | `bool` | True when no error (map may be empty) |
| `HasItems()` | `bool` | True when map has ≥1 entry |
| `IsEmpty()` | `bool` | True when map has zero entries |
| `Count()` | `int` | Number of entries |
| `Items()` | `map[K]V` | Returns underlying map (nil if error) |
| `Get(key)` | `Result[V]` | Value for key as `Result[V]`, or empty |
| `Has(key)` | `bool` | True if key exists |
| `Set(key, value)` | — | Adds/updates entry; no-op if error |
| `Remove(key)` | — | Deletes key; no-op if error |
| `Keys()` | `[]K` | All keys as slice |
| `Values()` | `[]V` | All values as slice |
| `AppError()` | `*AppError` | Returns underlying error, or nil |

### PHP — `DbResult<T>`

| Method | Returns | Description |
|--------|---------|-------------|
| `DbResult::of($value)` | `DbResult<T>` | Static: successful result with value |
| `DbResult::empty()` | `DbResult<T>` | Static: empty result (no row found) |
| `DbResult::error($e)` | `DbResult<T>` | Static: error result with stack trace |
| `isEmpty()` | `bool` | True when no row was found (not an error) |
| `isDefined()` | `bool` | True when a row was successfully mapped |
| `hasError()` | `bool` | True when the query failed |
| `isSafe()` | `bool` | True when value exists and no error |
| `value()` | `T\|null` | Returns mapped value (null if not defined) |
| `error()` | `Throwable\|null` | Returns underlying error, or null |
| `stackTrace()` | `string` | Captured stack trace if error occurred |

### PHP — `DbResultSet<T>`

| Method | Returns | Description |
|--------|---------|-------------|
| `DbResultSet::of($items)` | `DbResultSet<T>` | Static: successful result set |
| `DbResultSet::error($e)` | `DbResultSet<T>` | Static: error result with stack trace |
| `isEmpty()` | `bool` | True when zero items |
| `hasAny()` | `bool` | True when ≥1 item |
| `count()` | `int` | Number of items |
| `hasError()` | `bool` | True when the query failed |
| `isSafe()` | `bool` | True when no error (items may be empty) |
| `items()` | `array<T>` | Returns item array |
| `first()` | `DbResult<T>` | First item as `DbResult<T>`, or error/empty |
| `error()` | `Throwable\|null` | Returns underlying error, or null |
| `stackTrace()` | `string` | Captured stack trace if error occurred |

### PHP — `DbExecResult`

| Method | Returns | Description |
|--------|---------|-------------|
| `DbExecResult::of($rows, $id)` | `DbExecResult` | Static: successful exec result |
| `DbExecResult::error($e)` | `DbExecResult` | Static: error result with stack trace |
| `hasError()` | `bool` | True when the exec failed |
| `isSafe()` | `bool` | True when no error |
| `isEmpty()` | `bool` | True when zero rows affected |
| `affectedRows()` | `int` | Number of affected rows |
| `lastInsertId()` | `int` | Auto-increment ID from INSERT |
| `error()` | `Throwable\|null` | Returns underlying error, or null |
| `stackTrace()` | `string` | Captured stack trace if error occurred |

---

