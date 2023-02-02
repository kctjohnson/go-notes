package models

import tea "github.com/charmbracelet/bubbletea"

type Tag string

type TagList struct {
	Tags []Tag
}

func NewTagList() *TagList {
	return nil
}

func (m TagList) Init() tea.Cmd {
	return nil
}

func (m TagList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m TagList) View() string {
	return ""
}
