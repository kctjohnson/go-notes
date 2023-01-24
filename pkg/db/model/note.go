package model

import "time"

type Note struct {
	ID             int64     `db:"id"`
	Title          string    `db:"title"`
	Content        string    `db:"content"`
	CreatedDate    time.Time `db:"created_date"`
	LastEditedDate time.Time `db:"last_edited_date"`
}
