package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mebn/termfm/internal/cli"
	"github.com/mebn/termfm/internal/ui"
)

func main() {
	// handle potential command line arguments
	config := cli.NewConfig()
	configStatus := config.HandleConfig(os.Args[1:])

	if configStatus != cli.ConfigOk {
		switch configStatus {
		case cli.ConfigFlagNotFound:
			fmt.Fprintln(os.Stderr, "Invalid flag used.")
		}

		fmt.Fprint(os.Stderr, cli.HowToUse)
	}

	if config.Cli {
		config.ShowCLI()
	} else {
		handleTUI()
	}

}

func handleTUI() {
	model := ui.InitialModel()
	p := tea.NewProgram(&model)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
