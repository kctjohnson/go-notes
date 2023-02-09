package tags

import (
	"fmt"
	"go-notes/cmd/gonotes/utils"
	"go-notes/internal/db/model"
	"go-notes/internal/graphql"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type State int

const (
	LIST State = iota
	NEW
)

type Tags struct {
	Tags   []model.Tag
	width  int
	height int
	help   help.Model
	State  State

	ActiveFilterTag int
	cursor          int
	scrollIndex     int
	maxViewTags     int

	input textinput.Model

	gqlClient *graphql.Client
}

func New(gqlClient *graphql.Client) *Tags {
	maxHeight := 15
	width := 20

	textInput := textinput.New()
	textInput.Placeholder = "New tag name..."
	textInput.CharLimit = 20
	textInput.Width = width

	return &Tags{
		Tags:            []model.Tag{{ID: -1, Name: "All"}},
		width:           width,
		height:          maxHeight,
		help:            help.New(),
		State:           LIST,
		ActiveFilterTag: 0,
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

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.help.Width = msg.Width
	case tea.KeyMsg:
		if m.State == LIST {
			switch {
			case key.Matches(msg, Keys.Up):
				m.cursorUp()
			case key.Matches(msg, Keys.Down):
				m.cursorDown()
			case key.Matches(msg, Keys.Toggle):
				if m.cursor == m.ActiveFilterTag {
					m.ActiveFilterTag = 0
				} else {
					m.ActiveFilterTag = m.cursor
				}
			case key.Matches(msg, Keys.New):
				m.State = NEW
				m.input.Reset()
				m.input.Focus()
				return m, nil
			case key.Matches(msg, Keys.Delete):
				if m.cursor != 0 {
					toDelete := m.Tags[m.cursor].ID
					// Set the active filter back to all if they delete the active filter
					if m.cursor == m.ActiveFilterTag {
						m.ActiveFilterTag = 0
					}
					return m, utils.DeleteTagCmd(m.gqlClient, toDelete)
				}
			case key.Matches(msg, Keys.Help):
				m.help.ShowAll = !m.help.ShowAll
			case key.Matches(msg, Keys.Quit):
				return m, tea.Quit
			}
		} else if m.State == NEW {
			switch {
			case key.Matches(msg, InputKeys.Create):
				m.input.Blur()
				trimmedInput := strings.TrimSpace(m.input.Value())
				m.State = LIST
				return m, utils.CreateTagCmd(m.gqlClient, trimmedInput)
			case key.Matches(msg, InputKeys.Back):
				m.input.Blur()
				m.input.Reset()
				m.State = LIST
			case key.Matches(msg, InputKeys.Quit):
				m.input.Blur()
				return m, tea.Quit
			}
		}
	}

	// Update the input if in state
	if m.State == NEW {
		var c tea.Cmd
		m.input, c = m.input.Update(msg)
		cmd = tea.Batch(cmd, c)
	}

	return m, cmd
}

func (m Tags) View() string {
	switch m.State {
	case LIST:
		return m.viewList()
	case NEW:
		return m.viewInput()
	}
	return "Something went wrong! Press ctrl+c to exit!"
}

func (m Tags) viewList() string {
	str := utils.TitleStyle.Render("Tags:") + "\n"
	for row := m.scrollIndex; row < len(m.Tags) && row < m.scrollIndex+m.maxViewTags; row++ {
		line := m.Tags[row].Name
		if row == m.ActiveFilterTag && row == m.cursor {
			line = utils.FocusedLineStyle.Copy().Inherit(utils.UnderlinedStyle).Render(line)
		} else if row == m.ActiveFilterTag {
			line = utils.FocusedLineStyle.Render(line)
		} else if row == m.cursor {
			line = utils.UnderlinedStyle.Render(line)
		}

		if row == m.ActiveFilterTag {
			str += fmt.Sprintf("%s\n", line)
		} else {
			str += fmt.Sprintf("%s\n", line)
		}
	}
	str += "\n" + m.help.View(Keys)
	return str
}

func (m Tags) viewInput() string {
	str := m.input.View()
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
	if m.cursor == len(m.Tags)-1 {
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

	if m.cursor >= len(m.Tags) {
		m.cursor = len(m.Tags) - 1
	}
}
