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

# 先构建主程序
Write-Host "  Building main program..." -ForegroundColor Cyan
go build -ldflags "-s -w" -o "cmd/installer/ata.exe" "./cmd/ata"
if ($LASTEXITCODE -ne 0) {
    Write-Host "Windows main program build failed!" -ForegroundColor Red
    Read-Host "Press Enter to exit"
    exit 1
}

# 构建包含嵌入文件的安装程序
Write-Host "  Building installer with embedded executable..." -ForegroundColor Cyan
go build -ldflags "-s -w" -o "release/ata-installer-windows.exe" "./cmd/installer"
if ($LASTEXITCODE -ne 0) {
    Write-Host "Windows installer build failed!" -ForegroundColor Red
    Read-Host "Press Enter to exit"
    exit 1
}

# 清理临时文件
Remove-Item "cmd/installer/ata.exe" -Force

# 构建Linux版本
Write-Host "Building for Linux..." -ForegroundColor Yellow
$env:GOOS = "linux"
$env:GOARCH = "amd64"

# 确保清理之前的文件
if (Test-Path "cmd/installer/ata.exe") { Remove-Item "cmd/installer/ata.exe" -Force }
if (Test-Path "cmd/installer/ata") { Remove-Item "cmd/installer/ata" -Force }

# 先构建主程序
Write-Host "  Building main program..." -ForegroundColor Cyan
go build -ldflags "-s -w" -o "cmd/installer/ata" "./cmd/ata"
if ($LASTEXITCODE -ne 0) {
    Write-Host "Linux main program build failed!" -ForegroundColor Red
    Read-Host "Press Enter to exit"
    exit 1
}

# 构建包含嵌入文件的安装程序
Write-Host "  Building installer with embedded executable..." -ForegroundColor Cyan
go build -ldflags "-s -w" -o "release/ata-installer-linux" "./cmd/installer"
if ($LASTEXITCODE -ne 0) {
    Write-Host "Linux installer build failed!" -ForegroundColor Red
    Read-Host "Press Enter to exit"
    exit 1
}

# 清理临时文件
Remove-Item "cmd/installer/ata" -Force

# 构建macOS版本
Write-Host "Building for macOS..." -ForegroundColor Yellow
$env:GOOS = "darwin"
$env:GOARCH = "amd64"

# 确保清理之前的文件
if (Test-Path "cmd/installer/ata.exe") { Remove-Item "cmd/installer/ata.exe" -Force }
if (Test-Path "cmd/installer/ata") { Remove-Item "cmd/installer/ata" -Force }

# 先构建主程序
Write-Host "  Building main program..." -ForegroundColor Cyan
go build -ldflags "-s -w" -o "cmd/installer/ata" "./cmd/ata"
if ($LASTEXITCODE -ne 0) {
    Write-Host "macOS main program build failed!" -ForegroundColor Red
    Read-Host "Press Enter to exit"
    exit 1
}

# 构建包含嵌入文件的安装程序
Write-Host "  Building installer with embedded executable..." -ForegroundColor Cyan
go build -ldflags "-s -w" -o "release/ata-installer-macos" "./cmd/installer"
if ($LASTEXITCODE -ne 0) {
    Write-Host "macOS installer build failed!" -ForegroundColor Red
    Read-Host "Press Enter to exit"
    exit 1
}

# 清理临时文件
Remove-Item "cmd/installer/ata" -Force

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
Write-Host "- Windows users: Download and run ata-installer-windows.exe"
Write-Host "- Linux users: Download ata-installer-linux, make it executable (chmod +x ata-installer-linux)"
Write-Host "- macOS users: Download ata-installer-macos, make it executable (chmod +x ata-installer-macos)"
Write-Host ""
Write-Host "Each installer contains the embedded ata program and will:" -ForegroundColor Yellow
Write-Host "1. Extract and install ata to user's bin directory" -ForegroundColor Yellow
Write-Host "2. Set up PATH environment variable" -ForegroundColor Yellow
Write-Host "3. Check and install FFmpeg dependencies" -ForegroundColor Yellow
Write-Host "After installation, users can use 'ata' command from anywhere." -ForegroundColor Yellow