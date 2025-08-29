# ATA Unified Build Script
# Usage:
#   ./build.ps1 -dev     # Development build (no version increment)
#   ./build.ps1 -release # Release build (increment version)

param(
    [switch]$dev,
    [switch]$release
)

# Validate parameters
if (-not $dev -and -not $release) {
    Write-Host "Usage: ./build.ps1 [-dev | -release]" -ForegroundColor Red
    Write-Host "  -dev     Development build (no version increment)" -ForegroundColor Gray
    Write-Host "  -release Release build (increment version)" -ForegroundColor Gray
    exit 1
}

if ($dev -and $release) {
    Write-Host "Error: Cannot specify both -dev and -release" -ForegroundColor Red
    exit 1
}

# Version management
$versionFile = "version.txt"
if (!(Test-Path $versionFile)) {
    Write-Host "Creating initial version file..." -ForegroundColor Yellow
    "0.0.1" | Out-File -FilePath $versionFile -Encoding UTF8
}

# Read current version
$currentVersion = Get-Content $versionFile -Raw
$currentVersion = $currentVersion.Trim()

if ($dev) {
    # Development build
    Write-Host "Building ATA for development..." -ForegroundColor Green
    $buildVersion = "$currentVersion-dev"
    $outputFile = "ata-installer-dev.exe"
    $platforms = @("windows")
    $confirmRequired = $false
} else {
    # Release build
    Write-Host "Building ATA RELEASE version..." -ForegroundColor Red
    Write-Host "This will increment the version number and create official release builds." -ForegroundColor Yellow
    Write-Host ""
    $confirm = Read-Host "Are you sure you want to create a RELEASE build? (y/N)"
    if ($confirm -ne "y" -and $confirm -ne "Y") {
        Write-Host "Release build cancelled." -ForegroundColor Yellow
        exit 0
    }
    
    # Increment patch version
    if ($currentVersion -match '^(\d+)\.(\d+)\.(\d+)$') {
        $major = [int]$matches[1]
        $minor = [int]$matches[2]
        $patch = [int]$matches[3] + 1
        $buildVersion = "$major.$minor.$patch"
        
        # Update version file
        $buildVersion | Out-File -FilePath $versionFile -Encoding UTF8
        Write-Host "Version updated: $currentVersion -> $buildVersion" -ForegroundColor Green
    } else {
        Write-Host "Invalid version format in version.txt: $currentVersion" -ForegroundColor Red
        Read-Host "Press Enter to exit"
        exit 1
    }
    
    $platforms = @("windows", "linux", "macos")
    $confirmRequired = $true
}

Write-Host ""
Write-Host "Build version: $buildVersion" -ForegroundColor Cyan
Write-Host ""

# Clean old build files
if (Test-Path "ata.exe") { Remove-Item "ata.exe" }
if ($dev) {
    if (Test-Path "ata-installer-dev.exe") { Remove-Item "ata-installer-dev.exe" }
} else {
    if (!(Test-Path "release")) { New-Item -ItemType Directory -Name "release" }
    Get-ChildItem "release" | Remove-Item -Force
}

# Build function
function Build-Platform {
    param(
        [string]$Platform,
        [string]$Version,
        [bool]$IsDev
    )
    
    $platformName = $Platform
    $exeSuffix = if ($Platform -eq "windows") { ".exe" } else { "" }
    $mainExe = if ($Platform -eq "windows") { "ata.exe" } else { "ata" }
    
    Write-Host "Building for $platformName..." -ForegroundColor Yellow
    
    # Set environment variables
    $env:GOOS = if ($Platform -eq "macos") { "darwin" } else { $Platform }
    $env:GOARCH = "amd64"
    
    # Build flags
    $ldflags = if ($IsDev) { "-X main.version=$Version" } else { "-s -w -X main.version=$Version" }
    
    # Build main program first
    Write-Host "  Building main program..." -ForegroundColor Cyan
    go build -ldflags $ldflags -o "cmd/setup/$mainExe" "./cmd/ata"
    if ($LASTEXITCODE -ne 0) {
        Write-Host "$platformName main program build failed!" -ForegroundColor Red
        Read-Host "Press Enter to exit"
        exit 1
    }
    
    # Build installer
    Write-Host "  Building installer..." -ForegroundColor Cyan
    if ($IsDev) {
        $installerName = "ata-installer-dev$exeSuffix"
    } else {
        $installerName = "ata-installer-$platformName-$Version$exeSuffix"
    }
    
    go build -ldflags $ldflags -o $installerName "./cmd/setup"
    if ($LASTEXITCODE -ne 0) {
        Write-Host "$platformName installer build failed!" -ForegroundColor Red
        Read-Host "Press Enter to exit"
        exit 1
    }
    
    # Move to release directory for release builds
    if (-not $IsDev) {
        Move-Item $installerName "release/"
    }
    
    # Clean temporary files
    if (Test-Path "cmd/setup/$mainExe") { Remove-Item "cmd/setup/$mainExe" }
}

# Build for all required platforms
foreach ($platform in $platforms) {
    Build-Platform -Platform $platform -Version $buildVersion -IsDev $dev
}

# Reset environment variables
$env:GOOS = ""
$env:GOARCH = ""

Write-Host ""
Write-Host "Build completed successfully!" -ForegroundColor Green

if ($dev) {
    Write-Host "Development installer: ata-installer-dev.exe" -ForegroundColor Cyan
    Write-Host "Version: $buildVersion" -ForegroundColor White
} else {
    Write-Host "Release files:" -ForegroundColor Cyan
    Get-ChildItem "release" | ForEach-Object { Write-Host "  $($_.Name)" -ForegroundColor White }
    Write-Host ""
    Write-Host "Installation:" -ForegroundColor Yellow
    Write-Host "  Windows: run ata-installer-windows-$buildVersion.exe" -ForegroundColor White
    Write-Host "  Linux:   chmod +x ata-installer-linux-$buildVersion && ./ata-installer-linux-$buildVersion" -ForegroundColor White
    Write-Host "  macOS:   chmod +x ata-installer-macos-$buildVersion && ./ata-installer-macos-$buildVersion" -ForegroundColor White
}
Write-Host ""