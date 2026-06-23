package photoBatch

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func Run() {
	fmt.Println("📸 照片批量处理工具")
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println()
	fmt.Println("✨ 核心功能：")
	fmt.Println("1. 📂 自动创建【_NewFile_Photo】根目录，所有输出文件统一存放")
	fmt.Println("2. 🖼️ 支持批量压缩照片，优化文件体积")
	fmt.Println("3. 🙌 自动将所有照片统一转换为 JPG 格式")
	fmt.Println("4. ✌️ 支持自定义输出照片的命名规则与文件后缀")
	fmt.Println()
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println()
	fmt.Println("💡 操作提示：")
	fmt.Println("本工具必须在【待处理照片的源文件夹】内运行，否则无法识别图片文件")
	fmt.Println("✅ 已在目标目录？直接按步骤；未在？先移动程序到照片目录")
	fmt.Println()
	fmt.Println(strings.Repeat("─", 60))

	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("获取当前目录失败: %v\n", err)
		return
	}

	// 显示当前目录
	fmt.Printf("当前目录: %s\n\n", currentDir)

	// 统计当前目录下的图片文件数量
	fmt.Println("正在扫描图片文件...")
	imageFiles, folderCount, totalCount, err := scanDirectory(currentDir)
	if err != nil {
		fmt.Printf("扫描目录失败: %v\n", err)
		return
	}

	// 显示目录信息统计
	displayDirectoryStats(imageFiles, folderCount, totalCount)

	// 如果没有图片文件，询问是否继续
	if len(imageFiles) == 0 {
		fmt.Print("\n当前目录下没有找到图片文件，是否继续? (y: 继续, 其他: 退出): ")
		reader := bufio.NewReader(os.Stdin)
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(strings.ToLower(choice))

		if choice != "y" {
			fmt.Println("程序退出。")
			return
		}
	}

	// 提示用户输入图片名称前缀
	fmt.Print("\n请输入图片名称前缀: ")
	reader := bufio.NewReader(os.Stdin)
	prefix, _ := reader.ReadString('\n')
	prefix = strings.TrimSpace(prefix)

	if prefix == "" {
		fmt.Println("\n未输入名称前缀，使用默认前缀 'photo'")
		prefix = "photo"
	} else {
		fmt.Printf("\n使用名称前缀: '%s'\n", prefix)
	}

	// 获取起始编号
	startNum := 0          // 0表示使用默认补零模式
	useZeroPadding := true // 是否使用补零编号

	fmt.Print("\n请输入起始编号 (直接回车从001开始，输入数字如5则从5开始且不补零): ")
	numInput, _ := reader.ReadString('\n')
	numInput = strings.TrimSpace(numInput)

	if numInput == "" {
		fmt.Println("\n使用默认起始编号: 001 (补零格式)")
		startNum = 1
		useZeroPadding = true
	} else {
		// 尝试转换为整数
		num, err := strconv.Atoi(numInput)
		if err != nil {
			fmt.Println("\n请输入有效的数字!")
			return
		}

		if num < 1 {
			fmt.Println("\n起始编号必须大于0!")
			return
		}

		startNum = num
		useZeroPadding = false // 输入数字时不补零
		fmt.Printf("\n起始编号设置为: %d (不补零)\n", startNum)
	}

	// 显示命名规则
	if useZeroPadding {
		fmt.Printf("\n文件命名规则: %s-%03d.xxx (补零格式)\n", prefix, startNum)
	} else {
		fmt.Printf("\n文件命名规则: %s-%d.xxx (不补零)\n", prefix, startNum)
	}

	// 询问是否启用压缩
	fmt.Print("\n是否启用图片压缩? (输入 y 压缩并转换为JPG，输入 n 或回车跳过压缩): ")
	compressInput, _ := reader.ReadString('\n')
	compressInput = strings.TrimSpace(compressInput)

	var convertToJPG bool
	var quality int

	if strings.ToLower(compressInput) == "y" {
		// 压缩模式
		quality = getQualityFromUser(reader)
		convertToJPG = true
		fmt.Printf("\n启用图片压缩 (质量: %d%%)\n\n", quality)
	} else {
		// 不压缩模式
		quality = 100 // 不压缩时使用最高质量
		fmt.Println("\n禁用图片压缩")

		// 询问是否转换为JPG格式
		fmt.Print("\n是否将所有文件转为JPG格式? (输入 y 转换，输入 n 或回车保持原格式): ")
		convertInput, _ := reader.ReadString('\n')
		convertInput = strings.TrimSpace(convertInput)

		convertToJPG = (strings.ToLower(convertInput) == "y")
		if convertToJPG {
			fmt.Println("\n将转换为JPG格式 (不压缩)")
		} else {
			fmt.Println("\n保持原格式 (不转换)")
		}
		fmt.Println()
	}

	// 处理图片
	err = processImages(currentDir, prefix, startNum, convertToJPG, quality, useZeroPadding)
	if err != nil {
		fmt.Printf("\n处理图片失败: %v\n", err)
	} else {
		fmt.Println("\n图片处理完成! 原始文件保留，已生成新的文件到_NewFile文件夹。")
	}
}
