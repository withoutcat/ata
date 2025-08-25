# Dot-source Private helpers first
. $PSScriptRoot/../src/Private/SupportedExtensions.ps1
. $PSScriptRoot/../src/Private/Get-UniqueOutputPath.ps1
. $PSScriptRoot/../src/Private/Get-TargetFiles.ps1
. $PSScriptRoot/../src/Private/Detect-SubImages.ps1
. $PSScriptRoot/../src/Private/Confirm-Recursive.ps1
. $PSScriptRoot/../src/Private/Decide-DeleteOriginal.ps1
. $PSScriptRoot/../src/Private/Confirm-ContinueLargeCount.ps1
. $PSScriptRoot/../src/Private/Find-FFmpegPath.ps1
. $PSScriptRoot/../src/Private/Invoke-FFmpeg.ps1
. $PSScriptRoot/../src/Private/Prepare-AnimationFrames.ps1
. $PSScriptRoot/../src/Private/Encode-AvifAnimation.ps1

# 解析 ffmpeg 路径并缓存，避免重复查找
$script:FFmpegPath = Find-FFmpegPath

# Dot-source Public functions
. $PSScriptRoot/../src/Public/Show-Help.ps1
. $PSScriptRoot/../src/Public/Convert-Images.ps1
. $PSScriptRoot/../src/Public/Create-AvifAnimation.ps1

Export-ModuleMember -Function Show-Help, Convert-Images, Create-AvifAnimation 