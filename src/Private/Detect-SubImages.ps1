function Detect-SubImages {
    param(
        [string]$TargetDir,
        [string[]]$IncludePatterns
    )
    $subDirs = Get-ChildItem -Path $TargetDir -Directory
    $hasSubImages = $false
    foreach ($d in $subDirs) {
        if ((Get-ChildItem -Path $d.FullName -Include $IncludePatterns -File -Recurse | Measure-Object).Count -gt 0) {
            $hasSubImages = $true
            break
        }
    }
    return @{ HasSubImages = $hasSubImages; SubDirs = $subDirs }
} 