package dirRename

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// formatDirName 根据前缀、编号和补零标志生成文件夹名
func formatDirName(prefix string, num int, zeroPad bool) string {
	if zeroPad {
		return fmt.Sprintf("%s-%03d", prefix, num)
	}
	return fmt.Sprintf("%s-%d", prefix, num)
}

// copyDir 递归复制文件夹（纯标准库实现）
func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}
		return copyFile(path, targetPath, info.Mode())
	})
}

func copyFile(src, dst string, perm os.FileMode) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}
	return dstFile.Sync()
}
