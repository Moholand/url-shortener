package model

import "time"

type URL struct {
	ID          int64
	ShortCode   string
	OriginalURL string
	ExpiresAt   *time.Time
}

type Click struct {
	ID        int64
	ShortCode string
	IPAddress string
	UserAgent string
	Referer   string
	ClickedAt time.Time
}
