function Confirm-Recursive {
    param(
        [bool]$HasSubImages,
        [int]$SubDirCount,
        [switch]$Recursive,
        [switch]$Force
    )
    $doRecursive = $false
    if ($HasSubImages -and ($SubDirCount -gt 0)) {
        if ($Recursive -or $Force) {
            $doRecursive = $true
        } else {
            $choice = Read-Host "检测到子文件夹内有可转换图片，是否递归执行？(Y/N)"
            if ($choice.ToUpper() -eq "Y") { $doRecursive = $true }
        }
    }
    return $doRecursive
} 