param(
    [string]$TargetDir,
    [switch]$Debug,      # -debug 显示 ffmpeg 日志
    [switch]$Delete,     # -d 删除原图片
    [switch]$Force,      # -f 忽略数量和文件夹检查
    [switch]$Recursive   # -r 直接递归
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

# 生成唯一输出文件名
function Get-UniqueOutputPath {
    param (
        [string]$OutputPath
    )
    $directory = [System.IO.Path]::GetDirectoryName($OutputPath)
    $baseName = [System.IO.Path]::GetFileNameWithoutExtension($OutputPath)
    $extension = [System.IO.Path]::GetExtension($OutputPath)
    $counter = 0
    $newPath = $OutputPath

    while (Test-Path $newPath) {
        $counter++
        $newBaseName = "$baseName ($counter)"
        $newPath = Join-Path $directory "$newBaseName$extension"
    }
    return $newPath
}

# 默认显示帮助
if (-not $TargetDir -or $TargetDir -in @("/help", "-h", "--help")) {
    Show-Help
    exit 0
}

# 处理路径
if ($TargetDir -in @("./", "/")) {
    $TargetDir = Get-Location
}

$TargetDir = Convert-Path $TargetDir
if (-not (Test-Path $TargetDir)) {
    Write-Host "目录不存在：" $TargetDir -ForegroundColor Red
    exit 1
}

$TargetDir = $TargetDir.TrimEnd('\', '/')
Write-Host "准备转换目录：" $TargetDir -ForegroundColor Cyan

# 处理大小写不敏感 Debug 参数
$DebugMode = $Debug -or ($PSBoundParameters.Keys | Where-Object { $_.ToLower() -eq "debug" })

# 检测子文件夹是否有图片
$subDirs = Get-ChildItem -Path $TargetDir -Directory
$hasSubImages = $false
foreach ($d in $subDirs) {
    if ((Get-ChildItem -Path $d.FullName -Include *.jpg, *.jpeg, *.png -File -Recurse | Measure-Object).Count -gt 0) {
        $hasSubImages = $true
        break
    }
}

# 确定是否递归
$doRecursive = $false
if ($hasSubImages -and ($subDirs.Count -gt 0)) {
    if ($Recursive -or $Force) {
        $doRecursive = $true
    } else {
        $choice = Read-Host "检测到子文件夹内有可转换图片，是否递归执行？(Y/N)"
        if ($choice.ToUpper() -eq "Y") { $doRecursive = $true }
    }
}

# 是否删除原图片
$deleteOriginal = $Delete
if (-not $Delete) {
    $choice = Read-Host "是否删除原图片？(D 删除 / S 保留)"
    if ($choice.ToUpper() -eq "D") { 
        $deleteOriginal = $true 
    } else {
        $deleteOriginal = $false
    }
}
$null = $deleteOriginal

# 获取最终要处理的文件列表
if ($doRecursive) {
    $files = Get-ChildItem -Path $TargetDir -Recurse -File -Depth 3 | Where-Object {
        $_.Extension -in ".jpg", ".jpeg", ".png"
    }
} else {
    $files = Get-ChildItem -Path $TargetDir -File | Where-Object {
        $_.Extension -in ".jpg", ".jpeg", ".png"
    }
}

# 图片数量提示
$totalFiles = $files.Count
if ($totalFiles -gt 200 -and -not $Force) {
    $choice = Read-Host "发现 $totalFiles 张可转换图片，是否继续？(Y/N)"
    if ($choice.ToUpper() -ne "Y") { exit 0 }
}

# 转换计数器
$successCount = 0
$failCount = 0

# 开始处理每个文件
$index = 0
foreach ($f in $files) {
    $index++
    Write-Progress -Activity "转换图片" -Status "处理 $($f.Name) ($index/$totalFiles)" -PercentComplete ($index / $totalFiles * 100)
    $out = Join-Path $f.DirectoryName ($f.BaseName + ".avif")
    $out = Get-UniqueOutputPath -OutputPath $out
    $stopwatch = [System.Diagnostics.Stopwatch]::StartNew()

    try {
        $ffmpegArgs = "-i `"$($f.FullName)`" -c:v libaom-av1 -still-picture 1 -crf 28 -pix_fmt yuv420p `"$([System.IO.Path]::GetFullPath($out))`""

        # 调用 ffmpeg
        $processInfo = New-Object System.Diagnostics.ProcessStartInfo
        $processInfo.FileName = "ffmpeg"
        $processInfo.Arguments = $ffmpegArgs
        $processInfo.UseShellExecute = $false
        $processInfo.CreateNoWindow = $true
        $processInfo.RedirectStandardOutput = $true
        $processInfo.RedirectStandardError = $true

        $process = New-Object System.Diagnostics.Process
        $process.StartInfo = $processInfo
        $process.Start()

        $timeout = 30000  # 30秒超时
        if ($DebugMode) {
            $process.WaitForExit()
            $stdout = $process.StandardOutput.ReadToEnd()
            $stderr = $process.StandardError.ReadToEnd()
            if ($stdout) { Write-Output $stdout }
            if ($stderr) { Write-Output $stderr }
        } else {
            if (-not $process.WaitForExit($timeout)) {
                $process.Kill()
                $failCount++
                Write-Host "$($f.Name) 转换超时，耗时 $($stopwatch.ElapsedMilliseconds) ms" -ForegroundColor Red
                continue
            }
        }

        # 检查输出文件，删除原图
        if (Test-Path $out) {
            if ($deleteOriginal) {
                Remove-Item $f.FullName -Force
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

# 总结
Write-Host "--------------------------------------------------" -ForegroundColor Cyan
Write-Host "转换完成，总计成功：$successCount 张，失败：$failCount 张" -ForegroundColor Cyan