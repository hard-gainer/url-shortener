package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"time"

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

// GetOriginalURL получает оригинальный URL по короткому URL
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
	timestamp := time.Now().UnixNano()

	salt, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}

	result := timestamp + salt.Int64()
	return encodeNumber(result), nil
}

// encodeNumber encodes a number to a string with the use of Charset
func encodeNumber(num int64) string {
    base := int64(len(Charset))
    var result string
    
    for num > 0 {
        remainder := num % base
        result = string(Charset[remainder]) + result
        num /= base
    }
    
	for len(result) < ShortURLLength {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(Charset))))
		result = string(Charset[num.Int64()]) + result
	}
    
    return result
}
