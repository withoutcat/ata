# ATA Unified Build Script
# Usage:
#   ./build.ps1          # Development build (default)
#   ./build.ps1 -dev     # Development build (no version increment)
#   ./build.ps1 -release # Release build (increment version)

param(
    [switch]$dev,
    [switch]$release
)

# Validate parameters
if ($dev -and $release) {
    Write-Host "Error: Cannot specify both -dev and -release" -ForegroundColor Red
    exit 1
}

# Default to development build if no parameters specified
if (-not $dev -and -not $release) {
    $dev = $true
    Write-Host "No build type specified, defaulting to development build..." -ForegroundColor Yellow
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

# Create and clean build directories
if ($dev) {
    $buildDir = Join-Path "build" "dev"
} else {
    $buildDir = Join-Path "build" "release"
}

if (!(Test-Path $buildDir)) { 
    New-Item -ItemType Directory -Path $buildDir -Force | Out-Null
}
Get-ChildItem $buildDir -ErrorAction SilentlyContinue | Remove-Item -Force -Recurse

# Clean old build files from root directory (legacy cleanup)
if (Test-Path "ata.exe") { Remove-Item "ata.exe" }
if (Test-Path "ata-installer-dev.exe") { Remove-Item "ata-installer-dev.exe" }
if (Test-Path "release") { Remove-Item "release" -Recurse -Force }

# Build function
function Build-Platform {
    param(
        [string]$Platform,
        [string]$Version,
        [bool]$IsDev,
        [string]$BuildDir
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
    
    $installerPath = Join-Path $BuildDir $installerName
    go build -ldflags $ldflags -o $installerPath "./cmd/setup"
    if ($LASTEXITCODE -ne 0) {
        Write-Host "$platformName installer build failed!" -ForegroundColor Red
        Read-Host "Press Enter to exit"
        exit 1
    }
    
    # Clean temporary files
    if (Test-Path "cmd/setup/$mainExe") { Remove-Item "cmd/setup/$mainExe" }
}

# Build for all required platforms
foreach ($platform in $platforms) {
    Build-Platform -Platform $platform -Version $buildVersion -IsDev $dev -BuildDir $buildDir
}

# Reset environment variables
$env:GOOS = ""
$env:GOARCH = ""

Write-Host ""
Write-Host "Build completed successfully!" -ForegroundColor Green

if ($dev) {
    Write-Host "Development build directory: $buildDir" -ForegroundColor Cyan
    $devInstallerPath = Join-Path $buildDir "ata-installer-dev.exe"
    Write-Host "Development installer: $devInstallerPath" -ForegroundColor Cyan
    Write-Host "Version: $buildVersion" -ForegroundColor White
} else {
    Write-Host "Release build directory: $buildDir" -ForegroundColor Cyan
    Write-Host "Release files:" -ForegroundColor Cyan
    Get-ChildItem $buildDir | ForEach-Object { Write-Host "  $($_.Name)" -ForegroundColor White }
    Write-Host ""
    Write-Host "Installation:" -ForegroundColor Yellow
    $winInstaller = Join-Path $buildDir "ata-installer-windows-$buildVersion.exe"
    $linuxInstaller = Join-Path $buildDir "ata-installer-linux-$buildVersion"
    $macInstaller = Join-Path $buildDir "ata-installer-macos-$buildVersion"
    Write-Host "  Windows: run $winInstaller" -ForegroundColor White
    Write-Host "  Linux:   chmod +x $linuxInstaller && ./$linuxInstaller" -ForegroundColor White
    Write-Host "  macOS:   chmod +x $macInstaller && ./$macInstaller" -ForegroundColor White
}
Write-Host ""