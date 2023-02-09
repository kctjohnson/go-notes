package services

import (
	"go-notes/internal/db/model"
	"go-notes/internal/db/repositories"
)

type NotesService struct {
	NotesRepository *repositories.NotesRepository
	TagsRepository  *repositories.TagsRepository
}

func NewNotesService(nr *repositories.NotesRepository, tr *repositories.TagsRepository) *NotesService {
	return &NotesService{
		NotesRepository: nr,
		TagsRepository:  tr,
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

func (s *NotesService) GetTags() ([]model.Tag, error) {
	return s.TagsRepository.GetTags()
}

func (s *NotesService) GetTag(id int64) (model.Tag, error) {
	return s.TagsRepository.GetTag(id)
}

func (s *NotesService) CreateTag(name string) (model.Tag, error) {
	return s.TagsRepository.CreateTag(name)
}

func (s *NotesService) DeleteTag(id int64) error {
	// Clear out the tag id on all notes that currently are under that tag
	s.NotesRepository.ClearTagFromNotes(id)
	return s.TagsRepository.DeleteTag(id)
}
