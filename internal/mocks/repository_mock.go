package mocks

import (
	"context"

	"github.com/hard-gainer/url-shortener/internal/models"
	"github.com/stretchr/testify/mock"
)

// RepositoryMock is a mock of the storage 
type RepositoryMock struct {
	mock.Mock
}

// GetURL is a mock of GetURL
func (m *RepositoryMock) GetURL(ctx context.Context, shortURL string) (models.Url, error) {
	args := m.Called(ctx, shortURL)
	return args.Get(0).(models.Url), args.Error(1)
}

// SaveURL is a mock of SaveURL
func (m *RepositoryMock) SaveURL(ctx context.Context, shortURL, originalURL string) (int64, error) {
	args := m.Called(ctx, shortURL, originalURL)
	return args.Get(0).(int64), args.Error(1)
}

// OriginalURLExists is a mock of OriginalURLExists
func (m *RepositoryMock) OriginalURLExists(ctx context.Context, originalURL string) (string, bool, error) {
	args := m.Called(ctx, originalURL)
	return args.String(0), args.Bool(1), args.Error(2)
}

// Close is a mock of Close
func (m *RepositoryMock) Close() {
	m.Called()
}
