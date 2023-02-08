package models

import (
	"github.com/charmbracelet/bubbles/key"
)

type tagsKeymap struct {
	Up     key.Binding
	Down   key.Binding
	Toggle key.Binding
	Back   key.Binding
	New    key.Binding
	Delete key.Binding
	Help   key.Binding
	Quit   key.Binding
}

func (k tagsKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Toggle, k.Back, k.Help, k.Quit}
}

func (k tagsKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Toggle, k.New, k.Help},
		{k.Down, k.Back, k.Delete, k.Quit},
	}
}

var tagsKeys = tagsKeymap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "Up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "Down"),
	),
	Toggle: key.NewBinding(
		key.WithKeys(" ", "enter"),
		key.WithHelp("<space>/<enter>", "Toggle"),
	),
	New: key.NewBinding(
		key.WithKeys("N"),
		key.WithHelp("N", "New"),
	),
	Delete: key.NewBinding(
		key.WithKeys("D"),
		key.WithHelp("D", "Delete"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "q"),
		key.WithHelp("q/esc", "Back"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "Help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("esc", "Quit"),
	),
}
