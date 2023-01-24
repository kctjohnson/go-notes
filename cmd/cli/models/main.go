package models

import (
	"context"
	"go-notes/cmd/cli/utils"
	"go-notes/pkg/services"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type FocusEnum int

const (
	LIST FocusEnum = iota
	EDIT
	LOADING
)

type Main struct {
	noteService *services.NotesService
	list        List
	edit        Edit
	curFocus    FocusEnum
}

func NewMain(ns *services.NotesService) *Main {
	return &Main{
		noteService: ns,
		list:        List{},
		edit:        Edit{},
		curFocus:    LOADING,
	}
}

func (m Main) Init() tea.Cmd {
	return utils.LoadNotesCmd(m.noteService)
}

func (m Main) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Global keybindings
	switch msg := msg.(type) {
	case utils.FailedToLoadNotesMsg:
		log.Fatalf("Failed to load notes!\nError: %v\n", msg)
		return m, tea.Quit
	case utils.LoadedNotesMsg:
		m.list.notes = msg
		m.curFocus = LIST
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			_, err := m.noteService.CreateNote(context.Background(), "This is the title!")
			if err != nil {
				log.Fatalf("Error: %v\n", err)
				return m, tea.Quit
			}
			m.curFocus = LOADING
			return m, utils.LoadNotesCmd(m.noteService)
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	// Switch through context specific bindings
	// Message handlers
	return m, nil
}

func (m Main) View() string {
	switch m.curFocus {
	case LOADING:
		return "Loading the notes..."
	case LIST:
		return m.list.View()
	case EDIT:
		return m.edit.View()
	default:
		return "I'm not sure what else we should be doing, but here we are!"
	}
}
