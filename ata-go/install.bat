@echo off
echo Installing ATA using PowerShell script...
echo.

:: Check if PowerShell script exists
if not exist "%~dp0\install-powershell.ps1" (
    echo Error: install-powershell.ps1 not found
    echo Please ensure the PowerShell installation script is in the same directory
    exit /b 1
)

:: Run PowerShell installation script
powershell -ExecutionPolicy Bypass -File "%~dp0\install-powershell.ps1"

echo.
echo Batch script completed. Check above output for installation status.