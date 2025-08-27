package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// ANSI 颜色代码
var (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorGray   = "\033[90m"
)

// 检查是否支持颜色输出
func supportsColor() bool {
	// 在Windows上，检查是否有TERM环境变量或者是否在支持ANSI的终端中
	if runtime.GOOS == "windows" {
		// 简单检查，如果有TERM环境变量通常表示支持颜色
		if os.Getenv("TERM") != "" {
			return true
		}
		// 检查是否在Windows Terminal或其他支持ANSI的终端中
		if os.Getenv("WT_SESSION") != "" || os.Getenv("COLORTERM") != "" {
			return true
		}
		return false
	}
	return true // 非Windows系统通常支持颜色
}

// 初始化颜色支持
func initColors() {
	if !supportsColor() {
		// 如果不支持颜色，将所有颜色代码设为空字符串
		ColorReset = ""
		ColorRed = ""
		ColorGreen = ""
		ColorYellow = ""
		ColorBlue = ""
		ColorPurple = ""
		ColorCyan = ""
		ColorWhite = ""
		ColorGray = ""
	}
}

// Logger 结构体
type Logger struct {
	debugMode bool
	fileCounter int
	totalFiles int
	successCount int
	failureCount int
	showProgress bool
	lastProgressLen int
}

// 全局logger实例
var globalLogger *Logger

// Init 初始化logger
func Init(debugMode bool) {
	// 初始化颜色支持
	initColors()
	
	globalLogger = &Logger{
		debugMode: debugMode,
		fileCounter: 0,
		totalFiles: 0,
		successCount: 0,
		failureCount: 0,
		showProgress: true,
		lastProgressLen: 0,
	}
	
	// 设置标准log包的输出格式
	log.SetFlags(0) // 不显示时间戳，我们自己控制格式
}

// Debug 输出调试信息（仅在调试模式下显示）
func Debug(format string, args ...interface{}) {
	if globalLogger != nil && globalLogger.debugMode {
		msg := fmt.Sprintf(format, args...)
		fmt.Printf("%s[DEBUG]%s %s\n", ColorGray, ColorReset, msg)
	}
}

// Info 输出普通信息
func Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s[INFO]%s %s\n", ColorBlue, ColorReset, msg)
}

// Success 输出成功信息（绿色）
func Success(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s✓%s %s\n", ColorGreen, ColorReset, msg)
}

// Error 输出错误信息（红色）
func Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "%s✗%s %s\n", ColorRed, ColorReset, msg)
}

// Warning 输出警告信息（黄色）
func Warning(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s⚠%s %s\n", ColorYellow, ColorReset, msg)
}

// ProcessStart 开始处理文件（返回序号）
func ProcessStart(filePath string) int {
	if globalLogger == nil {
		return 0
	}
	
	globalLogger.fileCounter++
	fileName := filepath.Base(filePath)
	
	fmt.Printf("%s%d.%s %s", ColorCyan, globalLogger.fileCounter, ColorReset, fileName)
	return globalLogger.fileCounter
}

// ProcessSuccess 处理成功
func ProcessSuccess(duration time.Duration) {
	if globalLogger != nil {
		globalLogger.successCount++
	}
	fmt.Printf(" %s成功%s (耗时: %.2f秒)", ColorGreen, ColorReset, duration.Seconds())
	
	// 如果有多个文件，显示进度条；否则直接换行
	if globalLogger != nil && globalLogger.totalFiles > 1 {
		fmt.Print("\n")
		ShowProgress() // 每次文件处理完成都更新进度条
	} else {
		fmt.Print("\n")
	}
}

// ProcessError 处理失败
func ProcessError(err error, duration time.Duration) {
	if globalLogger != nil {
		globalLogger.failureCount++
	}
	fmt.Printf(" %s失败%s (耗时: %.2f秒) - %v", ColorRed, ColorReset, duration.Seconds(), err)
	
	// 如果有多个文件，显示进度条；否则直接换行
	if globalLogger != nil && globalLogger.totalFiles > 1 {
		fmt.Print("\n")
		ShowProgress() // 每次文件处理完成都更新进度条
	} else {
		fmt.Print("\n")
	}
}

// GetFileCounter 获取当前文件计数器
func GetFileCounter() int {
	if globalLogger == nil {
		return 0
	}
	return globalLogger.fileCounter
}

// ResetCounter 重置文件计数器
func ResetCounter() {
	if globalLogger != nil {
		globalLogger.fileCounter = 0
		globalLogger.successCount = 0
		globalLogger.failureCount = 0
		globalLogger.lastProgressLen = 0
	}
}

// SetTotalFiles 设置总文件数
func SetTotalFiles(total int) {
	if globalLogger != nil {
		globalLogger.totalFiles = total
	}
}

// ShowStartSummary 显示开始处理的摘要
func ShowStartSummary(total int) {
	SetTotalFiles(total)
	if total == 0 {
		Info("未找到可处理的文件")
		return
	}
	fmt.Printf("%s找到 %d 个文件，开始处理...%s\n", ColorCyan, total, ColorReset)
	
	// 如果有多个文件，显示初始进度条
	if total > 1 {
		ShowProgress()
	}
}

// ShowProgress 显示动态进度条
func ShowProgress() {
	if globalLogger == nil || !globalLogger.showProgress || globalLogger.totalFiles <= 1 {
		return
	}
	

	
	processed := globalLogger.successCount + globalLogger.failureCount
	percentage := float64(processed) / float64(globalLogger.totalFiles) * 100
	
	// 进度条长度
	barWidth := 25
	filledWidth := int(float64(barWidth) * float64(processed) / float64(globalLogger.totalFiles))
	
	// 使用Unicode字符构建更美观的进度条
	// █ (U+2588) 实心方块 - 已完成部分
	// ░ (U+2591) 浅色方块 - 未完成部分
	// 这些字符在支持连字的字体中会无缝连接
	var bar string
	for i := 0; i < barWidth; i++ {
		if i < filledWidth {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	
	// 构建进度信息，使用颜色让进度条更醒目
	progressText := fmt.Sprintf("\r%s进度: %s%s%s %.1f%% (%d/%d)%s", 
		ColorCyan, ColorGreen, bar, ColorCyan, percentage, processed, globalLogger.totalFiles, ColorReset)
	
	// 清除之前的进度条（如果更短）
	if len(progressText) < globalLogger.lastProgressLen {
		// 用空格填充到之前的长度
		for i := len(progressText); i < globalLogger.lastProgressLen; i++ {
			progressText += " "
		}
	}
	
	fmt.Print(progressText)
	globalLogger.lastProgressLen = len(progressText)
}

// ClearProgress 清除进度条
func ClearProgress() {
	if globalLogger == nil {
		return
	}
	

	
	// 清除基于print的进度条
	if globalLogger.lastProgressLen == 0 {
		return
	}
	
	// 用空格覆盖进度条
	clearText := "\r"
	for i := 0; i < globalLogger.lastProgressLen; i++ {
		clearText += " "
	}
	clearText += "\r"
	fmt.Print(clearText)
	globalLogger.lastProgressLen = 0
}

// SetProgressEnabled 设置是否显示进度条
func SetProgressEnabled(enabled bool) {
	if globalLogger != nil {
		globalLogger.showProgress = enabled
	}
}

// ShowFinalSummary 显示最终处理结果摘要
func ShowFinalSummary() {
	if globalLogger == nil {
		return
	}
	
	if globalLogger.totalFiles == 0 {
		return
	}
	
	// 如果有多个文件，先显示100%完成的进度条
	if globalLogger.totalFiles > 1 {
		// 构建100%完成的进度条
		barWidth := 25
		var bar string
		for i := 0; i < barWidth; i++ {
			bar += "█"
		}
		progressText := fmt.Sprintf("\r%s进度: %s%s%s 100.0%% (%d/%d)%s", 
			ColorCyan, ColorGreen, bar, ColorCyan, globalLogger.totalFiles, globalLogger.totalFiles, ColorReset)
		fmt.Print(progressText)
		
		// 等待一小段时间让用户看到100%完成
		time.Sleep(500 * time.Millisecond)
		
		// 清除进度条
		ClearProgress()
	}
	
	fmt.Printf("\n%s处理完成！%s\n", ColorCyan, ColorReset)
	fmt.Printf("总计: %d 个文件\n", globalLogger.totalFiles)
	fmt.Printf("%s成功: %d%s\n", ColorGreen, globalLogger.successCount, ColorReset)
	if globalLogger.failureCount > 0 {
		fmt.Printf("%s失败: %d%s\n", ColorRed, globalLogger.failureCount, ColorReset)
	}
}