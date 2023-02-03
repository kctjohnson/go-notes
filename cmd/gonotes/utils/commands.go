package utils

import (
	"go-notes/pkg/db/model"
	"go-notes/pkg/services"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type FailedToLoadNotesMsg error
type LoadedNotesMsg []model.Note

func LoadNotesCmd(ns *services.NotesService) tea.Cmd {
	return func() tea.Msg {
		notes, err := ns.GetNotes()
		if err != nil {
			return FailedToLoadNotesMsg(err)
		} else {
			return LoadedNotesMsg(notes)
		}
	}
}

type FailedToCreateNoteMsg error
type CreatedNoteMsg model.Note

func CreateNoteCmd(ns *services.NotesService, title string) tea.Cmd {
	return func() tea.Msg {
		note, err := ns.CreateNote("This is the title!")
		if err != nil {
			return FailedToCreateNoteMsg(err)
		}
		return CreatedNoteMsg(note)
	}
}

type FailedToEditNoteMsg error
type EditNoteMsg model.Note
type SaveEditsMsg struct {
	Note model.Note
	F    *os.File
	Err  error
}

func EditNoteCmd(ns *services.NotesService, id int64) tea.Cmd {
	return func() tea.Msg {
		note, err := ns.GetNote(id)
		if err != nil {
			return FailedToEditNoteMsg(err)
		}
		return EditNoteMsg(note)
	}
}

type FailedToSaveNoteMsg error
type SaveNoteMsg model.Note

func SaveNoteCmd(ns *services.NotesService, note model.Note) tea.Cmd {
	return func() tea.Msg {
		updatedNote, err := ns.SaveNote(note)
		if err != nil {
			return FailedToSaveNoteMsg(err)
		}
		return SaveNoteMsg(updatedNote)
	}
}

type FailedToDeleteNoteMsg error
type DeletedNoteMsg int64

func DeleteNoteCmd(ns *services.NotesService, id int64) tea.Cmd {
	return func() tea.Msg {
		err := ns.DeleteNote(id)
		if err != nil {
			return FailedToDeleteNoteMsg(err)
		}
		return DeletedNoteMsg(id)
	}
}
