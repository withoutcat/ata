package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/withoutcat/ata/internal/converter"
	"github.com/withoutcat/ata/internal/ffmpeg"
	"github.com/withoutcat/ata/pkg/cli"
	"github.com/withoutcat/ata/pkg/gui"
)

func main() {
	// 检查FFmpeg是否可用
	ffmpegPath, err := ffmpeg.FindFFmpegPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}

	// 设置FFmpeg路径
	ffmpeg.SetFFmpegPath(ffmpegPath)

	// 检查是否有命令行参数
	if len(os.Args) > 1 {
		// 命令行模式
		handleCommandLine()
	} else {
		// GUI模式
		gui.StartGUI()
	}
}

func handleCommandLine() {
	// 定义子命令
	convertCmd := flag.NewFlagSet("convert", flag.ExitOnError)
	aniCmd := flag.NewFlagSet("ani", flag.ExitOnError)
	pptCmd := flag.NewFlagSet("ppt", flag.ExitOnError)
	helpCmd := flag.NewFlagSet("help", flag.ExitOnError)

	// 通用选项
	debugMode := false
	deleteOriginal := false
	recursive := false
	force := false

	// convert子命令的选项
	convertCmd.BoolVar(&debugMode, "d", false, "启用调试模式")
	convertCmd.BoolVar(&deleteOriginal, "r", false, "删除原始文件")
	convertCmd.BoolVar(&recursive, "s", false, "递归处理子目录")
	convertCmd.BoolVar(&force, "f", false, "强制覆盖已存在的文件")

	// ani子命令的选项
	var fps int
	var crf int
	var speed int
	var threads int
	var alpha bool
	var width, height int
	var scale float64
	var background string

	aniCmd.BoolVar(&debugMode, "d", false, "启用调试模式")
	aniCmd.BoolVar(&deleteOriginal, "r", false, "删除原始文件")
	aniCmd.BoolVar(&force, "f", false, "强制覆盖已存在的文件")
	aniCmd.IntVar(&fps, "fps", 10, "帧率")
	aniCmd.IntVar(&crf, "crf", 30, "质量 (0-63, 越低质量越好)")
	aniCmd.IntVar(&speed, "speed", 8, "编码速度 (0-10, 越高越快)")
	aniCmd.IntVar(&threads, "threads", 0, "线程数 (0=自动)")
	aniCmd.BoolVar(&alpha, "alpha", false, "保留透明通道")
	aniCmd.IntVar(&width, "width", 0, "输出宽度")
	aniCmd.IntVar(&height, "height", 0, "输出高度")
	aniCmd.Float64Var(&scale, "scale", 1.0, "缩放比例")
	aniCmd.StringVar(&background, "bg", "white", "背景颜色")

	// ppt子命令的选项与ani相同
	pptCmd.BoolVar(&debugMode, "d", false, "启用调试模式")
	pptCmd.BoolVar(&deleteOriginal, "r", false, "删除原始文件")
	pptCmd.BoolVar(&force, "f", false, "强制覆盖已存在的文件")
	pptCmd.IntVar(&fps, "fps", 1, "帧率")
	pptCmd.IntVar(&crf, "crf", 30, "质量 (0-63, 越低质量越好)")
	pptCmd.IntVar(&speed, "speed", 8, "编码速度 (0-10, 越高越快)")
	pptCmd.IntVar(&threads, "threads", 0, "线程数 (0=自动)")
	pptCmd.BoolVar(&alpha, "alpha", false, "保留透明通道")
	pptCmd.IntVar(&width, "width", 0, "输出宽度")
	pptCmd.IntVar(&height, "height", 0, "输出高度")
	pptCmd.Float64Var(&scale, "scale", 1.0, "缩放比例")
	pptCmd.StringVar(&background, "bg", "white", "背景颜色")

	// 根据第一个参数选择子命令
	switch os.Args[1] {
	case "convert":
		convertCmd.Parse(os.Args[2:])
		if convertCmd.NArg() < 1 {
			fmt.Println("错误: 请指定输入路径")
			cli.ShowHelp()
			os.Exit(1)
		}
		path := convertCmd.Arg(0)
		converter.ConvertImages(path, debugMode, deleteOriginal, recursive, force)

	case "ani":
		aniCmd.Parse(os.Args[2:])
		if aniCmd.NArg() < 1 {
			fmt.Println("错误: 请指定输入路径")
			cli.ShowHelp()
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
		converter.CreateAnimation(path, outputPath, fps, crf, speed, threads, alpha, width, height, scale, background, debugMode, deleteOriginal, force)

	case "ppt":
		pptCmd.Parse(os.Args[2:])
		if pptCmd.NArg() < 1 {
			fmt.Println("错误: 请指定输入路径")
			cli.ShowHelp()
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
		converter.CreateAnimation(path, outputPath, fps, crf, speed, threads, alpha, width, height, scale, background, debugMode, deleteOriginal, force)

	case "help", "-h", "--help":
		cli.ShowHelp()

	default:
		// 如果第一个参数不是子命令，则假定为路径，使用convert子命令
		path := os.Args[1]
		// 解析剩余参数
		for i := 2; i < len(os.Args); i++ {
			switch os.Args[i] {
			case "-d":
				debugMode = true
			case "-r":
				deleteOriginal = true
			case "-s":
				recursive = true
			case "-f":
				force = true
			}
		}
		converter.ConvertImages(path, debugMode, deleteOriginal, recursive, force)
	}
}
