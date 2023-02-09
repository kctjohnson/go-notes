package repositories

import (
	"go-notes/internal/db/model"

	"github.com/jmoiron/sqlx"
)

type TagsRepository struct {
	db *sqlx.DB
}

func NewTagsRepository(db *sqlx.DB) *TagsRepository {
	return &TagsRepository{
		db: db,
	}
}

func (r *TagsRepository) GetTags() ([]model.Tag, error) {
	const getTagsSql = `SELECT * FROM tags`
	rows, err := r.db.Queryx(getTagsSql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []model.Tag{}
	for rows.Next() {
		var tag model.Tag
		err = rows.StructScan(&tag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *TagsRepository) GetTag(id int64) (model.Tag, error) {
	const getTagSql = `SELECT * FROM tags WHERE id=?`
	row := r.db.QueryRowx(getTagSql, id)
	var tag model.Tag
	err := row.StructScan(&tag)
	if err != nil {
		return tag, err
	}
	return tag, nil
}

func (r *TagsRepository) CreateTag(name string) (model.Tag, error) {
	const createTagSql = `INSERT INTO tags (name) VALUES (:name)`

	newTag := model.Tag{Name: name}
	res, err := r.db.NamedExec(createTagSql, newTag)
	if err != nil {
		return model.Tag{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return model.Tag{}, err
	}

	tag, err := r.GetTag(id)
	if err != nil {
		return tag, err
	}

	return tag, nil
}

func (r *TagsRepository) DeleteTag(id int64) error {
	const deleteTagSql = `DELETE FROM tags WHERE id=?`
	_, err := r.db.Exec(deleteTagSql, id)
	return err
}
