package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/hard-gainer/url-shortener/internal/config"
	"github.com/hard-gainer/url-shortener/internal/models"
	"github.com/hard-gainer/url-shortener/internal/storage"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// A postgres implementation of the repository
type PostgresRepository struct {
	db *pgxpool.Pool
}

// NewPostgres creates a new PostgreSQL repository with connection pool
func NewPostgres(cfg *config.Config) (storage.Repository, error) {
	connPool, err := pgxpool.New(context.Background(), cfg.Url)
	if err != nil {
		slog.Error("postgres.NewPostgres", "unable to create connection pool", err)
		return nil, err
	}
	return &PostgresRepository{db: connPool}, nil
}

// GetUrl retrieves the url from the storage by its short url
func (repo *PostgresRepository) GetURL(ctx context.Context, shortURL string) (models.Url, error) {
	const op = "storage.postgres.GetURL"
	var url models.Url

	err := repo.db.QueryRow(ctx,
		`SELECT id, short_url, original_url, created_at 
         FROM url_mappings  
         WHERE short_url = $1`,
		shortURL).Scan(&url.Id, &url.ShortURL, &url.OriginalURL, &url.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Url{}, fmt.Errorf("%s: %w", op, storage.ErrURLMappingNotFound)
		}
		return models.Url{}, fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}

// SaveUrl saves a new pair of short url and original url into the storage
func (repo *PostgresRepository) SaveURL(ctx context.Context, shortURL, originalURL string) (int64, error) {
	const op = "storage.postgres.SaveURL"
	var id int64

	err := repo.db.QueryRow(ctx,
		`INSERT INTO url_mappings(short_url, original_url)
         VALUES($1, $2)
         RETURNING id`,
		shortURL, originalURL).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLMappingExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// Close closes a connection with the storage
func (repo *PostgresRepository) Close() error {
	if repo.db != nil {
		repo.db.Close()
	}
	return nil
}
