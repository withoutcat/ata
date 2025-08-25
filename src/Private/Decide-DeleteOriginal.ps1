function Decide-DeleteOriginal {
    param(
        [switch]$Delete
    )
    $deleteOriginal = $Delete
    if (-not $Delete) {
        $choice = Read-Host "是否删除原图片？(D 删除 / S 保留)"
        if ($choice.ToUpper() -eq "D") { 
            $deleteOriginal = $true 
        } else {
            $deleteOriginal = $false
        }
    }
    return $deleteOriginal
} 