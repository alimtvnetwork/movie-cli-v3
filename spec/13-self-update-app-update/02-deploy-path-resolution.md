# 02 — Deploy Path Resolution

## Purpose

Define how the build/deploy script determines **where to install** the new binary. The resolution must handle first-time installs, existing installs, and explicit overrides.

> **Reference**: Adapted from gitmap-v2 ([02-deploy-path-resolution.md](https://github.com/alimtvnetwork/gitmap-v2/blob/main/spec/generic-update/02-deploy-path-resolution.md))

---

## 3-Tier Resolution Strategy

The deploy target is resolved with this priority:

| Priority | Source | When Used |
|----------|--------|-----------|
| 1 | **CLI flag** (`--deploy-path` or `-DeployPath`) | User explicitly specifies a path |
| 2 | **Global PATH lookup** | Binary is already installed and on PATH |
| 3 | **Config file default** | First-time install or binary not on PATH |

### Tier 1 — CLI Flag Override

If the user passes an explicit path, use it unconditionally:

```powershell
# PowerShell (run.ps1)
if ($DeployPath.Length -gt 0) {
    return $DeployPath
}
```

```bash
# Bash
if [[ -n "$DEPLOY_PATH" ]]; then
    echo "$DEPLOY_PATH"
    return
fi
```

### Tier 2 — Global PATH Lookup

If the binary is already installed and accessible via `PATH`, detect its current location and deploy there:

```powershell
# PowerShell
$activeCmd = Get-Command movie -ErrorAction SilentlyContinue
if ($activeCmd) {
    $resolvedPath = (Resolve-Path $activeCmd.Source).Path
    return Split-Path $resolvedPath -Parent
}
```

```bash
# Bash
active_cmd=$(command -v movie 2>/dev/null || true)
if [[ -n "$active_cmd" ]]; then
    resolved=$(readlink -f "$active_cmd" 2>/dev/null || realpath "$active_cmd" 2>/dev/null || echo "$active_cmd")
    echo "$(dirname "$resolved")"
    return
fi
```

### Tier 3 — Config File Default

Read from `powershell.json`:

```json
{
    "deployPath": "E:\\bin-run",
    "buildOutput": "./bin",
    "binaryName": "movie.exe"
}
```

Platform defaults if config is missing:
- **Windows**: `E:\bin-run`
- **macOS / Linux**: `/usr/local/bin`

---

## Mahin-Specific Defaults

| Platform | Default Deploy Path | Config Key |
|----------|-------------------|------------|
| Windows | `E:\bin-run` | `powershell.json → deployPath` |
| macOS | `/usr/local/bin` | `powershell.json → deployPath` |
| Linux | `/usr/local/bin` | `powershell.json → deployPath` |
| GitHub Release Install (Windows) | `$env:LOCALAPPDATA\movie` | Built into `install.ps1` |
| GitHub Release Install (Unix) | `~/.local/bin` | Built into `install.sh` |

---

## Acceptance Criteria

- GIVEN `--deploy-path /custom/dir` WHEN build runs THEN binary is deployed to `/custom/dir`
- GIVEN `movie` is on PATH at `/usr/local/bin/movie` WHEN build runs without `--deploy-path` THEN binary deploys to `/usr/local/bin`
- GIVEN `movie` is NOT on PATH WHEN build runs THEN config default is used
- GIVEN no config file exists WHEN build runs THEN platform default is used

---

*Deploy path resolution — updated: 2026-04-10*
