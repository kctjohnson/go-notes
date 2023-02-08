package models

import (
	"bufio"
	"fmt"
	"go-notes/cmd/gonotes/utils"
	"go-notes/internal/db/model"
	"go-notes/internal/graphql"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type StateEnum int

const (
	LIST StateEnum = iota
	TAGS
	LOADING
)

type Main struct {
	gqlClient *graphql.Client
	list      List
	tags      Tags
	curState  StateEnum
	prevState StateEnum
}

func NewMain(gqlClient *graphql.Client) *Main {
	return &Main{
		gqlClient: gqlClient,
		list:      *NewList(listKeys, gqlClient),
		tags:      *NewTags(tagsKeys, gqlClient),
		curState:  LOADING,
	}
}

func (m Main) Init() tea.Cmd {
	return tea.Batch(utils.LoadNotesCmd(m.gqlClient), utils.LoadTagsCmd(m.gqlClient))
}

func (m Main) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.curState {
		case LIST:
			switch {
			case key.Matches(msg, m.list.keys.Tags):
				m.curState = TAGS
				return m, nil
			}
		case TAGS:
			switch {
			case key.Matches(msg, m.tags.keys.Back):
				m.curState = LIST
				return m, nil
			}
		}
	case utils.FailedToLoadNotesMsg:
		log.Fatalf("Failed to load notes!\nError: %v\n", msg)
		return m, tea.Quit
	case utils.LoadedNotesMsg:
		m.list.notes = msg
		m.curState = LIST
	case utils.FailedToCreateNoteMsg:
		log.Fatalf("Failed to create note!\nError: %v\n", msg)
		return m, tea.Quit
	case utils.CreatedNoteMsg:
		m.curState = LOADING
		cmds = tea.Batch(cmds, utils.LoadNotesCmd(m.gqlClient))
	case utils.FailedToEditNoteMsg:
		log.Fatalf("Failed to edit note!\nError: %v\n", msg)
		return m, tea.Quit
	case utils.EditNoteMsg:
		cmd := m.editNote(model.Note(msg))
		return m, cmd
	case utils.FailedToSaveNoteMsg:
		log.Fatalf("Failed to save note!\nError: %v\n", msg)
		return m, tea.Quit
	case utils.SaveEditsMsg:
		if msg.Err != nil {
			log.Fatalf("Failed to process note!\nError: %v\n", msg.Err)
			return m, tea.Quit
		}
		return m, m.saveNote(msg.Note, msg.F)
	case utils.SaveNoteMsg:
		m.curState = LOADING
		cmds = tea.Batch(cmds, utils.LoadNotesCmd(m.gqlClient))
	case utils.FailedToDeleteNoteMsg:
		log.Fatalf("Failed to delete note!\nError: %v\n", msg)
		return m, tea.Quit
	case utils.DeletedNoteMsg:
		m.curState = LOADING
		cmds = tea.Batch(cmds, utils.LoadNotesCmd(m.gqlClient))
	case utils.FailedToLoadTagsMsg:
		log.Fatalf("Failed to load tags!\nError: %v\n", msg)
		return m, tea.Quit
	case utils.LoadedTagsMsg:
		m.tags.tags = append([]model.Tag{{ID: -1, Name: "All"}}, msg...)
		m.curState = TAGS
	case utils.FailedToCreateTagMsg:
		log.Fatalf("Failed to create tag!\nError: %v\n", msg)
		return m, tea.Quit
	case utils.CreatedTagMsg:
		m.curState = LOADING
		return m, tea.Batch(cmds, utils.LoadTagsCmd(m.gqlClient))
	case utils.FailedToDeleteTagMsg:
		log.Fatalf("Failed to delete tag!\nError: %v\n", msg)
		return m, tea.Quit
	case utils.DeletedTagMsg:
		m.curState = LOADING
		return m, tea.Batch(cmds, utils.LoadTagsCmd(m.gqlClient))
	}

	// Update all of the delegated models
	temp, cmd := m.modelUpdate(msg)
	m = temp.(Main)
	cmds = tea.Batch(cmds, cmd)

	return m, cmds
}

func (m Main) View() string {
	switch m.curState {
	case LOADING:
		return "Loading the notes..."
	case LIST:
		return m.list.View()
	case TAGS:
		return m.tags.View()
	default:
		return "I'm not sure what else we should be doing, but here we are!"
	}
}

// Sets the current and previous state of the model
func (m *Main) setState(state StateEnum) {
	if m.curState != state {
		m.prevState = m.curState
		m.curState = state
	}
}

// Keybindings that are local to each model
func (m Main) modelUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Delegate the window size messages
	if _, ok := msg.(tea.WindowSizeMsg); ok {
		listModel, listCmd := m.list.Update(msg)
		m.list = listModel.(List)
		tagsModel, tagsCmd := m.tags.Update(msg)
		m.tags = tagsModel.(Tags)
		return m, tea.Batch(listCmd, tagsCmd)
	}

	switch m.curState {
	case LIST:
		listModel, cmd := m.list.Update(msg)
		m.list = listModel.(List)
		return m, cmd
	case TAGS:
		tagsModel, cmd := m.tags.Update(msg)
		m.tags = tagsModel.(Tags)
		return m, cmd
	}
	return m, nil
}

func (m *Main) editNote(note model.Note) tea.Cmd {
	// Create the file header
	titleHeader := fmt.Sprintf("%s\n##########\n", note.Title)

	// Create a string holding the entire file
	fileContent := titleHeader + note.Content

	// Open the note in a temp file for editing
	f, err := os.CreateTemp("", "."+note.Title+"_*.txt")
	if err != nil {
		panic(err)
	}

	// Load the current note into the temp file
	f.WriteString(fileContent)
	f.Seek(0, 0)

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}

	cmd := exec.Command(editor, f.Name())
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return utils.SaveEditsMsg{
			Note: note,
			F:    f,
			Err:  err,
		}
	})
}

func (m *Main) saveNote(note model.Note, f *os.File) tea.Cmd {
	scanner := bufio.NewScanner(f)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	note.Title = lines[0]
	note.Content = strings.Join(lines[2:], "\n")
	note.LastEditedDate = time.Now()

	return utils.SaveNoteCmd(m.gqlClient, note)
}
