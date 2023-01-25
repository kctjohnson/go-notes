package utils

import "github.com/charmbracelet/lipgloss"

var (
	FocusedNoteStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("69")).
				Bold(true)
	TitleStyle = lipgloss.NewStyle().
			Underline(true)
)
