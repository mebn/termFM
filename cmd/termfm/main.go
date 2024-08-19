package main

import (
	"errors"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mebn/termfm/internal/cli"
	"github.com/mebn/termfm/internal/ui"
)

func main() {
	// handle potential command line arguments
	config := cli.NewConfig()
	err := config.HandleConfig(os.Args[1:])

	if err != nil {
		if errors.Is(err, cli.FlagNotFoundError) {
			fmt.Fprintln(os.Stderr, err)
		}

		fmt.Fprint(os.Stderr, "\n", cli.HowToUse, "\n")
		os.Exit(1)
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
