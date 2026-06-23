package main

import (
	"bufio"
	"fmt"
	"myKits-cli/tools/dirRename"
	"myKits-cli/tools/excelHeaderDirs"
	"myKits-cli/tools/photoBatch"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/term"
)

// 导出功能
type exportItem struct {
	name    string
	cmdDir  string
	binName string
}

type menuItem struct {
	label    string
	callback func()
}

var ProjectRoot string

func main() {
	ProjectRoot, _ = os.Getwd()
	for {
		showMainMenu()
		// 每次从工具返回后，主菜单重新显示，因此这里循环即可。
	}
}

func showMainMenu() {
	items := []menuItem{
		{"📷 1.批量重命名照片（压缩/后缀）", runPhoto},
		{"📁 2.批量重命名文件夹（自动创建/名称）", runDir},
		{"📊 3.读取 Excel 表头创建文件夹", runExcel},
		{"🏗️ 4.创建自定义系统目录结构", runStructure},
		{"💬 5.微信多开", runWechat},
		{"🗄️ 6.启动数据库服务", runDB},
		{"📤 7.导出工具为独立程序", runExport},
		{"❌ 8.退出", func() { os.Exit(0) }},
	}

	selected := 0
	renderMenu(items, selected)

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	buf := make([]byte, 3)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			break
		}
		if n == 0 {
			continue
		}

		switch {
		case buf[0] == 27 && n >= 2:
			if buf[1] == '[' && n == 3 {
				switch buf[2] {
				case 'A': // 上
					if selected > 0 {
						selected--
					}
				case 'B': // 下
					if selected < len(items)-1 {
						selected++
					}
				}
			}
		case buf[0] == 13 || buf[0] == 10: // Enter
			term.Restore(int(os.Stdin.Fd()), oldState)
			clearScreen()
			items[selected].callback()
			return // 执行完工具后返回，外层 main 循环会重新显示菜单
		case buf[0] == 'q', buf[0] == 'Q':
			term.Restore(int(os.Stdin.Fd()), oldState)
			os.Exit(0)
		case buf[0] >= '1' && buf[0] <= '8':
			idx := int(buf[0] - '1')
			if idx < len(items) {
				term.Restore(int(os.Stdin.Fd()), oldState)
				clearScreen()
				items[idx].callback()
				return
			}
		}

		renderMenu(items, selected)
	}
}

func renderMenu(items []menuItem, selected int) {
	clearScreen()
	fmt.Println("\n  📦 kits 交互式工具集")
	fmt.Println("  " + strings.Repeat("─", 30))
	for i, item := range items {
		if i == selected {
			fmt.Printf("\033[1;32m👉 %s\033[0m\n", item.label)
		} else {
			fmt.Printf("   %s\n", item.label)
		}
	}
	fmt.Println("\n  ↑↓选择  Enter确认  数字直达  q退出")
}

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

// ========== 工具调用函数 ==========
// 每个工具 Run() 结束后，等待用户按任意键，然后主菜单会再次出现

func runPhoto() {
	photoBatch.Run()
	waitToExitOrMenu()
}

func runDir() {
	dirRename.Run()
	waitToExitOrMenu()
}

func runExcel() {
	excelHeaderDirs.Run()
	waitToExitOrMenu()
}

func runStructure() {
	waitToExitOrMenu()
}

func runWechat() {
	waitToExitOrMenu()
}

func runDB() {
	waitToExitOrMenu()
}

// 等待任意按键（在 cook 模式下）
func waitToExitOrMenu() {
	fmt.Print("\n按 m 返回主菜单，按其他任意键退出程序...")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "m" {
		return // 回到主菜单
	}
	// 其他任何输入（包括空回车）都直接退出程序
	os.Exit(0)
}

func runExport() {
	exports := []exportItem{
		{"批量重命名照片", "./app/photoApp", "PhotoBatch"},
		{"批量重命名文件夹", "app/dirApp", "DirBatch"},
		{"读取 Excel 创建文件夹", "./app/excelHeaderDirs", "ExcelHeaderDirs"},
		{"微信多开", "./app/wechatMulti", "WechatMulti"},
		{"启动数据库服务", "./app/dbStart", "DBStart"},
	}

	selected := 0
	renderExportMenu(exports, selected)

	oldState, _ := term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	buf := make([]byte, 3)
	for {
		n, _ := os.Stdin.Read(buf)
		if n == 0 {
			continue
		}
		switch {
		case buf[0] == 27 && n >= 2:
			if buf[1] == '[' && n == 3 {
				switch buf[2] {
				case 'A':
					if selected > 0 {
						selected--
					}
				case 'B':
					if selected < len(exports)-1 {
						selected++
					}
				}
			}
		case buf[0] == 13 || buf[0] == 10: // Enter
			term.Restore(int(os.Stdin.Fd()), oldState)
			buildExport(exports[selected])
			// fmt.Print("\n按任意键返回主菜单...")
			waitToExitOrMenu()
			return
		case buf[0] == 'q', buf[0] == 'Q', buf[0] == 27: // ESC
			term.Restore(int(os.Stdin.Fd()), oldState)
			return
		}
		renderExportMenu(exports, selected)
	}
}

func renderExportMenu(items []exportItem, selected int) {
	clearScreen()
	fmt.Println("\n  📤 导出工具为独立程序")
	fmt.Println("  " + strings.Repeat("─", 30))
	for i, item := range items {
		if i == selected {
			fmt.Printf("\033[1;32m👉 %s\033[0m\n", item.name)
		} else {
			fmt.Printf("   %s\n", item.name)
		}
	}
	fmt.Println("\n  ↑↓选择  Enter编译  q/esc返回")
}

func buildExport(item exportItem) {
	outName := item.binName
	if os.PathSeparator == '\\' {
		outName += ".exe"
	}
	fmt.Printf("正在编译 %s -> %s ...\n", item.name, outName)
	cmd := exec.Command("go", "build", "-C", ProjectRoot, "-o", outName, item.cmdDir)
	fmt.Println(cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("❌ 编译失败: %v\n", err)
	} else {
		fmt.Printf("✅ 导出成功: %s\n", outName)
	}
}
