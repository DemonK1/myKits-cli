package models

type ViewType int

const (
	MenuView ViewType = iota
	PhotoView
	FolderView
	ExcelView
	StructureView
	WechatView
	DBView
)

// MenuModel 主菜单模型
type MenuModel struct {
	Cursor  int
	Options []string
}
