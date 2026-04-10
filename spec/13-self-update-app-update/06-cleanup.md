# 06 — Cleanup

## Purpose

Define how temporary artifacts from the update process are identified and removed after a successful update.

> **Reference**: Adapted from gitmap-v2 ([06-cleanup.md](https://github.com/alimtvnetwork/gitmap-v2/blob/main/spec/generic-update/06-cleanup.md))

---

## Artifacts That Need Cleanup

| Artifact | Created By | Location |
|----------|------------|----------|
| `movie.exe.old` | Rename-first deploy | Deploy directory |
| `movie.old` | Rename-first deploy (Unix) | Deploy directory |
| Build output (`./bin/`) | Build script | Repo root |
| Install temp dir | Install scripts | System temp |

---

## Automatic Cleanup

The update process should attempt automatic cleanup after a successful deploy:

### PowerShell (run.ps1 / build.ps1)

```powershell
# After successful deploy and verification
$oldFile = "$destFile.old"
if (Test-Path $oldFile) {
    Remove-Item $oldFile -Force -ErrorAction SilentlyContinue
}
```

### Bash

```bash
# After successful deploy
rm -f "$dest_file.old" 2>/dev/null || true
```

---

## Manual Cleanup

If automatic cleanup fails (e.g., file still locked), users can manually remove artifacts:

```powershell
# Windows
Remove-Item "$env:LOCALAPPDATA\movie\movie.exe.old" -Force

# Or via the binary's deploy directory
Remove-Item "E:\bin-run\movie.exe.old" -Force
```

```bash
# Unix
rm -f ~/.local/bin/movie.old
rm -f /usr/local/bin/movie.old
```

---

## Install Script Cleanup

Both `install.ps1` and `install.sh` use temporary directories that are cleaned up automatically:

### PowerShell

```powershell
$tmpDir = Join-Path $env:TEMP "movie-install-$(Get-Random)"
# ... download, verify, extract ...
finally { Remove-Item $tmpDir -Recurse -Force -ErrorAction SilentlyContinue }
```

### Bash

```bash
tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT
# ... download, verify, extract ...
```

The `trap` ensures cleanup even on errors.

---

## Acceptance Criteria

- GIVEN a successful deploy WHEN cleanup runs THEN `.old` backups are removed
- GIVEN a failed deploy WHEN rollback occurs THEN `.old` backup is restored, not deleted
- GIVEN an install script WHEN the script exits (success or error) THEN the temp directory is removed
- GIVEN a locked `.old` file on Windows WHEN cleanup fails THEN no error is shown (silent continue)

---

*Cleanup — updated: 2026-04-10*
