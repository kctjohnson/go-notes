package input

import "go-notes/internal/db/model"

type SaveNote struct {
	Note model.Note `graphql:"note"`
}
