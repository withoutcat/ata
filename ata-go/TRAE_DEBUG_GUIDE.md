# Trae中F5调试传参指南

在Trae中使用F5进行调试时，有多种方式传递参数：

## 方法1: 使用预定义的调试配置

按 `F5` 或 `Ctrl+Shift+P` 输入 "Debug: Start Debugging"，然后选择以下配置之一：

### 📋 可用的调试配置

1. **Debug ATA - Help**
   - 参数: `help`
   - 用途: 调试帮助信息显示

2. **Debug ATA - GUI Mode**
   - 参数: 无参数
   - 用途: 调试GUI模式启动

3. **Debug ATA - Convert**
   - 参数: `convert ./test -d`
   - 用途: 调试基本转换功能

4. **Debug ATA - Convert Recursive**
   - 参数: `convert ./images -d -s -f`
   - 用途: 调试递归转换，启用调试模式、递归处理、强制覆盖

5. **Debug ATA - Animation**
   - 参数: `ani ./frames output.avif -fps 24 -crf 20`
   - 用途: 调试动画创建功能，24fps，高质量

6. **Debug ATA - PPT Animation**
   - 参数: `ppt ./slides presentation.avif -fps 1`
   - 用途: 调试幻灯片动画创建

7. **Debug ATA - Custom Args** ⭐
   - 参数: 动态输入
   - 用途: 可以输入任意自定义参数

## 方法2: 使用自定义参数配置

选择 "Debug ATA - Custom Args" 配置时，Trae会弹出输入框让你输入自定义参数。

### 示例输入
```
# 基本转换
convert ./myimages -d

# 递归转换所有子目录
convert ./photos -d -s -r

# 创建高质量动画
ani ./sequence output.avif -fps 30 -crf 15 -threads 4

# 创建缩放动画
ani ./frames scaled.avif -scale 0.5 -fps 12

# 设置背景色的动画
ani ./images bg.avif -bg black -alpha
```

## 方法3: 修改launch.json添加新配置

如果你经常使用特定的参数组合，可以在 `.trae/launch.json` 中添加新的配置：

```json
{
    "name": "Debug ATA - My Custom Config",
    "type": "go",
    "request": "launch",
    "mode": "debug",
    "program": "${workspaceFolder}/cmd/ata",
    "args": ["convert", "./my-folder", "-d", "-s", "-crf", "25"],
    "env": {},
    "showLog": true
}
```

## 调试步骤

1. **设置断点**: 在代码中点击行号左侧设置断点
2. **启动调试**: 按 `F5` 选择调试配置
3. **输入参数**: 如果选择Custom Args，输入你的参数
4. **开始调试**: 程序会在断点处暂停
5. **调试操作**:
   - `F10`: 单步执行（不进入函数）
   - `F11`: 单步执行（进入函数）
   - `Shift+F11`: 跳出当前函数
   - `F5`: 继续执行
   - 鼠标悬停查看变量值
   - 在调试控制台中执行表达式

## 环境变量设置

如果需要设置环境变量，可以在配置中添加 `env` 字段：

```json
{
    "name": "Debug with Env Vars",
    "type": "go",
    "request": "launch",
    "mode": "debug",
    "program": "${workspaceFolder}/cmd/ata",
    "args": ["convert", "./test"],
    "env": {
        "DEBUG_MODE": "true",
        "LOG_LEVEL": "debug"
    },
    "showLog": true
}
```

## 工作目录设置

如果需要设置特定的工作目录，可以添加 `cwd` 字段：

```json
{
    "name": "Debug with Custom CWD",
    "type": "go",
    "request": "launch",
    "mode": "debug",
    "program": "${workspaceFolder}/cmd/ata",
    "args": ["convert", "./images"],
    "cwd": "${workspaceFolder}/testdata",
    "showLog": true
}
```

## 常用变量

在launch.json中可以使用以下变量：

- `${workspaceFolder}`: 工作区根目录
- `${file}`: 当前打开的文件
- `${fileBasename}`: 当前文件的文件名
- `${fileDirname}`: 当前文件的目录
- `${input:inputId}`: 用户输入变量

## 调试技巧

1. **条件断点**: 右键断点设置条件，只在特定条件下暂停
2. **日志断点**: 设置日志断点，不暂停程序但输出信息
3. **异常断点**: 在异常发生时自动暂停
4. **监视表达式**: 在调试面板中添加监视表达式
5. **调用栈**: 查看函数调用路径
6. **变量面板**: 实时查看所有变量值

这样你就可以在Trae中灵活地使用F5调试并传递各种参数了！