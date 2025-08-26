# PowerShell installation script for ATA
# This script adds ATA to PowerShell profile to avoid PATH length limitations

Write-Host "Installing ATA for PowerShell..." -ForegroundColor Green

# Check if ata.exe exists
$ataPath = Join-Path $PSScriptRoot "bin\ata.exe"
if (-not (Test-Path $ataPath)) {
    Write-Host "Error: ata.exe not found at $ataPath" -ForegroundColor Red
    Write-Host "Please build the project first, then run this install script" -ForegroundColor Red
    exit 1
}

# Create bin directory if it doesn't exist
$binDir = Join-Path $env:USERPROFILE "bin"
if (-not (Test-Path $binDir)) {
    New-Item -ItemType Directory -Path $binDir -Force | Out-Null
    Write-Host "Created $binDir directory" -ForegroundColor Yellow
}

# Copy ata.exe to bin directory
Copy-Item $ataPath $binDir -Force
Write-Host "Copied ata.exe to $binDir" -ForegroundColor Green

# Get PowerShell profile path
$profilePath = $PROFILE.CurrentUserCurrentHost
$profileDir = Split-Path $profilePath -Parent

# Create profile directory if it doesn't exist
if (-not (Test-Path $profileDir)) {
    New-Item -ItemType Directory -Path $profileDir -Force | Out-Null
    Write-Host "Created PowerShell profile directory: $profileDir" -ForegroundColor Yellow
}

# Check if profile exists and contains ata alias
$aliasLine = "Set-Alias -Name ata -Value '$binDir\ata.exe'"
$profileExists = Test-Path $profilePath
$aliasExists = $false

if ($profileExists) {
    $profileContent = Get-Content $profilePath -Raw
    $aliasExists = $profileContent -match "Set-Alias.*ata.*ata\.exe"
}

# Add to PATH environment variable (more universal solution)
$currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($currentPath -notlike "*$binDir*") {
    $newPath = "$currentPath;$binDir"
    [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
    Write-Host "Added $binDir to user PATH environment variable" -ForegroundColor Green
    Write-Host "Note: You may need to restart your terminal for PATH changes to take effect" -ForegroundColor Yellow
} else {
    Write-Host "$binDir already exists in user PATH" -ForegroundColor Yellow
}

# Also add PowerShell alias as fallback
if (-not $aliasExists) {
    if ($profileExists) {
        Add-Content $profilePath "`n# ATA alias (fallback)"
        Add-Content $profilePath $aliasLine
    } else {
        Set-Content $profilePath "# ATA alias (fallback)`n$aliasLine"
    }
    Write-Host "Added ata alias to PowerShell profile as fallback: $profilePath" -ForegroundColor Green
} else {
    Write-Host "ATA alias already exists in PowerShell profile" -ForegroundColor Yellow
}

Write-Host "`nInstallation complete!" -ForegroundColor Green
Write-Host "ATA has been added to your PATH environment variable." -ForegroundColor Cyan
Write-Host "You can now use 'ata' command in any terminal (CMD, PowerShell, Git Bash, etc.)" -ForegroundColor Cyan
Write-Host "If the command is not recognized immediately, please restart your terminal." -ForegroundColor Yellow
Write-Host "`nTest the installation by running: ata help" -ForegroundColor White