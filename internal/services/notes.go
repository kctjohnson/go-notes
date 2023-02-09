package services

import (
	"go-notes/internal/db/model"
	"go-notes/internal/db/repositories"
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

func (s *NotesService) CreateNote(title string, tagID *int64) (model.Note, error) {
	return s.NotesRepository.CreateNote(title, tagID)
}

func (s *NotesService) SaveNote(note model.Note) (model.Note, error) {
	return s.NotesRepository.SaveNote(note)
}

func (s *NotesService) SetNoteTag(noteID, tagID int64) (model.Note, error) {
	return s.NotesRepository.SetNoteTag(noteID, tagID)
}

func (s *NotesService) DeleteNote(id int64) error {
	return s.NotesRepository.DeleteNote(id)
}
