@echo off
echo 正在设置FFmpeg环境...

:: 创建ffmpeg目录（如果不存在）
if not exist "%~dp0\ffmpeg\bin" (
    mkdir "%~dp0\ffmpeg\bin"
    echo 已创建ffmpeg目录
)

:: 检查源ffmpeg文件是否存在
set "SOURCE_FFMPEG=%~dp0..\module\ffmpeg-n7.1-latest-win64-gpl-7.1\bin"
if not exist "%SOURCE_FFMPEG%\ffmpeg.exe" (
    echo 错误：找不到源FFmpeg文件！
    echo 请确保FFmpeg文件位于 %SOURCE_FFMPEG% 目录中
    exit /b 1
)

:: 复制ffmpeg文件
echo 正在复制FFmpeg文件...
copy "%SOURCE_FFMPEG%\ffmpeg.exe" "%~dp0\ffmpeg\bin\" /Y
copy "%SOURCE_FFMPEG%\ffprobe.exe" "%~dp0\ffmpeg\bin\" /Y
copy "%SOURCE_FFMPEG%\ffplay.exe" "%~dp0\ffmpeg\bin\" /Y

if %ERRORLEVEL% neq 0 (
    echo 复制FFmpeg文件失败！
    exit /b 1
)

echo FFmpeg环境设置完成！
echo 现在可以运行build.bat构建项目了