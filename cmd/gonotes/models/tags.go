package models

import (
	"fmt"
	"go-notes/cmd/gonotes/utils"
	"go-notes/internal/db/model"
	"go-notes/internal/graphql"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TagsState int

const (
	TAG_LIST TagsState = iota
	TAG_NEW
)

type Tags struct {
	tags   []model.Tag
	width  int
	height int
	keys   tagsKeymap
	help   help.Model
	state  TagsState

	activeFilterTag int
	cursor          int
	scrollIndex     int
	maxViewTags     int

	input textinput.Model

	gqlClient *graphql.Client
}

func NewTags(keys tagsKeymap, gqlClient *graphql.Client) *Tags {
	maxHeight := 15
	width := 20

	textInput := textinput.New()
	textInput.Placeholder = "New tag name..."
	textInput.CharLimit = 20
	textInput.Width = width

	return &Tags{
		tags:            []model.Tag{{ID: -1, Name: "All"}},
		width:           width,
		height:          maxHeight,
		keys:            keys,
		help:            help.New(),
		state:           TAG_LIST,
		activeFilterTag: 0,
		cursor:          0,
		scrollIndex:     0,
		maxViewTags:     maxHeight,
		input:           textInput,
		gqlClient:       gqlClient,
	}
}

func (m Tags) Init() tea.Cmd {
	return nil
}

func (m Tags) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Constrain the cursor to prevent it from going out of bounds
	m.constrainCursor()

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.help.Width = msg.Width
	case tea.KeyMsg:
		if m.state == TAG_LIST {
			switch {
			case key.Matches(msg, m.keys.Up):
				m.cursorUp()
			case key.Matches(msg, m.keys.Down):
				m.cursorDown()
			case key.Matches(msg, m.keys.Toggle):
				if m.cursor == m.activeFilterTag {
					m.activeFilterTag = 0
				} else {
					m.activeFilterTag = m.cursor
				}
			case key.Matches(msg, m.keys.New):
				m.state = TAG_NEW
				return m, utils.CreateTagCmd(m.gqlClient, "New Tag")
			case key.Matches(msg, m.keys.Delete):
				if m.cursor != 0 {
					toDelete := m.tags[m.cursor].ID
					// Set the active filter back to all if they delete the active filter
					if m.cursor == m.activeFilterTag {
						m.activeFilterTag = 0
					}
					return m, utils.DeleteTagCmd(m.gqlClient, toDelete)
				}
			case key.Matches(msg, m.keys.Help):
				m.help.ShowAll = !m.help.ShowAll
			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit
			}
		} else if m.state == TAG_NEW {
			switch {
			}
		}
	}
	return m, nil
}

func (m Tags) View() string {
	str := utils.TitleStyle.Render("Tags:") + "\n"
	for row := m.scrollIndex; row < len(m.tags) && row < m.scrollIndex+m.maxViewTags; row++ {
		line := m.tags[row].Name
		if row == m.activeFilterTag && row == m.cursor {
			line = utils.FocusedLineStyle.Copy().Inherit(utils.UnderlinedStyle).Render(line)
		} else if row == m.activeFilterTag {
			line = utils.FocusedLineStyle.Render(line)
		} else if row == m.cursor {
			line = utils.UnderlinedStyle.Render(line)
		}

		if row == m.activeFilterTag {
			str += fmt.Sprintf("%s\n", line)
		} else {
			str += fmt.Sprintf("%s\n", line)
		}
	}
	str += "\n" + m.help.View(m.keys)
	return str
}

func (m *Tags) cursorUp() {
	if m.cursor == 0 {
		return
	} else if m.cursor == m.scrollIndex {
		m.cursor--
		m.scrollIndex--
	} else {
		m.cursor--
	}
}

func (m *Tags) cursorDown() {
	if m.cursor == len(m.tags)-1 {
		return
	} else if m.cursor == m.scrollIndex+m.maxViewTags-1 {
		m.cursor++
		m.scrollIndex++
	} else {
		m.cursor++
	}
}

func (m *Tags) constrainCursor() {
	if m.cursor < 0 {
		m.cursor = 0
	}

	if m.cursor >= len(m.tags) {
		m.cursor = len(m.tags) - 1
	}
}
