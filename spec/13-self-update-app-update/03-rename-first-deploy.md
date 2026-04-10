# 03 — Rename-First Deploy Strategy

## Purpose

Define the file-replacement strategy that works on all platforms, including Windows where running executables cannot be overwritten.

> **Reference**: Adapted from gitmap-v2 ([03-rename-first-deploy.md](https://github.com/alimtvnetwork/gitmap-v2/blob/main/spec/generic-update/03-rename-first-deploy.md))

---

## The Problem

On Windows, copying a file over a running `.exe` fails with:

```
The process cannot access the file 'movie.exe' because it is being
used by another process.
```

This happens because:
- The OS holds a file lock on any running executable.
- `Copy-Item` / `cp` tries to open the destination for writing, which the lock prevents.

---

## The Solution: Rename-First

Windows allows **renaming** a running executable (the OS tracks the process by handle, not by filename). So instead of overwriting:

```
Step 1: Rename existing movie.exe → movie.exe.old
Step 2: Copy new binary → movie.exe (destination is now free)
```

After the rename, the original process continues running from the renamed file. The new binary occupies the original filename, ready for the next invocation.

---

## Implementation

### PowerShell

```powershell
$destFile = Join-Path $deployDir "movie.exe"
$backupFile = "$destFile.old"
$hasBackup = $false

if (Test-Path $destFile) {
    try {
        if (Test-Path $backupFile) {
            Remove-Item $backupFile -Force -ErrorAction SilentlyContinue
        }
        Rename-Item $destFile $backupFile -Force -ErrorAction Stop
        $hasBackup = $true
    } catch {
        Write-Warning "Rename-first failed: $_"
    }
}

# Copy new binary
Copy-Item $sourceBinary $destFile -Force

# Clean up backup
if ($hasBackup -and (Test-Path $backupFile)) {
    Remove-Item $backupFile -Force -ErrorAction SilentlyContinue
}
```

### Bash

```bash
dest_file="$deploy_dir/movie"
backup_file="$dest_file.old"

if [[ -f "$dest_file" ]]; then
    mv "$dest_file" "$backup_file" 2>/dev/null || true
fi

cp "$source_binary" "$dest_file"
chmod +x "$dest_file"

rm -f "$backup_file" 2>/dev/null || true
```

---

## Rollback

If the new binary fails verification:

```powershell
if ($hasBackup -and (Test-Path $backupFile)) {
    Remove-Item $destFile -Force -ErrorAction SilentlyContinue
    Rename-Item $backupFile $destFile -Force
    Write-Error "Rolled back to previous version"
}
```

---

## Cross-Drive Considerations

`os.Rename` (Go) and `mv` (shell) fail across filesystems with `EXDEV`. The mahin CLI already handles this in `cmd/movie_move_helpers.go` with a copy+delete fallback. The same pattern applies to deploy:

```go
func MoveFile(src, dst string) error {
    err := os.Rename(src, dst)
    if err != nil && isExdevError(err) {
        return copyAndRemove(src, dst)
    }
    return err
}
```

---

## Acceptance Criteria

- GIVEN a running `movie.exe` on Windows WHEN deploy runs THEN the old binary is renamed to `.old` and the new binary takes its place
- GIVEN deploy fails after rename WHEN rollback runs THEN the `.old` backup is restored
- GIVEN a cross-filesystem deploy WHEN `os.Rename` fails with EXDEV THEN the copy+delete fallback is used
- GIVEN a successful deploy WHEN cleanup runs THEN the `.old` backup is removed

---

*Rename-first deploy — updated: 2026-04-10*
