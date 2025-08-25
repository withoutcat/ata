function Show-Help {
    Write-Host "--------------------------------------------------" -ForegroundColor Cyan
    Write-Host "ata - 图片批量转换为 AVIF (PowerShell 版本)" -ForegroundColor Cyan
    Write-Host "用法：" -ForegroundColor Cyan
    Write-Host "  ata ./                  在当前目录转换所有支持的图片格式" -ForegroundColor Yellow
    Write-Host "  ata `"D:\MyPictures\2025-08-25`"  指定目录转换图片" -ForegroundColor Yellow
    Write-Host "  ata /help               显示此帮助文档" -ForegroundColor Yellow
    Write-Host "  ata ./ -debug           显示 ffmpeg 日志（或使用 -dbg）" -ForegroundColor Yellow
    Write-Host "  ata ./ -d -f -r         静默删除原图片、忽略数量、直接递归（-d 等同于 -Delete/-del）" -ForegroundColor Yellow
    Write-Host ("  支持的格式：{0}" -f ($SupportedImageExtensions -join ", ")) -ForegroundColor DarkGray
    Write-Host "--------------------------------------------------" -ForegroundColor Cyan
} 