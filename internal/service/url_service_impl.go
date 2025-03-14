package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/hard-gainer/url-shortener/internal/storage"
)

const (
	// ShortURLLength is a length of generated shortened URL
	ShortURLLength = 10
	// Charset is a set of symbols which will be used in shortened URL
	Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	// MaxRetries is a maximum amount of retries if the collision occured
	MaxRetries = 3
)

type URLServiceImpl struct {
	repo storage.Repository
}

// NewURLService creates a new instance of the URL service
func NewURLService(repo storage.Repository) URLService {
	return &URLServiceImpl{
		repo: repo,
	}
}

// ShortenURL creates a shortened URL for the original one
func (s *URLServiceImpl) ShortenURL(ctx context.Context, originalURL string) (string, error) {
	const op = "service.URLServiceImpl.ShortenURL"

	existingShort, exists, err := s.repo.OriginalURLExists(ctx, originalURL)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if exists {
		slog.Debug("URL already exists", "original_url", originalURL, "short_url", existingShort)
		return existingShort, nil
	}

	for i := 0; i < MaxRetries; i++ {
		shortURL, err := generateShortURL()
		if err != nil {
			return "", fmt.Errorf("%s: failed to generate short URL: %w", op, err)
		}

		_, err = s.repo.SaveURL(ctx, shortURL, originalURL)
		if err != nil {
			if errors.Is(err, storage.ErrURLMappingExists) {
				slog.Debug("URL collision, retrying", "attempt", i+1)
				continue
			}
			return "", fmt.Errorf("%s: %w", op, err)
		}

		return shortURL, nil
	}

	return "", fmt.Errorf("%s: failed to generate unique short URL after %d attempts", op, MaxRetries)
}

// GetOriginalURL gets original URL by a shortend one 
func (s *URLServiceImpl) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	const op = "service.URLServiceImpl.GetOriginalURL"

	url, err := s.repo.GetURL(ctx, shortURL)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url.OriginalURL, nil
}

// generateShortURL generates random string with specified length from the set of symbols
func generateShortURL() (string, error) {
	b := make([]byte, ShortURLLength)

	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(Charset))))
		if err != nil {
			return "", err
		}
		b[i] = Charset[n.Int64()]
	}

	return string(b), nil
}
