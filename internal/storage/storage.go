package storage

import (
	"context"

	"github.com/hard-gainer/url-shortener/internal/models"
)

type Repository interface {
	// GetUrl retrieves the url from the storage by its short url
	GetURL(ctx context.Context, shortURL string) (models.Url, error)
	// SaveUrl saves a new pair of short url and original url into the storage
	SaveURL(ctx context.Context, shortURL, originalURL string) (int64, error)
	// Close closes a connection with the storage
	Close()
}
