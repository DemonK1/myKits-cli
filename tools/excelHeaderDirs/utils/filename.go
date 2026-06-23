package utils

import "strings"

// 文件名清理等工具函数

func CleanFilename(name string) string {
	// 替换路径分隔符和常见非法字符
	replacer := strings.NewReplacer(
		"/", " ",
		"\\", " ",
		":", " ",
		"*", " ",
		"?", " ",
		"\"", " ",
		"<", " ",
		">", " ",
		"|", " ",
	)
	return replacer.Replace(name)
}
