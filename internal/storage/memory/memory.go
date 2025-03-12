package memory

import (
	"context"
	"sync"

	"github.com/hard-gainer/url-shortener/internal/storage"
)

// A memory implementation of the repository
type MemoryRepository struct {
	shortToOriginal map[string]string
	originalToShort map[string]string
	mutex           sync.RWMutex
}

// NewMemory creates a new memory repository with maps and rwmutex
func NewMemory() (storage.Repository, error) {
	return &MemoryRepository{
		shortToOriginal: make(map[string]string),
		originalToShort: make(map[string]string),
	}, nil
}

func (repo *MemoryRepository) GetUrl(ctx context.Context, shortUrl string) error {
	return nil
}

func (repo *MemoryRepository) SaveUrl(ctx context.Context, shortUrl, originalUrl string) error {
	return nil
}

func (repo *MemoryRepository) Close() error {
	return nil
}
