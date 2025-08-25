param(
    [string]$TargetDir,
    [Alias('dbg')][switch]$Debug,      # -debug 显示 ffmpeg 日志（别名：-dbg）
    [Alias('del','d')][switch]$Delete,     # -d 删除原图片（别名：-del, -d）
    [switch]$Force,      # -f 忽略数量和文件夹检查
    [switch]$Recursive   # -r 直接递归
)

Import-Module -Force -Name "$PSScriptRoot\module\ata.psm1"

# 默认显示帮助
if (-not $TargetDir -or $TargetDir -in @("/help", "-h", "--help")) {
    Show-Help
    exit 0
}

Convert-Images -TargetDir $TargetDir -Debug:$Debug -Delete:$Delete -Force:$Force -Recursive:$Recursive