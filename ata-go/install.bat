@echo off
setlocal enabledelayedexpansion

echo 正在安装ATA...

:: 检查ata.exe是否存在
if not exist "%~dp0\bin\ata.exe" (
    echo 错误: 找不到ata.exe文件
    echo 请先编译项目，然后再运行安装脚本
    exit /b 1
)

:: 创建bin目录（如果不存在）
if not exist "%USERPROFILE%\bin" (
    mkdir "%USERPROFILE%\bin"
    echo 已创建 %USERPROFILE%\bin 目录
)

:: 复制ata.exe到bin目录
copy /Y "%~dp0\bin\ata.exe" "%USERPROFILE%\bin\ata.exe" > nul
echo 已复制 ata.exe 到 %USERPROFILE%\bin

:: 创建ffmpeg目录（如果不存在）
if not exist "%USERPROFILE%\bin\ffmpeg\bin" (
    mkdir "%USERPROFILE%\bin\ffmpeg\bin"
    echo 已创建 %USERPROFILE%\bin\ffmpeg\bin 目录
)

:: 复制ffmpeg文件
copy /Y "%~dp0\ffmpeg\bin\ffmpeg.exe" "%USERPROFILE%\bin\ffmpeg\bin\" > nul
copy /Y "%~dp0\ffmpeg\bin\ffprobe.exe" "%USERPROFILE%\bin\ffmpeg\bin\" > nul
copy /Y "%~dp0\ffmpeg\bin\ffplay.exe" "%USERPROFILE%\bin\ffmpeg\bin\" > nul
echo 已复制 FFmpeg 文件到 %USERPROFILE%\bin\ffmpeg\bin

:: 检查PATH环境变量中是否已包含bin目录
echo 正在检查PATH环境变量...
set "binPath=%USERPROFILE%\bin"
set "found=0"

for /f "tokens=*" %%a in ('echo !PATH!') do (
    set "currentPath=%%a"
    if "!currentPath!"=="!binPath!" (
        set "found=1"
    )
    if "!currentPath!"=="!binPath!;" (
        set "found=1"
    )
)

:: 如果PATH中不包含bin目录，则添加
if "!found!"=="0" (
    setx PATH "%PATH%;%USERPROFILE%\bin"
    echo 已将 %USERPROFILE%\bin 添加到PATH环境变量
) else (
    echo %USERPROFILE%\bin 已在PATH环境变量中
)

echo.
echo 安装完成！
echo 请重新打开命令提示符或PowerShell窗口以使用ata命令
echo.

pause