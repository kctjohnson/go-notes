package graphql

import (
	"go-notes/cmd/gonotes-server/graphql/inputs"
	"go-notes/pkg/db/model"
	"go-notes/pkg/services"

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

	querySchemaObj.FieldFunc("note", func(args inputs.GetNotesInput) (model.Note, error) {
		return g.NotesService.GetNote(args.NoteID)
	})
}

func (g *NotesGql) registerMutation(mutationSchemaObj *schemabuilder.Object) {
	mutationSchemaObj.FieldFunc("createNote", func(args struct{ Title string }) (model.Note, error) {
		return g.NotesService.CreateNote(args.Title)
	})

	mutationSchemaObj.FieldFunc("saveNote", func(args struct{ Note model.Note }) (model.Note, error) {
		return g.NotesService.SaveNote(args.Note)
	})

	mutationSchemaObj.FieldFunc("deleteNote", func(args struct{ ID int64 }) error {
		return g.NotesService.DeleteNote(args.ID)
	})
}
