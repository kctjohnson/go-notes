package models

import (
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
		list:        *NewList(listKeys, ns),
		edit:        Edit{},
		curFocus:    LOADING,
	}
}

func (m Main) Init() tea.Cmd {
	return utils.LoadNotesCmd(m.noteService)
}

func (m Main) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds tea.Cmd
	switch msg := msg.(type) {
	case utils.FailedToLoadNotesMsg:
		log.Fatalf("Failed to load notes!\nError: %v\n", msg)
		return m, tea.Quit
	case utils.LoadedNotesMsg:
		m.list.notes = msg
		m.curFocus = LIST
	case utils.FailedToCreateNoteMsg:
		log.Fatalf("Failed to create note!\nError: %v\n", msg)
		return m, tea.Quit
	case utils.CreatedNoteMsg:
		m.curFocus = LOADING
		cmds = tea.Batch(cmds, utils.LoadNotesCmd(m.noteService))
	}

	// Update all of the delegated models
	temp, cmd := m.modelUpdate(msg)
	m = temp.(Main)
	cmds = tea.Batch(cmds, cmd)

	return m, cmds
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

// Keybindings that are local to each model
func (m Main) modelUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.curFocus {
	case LIST:
		listModel, cmd := m.list.Update(msg)
		m.list = listModel.(List)
		return m, cmd
	case EDIT:
		editModel, cmd := m.edit.Update(msg)
		m.edit = editModel.(Edit)
		return m, cmd
	}
	return m, nil
}
