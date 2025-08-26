package cli

import (
	"fmt"
)

// 支持的图像格式
var SupportedImageExtensions = []string{
	".jpg", ".jpeg", ".png", ".webp", ".gif", ".tiff", ".tif", ".bmp",
}

// ShowHelp 显示帮助信息
func ShowHelp() {
	fmt.Println("ATA - AVIF图像转换工具")
	fmt.Println("用法:")
	fmt.Println("  ata [路径] [选项]        - 将指定路径下的图像转换为AVIF格式")
	fmt.Println("  ata convert [路径] [选项] - 将指定路径下的图像转换为AVIF格式")
	fmt.Println("  ata ani [路径] [输出文件] [选项] - 从图像序列创建AVIF动画")
	fmt.Println("  ata ppt [路径] [输出文件] [选项] - 从图像创建幻灯片AVIF动画（低帧率）")
	fmt.Println("  ata help - 显示此帮助信息")
	fmt.Println("")
	fmt.Println("选项:")
	fmt.Println("  -d        启用调试模式")
	fmt.Println("  -r        删除原始文件")
	fmt.Println("  -f        强制覆盖已存在的文件")
	fmt.Println("  -s        递归处理子目录")
	fmt.Println("")
	fmt.Println("动画选项:")
	fmt.Println("  -fps N    设置帧率 (默认: ani=10, ppt=1)")
	fmt.Println("  -crf N    设置质量 (0-63, 越低质量越好, 默认: 30)")
	fmt.Println("  -speed N  设置编码速度 (0-10, 越高越快, 默认: 8)")
	fmt.Println("  -threads N 设置线程数 (0=自动, 默认: 0)")
	fmt.Println("  -alpha    保留透明通道")
	fmt.Println("  -width N  设置输出宽度")
	fmt.Println("  -height N 设置输出高度")
	fmt.Println("  -scale N  设置缩放比例 (默认: 1.0)")
	fmt.Println("  -bg COLOR 设置背景颜色 (默认: white)")
	fmt.Println("")
	fmt.Println("示例:")
	fmt.Println("  ata ./images -s -d       - 递归转换images目录下的所有图像，启用调试模式")
	fmt.Println("  ata ani ./frames output.avif -fps 24 -crf 20 - 从frames目录创建24fps的高质量动画")
	fmt.Println("  ata ppt ./slides presentation.avif - 从slides目录创建幻灯片动画")
	fmt.Println("")
	fmt.Println("支持的图像格式:")
	fmt.Print("  ")
	for i, ext := range SupportedImageExtensions {
		fmt.Print(ext)
		if i < len(SupportedImageExtensions)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Println("")
}