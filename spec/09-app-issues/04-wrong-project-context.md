# Issue #04: Wrong Project Context Applied

## Issue Summary

1. **What happened**: A request intended for a Chrome extension project was sent to this Go CLI project. The request referenced PHP, WordPress, TypeScript, React error components, and Chrome extension architecture — none of which exist here.
2. **Where**: User message / AI session context.
3. **Symptoms and impact**: If the AI had proceeded, it would have created irrelevant spec files, TypeScript types, and React components in a Go CLI project, causing severe confusion.
4. **How discovered**: AI identified the mismatch and asked a clarifying question. User confirmed it was the wrong project.

## Root Cause Analysis

1. **Direct cause**: User had multiple Lovable projects open and sent the message to the wrong one.
2. **Contributing factors**: The proofread prompt was highly detailed and specific, making it look intentional. AI could have blindly followed instructions.
3. **Triggering conditions**: Multi-project workflow without project name verification.
4. **Why spec did not prevent it**: Project overview stated "Go CLI, not a web app" but the AI needed to actively check this against incoming requests.

## Fix Description

1. **Spec change**: Added explicit scope statement to `spec/08-app/01-project-spec.md`: "This is a Go CLI project only — no web frontend, no PHP, no WordPress."
2. **New rule**: AI must validate that incoming requests match the project type before proceeding. If a request references technologies not in this project (PHP, WordPress, TypeScript, React, Chrome), stop and ask.
3. **Why it resolves root cause**: Forces a context check before any implementation.
4. **Config changes**: None.
5. **Diagnostics**: None required.

## Iterations History

1. **Iteration 1** (17-Mar-2026): AI asked clarifying question using `ask_questions` tool. User confirmed wrong project. No damage done.

## Prevention and Non-Regression

1. **Prevention rule**: Before processing any request, verify it matches the project type (Go CLI). If it references PHP, WordPress, TypeScript, React, Node.js, or Chrome extensions — stop and clarify.
2. **Acceptance criteria**: AI never creates files or specs for technologies outside the project scope without explicit user confirmation.
3. **Guardrails**: Project overview prominently states "This is NOT a web project." AI success plan Rule 1 requires reading project overview first.
4. **Spec references**: `spec/08-app/01-project-spec.md` (Scope section), `.lovable/memory/01-project-overview.md`.

## TODO and Follow-ups

- [x] Clarification obtained
- [x] Spec updated with scope statement
- [x] Memory updated
- [ ] No code changes needed — this was a process issue

## Done Checklist

- [x] Spec updated under `/spec/08-app/`
- [x] Issue write-up created under `/spec/09-app-issues/`
- [x] Memory updated with summary and prevention rule
- [x] Acceptance criteria updated
- [x] Iterations recorded
