package installer

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// ShowInteractiveMenu 显示交互式菜单
func ShowInteractiveMenu() {
	for {
		fmt.Println("=== ATA - AVIF图像转换工具 ===")
		fmt.Println("")
		fmt.Println("请选择操作:")
		fmt.Println("1. Install - 安装ATA到系统环境")
		fmt.Println("2. Help - 显示帮助信息")
		fmt.Println("3. Exit - 退出")
		fmt.Println("")
		fmt.Print("请输入选项 (1/2/3): ")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("读取输入失败: %v\n", err)
			os.Exit(1)
		}
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			RunInstaller()
			return
		case "2":
			ShowHelp()
			// ShowHelp会回到主菜单，所以这里continue
		case "3":
			fmt.Println("再见！")
			os.Exit(0)
		case "":
			// 空输入，退出程序
			fmt.Println("\n再见！")
			os.Exit(0)
		default:
			fmt.Printf("无效选项: '%s'，请输入 1、2 或 3\n\n", input)
		}
	}
}

// RunInstaller 运行安装程序
func RunInstaller() {
	fmt.Println("\n=== ATA 安装程序 ===")
	fmt.Println("")

	// 步骤1: 检查FFmpeg
	fmt.Println("步骤 1/2: 检查FFmpeg依赖...")
	if !checkFFmpeg() {
		fmt.Println("FFmpeg未安装，正在安装...")
		if !installFFmpeg() {
			fmt.Println("FFmpeg安装失败，请手动安装后重试")
			return
		}
	} else {
		fmt.Println("✓ FFmpeg已安装")
	}

	// 步骤2: 设置环境变量
	fmt.Println("\n步骤 2/2: 设置环境变量...")
	if !setupEnvironment() {
		fmt.Println("环境变量设置失败")
		return
	}

	fmt.Println("\n✓ 安装完成！")
	fmt.Println("现在您可以在任意位置使用 'ata' 命令了")
	fmt.Println("请重新打开终端以使环境变量生效")
}

// checkFFmpeg 检查FFmpeg是否已安装
func checkFFmpeg() bool {
	_, err := exec.LookPath("ffmpeg")
	return err == nil
}

// installFFmpeg 安装FFmpeg
func installFFmpeg() bool {
	fmt.Printf("正在为 %s 平台安装FFmpeg...\n", runtime.GOOS)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		// 检查是否有chocolatey
		if _, err := exec.LookPath("choco"); err != nil {
			fmt.Println("未检测到Chocolatey包管理器")
			fmt.Println("请手动安装FFmpeg或先安装Chocolatey:")
			fmt.Println("https://chocolatey.org/install")
			return false
		}
		cmd = exec.Command("choco", "install", "ffmpeg", "-y")
	case "darwin":
		// 检查是否有brew
		if _, err := exec.LookPath("brew"); err != nil {
			fmt.Println("未检测到Homebrew包管理器")
			fmt.Println("请手动安装FFmpeg或先安装Homebrew:")
			fmt.Println("https://brew.sh")
			return false
		}
		cmd = exec.Command("brew", "install", "ffmpeg")
	case "linux":
		// 尝试不同的包管理器
		if _, err := exec.LookPath("apt"); err == nil {
			cmd = exec.Command("sudo", "apt", "update", "&&", "sudo", "apt", "install", "-y", "ffmpeg")
		} else if _, err := exec.LookPath("yum"); err == nil {
			cmd = exec.Command("sudo", "yum", "install", "-y", "ffmpeg")
		} else if _, err := exec.LookPath("pacman"); err == nil {
			cmd = exec.Command("sudo", "pacman", "-S", "--noconfirm", "ffmpeg")
		} else {
			fmt.Println("未检测到支持的包管理器")
			fmt.Println("请手动安装FFmpeg")
			return false
		}
	default:
		fmt.Printf("不支持的操作系统: %s\n", runtime.GOOS)
		return false
	}

	fmt.Println("正在执行安装命令，这可能需要几分钟...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("安装失败: %v\n", err)
		return false
	}

	// 再次检查是否安装成功
	return checkFFmpeg()
}

// checkATAInPath 检查PATH中是否已存在ata命令
func checkATAInPath() (bool, string) {
	ataPath, err := exec.LookPath("ata")
	if err != nil {
		return false, ""
	}
	return true, ataPath
}

// setupEnvironment 设置环境变量
func setupEnvironment() bool {
	// 获取当前exe的目录
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("获取可执行文件路径失败: %v\n", err)
		return false
	}
	exeDir := filepath.Dir(exePath)

	// 检查是否已存在ata命令
	exists, existingPath := checkATAInPath()
	if exists {
		fmt.Printf("检测到系统中已存在ata命令: %s\n", existingPath)
		fmt.Print("是否覆盖现有安装? (y/N): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input != "y" && input != "yes" {
			fmt.Println("安装已取消")
			return false
		}
	}

	// 添加到PATH
	return addToPath(exeDir)
}

// addToPath 将目录添加到系统PATH
func addToPath(dir string) bool {
	switch runtime.GOOS {
	case "windows":
		return addToPathWindows(dir)
	case "darwin", "linux":
		return addToPathUnix(dir)
	default:
		fmt.Printf("不支持的操作系统: %s\n", runtime.GOOS)
		return false
	}
}

// addToPathWindows 在Windows上添加到PATH
func addToPathWindows(dir string) bool {
	// 使用setx命令设置用户环境变量
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf(
		"$path = [Environment]::GetEnvironmentVariable('PATH', 'User'); "+
			"if ($path -notlike '*%s*') { "+
			"[Environment]::SetEnvironmentVariable('PATH', $path + ';%s', 'User') }",
		dir, dir))

	err := cmd.Run()
	if err != nil {
		fmt.Printf("设置环境变量失败: %v\n", err)
		return false
	}

	fmt.Printf("✓ 已将 %s 添加到用户PATH\n", dir)
	return true
}

// addToPathUnix 在Unix系统上添加到PATH
func addToPathUnix(dir string) bool {
	// 检查shell类型并添加到相应的配置文件
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("获取用户目录失败: %v\n", err)
		return false
	}

	// 尝试不同的shell配置文件
	shellConfigs := []string{".bashrc", ".zshrc", ".profile"}
	var configFile string

	for _, config := range shellConfigs {
		path := filepath.Join(homeDir, config)
		if _, err := os.Stat(path); err == nil {
			configFile = path
			break
		}
	}

	if configFile == "" {
		// 如果没有找到现有配置文件，创建.bashrc
		configFile = filepath.Join(homeDir, ".bashrc")
	}

	// 检查是否已经添加过
	content, err := os.ReadFile(configFile)
	if err == nil && strings.Contains(string(content), dir) {
		fmt.Printf("✓ PATH中已包含 %s\n", dir)
		return true
	}

	// 添加到配置文件
	file, err := os.OpenFile(configFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("打开配置文件失败: %v\n", err)
		return false
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("\n# Added by ATA installer\nexport PATH=\"$PATH:%s\"\n", dir))
	if err != nil {
		fmt.Printf("写入配置文件失败: %v\n", err)
		return false
	}

	fmt.Printf("✓ 已将 %s 添加到 %s\n", dir, configFile)
	return true
}

// ShowHelp 显示帮助信息
func ShowHelp() {
	fmt.Println("\nATA - AVIF图像转换工具")
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
	fmt.Println("  ata -s -d ./images       - 递归转换images目录下的所有图像，启用调试模式")
	fmt.Println("  ata convert -f -d ./photos - 强制转换photos目录下的图像，启用调试模式")
	fmt.Println("  ata ani -fps 24 -crf 20 ./frames output.avif - 从frames目录创建24fps的高质量动画")
	fmt.Println("  ata ppt -fps 1 ./slides presentation.avif - 从slides目录创建幻灯片动画")
	fmt.Println("")
	fmt.Println("支持的图像格式:")
	fmt.Println("  .jpg, .jpeg, .png, .webp, .gif, .tiff, .tif, .bmp")
	fmt.Println("")
	fmt.Print("按回车键返回主菜单...")
	bufio.NewReader(os.Stdin).ReadString('\n')
	ShowInteractiveMenu()
}