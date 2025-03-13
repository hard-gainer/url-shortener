package memory

import (
	"context"
	"testing"
	"time"

	"github.com/hard-gainer/url-shortener/internal/models"
	"github.com/hard-gainer/url-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryRepository_GetURL_Success(t *testing.T) {
	repo, _ := NewMemory()
	ctx := context.Background()
	shortURL := "abc123"
	originalURL := "https://example.com"

	memRepo := repo.(*MemoryRepository)
	memRepo.lastID = 1
	url := models.Url{
		Id:          1,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
		CreatedAt:   time.Now(),
	}
	memRepo.urls[1] = url
	memRepo.shortToOriginal[shortURL] = originalURL
	memRepo.originalToShort[originalURL] = shortURL

	result, err := repo.GetURL(ctx, shortURL)

	require.NoError(t, err)
	assert.Equal(t, url.Id, result.Id)
	assert.Equal(t, shortURL, result.ShortURL)
	assert.Equal(t, originalURL, result.OriginalURL)
}

func TestMemoryRepository_GetURL_NotFound(t *testing.T) {
	repo, _ := NewMemory()
	ctx := context.Background()
	shortURL := "nonexistent"

	_, err := repo.GetURL(ctx, shortURL)

	require.Error(t, err)
	assert.ErrorIs(t, err, storage.ErrURLMappingNotFound)
}

func TestMemoryRepository_GetURL_CanceledContext(t *testing.T) {
	repo, _ := NewMemory()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	shortURL := "abc123"

	_, err := repo.GetURL(ctx, shortURL)

	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestMemoryRepository_SaveURL_Success(t *testing.T) {
	repo, _ := NewMemory()
	ctx := context.Background()
	shortURL := "abc123"
	originalURL := "https://example.com"

	id, err := repo.SaveURL(ctx, shortURL, originalURL)

	require.NoError(t, err)
	assert.Equal(t, int64(1), id)

	memRepo := repo.(*MemoryRepository)
	assert.Equal(t, originalURL, memRepo.shortToOriginal[shortURL])
	assert.Equal(t, shortURL, memRepo.originalToShort[originalURL])
	assert.Equal(t, originalURL, memRepo.urls[id].OriginalURL)
}

func TestMemoryRepository_SaveURL_ShortURLExists(t *testing.T) {
	repo, _ := NewMemory()
	ctx := context.Background()
	shortURL := "abc123"
	originalURL1 := "https://example1.com"
	originalURL2 := "https://example2.com"

	_, err := repo.SaveURL(ctx, shortURL, originalURL1)
	require.NoError(t, err)

	_, err = repo.SaveURL(ctx, shortURL, originalURL2)

	require.Error(t, err)
	assert.ErrorIs(t, err, storage.ErrURLMappingExists)
}

// func TestMemoryRepository_SaveURL_OriginalURLExists(t *testing.T) {
// 	repo, _ := NewMemory()
// 	ctx := context.Background()
// 	shortURL1 := "abc123"
// 	shortURL2 := "def456"
// 	originalURL := "https://example.com"

// 	_, err := repo.SaveURL(ctx, shortURL1, originalURL)
// 	require.NoError(t, err)

// 	_, err = repo.SaveURL(ctx, shortURL2, originalURL)

// 	require.Error(t, err)
// 	assert.Contains(t, err.Error(), storage.ErrOriginalURLExists.Error())
// 	assert.Contains(t, err.Error(), shortURL1)
// }

func TestMemoryRepository_SaveURL_CanceledContext(t *testing.T) {
	repo, _ := NewMemory()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	shortURL := "abc123"
	originalURL := "https://example.com"

	_, err := repo.SaveURL(ctx, shortURL, originalURL)

	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestMemoryRepository_OriginalURLExists_Found(t *testing.T) {
	repo, _ := NewMemory()
	ctx := context.Background()
	shortURL := "abc123"
	originalURL := "https://example.com"

	memRepo := repo.(*MemoryRepository)
	memRepo.originalToShort[originalURL] = shortURL

	resultShortURL, exists, err := repo.OriginalURLExists(ctx, originalURL)

	require.NoError(t, err)
	assert.True(t, exists)
	assert.Equal(t, shortURL, resultShortURL)
}

func TestMemoryRepository_OriginalURLExists_NotFound(t *testing.T) {
	repo, _ := NewMemory()
	ctx := context.Background()
	originalURL := "https://example.com"

	shortURL, exists, err := repo.OriginalURLExists(ctx, originalURL)

	require.NoError(t, err)
	assert.False(t, exists)
	assert.Empty(t, shortURL)
}

func TestMemoryRepository_OriginalURLExists_CanceledContext(t *testing.T) {
	repo, _ := NewMemory()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	originalURL := "https://example.com"

	_, _, err := repo.OriginalURLExists(ctx, originalURL)

	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}
