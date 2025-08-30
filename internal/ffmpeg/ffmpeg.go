package ffmpeg

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// FFmpeg路径
var ffmpegPath string

// SetFFmpegPath 设置FFmpeg可执行文件的路径
func SetFFmpegPath(path string) {
	ffmpegPath = path
}

// GetFFmpegPath 获取FFmpeg可执行文件的路径
func GetFFmpegPath() string {
	return ffmpegPath
}

// FindFFmpegPath 查找FFmpeg可执行文件的路径
func FindFFmpegPath() (string, error) {
	// 在系统PATH中查找FFmpeg
	pathFFmpeg, err := exec.LookPath("ffmpeg")
	if err == nil {
		return pathFFmpeg, nil
	}

	return "", errors.New("找不到FFmpeg可执行文件，请确保FFmpeg已安装并添加到系统PATH环境变量中")
}

// ExecuteFFmpeg 执行FFmpeg命令
func ExecuteFFmpeg(args []string) error {
	if ffmpegPath == "" {
		return errors.New("FFmpeg路径未设置")
	}

	// 构建完整命令
	cmd := exec.Command(ffmpegPath, args...)

	// 执行命令
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("启动FFmpeg失败: %v", err)
	}

	// 设置超时
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	// 等待命令完成或超时
	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("FFmpeg执行失败: %v", err)
		}
		return nil
	case <-time.After(30 * time.Minute): // 设置30分钟超时
		// 尝试终止进程
		cmd.Process.Kill()
		return errors.New("FFmpeg执行超时")
	}
}

// GetFFprobeInfo 使用FFprobe获取媒体文件信息
func GetFFprobeInfo(inputPath string, debugMode bool) (string, error) {
	if ffmpegPath == "" {
		return "", errors.New("FFmpeg路径未设置")
	}

	// 从ffmpeg路径推导ffprobe路径
	ffprobePath := strings.Replace(ffmpegPath, "ffmpeg", "ffprobe", 1)
	if _, err := os.Stat(ffprobePath); err != nil {
		return "", fmt.Errorf("找不到FFprobe: %v", err)
	}

	// 构建FFprobe命令
	args := []string{
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height",
		"-of", "csv=p=0",
		inputPath,
	}

	cmd := exec.Command(ffprobePath, args...)
	var output []byte
	var err error

	if debugMode {
		fmt.Println("执行命令:", ffprobePath, strings.Join(args, " "))
		cmd.Stderr = os.Stderr
		output, err = cmd.Output()
	} else {
		output, err = cmd.CombinedOutput()
	}

	if err != nil {
		return "", fmt.Errorf("FFprobe执行失败: %v, 输出: %s", err, output)
	}

	return strings.TrimSpace(string(output)), nil
}