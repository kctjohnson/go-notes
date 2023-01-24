package models

import tea "github.com/charmbracelet/bubbletea"

type Edit struct {
}

func (m Edit) Init() tea.Cmd {
	return nil
}

func (m Edit) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Edit) View() string {
	return ""
}
