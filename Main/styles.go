package main

import "github.com/charmbracelet/lipgloss"

func HomePageStyle1(termWidth int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("210")).
		PaddingRight(termWidth - 35).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("202")).
		Align(lipgloss.Center)
}

func HomePageStyle2(termWidth, termHeight int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("210")).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("202")).
		Height(termHeight - 7).
		PaddingRight(1).PaddingLeft(1).
		Align(lipgloss.Center).
		MarginLeft(2)
}

func OptionsStyle(termWidth int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("210")).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#3C3C3C")).
		Width(termWidth - 25).
		Padding(1)
}

var style4 = lipgloss.NewStyle().
	Foreground(lipgloss.Color("205")).
	Bold(true)

var style5 = lipgloss.NewStyle().
	Foreground(lipgloss.Color("150")).
	Bold(true)

var MethodStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("2")).
	Bold(true)

var UrlStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("4")).
	Bold(true)

var StatusOKStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))   // Green
var StatusErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")) // Red

func InputStyle(termWidth int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("210")).
		Width(termWidth - 3).
		Height(25).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("202"))

}
func ResponseStyle(termWidth int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("210")).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#3C3C3C")).
		Width(termWidth - 3).
		Padding(1)
}

var HelpTextStyle = lipgloss.NewStyle().
	Align(lipgloss.Center)
