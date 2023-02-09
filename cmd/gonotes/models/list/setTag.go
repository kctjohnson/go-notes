package list

import (
	"fmt"
	"go-notes/cmd/gonotes/utils"
	"go-notes/internal/db/model"
	"go-notes/internal/graphql"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type SetTag struct {
	Tags      []model.Tag
	gqlClient *graphql.Client
	help      help.Model

	cursor      int
	scrollIndex int
	maxViewTags int
}

func NewSetTag(gqlClient *graphql.Client, tags []model.Tag) *SetTag {
	return &SetTag{
		Tags:        tags,
		gqlClient:   gqlClient,
		help:        help.New(),
		cursor:      0,
		scrollIndex: 0,
		maxViewTags: 15,
	}
}

func (m SetTag) Init() tea.Cmd {
	return nil
}

func (m SetTag) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, SetTagKeys.Up):
			m.cursorUp()
		case key.Matches(msg, SetTagKeys.Down):
			m.cursorDown()
		case key.Matches(msg, SetTagKeys.Select):
		case key.Matches(msg, SetTagKeys.Quit):
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m SetTag) View() string {
	str := utils.TitleStyle.Render("Select Tag For Note:") + "\n"
	for row := m.scrollIndex; row < len(m.Tags) && row < m.scrollIndex+m.maxViewTags; row++ {
		if row == m.cursor {
			str += fmt.Sprintf("%s\n", utils.FocusedLineStyle.Render(m.Tags[row].Name))
		} else {
			str += fmt.Sprintf("%s\n", m.Tags[row].Name)
		}
	}
	str += "\n" + m.help.View(Keys)
	return str
}

func (m *SetTag) cursorUp() {
	if m.cursor == 0 {
		return
	} else if m.cursor == m.scrollIndex {
		m.cursor--
		m.scrollIndex--
	} else {
		m.cursor--
	}
}

func (m *SetTag) cursorDown() {
	if m.cursor == len(m.Tags)-1 {
		return
	} else if m.cursor == m.scrollIndex+m.maxViewTags-1 {
		m.cursor++
		m.scrollIndex++
	} else {
		m.cursor++
	}
}

func (m *SetTag) constrainCursor() {
	if m.cursor < 0 {
		m.cursor = 0
	}

	if m.cursor >= len(m.Tags) {
		m.cursor = len(m.Tags) - 1
	}
}
