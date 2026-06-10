package repository

import (
	"database/sql"
	"time"
)

type ClickRepository struct {
	DB *sql.DB
}

func NewClickRepository(db *sql.DB) *ClickRepository {
	return &ClickRepository{DB: db}
}

func (r *ClickRepository) Save(shortCode, ipAddress, userAgent, referer string) error {
	query := `
		INSERT INTO clicks (short_code, ip_address, user_agent, referer)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.DB.Exec(query, shortCode, ipAddress, userAgent, referer)
	return err
}

type ClickRecord struct {
	ShortCode string
	IPAddress string
	UserAgent string
	Referer   string
	ClickedAt time.Time
}

func (r *ClickRepository) GetByShortCode(shortCode string, limit, offset int) ([]ClickRecord, error) {
	query := `
		SELECT short_code, ip_address, user_agent, referer, clicked_at
		FROM clicks
		WHERE short_code = $1
		ORDER BY clicked_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.DB.Query(query, shortCode, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []ClickRecord
	for rows.Next() {
		var rec ClickRecord
		if err := rows.Scan(&rec.ShortCode, &rec.IPAddress, &rec.UserAgent, &rec.Referer, &rec.ClickedAt); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	return records, nil
}

func (r *ClickRepository) CountByShortCode(shortCode string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM clicks WHERE short_code = $1`
	err := r.DB.QueryRow(query, shortCode).Scan(&count)
	return count, err
}
