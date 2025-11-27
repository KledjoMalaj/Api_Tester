package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type View int

const (
	HomePage View = iota
	ApiPage
	RequestPage
)

type model struct {
	SelectedApi   Api
	Options       []Api
	CurrentPage   View
	termWidth     int
	termHeight    int
	pointer       int
	jsonInput     textinput.Model
	apiViewport   viewport.Model
	viewportReady bool
}

func NewModel(options []Api) model {
	ti := textinput.New()
	ti.Placeholder = "Enter JSON Body here..."
	ti.Focus()
	return model{
		CurrentPage:   HomePage,
		Options:       options,
		jsonInput:     ti,
		viewportReady: false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func main() {

	Options := ReadFile()
	m := NewModel(Options)
	p := tea.NewProgram(m, tea.WithAltScreen())

	watcher := watchFile(p)
	defer watcher.Close()

	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
