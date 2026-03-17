package main

import (
	"GoTuiFrontend/operations"
	"GoTuiFrontend/tui"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	storage, err := operations.ReadFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Cannot start application\n")
		fmt.Fprintf(os.Stderr, "Reason: %v \n", err)
		os.Exit(1)
	}

	m := tui.NewModel(storage)
	p := tea.NewProgram(m, tea.WithAltScreen())

	watcher, err := operations.WatchFile(p)
	if err != nil {
		log.Printf("Warning:File watcher failed: %v", err)
	} else {
		defer watcher.Close()
	}

	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
