# ATA Release Build Script
Write-Host "Building ATA release version..." -ForegroundColor Green
Write-Host ""

# 清理旧的构建文件
if (Test-Path "ata.exe") { Remove-Item "ata.exe" }
if (!(Test-Path "release")) { New-Item -ItemType Directory -Name "release" }
Get-ChildItem "release" | Remove-Item -Force

# 构建Windows版本
Write-Host "Building for Windows..." -ForegroundColor Yellow
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -ldflags "-s -w" -o "release/ata-windows.exe" "./cmd/ata"
if ($LASTEXITCODE -ne 0) {
    Write-Host "Windows build failed!" -ForegroundColor Red
    Read-Host "Press Enter to exit"
    exit 1
}

# 构建Linux版本
Write-Host "Building for Linux..." -ForegroundColor Yellow
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -ldflags "-s -w" -o "release/ata-linux" "./cmd/ata"
if ($LASTEXITCODE -ne 0) {
    Write-Host "Linux build failed!" -ForegroundColor Red
    Read-Host "Press Enter to exit"
    exit 1
}

# 构建macOS版本
Write-Host "Building for macOS..." -ForegroundColor Yellow
$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -ldflags "-s -w" -o "release/ata-macos" "./cmd/ata"
if ($LASTEXITCODE -ne 0) {
    Write-Host "macOS build failed!" -ForegroundColor Red
    Read-Host "Press Enter to exit"
    exit 1
}

# 重置环境变量
Remove-Item Env:GOOS
Remove-Item Env:GOARCH

Write-Host ""
Write-Host "✓ All builds completed successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "Release files created in 'release' directory:" -ForegroundColor Cyan
Get-ChildItem "release" | Format-Table Name, Length, LastWriteTime
Write-Host ""
Write-Host "Usage:" -ForegroundColor Cyan
Write-Host "- Windows users: Download and run ata-windows.exe"
Write-Host "- Linux users: Download ata-linux, make it executable (chmod +x ata-linux)"
Write-Host "- macOS users: Download ata-macos, make it executable (chmod +x ata-macos)"
Write-Host ""
Write-Host "When users run the executable without arguments, they will see the interactive installer." -ForegroundColor Yellow
Write-Host ""
Read-Host "Press Enter to exit"