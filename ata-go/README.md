# ATA - AVIF图像转换与动画工具

## 项目重构说明

本项目是对原PowerShell版本ATA工具的Go语言重构版本，现在是一个独立的软件包，功能包括：

- 批量将图像转换为AVIF格式
- 创建AVIF动画
- 制作幻灯片
- 交互式安装程序（自动安装FFmpeg依赖和配置环境变量）
- 嵌入式程序分发（单个安装文件包含完整程序）

## 项目特性

重构后的项目具有以下特性：

1. **独立软件包**：用户只需下载单个安装文件，无需克隆整个仓库
2. **嵌入式分发**：安装程序内嵌完整的ATA程序，无需额外下载
3. **交互式安装**：提供友好的安装界面，自动处理依赖和环境配置
4. **跨平台支持**：提供Windows、Linux、macOS三个平台的发布版本
5. **自动依赖管理**：自动检测和安装FFmpeg依赖

## 环境要求

### FFmpeg 安装

本工具依赖FFmpeg进行图像和视频处理，请确保FFmpeg已安装并添加到系统PATH环境变量中。

#### 方法一：使用Chocolatey安装（推荐）

1. 安装Chocolatey（如果尚未安装）：
   ```powershell
   Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
   ```

2. 使用Chocolatey安装FFmpeg：
   ```powershell
   choco install ffmpeg
   ```

#### 方法二：手动安装

1. 访问FFmpeg官网下载页面：https://ffmpeg.org/download.html
2. 选择Windows版本，推荐下载"release builds"
3. 解压下载的文件到合适位置（如`C:\ffmpeg`）
4. 将FFmpeg的bin目录添加到系统PATH环境变量：
   - 打开"系统属性" → "高级" → "环境变量"
   - 在"系统变量"中找到"Path"，点击"编辑"
   - 添加FFmpeg的bin目录路径（如`C:\ffmpeg\bin`）
   - 点击"确定"保存

#### 验证安装

打开命令提示符或PowerShell，运行以下命令验证FFmpeg是否正确安装：
```
ffmpeg -version
```

如果显示FFmpeg版本信息，说明安装成功。

## 安装与使用

### 快速开始

#### 用户使用
1. 从 [Releases](https://github.com/withoutcat/ata/releases) 下载对应平台的安装文件：
   - Windows: `ata-installer-windows.exe`
   - Linux: `ata-installer-linux`
   - macOS: `ata-installer-macos`
2. 运行安装文件，选择 "1. Install" 进行安装
3. 安装程序会自动：
   - 提取并安装ATA程序到用户bin目录
   - 设置PATH环境变量
   - 检查并安装FFmpeg依赖
4. 重启终端后即可使用 `ata` 命令

#### 开发者构建
1. 构建发布版本：
   ```powershell
   powershell -ExecutionPolicy Bypass -File ./build-release.ps1
   ```

2. 构建过程说明：
   - 脚本会为每个平台先构建主程序
   - 然后将主程序嵌入到安装程序中
   - 最终生成包含嵌入程序的单一安装文件

3. 测试构建结果：
   ```bash
   # 运行安装程序测试
   ./release/ata-installer-windows.exe
   ```

### 使用方法

#### 命令行模式

```
# 转换图像
ata convert [路径] [-d] [-del] [-r] [-f]

# 创建动画
ata ani [输入路径] [输出路径] [-fps 值] [-crf 值] [-speed 值] [-threads 值] [-alpha] [-width 值] [-height 值] [-scale 值] [-bg 颜色] [-d] [-del] [-f]

# 创建幻灯片
ata ppt [输入路径] [输出路径] [参数同ani命令]

# 显示帮助
ata help
```

#### 交互式安装模式

直接运行`ata`命令（不带参数）即可启动交互式安装界面，提供三个选项：

1. **Install** - 安装FFmpeg依赖并配置环境变量
2. **Help** - 显示详细的使用帮助信息
3. **Exit** - 退出程序

## 开发说明

项目使用Go语言开发，采用模块化设计。主要包结构：

- `cmd/ata`: 主程序入口
- `internal/ffmpeg`: FFmpeg调用相关功能
- `internal/converter`: 图像转换和动画创建核心功能
- `internal/animation`: 动画处理逻辑
- `internal/utils`: 工具函数
- `pkg/cli`: 命令行相关功能
- `pkg/installer`: 交互式安装程序
- `pkg/logger`: 日志系统

### 构建说明
- 使用 `build-release.ps1` 脚本可构建Windows、Linux、macOS三个平台的发布版本
- 构建产物位于 `release/` 目录下
- 详细构建说明请参考 `RELEASE_README.md`