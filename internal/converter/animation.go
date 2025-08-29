package converter

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/withoutcat/ata/internal/ffmpeg"
	"github.com/withoutcat/ata/internal/logger"
)

// CreateAnimation 从图像序列创建AVIF动画
func CreateAnimation(inputPath, outputPath string, fps, crf, speed, threads int, alpha bool, width, height int, scale float64, background string, debugMode, deleteOriginal, force bool) {
	// 初始化logger
	logger.Init(debugMode)
	logger.ResetCounter()
	
	// 检查输入路径是否存在
	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		logger.Error("无法访问路径 %s: %v", inputPath, err)
		return
	}
	
	// 统计可处理的文件数量
	var totalFiles int
	if fileInfo.IsDir() {
		totalFiles = countAnimationFiles(inputPath)
	} else {
		if isAnimationFile(inputPath) {
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

	// 如果输出路径为空，则使用默认输出路径
	if outputPath == "" {
		outputPath = filepath.Join(inputPath, "output.avif")
	}

	// 检查输出文件是否已存在
	if !force {
		if _, err := os.Stat(outputPath); err == nil {
			logger.Error("输出文件已存在: %s，使用-f选项强制覆盖", outputPath)
			return
		}
	}

	// 创建临时目录用于存放处理后的帧
	tempDir, err := os.MkdirTemp("", "ata-frames-*")
	if err != nil {
		logger.Error("无法创建临时目录: %v", err)
		return
	}
	defer os.RemoveAll(tempDir) // 确保在函数结束时删除临时目录

	// 准备动画帧
	if fileInfo.IsDir() {
		// 如果输入是目录，则处理目录中的所有图像文件
		err = prepareFramesFromDirectory(inputPath, tempDir, width, height, scale, background, debugMode)
	} else {
		// 如果输入是单个文件，则尝试将其作为动画处理
		err = prepareFramesFromAnimation(inputPath, tempDir, width, height, scale, background, debugMode)
	}

	if err != nil {
		logger.Error("准备帧失败: %v", err)
		return
	}

	// 编码AVIF动画
	err = encodeAvifAnimation(tempDir, outputPath, fps, crf, speed, threads, alpha, debugMode)
	if err != nil {
		logger.Error("编码动画失败: %v", err)
		return
	}

	logger.Success("成功创建AVIF动画: %s", outputPath)

	// 如果需要删除原始文件
	if deleteOriginal && fileInfo.IsDir() {
		// 如果输入是目录，则删除目录中的所有图像文件
		entries, err := os.ReadDir(inputPath)
		if err != nil {
			logger.Warning("无法读取目录 %s: %v", inputPath, err)
			return
		}

		for _, entry := range entries {
			if !entry.IsDir() && isSupportedImageFile(entry.Name()) {
				filePath := filepath.Join(inputPath, entry.Name())
				err = os.Remove(filePath)
				if err != nil {
					logger.Warning("无法删除文件 %s: %v", filePath, err)
				}
			}
		}
	} else if deleteOriginal && !fileInfo.IsDir() {
		// 如果输入是单个文件，则删除该文件
		err = os.Remove(inputPath)
		if err != nil {
			logger.Warning("无法删除文件 %s: %v", inputPath, err)
		}
	}
	
	// 显示最终处理结果摘要
	logger.ShowFinalSummary()
}

// 从目录准备帧
func prepareFramesFromDirectory(inputDir, tempDir string, width, height int, scale float64, background string, debugMode bool) error {
	// 读取目录中的所有文件
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		return fmt.Errorf("无法读取目录 %s: %v", inputDir, err)
	}

	// 过滤出支持的图像文件
	var imageFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && isSupportedImageFile(entry.Name()) {
			imageFiles = append(imageFiles, filepath.Join(inputDir, entry.Name()))
		}
	}

	if len(imageFiles) == 0 {
		return fmt.Errorf("目录 %s 中没有支持的图像文件", inputDir)
	}

	// 如果未指定宽度和高度，则从第一个图像获取
	if width == 0 || height == 0 {
		dimensions, err := ffmpeg.GetFFprobeInfo(imageFiles[0], debugMode)
		if err != nil {
			return fmt.Errorf("无法获取图像尺寸: %v", err)
		}

		parts := strings.Split(dimensions, ",")
		if len(parts) != 2 {
			return fmt.Errorf("无法解析图像尺寸: %s", dimensions)
		}

		width, err = strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("无法解析宽度: %v", err)
		}

		height, err = strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("无法解析高度: %v", err)
		}
	}

	// 应用缩放
	if scale != 1.0 {
		width = int(float64(width) * scale)
		height = int(float64(height) * scale)
	}

	// 处理每个图像文件
	for i, imagePath := range imageFiles {
		outputFrame := filepath.Join(tempDir, fmt.Sprintf("frame_%04d.png", i))
		err = processFrame(imagePath, outputFrame, width, height, background, debugMode)
		if err != nil {
			return fmt.Errorf("处理帧 %s 失败: %v", imagePath, err)
		}
		// 更新进度条
		logger.ShowProgress()
	}

	return nil
}

// 从动画文件准备帧
func prepareFramesFromAnimation(inputFile, tempDir string, width, height int, scale float64, background string, debugMode bool) error {
	// 使用FFmpeg提取帧
	logger.Info("正在从动画文件提取帧...")
	args := []string{
		"-i", inputFile,
		"-vsync", "0",
		filepath.Join(tempDir, "frame_%04d.png"),
	}

	logger.Debug("执行FFmpeg命令: %v", args)
	err := ffmpeg.ExecuteFFmpeg(args, debugMode)
	if err != nil {
		return fmt.Errorf("提取帧失败: %v", err)
	}

	logger.Info("帧提取完成")

	// 如果需要调整大小或应用背景，则处理每个提取的帧
	if width > 0 || height > 0 || scale != 1.0 || background != "white" {
		// 读取临时目录中的所有帧
		entries, err := os.ReadDir(tempDir)
		if err != nil {
			return fmt.Errorf("无法读取临时目录: %v", err)
		}

		// 如果未指定宽度和高度，则从第一个帧获取
		if (width == 0 || height == 0) && len(entries) > 0 {
			firstFrame := filepath.Join(tempDir, entries[0].Name())
			dimensions, err := ffmpeg.GetFFprobeInfo(firstFrame, debugMode)
			if err != nil {
				return fmt.Errorf("无法获取帧尺寸: %v", err)
			}

			parts := strings.Split(dimensions, ",")
			if len(parts) != 2 {
				return fmt.Errorf("无法解析帧尺寸: %s", dimensions)
			}

			width, err = strconv.Atoi(parts[0])
			if err != nil {
				return fmt.Errorf("无法解析宽度: %v", err)
			}

			height, err = strconv.Atoi(parts[1])
			if err != nil {
				return fmt.Errorf("无法解析高度: %v", err)
			}
		}

		// 应用缩放
		if scale != 1.0 {
			width = int(float64(width) * scale)
			height = int(float64(height) * scale)
		}

		// 处理每个帧
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasPrefix(entry.Name(), "frame_") {
				framePath := filepath.Join(tempDir, entry.Name())
				tempPath := filepath.Join(tempDir, "temp_"+entry.Name())

				// 处理帧
				err = processFrame(framePath, tempPath, width, height, background, debugMode)
				if err != nil {
					return fmt.Errorf("处理帧 %s 失败: %v", framePath, err)
				}

				// 用处理后的帧替换原始帧
				err = os.Remove(framePath)
				if err != nil {
					return fmt.Errorf("删除原始帧失败: %v", err)
				}

				err = os.Rename(tempPath, framePath)
				if err != nil {
					return fmt.Errorf("重命名处理后的帧失败: %v", err)
				}
			}
		}
	}

	return nil
}

// 处理单个帧
func processFrame(inputPath, outputPath string, width, height int, background string, debugMode bool) error {
	// 构建FFmpeg参数
	args := []string{
		"-i", inputPath,
		"-vf", fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease,pad=%d:%d:(ow-iw)/2:(oh-ih)/2:color=%s", width, height, width, height, background),
		"-y",
		outputPath,
	}

	// 执行FFmpeg命令
	err := ffmpeg.ExecuteFFmpeg(args, debugMode)
	if err != nil {
		return err
	}
	
	// 帧处理成功，更新计数器
	// 使用 ProcessSuccess 方法来更新计数器
	logger.ProcessSuccess(0)
	
	return nil
}

// countAnimationFiles 统计目录中可处理的动画文件数量
func countAnimationFiles(dir string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}

	count := 0
	for _, entry := range entries {
		if !entry.IsDir() && isSupportedImageFileForAnimation(entry.Name()) {
			count++
		}
	}
	return count
}

// isSupportedImageFileForAnimation 检查是否为支持的图像文件（用于动画）
func isSupportedImageFileForAnimation(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".bmp" || ext == ".tiff" || ext == ".tif" || ext == ".webp"
}

// isAnimationFile 检查文件是否为支持的动画文件
func isAnimationFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".gif" || ext == ".webp" || ext == ".mp4" || ext == ".avi" || ext == ".mov" || ext == ".mkv"
}

// 编码AVIF动画
func encodeAvifAnimation(framesDir, outputPath string, fps, crf, speed, threads int, alpha bool, debugMode bool) error {
	// 构建FFmpeg参数
	pixFmt := "yuv420p"
	if alpha {
		pixFmt = "yuva420p"
	}

	logger.Info("正在编码AVIF动画...")

	// 构建FFmpeg命令
	args := []string{
		"-framerate", fmt.Sprintf("%d", fps),
		"-i", filepath.Join(framesDir, "frame_%04d.png"),
	}

	if threads > 0 {
		args = append(args, "-threads", fmt.Sprintf("%d", threads))
	}

	args = append(args,
		"-c:v", "libaom-av1",
		"-crf", fmt.Sprintf("%d", crf),
		"-b:v", "0",
		"-cpu-used", fmt.Sprintf("%d", speed),
		"-pix_fmt", pixFmt,
		"-y",
		outputPath,
	)

	logger.Debug("执行FFmpeg命令: %v", args)
	// 执行FFmpeg命令
	err := ffmpeg.ExecuteFFmpeg(args, debugMode)
	if err != nil {
		return err
	}

	logger.Info("AVIF动画编码完成")
	return nil
}