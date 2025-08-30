# ata - 图片批量转换为 AVIF

一个功能强大的工具，支持批量将图片转换为 AVIF 格式，并可将多张 AVIF 合成为动图。现已使用 Go 语言重构，提供命令行和图形界面两种使用方式。

## ✨ 特性

- 🖼️ **批量转换**：支持 jpg、jpeg、png、webp、bmp、tiff、heic、heif 等格式
- 🎬 **动图合成**：将多张 AVIF 按文件名顺序合成为动画 AVIF
- 🔄 **递归处理**：支持子文件夹递归转换
- 🎯 **智能缩放**：动图支持 contain/cover/stretch 三种缩放模式
- 🚀 **高性能**：使用 libaom-av1 编码器，支持多线程
- 🎨 **灵活配置**：可调节质量、速度、帧率等参数
- 💾 **轻量级**：依赖系统 FFmpeg，保持项目精简

## 🚀 快速开始

### 1. 克隆项目
```bash
git clone <your-repo-url>
cd ata
```

### 2. 安装 FFmpeg

本工具依赖 FFmpeg 进行图像和视频处理，请确保 FFmpeg 已安装并添加到系统 PATH 环境变量中。

#### 使用 Chocolatey 安装（推荐）
```powershell
# 安装 Chocolatey（如果尚未安装）
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# 安装 FFmpeg
choco install ffmpeg
```

#### 手动安装
1. 访问 [FFmpeg 官网](https://ffmpeg.org/download.html) 下载 Windows 版本
2. 解压到合适位置（如 `C:\ffmpeg`）
3. 将 FFmpeg 的 bin 目录添加到系统 PATH 环境变量
4. 验证安装：`ffmpeg -version`

**安装脚本特性：**
- 自动检测 PowerShell 执行策略
- 智能处理权限问题
- 支持用户 PATH（推荐）或系统 PATH
- 自动验证安装结果

### 3. 下载和安装

#### 方法一：下载发布版本（推荐）
1. 从 [Releases](https://github.com/withoutcat/ata/releases) 页面下载对应平台的可执行文件：
   - Windows: `ata-windows.exe`
   - Linux: `ata-linux`
   - macOS: `ata-macos`

2. 运行安装程序：
   ```bash
   # Windows
   ./ata-windows.exe
   
   # Linux/macOS
   chmod +x ata-linux  # 或 ata-macos
   ./ata-linux         # 或 ./ata-macos
   ```

3. 选择 "1. Install" 进行安装，程序会自动：
   - 检查并安装FFmpeg依赖
   - 将ATA添加到系统PATH环境变量

#### 方法二：从源码构建

##### 🔧 构建脚本

使用统一的构建脚本，通过参数控制构建类型：

**开发构建**（日常开发、测试、调试）：
```powershell
./build.ps1 -dev
```
- ✅ 不会修改版本号
- ✅ 构建速度快（仅构建当前平台）
- ✅ 版本号带 `-dev` 后缀（如 `0.0.2-dev`）
- ✅ 生成 `ata-installer-dev.exe`

**正式发布构建**（正式版本发布）：
```powershell
./build.ps1 -release
```
- ⚠️ **会自动递增版本号**
- ⚠️ 需要用户确认（防止误操作）
- ✅ 构建所有平台（Windows/Linux/macOS）
- ✅ 生成正式版本号（如 `0.0.2` → `0.0.3`）
- ✅ 输出到 `release/` 目录
- ✅ **文件名包含版本号**：`ata-installer-{platform}-{version}.exe`

**构建产物**：
```
# 开发构建
ata-installer-dev.exe

# 正式发布构建
release/
├── ata-installer-windows-0.0.3.exe
├── ata-installer-linux-0.0.3
└── ata-installer-macos-0.0.3
```

**最佳实践**：
```powershell
# 日常开发
./build.ps1 -dev
./ata-installer-dev.exe

# 正式发布
./build.ps1 -release
# 输入 'y' 确认发布
```

### 4. 使用命令

安装完成后，重新打开终端，您就可以在任意位置使用 `ata` 命令了：

```bash
# 查看帮助
ata help

# 转换当前目录图片
ata ./

# 转换指定目录（启用调试模式）
ata "D:\MyPictures\2025-08-25" -d

# 递归转换所有子目录
ata ./photos -s

# 创建动画
ata ani ./frames output.avif

# 创建幻灯片动画
ata ppt ./slides presentation.avif
```

## 📖 使用说明

### 基本转换命令
```bash
ata <目录路径> [参数]
```

**参数说明：**
- `-r`：转换成功后删除原图片
- `-s`：递归处理子文件夹
- `-f`：强制覆盖已存在的文件
- `-d`：显示详细的调试信息

**示例：**
```bash
# 基本转换
ata "D:\Photos"

# 转换并删除原图
ata "D:\Photos" -r

# 递归转换子文件夹
ata "D:\Photos" -s

# 强制覆盖 + 递归 + 删除原图
ata "D:\Photos" -s -r -f

# 显示调试信息
ata "D:\Photos" -d
```

### 动画合成命令
```bash
ata ani <输入目录> <输出文件> [参数]
```

**参数说明：**
- `-fps <数值>`：设置帧率（默认：10）
- `-width <数值>`, `-height <数值>`：指定目标尺寸
- `-scale <数值>`：缩放比例（默认：1.0）
- `-bg <颜色>`：背景色（默认：white）
- `-alpha`：保留透明通道
- `-crf <数值>`：质量设置（0-63，默认：30）
- `-speed <数值>`：编码速度（0-10，默认：8）
- `-threads <数值>`：线程数（0=自动，默认：0）
- `-d`：显示详细日志

**示例：**
```bash
# 基本合成
ata ani "D:\Frames" output.avif

# 自定义帧率和质量
ata ani "D:\Frames" animation.avif -fps 24 -crf 20

# 指定尺寸和背景
ata ani "D:\Frames" result.avif -width 800 -height 600 -bg black

# 保留透明通道
ata ani "D:\Frames" transparent.avif -alpha

# 创建幻灯片（低帧率）
ata ppt "D:\Slides" presentation.avif -fps 1
```

## 🏗️ 项目结构

```
ata/
├── README.md               # 项目说明文档
└── ata-go/                 # Go 语言重构版本
    ├── README.md           # Go 版本说明
    ├── RELEASE_README.md   # 发布版本用户指南
    ├── build-release.ps1   # 跨平台构建脚本
    ├── go.mod              # Go 模块定义
    ├── go.sum              # Go 模块依赖锁定
    ├── cmd/
    │   └── ata/            # 主程序入口
    ├── internal/           # 内部包
    │   ├── animation/      # 动画处理
    │   ├── converter/      # 图像转换
    │   ├── ffmpeg/         # FFmpeg 集成
    │   └── utils/          # 工具函数
    ├── pkg/
    │   ├── cli/            # 命令行界面
    │   ├── gui/            # 图形界面（已移除）
    │   ├── setup/          # 安装程序启动器
    │   └── logger/         # 日志系统
    └── release/            # 发布版本文件
        ├── ata-windows.exe # Windows 可执行文件
        ├── ata-linux       # Linux 可执行文件
        └── ata-macos       # macOS 可执行文件
```

## ⚙️ 配置说明

### 支持的图片格式
脚本会自动检测以下格式：
- 静态图片：jpg、jpeg、png、webp、bmp、tiff、heic、heif
- 动图输入：仅支持 AVIF 格式

### 默认设置
- **转换质量**：CRF 28（平衡质量与大小）
- **编码速度**：5（平衡速度与压缩率）
- **动图帧率**：10 FPS
- **缩放模式**：contain（等比缩放，保持比例）
- **像素格式**：yuv420p（无透明）/ yuva420p（有透明）

### 环境要求
- Windows 10/11
- FFmpeg（已安装并添加到 PATH）
- 至少 2GB 可用内存（处理大图片时）

## 🔧 故障排除

### 常见问题

**Q: 提示"未找到 ffmpeg 可执行文件"**
A: 请按照安装说明安装 FFmpeg 并确保已添加到系统 PATH 环境变量中

**Q: 转换后的图片质量不佳**
A: 使用 `-Crf` 参数调整质量（数值越小质量越高，如 `-Crf 20`）

**Q: 动图编码速度慢**
A: 使用 `-Speed` 参数调整编码速度（数值越大速度越快，如 `-Speed 8`）

**Q: 中文路径显示乱码**
A: 确保 PowerShell 使用 UTF-8 编码，或使用英文路径

**Q: 编译失败**
A: 
1. 确保已安装 Go 1.19 或更高版本
2. 确保在 ata-go 目录下运行 build.bat
3. 检查网络连接，确保能下载 Go 模块依赖

### 性能优化建议
- 对于大量图片，使用 `-f` 参数跳过确认
- 动图编码时，适当提高 `-Speed` 值（如 6-8）
- 根据 CPU 核心数设置合适的 `-Threads` 值

## 📝 更新日志

### v10.0 (Go 重构版)
- 使用 Go 语言完全重构
- 提供命令行和图形界面两种模式
- 移除内置 FFmpeg，依赖系统安装
- 更好的性能和跨平台支持
- 现代化的命令行参数格式

### ~~v0.0.1(PowerShell 版)~~
- ~~基础图片批量转换功能~~
- ~~支持多种图片格式~~
- ~~递归处理子文件夹~~
- ~~动图 AVIF 合成功能~~

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

本项目采用 MIT 许可证。

## 🙏 致谢

- [FFmpeg](https://ffmpeg.org/) - 强大的多媒体处理工具
- [libaom-av1](https://aomedia.googlesource.com/aom/) - AV1 视频编码器