package excelHeaderDirs

import (
	"bufio"
	"fmt"
	"myKits-cli/tools/excelHeaderDirs/excel"
	"myKits-cli/tools/excelHeaderDirs/folder"
	"os"
	"strings"
)

func Run() {
	fmt.Println("📊 Excel 表头批量建文件夹工具")
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println("\n✨ 核心功能:")
	fmt.Println("1. 📂 在此目录下自动创建「_NewFile」根文件夹，所有生成的文件夹都会存放在这里，不会打乱原有文件")
	fmt.Println("2. 📑 支持读取 .xlsx / .xls 格式的 Excel 文件，自动匹配同名文件")
	fmt.Println("3. 📁 选中目标表头列后，会将该列所有【非空内容】作为名称，批量创建对应子文件夹")
	fmt.Println("💡 操作提示：按屏幕提示输入内容即可，无需手动添加 .xlsx/.xls 后缀")
	fmt.Println() // 空一行，和后续的输入提示隔开
	fmt.Println(strings.Repeat("─", 60))
	// 获取当前目录
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("\n❌ 获取当前目录失败: %v\n", err)
		return
	}
	fmt.Printf("\n当前工作目录: %s\n", dir)

	fmt.Printf("📊 读取 Excel 表头创建文件夹\n\n")
	fmt.Println(strings.Repeat("─", 60))
	reader := bufio.NewReader(os.Stdin)
	// 查找excel文件
	input, err := excel.FindExcelFile(reader)
	if err != nil {
		fmt.Printf("❌ 查找文件失败: %v\n", err)
		return // 返回 Run()，最终回到主菜单
	}
	// 	读取数据
	rows, colIdx, err := excel.ReadProducts(input, reader)
	if err != nil {
		fmt.Printf("❌ 读取 Excel 失败: %v\n", err)
		return
	}
	// 创建文件夹结构
	created, existed, err := folder.CreateProductFolders(rows, colIdx, dir)
	if err != nil {
		return
	}

	fmt.Printf("\n✅ 完成！共创建 %d 个文件夹，已有 %d 个文件夹，总计 %d 个。\n",
		created, existed, created+existed)

	// 输出汇总
	fmt.Println()
	fmt.Println("=" + "=========================================================")
	fmt.Println("任务完成!")
	fmt.Println("=" + "=========================================================")
	fmt.Printf("📊 Excel文件: %s\n", input)
	fmt.Println("📁 主文件夹: 产品图片/")
	fmt.Printf("📂 新建子文件夹: %d 个\n", created)
	fmt.Printf("📂 已有子文件夹: %d 个\n", existed)
	fmt.Printf("📂 总计子文件夹: %d 个\n", created+existed)
	fmt.Println("=" + "=========================================================")
	fmt.Println("✅ 所有文件夹创建完成!")
	// waitAndExit(0)
}

// func waitAndExit(code int) {
// 	fmt.Print("\n程序执行完毕，按回车键退出...")
// 	fmt.Scanln()
// 	os.Exit(code)
// }
