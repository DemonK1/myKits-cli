package dbStart

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Run 数据库服务管理主入口
func Run() {
	fmt.Println("🗄️ 数据库服务管理")
	fmt.Println(strings.Repeat("─", 60))

	reader := bufio.NewReader(os.Stdin)

	// 确保配置根目录存在
	root, err := configDir()
	if err != nil {
		fmt.Printf("❌ 获取配置目录失败: %v\n", err)
		return
	}
	if err := os.MkdirAll(root, 0755); err != nil {
		fmt.Printf("❌ 创建配置目录失败: %v\n", err)
		return
	}

	for {
		fmt.Println("\n请选择操作:")
		fmt.Println("  1. 查看已有数据库列表")
		fmt.Println("  2. 添加新的数据库")
		fmt.Println("  3. 返回主菜单")
		fmt.Print("请输入: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			configs, err := showList()
			if err != nil {
				fmt.Printf("❌ 读取列表失败: %v\n", err)
				continue
			}
			if len(configs) == 0 {
				continue
			}
			// 选择具体数据库进入操作
			fmt.Print("请选择要操作的数据库编号 (按回车返回): ")
			numStr, _ := reader.ReadString('\n')
			numStr = strings.TrimSpace(numStr)
			if numStr == "" {
				continue
			}
			idx, err := parseInt(numStr)
			if err != nil || idx < 1 || idx > len(configs) {
				fmt.Println("❌ 无效编号")
				continue
			}
			serviceMenu(configs[idx-1], reader)

		case "2":
			if err := addDatabase(reader); err != nil {
				fmt.Printf("❌ %v\n", err)
			}
		case "3":
			return // 返回主菜单
		default:
			fmt.Println("无效选择")
		}
	}
}
