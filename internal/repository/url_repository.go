package repository

import (
	"database/sql"
	"url-shortener/internal/model"
)

type URLRepository struct {
	DB *sql.DB
}

func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{DB: db}
}

func (r *URLRepository) Save(url *model.URL) error {
	query := `
		INSERT INTO urls (short_code, original_url, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := r.DB.QueryRow(query, url.ShortCode, url.OriginalURL, url.ExpiresAt).Scan(&url.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *URLRepository) GetByShortCode(shortCode string) (string, error) {
	var originalURL string
	query := `SELECT original_url FROM urls WHERE short_code = $1 AND (expires_at IS NULL OR expires_at > NOW())`
	
	err := r.DB.QueryRow(query, shortCode).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // Not Found!
		}
		return "", err
	}
	return originalURL, nil
}

func (r *URLRepository) GetURLByShortCode(shortCode string) (*model.URL, error) {
	query := `SELECT id, short_code, original_url, expires_at FROM urls WHERE short_code = $1 AND (expires_at IS NULL OR expires_at > NOW())`

	var u model.URL
	err := r.DB.QueryRow(query, shortCode).Scan(&u.ID, &u.ShortCode, &u.OriginalURL, &u.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
