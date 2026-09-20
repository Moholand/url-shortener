package model

import "time"

type URL struct {
	ID          int64
	ShortCode   string
	OriginalURL string
	CreatedAt   time.Time
	ExpiresAt   *time.Time
}

type URLListItem struct {
	URL
	ClickCount int
}

type Click struct {
	ID        int64
	ShortCode string
	IPAddress string
	UserAgent string
	Referer   string
	ClickedAt time.Time
}
