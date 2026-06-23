package excel

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

// excel 读取与结构定义

// Product 代表一行产品数据
type Product struct {
	Name        string
	Description string
	Category    string
	Detail      string
}

// ReadProducts 从excel文件读取产品列表，返回产品切片及可能的错误
func ReadProducts(input string, reader *bufio.Reader) ([][]string, string, error) {
	// 2. 打开 Excel 文件
	f, err := excelize.OpenFile(input)
	if err != nil {
		return nil, "", fmt.Errorf("❌ 打开文件失败: %v\n", err)
	}
	defer func(f *excelize.File) {
		err := f.Close()
		if err != nil {
			fmt.Printf("打开文件失败后关闭文件出现错误：%v\n", err)
		}
	}(f)

	// 获取第一个工作表
	sheet := f.GetSheetName(0)
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, "", fmt.Errorf("❌ 读取工作表失败: %v\n", err)
	}
	if len(rows) == 0 {
		return nil, "", fmt.Errorf("❌ 工作表为空")
	}

	// 3. 获取表头
	header := rows[0]
	if len(header) == 0 {
		return nil, "", fmt.Errorf("❌ 表头行为空")
	}

	fmt.Printf("\n📋 找到以下表头字段：\n\n")
	for i, h := range header {
		fmt.Printf("  %d. %s\n", i+1, h)
	}

	// 4. 用户选择要使用的字段
	fmt.Print("\n请输入要作为文件夹名称的字段编号（或字段名）：")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var colIdx = -1
	// 先尝试按编号解析
	if n, err := fmt.Sscanf(choice, "%d", &colIdx); err == nil && n == 1 {
		colIdx-- // 转为0基索引
		if colIdx < 0 || colIdx >= len(header) {
			return nil, "", fmt.Errorf("❌ 编号超出范围")
		}
	} else {
		// 按字段名匹配（不区分大小写，去除首尾空格）
		trimmed := strings.TrimSpace(choice)
		for i, h := range header {
			if strings.EqualFold(strings.TrimSpace(h), trimmed) {
				colIdx = i
				break
			}
		}
		if colIdx == -1 {
			return nil, "", fmt.Errorf("❌ 未找到字段: %s\n", choice)
		}
	}

	fieldName := header[colIdx]
	fmt.Printf("\n✅ 已选择字段: %s\n\n", fieldName)
	return rows, strconv.Itoa(colIdx), nil
}

// FindExcelFile 在当前目录查找 "*.xlsx" 文件
func FindExcelFile(reader *bufio.Reader) (string, error) {
	fmt.Print("请输入 Excel 文件名（不需要后缀，直接回车查看所有 .xlsx/.xls 文件）：")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	// 情况1：用户留空 → 列出当前目录所有 Excel 文件
	if input == "" {
		return findFileFromList(reader)
	}

	// 情况2：用户输入了文件名，查找匹配项
	return findFileByName(input, reader)
}

// findFileByName 根据用户输入的文件名（可能无后缀）查找匹配的 Excel 文件
// 如果只有一个匹配自动选用，多个则让用户选择
func findFileByName(name string, reader *bufio.Reader) (string, error) {
	// 如果用户已输入完整带后缀的文件名
	if strings.HasSuffix(name, ".xlsx") || strings.HasSuffix(name, ".xls") {
		if _, err := os.Stat(name); err == nil {
			return name, nil
		}
		return "", fmt.Errorf("文件不存在: %s", name)
	}

	// 收集所有可能的匹配文件
	patterns := []string{name + ".xlsx", name + ".xls"}
	var matches []string
	for _, pattern := range patterns {
		if _, err := os.Stat(pattern); err == nil {
			matches = append(matches, pattern)
		}
	}

	switch len(matches) {
	case 0:
		return "", fmt.Errorf("未找到文件 %s.xlsx 或 %s.xls", name, name)
	case 1:
		fmt.Printf("已自动选用文件: %s\n", matches[0])
		return matches[0], nil
	default:
		// 多个匹配，让用户选择
		fmt.Println("找到多个匹配的文件，请选择：")
		for i, f := range matches {
			fmt.Printf("\n  %d. %s\n\n", i+1, f)
		}
		fmt.Print("请输入编号：")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		var idx int
		_, err := fmt.Sscanf(choice, "%d", &idx)
		if err != nil || idx < 1 || idx > len(matches) {
			return "", fmt.Errorf("无效选择")
		}
		return matches[idx-1], nil
	}
}

// findFileFromList 列出所有 .xlsx 和 .xls 文件，让用户选择
func findFileFromList(reader *bufio.Reader) (string, error) {
	entries, err := os.ReadDir(".")
	if err != nil {
		return "", fmt.Errorf("读取当前目录失败: %w", err)
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".xlsx") || strings.HasSuffix(name, ".xls") {
			files = append(files, name)
		}
	}

	if len(files) == 0 {
		return "", fmt.Errorf("当前目录没有找到任何 .xlsx 或 .xls 文件")
	}

	if len(files) == 1 {
		fmt.Printf("自动选择唯一 Excel 文件: %s\n", files[0])
		return files[0], nil
	}

	fmt.Println("找到多个 Excel 文件，请选择：")
	for i, f := range files {
		fmt.Printf("  %d. %s\n", i+1, f)
	}
	fmt.Print("请输入编号：")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var idx int
	_, err = fmt.Sscanf(choice, "%d", &idx)
	if err != nil || idx < 1 || idx > len(files) {
		return "", fmt.Errorf("无效选择")
	}
	return files[idx-1], nil
}
