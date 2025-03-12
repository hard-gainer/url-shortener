package postgres

import (
	"context"
	"log/slog"

	"github.com/hard-gainer/url-shortener/internal/config"
	"github.com/hard-gainer/url-shortener/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

// A postgres implementation of the repository
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewPostgres creates a new PostgreSQL repository with connection pool
func NewPostgres(cfg *config.Config) (storage.Repository, error) {
	connPool, err := pgxpool.New(context.Background(), cfg.Url)
	if err != nil {
		slog.Error("postgres.NewPostgres", "unable to create connection pool", err)
		return nil, err
	}
	return &PostgresRepository{pool: connPool}, nil
}


func (repo *PostgresRepository) GetUrl(ctx context.Context, shortUrl string) error {
	return nil
}

func (repo *PostgresRepository) SaveUrl(ctx context.Context, shortUrl, originalUrl string) error {
	return nil
}

func (repo *PostgresRepository) Close() error {
	return nil
}
