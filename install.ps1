<#
.SYNOPSIS
 One-step bootstrap: clone (if needed), build, and deploy movie CLI.
.DESCRIPTION
 Detects OS, clones the repo if not already present, runs the full
 build pipeline via run.ps1, and verifies the installation.
.EXAMPLES
 # Fresh install (from anywhere)
 pwsh install.ps1

 # Re-install / update (from repo root)
 pwsh install.ps1

 # Custom deploy path
 pwsh install.ps1 -DeployPath ~/bin
.NOTES
 Requires: Git, Go 1.22+, PowerShell 5.1+ (Windows) or 7+ (cross-platform)
#>

[CmdletBinding()]
param(
    [string]$DeployPath = ""
)

$ErrorActionPreference = "Stop"

# -- Helpers ---------------------------------------------------

function Write-Banner {
    Write-Host ""
    Write-Host " +======================================+" -ForegroundColor DarkCyan
    Write-Host " | " -ForegroundColor DarkCyan -NoNewline
    Write-Host "movie installer" -ForegroundColor Cyan -NoNewline
    Write-Host "                  |" -ForegroundColor DarkCyan
    Write-Host " +======================================+" -ForegroundColor DarkCyan
    Write-Host ""
}

function Write-Ok    { param([string]$M) Write-Host "  OK " -ForegroundColor Green -NoNewline; Write-Host $M -ForegroundColor Green }
function Write-Info  { param([string]$M) Write-Host "  -> " -ForegroundColor Cyan -NoNewline; Write-Host $M -ForegroundColor Gray }
function Write-Err   { param([string]$M) Write-Host "  XX " -ForegroundColor Red -NoNewline; Write-Host $M -ForegroundColor Red }
function Write-ErrorAndExit {
    param([string]$Message, [string]$Hint = "")
    Write-Err $Message
    if ($Hint) { Write-Info $Hint }
    exit 1
}

# -- Pre-flight checks -----------------------------------------

Write-Banner

Write-Host " [1/4] Checking prerequisites" -ForegroundColor Magenta
Write-Host (" " + ("-" * 50)) -ForegroundColor DarkGray

# Check Git
$prevPref = $ErrorActionPreference; $ErrorActionPreference = "Continue"
$gitVer = git --version 2>&1; $gitExit = $LASTEXITCODE
$ErrorActionPreference = $prevPref

if ($gitExit -ne 0) {
    Write-ErrorAndExit "Git is not installed or not in PATH" "Install from https://git-scm.com/downloads"
}
Write-Ok "Git: $("$gitVer".Trim())"

# Check Go
$prevPref = $ErrorActionPreference; $ErrorActionPreference = "Continue"
$goVer = go version 2>&1; $goExit = $LASTEXITCODE
$ErrorActionPreference = $prevPref

if ($goExit -ne 0) {
    Write-ErrorAndExit "Go is not installed or not in PATH" "Install from https://go.dev/dl/"
}
Write-Ok "Go: $("$goVer".Trim())"

# -- Locate or clone repo -------------------------------------

Write-Host ""
Write-Host " [2/4] Locating repository" -ForegroundColor Magenta
Write-Host (" " + ("-" * 50)) -ForegroundColor DarkGray

$RepoName = "movie-cli-v2"
$RepoUrl  = "https://github.com/movie/movie-cli-v2.git"

# Check if we're already inside the repo
$inRepo = (Test-Path "go.mod") -and (Test-Path "run.ps1")

if ($inRepo) {
    $RepoRoot = (Get-Location).Path
    Write-Ok "Already inside repo: $RepoRoot"
} else {
    # Check if repo exists as a subdirectory
    $subDir = Join-Path (Get-Location).Path $RepoName
    if (Test-Path (Join-Path $subDir "go.mod")) {
        $RepoRoot = $subDir
        Write-Ok "Found repo: $RepoRoot"
    } else {
        Write-Info "Cloning $RepoUrl ..."
        $prevPref = $ErrorActionPreference; $ErrorActionPreference = "Continue"
        $cloneOutput = git clone $RepoUrl 2>&1
        $cloneExit = $LASTEXITCODE
        $ErrorActionPreference = $prevPref

        if ($cloneExit -ne 0) {
            Write-Err "Clone failed"
            foreach ($line in $cloneOutput) { Write-Host "  $line" -ForegroundColor Red }
            exit 1
        }
        $RepoRoot = $subDir
        Write-Ok "Cloned to: $RepoRoot"
    }
}

# -- Run build pipeline ----------------------------------------

Write-Host ""
Write-Host " [3/4] Running build pipeline" -ForegroundColor Magenta
Write-Host (" " + ("-" * 50)) -ForegroundColor DarkGray

$runScript = Join-Path $RepoRoot "run.ps1"

if (-not (Test-Path $runScript)) {
    Write-ErrorAndExit "run.ps1 not found at $runScript"
}

$runArgs = @()
if ($DeployPath) {
    $runArgs += "-DeployPath"
    $runArgs += $DeployPath
}

Push-Location $RepoRoot
try {
    & $runScript @runArgs
    if ($LASTEXITCODE -and $LASTEXITCODE -ne 0) {
        Write-ErrorAndExit "Build pipeline failed (exit $LASTEXITCODE)"
    }
} finally {
    Pop-Location
}

# -- Verify ----------------------------------------------------

Write-Host ""
Write-Host " [4/4] Verifying installation" -ForegroundColor Magenta
Write-Host (" " + ("-" * 50)) -ForegroundColor DarkGray

$prevPref = $ErrorActionPreference; $ErrorActionPreference = "Continue"
$verOutput = movie version 2>&1
$verExit = $LASTEXITCODE
$ErrorActionPreference = $prevPref

if ($verExit -eq 0) {
    Write-Ok "movie is ready: $("$verOutput".Trim())"
} else {
    foreach ($line in $verOutput) { Write-Host "    $line" -ForegroundColor Red }
    Write-ErrorAndExit "Verification failed — 'movie version' returned exit code $verExit" "Add the deploy directory to your PATH, then try again"
}

Write-Host ""
Write-Host " +======================================+" -ForegroundColor DarkCyan
Write-Host " | " -ForegroundColor DarkCyan -NoNewline
Write-Host "Installation complete" -ForegroundColor Green -NoNewline
Write-Host "            |" -ForegroundColor DarkCyan
Write-Host " +======================================+" -ForegroundColor DarkCyan
Write-Host ""
