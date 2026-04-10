# 02 — PowerShell Build & Deploy

> Summary spec for the `run.ps1` automation pipeline.

## Overview

`run.ps1` is the single-entry automation script for building, deploying, and running the **movie** CLI. It replaces ad-hoc manual steps with a repeatable, configurable pipeline.

## Pipeline Phases

| Phase | Step | Function | Flag to Skip |
|-------|------|----------|--------------|
| 1 | Git pull | `Invoke-GitPull` | `-NoPull` |
| 2 | Resolve deps | `Resolve-Dependencies` | — |
| 3 | Build binary | `Build-Binary` | — |
| 4 | Deploy | `Deploy-Binary` | `-NoDeploy` |
| — | Run | `Invoke-Run` | only with `-R` |
| — | Test | `Invoke-Tests` | only with `-t` |

## Configuration

All configuration is loaded from `powershell.json` at the repo root. See [03-config-reference.md](03-config-reference.md) for field details.

## Sub-Documents

| File | Content |
|------|---------|
| [01-install-guide.md](01-install-guide.md) | Prerequisites, setup, verification |
| [02-powershell-build-deploy.md](02-powershell-build-deploy.md) | This file — pipeline overview |
| [03-config-reference.md](03-config-reference.md) | `powershell.json` field reference |
| [04-run-guide.md](04-run-guide.md) | Usage, flags, examples, troubleshooting |

## Quick Reference

```powershell
# Full pipeline (pull → build → deploy)
.\run.ps1

# Build only, no git, no deploy
.\run.ps1 -NoPull -NoDeploy

# Build + deploy + run scan
.\run.ps1 -R scan D:\movies

# Run all unit tests
.\run.ps1 -t

# CI mode: force-pull, build, deploy
.\run.ps1 -ForcePull
```
