package tui

import (
	"GoTuiFrontend/models"

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
	DashBoard
)

type Model struct {
	NewApiInput        textinput.Model
	NewCollectionInput textinput.Model
	apiMethodOptions   []string
	apiMethodIndex     int
	apiMethodSelecting bool
	storage            models.Storage
	SelectedCollection models.Collection
	collectionIndex    int
	SelectedApi        models.Api
	Apis               []models.Api
	Collections        []models.Collection
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
	Headers        []models.Header
	ApiIndex       int

	newBodyFieldInput   textinput.Model
	bodyFiledValueInput textinput.Model
	editingBodyFields   textinput.Model
	BodyFields          []models.BodyField

	addQueryParamsKey   textinput.Model
	addQueryParamsValue textinput.Model
	editingQueryParams  textinput.Model
	QueryParams         []models.QueryParam

	Responses             []models.Response
	ResponseExpanded      map[string]bool
	LocalVariables        []models.LocalVariable
	VariablesFocus        bool
	addVariableKey        textinput.Model
	addVariableValue      textinput.Model
	editingLocalVariables textinput.Model

	ApiResponse models.ApiResponse

	errorMessage string
	hasError     bool

	responseScrollOffset int
	variableScrollOffset int

	pageScrollOffset int

	responseComponent    bool
	reqPageComponent     bool
	headersPageComponent bool
	paramsComponent      bool
	resComponent         bool

	secondPointer int
}

func NewModel(storage models.Storage) Model {
	ti := textinput.New()
	ti.Placeholder = "Enter JSON Body here..."
	ti.CharLimit = 50
	ti.Focus()

	ai := textinput.New()
	ai.Placeholder = "Enter request URL..."
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

	return Model{
		CurrentPage:          HomePage,
		jsonInput:            ti,
		viewportReady:        false,
		NewApiInput:          ai,
		NewCollectionInput:   collInput,
		apiMethodOptions:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		storage:              storage,
		Collections:          storage.Collections,
		addHeaderKey:         addHeaderKey,
		addHeaderValue:       addHeaderValue,
		newBodyFieldInput:    newBodyField,
		bodyFiledValueInput:  bodyFiledValue,
		addQueryParamsKey:    QueryParamsKey,
		addQueryParamsValue:  QueryParamsValue,
		addVariableKey:       VariableKey,
		addVariableValue:     VariableValue,
		ResponseExpanded:     map[string]bool{"$": true},
		responseComponent:    false,
		reqPageComponent:     false,
		headersPageComponent: false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
