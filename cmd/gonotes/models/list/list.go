package list

import (
	"fmt"
	"go-notes/cmd/gonotes/utils"
	"go-notes/internal/db/model"
	"go-notes/internal/graphql"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type List struct {
	Notes   []model.Note
	preview viewport.Model
	help    help.Model

	width        int
	cursor       int
	scrollIndex  int
	maxViewNotes int

	gqlClient *graphql.Client
}

func New(gqlClient *graphql.Client) *List {
	// Disable the built in viewport keybindings
	preview := viewport.New(40, 15)
	preview.KeyMap.Down.SetEnabled(false)
	preview.KeyMap.Up.SetEnabled(false)
	preview.KeyMap.PageDown.SetEnabled(false)
	preview.KeyMap.PageUp.SetEnabled(false)
	preview.KeyMap.HalfPageDown.SetEnabled(false)
	preview.KeyMap.HalfPageUp.SetEnabled(false)

	return &List{
		Notes:        []model.Note{},
		preview:      preview,
		help:         help.New(),
		width:        0,
		cursor:       0,
		scrollIndex:  0,
		maxViewNotes: 15,
		gqlClient:    gqlClient,
	}
}

func (m List) Init() tea.Cmd {
	return nil
}

func (m List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Constrain the cursor to prevent it from going out of bounds
	m.constrainCursor()

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.help.Width = msg.Width
		m.preview.Width = msg.Width - lipgloss.Width(m.View()) - 2
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Up):
			m.cursorUp()
			m.preview.GotoTop()
		case key.Matches(msg, Keys.Down):
			m.cursorDown()
			m.preview.GotoTop()
		case key.Matches(msg, Keys.ViewUp):
			m.preview.LineUp(1)
		case key.Matches(msg, Keys.ViewDown):
			m.preview.LineDown(1)
		case key.Matches(msg, Keys.New):
			return m, utils.CreateNoteCmd(m.gqlClient, "New Note Title")
		case key.Matches(msg, Keys.Delete):
			return m, utils.DeleteNoteCmd(m.gqlClient, m.Notes[m.cursor].ID)
		case key.Matches(msg, Keys.Select):
			return m, utils.EditNoteCmd(m.gqlClient, m.Notes[m.cursor].ID)
		case key.Matches(msg, Keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			m.preview.Width = m.width - lipgloss.Width(m.View()) - 2
		case key.Matches(msg, Keys.Quit):
			return m, tea.Quit
		}
	}

	// Update the height of the preview
	if len(m.Notes) > 0 {
		contentHeight := lipgloss.Height(m.listView())
		m.preview.Height = contentHeight
		m.preview.SetContent(m.Notes[m.cursor].Content)
	}
	previewModel, cmd := m.preview.Update(msg)
	m.preview = previewModel

	return m, cmd
}

func (m List) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top, m.listView(), m.previewView())
}

func (m List) listView() string {
	str := utils.TitleStyle.Render("Notes:") + "\n"
	for row := m.scrollIndex; row < len(m.Notes) && row < m.scrollIndex+m.maxViewNotes; row++ {
		if row == m.cursor {
			str += fmt.Sprintf("%s\n", utils.FocusedLineStyle.Render(m.Notes[row].Title))
		} else {
			str += fmt.Sprintf("%s\n", m.Notes[row].Title)
		}
	}
	str += "\n" + m.help.View(Keys)
	return str
}

func (m List) previewView() string {
	return utils.PreviewStyle.Render(m.preview.View())
}

func (m *List) cursorUp() {
	if m.cursor == 0 {
		return
	} else if m.cursor == m.scrollIndex {
		m.cursor--
		m.scrollIndex--
	} else {
		m.cursor--
	}
}

func (m *List) cursorDown() {
	if m.cursor == len(m.Notes)-1 {
		return
	} else if m.cursor == m.scrollIndex+m.maxViewNotes-1 {
		m.cursor++
		m.scrollIndex++
	} else {
		m.cursor++
	}
}

func (m *List) constrainCursor() {
	if m.cursor < 0 {
		m.cursor = 0
	}

	if m.cursor >= len(m.Notes) {
		m.cursor = len(m.Notes) - 1
	}
}
