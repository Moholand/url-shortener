package service

import (
	"context"
	"errors"
	"strings"
	"time"
	"url-shortener/internal/model"
	"url-shortener/internal/repository"
	"url-shortener/pkg/code"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type URLService struct {
	Repo       *repository.URLRepository
	ClickRepo  *repository.ClickRepository
	Redis      *redis.Client
}

func NewURLService(repo *repository.URLRepository, clickRepo *repository.ClickRepository, rdb *redis.Client) *URLService {
	return &URLService{
		Repo:      repo,
		ClickRepo: clickRepo,
		Redis:     rdb,
	}
}

func (s *URLService) Create(ctx context.Context, originalURL string) (*model.URL, error) {
	if !strings.HasPrefix(originalURL, "http://") &&
		!strings.HasPrefix(originalURL, "https://") {
		originalURL = "https://" + originalURL
	}

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

	if url != "" {
		s.Redis.Set(ctx, shortCode, url, 1*time.Hour)
	}

	return url, nil
}

func (s *URLService) RecordClick(ctx context.Context, shortCode, ipAddress, userAgent, referer string) {
	err := s.ClickRepo.Save(shortCode, ipAddress, userAgent, referer)
	if err != nil {
		return
	}
}

func (s *URLService) GetAnalytics(ctx context.Context, shortCode string) (int, []repository.ClickRecord, error) {
	clicks, err := s.ClickRepo.GetByShortCode(shortCode, 50, 0)
	if err != nil {
		return 0, nil, err
	}

	total, err := s.ClickRepo.CountByShortCode(shortCode)
	if err != nil {
		return 0, nil, err
	}

	return total, clicks, nil
}
