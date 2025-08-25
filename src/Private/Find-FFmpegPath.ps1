function Find-FFmpegPath {
    # 固定优先路径：<project_root>\module\ffmpeg-n7.1-latest-win64-gpl-7.1\bin\ffmpeg.exe
    # 若不存在则回退到系统 PATH，并在控制台提示所用路径；都不存在则返回 $null

    # 计算工程根目录：当前文件位于 <project_root>\src\Private
    $srcRoot = Split-Path -Parent $PSScriptRoot       # <project_root>\src
    $projectRoot = Split-Path -Parent $srcRoot         # <project_root>

    $fixedPath = Join-Path $projectRoot "module\ffmpeg-n7.1-latest-win64-gpl-7.1\bin\ffmpeg.exe"
    if (Test-Path $fixedPath) {
        return (Resolve-Path $fixedPath).Path
    }

    # 回退系统 PATH
    try {
        $ff = (Get-Command ffmpeg -ErrorAction Stop).Source
        if ($ff) {
            Write-Host "未在模块内发现内置 ffmpeg，已回退到系统 PATH：$ff" -ForegroundColor Yellow
            return $ff
        }
    } catch {}

    return $null
} 