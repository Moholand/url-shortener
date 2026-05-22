package service

import (
	"url-shortener/internal/model"
	"url-shortener/internal/repository"
)

type URLService struct {
	Repo *repository.URLRepository
}

func NewURLService(repo *repository.URLRepository) *URLService {
	return &URLService{Repo: repo}
}

func (s *URLService) Create(shortCode, originalURL string) (*model.URL, error) {

	urlData := &model.URL{
		ShortCode:   shortCode,
		OriginalURL: originalURL,
	}

	err := s.Repo.Save(urlData)
	if err != nil {
		return nil, err
	}

	return urlData, nil
}
