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

type State int

const (
	NOTES State = iota
	SET_TAG
)

type List struct {
	Notes   []model.Note
	preview viewport.Model
	SetTag  SetTag
	help    help.Model

	state        State
	width        int
	cursor       int
	scrollIndex  int
	maxViewNotes int

	currentFilter *model.Tag

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
		Notes:         []model.Note{},
		preview:       preview,
		SetTag:        *NewSetTag(gqlClient, []model.Tag{}),
		help:          help.New(),
		state:         NOTES,
		width:         0,
		cursor:        0,
		scrollIndex:   0,
		maxViewNotes:  15,
		currentFilter: nil,
		gqlClient:     gqlClient,
	}
}

func (m List) Init() tea.Cmd {
	return nil
}

func (m List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Constrain the cursor to prevent it from going out of bounds
	notes := m.getActiveNotes()
	m.constrainCursor()

	var cmds tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.help.Width = msg.Width
		m.preview.Width = msg.Width - lipgloss.Width(m.View()) - 2
	case tea.KeyMsg:
		if m.state == NOTES {
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
				if m.currentFilter != nil {
					return m, utils.CreateNoteCmd(m.gqlClient, "New Note Title", &m.currentFilter.ID)
				} else {
					return m, utils.CreateNoteCmd(m.gqlClient, "New Note Title", nil)
				}
			case key.Matches(msg, Keys.Delete):
				return m, utils.DeleteNoteCmd(m.gqlClient, notes[m.cursor].ID)
			case key.Matches(msg, Keys.Select):
				return m, utils.EditNoteCmd(m.gqlClient, notes[m.cursor].ID)
			case key.Matches(msg, Keys.SetTag):
				m.state = SET_TAG
				return m, utils.LoadSetTagsCmd(m.gqlClient)
			case key.Matches(msg, Keys.Help):
				m.help.ShowAll = !m.help.ShowAll
				m.preview.Width = m.width - lipgloss.Width(m.View()) - 2
			case key.Matches(msg, Keys.Quit):
				return m, tea.Quit
			}
		} else if m.state == SET_TAG {
			switch {
			case key.Matches(msg, SetTagKeys.Select):
				m.state = NOTES
				return m, utils.SetNoteTagCmd(m.gqlClient, notes[m.cursor].ID, m.SetTag.Tags[m.SetTag.cursor].ID)
			case key.Matches(msg, SetTagKeys.Back):
				m.state = NOTES
				return m, nil
			}

			model, newCmd := m.SetTag.Update(msg)
			m.SetTag = model.(SetTag)
			cmds = tea.Batch(cmds, newCmd)
		}
	}

	// Update the height of the preview
	if len(notes) > 0 {
		contentHeight := lipgloss.Height(m.listView())
		m.preview.Height = contentHeight
		m.preview.SetContent(notes[m.cursor].Content)
	}
	previewModel, cmd := m.preview.Update(msg)
	m.preview = previewModel

	return m, tea.Batch(cmds, cmd)
}

func (m List) View() string {
	if m.state == NOTES {
		return lipgloss.JoinHorizontal(lipgloss.Top, m.listView(), m.previewView())
	} else {
		return m.SetTag.View()
	}
}

func (m *List) SetFilter(tag *model.Tag) {
	m.currentFilter = tag
}

func (m *List) RemoveFilter() {
	m.currentFilter = nil
}

func (m List) listView() string {
	notes := m.getActiveNotes()
	str := utils.TitleStyle.Render("Notes:") + "\n"
	for row := m.scrollIndex; row < len(notes) && row < m.scrollIndex+m.maxViewNotes; row++ {
		if row == m.cursor {
			str += fmt.Sprintf("%s", utils.FocusedLineStyle.Render(notes[row].Title))
		} else {
			str += fmt.Sprintf("%s", notes[row].Title)
		}

		if notes[row].TagID != nil {
			if *notes[row].TagID != int64(-1) {
				tag := model.Tag{
					ID:   -1,
					Name: "All",
				}
				for _, t := range m.SetTag.Tags {
					if t.ID == *notes[row].TagID {
						tag = t
					}
				}
				if tag.ID != -1 {
					str += fmt.Sprintf(" [%s]", tag.Name)
				}
			}
		}
		str += "\n"
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
	notes := m.getActiveNotes()
	if m.cursor == len(notes)-1 {
		return
	} else if m.cursor == m.scrollIndex+m.maxViewNotes-1 {
		m.cursor++
		m.scrollIndex++
	} else {
		m.cursor++
	}
}

func (m *List) constrainCursor() {
	notes := m.getActiveNotes()
	if m.cursor < 0 {
		m.cursor = 0
	}

	if m.cursor >= len(notes) {
		m.cursor = len(notes) - 1
	}
}

func (m List) getActiveNotes() []model.Note {
	var notes []model.Note
	if m.currentFilter != nil {
		for _, note := range m.Notes {
			if note.TagID != nil {
				if *note.TagID == m.currentFilter.ID {
					notes = append(notes, note)
				}
			}
		}
	} else {
		notes = m.Notes
	}

	return notes
}
