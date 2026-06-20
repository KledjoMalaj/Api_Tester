package tui

import "github.com/charmbracelet/lipgloss"

const (
	colorText     = lipgloss.Color("252")
	colorMuted    = lipgloss.Color("245")
	colorSubtle   = lipgloss.Color("238")
	colorPanel    = lipgloss.Color("240")
	colorAccent   = lipgloss.Color("39")
	colorSelected = lipgloss.Color("154")
	colorInput    = lipgloss.Color("81")
	colorError    = lipgloss.Color("203")
	colorOK       = lipgloss.Color("42")
	colorWarning  = lipgloss.Color("214")
)

func safeWidth(width int, padding int) int {
	if width-padding < 12 {
		return 12
	}
	return width - padding
}

func safeHeight(height int, padding int) int {
	if height-padding < 1 {
		return 1
	}
	return height - padding
}

func HomePageStyle1(termWidth int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(colorText).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(colorAccent).
		Width(safeWidth(termWidth, 23)).
		Align(lipgloss.Center)
}

func TitleStyle(termWidth int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(colorText).
		Background(colorSubtle).
		Bold(true).
		Width(safeWidth(termWidth, 2)).
		Padding(0, 1).
		Align(lipgloss.Left)
}

func HomePageStyle2(termWidth, termHeight int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(colorMuted).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(colorPanel).
		Width(safeWidth(termWidth, 2)).
		Height(safeHeight(termHeight, 50)).
		Padding(0, 1).
		Align(lipgloss.Left)
}

func OptionsStyle(termWidth int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(colorText).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(colorPanel).
		Width(safeWidth(termWidth, 2)).
		Padding(1)
}

var style4 = lipgloss.NewStyle().
	Foreground(colorAccent).
	Bold(true)

var style5 = lipgloss.NewStyle().
	Foreground(colorSelected).
	Bold(true)

var MethodStyle = lipgloss.NewStyle().
	Foreground(colorOK).
	Bold(true)

func ApiMethodStyle(method string) lipgloss.Style {
	style := lipgloss.NewStyle().Bold(true)
	switch method {
	case "GET":
		return style.Foreground(lipgloss.Color("42"))
	case "POST":
		return style.Foreground(lipgloss.Color("39"))
	case "PUT":
		return style.Foreground(lipgloss.Color("214"))
	case "DELETE":
		return style.Foreground(lipgloss.Color("203"))
	case "PATCH":
		return style.Foreground(lipgloss.Color("171"))
	default:
		return style.Foreground(colorMuted)
	}
}

func MethodOptionStyle(method string, selected bool) lipgloss.Style {
	style := ApiMethodStyle(method).
		Width(8).
		Padding(0, 1)
	if selected {
		return style.
			Background(colorSubtle).
			BorderStyle(lipgloss.NormalBorder()).
			BorderLeft(true).
			BorderForeground(colorAccent)
	}
	return style.Foreground(colorMuted)
}

var UrlStyle = lipgloss.NewStyle().
	Foreground(colorAccent).
	Bold(true)

var StatusOKStyle = lipgloss.NewStyle().
	Foreground(colorOK).
	Bold(true)

var StatusErrorStyle = lipgloss.NewStyle().
	Foreground(colorError).
	Bold(true)

func InputStyle(termWidth int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(colorInput).
		Width(safeWidth(termWidth, 3)).
		Height(25).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(colorAccent)

}
func ResponseStyle(termWidth int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(colorText).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(colorPanel).
		Width(safeWidth(termWidth, 3)).
		Padding(1)
}

func statusLine(termWidth int) lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(colorPanel).
		BorderTop(true).
		Width(safeWidth(termWidth, 5))
}

var HelpTextStyle = lipgloss.NewStyle().
	Foreground(colorMuted).
	Align(lipgloss.Center)

var SectionTitleStyle = lipgloss.NewStyle().
	Foreground(colorMuted).
	Bold(true)

var EmptyStateStyle = lipgloss.NewStyle().
	Foreground(colorMuted).
	Italic(true)

var bodyElementStyle = lipgloss.NewStyle().
	Foreground(colorInput)

var bodyElementStyle2 = lipgloss.NewStyle().
	Foreground(colorSelected)

func loadingStyle(termWidth int, termHeight int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(colorOK).
		Bold(true).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(colorPanel).
		Width(safeWidth(termWidth, 3)).
		MarginTop(termHeight / 2).
		Align(lipgloss.Center)
}

func inputStyle(termWidth int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(colorInput).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(colorPanel).
		Width(safeWidth(termWidth, 2)).
		Padding(0, 1)
}

func errorStyle(termWidth int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(colorError).
		Bold(true).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(colorError).
		Padding(1).
		Width(safeWidth(termWidth, 2)).
		Align(lipgloss.Left)
}

func DashBoardResponseStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(colorPanel)
}
