# Go代码调试指南

在Trae中调试Go代码有多种方式，就像在IDEA中开发Java一样方便。

## 1. 直接运行代码（推荐用于快速测试）

使用 `go run` 命令可以直接运行Go代码而无需编译：

```bash
# 运行主程序
go run cmd/ata/main.go

# 带参数运行
go run cmd/ata/main.go help
go run cmd/ata/main.go convert ./images -d
```

## 2. 添加调试打印语句

在代码中添加 `fmt.Println()` 或 `fmt.Printf()` 来输出调试信息：

```go
func main() {
    fmt.Println("[DEBUG] 程序启动...")
    fmt.Printf("[DEBUG] 参数: %v\n", os.Args)
    
    // 你的代码...
    
    fmt.Println("[DEBUG] 程序结束")
}
```

## 3. 使用调试器（断点调试）

### 方法一：使用Trae调试配置

已为你创建了 `.trae/launch.json` 配置文件，包含以下调试配置：

#### F5调试传参方法

在Trae中使用F5调试时，有多种方式传递参数：

**方法1: 选择预定义配置**
按 `F5` 后选择以下配置之一：
- **Debug ATA - Help**: `help`
- **Debug ATA - GUI Mode**: 无参数（GUI模式）
- **Debug ATA - Convert**: `convert ./test -d`
- **Debug ATA - Convert Recursive**: `convert ./images -d -s -f`
- **Debug ATA - Animation**: `ani ./frames output.avif -fps 24 -crf 20`
- **Debug ATA - PPT Animation**: `ppt ./slides presentation.avif -fps 1`
- **Debug ATA - Custom Args**: 动态输入任意参数 ⭐

**方法2: 使用自定义参数**
选择 "Debug ATA - Custom Args" 配置，Trae会弹出输入框让你输入自定义参数，例如：
```
convert ./myimages -d -s -r
ani ./sequence output.avif -fps 30 -crf 15
```

#### 设置断点调试

使用方法：
1. 在代码中设置断点（点击行号左侧的空白处）
2. 按 `F5` 或使用 `Ctrl+Shift+P` 打开命令面板，输入 "Debug: Start Debugging"
3. 选择相应的调试配置（如 "Debug ATA - Help"）
4. 程序会在断点处暂停，你可以：
   - 查看变量值（鼠标悬停或在调试面板中查看）
   - 单步执行（F10: 逐过程，F11: 逐语句）
   - 继续执行（F5）
   - 查看调用栈
   - 在调试控制台中执行Go表达式

### 方法二：使用delve调试器

```bash
# 安装delve（如果还没安装）
go install github.com/go-delve/delve/cmd/dlv@latest

# 使用delve调试
dlv debug cmd/ata/main.go -- help
```

## 4. 单元测试调试

```bash
# 运行所有测试
go test ./...

# 运行特定测试并显示详细输出
go test -v ./internal/converter

# 调试特定测试
dlv test ./internal/converter -- -test.run TestConvertImage
```

## 5. 性能分析

```bash
# CPU性能分析
go run cmd/ata/main.go convert ./images -cpuprofile=cpu.prof

# 内存分析
go run cmd/ata/main.go convert ./images -memprofile=mem.prof

# 查看性能分析结果
go tool pprof cpu.prof
go tool pprof mem.prof
```

## 6. 实时代码重载

安装并使用 `air` 工具实现代码修改后自动重新运行：

```bash
# 安装air
go install github.com/cosmtrek/air@latest

# 在项目根目录运行
air
```

## 调试技巧

1. **使用条件断点**: 只在特定条件下暂停
2. **查看调用栈**: 了解函数调用路径
3. **监视变量**: 实时查看变量值变化
4. **使用日志级别**: 区分不同重要程度的调试信息
5. **单步调试**: 逐行执行代码，观察程序流程

## 常用调试命令

```bash
# 检查语法错误
go vet ./...

# 格式化代码
go fmt ./...

# 检查依赖
go mod tidy

# 构建但不运行（检查编译错误）
go build cmd/ata/main.go
```

这样你就可以像在IDEA中开发Java一样，在Trae中高效地调试Go代码了！