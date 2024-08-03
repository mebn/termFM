package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mebn/termfm/internal/ui"
)

func main() {
	// go PlayStation("http://rssand.se:4027/;")

	model := ui.InitialModel()
	p := tea.NewProgram(&model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
