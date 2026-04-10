# TypeScript ConnectionStatus Enum — `src/lib/enums/connection-status.ts`

> **Version**: 1.0.0  
> **Last updated**: 2026-02-27  
> **Tracks**: Issue #10 (`spec/23-how-app-issues-track/10-domain-status-magic-strings.md`)

---

## Purpose

Typed enum for WebSocket, SSE, and service connection lifecycle states. Replaces `connection.status === 'connected'` magic strings in frontend specs.

---

## Reference Implementation

```typescript
// src/lib/enums/connection-status.ts

export enum ConnectionStatus {
  Connected = "CONNECTED",
  Disconnected = "DISCONNECTED",
  Connecting = "CONNECTING",
  Reconnecting = "RECONNECTING",
  Error = "ERROR",
}
```

---

## Usage Patterns

### Status Comparisons

```typescript
// ❌ WRONG: Magic string
if (connection.status === 'connected') { ... }

// ✅ CORRECT: Enum constant
if (connection.status === ConnectionStatus.Connected) { ... }
```

### Conditional Rendering

```typescript
// ❌ WRONG
{wsStatus === 'disconnected' && <ReconnectBanner />}

// ✅ CORRECT
{wsStatus === ConnectionStatus.Disconnected && <ReconnectBanner />}
```

### Type Definitions

```typescript
// ❌ WRONG
interface WebSocketState {
  status: 'connected' | 'disconnected' | 'connecting' | 'error';
}

// ✅ CORRECT
interface WebSocketState {
  status: ConnectionStatus;
}
```

---

## Consuming Spec Files

| Spec File | Pattern Replaced |
|-----------|-----------------|
| `05-features/05-voice-input/06-voice-session-manager.md` | `connection.status === 'connected'` |
| `01-general-spec/09-api-integration/02-websocket-patterns-api-integration.md` | WebSocket connection status checks |
| `05-features/27-automation-pipeline/24-collaboration.md` | Participant connection state |
| `08-roadmap-overview/05-gap-analysis.md` | Connection status references |
| `16-ai-transcribe-cli/02-frontend/01-testing-ui.md` | Recording connection status |

---

## Cross-Language Parity

| Feature | Go | TypeScript |
|---------|-----|-----------|
| Package | `pkg/enums/connectionstatus` | `src/lib/enums/connection-status.ts` |
| Type | `byte` iota | String enum |
| Values | `Connected`, `Disconnected`, `Connecting`, `Reconnecting`, `Error` | Same |

---

## Cross-References

- Issue #10 — Domain Status Magic Strings <!-- legacy: spec/23-how-app-issues-track/10-domain-status-magic-strings.md — REMOVED — original project issue tracker no longer exists -->
- [HttpMethod Enum](./05-http-method-enum.md) — Sibling enum spec
- [TypeScript Standards](./08-typescript-standards-reference.md) — Parent spec

---

*ConnectionStatus enum v1.0.0 — 2026-02-27*
