package models

import (
	"bufio"
	"fmt"
	"go-notes/cmd/cli/utils"
	"go-notes/pkg/db/model"
	"go-notes/pkg/services"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type FocusEnum int

const (
	LIST FocusEnum = iota
	LOADING
)

type Main struct {
	noteService *services.NotesService
	list        List
	curFocus    FocusEnum
}

func NewMain(ns *services.NotesService) *Main {
	return &Main{
		noteService: ns,
		list:        *NewList(listKeys, ns),
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
		m.curFocus = LOADING
		cmds = tea.Batch(cmds, utils.LoadNotesCmd(m.noteService))
	case utils.FailedToDeleteNoteMsg:
		log.Fatalf("Failed to delete note!\nError: %v\n", msg)
		return m, tea.Quit
	case utils.DeletedNoteMsg:
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
	default:
		return "I'm not sure what else we should be doing, but here we are!"
	}
}

// Keybindings that are local to each model
func (m Main) modelUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Delegate the window size messages
	if _, ok := msg.(tea.WindowSizeMsg); ok {
		listModel, cmd := m.list.Update(msg)
		m.list = listModel.(List)
		return m, cmd
	}

	switch m.curFocus {
	case LIST:
		listModel, cmd := m.list.Update(msg)
		m.list = listModel.(List)
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

	return utils.SaveNoteCmd(m.noteService, note)
}
