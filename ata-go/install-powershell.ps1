# PowerShell installation script for ATA
# This script copies ATA to user bin directory and adds it to PATH

Write-Host "Installing ATA..." -ForegroundColor Green

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

# No PowerShell profile setup needed - using PATH only

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

# PATH-based installation complete

Write-Host "`nInstallation complete!" -ForegroundColor Green
Write-Host "ATA executable copied to: $binDir" -ForegroundColor Cyan
Write-Host "Added to PATH environment variable for current user" -ForegroundColor Cyan
Write-Host "You can now use 'ata' command in any terminal" -ForegroundColor Cyan
Write-Host "Note: Restart your terminal if the command is not recognized immediately" -ForegroundColor Yellow
Write-Host "`nTest the installation: ata help" -ForegroundColor White