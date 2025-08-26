package converter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/withoutcat/ata/internal/ffmpeg"
	"github.com/withoutcat/ata/pkg/cli"
)

// 检查文件是否为支持的图像格式
func isSupportedImageFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	for _, supportedExt := range cli.SupportedImageExtensions {
		if ext == supportedExt {
			return true
		}
	}
	return false
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
func ConvertImages(path string, debugMode, deleteOriginal, recursive, force bool) {
	// 获取文件信息
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: 无法访问路径 %s: %v\n", path, err)
		return
	}

	// 如果是目录，则处理目录中的所有文件
	if fileInfo.IsDir() {
		processDirectory(path, debugMode, deleteOriginal, recursive, force)
	} else {
		// 如果是文件，则直接处理该文件
		if isSupportedImageFile(path) {
			processFile(path, debugMode, deleteOriginal, force)
		} else {
			fmt.Fprintf(os.Stderr, "错误: 不支持的文件格式: %s\n", path)
		}
	}
}

// 处理目录中的所有文件
func processDirectory(dirPath string, debugMode, deleteOriginal, recursive, force bool) {
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
				processDirectory(path, debugMode, deleteOriginal, recursive, force)
			}
		} else {
			// 如果是文件，则检查是否为支持的图像格式
			if isSupportedImageFile(path) {
				processFile(path, debugMode, deleteOriginal, force)
			}
		}
	}
}

// 处理单个文件
func processFile(filePath string, debugMode, deleteOriginal, force bool) {
	// 获取输出路径
	outputPath, err := getOutputPath(filePath, force)
if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		return
	}

	fmt.Printf("转换: %s -> %s\n", filePath, outputPath)

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
	err = ffmpeg.ExecuteFFmpeg(args, debugMode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: 转换失败 %s: %v\n", filePath, err)
		// 如果转换失败，删除可能部分生成的输出文件
		os.Remove(outputPath)
		return
	}

	// 如果需要删除原始文件
	if deleteOriginal {
		err = os.Remove(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "警告: 无法删除原始文件 %s: %v\n", filePath, err)
		}
	}
}