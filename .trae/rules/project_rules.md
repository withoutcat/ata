# ATA 项目规则

## 响应格式
- 在所有 response 内容前加上 🐱 emoji，表示规则生效

## 代码风格与设计哲学
- 保持与现有代码一致的变量命名风格（Go 语言驼峰命名法）
- 遵循项目的设计哲学：简洁、高效、用户友好的 AVIF 图像转换工具
- 保持代码结构清晰，模块化设计
- 注释使用中文，保持简洁明了
- 错误处理要完善，提供有意义的错误信息

## 项目架构原则
- 命令行工具（cmd/ata）和安装程序（cmd/installer）分离
- 核心功能模块化（internal 包）
- 公共接口统一（pkg 包）
- 依赖管理清晰，避免循环依赖

## 任务完成要求
每次任务完成后，必须提供：

### 1. Git Commit Message
格式：`类型(范围): 简短描述`

示例：
- `feat(converter): 添加批量转换功能`
- `fix(installer): 修复 Windows 平台安装路径问题`
- `docs(readme): 更新安装说明`
- `refactor(build): 优化构建脚本性能`

### 2. Git 提交和推送命令
提供完整的 git 命令序列：
```bash
git add .
git commit -m "[生成的commit message]"
git push origin ${当前的分支}
```

## 开发规范
- 使用 Go 1.21+ 语法特性
- 单元测试覆盖核心功能
- 构建脚本支持跨平台（Windows/Linux/macOS）
- 安装程序要有友好的用户界面
- 日志记录要详细但不冗余