function Get-TargetFiles {
    param(
        [string]$TargetDir,
        [bool]$DoRecursive,
        [string[]]$SupportedImageExtensions
    )
    if ($DoRecursive) {
        return Get-ChildItem -Path $TargetDir -Recurse -File -Depth 3 | Where-Object {
            $SupportedImageExtensions -contains $_.Extension.ToLower()
        }
    } else {
        return Get-ChildItem -Path $TargetDir -File | Where-Object {
            $SupportedImageExtensions -contains $_.Extension.ToLower()
        }
    }
} 