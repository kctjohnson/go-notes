package utils

import (
	"go-notes/internal/db/model"
	"go-notes/internal/graphql"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type FailedToLoadNotesMsg error
type LoadedNotesMsg []model.Note

func LoadNotesCmd(gqlClient *graphql.Client) tea.Cmd {
	return func() tea.Msg {
		notes, err := gqlClient.GetNotes()
		if err != nil {
			return FailedToLoadNotesMsg(err)
		} else {
			return LoadedNotesMsg(notes)
		}
	}
}

type FailedToCreateNoteMsg error
type CreatedNoteMsg model.Note

func CreateNoteCmd(gqlClient *graphql.Client, title string) tea.Cmd {
	return func() tea.Msg {
		note, err := gqlClient.CreateNote("This is the title!")
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

func EditNoteCmd(gqlClient *graphql.Client, id int64) tea.Cmd {
	return func() tea.Msg {
		note, err := gqlClient.GetNote(id)
		if err != nil {
			return FailedToEditNoteMsg(err)
		}
		return EditNoteMsg(note)
	}
}

type FailedToSaveNoteMsg error
type SaveNoteMsg model.Note

func SaveNoteCmd(gqlClient *graphql.Client, note model.Note) tea.Cmd {
	return func() tea.Msg {
		updatedNote, err := gqlClient.SaveNote(note)
		if err != nil {
			return FailedToSaveNoteMsg(err)
		}
		return SaveNoteMsg(updatedNote)
	}
}

type FailedToDeleteNoteMsg error
type DeletedNoteMsg int64

func DeleteNoteCmd(gqlClient *graphql.Client, id int64) tea.Cmd {
	return func() tea.Msg {
		err := gqlClient.DeleteNote(id)
		if err != nil {
			return FailedToDeleteNoteMsg(err)
		}
		return DeletedNoteMsg(id)
	}
}
