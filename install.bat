@echo off
chcp 65001 >nul
title ata 安装程序

echo.
echo ========================================
echo            ata 安装程序
echo ========================================
echo.

REM 检查是否在项目根目录
if not exist "ata.ps1" (
    echo 错误：未找到 ata.ps1，请确保在项目根目录运行此脚本。
    echo.
    pause
    exit /b 1
)

echo 正在将 ata 添加到用户 PATH 环境变量...
echo.

REM 获取当前目录的完整路径
for %%i in ("%~dp0.") do set "projectPath=%%~fi"
set "projectPath=%projectPath:~0,-1%"

echo 项目路径：%projectPath%
echo.

REM 使用 setx 添加到用户 PATH
setx PATH "%PATH%;%projectPath%"

if %errorlevel% equ 0 (
    echo.
    echo ========================================
    echo           安装完成！
    echo ========================================
    echo.
    echo ata 已成功添加到用户 PATH 环境变量中。
    echo 请重新打开 PowerShell 或命令提示符以使更改生效。
    echo.
    echo 使用 'ata /help' 查看帮助信息。
    echo.
) else (
    echo.
    echo ========================================
    echo           安装失败！
    echo ========================================
    echo.
    echo 安装过程中出现错误。
    echo 请尝试以管理员身份运行此脚本。
    echo.
)

echo 按任意键退出...
pause >nul 