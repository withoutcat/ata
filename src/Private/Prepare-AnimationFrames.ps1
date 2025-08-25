function Prepare-AnimationFrames {
    param(
        [string[]]$InputFiles,
        [int]$Width,
        [int]$Height,
        [ValidateSet('contain','cover','stretch')][string]$ScaleMode = 'contain',
        [string]$Background = 'black',
        [switch]$Alpha,
        [switch]$KeepTemp,
        [switch]$Debug
    )

    if (-not $InputFiles -or $InputFiles.Count -lt 2) {
        throw "需要至少两张帧图"
    }

    $workDir = Join-Path $env:TEMP ("ata_ani_" + [Guid]::NewGuid().ToString("N"))
    New-Item -ItemType Directory -Path $workDir | Out-Null
    $framesDir = Join-Path $workDir "frames"
    New-Item -ItemType Directory -Path $framesDir | Out-Null

    # 若未指定尺寸，先扫描所有源图获取最大宽高
    if (-not $Width -or -not $Height) {
        $maxW = 0; $maxH = 0
        foreach ($f in $InputFiles) {
            $probe = ffprobe -v error -select_streams v:0 -show_entries stream=width,height -of csv=s=x:p=0 -- "$f" 2>$null
            if ($LASTEXITCODE -ne 0 -or -not $probe) {
                continue
            }
            $parts = $probe -split 'x'
            if ($parts.Count -eq 2) {
                $w = [int]$parts[0]; $h = [int]$parts[1]
                if ($w -gt $maxW) { $maxW = $w }
                if ($h -gt $maxH) { $maxH = $h }
            }
        }
        if (-not $Width) { $Width = $maxW }
        if (-not $Height) { $Height = $maxH }
    }

    $pixfmt = if ($Alpha) { 'rgba' } else { 'rgb24' }

    # 逐帧规范化为 PNG 序列
    $index = 0
    foreach ($src in $InputFiles) {
        $index++
        $dst = Join-Path $framesDir ("frame_" + $index.ToString("00000") + ".png")

        # 缩放/裁切/留边滤镜
        $vf = switch ($ScaleMode) {
            'contain' { "scale='min($Width,iw)':'min($Height,ih)':force_original_aspect_ratio=decrease,pad=$Width:$Height:(ow-iw)/2:(oh-ih)/2:color=$Background" }
            'cover'   { "scale='max($Width,iw)':'max($Height,ih)':force_original_aspect_ratio=increase,crop=$Width:$Height" }
            'stretch' { "scale=$Width:$Height" }
        }

        $args = @(
            '-y', '-i', "`"$src`"",
            '-vf', $vf,
            '-pix_fmt', $pixfmt,
            "`"$dst`""
        )

        $psi = New-Object System.Diagnostics.ProcessStartInfo
        $psi.FileName = $script:FFmpegPath
        $psi.Arguments = $args -join ' '
        $psi.UseShellExecute = $false
        $psi.CreateNoWindow = $true
        $psi.RedirectStandardError = $true
        $psi.RedirectStandardOutput = $true
        $psi.StandardOutputEncoding = [System.Text.Encoding]::UTF8
        $psi.StandardErrorEncoding  = [System.Text.Encoding]::UTF8

        $p = New-Object System.Diagnostics.Process
        $p.StartInfo = $psi
        $null = $p.Start()
        $p.WaitForExit()

        if ($Debug) {
            $out = $p.StandardOutput.ReadToEnd()
            $err = $p.StandardError.ReadToEnd()
            if ($out) { Write-Output $out }
            if ($err) { Write-Output $err }
        }

        if (-not (Test-Path $dst)) {
            throw "帧预处理失败：$src"
        }
    }

    return @{ WorkDir = $workDir; FramesDir = $framesDir; KeepTemp = [bool]$KeepTemp }
} 