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
	case HeadersPage:
		return HeadersPageView(m)
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

		if i == m.pointer && m.editing {
			line := style4.Render("> ") + m.editingCollection.View() + "\n"
			items = append(items, line)
			continue
		}

		if i == m.pointer {
			text = style4.Render("> ") + style5.Render(text+"\n")
		} else {
			text = "   " + text + "\n"
		}
		items = append(items, text)
	}

	leftBox := style3.Render(lipgloss.JoinVertical(lipgloss.Left, items...)) + "\n\n" + m.NewCollectionInput.View()
	rightBox := style2.Render("Commands\n----------------\nESC -> Quit\n\nk -> Up\n\nj -> Down\n\nEnter -> Open\n\n: -> Add New\n\nd -> Delete\n\ne -> Edit")
	layout := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)

	b.WriteString(layout)
	return b.String()
}

func Collectionpage(termWidth, termHeight int, m model) string {
	style1 := HomePageStyle1(termWidth + 21)
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
	rightBox := style2.Render("Commands\n----------------\nESC -> Quit\n\nk -> Up\n\nj -> Down\n\nEnter -> Open\n\n: -> Add New\n\nd -> Delete\n\ne -> Edit")
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
	style1 := TitleStyle(termWidth)
	style3 := ResponseStyle(termWidth)
	style2 := statusLine(termWidth)

	var b strings.Builder
	SelectedApi := m.SelectedApi

	var Response ApiResponse

	switch m.SelectedApi.Method {
	case "POST", "DELETE", "PUT", "PATCH":
		Response = PostAPiFunc(m)
	case "GET":
		Response = FetchData(m.SelectedApi)
	}

	statusStyle := StatusOKStyle
	if Response.StatusCode >= 400 {
		statusStyle = StatusErrorStyle
	}

	var resp strings.Builder

	resp.WriteString("Status: " + statusStyle.Render(Response.Status) + "\n")
	resp.WriteString(fmt.Sprintf("Status Code: %s\n", statusStyle.Render(fmt.Sprintf("%d \n", Response.StatusCode))))
	resp.WriteString(style2.Render(" "))
	resp.WriteString("Content Type: " + Response.ContentType + "\n")
	resp.WriteString(fmt.Sprintf("Content Length: %d\n", Response.ContentLength))

	resp.WriteString("\nRequestHeaders :\n")
	for i := 0; i < len(Response.RequestHeaders); i++ {
		resp.WriteString(Response.RequestHeaders[i].Key + " " + Response.RequestHeaders[i].Value + "\n")
	}

	resp.WriteString("\nHeaders :\n")
	for k, v := range Response.Headers {
		resp.WriteString(fmt.Sprintf("  %s: %s\n", k, strings.Join(v, ", ")))
	}

	resp.WriteString("\n" + style2.Render(" "))

	formattedBody := FormatJSON(Response.Body)
	resp.WriteString("\nBody:\n" + formattedBody + "\n")

	b.WriteString(style1.Render("This is the Api-Page !"))

	b.WriteString("\n\n")
	if m.editing {
		b.WriteString(style3.Render("editing..." + m.editingCurrentApi.View()))
	} else {
		b.WriteString(style3.Render(
			"Selected Api is : " +
				MethodStyle.Render(SelectedApi.Method) + " " + UrlStyle.Render(SelectedApi.Url),
		))
	}

	b.WriteString("\n\n")

	b.WriteString(style3.Render("Response:\n\n" + resp.String()))

	return b.String()
}

func ReqPage(m model) string {
	style1 := TitleStyle(m.termWidth)
	style2 := OptionsStyle(m.termWidth)
	style3 := HomePageStyle2(m.termWidth, m.termHeight)

	name := m.SelectedApi.Method + "  " + m.SelectedApi.Url

	var b strings.Builder

	bodyFields := m.BodyFields

	b.WriteString(style1.Render(name))
	b.WriteString("\n")

	var items []string

	if len(bodyFields) == 0 {
		line := style4.Render("No Request Fields\n\n")
		items = append(items, line)

	} else {
		for i := 0; i < len(bodyFields); i++ {
			var line string
			if bodyFields[i].Value != "" {
				if m.pointer == i {
					line = style4.Render("> ") + style5.Render(bodyFields[i].Key+" : "+bodyFields[i].Value+"\n")
				} else {
					line = style4.Render("   ") + (bodyFields[i].Key + " : " + bodyFields[i].Value + "\n")
				}

				items = append(items, line)
			} else {
				if m.pointer == i {
					line = style4.Render("> ") + style5.Render(bodyFields[i].Key+" : "+m.bodyFiledValueInput.View()+"\n")
				} else {
					line = style4.Render("   ") + (bodyFields[i].Key + " : " + bodyFields[i].Value + "\n")
				}
				items = append(items, line)
			}
		}
	}

	leftBox := style2.Render(lipgloss.JoinVertical(lipgloss.Left, items...))
	rightBox := style3.Render("Commands\n----------------\nESC -> Quit\n\nk -> Up\n\nj -> Down\n\nEnter -> Open\n\n: -> Add New\n\nd -> Delete\n\nv -> Add Value")
	layout := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)

	b.WriteString(layout)
	b.WriteString("\n")
	b.WriteString(m.newBodyFieldInput.View())

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

func HeadersPageView(m model) string {
	style1 := TitleStyle(m.termWidth)
	style2 := OptionsStyle(m.termWidth)
	style3 := HomePageStyle2(m.termWidth, m.termHeight)

	name := m.SelectedApi.Method + "  " + m.SelectedApi.Url

	var b strings.Builder
	b.WriteString(style1.Render(name))

	headers := m.Headers
	b.WriteString("\n")

	var items []string

	if len(headers) == 0 {
		line := style4.Render("No headers\n\n")
		items = append(items, line)

	} else {
		for i, h := range headers {
			var line string
			if h.Value != "" {
				if m.pointer == i {
					line = style4.Render("> ") + style5.Render(h.Key+" "+h.Value+"\n")
				} else {
					line = style4.Render("   ") + (h.Key + " " + h.Value + "\n")
				}
				items = append(items, line)

			} else {
				if m.pointer == i {
					line = style4.Render("> ") + style5.Render(h.Key+" "+m.addHeaderValue.View()+"\n")
				} else {
					line = style4.Render("   ") + (h.Key + " " + h.Value + "\n")
				}
				items = append(items, line)
			}
		}
	}

	leftBox := style2.Render(lipgloss.JoinVertical(lipgloss.Left, items...))
	rightBox := style3.Render("Commands\n----------------\nESC -> Quit\n\nk -> Up\n\nj -> Down\n\n: -> Add New\n\nd -> Delete\n\nEnter -> Add Val ")
	layout := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)

	b.WriteString(layout)
	b.WriteString("\n")
	b.WriteString(m.addHeaderKey.View())

	return b.String()
}
