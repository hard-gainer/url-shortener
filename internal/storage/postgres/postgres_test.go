package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hard-gainer/url-shortener/internal/config"
	"github.com/hard-gainer/url-shortener/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// setupPostgres sets up a PostgreSQL container for testing
func setupPostgres(t *testing.T) (string, func()) {
	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	require.NoError(t, err)

	host, err := postgresContainer.Host(ctx)
	require.NoError(t, err)
	port, err := postgresContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	connString := fmt.Sprintf(
		"postgres://testuser:testpass@%s:%s/testdb",
		host, port.Port(),
	)

	cleanup := func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	}

	pool, err := pgxpool.New(ctx, connString)
	require.NoError(t, err)
	defer pool.Close()

	_, err = pool.Exec(ctx, `
        CREATE TABLE IF NOT EXISTS url_mappings (
            id SERIAL PRIMARY KEY,
            short_url TEXT UNIQUE NOT NULL,
            original_url TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT NOW()
        )
    `)
	require.NoError(t, err)

	return connString, cleanup
}

// TestRepository tests all repository methods
func TestRepository(t *testing.T) {
	connString, cleanup := setupPostgres(t)
	defer cleanup()

	cfg := &config.Config{
		DBConfig: config.DBConfig{
			URL: connString,
		},
	}

	repo, err := NewPostgres(cfg)
	require.NoError(t, err)
	defer repo.Close()

	t.Run("SaveURL and GetURL", func(t *testing.T) {
		ctx := context.Background()
		shortURL := "abc123test"
		originalURL := "https://example.com"

		id, err := repo.SaveURL(ctx, shortURL, originalURL)
		require.NoError(t, err)
		assert.Greater(t, id, int64(0))

		url, err := repo.GetURL(ctx, shortURL)
		require.NoError(t, err)
		assert.Equal(t, shortURL, url.ShortURL)
		assert.Equal(t, originalURL, url.OriginalURL)
		assert.Equal(t, id, url.Id)
	})

	t.Run("GetURL Not Found", func(t *testing.T) {
		ctx := context.Background()
		shortURL := "nonexistent"

		_, err := repo.GetURL(ctx, shortURL)
		require.Error(t, err)
		assert.ErrorIs(t, err, storage.ErrURLMappingNotFound)
	})

	t.Run("SaveURL Duplicate Short URL", func(t *testing.T) {
		ctx := context.Background()
		shortURL := "duplicate"
		originalURL1 := "https://example1.com"
		originalURL2 := "https://example2.com"

		_, err := repo.SaveURL(ctx, shortURL, originalURL1)
		require.NoError(t, err)

		_, err = repo.SaveURL(ctx, shortURL, originalURL2)
		require.Error(t, err)
		assert.ErrorIs(t, err, storage.ErrURLMappingExists)
	})

	t.Run("SaveURL Duplicate Original URL", func(t *testing.T) {
		ctx := context.Background()
		shortURL1 := "original1"
		shortURL2 := "original2"
		originalURL := "https://duplicate-original.com"

		_, err := repo.SaveURL(ctx, shortURL1, originalURL)
		require.NoError(t, err)

		_, err = repo.SaveURL(ctx, shortURL2, originalURL)
		require.Error(t, err)
		assert.ErrorIs(t, err, storage.ErrOriginalURLExists)
		assert.Contains(t, err.Error(), shortURL1) // Должен содержать существующий shortURL
	})

	t.Run("OriginalURLExists", func(t *testing.T) {
		ctx := context.Background()
		shortURL := "exists_test"
		originalURL := "https://exists.example.com"

		_, err := repo.SaveURL(ctx, shortURL, originalURL)
		require.NoError(t, err)

		existingShort, exists, err := repo.OriginalURLExists(ctx, originalURL)
		require.NoError(t, err)
		assert.True(t, exists)
		assert.Equal(t, shortURL, existingShort)

		existingShort, exists, err = repo.OriginalURLExists(ctx, "https://nonexistent.com")
		require.NoError(t, err)
		assert.False(t, exists)
		assert.Empty(t, existingShort)
	})

	t.Run("Transaction Rollback on Error", func(t *testing.T) {
		ctx := context.Background()
		shortURL := "rollback_test"
		originalURL := "https://rollback.example.com"

		_, err := repo.SaveURL(ctx, shortURL, originalURL)
		require.NoError(t, err)

		pool, err := pgxpool.New(ctx, cfg.DBConfig.URL)
		require.NoError(t, err)
		defer pool.Close()

		_, err = pool.Exec(ctx, 
			"ALTER TABLE url_mappings ADD CONSTRAINT unique_original_url UNIQUE (original_url)",
		)
		require.NoError(t, err)

		_, err = repo.SaveURL(ctx, "another_short", originalURL)
		require.Error(t, err)

		existingShort, exists, err := repo.OriginalURLExists(ctx, originalURL)
		require.NoError(t, err)
		assert.True(t, exists)
		assert.Equal(t, shortURL, existingShort)
	})
}
