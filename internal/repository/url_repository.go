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
