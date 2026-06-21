package main

import (
	"fmt"
	"myKits-cli/internal/views"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(views.NewMenuModel())
	_, err := p.Run()
	if err != nil {
		fmt.Printf("程序运行出错: %v\n", err)
	}
	os.Exit(1)
}
