# Comments, Documentation & Dead Code

> **Version:** 4.0.0  
> **Updated:** 2026-03-31  
> **Applies to:** PHP, TypeScript, Go  
> **Rules covered:** 8, 14, 15, 16

---

## Rule 8: No Leading Backslash on Global Types

In catch blocks and type hints, use `Throwable` without the leading backslash, even in namespaced files. The same applies to other global types used in catch blocks or parameter hints.

```php
// ── PHP ──────────────────────────────────────────────────────

// ❌ FORBIDDEN
catch (\Throwable $e)
function foo(\Throwable $e): array

// ✅ REQUIRED
catch (Throwable $e)
function foo(Throwable $e): array
```

```typescript
// ── TypeScript / Go ─────────────────────────────────────────
// Not applicable — these languages don't have leading-backslash syntax.
```

---

## Rule 14: No Commented-Out or Dead Code

Commented-out code and unreachable (dead) code **must be deleted**, not left in the codebase. Version control preserves history — comments are not an archive tool.

```php
// ── PHP ──────────────────────────────────────────────────────

// ❌ FORBIDDEN: Commented-out code
// $oldValue = $this->legacyLookup($key);
// if ($oldValue !== null) {
//     return $oldValue;
// }
$value = $this->lookup($key);

// ❌ FORBIDDEN: Dead code after unconditional return
return $result;
$this->cleanup();   // never executes

// ✅ REQUIRED: Remove dead code entirely
$value = $this->lookup($key);
```

```typescript
// ── TypeScript ───────────────────────────────────────────────

// ❌ FORBIDDEN: Commented-out code
// const legacyResult = await fetchLegacy(id);
// return legacyResult;
const result = await fetchData(id);

// ✅ REQUIRED: Clean — no commented-out code
const result = await fetchData(id);
```

```go
// ── Go ───────────────────────────────────────────────────────

// ❌ FORBIDDEN: Commented-out code
// oldVal, err := legacyFetch(ctx, key)
// if err != nil {
//     return err
// }
val, err := fetch(ctx, key)

// ✅ REQUIRED: Clean
val, err := fetch(ctx, key)
```

### Exceptions

- **TODO/FIXME comments** with a ticket reference are allowed: `// TODO(PROJ-123): migrate to new API`
- **Intentional stubs** for future implementation are allowed when marked: `// STUB: placeholder for upcoming feature PROJ-456`

---

## Rule 15: Comment Formatting — Space After `//`

Every line comment **must** have a space between the comment marker and the text. Block comments must also have consistent spacing.

```php
// ── PHP ──────────────────────────────────────────────────────

// ❌ FORBIDDEN
//this is a comment
//$value = 5;

// ✅ REQUIRED
// This is a comment
// $value = 5;
```

```typescript
// ── TypeScript ───────────────────────────────────────────────

// ❌ FORBIDDEN
//calculate total
const total = items.reduce((sum, i) => sum + i.price, 0);

// ✅ REQUIRED
// Calculate total
const total = items.reduce((sum, i) => sum + i.price, 0);
```

```go
// ── Go ───────────────────────────────────────────────────────

// ❌ FORBIDDEN
//fetchData retrieves data from the store
func fetchData(ctx context.Context) error {

// ✅ REQUIRED
// fetchData retrieves data from the store.
func fetchData(ctx context.Context) error {
```

---

## Rule 16: Method and Function Documentation

Every **exported/public** function or method **must** have a doc comment describing its purpose. Internal/private helpers are encouraged but not required.

```php
// ── PHP ──────────────────────────────────────────────────────

// ❌ FORBIDDEN: No doc comment on public method
public function processUpload(UploadRequest $request): UploadResult {

// ✅ REQUIRED
/**
 * Validates and executes the upload, returning the stored file metadata.
 */
public function processUpload(UploadRequest $request): UploadResult {
```

```typescript
// ── TypeScript ───────────────────────────────────────────────

// ❌ FORBIDDEN: No doc comment on exported function
export const calculateDiscount = (price: number, tier: PricingTier): number => {

// ✅ REQUIRED
/** Calculates the discount amount based on the customer pricing tier. */
export const calculateDiscount = (price: number, tier: PricingTier): number => {
```

```go
// ── Go ───────────────────────────────────────────────────────

// ❌ FORBIDDEN: No doc comment on exported function
func ProcessUpload(ctx context.Context, req UploadRequest) (UploadResult, error) {

// ✅ REQUIRED — Go convention: comment starts with function name
// ProcessUpload validates and executes the upload request.
func ProcessUpload(ctx context.Context, req UploadRequest) (UploadResult, error) {
```

---

*Part of [Code Style](./00-overview.md) — Rules 8, 14, 15, 16*
