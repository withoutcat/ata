# Dot-source Private helpers first
. $PSScriptRoot/../src/Private/SupportedExtensions.ps1
. $PSScriptRoot/../src/Private/Get-UniqueOutputPath.ps1
. $PSScriptRoot/../src/Private/Get-TargetFiles.ps1
. $PSScriptRoot/../src/Private/Detect-SubImages.ps1
. $PSScriptRoot/../src/Private/Confirm-Recursive.ps1
. $PSScriptRoot/../src/Private/Decide-DeleteOriginal.ps1
. $PSScriptRoot/../src/Private/Confirm-ContinueLargeCount.ps1
. $PSScriptRoot/../src/Private/Invoke-FFmpeg.ps1

# Dot-source Public functions
. $PSScriptRoot/../src/Public/Show-Help.ps1
. $PSScriptRoot/../src/Public/Convert-Images.ps1

Export-ModuleMember -Function Show-Help, Convert-Images 