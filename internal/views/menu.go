package views

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type MenuModel struct {
	cursor  int
	options []string
}

func NewMenuModel() MenuModel {
	return MenuModel{
		options: []string{
			"📷 1.批量重命名照片（压缩/后缀）",
			"📁 2.批量重命名文件夹（新建/后缀）",
			"📊 3.读取 Excel 表头创建文件夹",
			"🏗️ 4.创建自定义系统目录结构",
			"💬 5.微信多开",
			"🗄️ 6.启动数据库服务",
			"❌ 7.退出",
		},
	}
}

func (m MenuModel) Init() tea.Cmd { return nil }

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// msg 是接口类型，通过类型断言判断具体是哪种消息
	switch msg := msg.(type) {

	// tea.KeyMsg 表示用户按下了键盘按键
	case tea.KeyMsg:

		// msg.String() 返回按键名称，如 "up"、"down"、"enter"、"q" 等
		switch msg.String() {

		// 退出程序
		case "ctrl+c", "q":
			return m, tea.Quit

			// 按 ↑ 或 k 键 → 向上移动光标（菜单上移）
		case "up", "k":
			// 如果光标不在第一项，则光标索引减 1
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			// 如果光标不在最后一项，则光标索引加 1
			// len(m.options)-1 是最后一项的索引
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}

			// 数字键直接跳转（1-7）
		case "1":
			return NewPhotoModel(), nil
		case "2":
			return NewFolderModel(), nil
		case "3":
			return NewExcelModel(), nil
		case "4":
			return NewStructureModel(), nil
		case "5":
			return NewWechatModel(), nil
		case "6":
			return NewDBModel(), nil
		case "7":
			return m, tea.Quit

		// 按 Enter 键 → 确认选择，跳转到对应的工具
		case "enter":
			// 根据当前光标所在位置（m.cursor）决定跳转到哪个工具
			switch m.cursor {
			case 0: // 第一项：照片重命名
				// NewPhotoModel() 返回照片工具的初始状态
				// 返回该模型后，Bubble Tea 会立即切换到照片工具的界面
				return NewPhotoModel(), nil
			case 1: // 第二项：文件夹重命名
				return NewFolderModel(), nil
			case 2: // 第三项：读取 Excel 创建文件夹
				return NewExcelModel(), nil
			case 3: // 第四项：创建自定义目录结构
				return NewStructureModel(), nil
			case 4: // 第五项：微信多开
				return NewWechatModel(), nil
			case 5: // 第六项：启动数据库服务
				return NewDBModel(), nil
			case 6: // 第七项：退出
				return m, tea.Quit

			}
		}

	}
	return m, nil
}

func (m MenuModel) View() string {
	// 使用 strings.Builder 高效拼接字符串（比 += 性能更好）
	s := "\n" // 顶部空一行，让界面不贴边

	// 标题行
	s += "  📦 kits 交互式工具集\n"

	// 分隔线：重复 30 次 "─" 字符
	s += "  " + strings.Repeat("─", 30) + "\n\n"

	// 遍历所有菜单选项
	for i, opt := range m.options {
		cursor := "  "     // 默认两个空格（未选中状态）
		if m.cursor == i { // 如果当前选项是光标所在位置
			cursor = "👉 " // 显示手指图标（选中状态）
		}
		// 拼接：光标标记 + 选项文字 + 换行
		// 例如："👉 📷 批量重命名照片\n"
		s += cursor + opt + "\n"
	}

	// 底部操作提示
	// 更新提示：加入数字快捷键
	s += "\n  ↑ ↓ 选择  •  Enter 确认  •  1-7 数字直达  •  q 退出\n"
	return s
}
