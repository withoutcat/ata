function Confirm-ContinueLargeCount {
    param(
        [int]$TotalFiles,
        [switch]$Force
    )
    if ($TotalFiles -gt 200 -and -not $Force) {
        $choice = Read-Host "发现 $TotalFiles 张可转换图片，是否继续？(Y/N)"
        if ($choice.ToUpper() -ne "Y") { return $false }
    }
    return $true
} 