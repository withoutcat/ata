package converter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/withoutcat/ata/internal/ffmpeg"
	"github.com/withoutcat/ata/internal/logger"
)

// 支持的图像格式(other to avif)
var SupportedImageExtensionsForConvertAvif = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
	".tiff": true,
	".tif":  true,
	".bmp":  true,
}

// 检查文件是否为支持的图像格式
func isSupportedImageFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	_, ok := SupportedImageExtensionsForConvertAvif[ext]
	return ok
}

// countSupportedFiles 统计目录中支持的图像文件数量
func countSupportedFiles(dirPath string, recursive bool) int {
	count := 0
	
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return 0
	}
	
	for _, entry := range entries {
		filePath := filepath.Join(dirPath, entry.Name())
		
		if entry.IsDir() {
			if recursive {
				count += countSupportedFiles(filePath, recursive)
			}
		} else {
			if isSupportedImageFile(filePath) {
				count++
			}
		}
	}
	
	return count
}

// 获取AVIF输出路径
func getOutputPath(inputPath string, force bool) (string, error) {
	// 获取不带扩展名的文件名
	base := filepath.Base(inputPath)
	dir := filepath.Dir(inputPath)
	fileNameWithoutExt := strings.TrimSuffix(base, filepath.Ext(base))

	// 构建输出路径
	outputPath := filepath.Join(dir, fileNameWithoutExt+".avif")

	// 检查输出文件是否已存在
	if !force {
		if _, err := os.Stat(outputPath); err == nil {
			return "", fmt.Errorf("输出文件已存在: %s，使用-f选项强制覆盖", outputPath)
		}
	}

	return outputPath, nil
}

// ConvertImages 将指定路径下的图像转换为AVIF格式
func ConvertImages(path string, deleteOriginal, recursive, force bool) {
	// 初始化logger
	logger.Init()
	logger.ResetCounter()
	
	// 获取文件信息
	fileInfo, err := os.Stat(path)
	if err != nil {
		logger.Error("无法访问路径 %s: %v", path, err)
		return
	}

	// 统计可处理的文件数量
	var totalFiles int
	if fileInfo.IsDir() {
		totalFiles = countSupportedFiles(path, recursive)
	} else {
		if isSupportedImageFile(path) {
			totalFiles = 1
		} else {
			totalFiles = 0
		}
	}
	
	// 显示开始处理的摘要
	logger.ShowStartSummary(totalFiles)
	
	// 如果没有可处理的文件，直接返回
	if totalFiles == 0 {
		return
	}

	// 如果是目录，则处理目录中的所有文件
	if fileInfo.IsDir() {
		processDirectory(path, deleteOriginal, recursive, force)
	} else {
		// 如果是文件，则直接处理该文件
		processFile(path, deleteOriginal, force)
	}
	
	// 显示最终处理结果摘要
	logger.ShowFinalSummary()
}

// 处理目录中的所有文件
func processDirectory(dirPath string, deleteOriginal, recursive, force bool) {
	// 读取目录内容
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: 无法读取目录 %s: %v\n", dirPath, err)
		return
	}

	// 处理目录中的每个条目
	for _, entry := range entries {
		path := filepath.Join(dirPath, entry.Name())

		if entry.IsDir() {
			// 如果是目录且启用了递归，则递归处理
			if recursive {
				processDirectory(path, deleteOriginal, recursive, force)
			}
		} else {
			// 如果是文件，则检查是否为支持的图像格式
			if isSupportedImageFile(path) {
				processFile(path, deleteOriginal, force)
			}
		}
	}
}

// 处理单个文件
func processFile(filePath string, deleteOriginal, force bool) {
	// 获取输出路径
	outputPath, err := getOutputPath(filePath, force)
	if err != nil {
		logger.Error("%v", err)
		return
	}

	// 开始处理文件，显示序号和文件名
	logger.ProcessStart(outputPath)
	
	// 计时器开始
	processStartTime := time.Now()

	// 构建FFmpeg参数
	args := []string{
		"-i", filePath,
		"-y", // 覆盖输出文件
		"-c:v", "libaom-av1",
		"-crf", "30",
		"-b:v", "0",
		"-pix_fmt", "yuv420p", // 使用yuv420p像素格式
		outputPath,
	}

	// 执行FFmpeg命令
	err = ffmpeg.ExecuteFFmpeg(args)
	
	// 计时器结束
	processEndTime := time.Now()
	processDuration := processEndTime.Sub(processStartTime)
	
	if err != nil {
		logger.ProcessError(err, processDuration)
		// 如果转换失败，删除可能部分生成的输出文件
		os.Remove(outputPath)
		return
	}

	// 打印转换成功信息
	logger.ProcessSuccess(nil)

	// 如果需要删除原始文件
	if deleteOriginal {
		err = os.Remove(filePath)
		if err != nil {
			logger.Warning("无法删除原始文件 %s: %v", filePath, err)
		}
	}
}