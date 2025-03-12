package models

import "time"

type Url struct {
	Id          int64
	ShortUrl    string
	OriginalUrl string
	CreatedAt   time.Time
}
