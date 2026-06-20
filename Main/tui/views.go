package tui

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	switch m.CurrentPage {
	case HomePage:
		return Homepage(m)
	case CollectionPage:
		return Collectionpage(m.termWidth, m.termHeight, m)
	case ApiPage:
		return ApipageWithViewport(m)
	case RequestPage:
		return ReqPage(m, m.termWidth)
	case HeadersPage:
		return HeadersPageView(m, m.termWidth)
	case QueryParamsPage:
		return QueryParamsPageView(m, m.termWidth)
	case LoadingPage:
		return loadingView(m)
	case ResponsePage:
		return ResponsePageView(m, m.termWidth)
	case VariablesPage:
		return VariablesPageView(m)
	case DashBoard:
		return DashBoardView(m)
	}
	return ""
}

func Homepage(m Model) string {
	styleInput := inputStyle(m.termWidth)
	style3 := OptionsStyle(m.termWidth)
	style2 := HomePageStyle2(m.termWidth, m.termHeight)
	style1 := TitleStyle(m.termWidth)

	var b strings.Builder

	b.WriteString(style1.Render("Collections "))
	b.WriteString("\n")

	collections := m.storage.Collections

	var items []string
	if len(collections) == 0 {
		line := EmptyStateStyle.Render("No collections")
		items = append(items, line)
	} else {

		maxVisible := 12
		start := m.pageScrollOffset
		end := start + maxVisible

		if end > len(m.Collections) {
			end = len(m.Collections)
		}

		for i := start; i < end; i++ {
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
	}

	var errorWarning string

	if m.hasError {
		errorStyle := errorStyle(m.termWidth)
		line := errorStyle.Render("ERROR: " + m.errorMessage + "\n\nPress 'x' to dismiss")
		errorWarning = line
	}

	leftBox := style3.Render(lipgloss.JoinVertical(lipgloss.Left, items...)) + "\n\n" + styleInput.Render(lipgloss.JoinVertical(lipgloss.Left, m.NewCollectionInput.View())) + "\n\n" + errorWarning
	rightBox := style2.Render("Keys: esc quit | k/j move | enter open | : add | e edit | d delete")
	topHeight := lipgloss.Height(leftBox)
	bottomHeight := lipgloss.Height(rightBox)

	space := m.termHeight - topHeight - bottomHeight - 4
	if space < 0 {
		space = 0
	}

	layout := lipgloss.JoinVertical(lipgloss.Top, leftBox, strings.Repeat("\n", space), rightBox)

	b.WriteString(layout)
	return b.String()
}

func Collectionpage(termWidth, termHeight int, m Model) string {
	style1 := TitleStyle(termWidth)
	style2 := HomePageStyle2(termWidth, termHeight)
	style3 := OptionsStyle(termWidth)
	styleInput := inputStyle(m.termWidth)

	var b strings.Builder
	collectionName := m.SelectedCollection.Name

	b.WriteString(style1.Render(collectionName))
	b.WriteString("\n")

	var items []string

	if len(m.Apis) == 0 {
		line := EmptyStateStyle.Render("No APIs")
		items = append(items, line)
	} else {

		maxVisible := 12
		start := m.pageScrollOffset
		end := start + maxVisible

		if end > len(m.Apis) {
			end = len(m.Apis)
		}

		for i := start; i < end; i++ {
			api := m.Apis[i]

			if i == m.pointer && m.editing {
				line := style4.Render("> ") + m.editingApi.View() + "\n"
				items = append(items, line)
				continue
			}

			if i == m.pointer {
				text := style4.Render("> ") + ApiMethodStyle(api.Method).Render(api.Method) + " " + style5.Render(api.Url) + "\n"
				items = append(items, text)
			} else {
				text := "   " + ApiMethodStyle(api.Method).Render(api.Method) + " " + api.Url + "\n"
				items = append(items, text)
			}
		}
	}

	var errorWarning string

	if m.hasError {
		errorStyle := errorStyle(m.termWidth)
		line := errorStyle.Render("ERROR: " + m.errorMessage + "\n\nPress 'x' to dismiss")
		errorWarning = line
	}

	leftBox := style3.Render(lipgloss.JoinVertical(lipgloss.Left, items...)) + "\n\n" + renderApiAddInput(m, styleInput) + "\n\n" + errorWarning
	rightBox := style2.Render("Keys: esc back | k/j move | enter send/open | : add | e edit | h headers | q params | v variables | d delete")

	topHeight := lipgloss.Height(leftBox)
	bottomHeight := lipgloss.Height(rightBox)

	space := m.termHeight - topHeight - bottomHeight - 4
	if space < 0 {
		space = 0
	}

	layout := lipgloss.JoinVertical(lipgloss.Top, leftBox, strings.Repeat("\n", space), rightBox)

	b.WriteString(layout)
	return b.String()
}

func renderApiAddInput(m Model, style lipgloss.Style) string {
	if !m.NewApiInput.Focused() && !m.apiMethodSelecting {
		return style.Render("Press ':' to add a request")
	}

	content := []string{
		SectionTitleStyle.Render("New Request"),
		renderApiMethodOptions(m),
	}

	if m.apiMethodSelecting {
		content = append(content, EmptyStateStyle.Render("Select method with k/j or up/down, then press enter"))
	} else {
		content = append(content, "URL", m.NewApiInput.View())
	}

	return style.Render(lipgloss.JoinVertical(lipgloss.Left, content...))
}

func renderApiMethodOptions(m Model) string {
	options := m.apiMethodOptions
	if len(options) == 0 {
		options = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	}

	lines := make([]string, 0, len(options))
	for i, method := range options {
		prefix := "  "
		if i == m.apiMethodIndex {
			prefix = "> "
		}
		lines = append(lines, prefix+MethodOptionStyle(method, i == m.apiMethodIndex).Render(method))
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func selectedApiMethod(m Model) string {
	if len(m.apiMethodOptions) == 0 {
		return "GET"
	}
	if m.apiMethodIndex < 0 || m.apiMethodIndex >= len(m.apiMethodOptions) {
		return m.apiMethodOptions[0]
	}
	return m.apiMethodOptions[m.apiMethodIndex]
}

func ApipageWithViewport(m Model) string {
	if !m.viewportReady {
		return "Loading..."
	}

	helpText := HelpTextStyle.Render("\n\nKeys: k/j scroll | space/f page down | b page up | g/G top/bottom | r response values | e edit | esc back")
	return m.apiViewport.View() + helpText
}

func BuildApiPageContent(m Model, termWidth int) string {
	style1 := TitleStyle(termWidth)
	style3 := ResponseStyle(termWidth)
	style2 := statusLine(termWidth)
	styleInput := inputStyle(termWidth)

	var b strings.Builder
	SelectedApi := m.SelectedApi

	Response := m.ApiResponse

	statusStyle := StatusOKStyle
	if Response.StatusCode >= 400 {
		statusStyle = StatusErrorStyle
	}

	var resp strings.Builder

	resp.WriteString("Status: " + statusStyle.Render(Response.Status) + "\n")
	resp.WriteString(fmt.Sprintf("Status Code: %s\n", statusStyle.Render(fmt.Sprintf("%d", Response.StatusCode))))
	resp.WriteString(style2.Render(" "))

	if !m.responseComponent {
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

	}

	formattedBody := FormatJSON(Response.Body, bodyElementStyle, bodyElementStyle2)
	resp.WriteString("\nBody:\n" + formattedBody + "\n")

	b.WriteString(style1.Render("API Response"))

	b.WriteString("\n\n")
	if m.editing {
		b.WriteString(style3.Render("editing..." + "\n" + styleInput.Render(m.editingCurrentApi.View())))
	} else {
		b.WriteString(style3.Render(
			"Request: " +
				MethodStyle.Render(SelectedApi.Method) + " " + UrlStyle.Render(SelectedApi.Url),
		))
	}

	b.WriteString("\n\n")

	b.WriteString(style3.Render("Response:\n\n" + resp.String()))

	return b.String()
}

func ReqPage(m Model, termWidth int) string {
	style1 := TitleStyle(termWidth)
	style2 := OptionsStyle(termWidth)
	style3 := HomePageStyle2(termWidth, m.termHeight)
	styleInput := inputStyle(termWidth)

	name := m.SelectedApi.Method + "  " + m.SelectedApi.Url

	var b strings.Builder

	bodyFields := m.BodyFields

	b.WriteString(style1.Render(name))
	b.WriteString("\n")

	var items []string

	if len(bodyFields) == 0 {
		line := EmptyStateStyle.Render("No request fields\n\n")
		items = append(items, line)

	} else {

		maxVisible := 12
		start := m.pageScrollOffset
		end := start + maxVisible

		if end > len(bodyFields) {
			end = len(bodyFields)
		}

		for i := start; i < end; i++ {
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
		line := errorStyle.Render("ERROR: " + m.errorMessage + "\n\nPress 'x' to dismiss")
		errorWarning = line
	}
	leftBox := style2.Render(lipgloss.JoinVertical(lipgloss.Left, items...)) + "\n\n" + styleInput.Render(m.newBodyFieldInput.View()) + "\n\n" + errorWarning

	var rightBox string
	var space int
	if !m.reqPageComponent {
		rightBox = style3.Render("Keys: esc back | k/j move | enter send | : add field | v add value | e edit | d delete")
		topHeight := lipgloss.Height(leftBox)
		bottomHeight := lipgloss.Height(rightBox)

		space = m.termHeight - topHeight - bottomHeight - 4
		if space < 0 {
			space = 0
		}
	}

	layout := lipgloss.JoinVertical(lipgloss.Top, leftBox, strings.Repeat("\n", space), rightBox)

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

func HeadersPageView(m Model, termWidth int) string {
	style1 := TitleStyle(termWidth)
	style2 := OptionsStyle(termWidth)
	style3 := HomePageStyle2(termWidth, m.termHeight)
	styleInput := inputStyle(termWidth)

	name := m.SelectedApi.Method + "  " + m.SelectedApi.Url

	var b strings.Builder
	b.WriteString(style1.Render(name))

	headers := m.Headers
	b.WriteString("\n")

	var items []string

	if len(headers) == 0 {
		line := EmptyStateStyle.Render("No headers\n\n")
		items = append(items, line)
	} else {

		maxVisible := 12
		start := m.pageScrollOffset
		end := start + maxVisible

		if end > len(headers) {
			end = len(headers)
		}

		for i := start; i < end; i++ {
			var line string

			if m.pointer == i && m.editing {
				line = style4.Render("> ") + style5.Render(headers[i].Key+" : "+m.editingHeader.View()+"\n")
			} else if m.pointer == i && headers[i].Value == "" {
				line = style4.Render("> ") + style5.Render(headers[i].Key+" : "+m.addHeaderValue.View()+"\n")
			} else if m.pointer == i {
				line = style4.Render("> ") + style5.Render(headers[i].Key+" : "+headers[i].Value+"\n")
			} else {
				line = style4.Render("   ") + (headers[i].Key + " : " + headers[i].Value + "\n")
			}

			items = append(items, line)
		}
	}

	var errorWarning string

	if m.hasError {
		errorStyle := errorStyle(m.termWidth)
		line := errorStyle.Render("ERROR: " + m.errorMessage + "\n\nPress 'x' to dismiss")
		errorWarning = line
	}

	leftBox := style2.Render(lipgloss.JoinVertical(lipgloss.Left, items...)) + "\n\n" + styleInput.Render(lipgloss.JoinVertical(lipgloss.Left, m.addHeaderKey.View())) + "\n\n" + errorWarning

	var rightBox string
	var space int
	if !m.headersPageComponent {
		rightBox = style3.Render("Keys: esc back | k/j move | enter add value | : add header | e edit | d delete")
		topHeight := lipgloss.Height(leftBox)
		bottomHeight := lipgloss.Height(rightBox)

		space = m.termHeight - topHeight - bottomHeight - 4
		if space < 0 {
			space = 0
		}
	}

	layout := lipgloss.JoinVertical(lipgloss.Top, leftBox, strings.Repeat("\n", space), rightBox)

	b.WriteString(layout)
	return b.String()
}

func QueryParamsPageView(m Model, termWidth int) string {
	style1 := TitleStyle(termWidth)
	style2 := OptionsStyle(termWidth)
	style3 := HomePageStyle2(termWidth, m.termHeight)
	styleInput := inputStyle(termWidth)

	var b strings.Builder
	b.WriteString(style1.Render("Query Params"))
	b.WriteString("\n")

	var items []string
	QueryParams := m.QueryParams

	if len(QueryParams) == 0 {
		line := EmptyStateStyle.Render("No query params\n\n")
		items = append(items, line)
	} else {

		maxVisible := 12
		start := m.pageScrollOffset
		end := start + maxVisible

		if end > len(QueryParams) {
			end = len(QueryParams)
		}

		for i := start; i < end; i++ {
			var line string
			if m.pointer == i && m.editing {
				line = style4.Render("> ") + style5.Render(QueryParams[i].Key+" : "+m.editingQueryParams.View()+"\n")
			} else if m.pointer == i && QueryParams[i].Value != "" {
				line = style4.Render("> ") + style5.Render(QueryParams[i].Key+" : "+QueryParams[i].Value+"\n")
			} else if m.pointer == i {
				line = style4.Render("> ") + style5.Render(QueryParams[i].Key+" : "+m.addQueryParamsValue.View()+"\n")
			} else {
				line = style4.Render("   ") + QueryParams[i].Key + " : " + QueryParams[i].Value + "\n"
			}
			items = append(items, line)
		}
	}

	var errorWarning string

	if m.hasError {
		errorStyle := errorStyle(m.termWidth)
		line := errorStyle.Render("ERROR: " + m.errorMessage + "\n\nPress 'x' to dismiss")
		errorWarning = line
	}

	leftBox := style2.Render(lipgloss.JoinVertical(lipgloss.Left, items...)) + "\n\n" + styleInput.Render(lipgloss.JoinVertical(lipgloss.Left, m.addQueryParamsKey.View())) + "\n\n" + errorWarning

	var rightBox string
	var space int

	if !m.headersPageComponent {
		rightBox = style3.Render("Keys: esc back | k/j move | enter add value | : add param | e edit | d delete")
		topHeight := lipgloss.Height(leftBox)
		bottomHeight := lipgloss.Height(rightBox)

		space = m.termHeight - topHeight - bottomHeight - 4
		if space < 0 {
			space = 0
		}
	}

	layout := lipgloss.JoinVertical(lipgloss.Top, leftBox, strings.Repeat("\n", space), rightBox)

	b.WriteString(layout)
	return b.String()
}

func loadingView(m Model) string {
	style1 := loadingStyle(m.termWidth, m.termHeight)
	var b strings.Builder
	b.WriteString(style1.Render("LOADING..."))
	return b.String()
}

func ResponsePageView(m Model, termWidth int) string {
	style1 := OptionsStyle(termWidth - 4)
	style3 := HomePageStyle2(termWidth, m.termHeight)
	titleStyle := TitleStyle(termWidth)

	var b strings.Builder

	b.WriteString(titleStyle.Render("Response Page"))
	b.WriteString("\n")

	var responses []string

	responseRows := buildResponseTreeRows(m.ApiResponse.Body, m.ResponseExpanded)
	if len(responseRows) == 0 {
		line := EmptyStateStyle.Render("No response loaded")
		responses = append(responses, line)
	} else {
		start, end := visibleResponseWindow(m, responseRows)
		valueWidth := termWidth - 24
		if valueWidth < 16 {
			valueWidth = 16
		}
		for i := start; i < end; i++ {
			row := responseRows[i]
			var line string

			toggle := " "
			if row.Expandable {
				if row.Expanded {
					toggle = "v"
				} else {
					toggle = ">"
				}
			}

			indent := strings.Repeat("  ", row.Depth)
			text := fmt.Sprintf("%s%s %s: %s", indent, toggle, row.Key, truncateText(row.Value, valueWidth-row.Depth*2))
			if m.pointer == i && !m.VariablesFocus {
				line = style4.Render("> ") + style5.Render(text) + "\n"
			} else {
				line = "   " + text + "\n"
			}
			responses = append(responses, line)
		}
	}
	var variables []string

	if len(m.LocalVariables) == 0 {
		line := EmptyStateStyle.Render("No variables loaded")
		variables = append(variables, line)
	} else {

		maxVisible := 5
		start := m.variableScrollOffset
		end := start + maxVisible

		if end > len(m.LocalVariables) {
			end = len(m.LocalVariables)
		}

		for i := start; i < end; i++ {
			v := m.LocalVariables[i]
			var line string
			if m.pointer == i && m.VariablesFocus {
				line = style4.Render("> ") + style5.Render(v.Key+" : "+v.Value) + "\n"
			} else {
				line = "   " + v.Key + " : " + v.Value + "\n"
			}
			variables = append(variables, line)
		}
	}

	responsePanel := lipgloss.JoinVertical(lipgloss.Left, SectionTitleStyle.Render("Response Values"), "", lipgloss.JoinVertical(lipgloss.Left, responses...))
	variablesPanel := lipgloss.JoinVertical(lipgloss.Left, SectionTitleStyle.Render("Local Variables"), "", lipgloss.JoinVertical(lipgloss.Left, variables...))
	leftBox := style1.Render(responsePanel) + "\n\n" + style1.Render(variablesPanel)

	var rightBox string
	var space int

	if !m.resComponent {
		rightBox = style3.Render("Keys: esc back | k/j move | o open/close | enter save variable | c copy | v/r switch pane | d delete variable")
		topHeight := lipgloss.Height(leftBox)
		bottomHeight := lipgloss.Height(rightBox)

		space = m.termHeight - topHeight - bottomHeight - 4
		if space < 0 {
			space = 0
		}
	}

	layout := lipgloss.JoinVertical(lipgloss.Top, leftBox, strings.Repeat("\n", space), rightBox)

	b.WriteString(layout)
	return b.String()
}

func VariablesPageView(m Model) string {
	style1 := OptionsStyle(m.termWidth - 4)
	style2 := TitleStyle(m.termWidth)
	style3 := HomePageStyle2(m.termWidth, m.termHeight)
	styleInput := inputStyle(m.termWidth)

	var b strings.Builder
	b.WriteString(style2.Render("Variables"))
	b.WriteString("\n")

	var items []string

	if len(m.LocalVariables) == 0 {
		line := EmptyStateStyle.Render("No variables\n\n")
		items = append(items, line)
	} else {

		maxVisible := 12
		start := m.pageScrollOffset
		end := start + maxVisible

		if end > len(m.LocalVariables) {
			end = len(m.LocalVariables)
		}

		for i := start; i < end; i++ {
			var line string
			if m.pointer == i && m.editing {
				line = style4.Render("> ") + style5.Render(m.LocalVariables[i].Key+" "+m.editingLocalVariables.View()+"\n")
			} else if m.pointer == i && m.LocalVariables[i].Value != "" {
				line = style4.Render("> ") + style5.Render(m.LocalVariables[i].Key+" : "+m.LocalVariables[i].Value+"\n")
			} else if m.pointer == i {
				line = style4.Render("> ") + style5.Render(m.LocalVariables[i].Key+" : "+m.addVariableValue.View()+"\n")
			} else {
				line = style4.Render("   ") + m.LocalVariables[i].Key + " : " + m.LocalVariables[i].Value + "\n"
			}
			items = append(items, line)
		}
	}

	var errorWarning string

	if m.hasError {
		errorStyle := errorStyle(m.termWidth)
		line := errorStyle.Render("ERROR: " + m.errorMessage + "\n\nPress 'x' to dismiss")
		errorWarning = line
	}

	leftBox := style1.Render(lipgloss.JoinVertical(lipgloss.Left, items...)) + "\n\n" + styleInput.Render(lipgloss.JoinVertical(lipgloss.Left, m.addVariableKey.View())) + "\n\n" + errorWarning
	rightBox := style3.Render("Keys: esc back | k/j move | enter add value | : add variable | e edit | d delete")
	topHeight := lipgloss.Height(leftBox)
	bottomHeight := lipgloss.Height(rightBox)

	space := m.termHeight - topHeight - bottomHeight - 4
	if space < 0 {
		space = 0
	}

	layout := lipgloss.JoinVertical(lipgloss.Top, leftBox, strings.Repeat("\n", space), rightBox)

	b.WriteString(layout)
	return b.String()
}

func DashBoardView(m Model) string {
	var b strings.Builder

	title := TitleStyle(m.termWidth).Render("Dashboard")
	height := safeHeight(m.termHeight, 3)

	requests := m.SelectedCollection.Requests
	var items []string

	if len(requests) == 0 {
		items = append(items, EmptyStateStyle.Render("No requests"))
	} else {
		maxItems := height
		if len(requests) < maxItems {
			maxItems = len(requests)
		}

		for i := 0; i < maxItems; i++ {
			if m.pointer == i {
				items = append(items, style4.Render("> ")+style5.Render(requests[i].Method+" : "+requests[i].Url)+"\n")
			} else {
				items = append(items, style4.Render("   "+requests[i].Method+" : "+requests[i].Url)+"\n")
			}
		}
	}

	var rightContent string
	useSplit := false

	leftWidth := m.termWidth/3 - 6
	if leftWidth < 18 {
		leftWidth = 18
	}
	if leftWidth > m.termWidth-32 {
		leftWidth = m.termWidth - 32
	}
	if leftWidth < 12 {
		leftWidth = 12
	}
	rightWidth := m.termWidth - leftWidth - 2

	switch {
	case m.responseComponent:
		rightContent = ApipageWithViewport(m)
		useSplit = true
	case m.reqPageComponent:
		rightContent = ReqPage(m, rightWidth)
		useSplit = true
	case m.headersPageComponent:
		rightContent = HeadersPageView(m, rightWidth)
		useSplit = true
	case m.paramsComponent:
		rightContent = QueryParamsPageView(m, rightWidth)
		useSplit = true
	case m.resComponent:
		rightContent = ResponsePageView(m, rightWidth)
		useSplit = true
	}

	if useSplit {
		leftBox := OptionsStyle(leftWidth).
			Render(lipgloss.JoinVertical(lipgloss.Left, items...))

		rightBox := DashBoardResponseStyle().Width(rightWidth).Height(safeHeight(height, 2)).Render(rightContent)

		layout := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)
		full := lipgloss.JoinVertical(lipgloss.Top, title, layout)

		b.WriteString(full)
		return b.String()
	}

	leftBox := OptionsStyle(m.termWidth).
		Render(lipgloss.JoinVertical(lipgloss.Left, items...))

	full := lipgloss.JoinVertical(lipgloss.Top, title, leftBox)

	b.WriteString(full)
	return b.String()
}
