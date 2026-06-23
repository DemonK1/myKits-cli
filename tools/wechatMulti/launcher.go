package wechatMulti

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// 配置结构

type Config struct {
	WeChatPath string `json:"wechat_path"`
}

const configFile = "wechat_config.json"

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

// LoadConfig 加载配置：有则直接返回，无则引导用户输入并保存
func LoadConfig() Config {

	// 获取最终路径
	cfgPath := configFilePath()
	// 尝试读取配置文件
	file, err := os.ReadFile(cfgPath)
	if err == nil {
		var cfg Config

		if json.Unmarshal(file, &cfg) == nil && IsValidWeChat(cfg.WeChatPath) {
			return cfg
		}
		fmt.Println("已保存的微信路径无效或不是微信程序，请重新配置。")
	}
	// 首次配置或路径失效
	gopher := "Powered by Go"
	fmt.Println(colorString(cyan, gopher))
	fmt.Println()
	fmt.Println(colorString(yellow, "微信多开助手 v1.0"))
	fmt.Println()
	// 【标题：红色加粗】首次使用必须配置微信程序路径（仅需配置1次，后续启动不再重复设置）
	fmt.Println(colorString(red+bold, "== 【首次初始化】微信路径配置，后续无需重复设置 =="))
	// 【操作步骤：绿色加粗】分步傻瓜指引
	fmt.Println(colorString(green+bold, "操作步骤："))
	fmt.Println(colorString(green, "1. 找到桌面微信图标 → 右键【属性】"))
	fmt.Println(colorString(green, "2. 在【目标】栏全选复制完整路径"))
	fmt.Println(colorString(green, "3. 粘贴到此"))
	fmt.Println(colorString(green, "参考格式示例：E:\\Apps\\WeChat\\WeChat.exe"))
	fmt.Println("")
	fmt.Print("请输入微信完整路径: ")
	var path string
	fmt.Scan(&path)
	for !IsValidWeChat(path) {
		fmt.Printf("路径无效或不是微信程序: %s\n", path)
		fmt.Print("请重新输入: ")
		fmt.Scanln(&path)
	}

	// 	保存配置
	cfg := Config{WeChatPath: path}
	indent, _ := json.MarshalIndent(cfg, "", "  ")

	os.WriteFile(cfgPath, indent, 0644)
	fmt.Println("配置已保存，下次启动将直接多开微信。")
	return cfg
}

// IsValidWeChat 验证给定路径是否指向微信程序
func IsValidWeChat(path string) bool {

	// 1. 路径不能为空
	if strings.TrimSpace(path) == "" {
		return false
	}

	stat, err := os.Stat(path)
	if err != nil || stat.IsDir() {
		return false
	}
	// 3. 文件名必须是 WeChat.exe 或 Weixin.exe（忽略大小写）
	name := strings.ToLower(filepath.Base(path))
	return name == "wechat.exe" || name == "weixin.exe"
}

// StartInstances 并发启动指定数量的进程
func StartInstances(path string, count int) {
	var wg sync.WaitGroup

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cmd := exec.Command(path)
			if err := cmd.Start(); err != nil {
				fmt.Printf("启动失败: %v\n", err)
			}
		}()
	}
	wg.Wait()
}

// 获取配置文件存放的完整路径，并确保目录存在
func configFilePath() string {
	// 返回系统认定的配置目录，无需自己拼接路径（C:\Users\用户名\AppData\Roaming）
	cfgDri, err := os.UserConfigDir()
	if err != nil {
		// 万一获取失败（极罕见），回退到 exe 所在目录
		executable, _ := os.Executable()
		return filepath.Join(filepath.Dir(executable), configFile)
	}
	// 在配置根目录下创建专属子文件夹
	appDri := filepath.Join(cfgDri, "WeChatMultiOpener")
	if errs := os.MkdirAll(appDri, 0700); errs != nil {
		// 创建失败也回退到 exe 所在目录
		exePath, _ := os.Executable()
		return filepath.Join(filepath.Dir(exePath), configFile)
	}
	return filepath.Join(appDri, configFile)
}

// 加点颜色
func colorString(color, text string) string {
	return color + text + reset
}
