package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	switch m.CurrentPage {
	case HomePage:
		return Homepage(m.termWidth, m.termHeight, m)
	case ApiPage:
		return Apipage(m, m.termWidth)
	case RequestPage:
		return ReqPage(m)
	}
	return ""
}

func Homepage(termWidth, termHeight int, m model) string {
	style1 := HomePageStyle1(termWidth)
	style2 := HomePageStyle2(termWidth, termHeight)
	style3 := OptionsStyle(termWidth)

	var b strings.Builder

	b.WriteString(style1.Render(" This is the HomePage !"))
	b.WriteString("\n\n")

	// Build options with pointer
	var items []string
	for i := 0; i < len(m.Options); i++ {
		api := m.Options[i]

		text := api.Method + " " + api.Url
		if i == m.pointer { // highlight selected option
			text = style4.Render("> ") + style5.Render(text+"\n")
		} else {
			text = "   " + text + "\n"
		}
		items = append(items, text)
	}

	// Render left box with styled options
	leftBox := style3.Render(lipgloss.JoinVertical(lipgloss.Left, items...))

	// Commands box
	rightBox := style2.Render("Commands\n----------------\nESC -> Quit\n\nk -> Up\n\nj -> Down\n\nEnter -> Open")

	// Combine side by side
	layout := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)

	b.WriteString(layout)
	return b.String()
}

func Apipage(m model, termWidth int) string {

	style1 := HomePageStyle1(termWidth)
	style3 := ResponseStyle(termWidth)

	var b strings.Builder
	SelectedApi := m.Options[m.pointer]

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

	resp.WriteString("\nBody:\n" + Response.Body + "\n")

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
