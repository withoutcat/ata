package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/withoutcat/ata/internal/converter"
	"github.com/withoutcat/ata/internal/ffmpeg"
	"github.com/withoutcat/ata/pkg/cli"
)

// 版本信息，在构建时通过 -ldflags 注入
var version = "dev"

func main() {
	// 检查FFmpeg依赖
	ffmpegPath, err := ffmpeg.FindFFmpegPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		fmt.Println("提示: 请先运行ata-installer-windows.exe来安装FFmpeg依赖")
		os.Exit(1)
	}
	// 设置FFmpeg路径
	ffmpeg.SetFFmpegPath(ffmpegPath)
	
	// 检查是否有命令行参数
	if len(os.Args) > 1 {
		handleCommandLine()
	} else {
		// 没有参数时显示帮助
		cli.ShowHelpWithVersion(version)
	}
}



func handleCommandLine() {
	// 定义子命令
	convertCmd := flag.NewFlagSet("convert", flag.ExitOnError)
	aniCmd := flag.NewFlagSet("ani", flag.ExitOnError)
	pptCmd := flag.NewFlagSet("ppt", flag.ExitOnError)

	// convert子命令的选项
	var convertDebugMode, convertDeleteOriginal, convertRecursive, convertForce bool
	convertCmd.BoolVar(&convertDebugMode, "d", false, "启用调试模式")
	convertCmd.BoolVar(&convertDeleteOriginal, "r", false, "删除原始文件")
	convertCmd.BoolVar(&convertRecursive, "s", false, "递归处理子目录")
	convertCmd.BoolVar(&convertForce, "f", false, "强制覆盖已存在的文件")

	// ani子命令的选项
	var aniDebugMode, aniDeleteOriginal, aniForce bool
	var fps int
	var crf int
	var speed int
	var threads int
	var alpha bool
	var width, height int
	var scale float64
	var background string

	aniCmd.BoolVar(&aniDebugMode, "d", false, "启用调试模式")
	aniCmd.BoolVar(&aniDeleteOriginal, "r", false, "删除原始文件")
	aniCmd.BoolVar(&aniForce, "f", false, "强制覆盖已存在的文件")
	aniCmd.IntVar(&fps, "fps", 10, "帧率")
	aniCmd.IntVar(&crf, "crf", 30, "质量 (0-63, 越低质量越好)")
	aniCmd.IntVar(&speed, "speed", 8, "编码速度 (0-10, 越高越快)")
	aniCmd.IntVar(&threads, "threads", 0, "线程数 (0=自动)")
	aniCmd.BoolVar(&alpha, "alpha", false, "保留透明通道")
	aniCmd.IntVar(&width, "width", 0, "输出宽度")
	aniCmd.IntVar(&height, "height", 0, "输出高度")
	aniCmd.Float64Var(&scale, "scale", 1.0, "缩放比例")
	aniCmd.StringVar(&background, "bg", "white", "背景颜色")

	// ppt子命令的选项
	var pptDebugMode, pptDeleteOriginal, pptForce bool
	var pptFps, pptCrf, pptSpeed, pptThreads int
	var pptAlpha bool
	var pptWidth, pptHeight int
	var pptScale float64
	var pptBackground string

	pptCmd.BoolVar(&pptDebugMode, "d", false, "启用调试模式")
	pptCmd.BoolVar(&pptDeleteOriginal, "r", false, "删除原始文件")
	pptCmd.BoolVar(&pptForce, "f", false, "强制覆盖已存在的文件")
	pptCmd.IntVar(&pptFps, "fps", 1, "帧率")
	pptCmd.IntVar(&pptCrf, "crf", 30, "质量 (0-63, 越低质量越好)")
	pptCmd.IntVar(&pptSpeed, "speed", 8, "编码速度 (0-10, 越高越快)")
	pptCmd.IntVar(&pptThreads, "threads", 0, "线程数 (0=自动)")
	pptCmd.BoolVar(&pptAlpha, "alpha", false, "保留透明通道")
	pptCmd.IntVar(&pptWidth, "width", 0, "输出宽度")
	pptCmd.IntVar(&pptHeight, "height", 0, "输出高度")
	pptCmd.Float64Var(&pptScale, "scale", 1.0, "缩放比例")
	pptCmd.StringVar(&pptBackground, "bg", "white", "背景颜色")

	// 根据第一个参数选择子命令
	switch os.Args[1] {
	case "convert":
		convertCmd.Parse(os.Args[2:])
		if convertCmd.NArg() < 1 {
			fmt.Println("错误: 请指定输入路径")
			cli.ShowHelpWithVersion(version)
			os.Exit(1)
		}
		path := convertCmd.Arg(0)
		converter.ConvertImages(path, convertDebugMode, convertDeleteOriginal, convertRecursive, convertForce)

	case "ani":
		aniCmd.Parse(os.Args[2:])
		if aniCmd.NArg() < 1 {
			fmt.Println("错误: 请指定输入路径")
			cli.ShowHelpWithVersion(version)
			os.Exit(1)
		}
		path := aniCmd.Arg(0)
		outputPath := ""
		if aniCmd.NArg() > 1 {
			outputPath = aniCmd.Arg(1)
		} else {
			// 默认输出路径为输入目录下的output.avif
			outputPath = filepath.Join(path, "output.avif")
		}
		converter.CreateAnimation(path, outputPath, fps, crf, speed, threads, alpha, width, height, scale, background, aniDebugMode, aniDeleteOriginal, aniForce)

	case "ppt":
		pptCmd.Parse(os.Args[2:])
		if pptCmd.NArg() < 1 {
			fmt.Println("错误: 请指定输入路径")
			cli.ShowHelpWithVersion(version)
			os.Exit(1)
		}
		path := pptCmd.Arg(0)
		outputPath := ""
		if pptCmd.NArg() > 1 {
			outputPath = pptCmd.Arg(1)
		} else {
			// 默认输出路径为输入目录下的output.avif
			outputPath = filepath.Join(path, "output.avif")
		}
		// PPT模式默认帧率为1
		converter.CreateAnimation(path, outputPath, pptFps, pptCrf, pptSpeed, pptThreads, pptAlpha, pptWidth, pptHeight, pptScale, pptBackground, pptDebugMode, pptDeleteOriginal, pptForce)

	case "help", "-h", "--help":
		cli.ShowHelpWithVersion(version)

	case "version", "-v", "--version":
		fmt.Printf("ATA v%s\n", version)

	default:
		// 如果第一个参数不是子命令，则假定为路径，使用convert子命令
		path := os.Args[1]
		// 解析剩余参数
		var defaultDebugMode, defaultDeleteOriginal, defaultRecursive, defaultForce bool
		for i := 2; i < len(os.Args); i++ {
			switch os.Args[i] {
			case "-d":
				defaultDebugMode = true
			case "-r":
				defaultDeleteOriginal = true
			case "-s":
				defaultRecursive = true
			case "-f":
				defaultForce = true
			}
		}
		converter.ConvertImages(path, defaultDebugMode, defaultDeleteOriginal, defaultRecursive, defaultForce)
	}
}
