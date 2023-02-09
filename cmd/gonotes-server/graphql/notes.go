package graphql

import (
	input "go-notes/cmd/gonotes-server/graphql/input/notes"
	"go-notes/internal/db/model"
	"go-notes/internal/services"

	"github.com/samsarahq/thunder/graphql/schemabuilder"
)

type NotesGql struct {
	NotesService *services.NotesService
}

func (g *NotesGql) registerNote(schema *schemabuilder.Schema) {
	schema.Object("Note", model.Note{})
}

func (g *NotesGql) registerQuery(querySchemaObj *schemabuilder.Object) {
	querySchemaObj.FieldFunc("Notes", func() ([]model.Note, error) {
		return g.NotesService.GetNotes()
	})

	querySchemaObj.FieldFunc("Note", func(args input.GetNote) (model.Note, error) {
		return g.NotesService.GetNote(args.NoteID)
	})
}

func (g *NotesGql) registerMutation(mutationSchemaObj *schemabuilder.Object) {
	mutationSchemaObj.FieldFunc("CreateNote", func(args input.CreateNote) (model.Note, error) {
		return g.NotesService.CreateNote(args.Title, args.TagID)
	})

	mutationSchemaObj.FieldFunc("SaveNote", func(args input.SaveNote) (model.Note, error) {
		return g.NotesService.SaveNote(args.Note)
	})

	mutationSchemaObj.FieldFunc("SetNoteTag", func(args input.SetNoteTag) (model.Note, error) {
		return g.NotesService.SetNoteTag(args.NoteID, args.TagID)
	})

	mutationSchemaObj.FieldFunc("DeleteNote", func(args input.DeleteNote) error {
		return g.NotesService.DeleteNote(args.ID)
	})
}
