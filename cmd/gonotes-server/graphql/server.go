package graphql

import (
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
)

type GQLServer struct {
	NotesGql *NotesGql
	TagsGql  *TagsGql
}

func (g *GQLServer) registerQuery(schema *schemabuilder.Schema) {
	obj := schema.Query()
	g.NotesGql.registerQuery(obj)
	g.TagsGql.registerQuery(obj)
}

func (g *GQLServer) registerMutation(schema *schemabuilder.Schema) {
	obj := schema.Mutation()
	g.NotesGql.registerMutation(obj)
	g.TagsGql.registerMutation(obj)
}

func (g *GQLServer) registerStructs(schema *schemabuilder.Schema) {
	g.NotesGql.registerNote(schema)
	g.TagsGql.registerTag(schema)
}

func (g *GQLServer) Schema() *graphql.Schema {
	builder := schemabuilder.NewSchema()
	g.registerQuery(builder)
	g.registerMutation(builder)
	g.registerStructs(builder)
	return builder.MustBuild()
}
