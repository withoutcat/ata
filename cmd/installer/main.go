package main

import (
	"embed"

	"github.com/withoutcat/ata/pkg/installer"
)

// 版本信息，在构建时通过 -ldflags 注入
var version = "dev"

//go:embed ata*
var ataFS embed.FS

func main() {
	// 安装程序入口 - 传入嵌入的ata可执行文件
	// 从embed.FS中读取ata可执行文件
	files, err := ataFS.ReadDir(".")
	if err != nil || len(files) == 0 {
		panic("No embedded executable found")
	}
	
	// 读取第一个文件（应该是ata或ata.exe）
	ataExecutable, err := ataFS.ReadFile(files[0].Name())
	if err != nil {
		panic("Failed to read embedded executable: " + err.Error())
	}
	
	installer.ShowInteractiveMenuWithVersion(ataExecutable, version)
}