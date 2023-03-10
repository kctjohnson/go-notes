package list

import "github.com/charmbracelet/bubbles/key"

type listKeymap struct {
	Up       key.Binding
	Down     key.Binding
	ViewUp   key.Binding
	ViewDown key.Binding
	Select   key.Binding
	New      key.Binding
	Delete   key.Binding
	Filter   key.Binding
	SetTag   key.Binding
	Help     key.Binding
	Quit     key.Binding
}

func (k listKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Help, k.Quit}
}

func (k listKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.New, k.ViewUp, k.Select, k.Help},
		{k.Down, k.Delete, k.ViewDown, k.Filter, k.SetTag, k.Quit},
	}
}

var Keys = listKeymap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "Up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "Down"),
	),
	ViewUp: key.NewBinding(
		key.WithKeys("shift+up", "K"),
		key.WithHelp("⇧+↑/K", "Preview Up"),
	),
	ViewDown: key.NewBinding(
		key.WithKeys("shift+down", "J"),
		key.WithHelp("⇧+↓/J", "Preview Down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("Space", "Select"),
	),
	New: key.NewBinding(
		key.WithKeys("N"),
		key.WithHelp("N", "New"),
	),
	Delete: key.NewBinding(
		key.WithKeys("D"),
		key.WithHelp("D", "Delete"),
	),
	Filter: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "Filter"),
	),
	SetTag: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "Set Tag"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "Help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("esc", "Quit"),
	),
}
