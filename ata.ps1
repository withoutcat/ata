param(
    [string]$TargetDir,
    [switch]$Debug,
    [switch]$Delete,   # -d 删除原图片
    [switch]$Force,    # -f 忽略数量和文件夹检查
    [switch]$Recursive # -r 直接递归
)

Add-Type -AssemblyName Microsoft.VisualBasic

# 帮助文档函数
function Show-Help {
    Write-Host "--------------------------------------------------" -ForegroundColor Cyan
    Write-Host "ata - 图片批量转换为 AVIF (PowerShell 版本)" -ForegroundColor Cyan
    Write-Host "用法：" -ForegroundColor Cyan
    Write-Host "  ata ./                  在当前目录转换所有 jpg/jpeg/png" -ForegroundColor Yellow
    Write-Host "  ata `"D:\MyPictures\2025-08-25`"  指定目录转换图片" -ForegroundColor Yellow
    Write-Host "  ata /help               显示此帮助文档" -ForegroundColor Yellow
    Write-Host "  ata ./ -debug           显示 ffmpeg 日志" -ForegroundColor Yellow
    Write-Host "  ata ./ -d -f -r         静默删除原图片、忽略数量、直接递归" -ForegroundColor Yellow
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
Write-Host "准备转换目录：" $TargetDir -ForegroundColor Cyan

# 处理大小写不敏感 Debug 参数
$DebugMode = $Debug -or ($PSBoundParameters.Keys | Where-Object { $_.ToLower() -eq "debug" })

# 检测目录下可转换图片数量
$filesInDir = Get-ChildItem -Path $TargetDir -Include *.jpg, *.jpeg, *.png -File
$fileCount = $filesInDir.Count

# 检测子文件夹及其可转换图片
$subDirs = Get-ChildItem -Path $TargetDir -Directory
$hasSubImages = $false
foreach ($d in $subDirs) {
    if ((Get-ChildItem -Path $d.FullName -Include *.jpg, *.jpeg, *.png -File -Recurse | Measure-Object).Count -gt 0) {
        $hasSubImages = $true
        break
    }
}

$doRecursive = $false
# 递归可选逻辑
if ($subDirs.Count -gt 0 -and $hasSubImages) {
    if ($Recursive -or $Force) {
        $doRecursive = $true
    } else {
        $choice = Read-Host "检测到子文件夹内有可转换图片，是否递归执行？(Y/N)"
        if ($choice.ToUpper() -eq "Y") { $doRecursive = $true }
    }
}

# 总计文件数（递归或非递归）
$totalFiles = if ($doRecursive) {
    (Get-ChildItem -Path $TargetDir -Include *.jpg, *.jpeg, *.png -File -Recurse).Count
} else { $fileCount }

# 图片数量提示
if ($totalFiles -gt 200 -and -not $Force) {
    $choice = Read-Host "发现 $totalFiles 张可转换图片，是否继续？(Y/N)"
    if ($choice.ToUpper() -ne "Y") { exit 0 }
}

# 确认是否删除原图片（当前目录执行）
$deleteOriginal = $Delete
if (-not $Delete -and $TargetDir -eq (Get-Location)) {
    $choice = Read-Host "是否删除原图片？(D 删除 / S 保留)"
    if ($choice.ToUpper() -eq "D") { $deleteOriginal = $true }
}

# 确定递归模式参数
$gciParams = @{
    Path = $TargetDir
    Include = '*.jpg','*.jpeg','*.png'
    File = $true
    Recurse = $doRecursive
}

# 转换计数器
$successCount = 0
$failCount = 0

foreach ($f in Get-ChildItem @gciParams) {
    $out = Join-Path $f.DirectoryName ($f.BaseName + ".avif")
    $stopwatch = [System.Diagnostics.Stopwatch]::StartNew()

    try {
        $ffmpegArgs = "-i `"$($f.FullName)`" -c:v libaom-av1 -still-picture 1 -crf 28 -pix_fmt yuv420p `"$out`""

        if ($DebugMode) {
            Start-Process ffmpeg -ArgumentList $ffmpegArgs -Wait -NoNewWindow -PassThru | Out-Null
        } else {
            Start-Process ffmpeg -ArgumentList $ffmpegArgs -Wait -NoNewWindow -PassThru *> $null
        }

        $stopwatch.Stop()

        if (Test-Path $out) {
            if ($deleteOriginal) {
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
