package utils

import (
	"context"
	"go-notes/pkg/db/model"
	"go-notes/pkg/services"

	tea "github.com/charmbracelet/bubbletea"
)

type FailedToLoadNotesMsg error
type LoadedNotesMsg []model.Note

func LoadNotesCmd(ns *services.NotesService) tea.Cmd {
	return func() tea.Msg {
		notes, err := ns.GetNotes(context.Background())
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
		note, err := ns.CreateNote(context.Background(), "This is the title!")
		if err != nil {
			return FailedToCreateNoteMsg(err)
		}
		return CreatedNoteMsg(note)
	}
}
