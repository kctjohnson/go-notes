package tags

import (
	"github.com/charmbracelet/bubbles/key"
)

type inputKeymap struct {
	Create key.Binding
	Back   key.Binding
	Quit   key.Binding
}

func (k inputKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Create, k.Back, k.Quit}
}

func (k inputKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Create, k.Back, k.Quit},
	}
}

var InputKeys = inputKeymap{
	Create: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("<enter>", "Create"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("esc", "Quit"),
	),
}
