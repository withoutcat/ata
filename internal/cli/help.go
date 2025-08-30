package cli

import (
	"fmt"

	"github.com/withoutcat/ata/internal/converter"
)

// ShowHelp 显示帮助信息
func ShowHelp() {
	ShowHelpWithVersion("dev")
}

// ShowHelpWithVersion 显示带版本信息的帮助
func ShowHelpWithVersion(version string) {
	fmt.Printf("ATA - AVIF图像转换工具 v%s\n", version)
	fmt.Println("用法:")
	fmt.Println("  ata [选项] [路径]        - 将指定路径下的图像转换为AVIF格式")
	fmt.Println("  ata convert [选项] [路径] - 将指定路径下的图像转换为AVIF格式")
	fmt.Println("  ata ani [选项] [路径] [输出文件] - 从图像序列创建AVIF动画")
	fmt.Println("  ata ppt [选项] [路径] [输出文件] - 从图像创建幻灯片AVIF动画（低帧率）")
	fmt.Println("  ata help - 显示此帮助信息")
	fmt.Println("")
	fmt.Println("注意: 选项必须在路径参数之前指定")
	fmt.Println("")
	fmt.Println("选项:")
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
	fmt.Println("  ata -s ./images       - 递归转换images目录下的所有图像")
	fmt.Println("  ata convert -f ./photos - 强制转换photos目录下的图像")
	fmt.Println("  ata ani -fps 24 -crf 20 ./frames output.avif - 从frames目录创建24fps的高质量动画")
	fmt.Println("  ata ppt -fps 1 ./slides presentation.avif - 从slides目录创建幻灯片动画")
	fmt.Println("")
	fmt.Println("支持的图像格式:")
	fmt.Print("  ")
	// 将map的key转换为slice以便排序和遍历
	exts := make([]string, 0, len(converter.SupportedImageExtensionsForConvertAvif))
	for ext := range converter.SupportedImageExtensionsForConvertAvif {
		exts = append(exts, ext)
	}
	
	for i, ext := range exts {
		fmt.Print(ext)
		if i < len(exts)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Println("")
}