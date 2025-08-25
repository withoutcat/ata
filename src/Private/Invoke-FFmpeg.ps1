function Invoke-FFmpeg {
    param(
        [string]$InputFullPath,
        [string]$OutputFullPath,
        [bool]$DebugMode
    )
    $ffmpegArgs = "-i `"$InputFullPath`" -c:v libaom-av1 -still-picture 1 -crf 28 -pix_fmt yuv420p `"$OutputFullPath`""

    $processInfo = New-Object System.Diagnostics.ProcessStartInfo
    $processInfo.FileName = "ffmpeg"
    $processInfo.Arguments = $ffmpegArgs
    $processInfo.UseShellExecute = $false
    $processInfo.CreateNoWindow = $true
    $processInfo.RedirectStandardOutput = $true
    $processInfo.RedirectStandardError = $true

    $process = New-Object System.Diagnostics.Process
    $process.StartInfo = $processInfo

    $null = $process.Start()

    $timeout = 30000
    $stdout = ""
    $stderr = ""
    $timedOut = $false

    if ($DebugMode) {
        $process.WaitForExit()
        $stdout = $process.StandardOutput.ReadToEnd()
        $stderr = $process.StandardError.ReadToEnd()
    } else {
        if (-not $process.WaitForExit($timeout)) {
            $process.Kill()
            $timedOut = $true
        }
    }

    return @{ TimedOut = $timedOut; Stdout = $stdout; Stderr = $stderr }
} 