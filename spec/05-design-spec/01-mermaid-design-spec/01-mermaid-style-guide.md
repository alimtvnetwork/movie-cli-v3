# Mermaid Flowchart Style Guide

**Version:** 1.0.0  
**Updated:** 2026-04-06  
**Status:** Active

---

## 1. Purpose

Define consistent visual styling for all Mermaid flowchart diagrams across the movie project spec. Colors derive from the project's design tokens (see [Design Tokens](../../02-error-manage-spec/02-error-architecture/04-error-modal/04-color-themes/01-design-tokens.md)).

---

## 2. Color Palette

Derived from the project's HSL design tokens, converted to hex for Mermaid compatibility.

| Role | Token | Hex | Usage |
|------|-------|-----|-------|
| **Primary** | `--primary` | `#1e293b` | Entry/exit nodes, command start |
| **Primary text** | `--primary-foreground` | `#f1f5f9` | Text on primary nodes |
| **Success** | `--success` | `#22c55e` | Happy-path outcomes, done states |
| **Success text** | | `#ffffff` | Text on success nodes |
| **Error** | `--destructive` | `#ef4444` | Error outcomes, failures |
| **Error text** | | `#ffffff` | Text on error nodes |
| **Warning** | `--warning` | `#f59e0b` | Warnings, skipped states |
| **Warning text** | | `#ffffff` | Text on warning nodes |
| **Info** | `--info` | `#3b82f6` | Informational, fetch/API steps |
| **Info text** | | `#ffffff` | Text on info nodes |
| **Process** | `--muted` | `#f1f5f9` | Standard processing steps |
| **Process text** | | `#1e293b` | Text on process nodes |
| **Decision** | | `#e2e8f0` | Decision diamonds |
| **Decision text** | | `#1e293b` | Text on decision nodes |

---

## 3. Node Shape Conventions

| Shape | Mermaid Syntax | Usage |
|-------|---------------|-------|
| Stadium | `([text])` | Entry and exit points |
| Rectangle | `[text]` | Standard process steps |
| Diamond | `{text}` | Decision / conditional |
| Parallelogram | `[/text/]` | Output / display / error message |
| Rounded rect | `(text)` | Prompts / user interaction |

---

## 4. classDef Definitions

Place these at the bottom of every `.mmd` file:

```mermaid
classDef entryExit fill:#1e293b,stroke:#1e293b,color:#f1f5f9,stroke-width:2px
classDef process fill:#f1f5f9,stroke:#cbd5e1,color:#1e293b
classDef decision fill:#e2e8f0,stroke:#94a3b8,color:#1e293b
classDef success fill:#22c55e,stroke:#16a34a,color:#ffffff
classDef error fill:#ef4444,stroke:#dc2626,color:#ffffff
classDef warning fill:#f59e0b,stroke:#d97706,color:#ffffff
classDef info fill:#3b82f6,stroke:#2563eb,color:#ffffff
classDef prompt fill:#8b5cf6,stroke:#7c3aed,color:#ffffff
```

---

## 5. Class Assignment Rules

| Node type | Class | Examples |
|-----------|-------|----------|
| Start / End | `entryExit` | `START`, `DONE`, `BYE` |
| Process step | `process` | `Open DB`, `Build record`, `Count media` |
| Decision | `decision` | `{Found?}`, `{Valid?}`, `{Type?}` |
| Success outcome | `success` | `Saved!`, `Move successful`, `Done` |
| Error / failure | `error` | `Error: not found`, `Move failed` |
| Warning / skip | `warning` | `Warn: no match`, `Skip`, `Already exists` |
| API / fetch | `info` | `SearchMulti`, `Fetch trending`, `Download poster` |
| User prompt | `prompt` | `Prompt: confirm`, `Wait for input` |

---

## 6. Applying Classes

Use the `:::className` syntax inline or batch `class` statements at the end:

```mermaid
%% Inline syntax
A([Start]):::entryExit --> B[Open DB]:::process
B --> C{Found?}:::decision
C -->|yes| D[/Saved!/]:::success
C -->|no| E[/Error: not found/]:::error

%% Or batch syntax at end
class A entryExit
class B process
class C decision
class D success
class E error
```

**Preferred:** Batch `class` statements at the bottom for readability.

---

## 7. Edge Labels

- Keep labels short: `yes`, `no`, `valid`, `invalid`, `movie`, `tv`
- Use `|label|` syntax on edges
- No styling on edges (use default Mermaid arrows)

---

## 8. Header Comment

Every diagram file must start with:

```
%% filename.mmd — command description
%% Version: X.Y.Z | Updated: YYYY-MM-DD
```

---

## 9. Checklist for New Diagrams

- [ ] Header comment with version and date
- [ ] `flowchart TD` direction
- [ ] Stadium shapes for entry/exit
- [ ] Parallelogram for output/error messages
- [ ] All classDef definitions present at bottom
- [ ] Every node assigned a class
- [ ] No emoji in node text
- [ ] Renders without errors in Mermaid live editor

---

## Cross-References

- [Design Tokens](../../02-error-manage-spec/02-error-architecture/04-error-modal/04-color-themes/01-design-tokens.md)
- [App Diagrams Index](../../08-app/04-diagrams/00-overview.md)

---

*Mermaid style guide — updated: 2026-04-06*
