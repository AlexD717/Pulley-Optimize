package main

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#009bce")).Bold(true)
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#9b9b9b")).Italic(true)
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#0a9ee8"))

	errorStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff3700"))
	activeStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#c75d00")).Bold(true)
	unActiveStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#cccccc"))
	activeTextStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#d2d2d2")).Bold(true)
	unActiveTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#bebebe"))
	fixedWidthText    = lipgloss.NewStyle().Width(35)

	boxTitleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
	resultBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#0707b8")).
			Padding(1, 2)
	resultsTableStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
	leftColumnStyle   = lipgloss.NewStyle().PaddingRight(5)
)

func RenderWithMinSize(content string, minWidth int, minHeight int) string {
	actualWidth := lipgloss.Width(content) + 4 // + 4 to account for default table padding
	actualHeight := lipgloss.Height(content) + 2

	finalWidth := max(actualWidth, minWidth)
	finalHeight := max(actualHeight, minHeight)

	return resultBoxStyle.Width(finalWidth).Height(finalHeight).Render(content)
}
