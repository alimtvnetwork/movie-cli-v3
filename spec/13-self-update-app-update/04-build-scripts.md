# 04 — Build Scripts

## Purpose

Define the cross-platform build scripts (`run.ps1` and `build.ps1`) that automate the full pipeline: pull → build → deploy.

> **Reference**: Adapted from gitmap-v2 ([04-build-scripts.md](https://github.com/alimtvnetwork/gitmap-v2/blob/main/spec/generic-update/04-build-scripts.md))

---

## Script Inventory

| Script | Platform | Purpose |
|--------|----------|---------|
| `run.ps1` | Windows (PowerShell 5.1+) | Full pipeline: pull → tidy → build → deploy |
| `build.ps1` | Windows (PowerShell 5.1+) | Build + deploy without git pull |

Both scripts can also run on macOS/Linux via PowerShell 7+ (`pwsh`).

---

## Pipeline Steps

Both scripts implement the same 4-step pipeline:

```
[1/4] Pull latest changes (git pull, branch check)
[2/4] Resolve dependencies (go mod tidy)
[3/4] Build binary (go build with ldflags)
[4/4] Deploy (rename-first to resolved target)
```

---

## Configuration

Read build settings from `powershell.json`:

```json
{
    "deployPath": "E:\\bin-run",
    "buildOutput": "./bin",
    "binaryName": "movie.exe",
    "copyData": false
}
```

| Key | Default | Purpose |
|-----|---------|---------|
| `deployPath` | `E:\bin-run` (Win) / `/usr/local/bin` (Unix) | Where to deploy the built binary |
| `buildOutput` | `./bin` | Intermediate build output directory |
| `binaryName` | `movie.exe` | Output binary name |
| `copyData` | `false` | Whether to copy data files alongside binary |

---

## Build Step

```powershell
$version = git describe --tags --always 2>$null
$commit = (git rev-parse --short HEAD)
$buildDate = (Get-Date -Format "yyyy-MM-dd")

$ldflags = "-s -w " +
    "-X 'github.com/mahin/mahin-cli-v2/version.Version=$version' " +
    "-X 'github.com/mahin/mahin-cli-v2/version.Commit=$commit' " +
    "-X 'github.com/mahin/mahin-cli-v2/version.BuildDate=$buildDate'"

go build -ldflags $ldflags -o "$binDir/$binaryName" .
```

---

## Deploy Step

Uses rename-first strategy (see [03-rename-first-deploy.md](./03-rename-first-deploy.md)):

1. Resolve deploy path (3-tier strategy)
2. Rename existing binary → `.old`
3. Copy new binary to deploy path
4. Verify with `movie version`
5. Clean up `.old` backup

---

## Acceptance Criteria

- GIVEN `run.ps1` is executed WHEN all steps pass THEN the binary is deployed and `movie version` shows the new version
- GIVEN `powershell.json` exists WHEN the script reads config THEN custom deploy path is used
- GIVEN no `powershell.json` WHEN the script reads config THEN platform defaults are used
- GIVEN `build.ps1` is executed WHEN git pull is skipped THEN the build still completes with the current code

---

*Build scripts — updated: 2026-04-10*
