package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/hard-gainer/url-shortener/internal/mocks"
	"github.com/hard-gainer/url-shortener/internal/models"
	"github.com/hard-gainer/url-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestShortenURL_Success(t *testing.T) {
	mockRepo := new(mocks.RepositoryMock)
	service := NewURLService(mockRepo)
	ctx := context.Background()
	originalURL := "https://example.com"

	mockRepo.On("OriginalURLExists", ctx, originalURL).
		Return("", false, nil)

	mockRepo.On("SaveURL", ctx, mock.AnythingOfType("string"), originalURL).
		Return(int64(1), nil)

	shortURL, err := service.ShortenURL(ctx, originalURL)

	require.NoError(t, err)
	assert.Len(t, shortURL, ShortURLLength + 1)

	for _, char := range shortURL {
		assert.True(t, strings.ContainsRune(Charset, char),
			"short URL contains invalid character: %c", char)
	}

	mockRepo.AssertExpectations(t)
}

func TestShortenURL_ExistingURL(t *testing.T) {
	mockRepo := new(mocks.RepositoryMock)
	service := NewURLService(mockRepo)
	ctx := context.Background()
	originalURL := "https://example.com"
	existingShortURL := "existing123"

	mockRepo.On("OriginalURLExists", ctx, originalURL).
		Return(existingShortURL, true, nil)

	shortURL, err := service.ShortenURL(ctx, originalURL)

	require.NoError(t, err)
	assert.Equal(t, existingShortURL, shortURL)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "SaveURL")
}

func TestShortenURL_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.RepositoryMock)
	service := NewURLService(mockRepo)
	ctx := context.Background()
	originalURL := "https://example.com"
	expectedError := errors.New("database error")

	mockRepo.On("OriginalURLExists", ctx, originalURL).
		Return("", false, expectedError)

	_, err := service.ShortenURL(ctx, originalURL)

	require.Error(t, err)
	assert.ErrorContains(t, err, expectedError.Error())

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "SaveURL")
}

func TestGetOriginalURL_Success(t *testing.T) {
	mockRepo := new(mocks.RepositoryMock)
	service := NewURLService(mockRepo)
	ctx := context.Background()
	shortURL := "abc123"
	expectedOriginalURL := "https://example.com"

	mockRepo.On("GetURL", ctx, shortURL).
		Return(models.Url{
			Id:          1,
			ShortURL:    shortURL,
			OriginalURL: expectedOriginalURL,
			CreatedAt:   time.Now(),
		}, nil)

	originalURL, err := service.GetOriginalURL(ctx, shortURL)

	require.NoError(t, err)
	assert.Equal(t, expectedOriginalURL, originalURL)

	mockRepo.AssertExpectations(t)
}

func TestGetOriginalURL_NotFound(t *testing.T) {
	mockRepo := new(mocks.RepositoryMock)
	service := NewURLService(mockRepo)
	ctx := context.Background()
	shortURL := "nonexistent"

	mockRepo.On("GetURL", ctx, shortURL).
		Return(models.Url{}, storage.ErrURLMappingNotFound)

	_, err := service.GetOriginalURL(ctx, shortURL)

	require.Error(t, err)
	assert.ErrorIs(t, errors.Unwrap(err), storage.ErrURLMappingNotFound)

	mockRepo.AssertExpectations(t)
}

func TestGetOriginalURL_DatabaseError(t *testing.T) {
	mockRepo := new(mocks.RepositoryMock)
	service := NewURLService(mockRepo)
	ctx := context.Background()
	shortURL := "abc123"
	expectedError := errors.New("database connection failed")

	mockRepo.On("GetURL", ctx, shortURL).
		Return(models.Url{}, expectedError)

	_, err := service.GetOriginalURL(ctx, shortURL)

	require.Error(t, err)
	assert.ErrorIs(t, errors.Unwrap(err), expectedError)

	mockRepo.AssertExpectations(t)
}
