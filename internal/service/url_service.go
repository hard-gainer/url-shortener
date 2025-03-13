package service

import (
	"context"
)

// URLService represents a main interface for the service for the url shortening
type URLService interface {
	ShortenURL(ctx context.Context, originalURL string) (string, error)
	GetOriginalURL(ctx context.Context, shortURL string) (string, error)
}
