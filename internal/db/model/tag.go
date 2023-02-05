package model

type Tag struct {
	ID   int64  `db:"id" graphql:"id"`
	Name string `db:"name" graphql:"name"`
}
