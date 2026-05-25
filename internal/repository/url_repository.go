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
		INSERT INTO urls (short_code, original_url)
		VALUES ($1, $2)
		RETURNING id;
	`

	return r.DB.QueryRow(query, url.ShortCode, url.OriginalURL).Scan(&url.ID)
}

func (r *URLRepository) GetByShortCode(shortCode string) (string, error) {
	var originalURL string
	query := `SELECT original_url FROM urls WHERE short_code = $1`
	
	err := r.DB.QueryRow(query, shortCode).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // Not Found!
		}
		return "", err
	}
	return originalURL, nil
}
