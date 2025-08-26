@echo off
echo 正在构建ATA...

:: 创建bin目录（如果不存在）
if not exist "%~dp0\bin" (
    mkdir "%~dp0\bin"
    echo 已创建bin目录
)

:: 编译项目
echo 正在编译...
cd /d "%~dp0"
go build -o bin/ata.exe cmd/ata/main.go

if %ERRORLEVEL% neq 0 (
    echo 构建失败！
    exit /b 1
)

echo 构建成功！
echo 可执行文件已生成: %~dp0\bin\ata.exe

echo.
echo 要安装ATA，请运行install.bat
echo.

pause