$SupportedImageExtensions = @(
    ".jpg", ".jpeg", ".png", ".webp", ".bmp", ".tiff", ".heic", ".heif"
)
$IncludePatterns = $SupportedImageExtensions | ForEach-Object { "*$_" } 