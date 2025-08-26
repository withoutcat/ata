# ATA - AVIF图像转换与动画工具

## 项目重构说明

本项目是对原PowerShell版本ATA工具的Go语言重构版本，提供了命令行和GUI两种使用方式，功能包括：

- 批量将图像转换为AVIF格式
- 创建AVIF动画
- 制作幻灯片

## 目录结构变更

重构后的项目结构做了以下调整：

1. **FFmpeg位置变更**：将FFmpeg从原来的`module/ffmpeg-n7.1-latest-win64-gpl-7.1`目录移动到了`ata-go/ffmpeg`目录下，使项目结构更加清晰。

2. **PowerShell脚本**：原有的PowerShell脚本(`*.ps1`文件)仅作为重构参考，不再参与实际代码运行。新版本完全使用Go语言实现所有功能。

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

### 构建和安装

1. 构建项目：
   ```
   build.bat
   ```

2. 安装到系统：
   ```
   install.bat
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

#### GUI模式

直接运行`ata`命令（不带参数）即可启动GUI界面，包含三个功能选项卡：

1. 批量转换 - 将图像批量转换为AVIF格式
2. 动画合成 - 创建AVIF动画
3. 幻灯片 - 制作AVIF幻灯片

## 开发说明

项目使用Go语言开发，GUI部分使用GoVCL框架。主要包结构：

- `cmd/ata`: 主程序入口
- `internal/ffmpeg`: FFmpeg调用相关功能
- `internal/converter`: 图像转换和动画创建核心功能
- `pkg/cli`: 命令行相关功能
- `pkg/gui`: GUI相关功能