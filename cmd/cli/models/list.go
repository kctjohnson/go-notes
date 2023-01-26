package models

import (
	"fmt"
	"go-notes/cmd/cli/utils"
	"go-notes/pkg/db/model"
	"go-notes/pkg/services"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type List struct {
	notes []model.Note
	keys  listKeymap
	help  help.Model

	cursor       int
	scrollIndex  int
	maxViewNotes int

	noteService *services.NotesService
}

func NewList(keys listKeymap, ns *services.NotesService) *List {
	return &List{
		notes:        []model.Note{},
		keys:         keys,
		help:         help.New(),
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
		m.help.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			m.cursorUp()
		case key.Matches(msg, m.keys.Down):
			m.cursorDown()
		case key.Matches(msg, m.keys.New):
			return m, utils.CreateNoteCmd(m.noteService, "New Note Title")
		case key.Matches(msg, m.keys.Delete):
			return m, utils.DeleteNoteCmd(m.noteService, m.notes[m.cursor].ID)
		case key.Matches(msg, m.keys.Select):
			return m, utils.EditNoteCmd(m.noteService, m.notes[m.cursor].ID)
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m List) View() string {
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
