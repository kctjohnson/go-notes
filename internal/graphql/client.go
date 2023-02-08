package graphql

import (
	"context"
	"go-notes/internal/db/model"

	"github.com/Khan/genqlient/graphql"
)

type Client struct {
	client graphql.Client
}

func NewClient(endpoint string, httpclient graphql.Doer) *Client {
	return &Client{
		client: graphql.NewClient(endpoint, httpclient),
	}
}

func (c *Client) GetTags() ([]model.Tag, error) {
	ctx := context.Background()
	resp, err := GetTags(ctx, c.client)
	if err != nil {
		return nil, err
	}

	// Convert it to a model note slice
	tags := []model.Tag{}
	for _, n := range resp.Tags {
		tag := model.Tag{
			ID:   n.Id,
			Name: n.Name,
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (c *Client) GetTag(id int64) (model.Tag, error) {
	ctx := context.Background()
	resp, err := GetTag(ctx, c.client, id)
	if err != nil {
		return model.Tag{}, err
	}

	// Convert it to a model tag
	tag := model.Tag{
		ID:   resp.Tag.Id,
		Name: resp.Tag.Name,
	}

	return tag, nil
}

func (c *Client) CreateTag(name string) (model.Tag, error) {
	ctx := context.Background()
	resp, err := CreateTag(ctx, c.client, name)
	if err != nil {
		return model.Tag{}, err
	}

	// Convert it to a model note
	note := model.Tag{
		ID:   resp.CreateTag.Id,
		Name: resp.CreateTag.Name,
	}

	return note, nil
}

func (c *Client) DeleteTag(id int64) error {
	ctx := context.Background()
	_, err := DeleteTag(ctx, c.client, id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetNotes() ([]model.Note, error) {
	ctx := context.Background()
	resp, err := GetNotes(ctx, c.client)
	if err != nil {
		return nil, err
	}

	// Convert it to a model note slice
	notes := []model.Note{}
	for _, n := range resp.Notes {
		note := model.Note{
			ID:             n.Id,
			Title:          n.Title,
			Content:        n.Content,
			CreatedDate:    n.Created_date,
			LastEditedDate: n.Last_edited_date,
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (c *Client) GetNote(id int64) (model.Note, error) {
	ctx := context.Background()
	resp, err := GetNote(ctx, c.client, id)
	if err != nil {
		return model.Note{}, err
	}

	// Convert it to a model note
	note := model.Note{
		ID:             resp.Note.Id,
		Title:          resp.Note.Title,
		Content:        resp.Note.Content,
		CreatedDate:    resp.Note.Created_date,
		LastEditedDate: resp.Note.Last_edited_date,
	}

	return note, nil
}

func (c *Client) CreateNote(title string) (model.Note, error) {
	ctx := context.Background()
	resp, err := CreateNote(ctx, c.client, title)
	if err != nil {
		return model.Note{}, err
	}

	// Convert it to a model note
	note := model.Note{
		ID:             resp.CreateNote.Id,
		Title:          resp.CreateNote.Title,
		Content:        resp.CreateNote.Content,
		CreatedDate:    resp.CreateNote.Created_date,
		LastEditedDate: resp.CreateNote.Last_edited_date,
	}

	return note, nil
}

func (c *Client) SaveNote(note model.Note) (model.Note, error) {
	ctx := context.Background()
	resp, err := SaveNote(ctx, c.client, Note_InputObject{
		Id:               note.ID,
		Title:            note.Title,
		Content:          note.Content,
		Created_date:     note.CreatedDate,
		Last_edited_date: note.LastEditedDate,
	})
	if err != nil {
		return model.Note{}, err
	}

	// Convert it to a model note
	cnvNote := model.Note{
		ID:             resp.SaveNote.Id,
		Title:          resp.SaveNote.Title,
		Content:        resp.SaveNote.Content,
		CreatedDate:    resp.SaveNote.Created_date,
		LastEditedDate: resp.SaveNote.Last_edited_date,
	}

	return cnvNote, nil
}

func (c *Client) DeleteNote(id int64) error {
	ctx := context.Background()
	_, err := DeleteNote(ctx, c.client, id)
	if err != nil {
		return err
	}
	return nil
}
