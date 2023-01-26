package services

import (
	"context"
	"go-notes/pkg/db/model"
	"go-notes/pkg/db/repositories"
)

type NotesService struct {
	NotesRepository *repositories.NotesRepository
}

func NewNotesService(r *repositories.NotesRepository) *NotesService {
	return &NotesService{
		NotesRepository: r,
	}
}

func (s *NotesService) GetNotes(ctx context.Context) ([]model.Note, error) {
	return s.NotesRepository.GetNotes(ctx)
}

func (s *NotesService) GetNote(ctx context.Context, id int64) (model.Note, error) {
	return s.NotesRepository.GetNote(ctx, id)
}

func (s *NotesService) CreateNote(ctx context.Context, title string) (model.Note, error) {
	return s.NotesRepository.CreateNote(ctx, title)
}

func (s *NotesService) SaveNote(ctx context.Context, note model.Note) (model.Note, error) {
	return s.NotesRepository.SaveNote(ctx, note)
}

func (s *NotesService) DeleteNote(ctx context.Context, id int64) error {
	return s.NotesRepository.DeleteNote(ctx, id)
}
