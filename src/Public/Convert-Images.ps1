function Convert-Images {
    param(
        [string]$TargetDir,
        [switch]$Debug,
        [switch]$Delete,
        [switch]$Force,
        [switch]$Recursive
    )

    Add-Type -AssemblyName Microsoft.VisualBasic

    # 处理路径
    if ($TargetDir -in @("./", "/")) {
        $TargetDir = Get-Location
    }

    $TargetDir = Convert-Path $TargetDir
    if (-not (Test-Path $TargetDir)) {
        Write-Host "目录不存在：" $TargetDir -ForegroundColor Red
        return
    }

    $TargetDir = $TargetDir.TrimEnd([char[]]@('\','/'))

    Write-Host "准备转换目录：" $TargetDir -ForegroundColor Cyan

    # 处理 Debug 参数（仅由开关决定）
    $DebugMode = [bool]$Debug

    # 检测子文件夹是否有图片
    $detect = Detect-SubImages -TargetDir $TargetDir -IncludePatterns $IncludePatterns

    # 确定是否递归
    $doRecursive = Confirm-Recursive -HasSubImages:$($detect.HasSubImages) -SubDirCount $($detect.SubDirs.Count) -Recursive:$Recursive -Force:$Force

    # 是否删除原图片
    $deleteOriginal = Decide-DeleteOriginal -Delete:$Delete

    # 获取最终要处理的文件列表
    $files = Get-TargetFiles -TargetDir $TargetDir -DoRecursive:$doRecursive -SupportedImageExtensions $SupportedImageExtensions

    # 图片数量提示
    $totalFiles = $files.Count
    $shouldContinue = Confirm-ContinueLargeCount -TotalFiles $totalFiles -Force:$Force
    if (-not $shouldContinue) { return }

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
            $result = Invoke-FFmpeg -InputFullPath $f.FullName -OutputFullPath ([System.IO.Path]::GetFullPath($out)) -DebugMode:$DebugMode

            if ($DebugMode) {
                if ($result.Stdout) { Write-Output $result.Stdout }
                if ($result.Stderr) { Write-Output $result.Stderr }
            }

            if ($result.TimedOut) {
                $failCount++
                Write-Host "$($f.Name) 转换超时，耗时 $($stopwatch.ElapsedMilliseconds) ms" -ForegroundColor Red
                continue
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
} 