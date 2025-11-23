package main

import tea "github.com/charmbracelet/bubbletea"

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height

	case tea.KeyMsg:
		switch m.CurrentPage {
		case HomePage:
			m, cmd := UpdateHomePage(m, msg)
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
		case "up", "k":
			if m.pointer > 0 {
				m.pointer--
			}
		case "down", "j":
			if m.pointer < len(m.Options)-1 {
				m.pointer++
			}
		case "enter":
			m.SelectedApi = m.Options[m.pointer]

			if m.SelectedApi.Method == "POST" {
				m.CurrentPage = RequestPage
			}

			if m.SelectedApi.Method == "GET" {
				m.CurrentPage = ApiPage
			}
		case "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func UpdateApiPage(m model, msg tea.Msg) (model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.CurrentPage = HomePage
		}
	}

	return m, nil
}

func UpdateReqPage(m model, msg tea.Msg) (model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.CurrentPage = ApiPage
		case "esc":
			m.CurrentPage = HomePage
		}
	}

	m.jsonInput, cmd = m.jsonInput.Update(msg)
	return m, cmd

}
