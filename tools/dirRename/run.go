package dirRename

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Run() {
	fmt.Println("📁 批量新建/复制文件夹（带重命名）")
	fmt.Println(strings.Repeat("─", 40))

	reader := bufio.NewReader(os.Stdin)

	// 0. 选择操作模式
	fmt.Println("\n请选择操作模式：")
	fmt.Println("  1. 仅新建文件夹（空文件夹）")
	fmt.Println("  2. 仅重命名磁盘原有文件夹（保留内容）")
	fmt.Println("  3. 一起都重命名（新建 + 原有，统一编号）")
	fmt.Println("  4. 都不重命名（跳过）")
	fmt.Print("请输入编号 (1-4): ")
	modeStr, _ := reader.ReadString('\n')
	modeStr = strings.TrimSpace(modeStr)
	mode, err := strconv.Atoi(modeStr)
	if err != nil || mode < 1 || mode > 4 {
		fmt.Println("❌ 无效选择")
		return
	}
	if mode == 4 {
		fmt.Println("操作已取消。")
		return
	}

	// 获取当前目录
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("❌ 获取当前目录失败: %v\n", err)
		return
	}
	fmt.Printf("\n当前工作目录: %s\n\n", currentDir)

	// 创建主目录
	baseDir := filepath.Join(currentDir, "_NewFile_Dir")
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		fmt.Printf("❌ 创建主目录失败: %v\n", err)
		return
	}

	// 1. 根据模式决定是否需要新建空文件夹
	newCount := 0
	if mode == 1 || mode == 3 {
		fmt.Print("请输入要新建的文件夹数量: ")
		countStr, _ := reader.ReadString('\n')
		countStr = strings.TrimSpace(countStr)
		n, err := strconv.Atoi(countStr)
		if err != nil || n <= 0 {
			fmt.Println("❌ 请输入有效的正整数")
			return
		}
		newCount = n
	}

	// 2. 根据模式决定是否需要扫描磁盘文件夹
	var diskFolders []string
	if mode == 2 || mode == 3 {
		entries, err := os.ReadDir(currentDir)
		if err != nil {
			fmt.Printf("❌ 读取目录失败: %v\n", err)
			return
		}
		for _, e := range entries {
			if e.IsDir() && e.Name() != "_NewFile_Dir" {
				diskFolders = append(diskFolders, e.Name())
			}
		}
		if len(diskFolders) == 0 {
			fmt.Println("当前目录没有其他文件夹可复制。")
			if mode == 2 { // 仅复制模式，无文件夹则直接返回
				return
			}
			// 如果是一起模式，至少还可以新建空文件夹
		}
	}

	// 3. 输入名称前缀和起始编号（对所有模式适用）
	fmt.Print("请输入文件夹名称前缀（留空默认使用 'folder'）: ")
	prefix, _ := reader.ReadString('\n')
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		prefix = "folder"
	}
	fmt.Printf("使用前缀: %s\n", prefix)

	fmt.Print("请输入起始编号（直接回车从001开始，输入数字如5则从5开始且不补零）: ")
	startInput, _ := reader.ReadString('\n')
	startInput = strings.TrimSpace(startInput)

	var startNum int
	useZeroPadding := true
	if startInput == "" {
		startNum = 1
		fmt.Println("使用默认起始编号: 001")
	} else {
		n, err := strconv.Atoi(startInput)
		if err != nil || n < 1 {
			fmt.Println("❌ 请输入大于0的数字")
			return
		}
		startNum = n
		useZeroPadding = false
		fmt.Printf("起始编号设置为: %d (不补零)\n", startNum)
	}

	// 4. 执行操作，统一编号
	created := 0
	existed := 0
	errors := 0
	currentIndex := startNum

	// 先处理新建空文件夹（如果有）
	for i := 0; i < newCount; i++ {
		num := currentIndex
		currentIndex++
		dirName := formatDirName(prefix, num, useZeroPadding)
		fullPath := filepath.Join(baseDir, dirName)
		err := os.Mkdir(fullPath, 0755)
		if err != nil {
			if os.IsExist(err) {
				fmt.Printf("  ⚠️  已存在，跳过: %s\n", dirName)
				existed++
			} else {
				fmt.Printf("  ❌ 创建失败: %s (%v)\n", dirName, err)
				errors++
			}
		} else {
			fmt.Printf("  📁 新建: %s\n", dirName)
			created++
		}
	}

	// 再处理复制磁盘文件夹（如果有）
	for _, oldName := range diskFolders {
		num := currentIndex
		currentIndex++
		newName := formatDirName(prefix, num, useZeroPadding)
		newPath := filepath.Join(baseDir, newName)

		// 检查目标是否已存在
		if _, err := os.Stat(newPath); err == nil {
			fmt.Printf("  ⚠️  目标已存在，跳过: %s\n", newName)
			existed++
			continue
		}

		oldPath := filepath.Join(currentDir, oldName)
		if err := copyDir(oldPath, newPath); err != nil {
			fmt.Printf("  ❌ 复制失败 %s -> %s: %v\n", oldName, newName, err)
			errors++
		} else {
			fmt.Printf("  📁 复制: %s -> %s\n", oldName, newName)
			created++
		}
	}

	// 5. 输出汇总
	fmt.Printf("\n✅ 完成！成功 %d 个，跳过 %d 个，失败 %d 个。\n", created, existed, errors)
	fmt.Printf("所有操作均在: %s\n", baseDir)
	// waitAndExit(0)
}

// func waitAndExit(code int) {
// 	fmt.Print("\n程序执行完毕，按回车键退出...")
// 	fmt.Scanln()
// 	os.Exit(code)
// }
