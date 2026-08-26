package main

import "github.com/charmbracelet/lipgloss"

var (
	activeStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#980aeb")).Bold(true)
	unActiveStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#cccccc"))
	activeTextStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#d2d2d2")).Bold(true)
	unActiveTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#bebebe"))

	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#9b9b9b")).Italic(true)
)
