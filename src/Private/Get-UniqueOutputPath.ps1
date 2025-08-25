function Get-UniqueOutputPath {
    param (
        [string]$OutputPath
    )
    $directory = [System.IO.Path]::GetDirectoryName($OutputPath)
    $baseName = [System.IO.Path]::GetFileNameWithoutExtension($OutputPath)
    $extension = [System.IO.Path]::GetExtension($OutputPath)
    $counter = 0
    $newPath = $OutputPath

    while (Test-Path $newPath) {
        $counter++
        $newBaseName = "$baseName ($counter)"
        $newPath = Join-Path $directory "$newBaseName$extension"
    }
    return $newPath
} 