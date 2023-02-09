package utils

import (
	"go-notes/internal/db/model"
	"go-notes/internal/graphql"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type (
	FailedToLoadNotesMsg error
	LoadedNotesMsg       []model.Note
)

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

type (
	FailedToCreateNoteMsg error
	CreatedNoteMsg        model.Note
)

func CreateNoteCmd(gqlClient *graphql.Client, title string, tagID *int64) tea.Cmd {
	return func() tea.Msg {
		note, err := gqlClient.CreateNote("This is the title!", tagID)
		if err != nil {
			return FailedToCreateNoteMsg(err)
		}
		return CreatedNoteMsg(note)
	}
}

type (
	FailedToEditNoteMsg error
	EditNoteMsg         model.Note
	SaveEditsMsg        struct {
		Note model.Note
		F    *os.File
		Err  error
	}
)

func EditNoteCmd(gqlClient *graphql.Client, id int64) tea.Cmd {
	return func() tea.Msg {
		note, err := gqlClient.GetNote(id)
		if err != nil {
			return FailedToEditNoteMsg(err)
		}
		return EditNoteMsg(note)
	}
}

type (
	FailedToSaveNoteMsg error
	SaveNoteMsg         model.Note
)

func SaveNoteCmd(gqlClient *graphql.Client, note model.Note) tea.Cmd {
	return func() tea.Msg {
		updatedNote, err := gqlClient.SaveNote(note)
		if err != nil {
			return FailedToSaveNoteMsg(err)
		}
		return SaveNoteMsg(updatedNote)
	}
}

type (
	FailedToSetNoteTagMsg error
	SetNoteTagMsg         model.Note
)

func SetNoteTagCmd(gqlClient *graphql.Client, noteID, tagID int64) tea.Cmd {
	return func() tea.Msg {
		updatedNote, err := gqlClient.SetNoteTag(noteID, tagID)
		if err != nil {
			return FailedToSetNoteTagMsg(err)
		}
		return SetNoteTagMsg(updatedNote)
	}
}

type (
	FailedToDeleteNoteMsg error
	DeletedNoteMsg        int64
)

func DeleteNoteCmd(gqlClient *graphql.Client, id int64) tea.Cmd {
	return func() tea.Msg {
		err := gqlClient.DeleteNote(id)
		if err != nil {
			return FailedToDeleteNoteMsg(err)
		}
		return DeletedNoteMsg(id)
	}
}

type (
	FailedToLoadTagsMsg error
	LoadedTagsMsg       []model.Tag
)

func LoadTagsCmd(gqlClient *graphql.Client) tea.Cmd {
	return func() tea.Msg {
		tags, err := gqlClient.GetTags()
		if err != nil {
			return FailedToLoadTagsMsg(err)
		} else {
			return LoadedTagsMsg(tags)
		}
	}
}

type (
	FailedToLoadSetTagsMsg error
	LoadedSetTagsMsg       []model.Tag
)

func LoadSetTagsCmd(gqlClient *graphql.Client) tea.Cmd {
	return func() tea.Msg {
		tags, err := gqlClient.GetTags()
		if err != nil {
			return FailedToLoadSetTagsMsg(err)
		} else {
			return LoadedSetTagsMsg(tags)
		}
	}
}

type (
	FailedToCreateTagMsg error
	CreatedTagMsg        model.Tag
)

func CreateTagCmd(gqlClient *graphql.Client, name string) tea.Cmd {
	return func() tea.Msg {
		tag, err := gqlClient.CreateTag(name)
		if err != nil {
			return FailedToCreateTagMsg(err)
		} else {
			return CreatedTagMsg(tag)
		}
	}
}

type (
	FailedToDeleteTagMsg error
	DeletedTagMsg        int64
)

func DeleteTagCmd(gqlClient *graphql.Client, id int64) tea.Cmd {
	return func() tea.Msg {
		err := gqlClient.DeleteTag(id)
		if err != nil {
			return FailedToDeleteTagMsg(err)
		} else {
			return DeletedTagMsg(id)
		}
	}
}
