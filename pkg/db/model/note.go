package model

import "time"

type Note struct {
	ID             int64     `db:"id" graphql:"id"`
	Title          string    `db:"title" graphql:"title"`
	Content        string    `db:"content" graphql:"content"`
	CreatedDate    time.Time `db:"created_date" graphql:"created_date"`
	LastEditedDate time.Time `db:"last_edited_date" graphql:"last_edited_date"`
}
