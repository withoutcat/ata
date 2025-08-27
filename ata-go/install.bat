@echo off
echo Installing ATA...
echo.

:: First, build the latest version
echo Step 1: Building latest version...
call "%~dp0\build.bat"
if %ERRORLEVEL% neq 0 (
    echo Build failed! Cannot proceed with installation.
    exit /b 1
)

echo.
echo Step 2: Installing to user bin directory...

:: Check if PowerShell script exists
if not exist "%~dp0\install-powershell.ps1" (
    echo Error: install-powershell.ps1 not found
    echo Please ensure the PowerShell installation script is in the same directory
    exit /b 1
)

:: Run PowerShell installation script
powershell -ExecutionPolicy Bypass -File "%~dp0\install-powershell.ps1"

echo.
echo Installation completed. Check above output for status.