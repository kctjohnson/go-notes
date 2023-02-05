package graphql

import (
	"go-notes/cmd/gonotes-server/graphql/input"
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
	querySchemaObj.FieldFunc("notes", func() ([]model.Note, error) {
		return g.NotesService.GetNotes()
	})

	querySchemaObj.FieldFunc("note", func(args input.GetNotes) (model.Note, error) {
		return g.NotesService.GetNote(args.NoteID)
	})
}

func (g *NotesGql) registerMutation(mutationSchemaObj *schemabuilder.Object) {
	mutationSchemaObj.FieldFunc("createNote", func(args input.CreateNote) (model.Note, error) {
		return g.NotesService.CreateNote(args.Title)
	})

	mutationSchemaObj.FieldFunc("saveNote", func(args input.SaveNote) (model.Note, error) {
		return g.NotesService.SaveNote(args.Note)
	})

	mutationSchemaObj.FieldFunc("deleteNote", func(args input.DeleteNote) error {
		return g.NotesService.DeleteNote(args.ID)
	})
}
