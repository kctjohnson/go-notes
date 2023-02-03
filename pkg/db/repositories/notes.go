package repositories

import (
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

func (r *NotesRepository) GetNotes() ([]model.Note, error) {
	const getNotesSql = `SELECT * FROM notes`
	rows, err := r.db.Queryx(getNotesSql)
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

func (r *NotesRepository) GetNote(id int64) (model.Note, error) {
	const getNoteSql = `SELECT * FROM notes WHERE id=?`
	row := r.db.QueryRowx(getNoteSql, id)
	var note model.Note
	err := row.StructScan(&note)
	if err != nil {
		return note, err
	}
	return note, nil
}

func (r *NotesRepository) CreateNote(title string) (model.Note, error) {
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

	res, err := r.db.NamedExec(createNoteSql, newNote)
	if err != nil {
		return model.Note{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return model.Note{}, err
	}

	note, err := r.GetNote(id)
	if err != nil {
		return note, err
	}

	return note, nil
}

func (r *NotesRepository) SaveNote(note model.Note) (model.Note, error) {
	const saveNoteSql = `
		UPDATE notes
		SET title=:title, content=:content, last_edited_date=:last_edited_date
		WHERE id=:id
	`

	note.LastEditedDate = time.Now()

	_, err := r.db.NamedExec(saveNoteSql, note)
	if err != nil {
		return note, err
	}

	newNote, err := r.GetNote(note.ID)
	if err != nil {
		return newNote, err
	}

	return newNote, nil
}

func (r *NotesRepository) DeleteNote(id int64) error {
	const deleteNoteSql = `DELETE FROM notes WHERE id=?`
	_, err := r.db.Exec(deleteNoteSql, id)
	return err
}
