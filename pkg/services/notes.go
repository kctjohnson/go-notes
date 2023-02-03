package services

import (
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

func (s *NotesService) GetNotes() ([]model.Note, error) {
	return s.NotesRepository.GetNotes()
}

func (s *NotesService) GetNote(id int64) (model.Note, error) {
	return s.NotesRepository.GetNote(id)
}

func (s *NotesService) CreateNote(title string) (model.Note, error) {
	return s.NotesRepository.CreateNote(title)
}

func (s *NotesService) SaveNote(note model.Note) (model.Note, error) {
	return s.NotesRepository.SaveNote(note)
}

func (s *NotesService) DeleteNote(id int64) error {
	return s.NotesRepository.DeleteNote(id)
}
