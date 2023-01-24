package repositories

import (
	"context"
	"go-notes/pkg/db/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type NotesRepository struct {
	db *sqlx.DB
}

func NewNotesRepository(db *sqlx.DB) *NotesRepository {
	return &NotesRepository{
		db: db,
	}
}

func (r *NotesRepository) GetNotes(ctx context.Context) ([]model.Note, error) {
	const getNotesSql = `SELECT * FROM notes`
	rows, err := r.db.QueryxContext(ctx, getNotesSql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := []model.Note{}
	for rows.Next() {
		var note model.Note
		err = rows.StructScan(&note)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

func (r *NotesRepository) GetNote(ctx context.Context, id int64) (model.Note, error) {
	const getNoteSql = `SELECT * FROM notes WHERE id=?`
	row := r.db.QueryRowxContext(ctx, getNoteSql, id)
	var note model.Note
	err := row.StructScan(&note)
	if err != nil {
		return note, err
	}
	return note, nil
}

func (r *NotesRepository) CreateNote(ctx context.Context, title string) (model.Note, error) {
	const createNoteSql = `
	INSERT INTO notes (title, content, created_date, last_edited_date)
	VALUES (:title, :content, :created_date, :last_edited_date)
	`

	curDate := time.Now()
	newNote := model.Note{
		Title:          title,
		CreatedDate:    curDate,
		LastEditedDate: curDate,
	}

	res, err := r.db.NamedExecContext(ctx, createNoteSql, newNote)
	if err != nil {
		return model.Note{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return model.Note{}, err
	}

	note, err := r.GetNote(ctx, id)
	if err != nil {
		return note, err
	}

	return note, nil
}

func (r *NotesRepository) DeleteNote(ctx context.Context) error {
	return nil
}
