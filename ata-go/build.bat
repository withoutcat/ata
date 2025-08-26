@echo off
echo Building ATA...

:: Create bin directory if it doesn't exist
if not exist "%~dp0\bin" (
    mkdir "%~dp0\bin"
    echo Created bin directory
)

:: Compile project
echo Compiling...
cd /d "%~dp0"
go build -o bin/ata.exe cmd/ata/main.go

if %ERRORLEVEL% neq 0 (
    echo Build failed!
    exit /b 1
)

echo Build successful!
echo Executable generated: %~dp0\bin\ata.exe

echo.
echo To install ATA, run install.bat
echo.