package storage

import "context"

type Repository interface {
	// GetUrl retrieves the original url from the storage for a short url 
	GetUrl(ctx context.Context, shortUrl string) error
	// SaveUrl saves a new pair of short url and original url into the storage
	SaveUrl(ctx context.Context, shortUrl, originalUrl string) error
	// Close closes a connection with the storage
	Close() error
}
