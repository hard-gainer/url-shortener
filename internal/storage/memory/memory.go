package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hard-gainer/url-shortener/internal/models"
	"github.com/hard-gainer/url-shortener/internal/storage"
)

// A memory implementation of the repository
type MemoryRepository struct {
	shortToOriginal map[string]string
	originalToShort map[string]string
	urls            map[int64]models.Url
	mutex           sync.RWMutex
	lastID          int64
}

// NewMemory creates a new memory repository with maps and rwmutex
func NewMemory() (storage.Repository, error) {
	return &MemoryRepository{
		shortToOriginal: make(map[string]string),
		originalToShort: make(map[string]string),
		urls:            make(map[int64]models.Url),
		lastID:          0,
	}, nil
}

func (repo *MemoryRepository) GetURL(ctx context.Context, shortURL string) (models.Url, error) {
	const op = "storage.memory.GetURL"

	select {
	case <-ctx.Done():
		return models.Url{}, ctx.Err()
	default:
	}

	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	_, exists := repo.shortToOriginal[shortURL]
	if !exists {
		return models.Url{}, fmt.Errorf("%s: %w", op, storage.ErrURLMappingNotFound)
	}

	var url models.Url
	for _, u := range repo.urls {
		if u.ShortURL == shortURL {
			url = u
			break
		}
	}

	return url, nil
}

func (repo *MemoryRepository) SaveURL(ctx context.Context, shortURL, originalURL string) (int64, error) {
	const op = "storage.memory.SaveURL"

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	if _, exists := repo.shortToOriginal[shortURL]; exists {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrURLMappingExists)
	}

	if existingShort, exists := repo.originalToShort[originalURL]; exists {
		for _, url := range repo.urls {
			if url.ShortURL == existingShort {
				return url.Id, fmt.Errorf("%s: %w", op, storage.ErrOriginalURLExists)
			}
		}
	}

	repo.lastID++

	urlModel := models.Url{
		Id:          repo.lastID,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
		CreatedAt:   time.Now(),
	}

	repo.shortToOriginal[shortURL] = originalURL
	repo.originalToShort[originalURL] = shortURL
	repo.urls[repo.lastID] = urlModel

	return repo.lastID, nil
}

func (repo *MemoryRepository) Close() error {
	return nil
}
