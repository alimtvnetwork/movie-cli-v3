<#
.SYNOPSIS
 Build, deploy, and run movie CLI from the repo root.
.DESCRIPTION
 Pulls latest code, resolves Go dependencies, builds the binary
 into ./bin, deploys to a target directory, and optionally runs
 movie with any arguments.
.EXAMPLES
 .\run.ps1                        # pull, build, deploy
 .\run.ps1 -NoPull                # skip git pull
 .\run.ps1 -ForcePull             # discard local changes + pull (no prompt)
 .\run.ps1 -NoDeploy              # skip deploy step
 .\run.ps1 -R scan                # build + scan parent folder
 .\run.ps1 -R scan D:\movies      # build + scan specific path
 .\run.ps1 -R help                # build + show help
 .\run.ps1 -NoPull -NoDeploy -R scan  # just build and scan
 .\run.ps1 -t                     # run all unit tests with reports
.NOTES
 Configuration is read from powershell.json.
 -R accepts ALL movie CLI arguments after it (scan, ls, move, help, flags, paths).
  If -R is used with no arguments, it defaults to: movie scan <parentDir>
 -t runs all Go unit tests and writes reports to data/unit-test-reports/.
 -ForcePull automatically discards local changes and removes untracked files
 before pulling. Useful for CI or unattended builds.
#>

[CmdletBinding(PositionalBinding=$false)]
param(
    [switch]$NoPull,
    [switch]$NoDeploy,
    [switch]$ForcePull,
    [string]$DeployPath = "",
    [switch]$R,
    [Alias("t")]
    [switch]$Test,
    [Parameter(ValueFromRemainingArguments=$true)]
    [string[]]$RunArgs
)

$ErrorActionPreference = "Stop"
$RepoRoot = Split-Path -Parent $MyInvocation.MyCommand.Path

# -- Logging helpers -------------------------------------------

function Write-Step {
    param([string]$Step, [string]$Message)
    Write-Host ""
    Write-Host " [$Step] " -ForegroundColor Magenta -NoNewline
    Write-Host $Message -ForegroundColor White
    Write-Host (" " + ("-" * 50)) -ForegroundColor DarkGray
}

function Write-Success {
    param([string]$Message)
    Write-Host "  OK " -ForegroundColor Green -NoNewline
    Write-Host $Message -ForegroundColor Green
}

function Write-Info {
    param([string]$Message)
    Write-Host "  -> " -ForegroundColor Cyan -NoNewline
    Write-Host $Message -ForegroundColor Gray
}

function Write-Warn {
    param([string]$Message)
    Write-Host "  !! " -ForegroundColor Yellow -NoNewline
    Write-Host $Message -ForegroundColor Yellow
}

function Write-Fail {
    param([string]$Message)
    Write-Host "  XX " -ForegroundColor Red -NoNewline
    Write-Host $Message -ForegroundColor Red
}

function Write-ErrorAndExit {
    param([string]$Message, [string]$Hint = "")
    Write-Fail $Message
    if ($Hint) { Write-Info $Hint }
    exit 1
}

# -- Banner ----------------------------------------------------

function Show-Banner {
    Write-Host ""
    Write-Host " +======================================+" -ForegroundColor DarkCyan
    Write-Host " | " -ForegroundColor DarkCyan -NoNewline
    Write-Host "movie-cli builder" -ForegroundColor Cyan -NoNewline
    Write-Host "                |" -ForegroundColor DarkCyan
    Write-Host " +======================================+" -ForegroundColor DarkCyan
    Write-Host ""
}

# -- Load config -----------------------------------------------

function Load-Config {
    $configPath = Join-Path $RepoRoot "powershell.json"
    if (Test-Path $configPath) {
        Write-Info "Config loaded from powershell.json"
        return Get-Content $configPath | ConvertFrom-Json
    }
    Write-Warn "No powershell.json found, using defaults"
    return @{
        deployPath  = "E:\bin-run"
        buildOutput = "./bin"
        binaryName  = "movie.exe"
        copyData    = $false
    }
}

# -- Ensure main branch ----------------------------------------

function Ensure-MainBranch {
    Push-Location $RepoRoot
    try {
        $prevPref = $ErrorActionPreference
        $ErrorActionPreference = "Continue"
        $currentBranch = (git rev-parse --abbrev-ref HEAD 2>&1).Trim()
        $ErrorActionPreference = $prevPref

        if ($currentBranch -ne "main") {
            Write-Warn "Currently on branch '$currentBranch', switching to main..."
            $ErrorActionPreference = "Continue"
            $checkoutOutput = git checkout main 2>&1
            $checkoutExit = $LASTEXITCODE
            $ErrorActionPreference = $prevPref

            if ($checkoutExit -ne 0) {
                Write-Fail "Failed to switch to main branch"
                foreach ($line in $checkoutOutput) {
                    Write-Host "  $line" -ForegroundColor Red
                }
                exit 1
            }
            Write-Success "Switched to main branch"
        }
    } finally {
        Pop-Location
    }
}

# -- Retry git pull after stash/discard -----------------------

function Retry-GitPull {
    Write-Info "Retrying git pull..."
    $prevPref = $ErrorActionPreference
    $ErrorActionPreference = "Continue"
    $retryOutput = git pull 2>&1
    $retryExit = $LASTEXITCODE
    $ErrorActionPreference = $prevPref

    foreach ($line in $retryOutput) {
        $text = "$line".Trim()
        if ($text.Length -gt 0) {
            Write-Info $text
        }
    }

    if ($retryExit -ne 0) {
        Write-ErrorAndExit "Git pull failed again (exit code $retryExit)"
    }

    Write-Success "Pull complete"
}

# -- Resolve pull conflict with local changes ------------------

function Resolve-PullConflict {
    Write-Warn "Git pull failed due to local changes"
    Write-Host ""
    Write-Host "  Choose how to proceed:" -ForegroundColor Yellow
    Write-Host "  [S] Stash changes (save for later, then pull)" -ForegroundColor Cyan
    Write-Host "  [D] Discard changes (reset working tree, then pull)" -ForegroundColor Cyan
    Write-Host "  [C] Clean all (discard changes + remove untracked files, then pull)" -ForegroundColor Cyan
    Write-Host "  [Q] Quit (abort without changes)" -ForegroundColor Cyan
    Write-Host ""

    $choice = Read-Host "  Enter choice (S/D/C/Q)"

    switch ($choice.ToUpper()) {
        "S" {
            Write-Info "Stashing local changes..."
            $prevPref = $ErrorActionPreference
            $ErrorActionPreference = "Continue"
            $stashOutput = git stash push -m "auto-stash before run.ps1 pull" 2>&1
            $stashExit = $LASTEXITCODE
            $ErrorActionPreference = $prevPref

            if ($stashExit -ne 0) {
                Write-Fail "Git stash failed"
                foreach ($line in $stashOutput) {
                    Write-Host "  $line" -ForegroundColor Red
                }
                exit 1
            }
            Write-Success "Changes stashed"
            Write-Info "Run 'git stash pop' later to restore your changes"

            Retry-GitPull
        }
        "D" {
            Write-Warn "Discarding all local changes..."
            $prevPref = $ErrorActionPreference
            $ErrorActionPreference = "Continue"
            $resetOutput = git checkout -- . 2>&1
            $resetExit = $LASTEXITCODE
            $ErrorActionPreference = $prevPref

            if ($resetExit -ne 0) {
                Write-Fail "Git checkout failed"
                foreach ($line in $resetOutput) {
                    Write-Host "  $line" -ForegroundColor Red
                }
                exit 1
            }
            Write-Success "Local changes discarded"

            Retry-GitPull
        }
        "C" {
            Write-Warn "Discarding all local changes and removing untracked files..."
            $prevPref = $ErrorActionPreference
            $ErrorActionPreference = "Continue"

            $resetOutput = git checkout -- . 2>&1
            $resetExit = $LASTEXITCODE

            if ($resetExit -ne 0) {
                Write-Fail "Git checkout failed"
                foreach ($line in $resetOutput) {
                    Write-Host "  $line" -ForegroundColor Red
                }
                $ErrorActionPreference = $prevPref
                exit 1
            }
            Write-Success "Local changes discarded"

            $cleanOutput = git clean -fd 2>&1
            $cleanExit = $LASTEXITCODE
            $ErrorActionPreference = $prevPref

            if ($cleanExit -ne 0) {
                Write-Fail "Git clean failed"
                foreach ($line in $cleanOutput) {
                    Write-Host "  $line" -ForegroundColor Red
                }
                exit 1
            }

            $cleanedFiles = @($cleanOutput | ForEach-Object { "$_".Trim() } | Where-Object { $_.Length -gt 0 })
            if ($cleanedFiles.Count -gt 0) {
                foreach ($line in $cleanedFiles) {
                    Write-Info $line
                }
                Write-Success "Removed $($cleanedFiles.Count) untracked file(s)"
            } else {
                Write-Info "No untracked files to remove"
            }

            Retry-GitPull
        }
        default {
            Write-Info "Aborted by user"
            exit 0
        }
    }
}

# -- Git pull --------------------------------------------------

function Invoke-GitPull {
    Write-Step "1/4" "Pulling latest changes"

    Ensure-MainBranch

    Push-Location $RepoRoot
    try {
        $prevPref = $ErrorActionPreference
        $ErrorActionPreference = "Continue"
        $output = git pull 2>&1
        $pullExit = $LASTEXITCODE
        $ErrorActionPreference = $prevPref

        foreach ($line in $output) {
            $text = "$line".Trim()
            if ($text.Length -gt 0) {
                Write-Info $text
            }
        }

        if ($pullExit -ne 0) {
            $outputText = ($output | ForEach-Object { "$_" }) -join "`n"
            $hasConflict = $outputText -match "Your local changes" -or
                           $outputText -match "overwritten by merge" -or
                           $outputText -match "not possible because you have unmerged" -or
                           $outputText -match "Please commit your changes or stash them"

            if ($hasConflict) {
                if ($ForcePull) {
                    Write-Warn "Force-pull: discarding local changes and removing untracked files..."
                    $prevPref = $ErrorActionPreference
                    $ErrorActionPreference = "Continue"

                    $resetOutput = git checkout -- . 2>&1
                    $resetExit = $LASTEXITCODE
                    if ($resetExit -ne 0) {
                        Write-Fail "Git checkout failed"
                        $ErrorActionPreference = $prevPref
                        Write-ErrorAndExit "Git checkout failed"
                    }
                    Write-Success "Local changes discarded"

                    $cleanOutput = git clean -fd 2>&1
                    $cleanExit = $LASTEXITCODE
                    $ErrorActionPreference = $prevPref

                    if ($cleanExit -ne 0) {
                        Write-ErrorAndExit "Git clean failed"
                    }

                    $cleanedFiles = @($cleanOutput | ForEach-Object { "$_".Trim() } | Where-Object { $_.Length -gt 0 })
                    if ($cleanedFiles.Count -gt 0) {
                        Write-Success "Removed $($cleanedFiles.Count) untracked file(s)"
                    }

                    Retry-GitPull
                } else {
                    Resolve-PullConflict
                }
            } else {
                Write-ErrorAndExit "Git pull failed (exit code $pullExit)"
            }
        } else {
            Write-Success "Pull complete"
        }
    } finally {
        Pop-Location
    }
}

# -- Resolve dependencies --------------------------------------

function Resolve-Dependencies {
    Write-Step "2/4" "Resolving Go dependencies"

    Push-Location $RepoRoot
    try {
        # Verify Go is installed
        $prevPref = $ErrorActionPreference
        $ErrorActionPreference = "Continue"
        $goVersion = go version 2>&1
        $goExit = $LASTEXITCODE
        $ErrorActionPreference = $prevPref

        if ($goExit -ne 0) {
            Write-ErrorAndExit "Go is not installed or not in PATH" -Hint "Install Go from https://go.dev/dl/"
        }
        Write-Info "$goVersion"

        # Tidy modules
        Write-Info "Running go mod tidy..."
        $prevPref = $ErrorActionPreference
        $ErrorActionPreference = "Continue"
        $tidyOutput = go mod tidy 2>&1
        $tidyExit = $LASTEXITCODE
        $ErrorActionPreference = $prevPref

        if ($tidyExit -ne 0) {
            Write-Fail "go mod tidy failed"
            foreach ($line in $tidyOutput) {
                Write-Host "  $line" -ForegroundColor Red
            }
            exit 1
        }
        Write-Success "Dependencies resolved"
    } finally {
        Pop-Location
    }
}

# -- Build binary ----------------------------------------------

function Build-Binary {
    Write-Step "3/4" "Building binary"

    Push-Location $RepoRoot
    try {
        $binaryName = $config.binaryName
    if (-not $binaryName) { $binaryName = "movie.exe" }
    $buildOutput = $config.buildOutput
        if (-not $buildOutput) { $buildOutput = "./bin" }

        # Create build output directory
        $buildDir = Join-Path $RepoRoot ($buildOutput -replace '^\.\/', '')
        if (-not (Test-Path $buildDir)) {
            New-Item -ItemType Directory -Path $buildDir -Force | Out-Null
            Write-Info "Created build directory: $buildDir"
        }

        $outputPath = Join-Path $buildDir $binaryName

        # Detect OS and architecture
        $IsWindows_ = ($PSVersionTable.PSEdition -eq "Desktop") -or ($IsWindows -eq $true)
        $IsMac_ = ($IsMacOS -eq $true)

        if ($IsWindows_) {
            $env:GOOS = "windows"
            $env:GOARCH = "amd64"
        } elseif ($IsMac_) {
            $env:GOOS = "darwin"
            $prevPref = $ErrorActionPreference
            $ErrorActionPreference = "Continue"
            $arch = & uname -m 2>&1
            $ErrorActionPreference = $prevPref
            if ("$arch".Trim() -eq "arm64") {
                $env:GOARCH = "arm64"
            } else {
                $env:GOARCH = "amd64"
            }
        } else {
            $env:GOOS = "linux"
            $env:GOARCH = "amd64"
        }

        Write-Info "Target: $($env:GOOS)/$($env:GOARCH)"

        # Gather ldflags: version, commit, build date
        $prevPref = $ErrorActionPreference
        $ErrorActionPreference = "Continue"

        $gitCommit = (git rev-parse --short HEAD 2>&1)
        if ($LASTEXITCODE -ne 0) { $gitCommit = "unknown" }
        $gitCommit = "$gitCommit".Trim()

        $buildDate = (Get-Date -Format "yyyy-MM-ddTHH:mm:ssK")

        $sourceVersion = "dev"
        $versionFile = Join-Path $RepoRoot "version" "version.go"
        if (Test-Path $versionFile) {
            $versionMatch = Select-String -Path $versionFile -Pattern 'Version\s*=\s*"([^"]+)"' -AllMatches
            if ($versionMatch -and $versionMatch.Matches.Count -gt 0) {
                $sourceVersion = $versionMatch.Matches[0].Groups[1].Value.Trim()
            }
        }

        $ErrorActionPreference = $prevPref

        $ldflags = "-s -w -X github.com/alimtvnetwork/movie-cli-v3/version.Commit=$gitCommit -X github.com/alimtvnetwork/movie-cli-v3/version.BuildDate=$buildDate"
        Write-Info "Version: $sourceVersion | Commit: $gitCommit"

        # Build
        Write-Info "Compiling to $outputPath..."
        $prevPref = $ErrorActionPreference
        $ErrorActionPreference = "Continue"
        $buildOutput2 = go build -ldflags="$ldflags" -o $outputPath . 2>&1
        $buildExit = $LASTEXITCODE
        $ErrorActionPreference = $prevPref

        if ($buildExit -ne 0) {
            Write-Fail "Build failed"
            foreach ($line in $buildOutput2) {
                Write-Host "  $line" -ForegroundColor Red
            }
            exit 1
        }

        # Verify binary exists and get size
        if (Test-Path $outputPath) {
            $fileInfo = Get-Item $outputPath
            $sizeMB = [math]::Round($fileInfo.Length / 1MB, 2)
            Write-Success "Built: $binaryName ($sizeMB MB)"
        } else {
            Write-ErrorAndExit "Binary not found after build: $outputPath"
        }

        # Post-build verification: run version command
        $prevPref = $ErrorActionPreference
        $ErrorActionPreference = "Continue"
        $verifyOutput = & $outputPath version 2>&1
        $verifyExit = $LASTEXITCODE
        $ErrorActionPreference = $prevPref

        if ($verifyExit -eq 0) {
            Write-Success "Verified: $("$verifyOutput".Trim())"
        } else {
            Write-Warn "Binary built but 'version' command failed — may be OK if not implemented yet"
        }

        return $outputPath
    } finally {
        Pop-Location
    }
}

# -- Deploy binary ---------------------------------------------

function Deploy-Binary {
    param([string]$SourceBinary)

    Write-Step "4/4" "Deploying binary"

    $deployPath = $config.deployPath
    if (-not $deployPath) {
        $deployPath = if (($PSVersionTable.PSEdition -eq "Desktop") -or ($IsWindows -eq $true)) {
            "E:\bin-run"
        } else {
            "/usr/local/bin"
        }
    }

    $binaryName = $config.binaryName
    if (-not $binaryName) { $binaryName = "movie.exe" }

    # Create deploy directory if needed
    if (-not (Test-Path $deployPath)) {
        New-Item -ItemType Directory -Path $deployPath -Force | Out-Null
        Write-Info "Created deploy directory: $deployPath"
    }

    $destFile = Join-Path $deployPath $binaryName
    $backupFile = "$destFile.bak"

    # Safe deploy: rename existing binary first (rollback on failure)
    $hadExisting = Test-Path $destFile
    if ($hadExisting) {
        try {
            Rename-Item -Path $destFile -NewName "$binaryName.bak" -Force
            Write-Info "Backed up existing binary"
        } catch {
            Write-ErrorAndExit "Could not back up existing binary: $_ — cannot guarantee rollback safety"
        }
    }

    try {
        Copy-Item -Path $SourceBinary -Destination $destFile -Force
        Write-Success "Deployed to: $destFile"
    } catch {
        Write-Fail "Deploy failed: $_"
        # Rollback
        if ($hadExisting -and (Test-Path $backupFile)) {
            Rename-Item -Path $backupFile -NewName $binaryName -Force
            Write-Warn "Rolled back to previous binary"
        }
        exit 1
    }

    # On Mac/Linux, ensure it's executable
    $IsWindows_ = ($PSVersionTable.PSEdition -eq "Desktop") -or ($IsWindows -eq $true)
    if (-not $IsWindows_) {
        $prevPref = $ErrorActionPreference
        $ErrorActionPreference = "Continue"
        chmod +x $destFile 2>&1 | Out-Null
        $chmodExit = $LASTEXITCODE
        $ErrorActionPreference = $prevPref
        if ($chmodExit -ne 0) {
            Write-Warn "chmod +x failed on $destFile (exit $chmodExit) — binary may not be executable"
        }
    }

    # Copy changelog next to the deployed binary so `movie changelog` works outside the repo root
    $repoChangelog = Join-Path $RepoRoot "CHANGELOG.md"
    $deployChangelog = Join-Path $deployPath "CHANGELOG.md"
    if (Test-Path $repoChangelog) {
        try {
            Copy-Item -Path $repoChangelog -Destination $deployChangelog -Force
            Write-Info "Copied CHANGELOG.md to deploy directory"
        } catch {
            Write-Warn "Could not copy CHANGELOG.md to deploy directory: $_"
        }
    } else {
        Write-Warn "CHANGELOG.md not found in repo root"
    }

    # Clean up backup on success
    if (Test-Path $backupFile) {
        try {
            Remove-Item -Path $backupFile -Force
            Write-Info "Cleaned up backup"
        } catch {
            Write-Warn "Could not remove backup file '$backupFile': $_"
        }
    }

    # Verify deployed binary
    $prevPref = $ErrorActionPreference
    $ErrorActionPreference = "Continue"
    $verifyOutput = & $destFile version 2>&1
    $verifyExit = $LASTEXITCODE
    $ErrorActionPreference = $prevPref

    if ($verifyExit -eq 0) {
        Write-Success "Verified deployed binary: $("$verifyOutput".Trim())"
    } else {
        Write-Warn "Deployed but 'version' check failed — may be OK"
    }

    if (Test-Path $deployChangelog) {
        $prevPref = $ErrorActionPreference
        $ErrorActionPreference = "Continue"
        $changelogVerify = & $destFile changelog --latest 2>&1
        $changelogExit = $LASTEXITCODE
        $ErrorActionPreference = $prevPref

        if ($changelogExit -eq 0) {
            Write-Success "Verified deployed changelog"
        } else {
            Write-Warn "CHANGELOG.md was deployed, but 'changelog --latest' failed"
        }
    }

    # PATH check
    $pathDirs = $env:PATH -split [IO.Path]::PathSeparator
    $inPath = $pathDirs | Where-Object { $_ -eq $deployPath }
    if (-not $inPath) {
        Write-Warn "Deploy path is NOT in PATH: $deployPath"
        Write-Info "Add it to PATH to run '$($binaryName -replace '\.exe$','')' from anywhere"
    } else {
        Write-Success "Deploy path is in PATH"
    }

    $resolvedMovie = Get-Command movie -ErrorAction SilentlyContinue
    if ($resolvedMovie -and $resolvedMovie.Source -and (Test-Path $resolvedMovie.Source)) {
        $pathMovie = (Resolve-Path $resolvedMovie.Source).Path
        $deployMovie = (Resolve-Path $destFile).Path
        if ($pathMovie -ne $deployMovie) {
            Write-Warn "PATH resolves 'movie' to a different binary: $pathMovie"
            Write-Info "Run '$destFile version' or move $deployPath earlier in PATH to use the newly deployed build"
        }
    }

    # Copy data directory if configured
    $copyData = $config.copyData
    if ($copyData -eq $true) {
        $dataSource = Join-Path $RepoRoot "data"
        if (Test-Path $dataSource) {
            $dataDest = Join-Path $deployPath "data"
            Copy-Item -Path $dataSource -Destination $dataDest -Recurse -Force
            Write-Info "Copied data directory to $dataDest"
        }
    }

    return $destFile
}

# -- Run mode --------------------------------------------------

function Invoke-Run {
    param([string]$BinaryPath, [string[]]$Arguments)

    $parentDir = Split-Path -Parent $RepoRoot

    # Default: if no arguments, run "movie scan <parentDir>"
    if (-not $Arguments -or $Arguments.Count -eq 0) {
        $Arguments = @("movie", "scan", $parentDir)
        Write-Info "No arguments provided, defaulting to: movie scan $parentDir"
    }

    # Resolve relative paths in arguments
    $resolvedArgs = @()
    foreach ($arg in $Arguments) {
        if ($arg -match '^[a-zA-Z]' -and (Test-Path $arg -ErrorAction SilentlyContinue)) {
            $resolvedArgs += (Resolve-Path $arg).Path
        } else {
            $resolvedArgs += $arg
        }
    }

    Write-Info "Running: $BinaryPath $($resolvedArgs -join ' ')"
    Write-Host (" " + ("-" * 50)) -ForegroundColor DarkGray

    & $BinaryPath @resolvedArgs
    $runExit = $LASTEXITCODE

    Write-Host (" " + ("-" * 50)) -ForegroundColor DarkGray
    if ($runExit -eq 0) {
        Write-Success "Run complete (exit 0)"
    } else {
        Write-Fail "Run exited with code $runExit"
    }
}

# -- Test mode -------------------------------------------------

function Invoke-Tests {
    Write-Step "TEST" "Running all Go unit tests"

    Push-Location $RepoRoot
    try {
        $reportDir = Join-Path $RepoRoot "data" "unit-test-reports"
        if (-not (Test-Path $reportDir)) {
            New-Item -ItemType Directory -Path $reportDir -Force | Out-Null
        }

        $timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
        $reportFile = Join-Path $reportDir "test-report-$timestamp.txt"

        Write-Info "Output: $reportFile"

        $prevPref = $ErrorActionPreference
        $ErrorActionPreference = "Continue"
        $testOutput = go test ./... -v 2>&1
        $testExit = $LASTEXITCODE
        $ErrorActionPreference = $prevPref

        # Write report
        $testOutput | Out-File -FilePath $reportFile -Encoding utf8

        # Display summary
        $passed = @($testOutput | Where-Object { "$_" -match "^--- PASS:" }).Count
        $failed = @($testOutput | Where-Object { "$_" -match "^--- FAIL:" }).Count
        $skipped = @($testOutput | Where-Object { "$_" -match "^--- SKIP:" }).Count

        Write-Info "Passed: $passed | Failed: $failed | Skipped: $skipped"

        if ($testExit -eq 0) {
            Write-Success "All tests passed"
        } else {
            Write-Fail "Some tests failed (exit code $testExit)"
            # Show failed test lines
            $failLines = @($testOutput | Where-Object { "$_" -match "FAIL" })
            foreach ($line in $failLines) {
                Write-Host "  $line" -ForegroundColor Red
            }
        }

        Write-Info "Full report: $reportFile"
    } finally {
        Pop-Location
    }
}

# -- Main ------------------------------------------------------

Show-Banner
$config = Load-Config

# Test mode — run tests and exit
if ($Test) {
    Invoke-Tests
    exit 0
}

if (-not $NoPull) {
    Invoke-GitPull
} else {
    Write-Info "Skipping git pull (-NoPull)"
}

Resolve-Dependencies

$builtBinary = Build-Binary

$deployedBinary = $null
if (-not $NoDeploy) {
    $deployedBinary = Deploy-Binary -SourceBinary $builtBinary
} else {
    Write-Info "Skipping deploy (-NoDeploy)"
    $deployedBinary = $builtBinary
}

# Show the version from the binary that was just built/deployed
$activeBinary = $deployedBinary
if (-not $activeBinary -or -not (Test-Path $activeBinary)) {
    $activeCmd = Get-Command movie -ErrorAction SilentlyContinue
    if ($activeCmd -and $activeCmd.Source -and (Test-Path $activeCmd.Source)) {
        $activeBinary = $activeCmd.Source
    }
}

if ($activeBinary -and (Test-Path $activeBinary)) {
    $prevPref = $ErrorActionPreference
    $ErrorActionPreference = "Continue"
    $finalVersion = & $activeBinary version 2>&1
    $ErrorActionPreference = $prevPref
    $finalVersionText = "$finalVersion".Trim()
} else {
    $finalVersionText = "unknown"
}

# Run mode
if ($R) {
    Write-Host ""
    Invoke-Run -BinaryPath $deployedBinary -Arguments $RunArgs
}

# -- Summary ---------------------------------------------------

# Parse version output into separate lines
$versionLines = @()
if ($finalVersionText -and $finalVersionText -ne "unknown") {
    $rawLines = $finalVersionText -split "`n" | ForEach-Object { $_.Trim() } | Where-Object { $_.Length -gt 0 }
    foreach ($l in $rawLines) { $versionLines += $l }
} else {
    $versionLines += "unknown"
}

$boxWidth = 40
$innerWidth = $boxWidth - 4  # account for " | " and " |"

Write-Host ""
Write-Host " +$('=' * ($boxWidth - 2))+" -ForegroundColor DarkCyan
# Title line
$title = "All done!"
$titlePad = $innerWidth - $title.Length
if ($titlePad -lt 0) { $titlePad = 0 }
Write-Host " | " -ForegroundColor DarkCyan -NoNewline
Write-Host $title -ForegroundColor Green -NoNewline
Write-Host "$(' ' * $titlePad) |" -ForegroundColor DarkCyan
Write-Host " +$('-' * ($boxWidth - 2))+" -ForegroundColor DarkCyan

# Version info lines
foreach ($vLine in $versionLines) {
    $displayLine = $vLine
    if ($displayLine.Length -gt $innerWidth) {
        $displayLine = $displayLine.Substring(0, $innerWidth)
    }
    $linePad = $innerWidth - $displayLine.Length
    if ($linePad -lt 0) { $linePad = 0 }
    Write-Host " | " -ForegroundColor DarkCyan -NoNewline
    Write-Host $displayLine -ForegroundColor Cyan -NoNewline
    Write-Host "$(' ' * $linePad) |" -ForegroundColor DarkCyan
}

Write-Host " +$('=' * ($boxWidth - 2))+" -ForegroundColor DarkCyan
Write-Host ""

# -- Latest changelog ------------------------------------------
$changelogPath = Join-Path $RepoRoot "CHANGELOG.md"
if (Test-Path $changelogPath) {
    $clLines = Get-Content $changelogPath -Encoding UTF8
    $latestBlock = @()
    $inBlock = $false
    $blockCount = 0
    foreach ($cl in $clLines) {
        if ($cl -match '^## ') {
            $blockCount++
            if ($blockCount -eq 1) {
                $inBlock = $true
                $latestBlock += $cl
                continue
            } else {
                break
            }
        }
        if ($inBlock) { $latestBlock += $cl }
    }
    if ($latestBlock.Count -gt 0) {
        Write-Host ""
        Write-Host "  Latest changelog:" -ForegroundColor DarkCyan
        Write-Host "  --------------------------------------------------" -ForegroundColor DarkGray
        foreach ($cl in $latestBlock) {
            if ($cl -match '^## ') {
                Write-Host "  $cl" -ForegroundColor Cyan
            } elseif ($cl -match '^### ') {
                Write-Host "  $cl" -ForegroundColor Yellow
            } elseif ($cl -match '^- ') {
                Write-Host "  $cl" -ForegroundColor Gray
            } elseif ($cl.Trim().Length -gt 0) {
                Write-Host "  $cl" -ForegroundColor DarkGray
            }
        }
        Write-Host ""
    }
}
