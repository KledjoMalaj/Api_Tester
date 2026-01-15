package main

import (
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case errorMsg:
		m.errorMessage = msg.message
		m.hasError = true
		return m, nil

	case apiResponseMsg:
		m.apiResponse = msg.response
		m.CurrentPage = ApiPage
		if m.viewportReady {
			m.apiViewport.SetContent(BuildApiPageContent(m, m.termWidth))
			m.apiViewport.GotoTop()
		}
		return m, nil

	case fileChangedMsg:
		m.storage = Storage(msg)
		m.Collections = m.storage.Collections
		m.LocalVariables = m.storage.Collections[m.collectionIndex].LocalVariables

		if m.CurrentPage == CollectionPage || m.CurrentPage == HeadersPage ||
			m.CurrentPage == RequestPage || m.CurrentPage == QueryParamsPage {

			if m.collectionIndex >= 0 && m.collectionIndex < len(m.Collections) {
				m.SelectedCollection = m.Collections[m.collectionIndex]
				m.Apis = m.SelectedCollection.Requests

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
		case VariablesPage:
			m, cmd := UpdateVariablesPage(m, msg)
			return m, cmd
		}
	}

	return m, nil
}

func UpdateHomePage(m model, msg tea.Msg) (model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:

		if m.editing {
			switch msg.String() {

			case "esc":
				m.editingCollection.Blur()
				m.editing = false
			case "enter":
				if err := editCollection(m.storage, m.SelectedCollection, m.editingCollection.Value()); err != nil {
					return m, showErrorCommand("Failed to edit Collection: " + err.Error())
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

				if err := AddCollection(m.storage, m.Collections, m.NewCollectionInput.Value()); err != nil {
					return m, showErrorCommand("Failed to add collection: " + err.Error())
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
			}
		case "down", "j":
			if m.pointer < len(m.storage.Collections)-1 {
				m.pointer++
			}
		case "enter":
			m.CurrentPage = CollectionPage
			m.SelectedCollection = m.storage.Collections[m.pointer]
			m.Apis = m.SelectedCollection.Requests
			m.collectionIndex = m.pointer
			m.pointer = 0

		case ":":
			m.NewCollectionInput.Focus()
			return m, nil

		case "d":
			if len(m.Collections) > 0 {
				selectedCollection := m.storage.Collections[m.pointer]
				newCollections, err := deleteCollection(selectedCollection, m.storage)
				if err != nil {
					return m, showErrorCommand("Failed to delete collection: " + err.Error())
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

func UpdateCollectionPage(m model, msg tea.Msg) (model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.editing {
			switch msg.String() {
			case "enter":
				if err := editApi(m.storage, m.collectionIndex, m.SelectedApi, m.editingApi.Value()); err != nil {
					return m, showErrorCommand("Failed to edit api: " + err.Error())
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

				if err := AddApi(m.storage, m.collectionIndex, m.Apis, m.NewApiInput.Value()); err != nil {
					return m, showErrorCommand("Failed to add api: " + err.Error())
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
			}
		case "down", "j":
			if m.pointer < len(m.Apis)-1 {
				m.pointer++
			}
		case "enter":
			m.SelectedApi = m.Apis[m.pointer]

			switch m.SelectedApi.Method {
			case "POST", "DELETE", "PUT", "PATCH":
				m.SelectedApi = m.Apis[m.pointer]
				m.BodyFields = m.SelectedApi.BodyField
				m.ApiIndex = m.pointer
				m.CurrentPage = RequestPage
				m.pointer = 0

			case "GET":
				m.CurrentPage = LoadingPage
				m.ApiIndex = m.pointer
				m.apiResponse = FetchData(m.SelectedApi, m)
				m.Responses, _ = HandleJson(m.apiResponse)
				return m, fetchApiCommand(m.SelectedApi, m)
			}

		case ":":
			m.NewApiInput.Focus()
			return m, nil

		case "d":
			if len(m.Apis) > 0 {
				selectedApi := m.Apis[m.pointer]
				newApis, err := deleteApi(selectedApi, m.storage, m.collectionIndex)
				if err != nil {
					return m, showErrorCommand("Failed to delete api: " + err.Error())
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
		}
	}

	m.NewApiInput, cmd = m.NewApiInput.Update(msg)
	return m, cmd
}

func UpdateApiPage(m model, msg tea.Msg) (model, tea.Cmd) {
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
				if err := editApi(m.storage, m.collectionIndex, m.SelectedApi, m.editingCurrentApi.Value()); err != nil {
					return m, showErrorCommand("Failed to edit api: " + err.Error())
				}

				// Update local state
				m.storage, _ = ReadFile()
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
			m.CurrentPage = VariablesPage
			m.pointer = 0
		}
	}

	m.apiViewport, cmd = m.apiViewport.Update(msg)
	return m, cmd
}

func UpdateReqPage(m model, msg tea.Msg) (model, tea.Cmd) {
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
				newBodyFields, err := addBodyField(m.storage, m.collectionIndex, m.ApiIndex, m.BodyFields)
				if err != nil {
					return m, showErrorCommand("Failed to edit body field: " + err.Error())
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
				newBodyFiled := BodyField{
					Key:   newBodyFieldKey,
					Value: "",
				}
				m.BodyFields = append(m.BodyFields, newBodyFiled)
				newBodyFields, err := addBodyField(m.storage, m.collectionIndex, m.ApiIndex, m.BodyFields)
				if err != nil {
					return m, showErrorCommand("Failed to add body field: " + err.Error())
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
				_, err := addBodyField(m.storage, m.collectionIndex, m.ApiIndex, m.BodyFields)
				if err != nil {
					return m, showErrorCommand("Failed to add body field value: " + err.Error())
				}
				m.bodyFiledValueInput.SetValue("")
				m.bodyFiledValueInput.Blur()
			}
			m.bodyFiledValueInput, cmd = m.bodyFiledValueInput.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "enter":
			m.CurrentPage = LoadingPage
			m.apiResponse = PostAPiFunc(m)
			m.Responses, _ = HandleJson(m.apiResponse)
			return m, postApiCommand(m)

		case "v":
			m.bodyFiledValueInput.Focus()
		case ":":
			m.newBodyFieldInput.Focus()
		case "esc":
			m.CurrentPage = CollectionPage
		case "up", "k":
			if m.pointer > 0 {
				m.pointer--
			}
		case "down", "j":
			if m.pointer < len(m.BodyFields)-1 {
				m.pointer++
			}
		case "d":
			if len(m.BodyFields) > 0 {
				selectedBodyField := m.BodyFields[m.pointer]
				newBodyFields, err := deleteBodyField(selectedBodyField, m.storage, m.collectionIndex, m.ApiIndex)
				if err != nil {
					return m, showErrorCommand("Failed to delete body field: " + err.Error())
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

func UpdateHeadersPage(m model, msg tea.Msg) (model, tea.Cmd) {
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
				if err := addHeader(m.Headers, m.storage, m.collectionIndex, m.ApiIndex); err != nil {
					return m, showErrorCommand("Failed to add new header: " + err.Error())
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
				newHeder := Header{
					Key: headerKey,
				}
				m.Headers = append(m.Headers, newHeder)
				if err := addHeader(m.Headers, m.storage, m.collectionIndex, m.ApiIndex); err != nil {
					return m, showErrorCommand("Failed to add Header: " + err.Error())
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
				if err := addHeader(m.Headers, m.storage, m.collectionIndex, m.ApiIndex); err != nil {
					return m, showErrorCommand("Failed to add header value: " + err.Error())
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

			m.SelectedApi = m.Apis[m.ApiIndex]
			m.Headers = m.SelectedApi.Headers

		case ":":
			m.addHeaderKey.Focus()

		case "enter":
			m.addHeaderValue.Focus()

		case "d":
			if len(m.Headers) > 0 {
				selectedHeader := m.Headers[m.pointer]
				newHeaders, err := deleteHeader(selectedHeader, m.storage, m.collectionIndex, m.ApiIndex)
				if err != nil {
					return m, showErrorCommand("Failed to delete header: " + err.Error())
				}
				m.Headers = newHeaders
				if m.pointer >= len(m.Headers) && m.pointer > 0 {
					m.pointer--
				}
			}

		case "up", "k":
			if m.pointer > 0 {
				m.pointer--
			}

		case "down", "j":
			if m.pointer < len(m.Headers)-1 {
				m.pointer++
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

func UpdateQueryParamsPage(m model, msg tea.Msg) (model, tea.Cmd) {
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
				if err := addQueryParam(m.QueryParams, m.storage, m.collectionIndex, m.ApiIndex); err != nil {
					return m, showErrorCommand("Failed to edit query params: " + err.Error())
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
				newQueryParam := QueryParam{
					Key:   key,
					Value: "",
				}
				m.QueryParams = append(m.QueryParams, newQueryParam)
				if err := addQueryParam(m.QueryParams, m.storage, m.collectionIndex, m.ApiIndex); err != nil {
					return m, showErrorCommand("Failed to add query param: " + err.Error())
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
				if err := addQueryParam(m.QueryParams, m.storage, m.collectionIndex, m.ApiIndex); err != nil {
					return m, showErrorCommand("Failed to add query param value: " + err.Error())
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
		case ":":
			m.addQueryParamsKey.Focus()
		case "enter":
			m.addQueryParamsValue.Focus()
		case "up", "k":
			if m.pointer > 0 {
				m.pointer--
			}
		case "down", "j":
			if m.pointer < len(m.QueryParams)-1 {
				m.pointer++
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
				newQueryParams, err := deleteQueryParam(selectedQueryParam, m.storage, m.collectionIndex, m.ApiIndex)
				if err != nil {
					return m, showErrorCommand("Failed to delete query param: " + err.Error())
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

func UpdateLoadingPage(m model, msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.CurrentPage = CollectionPage
		}
	}
	return m, nil
}

func UpdateVariablesPage(m model, msg tea.Msg) (model, tea.Cmd) {
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
			case "up", "k":
				if m.pointer > 0 {
					m.pointer--
				}
			case "down", "j":
				if m.pointer < len(m.LocalVariables)-1 {
					m.pointer++
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
				}
			case "down", "j":
				if m.pointer < len(m.Responses)-1 {
					m.pointer++
				}
			case "enter":
				selectedResponse := m.Responses[m.pointer]
				err := addLocalVariable(m.storage, m.collectionIndex, selectedResponse, m.LocalVariables)
				if err != nil {
					return m, showErrorCommand("Failed to add local Variable: " + err.Error())
				}
			case "v":
				m.VariablesFocus = true
				m.pointer = 0
			case "c":
				selectedResponse := m.Responses[m.pointer]
				err := clipboard.WriteAll(selectedResponse.Value)
				if err != nil {
					return m, showErrorCommand("Failed to Copy Response: " + err.Error())
				}
			}
		}
	}
	return m, nil
}
