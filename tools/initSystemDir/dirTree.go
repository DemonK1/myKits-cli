package initSystemDir

import (
	"fmt"
	"os"
	"path/filepath"
)

// 定义目录树结构
var dirTree = map[string]interface{}{
	"00-Dev": map[string]interface{}{
		"DB": map[string]interface{}{
			"MySQL":      nil,
			"Redis":      nil,
			"MongoDB":    nil,
			"PostgreSQL": nil,
		},
		"Run": map[string]interface{}{
			"GO": map[string]interface{}{
				"go":  nil,
				"sdk": nil,
			},
			"Python": nil,
		},
		"Tool": map[string]interface{}{
			"Git":      nil,
			"Postman ": nil,
			"VSCode":   nil,
			"GoLand":   nil,
			"PyCharm":  nil,
		},
	},
	"01-Downloads": map[string]interface{}{
		"Archive": nil,
		"Chrome":  nil,
		"Drive":   nil,
		"Edge":    nil,
	},
	"02-Projects": map[string]interface{}{
		"01-Code": map[string]interface{}{
			"GO":     nil,
			"Python": nil,
		},
		"02-采购客户项目": nil,
		"03-某工程项目":  nil,
	},
	"03-Life": map[string]interface{}{
		"01-Documents": nil,
		"02-Photos": map[string]interface{}{
			"01-Normal":      nil,
			"02-Cryptomator": nil,
			"壁纸":             nil,
			"身份证":            nil,
			"头像":             nil,
		},
		"03-Videos":      nil,
		"04-YiXue":       nil,
		"05-Investment":  nil,
		"06-iCloudDrive": nil,
	},
	"04-Work": map[string]interface{}{
		"01-Company": map[string]interface{}{
			"01-腾讯": nil,
			"02-阿里": nil,
		},
		"02-YiXue-Clients": nil,
	},
	"05-Apps": map[string]interface{}{
		"01-Installers": map[string]interface{}{
			// 安装普通软件（微信，qq，网易云等）时统一的模板，01-Data是安装位置，download是下载位置（图片和文件），cache是缓存位置
			"_template": map[string]interface{}{
				"01-Data":     nil,
				"02-Download": nil,
				"03-Cache":    nil,
			},
			"CloudMusic": map[string]interface{}{
				"01-Data":     nil,
				"02-Download": nil,
				"01-Cache":    nil,
			},
			"WeChat": map[string]interface{}{
				"01-Data":     nil,
				"02-Download": nil,
				"01-Cache":    nil,
			},
		},
		// 绿色版，免安装软件
		"02-Portable": nil,
	},
	"06-Temp":    nil,
	"07-Archive": nil,
	"08-Backup":  nil,
	"09-Others":  nil,
}

// 创建文件夹，递归函数
func createDirTree(basePath string, tree map[string]interface{}) ([]string, error) {
	var created []string
	for name, children := range tree {
		// 跳过注释键，不创建目录
		if name == "_comment" {
			continue
		}

		path := filepath.Join(basePath, name)
		if err := os.MkdirAll(path, 0755); err != nil {
			return created, fmt.Errorf("创建 %s 失败：%w", path, err)
		}
		created = append(created, path)

		// 	如果有子目录，递归创建
		if children != nil {
			subTree, ok := children.(map[string]interface{})
			if ok {
				subCreated, err := createDirTree(path, subTree)
				if err != nil {
					return created, err
				}
				created = append(created, subCreated...)
			}
		}
	}
	return created, nil
}
