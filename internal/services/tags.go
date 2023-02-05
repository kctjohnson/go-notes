package services

import (
	"go-notes/internal/db/model"
	"go-notes/internal/db/repositories"
)

type TagsService struct {
	TagsRepository *repositories.TagsRepository
}

func NewTagsService(r *repositories.TagsRepository) *TagsService {
	return &TagsService{
		TagsRepository: r,
	}
}

func (s *TagsService) GetTags() ([]model.Tag, error) {
	return s.TagsRepository.GetTags()
}

func (s *TagsService) GetTag(id int64) (model.Tag, error) {
	return s.TagsRepository.GetTag(id)
}

func (s *TagsService) CreateTag(name string) (model.Tag, error) {
	return s.TagsRepository.CreateTag(name)
}

func (s *TagsService) DeleteTag(id int64) error {
	return s.TagsRepository.DeleteTag(id)
}
