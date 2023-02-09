package input

type CreateNote struct {
	Title string `graphql:"title"`
	TagID *int64 `graphql:"tag_id"`
}
