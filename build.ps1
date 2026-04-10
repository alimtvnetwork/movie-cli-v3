# build.ps1 — Pull, build, and deploy movie CLI
# Usage: pwsh build.ps1

param(
    [string]$BinDir
)

$ErrorActionPreference = "Stop"
$BinaryName = "movie"

function Write-ErrorAndExit {
    param([string]$Message, [string]$Hint = "")
    Write-Host "❌ $Message" -ForegroundColor Red
    if ($Hint) { Write-Host "  $Hint" -ForegroundColor Yellow }
    exit 1
}

# Detect OS
$IsWindows_ = ($PSVersionTable.PSEdition -eq "Desktop") -or ($IsWindows -eq $true)
$IsMac_ = ($IsMacOS -eq $true)

# Set default bin directory based on OS
if ($BinDir) {
    $DestDir = $BinDir
} elseif ($IsWindows_) {
    $DestDir = "E:\bin-run"
} elseif ($IsMac_) {
    $DestDir = "/usr/local/bin"
} else {
    # Linux
    $DestDir = "/usr/local/bin"
}

Write-Host ""
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "  🔧 Movie CLI — Build & Deploy" -ForegroundColor Cyan
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host ""

# Step 1: Git pull
Write-Host "📥 Pulling latest changes..." -ForegroundColor Yellow
git pull
if ($LASTEXITCODE -ne 0) {
    Write-ErrorAndExit "Git pull failed!"
}
Write-Host "✅ Pull complete" -ForegroundColor Green
Write-Host ""

# Step 2: Tidy modules
Write-Host "📦 Tidying Go modules..." -ForegroundColor Yellow
go mod tidy
Write-Host "✅ Modules ready" -ForegroundColor Green
Write-Host ""

# Step 3: Build
Write-Host "🔨 Building binary..." -ForegroundColor Yellow
if ($IsWindows_) {
    $BinaryFile = "$BinaryName.exe"
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
} else {
    $BinaryFile = $BinaryName
    if ($IsMac_) {
        $env:GOOS = "darwin"
        # Detect architecture
        $Arch = & uname -m
        if ($Arch -eq "arm64") {
            $env:GOARCH = "arm64"
        } else {
            $env:GOARCH = "amd64"
        }
    } else {
        $env:GOOS = "linux"
        $env:GOARCH = "amd64"
    }
}

# Embed icon for Windows builds
if ($IsWindows_) {
    Write-Host "🎨 Embedding icon into binary..." -ForegroundColor Yellow
    $iconPath = Join-Path $PSScriptRoot "assets" "icon.ico"
    if (Test-Path $iconPath) {
        # Install go-winres if not available
        $winresPath = & go env GOPATH
        $winresBin = Join-Path $winresPath "bin" "go-winres.exe"
        if (-not (Test-Path $winresBin)) {
            Write-Host "   Installing go-winres..." -ForegroundColor Gray
            go install github.com/tc-hib/go-winres@v0.3.3
            if ($LASTEXITCODE -ne 0) {
                Write-Host "⚠️  go-winres install failed — building without icon" -ForegroundColor Yellow
            }
        }
        if (Test-Path $winresBin) {
            & $winresBin init --icon $iconPath 2>$null
            Write-Host "✅ Icon embedded" -ForegroundColor Green
        }
    } else {
        Write-Host "⚠️  Icon not found at $iconPath — building without icon" -ForegroundColor Yellow
    }
}

go build -ldflags="-s -w" -o $BinaryFile .
if ($LASTEXITCODE -ne 0) {
    Write-ErrorAndExit "Build failed!"
}

# Clean up winres artifacts
if ($IsWindows_) {
    Remove-Item -Path "*.syso" -Force -ErrorAction SilentlyContinue
    Remove-Item -Path "winres" -Recurse -Force -ErrorAction SilentlyContinue
}
Write-Host "✅ Built: $BinaryFile" -ForegroundColor Green
Write-Host ""

# Step 4: Deploy to bin directory
Write-Host "📁 Deploying to: $DestDir" -ForegroundColor Yellow

# Create destination directory if it doesn't exist
if (-not (Test-Path $DestDir)) {
    New-Item -ItemType Directory -Path $DestDir -Force | Out-Null
    Write-Host "   Created directory: $DestDir" -ForegroundColor Gray
}

$DestPath = Join-Path $DestDir $BinaryFile
Copy-Item -Path $BinaryFile -Destination $DestPath -Force

# On Mac/Linux, ensure it's executable
if (-not $IsWindows_) {
    chmod +x $DestPath
}

Write-Host "✅ Deployed to: $DestPath" -ForegroundColor Green
Write-Host ""

# Step 5: Clean up local binary
try {
    Remove-Item -Path $BinaryFile -Force
    Write-Host "🧹 Cleaned local build artifact" -ForegroundColor Gray
} catch {
    Write-Host "⚠️  Could not remove local build artifact '$BinaryFile': $_" -ForegroundColor Yellow
}

# Step 6: Verify
Write-Host ""
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "  ✅ Done! Run 'movie' to use." -ForegroundColor Green
Write-Host "  📍 Binary at: $DestPath" -ForegroundColor Gray
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host ""

# Verify it runs
$prevPref = $ErrorActionPreference
$ErrorActionPreference = "Continue"
$verifyOutput = & $DestPath version 2>&1
$verifyExit = $LASTEXITCODE
$ErrorActionPreference = $prevPref

if ($verifyExit -eq 0) {
    Write-Host "  🎉 Verified — movie-cli is ready!" -ForegroundColor Green
} else {
    Write-Host "  ⚠️  Binary deployed but 'movie version' failed (exit $verifyExit)." -ForegroundColor Yellow
    foreach ($line in $verifyOutput) { Write-Host "    $line" -ForegroundColor Yellow }
    Write-Host "  Make sure $DestDir is in your PATH." -ForegroundColor Yellow
}
