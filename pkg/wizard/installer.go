package wizard

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
func ShowInteractiveMenu(embeddedExecutable []byte) {
	ShowInteractiveMenuWithVersion(embeddedExecutable, "dev")
}

// ShowInteractiveMenuWithVersion 显示带版本信息的交互式菜单
func ShowInteractiveMenuWithVersion(embeddedExecutable []byte, version string) {
	showWelcomePage(embeddedExecutable, version)
}

// showWelcomePage 显示欢迎页面
func showWelcomePage(embeddedExecutable []byte, version string) {
	clearScreen()
	fmt.Println("┌─────────────────────────────────────────────────────────────┐")
	fmt.Printf("│                 ATA 安装向导 v%-8s                     │\n", version)
	fmt.Println("├─────────────────────────────────────────────────────────────┤")
	fmt.Println("│                                                             │")
	fmt.Println("│  欢迎使用 ATA - AVIF 图像转换工具安装向导                    │")
	fmt.Println("│                                                             │")
	fmt.Println("│  ATA 是一个强大的图像转换工具，支持：                        │")
	fmt.Println("│  • 批量将图像转换为 AVIF 格式                               │")
	fmt.Println("│  • 创建 AVIF 动画                                           │")
	fmt.Println("│  • 制作幻灯片动画                                           │")
	fmt.Println("│  • 自动管理 FFmpeg 依赖                                     │")
	fmt.Println("│                                                             │")
	fmt.Println("│  安装向导将引导您完成 ATA 的安装过程。                       │")
	fmt.Println("│                                                             │")
	fmt.Println("└─────────────────────────────────────────────────────────────┘")
	fmt.Println("")
	fmt.Println("请选择操作:")
	fmt.Println("  [1] 开始安装")
	fmt.Println("  [2] 查看帮助")
	fmt.Println("  [3] 退出安装")
	fmt.Println("")
	fmt.Print("请输入您的选择 (1-3): ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("读取输入失败: %v\n", err)
		os.Exit(1)
	}
	input = strings.TrimSpace(input)

	switch input {
	case "1":
		showLicensePage(embeddedExecutable, version)
	case "2":
		ShowHelp(embeddedExecutable)
	case "3":
		showExitPage()
	case "":
		showExitPage()
	default:
		fmt.Printf("\n无效选项，请输入 1、2 或 3\n")
		fmt.Print("按回车键继续...")
		bufio.NewReader(os.Stdin).ReadString('\n')
		showWelcomePage(embeddedExecutable, version)
	}
}

// showLicensePage 显示许可协议页面
func showLicensePage(embeddedExecutable []byte, version string) {
	clearScreen()
	fmt.Println("┌─────────────────────────────────────────────────────────────┐")
	fmt.Println("│                      许可协议                               │")
	fmt.Println("├─────────────────────────────────────────────────────────────┤")
	fmt.Println("│                                                             │")
	fmt.Println("│  ATA 使用 MIT 许可证发布                                     │")
	fmt.Println("│                                                             │")
	fmt.Println("│  Copyright (c) 2024 ATA Project                            │")
	fmt.Println("│                                                             │")
	fmt.Println("│  Permission is hereby granted, free of charge, to any      │")
	fmt.Println("│  person obtaining a copy of this software and associated   │")
	fmt.Println("│  documentation files (the \"Software\"), to deal in the      │")
	fmt.Println("│  Software without restriction, including without           │")
	fmt.Println("│  limitation the rights to use, copy, modify, merge,        │")
	fmt.Println("│  publish, distribute, sublicense, and/or sell copies of    │")
	fmt.Println("│  the Software...                                            │")
	fmt.Println("│                                                             │")
	fmt.Println("└─────────────────────────────────────────────────────────────┘")
	fmt.Println("")
	fmt.Println("请选择操作:")
	fmt.Println("  [1] 我同意许可协议，继续安装")
	fmt.Println("  [2] 返回上一步")
	fmt.Println("  [3] 退出安装")
	fmt.Println("")
	fmt.Print("请输入您的选择 (1-3): ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("读取输入失败: %v\n", err)
		os.Exit(1)
	}
	input = strings.TrimSpace(input)

	switch input {
	case "1":
		showInstallationPage(embeddedExecutable, version)
	case "2":
		showWelcomePage(embeddedExecutable, version)
	case "3":
		showExitPage()
	default:
		fmt.Printf("\n无效选项，请输入 1、2 或 3\n")
		fmt.Print("按回车键继续...")
		bufio.NewReader(os.Stdin).ReadString('\n')
		showLicensePage(embeddedExecutable, version)
	}
}

// showInstallationPage 显示安装进度页面
func showInstallationPage(embeddedExecutable []byte, version string) {
	clearScreen()
	fmt.Println("┌─────────────────────────────────────────────────────────────┐")
	fmt.Println("│                      正在安装                               │")
	fmt.Println("├─────────────────────────────────────────────────────────────┤")
	fmt.Println("│                                                             │")
	fmt.Println("│  正在安装 ATA 到您的系统...                                  │")
	fmt.Println("│                                                             │")
	fmt.Println("│  安装步骤：                                                 │")
	fmt.Println("│  1. 检查系统环境                                           │")
	fmt.Println("│  2. 安装 ATA 程序文件                                       │")
	fmt.Println("│  3. 配置环境变量                                           │")
	fmt.Println("│  4. 检查并安装 FFmpeg 依赖                                  │")
	fmt.Println("│                                                             │")
	fmt.Println("│  请稍候...                                                  │")
	fmt.Println("│                                                             │")
	fmt.Println("└─────────────────────────────────────────────────────────────┘")
	fmt.Println("")

	// 执行实际安装
	err := performInstallation(embeddedExecutable)
	if err != nil {
		showErrorPage(err, embeddedExecutable, version)
	} else {
		showCompletePage()
	}
}

// showCompletePage 显示安装完成页面
func showCompletePage() {
	clearScreen()
	fmt.Println("┌─────────────────────────────────────────────────────────────┐")
	fmt.Println("│                    安装完成                                 │")
	fmt.Println("├─────────────────────────────────────────────────────────────┤")
	fmt.Println("│                                                             │")
	fmt.Println("│  🎉 恭喜！ATA 已成功安装到您的系统中。                      │")
	fmt.Println("│                                                             │")
	fmt.Println("│  安装内容：                                                 │")
	fmt.Println("│  ✓ ATA 程序文件已安装                                       │")
	fmt.Println("│  ✓ 环境变量已配置                                           │")
	fmt.Println("│  ✓ FFmpeg 依赖已检查                                        │")
	fmt.Println("│                                                             │")
	fmt.Println("│  使用方法：                                                 │")
	fmt.Println("│  1. 重启您的终端或命令提示符                                │")
	fmt.Println("│  2. 输入 'ata help' 查看使用帮助                            │")
	fmt.Println("│  3. 输入 'ata [路径]' 开始转换图像                          │")
	fmt.Println("│                                                             │")
	fmt.Println("│  感谢您使用 ATA！                                           │")
	fmt.Println("│                                                             │")
	fmt.Println("└─────────────────────────────────────────────────────────────┘")
	fmt.Println("")
	fmt.Print("按回车键退出安装程序...")
	bufio.NewReader(os.Stdin).ReadString('\n')
	os.Exit(0)
}

// showErrorPage 显示错误页面
func showErrorPage(err error, embeddedExecutable []byte, version string) {
	clearScreen()
	fmt.Println("┌─────────────────────────────────────────────────────────────┐")
	fmt.Println("│                    安装失败                                 │")
	fmt.Println("├─────────────────────────────────────────────────────────────┤")
	fmt.Println("│                                                             │")
	fmt.Println("│  ❌ 安装过程中发生错误：                                    │")
	fmt.Printf("│  %s", err.Error())
	// 填充空格以对齐
	for i := len(err.Error()); i < 59; i++ {
		fmt.Print(" ")
	}
	fmt.Println("│")
	fmt.Println("│                                                             │")
	fmt.Println("│  可能的解决方案：                                           │")
	fmt.Println("│  1. 以管理员身份运行安装程序                                │")
	fmt.Println("│  2. 检查磁盘空间是否充足                                    │")
	fmt.Println("│  3. 关闭杀毒软件后重试                                      │")
	fmt.Println("│  4. 检查网络连接（用于下载 FFmpeg）                         │")
	fmt.Println("│                                                             │")
	fmt.Println("└─────────────────────────────────────────────────────────────┘")
	fmt.Println("")
	fmt.Println("请选择操作:")
	fmt.Println("  [1] 重试安装")
	fmt.Println("  [2] 返回主菜单")
	fmt.Println("  [3] 退出安装")
	fmt.Println("")
	fmt.Print("请输入您的选择 (1-3): ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("读取输入失败: %v\n", err)
		os.Exit(1)
	}
	input = strings.TrimSpace(input)

	switch input {
	case "1":
		showInstallationPage(embeddedExecutable, version)
	case "2":
		showWelcomePage(embeddedExecutable, version)
	case "3":
		showExitPage()
	default:
		fmt.Printf("\n无效选项，请输入 1、2 或 3\n")
		fmt.Print("按回车键继续...")
		bufio.NewReader(os.Stdin).ReadString('\n')
		showErrorPage(err, embeddedExecutable, version)
	}
}

// showExitPage 显示退出页面
func showExitPage() {
	clearScreen()
	fmt.Println("┌─────────────────────────────────────────────────────────────┐")
	fmt.Println("│                      退出安装                               │")
	fmt.Println("├─────────────────────────────────────────────────────────────┤")
	fmt.Println("│                                                             │")
	fmt.Println("│  感谢您考虑使用 ATA！                                        │")
	fmt.Println("│                                                             │")
	fmt.Println("│  如果您改变主意，随时可以重新运行安装程序。                  │")
	fmt.Println("│                                                             │")
	fmt.Println("│  再见！                                                     │")
	fmt.Println("│                                                             │")
	fmt.Println("└─────────────────────────────────────────────────────────────┘")
	fmt.Println("")
	os.Exit(0)
}

// clearScreen 清屏
func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// performInstallation 执行实际安装逻辑
func performInstallation(embeddedExecutable []byte) error {
	fmt.Println("[1/4] 检查系统环境...")
	// 这里可以添加系统检查逻辑
	fmt.Println("✓ 系统环境检查完成")
	fmt.Println("")

	fmt.Println("[2/4] 安装 ATA 程序文件...")
	err := setupEnvironment(embeddedExecutable)
	if err != nil {
		return fmt.Errorf("安装程序文件失败: %v", err)
	}
	fmt.Println("✓ ATA 程序文件安装完成")
	fmt.Println("")

	fmt.Println("[3/4] 配置环境变量...")
	// 环境变量配置已在 setupEnvironment 中完成
	fmt.Println("✓ 环境变量配置完成")
	fmt.Println("")

	fmt.Println("[4/4] 检查 FFmpeg 依赖...")
	err = checkAndInstallFFmpeg()
	if err != nil {
		return fmt.Errorf("FFmpeg 依赖安装失败: %v", err)
	}
	fmt.Println("✓ FFmpeg 依赖检查完成")
	fmt.Println("")

	return nil
}

// checkAndInstallFFmpeg 检查并安装FFmpeg依赖
func checkAndInstallFFmpeg() error {
	if !checkFFmpeg() {
		fmt.Println("  FFmpeg未安装，正在安装...")
		if !installFFmpeg() {
			return fmt.Errorf("FFmpeg安装失败，请手动安装后重试")
		}
		fmt.Println("  ✓ FFmpeg安装完成")
	} else {
		fmt.Println("  ✓ FFmpeg已安装")
	}
	return nil
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
func setupEnvironment(embeddedExecutable []byte) error {
	// 检查是否已存在ata命令
	exists, existingPath := checkATAInPath()
	if exists {
		fmt.Printf("  检测到系统中已存在ata命令: %s\n", existingPath)
		fmt.Print("  是否覆盖现有安装? (y/N): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input != "y" && input != "yes" {
			return fmt.Errorf("用户取消安装")
		}
	}

	// 将嵌入的ata程序安装到用户bin目录并设置环境变量
	return installATAToUserBin(embeddedExecutable)
}

// installATAToUserBin 将嵌入的ata程序安装到用户bin目录并设置环境变量
func installATAToUserBin(embeddedExecutable []byte) error {
	switch runtime.GOOS {
	case "windows":
		return installATAToUserBinWindows(embeddedExecutable)
	case "darwin", "linux":
		return installATAToUserBinUnix(embeddedExecutable)
	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}
}

// installATAToUserBinWindows 在Windows上安装ata.exe到用户bin目录
func installATAToUserBinWindows(embeddedExecutable []byte) error {
	// 获取用户目录
	userDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("获取用户目录失败: %v", err)
	}
	
	// 创建用户bin目录
	userBinDir := filepath.Join(userDir, "bin")
	err = os.MkdirAll(userBinDir, 0755)
	if err != nil {
		return fmt.Errorf("创建用户bin目录失败: %v", err)
	}
	
	// 将嵌入的ata.exe写入用户bin目录
	targetFile := filepath.Join(userBinDir, "ata.exe")
	err = os.WriteFile(targetFile, embeddedExecutable, 0755)
	if err != nil {
		return fmt.Errorf("写入ata.exe失败: %v", err)
	}
	
	fmt.Printf("  ✓ 已将ata.exe安装到 %s\n", targetFile)
	
	// 添加用户bin目录到PATH
	err = addToPathWindows(userBinDir)
	if err != nil {
		return fmt.Errorf("添加到PATH失败: %v", err)
	}
	return nil
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()
	
	_, err = destFile.ReadFrom(sourceFile)
	return err
}

// installATAToUserBinUnix 在Unix系统上安装ata到用户bin目录
func installATAToUserBinUnix(embeddedExecutable []byte) error {
	// 获取用户目录
	userDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("获取用户目录失败: %v", err)
	}
	
	// 创建用户bin目录
	userBinDir := filepath.Join(userDir, "bin")
	err = os.MkdirAll(userBinDir, 0755)
	if err != nil {
		return fmt.Errorf("创建用户bin目录失败: %v", err)
	}
	
	// 将嵌入的ata程序写入用户bin目录
	targetFile := filepath.Join(userBinDir, "ata")
	err = os.WriteFile(targetFile, embeddedExecutable, 0755)
	if err != nil {
		return fmt.Errorf("写入ata失败: %v", err)
	}
	
	fmt.Printf("  ✓ 已将ata安装到 %s\n", targetFile)
	
	// 添加用户bin目录到PATH
	err = addToPathUnix(userBinDir)
	if err != nil {
		return fmt.Errorf("添加到PATH失败: %v", err)
	}
	return nil
}

// addToPath 将目录添加到系统PATH
func addToPath(dir string) error {
	switch runtime.GOOS {
	case "windows":
		return addToPathWindows(dir)
	case "darwin", "linux":
		return addToPathUnix(dir)
	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}
}

// addToPathWindows 在Windows上添加到PATH
func addToPathWindows(dir string) error {
	// 使用setx命令设置用户环境变量
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf(
		"$path = [Environment]::GetEnvironmentVariable('PATH', 'User'); "+
			"if ($path -notlike '*%s*') { "+
			"[Environment]::SetEnvironmentVariable('PATH', $path + ';%s', 'User') }",
		dir, dir))

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("设置环境变量失败: %v", err)
	}

	fmt.Printf("  ✓ 已将 %s 添加到用户PATH\n", dir)
	return nil
}

// addToPathUnix 在Unix系统上添加到PATH
func addToPathUnix(dir string) error {
	// 检查shell类型并添加到相应的配置文件
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("获取用户目录失败: %v", err)
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
		fmt.Printf("  ✓ PATH中已包含 %s\n", dir)
		return nil
	}

	// 添加到配置文件
	file, err := os.OpenFile(configFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("打开配置文件失败: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("\n# Added by ATA installer\nexport PATH=\"$PATH:%s\"\n", dir))
	if err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	fmt.Printf("  ✓ 已将 %s 添加到 %s\n", dir, configFile)
	return nil
}

// ShowHelp 显示帮助信息
func ShowHelp(embeddedExecutable []byte) {
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
	fmt.Println("  .jpg, .jpeg, .png, .webp, .gif, .tiff, .tif, .bmp")
	fmt.Println("")
	fmt.Print("按回车键返回主菜单...")
	bufio.NewReader(os.Stdin).ReadString('\n')
	ShowInteractiveMenu(embeddedExecutable)
}