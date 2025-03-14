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

// GetURL retrieves the url from the storage by its short url
func (repo *MemoryRepository) GetURL(ctx context.Context, shortURL string) (models.Url, error) {
	const op = "storage.memory.GetURL"

	if err := ctx.Err(); err != nil {
		return models.Url{}, err
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

// SaveURL saves a new pair of short url and original url into the storage
func (repo *MemoryRepository) SaveURL(ctx context.Context, shortURL, originalURL string) (int64, error) {
	const op = "storage.memory.SaveURL"

	if err := ctx.Err(); err != nil {
		return 0, err
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

// OriginalURLExists checks if an original URL already exists in storage
func (repo *MemoryRepository) OriginalURLExists(ctx context.Context, originalURL string) (string, bool, error) {
    if err := ctx.Err(); err != nil {
        return "", false, err
    }

    repo.mutex.RLock()
    defer repo.mutex.RUnlock()

    shortURL, exists := repo.originalToShort[originalURL]
    return shortURL, exists, nil
}

func (repo *MemoryRepository) Close() {
}
