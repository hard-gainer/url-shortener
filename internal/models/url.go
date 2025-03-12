package models

import "time"

type Url struct {
	Id          int64
	ShortURL    string
	OriginalURL string
	CreatedAt   time.Time
}
