package models

import (
	"fmt"
	"go-notes/pkg/db/model"

	tea "github.com/charmbracelet/bubbletea"
)

type List struct {
	notes []model.Note
}

func (m List) Init() tea.Cmd {
	return nil
}

func (m List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m List) View() string {
	str := "Notes:\n"
	for _, n := range m.notes {
		str += fmt.Sprintf("%s\n", n.Title)
	}
	return str
}
