package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case fileChangedMsg:
		m.Apis = msg

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
		}
	}

	return m, nil
}

func UpdateHomePage(m model, msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
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
			m.Apis = m.storage.Collections[m.pointer].Requests
			m.SelectedCollection = m.storage.Collections[m.pointer]
			m.collectionIndex = m.pointer
			m.pointer = 0
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
				parts := strings.SplitN(m.NewApiInput.Value(), " ", 2)
				newApi := Api{
					Method: parts[0],
					Url:    parts[1],
				}
				m.Apis = append(m.Apis, newApi)
				AddApi(m.storage, m.collectionIndex, m.Apis)
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
				m.CurrentPage = RequestPage

			case "GET":
				m.CurrentPage = ApiPage
				if m.viewportReady {
					m.apiViewport.SetContent(BuildApiPageContent(m, m.termWidth))
					m.apiViewport.GotoTop()
				}
			}

		case ":":
			m.NewApiInput.Focus()

		case "d":
			if len(m.Apis) > 0 {
				m.SelectedApi = m.Apis[m.pointer]
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
			m.pointer = 0
		}
	}

	m.NewApiInput, cmd = m.NewApiInput.Update(msg)
	return m, cmd
}

func UpdateApiPage(m model, msg tea.Msg) (model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.CurrentPage = CollectionPage
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
		}
	}

	m.apiViewport, cmd = m.apiViewport.Update(msg)
	return m, cmd
}

func UpdateReqPage(m model, msg tea.Msg) (model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.CurrentPage = ApiPage

			// Load content into viewport when entering ApiPage from RequestPage
			if m.viewportReady {
				m.apiViewport.SetContent(BuildApiPageContent(m, m.termWidth))
				m.apiViewport.GotoTop()
			}
		case "esc":
			m.CurrentPage = CollectionPage
		}
	}

	m.jsonInput, cmd = m.jsonInput.Update(msg)
	return m, cmd
}
