param(
    [string]$TargetDir,
    [switch]$Debug,
    [switch]$d,  # 删除原图片
    [switch]$f   # 忽略文件夹/文件数量检查
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
    Write-Host "  ata ./ -d               转换完成后删除原图片（静默执行）" -ForegroundColor Yellow
    Write-Host "  ata ./ -f               忽略文件夹/文件数量检查（静默执行）" -ForegroundColor Yellow
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

# 统计目录和文件数量
$subFolderCount = (Get-ChildItem -Path $TargetDir -Directory).Count
$imageFiles = Get-ChildItem -Recurse -Path $TargetDir -Include *.jpg, *.jpeg, *.png -File
$imageCount = $imageFiles.Count

# 提示子文件夹过多
if (-not $f -and $subFolderCount -gt 10) {
    $resp = Read-Host "警告：当前目录包含 $subFolderCount 个子文件夹，递归操作可能影响大量文件。是否继续？(Y/N)"
    if ($resp -notin @("Y","y")) { exit 0 }
}

# 提示图片过多
if (-not $f -and $imageCount -gt 200) {
    $resp = Read-Host "检测到 $imageCount 张可转换图片，处理可能需要较长时间。是否继续？(Y/N)"
    if ($resp -notin @("Y","y")) { exit 0 }
}

# 是否删除原图片询问（只在当前目录并且未使用 -d 参数时）
$DeleteOriginal = $d.IsPresent
if (-not $DeleteOriginal -and $TargetDir -eq (Get-Location)) {
    $resp = Read-Host "是否删除原图片？ 输入 D 删除并移动到回收站，输入 S 保留原图片"
    if ($resp -in @("D","d")) { $DeleteOriginal = $true }
    else { $DeleteOriginal = $false }
}

$successCount = 0
$failCount = 0

foreach ($f in $imageFiles) {
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
            if ($DeleteOriginal) {
                [Microsoft.VisualBasic.FileIO.FileSystem]::DeleteFile($f.FullName, 'OnlyErrorDialogs', 'SendToRecycleBin')
            }
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
