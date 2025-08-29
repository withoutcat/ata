# ATA 构建指南

本项目提供两种构建方式，用于区分开发测试和正式发布：

## 🔧 统一构建脚本 (build.ps1)

使用统一的构建脚本，通过参数控制构建类型：

### 开发构建 (Development Build)

**用途**: 日常开发、测试、调试

```powershell
.\build.ps1 -dev
```

**特点**:
- ✅ 不会修改版本号
- ✅ 构建速度快（仅构建当前平台）
- ✅ 版本号带 `-dev` 后缀（如 `0.0.2-dev`）
- ✅ 生成 `ata-installer-dev.exe`
- ✅ 适合频繁测试

### 正式发布构建 (Release Build)

**用途**: 正式版本发布

```powershell
.\build.ps1 -release
```

**特点**:
- ⚠️ **会自动递增版本号**
- ⚠️ 需要用户确认（防止误操作）
- ✅ 构建所有平台（Windows/Linux/macOS）
- ✅ 生成正式版本号（如 `0.0.2` → `0.0.3`）
- ✅ 输出到 `release/` 目录
- ✅ **文件名包含版本号**：`ata-installer-{platform}-{version}.exe`
- ✅ 适合正式发布

## 📁 构建产物

### 开发构建
```
ata-installer-dev.exe    # 开发版安装程序
```

### 正式发布构建
```
release/
├── ata-installer-windows-0.0.3.exe
├── ata-installer-linux-0.0.3
└── ata-installer-macos-0.0.3
```

## 🔄 版本管理

- 版本号存储在 `version.txt` 文件中
- 开发构建：读取版本号，添加 `-dev` 后缀
- 正式构建：自动递增 patch 版本号（0.0.1 → 0.0.2）

## 💡 最佳实践

1. **日常开发**: 使用 `build.ps1 -dev`
   ```powershell
   # 修改代码后测试
   .\build.ps1 -dev
   .\ata-installer-dev.exe
   ```

2. **正式发布**: 使用 `build.ps1 -release`
   ```powershell
   # 确认代码无误后发布
   .\build.ps1 -release
   # 输入 'y' 确认发布
   ```

3. **Git 工作流**:
   ```bash
   # 开发阶段 - 不提交构建产物
   git add .
   git commit -m "feat: 添加新功能"
   
   # 发布阶段 - 提交版本文件
   .\build.ps1 -release
   git add version.txt
   git commit -m "chore: 发布版本 v0.0.3"
   git tag v0.0.3
   ```

## ⚠️ 注意事项

- 开发构建产物 (`ata-installer-dev.exe`) 不应提交到 Git
- 正式构建产物 (`release/`) 不应提交到 Git
- 只有 `version.txt` 需要在发布时提交
- 使用 `.gitignore` 忽略构建产物