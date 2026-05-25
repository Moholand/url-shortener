package service

import (
	"context"
	"time"
	"url-shortener/internal/model"
	"url-shortener/internal/repository"
	"github.com/redis/go-redis/v9"
)

type URLService struct {
	Repo *repository.URLRepository
	Redis *redis.Client
}

func NewURLService(repo *repository.URLRepository, rdb *redis.Client) *URLService {
	return &URLService{
		Repo: repo, 
		Redis: rdb,
	}
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

func (s *URLService) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
    val, err := s.Redis.Get(ctx, shortCode).Result()
    if err == nil {
        return val, nil
    }

    url, err := s.Repo.GetByShortCode(shortCode)
    if err != nil {
        return "", err
    }

    s.Redis.Set(ctx, shortCode, url, 1*time.Hour)
    
    return url, nil
}
