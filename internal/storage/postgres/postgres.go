package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/hard-gainer/url-shortener/internal/config"
	"github.com/hard-gainer/url-shortener/internal/models"
	"github.com/hard-gainer/url-shortener/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type qurier interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Begin(ctx context.Context) (pgx.Tx, error)
	Close()
}

// A postgres implementation of the repository
type PostgresRepository struct {
	db qurier
}

// NewPostgres creates a new PostgreSQL repository with connection pool
func NewPostgres(cfg *config.Config) (storage.Repository, error) {
	const op = "storage.postgres.NewPostgres"

	connPool, err := pgxpool.New(context.Background(), cfg.Url)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := connPool.Ping(context.Background()); err != nil {
		connPool.Close()
		return nil, fmt.Errorf("%s: could not ping database: %w", op, err)
	}

	return &PostgresRepository{db: connPool}, nil
}

// GetURL retrieves the url from the storage by its short url
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

// SaveURL saves a new pair of short url and original url into the storage
func (repo *PostgresRepository) SaveURL(ctx context.Context, shortURL, originalURL string) (int64, error) {
	const op = "storage.postgres.SaveURL"

	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	var existingShort string
	err = tx.QueryRow(ctx,
		`SELECT short_url FROM url_mappings WHERE original_url = $1`,
		originalURL).Scan(&existingShort)

	if err == nil {
		return 0, fmt.Errorf("%s: %w: %s", op, storage.ErrOriginalURLExists, existingShort)
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf("%s: checking existing URL: %w", op, err)
	}

	var id int64
	err = tx.QueryRow(ctx,
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

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return id, nil
}

// Close closes a connection with the storage
func (repo *PostgresRepository) Close() {
	repo.db.Close()
}
