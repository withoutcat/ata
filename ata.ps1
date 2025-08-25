param(
    [string]$TargetDir,
    [switch]$Debug
)

Add-Type -AssemblyName Microsoft.VisualBasic

# 帮助文档函数
function Show-Help {
    Write-Host "--------------------------------------------------" -ForegroundColor Cyan
    Write-Host "ata - 图片批量转换为 AVIF (PowerShell 版本)" -ForegroundColor Cyan
    Write-Host "用法：" -ForegroundColor Cyan
    Write-Host "  ata ./                  在当前目录递归转换所有 jpg/jpeg/png" -ForegroundColor Yellow
    Write-Host "  ata `"D:\MyPictures\2025-08-25`"  指定目录递归转换图片" -ForegroundColor Yellow
    Write-Host "  ata /help               显示此帮助文档" -ForegroundColor Yellow
    Write-Host "  ata ./ -debug           显示包含 ffmpeg 的详细日志（不区分大小写）" -ForegroundColor Yellow
    Write-Host "--------------------------------------------------" -ForegroundColor Cyan
}

# 默认显示帮助或 /help 参数
if (-not $TargetDir -or $TargetDir -in @("/help","-h","--help")) {
    Show-Help
    exit 0
}

# 处理路径，./ 或 / 表示当前目录
if ($TargetDir -in @("./","/")) {
    $TargetDir = Get-Location
}

# 检查目录是否存在
if (-not (Test-Path $TargetDir)) {
    Write-Host "目录不存在：" $TargetDir -ForegroundColor Red
    exit 1
}

$TargetDir = $TargetDir.TrimEnd('\','/')
Write-Host "开始递归转换目录：" $TargetDir -ForegroundColor Cyan

# 处理大小写不敏感 Debug 参数
$DebugMode = $Debug -or ($PSBoundParameters.Keys | Where-Object { $_.ToLower() -eq "debug" })

$successCount = 0
$failCount = 0

foreach ($f in Get-ChildItem -Recurse -Path $TargetDir -Include *.jpg, *.jpeg, *.png -File) {
    $out = Join-Path $f.DirectoryName ($f.BaseName + ".avif")
    $stopwatch = [System.Diagnostics.Stopwatch]::StartNew()

    try {
        $ffmpegArgs = "-i `"$($f.FullName)`" -c:v libaom-av1 -still-picture 1 -crf 28 -pix_fmt yuv420p `"$out`""

        if ($DebugMode) {
            # 调试模式：显示 ffmpeg 日志
            Start-Process ffmpeg -ArgumentList $ffmpegArgs -Wait -NoNewWindow -PassThru | Out-Null
        } else {
            # 普通模式：屏蔽 ffmpeg 输出
            Start-Process ffmpeg -ArgumentList $ffmpegArgs -Wait -NoNewWindow -PassThru *> $null
        }

        $stopwatch.Stop()

        if (Test-Path $out) {
            [Microsoft.VisualBasic.FileIO.FileSystem]::DeleteFile($f.FullName, 'OnlyErrorDialogs', 'SendToRecycleBin')
            $successCount++
            Write-Host "$($f.Name) 转换成功，耗时 $($stopwatch.ElapsedMilliseconds) ms" -ForegroundColor Green
        } else {
            $failCount++
            Write-Host "$($f.Name) 转换失败，耗时 $($stopwatch.ElapsedMilliseconds) ms" -ForegroundColor Red
        }
    }
    catch {
        $stopwatch.Stop()
        $failCount++
        Write-Host "$($f.Name) 转换失败，耗时 $($stopwatch.ElapsedMilliseconds) ms" -ForegroundColor Red
    }
}

Write-Host "--------------------------------------------------" -ForegroundColor Cyan
Write-Host "转换完成，总计成功：$successCount 张，失败：$failCount 张" -ForegroundColor Cyan
