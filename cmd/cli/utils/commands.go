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
	notes, err := ns.GetNotes(context.Background())
	return func() tea.Msg {
		if err != nil {
			return FailedToLoadNotesMsg(err)
		} else {
			return LoadedNotesMsg(notes)
		}
	}
}
