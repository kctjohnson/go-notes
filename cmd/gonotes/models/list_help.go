package models

import "github.com/charmbracelet/bubbles/key"

type listKeymap struct {
	Up       key.Binding
	Down     key.Binding
	ViewUp   key.Binding
	ViewDown key.Binding
	Select   key.Binding
	New      key.Binding
	Delete   key.Binding
	Help     key.Binding
	Quit     key.Binding
}

func (k listKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Help, k.Quit}
}

func (k listKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.New, k.ViewUp, k.Select},
		{k.Down, k.Delete, k.ViewDown, k.Help, k.Quit},
	}
}

var listKeys = listKeymap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	ViewUp: key.NewBinding(
		key.WithKeys("shift+up", "K"),
		key.WithHelp("⇧+↑/K", "preview up"),
	),
	ViewDown: key.NewBinding(
		key.WithKeys("shift+down", "J"),
		key.WithHelp("⇧+↓/J", "preview down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("Enter", "Select"),
	),
	New: key.NewBinding(
		key.WithKeys("N"),
		key.WithHelp("N", "New"),
	),
	Delete: key.NewBinding(
		key.WithKeys("D"),
		key.WithHelp("D", "Delete"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("esc", "quit"),
	),
}