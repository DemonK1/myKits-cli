package photoBatch

import (
	"bufio"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// 扫描目录，返回图片文件、文件夹数量和总条目数
func scanDirectory(dir string) ([]string, int, int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, 0, 0, err
	}

	// 支持的图片格式
	imageExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".tif":  true,
		".tiff": true,
		".webp": true,
		".ico":  true,
		".heic": true,
		".heif": true,
		".avif": true,
		".svg":  true,
		".psd":  true,
		".raw":  true,
		".cr2":  true,
		".nef":  true,
		".arw":  true,
		".dng":  true,
		".orf":  true,
		".raf":  true,
		".rw2":  true,
		".pcx":  true,
		".tga":  true,
		".pnm":  true,
		".pbm":  true,
		".pgm":  true,
		".ppm":  true,
		".hdr":  true,
		".exr":  true,
		".jxr":  true,
		".jp2":  true,
		".j2k":  true,
		".jpf":  true,
		".jpx":  true,
		".jpm":  true,
		".mj2":  true,
		".apng": true,
		".bpg":  true,
		".dds":  true,
		".xbm":  true,
		".xpm":  true,
		".cur":  true,
		".icns": true,
		".jfif": true,
		".jpe":  true,
		".jif":  true,
		".jfi":  true,
		".dib":  true,
		".sr":   true,
		".ras":  true,
		".emf":  true,
		".wmf":  true,
	}

	var imageFiles []string
	folderCount := 0
	totalCount := len(entries)

	for _, entry := range entries {
		if entry.IsDir() {
			folderCount++
		} else {
			// 获取文件扩展名并转换为小写
			ext := strings.ToLower(filepath.Ext(entry.Name()))

			// 检查是否是图片格式
			if imageExtensions[ext] {
				imageFiles = append(imageFiles, entry.Name())
			}
		}
	}

	return imageFiles, folderCount, totalCount, nil
}

// 显示目录统计信息
func displayDirectoryStats(imageFiles []string, folderCount, totalCount int) {
	// 统计图片格式分布
	formatStats := make(map[string]int)
	for _, file := range imageFiles {
		ext := strings.ToLower(filepath.Ext(file))
		if ext != "" {
			// 去掉点号
			format := strings.TrimPrefix(ext, ".")
			formatStats[format]++
		}
	}

	// 显示统计信息
	fmt.Println("目录信息统计:")
	fmt.Printf("  - 文件夹数量: %d 个\n", folderCount)
	fmt.Printf("  - 图片文件数量: %d 个\n", len(imageFiles))
	fmt.Printf("  - 其他文件数量: %d 个\n", totalCount-folderCount-len(imageFiles))
	fmt.Printf("  - 总计: %d 个条目\n", totalCount)

	// 显示图片格式分布
	if len(imageFiles) > 0 {
		fmt.Printf("\n图片格式分布:\n")
		// 按数量从多到少排序显示
		sortedFormats := sortFormatsByCount(formatStats)
		for _, format := range sortedFormats {
			count := formatStats[format]
			percentage := float64(count) / float64(len(imageFiles)) * 100
			fmt.Printf("  - %s: %d 个 (%.1f%%)\n",
				strings.ToUpper(format), count, percentage)
		}
	}
}

// 按数量从多到少排序格式
func sortFormatsByCount(stats map[string]int) []string {
	type formatCount struct {
		format string
		count  int
	}

	var formats []formatCount
	for format, count := range stats {
		formats = append(formats, formatCount{format, count})
	}

	// 按数量从多到少排序
	for i := 0; i < len(formats)-1; i++ {
		for j := i + 1; j < len(formats); j++ {
			if formats[j].count > formats[i].count {
				formats[i], formats[j] = formats[j], formats[i]
			}
		}
	}

	// 提取格式字符串
	result := make([]string, len(formats))
	for i, fc := range formats {
		result[i] = fc.format
	}

	return result
}

// 获取用户输入的压缩质量
func getQualityFromUser(reader *bufio.Reader) int {
	for {
		fmt.Print("\n请输入压缩质量 (10-100，推荐30-80): ")
		qualityInput, _ := reader.ReadString('\n')
		qualityInput = strings.TrimSpace(qualityInput)

		// 尝试转换为整数
		quality, err := strconv.Atoi(qualityInput)
		if err != nil {
			fmt.Println("\n请输入有效的数字!")
			continue
		}

		// 检查质量是否在有效范围内
		if quality < 10 || quality > 100 {
			fmt.Println("质量必须在10-100之间!")
			continue
		}

		return quality
	}
}

// 处理目录中的所有图片
func processImages(dir, prefix string, startNum int, convertToJPG bool, quality int, useZeroPadding bool) error {
	// 支持的图片格式
	supportedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".tif", ".webp"}

	// 创建新文件夹
	newDir := filepath.Join(dir, "_NewFile")
	err := os.MkdirAll(newDir, 0755)
	if err != nil {
		return fmt.Errorf("\n创建文件夹失败: %v", err)
	}

	// 读取目录中的所有文件
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	// 计数器
	count := startNum
	processedCount := 0
	skippedCount := 0

	fmt.Printf("\n开始处理图片...\n")

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// 检查文件扩展名
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if !contains(supportedExts, ext) {
			continue
		}

		// 原始文件路径
		oldPath := filepath.Join(dir, file.Name())

		// 声明变量
		var newName, newPath string

		if convertToJPG {
			// 转换为JPG格式
			if useZeroPadding {
				newName = fmt.Sprintf("%s-%03d.jpg", prefix, count)
			} else {
				newName = fmt.Sprintf("%s-%d.jpg", prefix, count)
			}
			newPath = filepath.Join(newDir, newName)

			// 转换图片格式
			err = convertImageToJPG(oldPath, newPath, quality)
			if err != nil {
				fmt.Printf("\n处理图片失败: %s (%v)\n", file.Name(), err)
				skippedCount++
				continue
			}
		} else {
			// 不转换格式，保持原格式
			if useZeroPadding {
				newName = fmt.Sprintf("%s-%03d%s", prefix, count, ext)
			} else {
				newName = fmt.Sprintf("%s-%d%s", prefix, count, ext)
			}
			newPath = filepath.Join(newDir, newName)

			// 直接复制文件，不进行任何处理
			err = copyFile(oldPath, newPath)
			if err != nil {
				fmt.Printf("\n复制图片失败: %s (%v)\n", file.Name(), err)
				skippedCount++
				continue
			}
		}

		// 获取文件大小信息
		oldFileInfo, _ := os.Stat(oldPath)
		newFileInfo, _ := os.Stat(newPath)
		oldSizeKB := oldFileInfo.Size() / 1024
		newSizeKB := newFileInfo.Size() / 1024

		fmt.Printf("已处理: %s -> %s (大小: %dKB -> %dKB)\n",
			file.Name(), newName, oldSizeKB, newSizeKB)
		count++
		processedCount++
	}

	// 显示处理结果
	fmt.Printf("\n处理结果:\n")
	fmt.Printf("  - 成功处理: %d 个图片\n", processedCount)
	if skippedCount > 0 {
		fmt.Printf("  - 跳过: %d 个图片\n", skippedCount)
	}
	fmt.Printf("  - 新文件保存到: %s\n", newDir)

	return nil
}

// 将任何图片格式转换为JPG
func convertImageToJPG(inputPath, outputPath string, quality int) error {
	// 打开原始图片文件
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 解码图片
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// 创建输出文件
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// 保存为JPG格式
	return jpeg.Encode(outFile, img, &jpeg.Options{Quality: quality})
}

// 复制文件（不进行任何处理）
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// 检查切片是否包含某个元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
