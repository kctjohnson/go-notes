package utils

import "github.com/charmbracelet/lipgloss"

var (
	UnderlinedStyle  = lipgloss.NewStyle().Underline(true)
	FocusedLineStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("69")).
				Bold(true)
	TitleStyle = lipgloss.NewStyle().
			Underline(true)
	PreviewStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), false, false, false, true).
			PaddingLeft(1).MarginLeft(1)
)
