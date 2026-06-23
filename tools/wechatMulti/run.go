package wechatMulti

import (
	"fmt"
)

func Run() {
	const instanceCount = 2
	// 1.加载配置，如果不存在则引导用户首次配置
	cfg := LoadConfig()
	// 2. 并发启动微信实例（无交互，极速）
	fmt.Printf("正在快速启动 %d 个微信...\n", instanceCount)

	StartInstances(cfg.WeChatPath, instanceCount)
	fmt.Println("启动完成！")
}
