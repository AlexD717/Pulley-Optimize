package main

import "github.com/charmbracelet/lipgloss"

var (
	activeStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#7e00c7")).Bold(true)
	unActiveStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#cccccc"))
	activeTextStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#d2d2d2")).Bold(true)
	unActiveTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#bebebe"))
	fixedWidthText    = lipgloss.NewStyle().Width(40)

	titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#530fe6")).Bold(true)
	slashStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#530fe6"))
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#9b9b9b")).Italic(true)

	boxTitleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
	resultBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#530fe6")).
			Padding(1, 2).
			Width(45).
			Height(12)
	leftColumnStyle = lipgloss.NewStyle().PaddingRight(5)
)
