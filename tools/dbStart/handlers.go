package dbStart

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// 颜色码
const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	cyan   = "\033[36m"
	bold   = "\033[1m"
)

// DBConfig 改为存储完整命令行
type DBConfig struct {
	Name    string `json:"name"`
	CmdLine string `json:"cmdLine"` // 完整启动命令（含参数）
}

// 加点颜色
func colorString(color, text string) string {
	return color + text + reset
}

// configDir 返回配置根目录
func configDir() (string, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cfgDir, "myKits_DBServe_Config"), nil
}

// listConfigs 读取所有已保存的数据库配置
func listConfigs() ([]DBConfig, error) {
	root, err := configDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var configs []DBConfig
	for _, e := range entries {
		if !e.IsDir() || !strings.HasSuffix(e.Name(), "_Opener") {
			continue
		}
		configFile := filepath.Join(root, e.Name(), "config.json")
		data, err := os.ReadFile(configFile)
		if err != nil {
			continue
		}
		var cfg DBConfig
		if err := json.Unmarshal(data, &cfg); err != nil {
			continue
		}
		configs = append(configs, cfg)
	}
	return configs, nil
}

// saveConfig 保存数据库配置到对应的 _Opener 文件夹
func saveConfig(cfg DBConfig) error {
	root, err := configDir()
	if err != nil {
		return err
	}
	dir := filepath.Join(root, cfg.Name+"_Opener")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "config.json"), data, 0644)
}

// startService 启动数据库服务
func startService(cfg DBConfig) error {
	cmd := exec.Command("cmd", "/C", cfg.CmdLine) // Windows
	// Linux: exec.Command("sh", "-c", cfg.CmdLine)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动失败: %w", err)
	}
	fmt.Printf("✅ 服务已启动 (PID: %d)\n", cmd.Process.Pid)
	return nil
}

// serviceMenu 二级菜单：针对选定数据库的操作
func serviceMenu(cfg DBConfig, reader *bufio.Reader) {
	for {
		fmt.Println("\n数据库操作:")
		fmt.Println("  1. 开启服务")
		fmt.Println("  2. 关闭服务")
		fmt.Println("  3. 返回上级")
		fmt.Print("请选择: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			if err := startService(cfg); err != nil {
				fmt.Printf("❌ %v\n", err)
			}
		case "2":
			if err := stopService(cfg); err != nil {
				fmt.Printf("❌ %v\n", err)
			}
		case "3":
			return
		default:
			fmt.Println("无效选择")
		}
	}
}

// extractFeature 从路径提取特征名
func extractFeature(path string) string {
	base := filepath.Base(path)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	return strings.ToLower(name)
}

// listExecutables 列出目录下的可执行文件（跨平台）
func listExecutables(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var exes []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if os.PathSeparator == '\\' {
			ext := strings.ToLower(filepath.Ext(name))
			if ext == ".exe" || ext == ".bat" || ext == ".cmd" {
				exes = append(exes, name)
			}
		} else {
			info, err := e.Info()
			if err != nil {
				continue
			}
			if info.Mode()&0111 != 0 {
				exes = append(exes, name)
			}
		}
	}
	return exes, nil
}

// describeExe 尝试获取可执行文件的版本描述
func describeExe(path string) string {
	tryArgs := [][]string{
		{"--version"},
		{"-V"},
		{"version"},
	}
	for _, args := range tryArgs {
		cmd := exec.Command(path, args...)
		out, err := cmd.CombinedOutput()
		if err == nil && len(out) > 0 {
			return strings.TrimSpace(string(out))
		}
	}
	return "无法获取描述"
}

// addDatabase 交互式添加新数据库（改进版）
func addDatabase(reader *bufio.Reader) error {
	// 第一步：输入 bin 目录

	fmt.Print(colorString(yellow+bold, "请输入数据库可执行文件所在目录（如 PostgreSQL 的 bin 目录）: "))
	dirPath, _ := reader.ReadString('\n')
	dirPath = strings.TrimSpace(dirPath)
	if dirPath == "" {
		return fmt.Errorf("路径不能为空")
	}
	info, err := os.Stat(dirPath)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("目录不存在或不是目录: %s", dirPath)
	}

	// 第二步：列出可执行文件
	exes, err := listExecutables(dirPath)
	if err != nil {
		return fmt.Errorf("读取目录失败: %w", err)
	}
	if len(exes) == 0 {
		return fmt.Errorf("该目录下未找到任何可执行文件")
	}

	fmt.Println(colorString(red+bold, "\n找到以下可执行文件（可查看描述）:"))
	for i, name := range exes {
		fullPath := filepath.Join(dirPath, name)
		desc := describeExe(fullPath)
		fmt.Printf(colorString(green+bold, "  %d. %-30s → %s\n"), i+1, name, desc)
	}
	fmt.Print(colorString(cyan+bold, "\n请选择要使用的可执行文件编号（仅作参考，后续需输入完整命令）: "))
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	idx, err := parseInt(choice)
	if err != nil || idx < 1 || idx > len(exes) {
		return fmt.Errorf("无效选择")
	}
	selectedFile := exes[idx-1]
	fullPath := filepath.Join(dirPath, selectedFile)

	// 第三步：提取特征名作为默认名称，允许修改
	defaultFeature := extractFeature(selectedFile)
	fmt.Printf("\n自动识别的数据库名称: [%s]\n", defaultFeature)
	fmt.Print("如需修改请输入新名称（直接回车使用自动名称）: ")
	customName, _ := reader.ReadString('\n')
	customName = strings.TrimSpace(customName)
	if customName != "" {
		defaultFeature = customName
	}

	// 第四步：提示输入完整启动命令
	fmt.Println("\n现在请输入完整的启动命令（包括路径和参数）。")
	fmt.Printf("示例: \"%s\" -D \"C:\\pgdata\"\n", fullPath)
	fmt.Print("完整启动命令: ")
	cmdLine, _ := reader.ReadString('\n')
	cmdLine = strings.TrimSpace(cmdLine)
	if cmdLine == "" {
		return fmt.Errorf("启动命令不能为空")
	}

	cfg := DBConfig{
		Name:    defaultFeature,
		CmdLine: cmdLine,
	}
	if err := saveConfig(cfg); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}
	fmt.Printf("✅ 已添加数据库: %s (配置文件夹: %s_Opener)\n", cfg.Name, cfg.Name)
	return nil
}

// parseInt 字符串转整数（简单实现）
func parseInt(s string) (int, error) {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("非数字字符")
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}

// stopService 尝试关闭数据库服务（通过进程名终止）
func stopService(cfg DBConfig) error {
	// 从命令行提取可执行文件名（去除参数和引号）
	parts := strings.Fields(cfg.CmdLine)
	if len(parts) == 0 {
		return fmt.Errorf("命令行格式错误")
	}
	exePath := strings.Trim(parts[0], "\"") // 去掉可能存在的双引号
	name := filepath.Base(exePath)          // 获取文件名，如 pg_ctl.exe

	var cmd *exec.Cmd
	if os.PathSeparator == '\\' { // Windows
		cmd = exec.Command("taskkill", "/IM", name, "/F")
	} else { // Unix/Linux/macOS
		cmd = exec.Command("pkill", "-f", name)
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("关闭失败: %v, 输出: %s", err, string(output))
	}
	fmt.Printf("✅ 已尝试关闭 %s\n", name)
	return nil
}

// showList 显示所有已保存的数据库，返回配置切片供后续选择
func showList() ([]DBConfig, error) {
	configs, err := listConfigs()
	if err != nil {
		return nil, err
	}
	if len(configs) == 0 {
		fmt.Println(colorString(blue+bold, "暂无已保存的数据库。"))
		return nil, nil
	}

	fmt.Println("\n已保存的数据库列表:")
	for i, c := range configs {
		fmt.Printf("  %d. %-20s → %s\n", i+1, c.Name, c.CmdLine)
	}
	return configs, nil
}
