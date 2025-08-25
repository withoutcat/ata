param(
    [string]$TargetDir,
    [Alias('dbg')][switch]$ShowDebug,      # -debug 显示 ffmpeg 日志（别名：-dbg）
    [Alias('del','d')][switch]$Delete,     # -d 删除原图片（别名：-del, -d）
    [switch]$Force,      # -f 忽略数量和文件夹检查
    [switch]$Recursive   # -r 直接递归
)

Import-Module -Force -Name "$PSScriptRoot\module\ata.psm1"

# 子命令：ani 动图合成
if ($TargetDir -and ($TargetDir -eq 'ani')) {
    $nextArg = $args | Select-Object -First 1
    if (-not $nextArg) { Write-Host "用法：ata ani <目录>" -ForegroundColor Yellow; exit 1 }
    Create-AvifAnimation -InputDir $nextArg -ShowDebug:$ShowDebug
    exit 0
}

# 子命令：ppt 演示文稿动图合成（0.4秒帧间隔）
if ($TargetDir -and ($TargetDir -eq 'ppt')) {
    $nextArg = $args | Select-Object -First 1
    if (-not $nextArg) { Write-Host "用法：ata ppt <目录>" -ForegroundColor Yellow; exit 1 }
    Create-AvifAnimation -InputDir $nextArg -Fps 2.5 -ShowDebug:$ShowDebug
    exit 0
}

# 默认显示帮助
if (-not $TargetDir -or $TargetDir -in @("/help", "-h", "--help")) {
    Show-Help
    exit 0
}

Convert-Images -TargetDir $TargetDir -ShowDebug:$ShowDebug -Delete:$Delete -Force:$Force -Recursive:$Recursive