package main

import (
	"encoding/json"
	"fmt"
	"regexp"
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
	case QueryParamsPage:
		return QueryParamsPageView(m)
	case LoadingPage:
		return loadingView(m)
	case VariablesPage:
		return VariablePageView(m)
	}
	return ""
}

func Homepage(m model) string {
	styleInput := inputStyle(m.termWidth)
	style3 := OptionsStyle(m.termWidth)
	style2 := HomePageStyle2(m.termWidth, m.termHeight)
	style1 := TitleStyle(m.termWidth)

	var b strings.Builder

	b.WriteString(style1.Render("Collections "))
	b.WriteString("\n")

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

	var errorWarning string

	if m.hasError {
		errorStyle := errorStyle(m.termWidth)
		line := errorStyle.Render("⚠ ERROR: " + m.errorMessage + "\n\nPress 'x' to dismiss")
		errorWarning = line
	}

	leftBox := style3.Render(lipgloss.JoinVertical(lipgloss.Left, items...)) + "\n\n" + styleInput.Render(lipgloss.JoinVertical(lipgloss.Left, m.NewCollectionInput.View())) + "\n\n" + errorWarning
	rightBox := style2.Render("Commands\n----------------\nESC -> Quit\n\nk -> Up\n\nj -> Down\n\nEnter -> Open\n\n: -> Add New\n\nd -> Delete\n\ne -> Edit")
	layout := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)

	b.WriteString(layout)
	return b.String()
}

func Collectionpage(termWidth, termHeight int, m model) string {
	style1 := TitleStyle(termWidth)
	style2 := HomePageStyle2(termWidth, termHeight)
	style3 := OptionsStyle(termWidth)
	styleInput := inputStyle(m.termWidth)

	var b strings.Builder
	collectionName := m.SelectedCollection.Name

	b.WriteString(style1.Render(collectionName))
	b.WriteString("\n")

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

	var errorWarning string

	if m.hasError {
		errorStyle := errorStyle(m.termWidth)
		line := errorStyle.Render("⚠ ERROR: " + m.errorMessage + "\n\nPress 'x' to dismiss")
		errorWarning = line
	}

	leftBox := style3.Render(lipgloss.JoinVertical(lipgloss.Left, items...)) + "\n\n" + styleInput.Render(lipgloss.JoinVertical(lipgloss.Left, m.NewApiInput.View())) + "\n\n" + errorWarning
	rightBox := style2.Render("Commands\n----------------\nESC -> Quit\n\nk -> Up\n\nj -> Down\n\nEnter -> Open\n\n: -> Add New\n\nd -> Delete\n\ne -> Edit\n\nh -> Headers\n\nq -> QueryParams")
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
	styleInput := inputStyle(termWidth)

	var b strings.Builder
	SelectedApi := m.SelectedApi

	Response := m.apiResponse

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
		resp.WriteString(" " + Response.RequestHeaders[i].Key + " : " + Response.RequestHeaders[i].Value + "\n")
	}

	resp.WriteString("\nHeaders :\n")
	for k, v := range Response.Headers {
		resp.WriteString(fmt.Sprintf("  %s : %s\n", k, strings.Join(v, ", ")))
	}

	resp.WriteString("\n" + style2.Render(" "))

	formattedBody := FormatJSON(Response.Body, bodyElementStyle, bodyElementStyle2)
	resp.WriteString("\nBody:\n" + formattedBody + "\n")

	b.WriteString(style1.Render("This is the Api-Page !"))

	b.WriteString("\n\n")
	if m.editing {
		b.WriteString(style3.Render("editing..." + "\n" + styleInput.Render(m.editingCurrentApi.View())))
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
	styleInput := inputStyle(m.termWidth)

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

			if m.pointer == i && m.editing {
				line = style4.Render("> ") + style5.Render(bodyFields[i].Key+" : "+m.editingBodyFields.View()+"\n")
			} else if m.pointer == i && bodyFields[i].Value == "" {
				line = style4.Render("> ") + style5.Render(bodyFields[i].Key+" : "+m.bodyFiledValueInput.View()+"\n")
			} else if m.pointer == i {
				line = style4.Render("> ") + style5.Render(bodyFields[i].Key+" : "+bodyFields[i].Value+"\n")
			} else {
				line = style4.Render("   ") + (bodyFields[i].Key + " : " + bodyFields[i].Value + "\n")
			}
			items = append(items, line)
		}
	}

	var errorWarning string

	if m.hasError {
		errorStyle := errorStyle(m.termWidth)
		line := errorStyle.Render("⚠ ERROR: " + m.errorMessage + "\n\nPress 'x' to dismiss")
		errorWarning = line
	}

	leftBox := style2.Render(lipgloss.JoinVertical(lipgloss.Left, items...)) + "\n\n" + styleInput.Render(lipgloss.JoinVertical(lipgloss.Left, m.newBodyFieldInput.View())) + "\n\n" + errorWarning
	rightBox := style3.Render("Commands\n----------------\nESC -> Quit\n\nk -> Up\n\nj -> Down\n\nEnter -> Open\n\n: -> Add New\n\nd -> Delete\n\nv -> Add Value\n\ne -> edit")
	layout := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)

	b.WriteString(layout)

	return b.String()
}

func FormatJSON(body string, keyStyle lipgloss.Style, valueStyle lipgloss.Style) string {
	var jsonData interface{}
	if err := json.Unmarshal([]byte(body), &jsonData); err != nil {
		return body
	}

	formatted, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return body
	}

	formattedStr := string(formatted)

	keyRe := regexp.MustCompile(`"([^"]+)":`)
	styled := keyRe.ReplaceAllStringFunc(formattedStr, func(match string) string {
		keyMatch := keyRe.FindStringSubmatch(match)
		if len(keyMatch) > 1 {
			return keyStyle.Render(`"`+keyMatch[1]+`"`) + ":"
		}
		return match
	})

	stringValueRe := regexp.MustCompile(`:\s*"([^"]+)"`)
	styled = stringValueRe.ReplaceAllStringFunc(styled, func(match string) string {
		valueMatch := stringValueRe.FindStringSubmatch(match)
		if len(valueMatch) > 1 {
			prefix := match[:strings.Index(match, `"`)]
			return prefix + valueStyle.Render(`"`+valueMatch[1]+`"`)
		}
		return match
	})

	numberRe := regexp.MustCompile(`:\s*(-?\d+\.?\d*)`)
	styled = numberRe.ReplaceAllStringFunc(styled, func(match string) string {
		valueMatch := numberRe.FindStringSubmatch(match)
		if len(valueMatch) > 1 {
			prefix := match[:strings.Index(match, valueMatch[1])]
			return prefix + valueStyle.Render(valueMatch[1])
		}
		return match
	})

	boolNullRe := regexp.MustCompile(`:\s*(true|false|null)`)
	styled = boolNullRe.ReplaceAllStringFunc(styled, func(match string) string {
		valueMatch := boolNullRe.FindStringSubmatch(match)
		if len(valueMatch) > 1 {
			prefix := match[:strings.Index(match, valueMatch[1])]
			return prefix + valueStyle.Render(valueMatch[1])
		}
		return match
	})

	return styled
}

func HeadersPageView(m model) string {
	style1 := TitleStyle(m.termWidth)
	style2 := OptionsStyle(m.termWidth)
	style3 := HomePageStyle2(m.termWidth, m.termHeight)
	styleInput := inputStyle(m.termWidth)

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

			if m.pointer == i && m.editing {
				line = style4.Render("> ") + style5.Render(h.Key+" "+m.editingHeader.View()+"\n")
			} else if m.pointer == i && h.Value == "" {
				line = style4.Render("> ") + style5.Render(h.Key+" "+m.addHeaderValue.View()+"\n")
			} else if m.pointer == i {
				line = style4.Render("> ") + style5.Render(h.Key+" "+h.Value+"\n")
			} else {
				line = style4.Render("   ") + (h.Key + " " + h.Value + "\n")
			}

			items = append(items, line)
		}
	}

	var errorWarning string

	if m.hasError {
		errorStyle := errorStyle(m.termWidth)
		line := errorStyle.Render("⚠ ERROR: " + m.errorMessage + "\n\nPress 'x' to dismiss")
		errorWarning = line
	}

	leftBox := style2.Render(lipgloss.JoinVertical(lipgloss.Left, items...)) + "\n\n" + styleInput.Render(lipgloss.JoinVertical(lipgloss.Left, m.addHeaderKey.View())) + "\n\n" + errorWarning
	rightBox := style3.Render("Commands\n----------------\nESC -> Quit\n\nk -> Up\n\nj -> Down\n\n: -> Add New\n\nd -> Delete\n\nEnter -> Add Val\n\ne -> edit")
	layout := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)

	b.WriteString(layout)

	return b.String()
}

func QueryParamsPageView(m model) string {
	style1 := TitleStyle(m.termWidth)
	style2 := OptionsStyle(m.termWidth)
	style3 := HomePageStyle2(m.termWidth, m.termHeight)
	styleInput := inputStyle(m.termWidth)

	var b strings.Builder
	b.WriteString(style1.Render("QueryParams Page"))
	b.WriteString("\n")

	var items []string
	QueryParams := m.QueryParams

	if len(QueryParams) == 0 {
		line := "No Query Params\n\n"
		items = append(items, line)
	} else {
		for i, h := range QueryParams {
			var line string
			if m.pointer == i && m.editing {
				line = style4.Render("> ") + style5.Render(h.Key+" "+m.editingQueryParams.View()+"\n")
			} else if m.pointer == i && h.Value != "" {
				line = style4.Render("> ") + style5.Render(h.Key+" : "+h.Value+"\n")
			} else if m.pointer == i {
				line = style4.Render("> ") + style5.Render(h.Key+" : "+m.addQueryParamsValue.View()+"\n")
			} else {
				line = style4.Render("   ") + h.Key + " : " + h.Value + "\n"
			}
			items = append(items, line)
		}
	}

	var errorWarning string

	if m.hasError {
		errorStyle := errorStyle(m.termWidth)
		line := errorStyle.Render("⚠ ERROR: " + m.errorMessage + "\n\nPress 'x' to dismiss")
		errorWarning = line
	}

	leftBox := style2.Render(lipgloss.JoinVertical(lipgloss.Left, items...)) + "\n\n" + styleInput.Render(lipgloss.JoinVertical(lipgloss.Left, m.addQueryParamsKey.View())) + "\n\n" + errorWarning
	rightBox := style3.Render("Commands\n----------------\nESC -> Quit\n\nk -> Up\n\nj -> Down\n\n: -> Add New\n\nd -> Delete\n\nEnter -> Add Val\n\ne -> edit")
	layout := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)

	b.WriteString(layout)

	return b.String()
}

func loadingView(m model) string {
	style1 := loadingStyle(m.termWidth, m.termHeight)
	var b strings.Builder
	b.WriteString(style1.Render("LOADING..."))
	return b.String()
}

func VariablePageView(m model) string {
	style1 := OptionsStyle(m.termWidth)

	var b strings.Builder

	var responses []string

	if len(m.Responses) == 0 {
		line := "No Response loaded"
		responses = append(responses, line)
	}

	for i, v := range m.Responses {
		var line string
		if m.pointer == i && !m.VariablesFocus {
			line = style4.Render("> ") + style5.Render(v.Key+" : "+v.Value+"   press C to copy value... "+"\n")
		} else {
			line = "   " + v.Key + " : " + v.Value + "\n"
		}
		responses = append(responses, line)
	}
	var variables []string

	if len(m.LocalVariables) == 0 {
		line := "No Variables loaded"
		variables = append(variables, line)
	}

	for i, v := range m.LocalVariables {
		var line string
		if m.pointer == i && m.VariablesFocus {
			line = style4.Render("> " + style5.Render(v.Key+" : "+v.Value+"\n"))
		} else {
			line = "  " + v.Key + " : " + v.Value + "\n"
		}
		variables = append(variables, line)
	}

	b.WriteString(style1.Render(lipgloss.JoinVertical(lipgloss.Left, responses...)) + "\n\n" +
		style1.Render(lipgloss.JoinVertical(lipgloss.Left, variables...)))

	return b.String()
}
