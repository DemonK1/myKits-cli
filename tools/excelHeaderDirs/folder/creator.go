package folder

import (
	"fmt"
	"myKits-cli/tools/excelHeaderDirs/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// 文件夹创建逻辑

// CreateProductFolders 根据产品列表创建文件夹结构
// 返回: 新建数、已存在数、主文件夹路径
func CreateProductFolders(rows [][]string, colIdxs string, dir string) (created, existed int, err error) {
	// 1. 先创建「_NewFile_Excel」父目录（必须第一步做！）
	newDir := filepath.Join(dir, "_NewFile_Excel")
	// 递归创建目录，避免父目录不存在导致后续子文件夹创建失败
	if err = os.MkdirAll(newDir, 0755); err != nil {
		return 0, 0, fmt.Errorf("创建_NewFile目录失败: %v", err)
	}

	created = 0
	existed = 0
	colIdx, err := strconv.Atoi(colIdxs)

	if err != nil {
		return 0, 0, err
	}
	seen := map[string]bool{} // 可选：避免重复创建同名文件夹
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		var cellValue string
		if colIdx < len(row) {
			cellValue = strings.TrimSpace(row[colIdx])
		}
		if cellValue == "" {
			continue
		}
		// 可选：如果文件夹名称可能包含非法字符，需要清理
		// 这里简单处理：替换路径分隔符为空格
		dirName := utils.CleanFilename(cellValue)

		// 避免重复创建
		if seen[dirName] {
			continue
		}
		seen[dirName] = true

		// 2. 关键改动：把子文件夹路径拼到_NewFile目录下！！！
		targetPath := filepath.Join(newDir, dirName)
		err = os.Mkdir(targetPath, 0755)

		if err != nil {
			if os.IsExist(err) {
				fmt.Printf("  ⚠️  文件夹已存在，跳过: %s\n", targetPath)
				existed++ // 已存在，计数+1
			} else {
				fmt.Printf("  ❌ 创建失败: %s (%v)\n", targetPath, err)
			}
		} else {
			fmt.Printf("  📁 创建: %s\n", targetPath)
			created++
		}
	}
	return created, existed, nil
}
