package tui

import (
	"GoTuiFrontend/models"
	"GoTuiFrontend/operations"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case models.ErrorMsg:
		m.errorMessage = msg.Message
		m.hasError = true
		return m, nil

	case models.ApiResponseMsg:
		m.ApiResponse = msg.Response
		m.CurrentPage = ApiPage
		if m.viewportReady {
			m.apiViewport.SetContent(BuildApiPageContent(m, m.termWidth))
			m.apiViewport.GotoTop()
		}
		return m, nil

	case operations.FileChangedMsg:
		m.storage = models.Storage(msg)
		m.Collections = m.storage.Collections

		if m.CurrentPage == CollectionPage || m.CurrentPage == HeadersPage ||
			m.CurrentPage == RequestPage || m.CurrentPage == QueryParamsPage ||
			m.CurrentPage == ResponsePage || m.CurrentPage == VariablesPage {

			if m.collectionIndex >= 0 && m.collectionIndex < len(m.Collections) {
				m.SelectedCollection = m.Collections[m.collectionIndex]
				m.Apis = m.SelectedCollection.Requests
				m.LocalVariables = m.SelectedCollection.LocalVariables

				if m.ApiIndex >= 0 && m.ApiIndex < len(m.Apis) {
					m.SelectedApi = m.Apis[m.ApiIndex]

					m.Headers = m.SelectedApi.Headers
					m.BodyFields = m.SelectedApi.BodyField
					m.QueryParams = m.SelectedApi.QueryParams
				}
			}
		}

	case tea.WindowSizeMsg:

		m.termWidth = msg.Width
		m.termHeight = msg.Height

		// Initialize viewport when we have terminal dimensions
		if !m.viewportReady {
			m.apiViewport = viewport.New(msg.Width, msg.Height-4)
			m.viewportReady = true
		} else {
			m.apiViewport.Width = msg.Width
			m.apiViewport.Height = msg.Height - 4
		}

		// Update viewport content if we're on ApiPage
		if m.CurrentPage == ApiPage {
			m.apiViewport.SetContent(BuildApiPageContent(m, m.termWidth))
		}

	case tea.KeyMsg:
		switch m.CurrentPage {
		case HomePage:
			m, cmd := UpdateHomePage(m, msg)
			return m, cmd
		case CollectionPage:
			m, cmd := UpdateCollectionPage(m, msg)
			return m, cmd
		case ApiPage:
			m, cmd := UpdateApiPage(m, msg)
			return m, cmd
		case RequestPage:
			m, cmd := UpdateReqPage(m, msg)
			return m, cmd
		case HeadersPage:
			m, cmd := UpdateHeadersPage(m, msg)
			return m, cmd
		case QueryParamsPage:
			m, cmd := UpdateQueryParamsPage(m, msg)
			return m, cmd
		case LoadingPage:
			m, cmd := UpdateLoadingPage(m, msg)
			return m, cmd
		case ResponsePage:
			m, cmd := UpdateResponsePage(m, msg)
			return m, cmd
		case VariablesPage:
			m, cmd := UpdateVariablesPage(m, msg)
			return m, cmd
		case DashBoard:
			m, cmd := UpdateDashBoard(m, msg)
			return m, cmd
		}
	}

	return m, nil
}

func UpdateHomePage(m Model, msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:

		if m.editing {
			switch msg.String() {

			case "esc":
				m.editingCollection.Blur()
				m.editing = false
			case "enter":
				if err := operations.EditCollection(m.storage, m.SelectedCollection, m.editingCollection.Value()); err != nil {
					return m, models.ShowErrorCommand("Failed to edit Collection: " + err.Error())
				}
				m.editingApi.Blur()
				m.editing = false
			}

			m.editingCollection, cmd = m.editingCollection.Update(msg)
			return m, cmd
		}

		if m.NewCollectionInput.Focused() {
			switch msg.String() {
			case "esc":
				m.NewCollectionInput.Blur()
				return m, nil
			case "enter":

				if err := operations.AddCollection(m.storage, m.Collections, m.NewCollectionInput.Value()); err != nil {
					return m, models.ShowErrorCommand("Failed to add collection: " + err.Error())
				}
				m.NewCollectionInput.SetValue("")
				m.NewCollectionInput.Blur()

			}
			m.NewCollectionInput, cmd = m.NewCollectionInput.Update(msg)
			return m, cmd
		}

		switch msg.String() {

		case "esc":
			return m, tea.Quit

		case "up", "k":
			if m.pointer > 0 {
				m.pointer--

				if m.pointer < m.pageScrollOffset {
					m.pageScrollOffset--
				}
			}
		case "down", "j":
			if m.pointer < len(m.storage.Collections)-1 {
				m.pointer++

				maxVisible := 12
				if m.pointer >= m.pageScrollOffset+maxVisible {
					m.pageScrollOffset++
				}
			}
		case "enter":
			m.CurrentPage = CollectionPage
			m.SelectedCollection = m.storage.Collections[m.pointer]
			m.Apis = m.SelectedCollection.Requests
			m.LocalVariables = m.SelectedCollection.LocalVariables
			m.collectionIndex = m.pointer
			m.pointer = 0
			m.pageScrollOffset = 0

		case "t":
			m.CurrentPage = DashBoard
			m.SelectedCollection = m.storage.Collections[m.pointer]
			m.pointer = 0

		case ":":
			m.NewCollectionInput.Focus()
			return m, nil

		case "d":
			if len(m.Collections) > 0 {
				selectedCollection := m.storage.Collections[m.pointer]
				newCollections, err := operations.DeleteCollection(selectedCollection, m.storage)
				if err != nil {
					return m, models.ShowErrorCommand("Failed to delete collection: " + err.Error())
				}
				m.Collections = newCollections
				if m.pointer >= len(m.Collections) && m.pointer > 0 {
					m.pointer--
				}
			}

		case "e":
			m.editing = true
			m.editingCollection = textinput.New()
			m.SelectedCollection = m.Collections[m.pointer]
			m.editingCollection.SetValue(m.SelectedCollection.Name)
			m.editingCollection.Focus()

		case "x":
			if m.hasError {
				m.hasError = false
				m.errorMessage = ""
				return m, nil
			}
		}

	}
	return m, nil
}

func UpdateCollectionPage(m Model, msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.editing {
			switch msg.String() {
			case "enter":
				if err := operations.EditApi(m.storage, m.collectionIndex, m.SelectedApi, m.editingApi.Value()); err != nil {
					return m, models.ShowErrorCommand("Failed to edit api: " + err.Error())
				}
				m.editingApi.Blur()
				m.editing = false
			case "esc":
				m.editingApi.Blur()
				m.editing = false
			}

			m.editingApi, cmd = m.editingApi.Update(msg)
			return m, cmd
		}

		if m.NewApiInput.Focused() {
			switch msg.String() {
			case "esc":
				m.NewApiInput.Blur()
				return m, nil
			case "enter":

				if err := operations.AddApi(m.storage, m.collectionIndex, m.Apis, m.NewApiInput.Value()); err != nil {
					return m, models.ShowErrorCommand("Failed to add api: " + err.Error())
				}
				m.NewApiInput.SetValue("")
				m.NewApiInput.Blur()
			}

			m.NewApiInput, cmd = m.NewApiInput.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "up", "k":
			if m.pointer > 0 {
				m.pointer--

				if m.pointer < m.pageScrollOffset {
					m.pageScrollOffset--
				}
			}
		case "down", "j":
			if m.pointer < len(m.Apis)-1 {
				m.pointer++

				maxVisible := 12
				if m.pointer >= m.pageScrollOffset+maxVisible {
					m.pageScrollOffset++
				}
			}
		case "enter":
			m.SelectedApi = m.Apis[m.pointer]

			processedApi := operations.ProcessRequest(m.SelectedApi, m.SelectedCollection.LocalVariables)

			switch processedApi.Method {
			case "POST", "DELETE", "PUT", "PATCH":
				m.SelectedApi = processedApi
				m.BodyFields = processedApi.BodyField
				m.ApiIndex = m.pointer
				m.CurrentPage = RequestPage
				m.pointer = 0

			case "GET":
				TuiModel := models.TuiModel{
					SelectedApi:        m.SelectedApi,
					SelectedCollection: m.SelectedCollection,
					LocalVariables:     m.LocalVariables,
				}
				m.CurrentPage = LoadingPage
				m.ApiIndex = m.pointer
				m.ApiResponse = operations.FetchData(m.SelectedApi, TuiModel)
				m.Responses, _ = operations.HandleJson(m.ApiResponse)
				return m, operations.FetchApiCommand(m.SelectedApi, TuiModel)
			}

		case ":":
			m.NewApiInput.Focus()
			return m, nil

		case "d":
			if len(m.Apis) > 0 {
				selectedApi := m.Apis[m.pointer]
				newApis, err := operations.DeleteApi(selectedApi, m.storage, m.collectionIndex)
				if err != nil {
					return m, models.ShowErrorCommand("Failed to delete api: " + err.Error())
				}
				m.Apis = newApis
				if m.pointer >= len(m.Apis) && m.pointer > 0 {
					m.pointer--
				}
			}

		case "e":
			m.editing = true
			m.editingApi = textinput.New()
			m.SelectedApi = m.Apis[m.pointer]
			m.editingApi.SetValue(m.SelectedApi.Method + " " + m.SelectedApi.Url)
			m.editingApi.Focus()

		case "esc":
			m.CurrentPage = HomePage
			m.pointer = m.collectionIndex
			m.pageScrollOffset = 0

		case "h":
			m.CurrentPage = HeadersPage
			m.SelectedApi = m.Apis[m.pointer]
			m.Headers = m.SelectedApi.Headers
			m.ApiIndex = m.pointer
			m.pointer = 0

		case "q":
			m.CurrentPage = QueryParamsPage
			m.SelectedApi = m.Apis[m.pointer]
			m.ApiIndex = m.pointer
			m.QueryParams = m.SelectedApi.QueryParams
			m.pointer = 0

		case "x":
			if m.hasError {
				m.hasError = false
				m.errorMessage = ""
				return m, nil
			}
		case "v":
			m.CurrentPage = VariablesPage
			m.LocalVariables = m.SelectedCollection.LocalVariables
			m.pointer = 0
		}
	}

	m.NewApiInput, cmd = m.NewApiInput.Update(msg)
	return m, cmd
}

func UpdateApiPage(m Model, msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:

		if m.editing {
			switch msg.String() {
			case "esc":
				m.editingCurrentApi.Blur()
				m.editing = false
				// Rebuild to hide the input
				if m.viewportReady {
					m.apiViewport.SetContent(BuildApiPageContent(m, m.termWidth))
				}
				return m, nil

			case "enter":
				if err := operations.EditApi(m.storage, m.collectionIndex, m.SelectedApi, m.editingCurrentApi.Value()); err != nil {
					return m, models.ShowErrorCommand("Failed to edit api: " + err.Error())
				}

				// Update local state
				m.storage, _ = operations.ReadFile()
				m.Collections = m.storage.Collections
				m.SelectedCollection = m.Collections[m.collectionIndex]
				m.Apis = m.SelectedCollection.Requests
				m.SelectedApi = m.Apis[m.pointer]

				m.editingCurrentApi.Blur()
				m.editing = false

				// Rebuild content ONLY here with the new API - this will re-fetch
				if m.viewportReady {
					m.apiViewport.SetContent(BuildApiPageContent(m, m.termWidth))
				}
				return m, nil
			}

			m.editingCurrentApi, cmd = m.editingCurrentApi.Update(msg)

			// Show typing but don't re-fetch API yet
			if m.viewportReady {
				m.apiViewport.SetContent(BuildApiPageContent(m, m.termWidth))
			}

			return m, cmd
		}

		switch msg.String() {
		case "esc":
			m.CurrentPage = CollectionPage
			m.pointer = m.ApiIndex
			return m, nil
		case "up", "k":
			m.apiViewport.LineUp(1)
		case "down", "j":
			m.apiViewport.LineDown(1)
		case "pgup", "b":
			m.apiViewport.ViewUp()
		case "pgdown", "f", " ":
			m.apiViewport.ViewDown()
		case "home", "g":
			m.apiViewport.GotoTop()
		case "end", "G":
			m.apiViewport.GotoBottom()

		case "e":
			m.editing = true
			m.editingCurrentApi = textinput.New()
			m.editingCurrentApi.SetValue(m.SelectedApi.Method + " " + m.SelectedApi.Url)
			m.editingCurrentApi.Focus()

			// Rebuild viewport to show the editing input
			if m.viewportReady {
				m.apiViewport.SetContent(BuildApiPageContent(m, m.termWidth))
			}
		case "r":
			m.LocalVariables = m.SelectedCollection.LocalVariables
			m.CurrentPage = ResponsePage
			m.pointer = 0
			m.responseScrollOffset = 0
			m.variableScrollOffset = 0
		}
	}

	m.apiViewport, cmd = m.apiViewport.Update(msg)
	return m, cmd
}

func UpdateReqPage(m Model, msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.editingBodyFields.Focused() {
			switch msg.String() {
			case "esc":
				m.editing = false
				m.editingBodyFields.Blur()
				return m, nil
			case "enter":
				m.BodyFields[m.pointer].Value = m.editingBodyFields.Value()
				newBodyFields, err := operations.AddBodyField(m.storage, m.collectionIndex, m.ApiIndex, m.BodyFields)
				if err != nil {
					return m, models.ShowErrorCommand("Failed to edit body field: " + err.Error())
				}
				m.BodyFields = newBodyFields
				m.editing = false
				m.editingBodyFields.Blur()
			}
			m.editingBodyFields, cmd = m.editingBodyFields.Update(msg)
			return m, cmd
		}

		if m.newBodyFieldInput.Focused() {
			switch msg.String() {
			case "esc":
				m.newBodyFieldInput.Blur()
				m.newBodyFieldInput.SetValue("")
			case "enter":
				newBodyFieldKey := m.newBodyFieldInput.Value()
				newBodyFiled := models.BodyField{
					Key:   newBodyFieldKey,
					Value: "",
				}
				m.BodyFields = append(m.BodyFields, newBodyFiled)
				newBodyFields, err := operations.AddBodyField(m.storage, m.collectionIndex, m.ApiIndex, m.BodyFields)
				if err != nil {
					return m, models.ShowErrorCommand("Failed to add body field: " + err.Error())
				}
				m.BodyFields = newBodyFields
				m.newBodyFieldInput.SetValue("")
				m.newBodyFieldInput.Blur()
			}
			m.newBodyFieldInput, cmd = m.newBodyFieldInput.Update(msg)
			return m, cmd
		}
		if m.bodyFiledValueInput.Focused() {
			switch msg.String() {
			case "esc":
				m.bodyFiledValueInput.Blur()
				m.bodyFiledValueInput.SetValue("")
			case "enter":
				newBodyFieldValue := m.bodyFiledValueInput.Value()
				m.BodyFields[m.pointer].Value = newBodyFieldValue
				_, err := operations.AddBodyField(m.storage, m.collectionIndex, m.ApiIndex, m.BodyFields)
				if err != nil {
					return m, models.ShowErrorCommand("Failed to add body field value: " + err.Error())
				}
				m.bodyFiledValueInput.SetValue("")
				m.bodyFiledValueInput.Blur()
			}
			m.bodyFiledValueInput, cmd = m.bodyFiledValueInput.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "enter":
			TuiModel := models.TuiModel{
				SelectedApi:        m.SelectedApi,
				SelectedCollection: m.SelectedCollection,
				LocalVariables:     m.LocalVariables,
			}

			m.pageScrollOffset = 0
			m.CurrentPage = LoadingPage
			m.ApiResponse = operations.PostAPiFunc(TuiModel)
			m.Responses, _ = operations.HandleJson(m.ApiResponse)
			return m, operations.PostApiCommand(TuiModel)

		case "v":
			m.bodyFiledValueInput.Focus()
		case ":":
			m.newBodyFieldInput.Focus()
		case "esc":
			if m.reqPageComponent {
				m.reqPageComponent = false
			} else {
				m.CurrentPage = CollectionPage
				m.pageScrollOffset = 0
			}
		case "up", "k":
			if m.pointer > 0 {
				m.pointer--

				if m.pointer < m.pageScrollOffset {
					m.pageScrollOffset--
				}
			}
		case "down", "j":
			if m.pointer < len(m.BodyFields)-1 {
				m.pointer++

				maxVisible := 12
				if m.pointer >= m.pageScrollOffset+maxVisible {
					m.pageScrollOffset++
				}
			}
		case "d":
			if len(m.BodyFields) > 0 {
				selectedBodyField := m.BodyFields[m.pointer]
				newBodyFields, err := operations.DeleteBodyField(selectedBodyField, m.storage, m.collectionIndex, m.ApiIndex)
				if err != nil {
					return m, models.ShowErrorCommand("Failed to delete body field: " + err.Error())
				}

				m.BodyFields = newBodyFields
				if m.pointer >= len(m.BodyFields) && m.pointer > 0 {
					m.pointer--
				}
			}
		case "e":
			m.editing = true
			value := m.BodyFields[m.pointer].Value
			m.editingBodyFields = textinput.New()
			m.editingBodyFields.SetValue(value)
			m.editingBodyFields.Focus()

		case "x":
			if m.hasError {
				m.hasError = false
				m.errorMessage = ""
				return m, nil
			}
		}
	}

	m.jsonInput, cmd = m.jsonInput.Update(msg)
	return m, cmd
}

func UpdateHeadersPage(m Model, msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:

		if m.editingHeader.Focused() {
			switch msg.String() {
			case "esc":
				m.editing = false
				m.editingHeader.Blur()
				return m, nil
			case "enter":
				m.Headers[m.pointer].Value = m.editingHeader.Value()
				if err := operations.AddHeader(m.Headers, m.storage, m.collectionIndex, m.ApiIndex); err != nil {
					return m, models.ShowErrorCommand("Failed to add new header: " + err.Error())
				}
				m.editing = false
				m.editingHeader.Blur()
			}
			m.editingHeader, cmd = m.editingHeader.Update(msg)
			return m, cmd
		}

		if m.addHeaderKey.Focused() {
			switch msg.String() {
			case "esc":
				m.addHeaderKey.SetValue("")
				m.addHeaderKey.Blur()
				return m, nil
			case "enter":
				headerKey := m.addHeaderKey.Value()
				newHeder := models.Header{
					Key: headerKey,
				}
				m.Headers = append(m.Headers, newHeder)
				if err := operations.AddHeader(m.Headers, m.storage, m.collectionIndex, m.ApiIndex); err != nil {
					return m, models.ShowErrorCommand("Failed to add Header: " + err.Error())
				}
				m.addHeaderKey.SetValue("")
				m.addHeaderKey.Blur()
			}
			m.addHeaderKey, cmd = m.addHeaderKey.Update(msg)
			return m, cmd
		}
		if m.addHeaderValue.Focused() {
			switch msg.String() {
			case "esc":
				m.addHeaderValue.SetValue("")
				m.addHeaderValue.Blur()
				return m, nil
			case "enter":
				m.Headers[m.pointer].Value = m.addHeaderValue.Value()
				if err := operations.AddHeader(m.Headers, m.storage, m.collectionIndex, m.ApiIndex); err != nil {
					return m, models.ShowErrorCommand("Failed to add header value: " + err.Error())
				}
				m.addHeaderValue.SetValue("")
				m.addHeaderValue.Blur()
			}
			m.addHeaderValue, cmd = m.addHeaderValue.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "esc":
			m.CurrentPage = CollectionPage
			m.pointer = m.ApiIndex
			m.pageScrollOffset = 0

			m.SelectedApi = m.Apis[m.ApiIndex]
			m.Headers = m.SelectedApi.Headers

		case ":":
			m.addHeaderKey.Focus()

		case "enter":
			m.addHeaderValue.Focus()

		case "d":
			if len(m.Headers) > 0 {
				selectedHeader := m.Headers[m.pointer]
				newHeaders, err := operations.DeleteHeader(selectedHeader, m.storage, m.collectionIndex, m.ApiIndex)
				if err != nil {
					return m, models.ShowErrorCommand("Failed to delete header: " + err.Error())
				}
				m.Headers = newHeaders
				if m.pointer >= len(m.Headers) && m.pointer > 0 {
					m.pointer--
				}
			}

		case "up", "k":
			if m.pointer > 0 {
				m.pointer--

				if m.pointer < m.pageScrollOffset {
					m.pageScrollOffset--
				}
			}

		case "down", "j":
			if m.pointer < len(m.Headers)-1 {
				m.pointer++

				maxVisible := 12
				if m.pointer >= m.pageScrollOffset+maxVisible {
					m.pageScrollOffset++
				}
			}

		case "e":
			m.editing = true
			value := m.Headers[m.pointer].Value
			m.editingHeader = textinput.New()
			m.editingHeader.SetValue(value)
			m.editingHeader.Focus()

		case "x":
			if m.hasError {
				m.hasError = false
				m.errorMessage = ""
				return m, nil
			}
		}
	}
	return m, nil
}

func UpdateQueryParamsPage(m Model, msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:

		if m.editingQueryParams.Focused() {
			switch msg.String() {
			case "esc":
				m.editing = false
				m.editingQueryParams.Blur()
				return m, nil
			case "enter":
				m.QueryParams[m.pointer].Value = m.editingQueryParams.Value()
				if err := operations.AddQueryParam(m.QueryParams, m.storage, m.collectionIndex, m.ApiIndex); err != nil {
					return m, models.ShowErrorCommand("Failed to edit query params: " + err.Error())
				}
				m.editing = false
				m.editingQueryParams.Blur()
			}
			m.editingQueryParams, cmd = m.editingQueryParams.Update(msg)
			return m, cmd
		}

		if m.addQueryParamsKey.Focused() {
			switch msg.String() {
			case "esc":
				m.addQueryParamsKey.SetValue("")
				m.addQueryParamsKey.Blur()
			case "enter":
				key := m.addQueryParamsKey.Value()
				newQueryParam := models.QueryParam{
					Key:   key,
					Value: "",
				}
				m.QueryParams = append(m.QueryParams, newQueryParam)
				if err := operations.AddQueryParam(m.QueryParams, m.storage, m.collectionIndex, m.ApiIndex); err != nil {
					return m, models.ShowErrorCommand("Failed to add query param: " + err.Error())
				}
				m.addQueryParamsKey.SetValue("")
				m.addQueryParamsKey.Blur()
			}
			m.addQueryParamsKey, cmd = m.addQueryParamsKey.Update(msg)
			return m, cmd
		}

		if m.addQueryParamsValue.Focused() {
			switch msg.String() {
			case "esc":
				m.addQueryParamsValue.SetValue("")
				m.addQueryParamsValue.Blur()
			case "enter":
				m.QueryParams[m.pointer].Value = m.addQueryParamsValue.Value()
				if err := operations.AddQueryParam(m.QueryParams, m.storage, m.collectionIndex, m.ApiIndex); err != nil {
					return m, models.ShowErrorCommand("Failed to add query param value: " + err.Error())
				}
				m.addQueryParamsValue.SetValue("")
				m.addQueryParamsValue.Blur()
			}
			m.addQueryParamsValue, cmd = m.addQueryParamsValue.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "esc":
			m.CurrentPage = CollectionPage
			m.pointer = m.ApiIndex
			m.pageScrollOffset = 0
		case ":":
			m.addQueryParamsKey.Focus()
		case "enter":
			m.addQueryParamsValue.Focus()
		case "up", "k":
			if m.pointer > 0 {
				m.pointer--

				if m.pointer < m.pageScrollOffset {
					m.pageScrollOffset--
				}
			}
		case "down", "j":
			if m.pointer < len(m.QueryParams)-1 {
				m.pointer++

				maxVisible := 12
				if m.pointer >= m.pageScrollOffset+maxVisible {
					m.pageScrollOffset++
				}
			}
		case "e":
			m.editing = true
			value := m.QueryParams[m.pointer].Value
			m.editingQueryParams = textinput.New()
			m.editingQueryParams.SetValue(value)
			m.editingQueryParams.Focus()
		case "d":
			if len(m.QueryParams) > 0 {
				selectedQueryParam := m.QueryParams[m.pointer]
				newQueryParams, err := operations.DeleteQueryParam(selectedQueryParam, m.storage, m.collectionIndex, m.ApiIndex)
				if err != nil {
					return m, models.ShowErrorCommand("Failed to delete query param: " + err.Error())
				}
				m.QueryParams = newQueryParams
				if m.pointer >= len(m.QueryParams) && m.pointer > 0 {
					m.pointer--
				}
			}

		case "x":
			if m.hasError {
				m.hasError = false
				m.errorMessage = ""
				return m, nil
			}
		}
	}
	return m, nil
}

func UpdateLoadingPage(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.CurrentPage = CollectionPage
		}
	}
	return m, nil
}

func UpdateResponsePage(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		if m.VariablesFocus {
			switch msg.String() {
			case "esc":
				m.CurrentPage = ApiPage
				m.pointer = m.ApiIndex
			case "r":
				m.VariablesFocus = false
				m.pointer = 0
				m.responseScrollOffset = 0
			case "up", "k":
				if m.pointer > 0 {
					m.pointer--

					if m.pointer < m.variableScrollOffset {
						m.variableScrollOffset--
					}
				}
			case "down", "j":
				if m.pointer < len(m.LocalVariables)-1 {
					m.pointer++

					maxVisible := 5
					if m.pointer >= m.variableScrollOffset+maxVisible {
						m.variableScrollOffset++
					}
				}
			case "d":
				if len(m.LocalVariables) > 0 {
					selectedVariable := m.LocalVariables[m.pointer]
					newLocalVariables, err := operations.DeleteLocalVariable(selectedVariable, m.storage, m.collectionIndex)
					if err != nil {
						return m, models.ShowErrorCommand("Failed to delete Local Variable : " + err.Error())
					}
					m.LocalVariables = newLocalVariables
					if m.pointer >= len(m.LocalVariables) && m.pointer > 0 {
						m.pointer--
					}
				}
			}
		}

		if !m.VariablesFocus {
			switch msg.String() {
			case "esc":
				m.CurrentPage = ApiPage
				m.pointer = m.ApiIndex

			case "up", "k":
				if m.pointer > 0 {
					m.pointer--

					if m.pointer < m.responseScrollOffset {
						m.responseScrollOffset--
					}
				}
			case "down", "j":
				if m.pointer < len(m.Responses)-1 {
					m.pointer++

					maxVisible := 5
					if m.pointer >= m.responseScrollOffset+maxVisible {
						m.responseScrollOffset++
					}
				}
			case "enter":
				selectedResponse := m.Responses[m.pointer]
				newLocalVariable := models.LocalVariable{
					Key:   selectedResponse.Key,
					Value: selectedResponse.Value,
				}
				m.LocalVariables = append(m.LocalVariables, newLocalVariable)
				err := operations.AddLocalVariable(m.storage, m.collectionIndex, m.LocalVariables)
				if err != nil {
					return m, models.ShowErrorCommand("Failed to add local Variable: " + err.Error())
				}
			case "v":
				m.VariablesFocus = true
				m.pointer = 0
				m.variableScrollOffset = 0
			case "c":
				selectedResponse := m.Responses[m.pointer]
				err := clipboard.WriteAll(selectedResponse.Value)
				if err != nil {
					return m, models.ShowErrorCommand("Failed to Copy Response: " + err.Error())
				}
			}
		}
	}
	return m, nil
}

func UpdateVariablesPage(m Model, msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.editingLocalVariables.Focused() {
			switch msg.String() {
			case "esc":
				m.editing = false
				m.editingLocalVariables.Blur()
				return m, nil
			case "enter":
				m.LocalVariables[m.pointer].Value = m.editingLocalVariables.Value()
				if err := operations.AddLocalVariable(m.storage, m.collectionIndex, m.LocalVariables); err != nil {
					return m, models.ShowErrorCommand("Failed to edit Local Variable : " + err.Error())
				}
				m.editing = false
				m.editingLocalVariables.Blur()
			}
			m.editingLocalVariables, cmd = m.editingLocalVariables.Update(msg)
			return m, cmd
		}
		if m.addVariableValue.Focused() {
			switch msg.String() {
			case "esc":
				m.addVariableValue.Blur()
				m.addVariableValue.SetValue("")
			case "enter":
				m.LocalVariables[m.pointer].Value = m.addVariableValue.Value()
				err := operations.AddLocalVariable(m.storage, m.collectionIndex, m.LocalVariables)
				if err != nil {
					return m, models.ShowErrorCommand("Failed to add local Variable: " + err.Error())
				}
				m.addVariableValue.Blur()
				m.addVariableValue.SetValue("")
			}
			m.addVariableValue, cmd = m.addVariableValue.Update(msg)
			return m, cmd
		}
		if m.addVariableKey.Focused() {
			switch msg.String() {
			case "esc":
				m.addVariableKey.Blur()
				m.addVariableKey.SetValue("")
			case "enter":
				NewResponse := models.LocalVariable{
					Key:   m.addVariableKey.Value(),
					Value: "",
				}
				m.LocalVariables = append(m.LocalVariables, NewResponse)
				err := operations.AddLocalVariable(m.storage, m.collectionIndex, m.LocalVariables)
				if err != nil {
					return m, models.ShowErrorCommand("Failed to add local Variable: " + err.Error())
				}
				m.addVariableKey.Blur()
				m.addVariableKey.SetValue("")
			}
			m.addVariableKey, cmd = m.addVariableKey.Update(msg)
			return m, cmd
		}
		switch msg.String() {
		case "esc":
			m.CurrentPage = CollectionPage
			m.pageScrollOffset = 0
		case "up", "k":
			if m.pointer > 0 {
				m.pointer--

				if m.pointer < m.pageScrollOffset {
					m.pageScrollOffset--
				}
			}
		case "down", "j":
			if m.pointer < len(m.LocalVariables)-1 {
				m.pointer++

				maxVisible := 12
				if m.pointer >= m.pageScrollOffset+maxVisible {
					m.pageScrollOffset++
				}
			}
		case "d":
			if len(m.LocalVariables) > 0 {
				selectedVariable := m.LocalVariables[m.pointer]
				newLocalVariables, err := operations.DeleteLocalVariable(selectedVariable, m.storage, m.collectionIndex)
				if err != nil {
					return m, models.ShowErrorCommand("Failed to delete Local Variable : " + err.Error())
				}
				m.LocalVariables = newLocalVariables
				if m.pointer >= len(m.LocalVariables) && m.pointer > 0 {
					m.pointer--
				}
			}
		case ":":
			m.addVariableKey.Focus()
		case "enter":
			m.addVariableValue.Focus()
		case "e":
			m.editing = true
			value := m.LocalVariables[m.pointer].Value
			m.editingLocalVariables = textinput.New()
			m.editingLocalVariables.SetValue(value)
			m.editingLocalVariables.Focus()
		}
	}
	return m, nil
}

func UpdateDashBoard(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.responseComponent {
			switch msg.String() {
			case "esc":
				m.responseComponent = false
			}
			return m, nil
		}
		if m.reqPageComponent {
			switch msg.String() {
			// need to fix all the update functions so i can use them in the dashboard
			// propably the whole app will change , this is a test so idc

			case "enter":
				TuiModel := models.TuiModel{
					SelectedApi:        m.SelectedApi,
					SelectedCollection: m.SelectedCollection,
					LocalVariables:     m.LocalVariables,
				}

				m.pageScrollOffset = 0
				m.ApiResponse = operations.PostAPiFunc(TuiModel)
				m.Responses, _ = operations.HandleJson(m.ApiResponse)
				m.reqPageComponent = false
				m.responseComponent = true
			}
			return UpdateReqPage(m, msg)
		}

		switch msg.String() {
		case "esc":
			m.CurrentPage = HomePage
		case "up", "k":
			if m.pointer > 0 {
				m.pointer--
			}
		case "down", "j":
			if m.pointer < len(m.SelectedCollection.Requests)-1 {
				m.pointer++
			}
		case "enter":
			m.SelectedApi = m.SelectedCollection.Requests[m.pointer]

			processedApi := operations.ProcessRequest(m.SelectedApi, m.SelectedCollection.LocalVariables)

			switch processedApi.Method {
			case "POST", "DELETE", "PUT", "PATCH":
				m.SelectedApi = processedApi
				m.BodyFields = processedApi.BodyField
				m.ApiIndex = m.pointer
				m.reqPageComponent = true

			case "GET":
				TuiModel := models.TuiModel{
					SelectedApi:        m.SelectedApi,
					SelectedCollection: m.SelectedCollection,
					LocalVariables:     m.LocalVariables,
				}
				m.responseComponent = true
				m.ApiIndex = m.pointer
				m.ApiResponse = operations.FetchData(m.SelectedApi, TuiModel)
				m.Responses, _ = operations.HandleJson(m.ApiResponse)
				return m, nil
			}
		}
	}
	return m, nil
}
