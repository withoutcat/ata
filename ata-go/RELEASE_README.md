# ATA - AVIF图像转换工具 发布版本

## 简介

ATA是一个强大的AVIF图像转换工具，支持批量转换图像为AVIF格式，以及创建AVIF动画。

## 快速开始

### 1. 下载

根据您的操作系统下载对应的可执行文件：
- **Windows**: `ata-windows.exe`
- **Linux**: `ata-linux`
- **macOS**: `ata-macos`

### 2. 安装

#### Windows用户
1. 下载 `ata-windows.exe`
2. 双击运行，或在命令行中运行
3. 选择 "1. Install" 进行安装
4. 安装程序会自动：
   - 检查并安装FFmpeg依赖（需要Chocolatey）
   - 将ATA添加到系统PATH环境变量

#### Linux/macOS用户
1. 下载对应的可执行文件
2. 添加执行权限：`chmod +x ata-linux` 或 `chmod +x ata-macos`
3. 运行安装程序：`./ata-linux` 或 `./ata-macos`
4. 选择 "1. Install" 进行安装
5. 安装程序会自动：
   - 检查并安装FFmpeg依赖
   - 将ATA添加到PATH环境变量

### 3. 使用

安装完成后，重新打开终端，您就可以在任意位置使用 `ata` 命令了：

```bash
# 转换单个目录下的图像
ata ./images

# 递归转换所有子目录
ata -s ./photos

# 创建动画
ata ani ./frames output.avif

# 查看帮助
ata help
```

## 功能特性

- ✅ **批量转换**: 支持批量转换多种图像格式到AVIF
- ✅ **动画创建**: 从图像序列创建AVIF动画
- ✅ **进度显示**: 实时显示转换进度
- ✅ **跨平台**: 支持Windows、Linux、macOS
- ✅ **自动安装**: 一键安装FFmpeg依赖和环境变量
- ✅ **彩色输出**: 清晰的成功/失败状态显示

## 支持的图像格式

输入格式：`.jpg`, `.jpeg`, `.png`, `.webp`, `.gif`, `.tiff`, `.tif`, `.bmp`
输出格式：`.avif`

## 系统要求

- **FFmpeg**: 自动安装（需要包管理器）
  - Windows: Chocolatey
  - macOS: Homebrew
  - Linux: apt/yum/pacman

## 命令行选项

```
选项:
  -d        启用调试模式
  -r        删除原始文件
  -f        强制覆盖已存在的文件
  -s        递归处理子目录

动画选项:
  -fps N    设置帧率 (默认: ani=10, ppt=1)
  -crf N    设置质量 (0-63, 越低质量越好, 默认: 30)
  -speed N  设置编码速度 (0-10, 越高越快, 默认: 8)
  -threads N 设置线程数 (0=自动, 默认: 0)
  -alpha    保留透明通道
  -width N  设置输出宽度
  -height N 设置输出高度
  -scale N  设置缩放比例 (默认: 1.0)
  -bg COLOR 设置背景颜色 (默认: white)
```

## 使用示例

```bash
# 基本转换
ata ./images

# 递归转换，启用调试模式
ata -s -d ./photos

# 强制覆盖现有文件
ata -f ./images

# 创建高质量动画
ata ani -fps 24 -crf 20 ./frames output.avif

# 创建幻灯片动画（低帧率）
ata ppt -fps 1 ./slides presentation.avif
```

## 故障排除

### FFmpeg安装失败
如果自动安装FFmpeg失败，请手动安装：

**Windows**:
1. 安装Chocolatey: https://chocolatey.org/install
2. 运行: `choco install ffmpeg`

**macOS**:
1. 安装Homebrew: https://brew.sh
2. 运行: `brew install ffmpeg`

**Linux**:
```bash
# Ubuntu/Debian
sudo apt update && sudo apt install ffmpeg

# CentOS/RHEL
sudo yum install ffmpeg

# Arch Linux
sudo pacman -S ffmpeg
```

### 环境变量未生效
安装完成后，请重新打开终端窗口以使环境变量生效。

## 开发者构建指南

如果您是开发者并且修改了代码，需要重新构建发布版本，请按照以下步骤操作：

### 前提条件
- 安装Go语言环境 (Go 1.19+)
- 确保在项目根目录 `ata-go/` 下执行构建命令

### Windows开发者
```powershell
# 在PowerShell中运行（推荐）
powershell -ExecutionPolicy Bypass -File ./build-release.ps1

# 或者直接运行（如果执行策略允许）
./build-release.ps1
```

### Linux/macOS开发者
```bash
# 方法1: 使用PowerShell Core（推荐，如果已安装）
pwsh ./build-release.ps1

# 方法2: 手动构建各平台版本
# 清理旧文件
rm -rf release
mkdir -p release

# 构建Windows版本
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o release/ata-windows.exe ./cmd/ata

# 构建Linux版本
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o release/ata-linux ./cmd/ata

# 构建macOS版本
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o release/ata-macos ./cmd/ata

echo "✓ 所有平台构建完成！"
ls -la release/
```

### 构建说明
- 构建脚本会自动创建 `release/` 目录
- 生成三个平台的可执行文件：
  - `ata-windows.exe` (Windows)
  - `ata-linux` (Linux)
  - `ata-macos` (macOS)
- 使用 `-ldflags "-s -w"` 参数减小文件体积
- 构建完成后可直接分发 `release/` 目录中的文件

### 验证构建
```bash
# 测试Windows版本（在Windows上）
./release/ata-windows.exe help

# 测试Linux版本（在Linux上）
./release/ata-linux help

# 测试macOS版本（在macOS上）
./release/ata-macos help
```

## 开发者信息

- **项目地址**: https://github.com/withoutcat/ata
- **问题反馈**: https://github.com/withoutcat/ata/issues
- **许可证**: MIT License

---

**注意**: 首次运行时选择安装选项，后续可直接使用命令行参数进行转换操作。