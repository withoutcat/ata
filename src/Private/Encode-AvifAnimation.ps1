function Encode-AvifAnimation {
    param(
        [Parameter(Mandatory=$true)][string]$FramesDir,
        [Parameter(Mandatory=$true)][string]$Output,
        [int]$Fps = 10,
        [switch]$Alpha,
        [int]$Crf = 28,
        [int]$Speed = 5,
        [int]$Threads,
        [switch]$ShowDebug
    )

    $pix = if ($Alpha) { 'yuva420p' } else { 'yuv420p' }

    $pattern = Join-Path $FramesDir 'frame_%05d.png'
    $args = @(
        '-y', '-framerate', $Fps, '-i', "`"$pattern`"",
        '-c:v', 'libaom-av1', '-pix_fmt', $pix,
        '-crf', $Crf, '-b:v', 0, '-cpu-used', $Speed,
        '-row-mt', 1, '-tiles', '2x2'
    )
    if ($Threads) { $args += @('-threads', $Threads) }
    $args += @('-an', '-f', 'avif', "`"$Output`"")

    $psi = New-Object System.Diagnostics.ProcessStartInfo
    $psi.FileName = $script:FFmpegPath
    $psi.Arguments = ($args | ForEach-Object { $_.ToString() }) -join ' '
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

    if ($ShowDebug) {
        $out = $p.StandardOutput.ReadToEnd()
        $err = $p.StandardError.ReadToEnd()
        if ($out) { Write-Output $out }
        if ($err) { Write-Output $err }
    }

    if (-not (Test-Path $Output)) {
        throw "AVIF 动画编码失败"
    }
} 