package graphql

import (
	input "go-notes/cmd/gonotes-server/graphql/input/tags"
	"go-notes/internal/db/model"
	"go-notes/internal/services"

	"github.com/samsarahq/thunder/graphql/schemabuilder"
)

type TagsGql struct {
	NotesService *services.NotesService
}

func (g *TagsGql) registerTag(schema *schemabuilder.Schema) {
	schema.Object("Tag", model.Tag{})
}

func (g *TagsGql) registerQuery(querySchemaObj *schemabuilder.Object) {
	querySchemaObj.FieldFunc("Tags", func() ([]model.Tag, error) {
		return g.NotesService.GetTags()
	})

	querySchemaObj.FieldFunc("Tag", func(args input.GetTag) (model.Tag, error) {
		return g.NotesService.GetTag(args.TagID)
	})
}

func (g *TagsGql) registerMutation(mutationSchemaObj *schemabuilder.Object) {
	mutationSchemaObj.FieldFunc("CreateTag", func(args input.CreateTag) (model.Tag, error) {
		return g.NotesService.CreateTag(args.Name)
	})

	mutationSchemaObj.FieldFunc("DeleteTag", func(args input.DeleteTag) error {
		return g.NotesService.DeleteTag(args.ID)
	})
}
