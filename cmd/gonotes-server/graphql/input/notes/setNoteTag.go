package input

type SetNoteTag struct {
	NoteID int64 `graphql:"note_id"`
	TagID  int64 `graphql:"tag_id"`
}
