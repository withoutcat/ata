function Create-AvifAnimation {
    param(
        [Parameter(Mandatory=$true)][string]$InputDir,
        [string]$Output,
        [int]$Fps = 10,
        [int]$Width,
        [int]$Height,
        [ValidateSet('contain','cover','stretch')][string]$ScaleMode = 'contain',
        [string]$Background = 'black',
        [switch]$Alpha,
        [int]$Crf = 28,
        [int]$Speed = 5,
        [int]$Threads,
        [Alias('dbg')][switch]$ShowDebug
    )

    if (-not (Test-Path $InputDir)) {
        throw "输入目录不存在：$InputDir"
    }

    $avifFiles = Get-ChildItem -Path $InputDir -File | Where-Object { $_.Extension.ToLower() -eq '.avif' } | Sort-Object Name
    if ($avifFiles.Count -lt 2) {
        throw "动图至少需要 2 张 AVIF 帧。目录下找到 $($avifFiles.Count) 张。"
    }

    if (-not $Output) {
        $baseName = Split-Path -Leaf $InputDir
        $Output = Join-Path $InputDir ($baseName + '.avif')
        $Output = Get-UniqueOutputPath -OutputPath $Output
    }

    # 准备帧序列（统一尺寸与像素格式，输出为 PNG 序列）
    $prep = Prepare-AnimationFrames -InputFiles $avifFiles.FullName -Width $Width -Height $Height -ScaleMode $ScaleMode -Background $Background -Alpha:$Alpha -ShowDebug:$ShowDebug

    # 编码为 AVIF 动画
    Encode-AvifAnimation -FramesDir $prep.FramesDir -Fps $Fps -Output $Output -Alpha:$Alpha -Crf $Crf -Speed $Speed -Threads $Threads -ShowDebug:$ShowDebug

    # 清理临时目录
    if (-not $prep.KeepTemp) {
        Remove-Item -Recurse -Force $prep.WorkDir
    }

    Write-Host "动图生成成功：$Output" -ForegroundColor Green
} 