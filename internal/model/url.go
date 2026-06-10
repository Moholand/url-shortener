package model

import "time"

type URL struct {
	ID          int64
	ShortCode   string
	OriginalURL string
}

type Click struct {
	ID        int64
	ShortCode string
	IPAddress string
	UserAgent string
	Referer   string
	ClickedAt time.Time
}
