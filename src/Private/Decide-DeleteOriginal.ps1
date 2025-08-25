function Decide-DeleteOriginal {
    param(
        [switch]$Delete
    )
    $deleteOriginal = $Delete
    if (-not $Delete) {
        $choice = Read-Host "是否删除原图片？(Y 删除 / N 保留)"
        if ($choice.ToUpper() -eq "Y") { 
            $deleteOriginal = $true 
        } else {
            $deleteOriginal = $false
        }
    }
    return $deleteOriginal
} 