package UI

import "github.com/charmbracelet/lipgloss"

var ErrorStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#e74c3c"))

var DebugStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#95a5a6"))

var InfoStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#3498db"))

var CodeStyle = lipgloss.NewStyle().
	Bold(true).Foreground(lipgloss.Color("#2980b9"))
