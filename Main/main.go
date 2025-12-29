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
	HeadersPage
	QueryParamsPage
	LoadingPage
)

type model struct {
	NewApiInput        textinput.Model
	NewCollectionInput textinput.Model
	storage            Storage
	SelectedCollection Collection
	collectionIndex    int
	SelectedApi        Api
	Apis               []Api
	Collections        []Collection
	CurrentPage        View
	termWidth          int
	termHeight         int
	pointer            int
	jsonInput          textinput.Model
	apiViewport        viewport.Model
	viewportReady      bool
	editingApi         textinput.Model
	editingCollection  textinput.Model
	editingCurrentApi  textinput.Model
	editing            bool

	addHeaderKey   textinput.Model
	addHeaderValue textinput.Model
	editingHeader  textinput.Model
	Headers        []Header
	ApiIndex       int

	newBodyFieldInput   textinput.Model
	bodyFiledValueInput textinput.Model
	editingBodyFields   textinput.Model
	BodyFields          []BodyField

	addQueryParamsKey  textinput.Model
	addQueryParmsValue textinput.Model
	editingQueryParams textinput.Model
	QueryParams        []QueryParam

	apiResponse ApiResponse
}

func NewModel(storage Storage) model {
	ti := textinput.New()
	ti.Placeholder = "Enter JSON Body here..."
	ti.Focus()

	ai := textinput.New()
	ai.Placeholder = "Add New Api..."

	collInput := textinput.New()
	collInput.Placeholder = "Add New Collection..."

	addHeaderKey := textinput.New()
	addHeaderKey.Placeholder = "Add Header Key..."

	addHeaderValue := textinput.New()
	addHeaderValue.Placeholder = "Add Header Value..."

	newBodyField := textinput.New()
	newBodyField.Placeholder = "Add New Body Field..."

	bodyFiledValue := textinput.New()
	bodyFiledValue.Placeholder = "Add new Body Field Value..."

	QueryParamsKey := textinput.New()
	QueryParamsKey.Placeholder = "Add new Query Param Key..."

	QueryParamsValue := textinput.New()
	QueryParamsValue.Placeholder = "Add new Query Params Value..."

	return model{
		CurrentPage:         HomePage,
		jsonInput:           ti,
		viewportReady:       false,
		NewApiInput:         ai,
		NewCollectionInput:  collInput,
		storage:             storage,
		Collections:         storage.Collections,
		addHeaderKey:        addHeaderKey,
		addHeaderValue:      addHeaderValue,
		newBodyFieldInput:   newBodyField,
		bodyFiledValueInput: bodyFiledValue,
		addQueryParamsKey:   QueryParamsKey,
		addQueryParmsValue:  QueryParamsValue,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func main() {

	storage := ReadFile()

	m := NewModel(storage)
	p := tea.NewProgram(m, tea.WithAltScreen())

	watcher := watchFile(p)
	defer watcher.Close()

	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
