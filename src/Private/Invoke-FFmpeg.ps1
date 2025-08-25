function Invoke-FFmpeg {
    param(
        [string]$InputFullPath,
        [string]$OutputFullPath,
        [bool]$DebugMode
    )
    $ffmpegPath = $script:FFmpegPath
    if (-not $ffmpegPath) {
        throw "未找到 ffmpeg 可执行文件。请在工程目录下放置 module\\ffmpeg-n7.1-latest-win64-gpl-7.1\\bin\\ffmpeg.exe 或将 ffmpeg 加入 PATH。"
    }

    $ffmpegArgs = "-i `"$InputFullPath`" -c:v libaom-av1 -still-picture 1 -crf 28 -pix_fmt yuv420p `"$OutputFullPath`""

    $processInfo = New-Object System.Diagnostics.ProcessStartInfo
    $processInfo.FileName = $ffmpegPath
    $processInfo.Arguments = $ffmpegArgs
    $processInfo.UseShellExecute = $false
    $processInfo.CreateNoWindow = $true
    $processInfo.RedirectStandardOutput = $true
    $processInfo.RedirectStandardError = $true
    $processInfo.StandardOutputEncoding = [System.Text.Encoding]::UTF8
    $processInfo.StandardErrorEncoding  = [System.Text.Encoding]::UTF8

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