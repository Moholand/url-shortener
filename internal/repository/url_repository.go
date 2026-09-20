package repository

import (
	"database/sql"
	"fmt"
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

func (r *URLRepository) DeleteByShortCode(shortCode string) error {
	query := `DELETE FROM urls WHERE short_code = $1`
	_, err := r.DB.Exec(query, shortCode)
	if err != nil {
		return err
	}
	return nil
}

func (r *URLRepository) GetURLByShortCode(shortCode string) (*model.URL, error) {
	query := `SELECT id, short_code, original_url, created_at, expires_at FROM urls WHERE short_code = $1 AND (expires_at IS NULL OR expires_at > NOW())`

	var u model.URL
	err := r.DB.QueryRow(query, shortCode).Scan(&u.ID, &u.ShortCode, &u.OriginalURL, &u.CreatedAt, &u.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *URLRepository) ListAll(limit, offset int, search string) ([]model.URLListItem, error) {
	query := `
		SELECT u.id, u.short_code, u.original_url, u.created_at, u.expires_at,
		       COALESCE(c.cnt, 0) AS click_count
		FROM urls u
		LEFT JOIN (
			SELECT short_code, COUNT(*) AS cnt
			FROM clicks
			GROUP BY short_code
		) c ON u.short_code = c.short_code
	`

	var args []interface{}
	argIdx := 1

	if search != "" {
		query += fmt.Sprintf(` WHERE u.original_url ILIKE '%%' || $%d || '%%' OR u.short_code ILIKE '%%' || $%d || '%%'`, argIdx, argIdx)
		args = append(args, search)
		argIdx++
	}

	query += fmt.Sprintf(` ORDER BY u.created_at DESC LIMIT $%d OFFSET $%d`, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.URLListItem
	for rows.Next() {
		var item model.URLListItem
		if err := rows.Scan(&item.ID, &item.ShortCode, &item.OriginalURL, &item.CreatedAt, &item.ExpiresAt, &item.ClickCount); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *URLRepository) CountAll(search string) (int, error) {
	query := `SELECT COUNT(*) FROM urls`
	var args []interface{}

	if search != "" {
		query += ` WHERE original_url ILIKE '%' || $1 || '%' OR short_code ILIKE '%' || $1 || '%'`
		args = append(args, search)
	}

	var count int
	err := r.DB.QueryRow(query, args...).Scan(&count)
	return count, err
}
