package service

import (
	"context"
	"errors"
	"time"
	"url-shortener/internal/model"
	"url-shortener/internal/repository"
	"url-shortener/pkg/code"
	"github.com/lib/pq"
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

func (s *URLService) Create(ctx context.Context, originalURL string) (*model.URL, error) {
	for i := 0; i < 3; i++ {
		shortCode := code.Generate(6)

		urlData := &model.URL{
			ShortCode:   shortCode,
			OriginalURL: originalURL,
		}

		err := s.Repo.Save(urlData)
		if err == nil {
			s.Redis.Set(ctx, shortCode, originalURL, 24*time.Hour)
			return urlData, nil
		}

		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			continue
		}

		return nil, err
	}

	return nil, errors.New("could not generate a unique short code after 3 attempts")
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
