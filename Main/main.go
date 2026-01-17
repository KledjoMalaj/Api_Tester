package main

import (
	"fmt"
	"log"
	"os"

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
	ResponsePage
	VariablesPage
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

	editingApi        textinput.Model
	editingCollection textinput.Model
	editingCurrentApi textinput.Model
	editing           bool

	addHeaderKey   textinput.Model
	addHeaderValue textinput.Model
	editingHeader  textinput.Model
	Headers        []Header
	ApiIndex       int

	newBodyFieldInput   textinput.Model
	bodyFiledValueInput textinput.Model
	editingBodyFields   textinput.Model
	BodyFields          []BodyField

	addQueryParamsKey   textinput.Model
	addQueryParamsValue textinput.Model
	editingQueryParams  textinput.Model
	QueryParams         []QueryParam

	Responses        []Response
	LocalVariables   []LocalVariable
	VariablesFocus   bool
	addVariableKey   textinput.Model
	addVariableValue textinput.Model

	apiResponse ApiResponse

	errorMessage string
	hasError     bool
}

func NewModel(storage Storage) model {
	ti := textinput.New()
	ti.Placeholder = "Enter JSON Body here..."
	ti.CharLimit = 50
	ti.Focus()

	ai := textinput.New()
	ai.Placeholder = "Add New Api..."
	ai.Width = 50

	collInput := textinput.New()
	collInput.Placeholder = "Add New Collection..."
	collInput.Width = 50

	addHeaderKey := textinput.New()
	addHeaderKey.Placeholder = "Add Header Key..."
	addHeaderKey.Width = 50

	addHeaderValue := textinput.New()
	addHeaderValue.Placeholder = "Add Header Value..."
	addHeaderValue.Width = 50

	newBodyField := textinput.New()
	newBodyField.Placeholder = "Add Body Field..."
	newBodyField.Width = 50

	bodyFiledValue := textinput.New()
	bodyFiledValue.Placeholder = "Add Body Field Value..."
	bodyFiledValue.Width = 50

	QueryParamsKey := textinput.New()
	QueryParamsKey.Placeholder = "Add Query Param Key..."
	QueryParamsKey.Width = 50

	QueryParamsValue := textinput.New()
	QueryParamsValue.Placeholder = "Add Query Params Value..."
	QueryParamsValue.Width = 50

	VariableKey := textinput.New()
	VariableKey.Placeholder = "Add New Variable Key..."
	VariableKey.Width = 50

	VariableValue := textinput.New()
	VariableValue.Placeholder = "Add New Variable Value..."
	VariableValue.Width = 50

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
		addQueryParamsValue: QueryParamsValue,
		addVariableKey:      VariableKey,
		addVariableValue:    VariableValue,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func main() {

	storage, err := ReadFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Cannot start application\n")
		fmt.Fprintf(os.Stderr, "Reason: %v \n", err)
		os.Exit(1)
	}

	m := NewModel(storage)
	p := tea.NewProgram(m, tea.WithAltScreen())

	watcher, err := watchFile(p)
	if err != nil {
		log.Printf("Warning:File watcher failed: %v", err)
	} else {
		defer watcher.Close()
	}

	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
