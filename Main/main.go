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
	CollectionPage
	ApiPage
	RequestPage
)

type model struct {
	NewApiInput        textinput.Model
	storage            Storage
	SelectedCollection Collection
	collectionIndex    int
	SelectedApi        Api
	Apis               []Api
	CurrentPage        View
	termWidth          int
	termHeight         int
	pointer            int
	jsonInput          textinput.Model
	apiViewport        viewport.Model
	viewportReady      bool
	editingApi         textinput.Model
	editing            bool
}

func NewModel(storage Storage) model {
	ti := textinput.New()
	ti.Placeholder = "Enter JSON Body here..."
	ti.Focus()

	ai := textinput.New()
	ai.Placeholder = "Add New Api..."

	return model{
		CurrentPage:   HomePage,
		jsonInput:     ti,
		viewportReady: false,
		NewApiInput:   ai,
		storage:       storage,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func main() {

	storage := ReadFilenew()

	m := NewModel(storage)
	p := tea.NewProgram(m, tea.WithAltScreen())

	watcher := watchFile(p, m.collectionIndex)
	defer watcher.Close()

	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
