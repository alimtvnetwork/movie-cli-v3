# 05 — Release Distribution

## Purpose

Define how release binaries are cross-compiled, packaged, checksummed, and distributed to end users via install scripts and GitHub Releases.

> **Reference**: Adapted from gitmap-v2 ([generic-release spec](https://github.com/alimtvnetwork/gitmap-v2/tree/main/spec/generic-release))

---

## Cross-Compilation

### Default Targets

| OS | Architecture | Binary Name | Archive Format |
|----|-------------|-------------|----------------|
| `windows` | `amd64` | `movie.exe` | `.zip` |
| `windows` | `arm64` | `movie.exe` | `.zip` |
| `linux` | `amd64` | `movie` | `.tar.gz` |
| `linux` | `arm64` | `movie` | `.tar.gz` |
| `darwin` | `amd64` | `movie` | `.tar.gz` |
| `darwin` | `arm64` | `movie` | `.tar.gz` |

Windows uses `.zip` because it is natively supported by PowerShell and Explorer. Unix platforms use `.tar.gz` for permission preservation.

### Build Command

```bash
CGO_ENABLED=0 GOOS=<os> GOARCH=<arch> go build \
    -ldflags "-s -w \
        -X 'github.com/mahin/mahin-cli-v2/version.Version=<version>' \
        -X 'github.com/mahin/mahin-cli-v2/version.Commit=<commit>' \
        -X 'github.com/mahin/mahin-cli-v2/version.BuildDate=<date>'" \
    -o <output> .
```

### Key Flags

| Flag | Purpose |
|------|---------|
| `CGO_ENABLED=0` | Produce a fully static binary (no C dependencies) |
| `-s -w` | Strip debug info and DWARF symbols (smaller binaries) |
| `-X` | Embed build-time constants |

---

## Checksums

After compressing all binaries, generate a single `checksums.txt`:

```bash
sha256sum *.zip *.tar.gz > checksums.txt
```

### Verification — PowerShell

```powershell
$expectedHash = (Get-Content checksums.txt |
    Where-Object { $_ -match $archiveName } |
    ForEach-Object { ($_ -split '\s+')[0] })

$actualHash = (Get-FileHash $archivePath -Algorithm SHA256).Hash.ToLower()

if ($actualHash -ne $expectedHash.ToLower()) {
    Write-Host "❌ Checksum mismatch!" -ForegroundColor Red
    exit 1
}
```

### Verification — Bash

```bash
expected=$(grep "$archive_name" checksums.txt | awk '{print $1}')
if command -v sha256sum &>/dev/null; then
    actual=$(sha256sum "$archive_path" | awk '{print $1}')
else
    actual=$(shasum -a 256 "$archive_path" | awk '{print $1}')
fi

if [[ "$actual" != "$expected" ]]; then
    echo "❌ Checksum mismatch!" >&2
    exit 1
fi
```

---

## Install Scripts

Two scripts are generated per release:

| Script | Platform | Invocation |
|--------|----------|------------|
| `install.ps1` | Windows (PowerShell 5.1+) | `irm .../install.ps1 \| iex` |
| `install.sh` | Linux / macOS (Bash 4+) | `curl -fsSL .../install.sh \| bash` |

Both scripts are **version-pinned** at generation time — they download a specific release version, not "latest".

### Pipeline

```
1. Resolve version (pinned at generation time)
2. Detect architecture (amd64 or arm64)
3. Download archive from release assets
4. Download checksums.txt
5. Verify SHA-256 hash
6. Extract binary to install directory
7. Add install directory to user PATH
8. Verify installation: movie version
9. Print summary
```

### Default Install Locations

| Platform | Default Path |
|----------|-------------|
| Windows | `$env:LOCALAPPDATA\movie` |
| Linux | `~/.local/bin` |
| macOS | `~/.local/bin` |

---

## Release Assets

Each GitHub Release includes:

| Asset | Description |
|-------|-------------|
| `movie-<version>-windows-amd64.zip` | Windows x64 binary |
| `movie-<version>-windows-arm64.zip` | Windows ARM64 binary |
| `movie-<version>-linux-amd64.tar.gz` | Linux x64 binary |
| `movie-<version>-linux-arm64.tar.gz` | Linux ARM64 binary |
| `movie-<version>-darwin-amd64.tar.gz` | macOS x64 binary |
| `movie-<version>-darwin-arm64.tar.gz` | macOS ARM64 binary |
| `checksums.txt` | SHA256 hashes for all archives |
| `install.ps1` | Version-pinned Windows installer |
| `install.sh` | Version-pinned Unix installer |

---

## Acceptance Criteria

- GIVEN 6 build targets WHEN cross-compilation runs THEN 6 static binaries are produced
- GIVEN all archives WHEN `sha256sum` generates checksums THEN every archive has an entry
- GIVEN `install.ps1` WHEN a Windows user runs `irm | iex` THEN the binary is installed and verified
- GIVEN `install.sh` WHEN a Linux user runs `curl | bash` THEN the binary is installed and verified
- GIVEN a tampered archive WHEN checksum verification runs THEN the install aborts with an error

---

*Release distribution — updated: 2026-04-10*
