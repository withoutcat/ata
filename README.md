# ata - 图片批量转换为 AVIF (PowerShell 版本)

一个功能强大的 PowerShell 脚本，支持批量将图片转换为 AVIF 格式，并可将多张 AVIF 合成为动图。

## ✨ 特性

- 🖼️ **批量转换**：支持 jpg、jpeg、png、webp、bmp、tiff、heic、heif 等格式
- 🎬 **动图合成**：将多张 AVIF 按文件名顺序合成为动画 AVIF
- 🔄 **递归处理**：支持子文件夹递归转换
- 🎯 **智能缩放**：动图支持 contain/cover/stretch 三种缩放模式
- 🚀 **高性能**：使用 libaom-av1 编码器，支持多线程
- 🎨 **灵活配置**：可调节质量、速度、帧率等参数
- 💾 **自包含**：内置 ffmpeg，无需额外安装

## 🚀 快速开始

### 1. 克隆项目
```bash
git clone <your-repo-url>
cd ata
```

### 2. 安装到系统 PATH
```bash
# 方式一：双击 install.bat（推荐）
# 方式二：右键 install.bat -> "以管理员身份运行"（添加到系统 PATH）
# 方式三：手动运行 PowerShell 命令
powershell -ExecutionPolicy Bypass -File install.ps1
```

**安装脚本特性：**
- 自动检测 PowerShell 执行策略
- 智能处理权限问题
- 支持用户 PATH（推荐）或系统 PATH
- 自动验证安装结果

### 3. 使用命令
```powershell
# 查看帮助
ata /help

# 转换当前目录图片
ata ./

# 转换指定目录
ata "D:\MyPictures\2025-08-25"

# 合成动图
ata ani ./
```

## 📖 使用说明

### 基本转换命令
```powershell
ata <目录路径> [参数]
```

**参数说明：**
- `-d, -del, -Delete`：转换成功后删除原图片
- `-r, -Recursive`：递归处理子文件夹
- `-f, -Force`：忽略数量检查和文件夹确认
- `-debug, -dbg`：显示详细的 ffmpeg 日志

**示例：**
```powershell
# 转换并删除原图
ata "D:\Photos" -d

# 递归转换子文件夹
ata "D:\Photos" -r

# 静默模式（不询问）
ata "D:\Photos" -d -f -r

# 显示调试信息
ata "D:\Photos" -debug
```

### 动图合成命令
```powershell
ata ani <目录路径> [参数]
```

**参数说明：**
- `-Output <文件路径>`：指定输出文件（默认：目录名.avif）
- `-Fps <数值>`：设置帧率（默认：10）
- `-Width <数值>`, `-Height <数值>`：指定目标尺寸
- `-ScaleMode <模式>`：缩放模式（contain/cover/stretch，默认：contain）
- `-Background <颜色>`：背景色（默认：black）
- `-Alpha`：保留透明通道
- `-Crf <数值>`：质量设置（0-63，默认：28）
- `-Speed <数值>`：编码速度（0-8，默认：5）
- `-Threads <数值>`：线程数
- `-debug, -dbg`：显示详细日志

**示例：**
```powershell
# 基本合成
ata ani "D:\Frames"

# 自定义参数
ata ani "D:\Frames" -Fps 15 -Crf 20 -ScaleMode cover

# 指定尺寸和背景
ata ani "D:\Frames" -Width 800 -Height 600 -Background white

# 保留透明通道
ata ani "D:\Frames" -Alpha -ScaleMode contain
```

## 🏗️ 项目结构

```
ata/
├── ata.ps1                 # 主入口脚本
├── install.bat             # 安装脚本（推荐）
├── install.ps1             # 安装脚本（备用）
├── README.md               # 说明文档
├── module/
│   └── ata.psm1           # PowerShell 模块
├── src/
│   ├── Public/            # 公共函数
│   │   ├── Show-Help.ps1
│   │   ├── Convert-Images.ps1
│   │   └── Create-AvifAnimation.ps1
│   └── Private/           # 私有函数
│       ├── SupportedExtensions.ps1
│       ├── Find-FFmpegPath.ps1
│       ├── Invoke-FFmpeg.ps1
│       ├── Prepare-AnimationFrames.ps1
│       └── Encode-AvifAnimation.ps1
└── module/
    └── ffmpeg-n7.1-latest-win64-gpl-7.1/
        └── bin/
            └── ffmpeg.exe  # 内置 ffmpeg
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
- PowerShell 5.1 或更高版本
- 至少 2GB 可用内存（处理大图片时）

## 🔧 故障排除

### 常见问题

**Q: 提示"未找到 ffmpeg 可执行文件"**
A: 确保 `module/ffmpeg-n7.1-latest-win64-gpl-7.1/bin/ffmpeg.exe` 存在，或系统 PATH 中有 ffmpeg

**Q: 转换后的图片质量不佳**
A: 使用 `-Crf` 参数调整质量（数值越小质量越高，如 `-Crf 20`）

**Q: 动图编码速度慢**
A: 使用 `-Speed` 参数调整编码速度（数值越大速度越快，如 `-Speed 8`）

**Q: 中文路径显示乱码**
A: 确保 PowerShell 使用 UTF-8 编码，或使用英文路径

**Q: 安装脚本无法运行**
A: 
1. 确保在项目根目录运行
2. 以管理员身份运行（如需添加到系统 PATH）
3. 手动运行：`powershell -ExecutionPolicy Bypass -File install.ps1`

### 性能优化建议
- 对于大量图片，使用 `-f` 参数跳过确认
- 动图编码时，适当提高 `-Speed` 值（如 6-8）
- 根据 CPU 核心数设置合适的 `-Threads` 值

## 📝 更新日志

### v1.0.0
- 基础图片批量转换功能
- 支持多种图片格式
- 递归处理子文件夹
- 动图 AVIF 合成功能
- 自包含 ffmpeg 依赖
- 一键安装脚本（.bat 和 .ps1）

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

本项目采用 MIT 许可证。

## 🙏 致谢

- [FFmpeg](https://ffmpeg.org/) - 强大的多媒体处理工具
- [libaom-av1](https://aomedia.googlesource.com/aom/) - AV1 视频编码器 