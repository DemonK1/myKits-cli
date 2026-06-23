package initSystemDir

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

// 颜色码
const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	cyan   = "\033[36m"
	bold   = "\033[1m"
)

//go:embed "README.pdf"
var readmeData []byte // 将现有的 README.pdf 文档集成在 exe 程序中，而不是根据结构体创建

func Run() {
	root := "."

	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	// 创建所有目录
	createdDirs, err := createDirTree(root, dirTree)
	if err != nil {
		fmt.Println("错误：", err)
		waitAndExit()
	}

	// 将嵌入的 README.pdf 释放到根目录
	readmePath := filepath.Join(root, "README.pdf")
	if err = os.WriteFile(readmePath, readmeData, 0644); err != nil {
		fmt.Printf("警告：无法写入 README.pdf：%v\n", err)
	} else {
		createdDirs = append(createdDirs, createdDirs...)
	}

	// 打印所有已创建的文件夹和文件
	fmt.Println("已创建以下文件夹：")
	for _, dir := range createdDirs {
		fmt.Println("  ", dir)
	}
	fmt.Println("\n所有目录创建完成✅！")
	fmt.Println()
	fmt.Println(colorString(green+bold, "嵌入时读取 README.pdf 的目录："), readmePath)
	fmt.Println()

	// 可选：如果需要生成 README.md，可以在这里调用 generateReadme
	waitAndExit()
}
func waitAndExit() {
	fmt.Print("\n按回车键退出...")
	_, err := fmt.Scanln()
	if err != nil {
		return
	}
	os.Exit(0)
}

// 加点颜色
func colorString(color, text string) string {
	return color + text + reset
}
