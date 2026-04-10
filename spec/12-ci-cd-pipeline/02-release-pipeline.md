# Release Pipeline

## Overview

The release pipeline automates binary production, packaging, and GitHub Release creation whenever code is pushed to a `release/**` branch or a `v*` tag. It produces cross-compiled binaries for 6 platform/architecture targets, platform-specific install scripts, checksums, and a fully formatted release page.

> **Reference**: Adapted from gitmap-v2 release patterns ([source](https://github.com/alimtvnetwork/gitmap-v2/blob/main/spec/generic-release/02-release-pipeline.md))

---

## Trigger and Concurrency

### Trigger

```yaml
on:
  push:
    branches:
      - "release/**"
    tags:
      - "v*"
```

### Concurrency

Release pipelines **never cancel** in-progress runs — every release commit must produce complete artifacts:

```yaml
concurrency:
  group: release-${{ github.ref }}
  cancel-in-progress: false
```

### Permissions

The workflow needs write access to create GitHub Releases:

```yaml
permissions:
  contents: write
```

---

## Version Resolution

The version is extracted from the Git ref:
- **Tags** (`refs/tags/v1.2.3`): use the tag name directly
- **Branches** (`refs/heads/release/v1.2.3`): strip the `release/` prefix

```bash
if [[ "$GITHUB_REF" == refs/tags/* ]]; then
  VERSION="${GITHUB_REF_NAME}"
elif [[ "$GITHUB_REF" == refs/heads/release/* ]]; then
  VERSION="${GITHUB_REF_NAME#release/}"
fi
```

---

## Icon Embedding

Windows binaries include an embedded application icon (`assets/icon.ico`). The pipeline uses `go-winres` (v0.3.3) to generate a `.syso` resource file before compilation:

```bash
go install github.com/tc-hib/go-winres@v0.3.3
go-winres init --icon assets/icon.ico
```

This runs once before the build loop. The `.syso` file is automatically linked by `go build` for Windows targets. Non-Windows binaries are unaffected.

---

## Binary Building

### Targets

Build for 6 platform/architecture combinations:

| Platform | Architecture | Extension |
|----------|-------------|-----------|
| Windows | amd64 | `.exe` |
| Windows | arm64 | `.exe` |
| Linux | amd64 | (none) |
| Linux | arm64 | (none) |
| macOS | amd64 | (none) |
| macOS | arm64 | (none) |

### Build Command Pattern

```bash
CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" go build \
  -ldflags "-s -w \
    -X 'github.com/mahin/mahin-cli-v2/version.Version=$VERSION' \
    -X 'github.com/mahin/mahin-cli-v2/version.Commit=$COMMIT' \
    -X 'github.com/mahin/mahin-cli-v2/version.BuildDate=$BUILD_DATE'" \
  -o "dist/movie-${VERSION}-${os}-${arch}${ext}" .
```

### Naming Convention

- Binary: `movie-{version}-{os}-{arch}{ext}`
- Windows zip: `movie-{version}-windows-{arch}.zip`
- Unix tarball: `movie-{version}-{os}-{arch}.tar.gz`

### Build Once Rule

Binaries are compiled **exactly once**. All downstream steps (compress, checksum, publish) reuse the same artifacts and must never trigger a rebuild.

---

## Packaging

### Compression

- **Windows**: `.exe` → `.zip` (using `zip`)
- **Linux / macOS**: binary → `.tar.gz` (using `tar czf`)

### Checksums

After compression, generate SHA256 checksums for all artifacts:

```bash
sha256sum * > checksums.txt
```

Output format: `<hash>  <filename>` (two spaces between hash and name).

---

## Install Scripts

Two platform-specific install scripts are generated with each release and uploaded as release assets.

### PowerShell (`install.ps1`) — Windows

One-liner install:
```powershell
irm https://github.com/REPO/releases/download/VERSION/install.ps1 | iex
```

**Important**: The script must NOT use a top-level `param()` block — `irm | iex` pipes content to `Invoke-Expression`, which cannot bind parameters.

Features:
- Auto-detect CPU architecture (amd64/arm64)
- Download + SHA256 verification
- Install to `$env:LOCALAPPDATA\movie`
- Add to user PATH
- Verify installation by running `movie version`

### Bash (`install.sh`) — Linux / macOS

One-liner install:
```bash
curl -fsSL https://github.com/REPO/releases/download/VERSION/install.sh | bash
```

Features:
- Auto-detect OS (linux/darwin) and architecture (amd64/arm64)
- Download + SHA256 verification via `sha256sum` or `shasum`
- Install to `~/.local/bin`
- Add to shell rc file (bash/zsh/fish)
- Verify installation by running `movie version`

### Placeholder Substitution

Both scripts are generated with `VERSION_PLACEHOLDER` and `REPO_PLACEHOLDER` tokens, which are replaced via `sed` during the pipeline run.

---

## Release Page

The release page includes:

1. **Changelog entry** — extracted from `CHANGELOG.md` (if present)
2. **Release info table** — version, commit, branch, build date, Go version
3. **SHA256 checksums** — full `checksums.txt` content
4. **Install instructions** — one-liners for Windows and Unix
5. **Asset table** — platform × architecture matrix with filenames

Created using `softprops/action-gh-release@v2`.

---

## How to Create a Release

### Option A: Push a release branch

```bash
git checkout -b release/v1.3.0
git push origin release/v1.3.0
```

### Option B: Push a tag

```bash
git tag v1.3.0
git push origin v1.3.0
```

Both trigger the same pipeline. The version is resolved from the ref name.

---

## Acceptance Criteria

- GIVEN a `release/v1.3.0` branch push WHEN pipeline runs THEN 6 binaries are built with version `v1.3.0`
- GIVEN a `v1.3.0` tag push WHEN pipeline runs THEN the same release artifacts are produced
- GIVEN a Windows binary WHEN packaged THEN it is compressed as `.zip`
- GIVEN all binaries WHEN checksums are generated THEN `checksums.txt` contains one entry per archive
- GIVEN install scripts WHEN generated THEN `VERSION_PLACEHOLDER` is replaced with the actual version
- GIVEN a CHANGELOG.md entry WHEN release page is built THEN the entry appears in the release body

---

*Release pipeline spec — updated: 2026-04-10*
