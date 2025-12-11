package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	switch m.CurrentPage {
	case HomePage:
		return Homepage(m)
	case CollectionPage:
		return Collectionpage(m.termWidth, m.termHeight, m)
	case ApiPage:
		return ApipageWithViewport(m)
	case RequestPage:
		return ReqPage(m)
	}
	return ""
}

func Homepage(m model) string {
	style3 := OptionsStyle(m.termWidth)
	style2 := HomePageStyle2(m.termWidth, m.termHeight)

	var b strings.Builder

	collections := m.storage.Collections

	var items []string

	for i := 0; i < len(collections); i++ {
		text := collections[i].Name
		if i == m.pointer {
			text = style4.Render("> ") + style5.Render(text+"\n")
		} else {
			text = "   " + text + "\n"
		}
		items = append(items, text)
	}

	leftBox := style3.Render(lipgloss.JoinVertical(lipgloss.Left, items...)) + "\n\n" + m.NewCollectionInput.View()
	rightBox := style2.Render("Commands\n----------------\nESC -> Quit\n\nk -> Up\n\nj -> Down\n\nEnter -> Open\n\n: -> Add New")
	layout := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)

	b.WriteString(layout)
	return b.String()
}

func Collectionpage(termWidth, termHeight int, m model) string {
	style1 := HomePageStyle1(termWidth)
	style2 := HomePageStyle2(termWidth, termHeight)
	style3 := OptionsStyle(termWidth)

	var b strings.Builder
	collectionName := m.SelectedCollection.Name

	b.WriteString(style1.Render(collectionName))
	b.WriteString("\n\n")

	var items []string

	for i := 0; i < len(m.Apis); i++ {
		api := m.Apis[i]

		if i == m.pointer && m.editing {
			line := style4.Render("> ") + m.editingApi.View() + "\n"
			items = append(items, line)
			continue
		}

		text := api.Method + " " + api.Url
		if i == m.pointer {
			text = style4.Render("> ") + style5.Render(text+"\n")
		} else {
			text = "   " + text + "\n"
		}
		items = append(items, text)
	}

	leftBox := style3.Render(lipgloss.JoinVertical(lipgloss.Left, items...)) + "\n\n" + m.NewApiInput.View()
	rightBox := style2.Render("Commands\n----------------\nESC -> Quit\n\nk -> Up\n\nj -> Down\n\nEnter -> Open\n\n: -> Add New")
	layout := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)

	b.WriteString(layout)
	return b.String()
}

func ApipageWithViewport(m model) string {
	if !m.viewportReady {
		return "Loading..."
	}

	helpText := HelpTextStyle.Render("\n\n↑/↓ j/k: scroll • space/b: page up/down • g/G: top/bottom • esc: back")
	return m.apiViewport.View() + helpText
}

func BuildApiPageContent(m model, termWidth int) string {
	style1 := HomePageStyle1(termWidth)
	style3 := ResponseStyle(termWidth)

	var b strings.Builder
	SelectedApi := m.Apis[m.pointer]

	var Response ApiResponse
	if m.SelectedApi.Method == "POST" {
		Response = PostAPiFunc(m)
	}
	if m.SelectedApi.Method == "GET" {
		Response = FetchData(SelectedApi)
	}

	statusStyle := StatusOKStyle
	if Response.StatusCode >= 400 {
		statusStyle = StatusErrorStyle
	}

	var resp strings.Builder

	resp.WriteString("Status: " + statusStyle.Render(Response.Status) + "\n")
	resp.WriteString(fmt.Sprintf("Status Code: %s\n", statusStyle.Render(fmt.Sprintf("%d", Response.StatusCode))))
	resp.WriteString("Content Type: " + Response.ContentType + "\n")
	resp.WriteString(fmt.Sprintf("Content Length: %d\n", Response.ContentLength))
	resp.WriteString("\nHeaders:\n")

	for k, v := range Response.Headers {
		resp.WriteString(fmt.Sprintf("  %s: %s\n", k, strings.Join(v, ", ")))
	}

	formattedBody := FormatJSON(Response.Body)
	resp.WriteString("\nBody:\n" + formattedBody + "\n")

	b.WriteString(style1.Render("This is the Api-Page !"))
	b.WriteString("\n\n")

	b.WriteString(style3.Render(
		"Selected Api is : " +
			MethodStyle.Render(SelectedApi.Method) + " " +
			UrlStyle.Render(SelectedApi.Url),
	))
	b.WriteString("\n\n")

	b.WriteString(style3.Render("Response:\n\n" + resp.String()))

	return b.String()
}

func ReqPage(m model) string {
	var b strings.Builder
	style1 := InputStyle(m.termWidth)
	b.WriteString("\n")

	b.WriteString("Edit JSON body below:\n\n")
	b.WriteString(style1.Render(m.jsonInput.View()))
	b.WriteString("\n\n")
	b.WriteString("Press Enter to send POST request...\n")

	return b.String()
}

func FormatJSON(body string) string {
	var formatted strings.Builder
	indent := 0

	for i := 0; i < len(body); i++ {
		char := body[i]

		switch char {
		case '{', '[':
			formatted.WriteByte(char)
			formatted.WriteByte('\n')
			indent++
			formatted.WriteString(strings.Repeat("  ", indent))

		case '}', ']':
			formatted.WriteByte('\n')
			indent--
			formatted.WriteString(strings.Repeat("  ", indent))
			formatted.WriteByte(char)

		case ',':
			formatted.WriteByte(char)
			formatted.WriteByte('\n')
			formatted.WriteString(strings.Repeat("  ", indent))

		case ':':
			formatted.WriteString(": ")

		case ' ', '\t', '\n', '\r':
			continue

		default:
			formatted.WriteByte(char)
		}
	}

	return formatted.String()
}
