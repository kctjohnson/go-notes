package models

import (
	"fmt"
	"go-notes/cmd/gonotes/utils"
	"go-notes/pkg/db/model"
	"go-notes/pkg/services"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type List struct {
	notes   []model.Note
	keys    listKeymap
	preview viewport.Model
	help    help.Model

	width        int
	cursor       int
	scrollIndex  int
	maxViewNotes int

	noteService *services.NotesService
}

func NewList(keys listKeymap, ns *services.NotesService) *List {
	// Disable the built in viewport keybindings
	preview := viewport.New(40, 15)
	preview.KeyMap.Down.SetEnabled(false)
	preview.KeyMap.Up.SetEnabled(false)
	preview.KeyMap.PageDown.SetEnabled(false)
	preview.KeyMap.PageUp.SetEnabled(false)
	preview.KeyMap.HalfPageDown.SetEnabled(false)
	preview.KeyMap.HalfPageUp.SetEnabled(false)

	return &List{
		notes:        []model.Note{},
		keys:         keys,
		preview:      preview,
		help:         help.New(),
		width:        0,
		cursor:       0,
		scrollIndex:  0,
		maxViewNotes: 15,
		noteService:  ns,
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
		case key.Matches(msg, m.keys.Up):
			m.cursorUp()
			m.preview.GotoTop()
		case key.Matches(msg, m.keys.Down):
			m.cursorDown()
			m.preview.GotoTop()
		case key.Matches(msg, m.keys.ViewUp):
			m.preview.LineUp(1)
		case key.Matches(msg, m.keys.ViewDown):
			m.preview.LineDown(1)
		case key.Matches(msg, m.keys.New):
			return m, utils.CreateNoteCmd(m.noteService, "New Note Title")
		case key.Matches(msg, m.keys.Delete):
			return m, utils.DeleteNoteCmd(m.noteService, m.notes[m.cursor].ID)
		case key.Matches(msg, m.keys.Select):
			return m, utils.EditNoteCmd(m.noteService, m.notes[m.cursor].ID)
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			m.preview.Width = m.width - lipgloss.Width(m.View()) - 2
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
	}

	// Update the height of the preview
	if len(m.notes) > 0 {
		contentHeight := lipgloss.Height(m.listView())
		m.preview.Height = contentHeight
		m.preview.SetContent(m.notes[m.cursor].Content)
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
	for row := m.scrollIndex; row < len(m.notes) && row < m.scrollIndex+m.maxViewNotes; row++ {
		if row == m.cursor {
			str += fmt.Sprintf("%s\n", utils.FocusedNoteStyle.Render(m.notes[row].Title))
		} else {
			str += fmt.Sprintf("%s\n", m.notes[row].Title)
		}
	}
	str += "\n" + m.help.View(listKeys)
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
	if m.cursor == len(m.notes)-1 {
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

	if m.cursor >= len(m.notes) {
		m.cursor = len(m.notes) - 1
	}
}
