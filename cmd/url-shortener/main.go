package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/hard-gainer/url-shortener/internal/config"
	"github.com/hard-gainer/url-shortener/internal/logger"
	"github.com/hard-gainer/url-shortener/internal/storage"
	"github.com/hard-gainer/url-shortener/internal/storage/memory"
	"github.com/hard-gainer/url-shortener/internal/storage/postgres"
)

func main() {
	logger.InitLogger()
	cfg := config.InitConfig()

	storageType := flag.String("storage", "postgres", "Storage type")
	flag.Parse()
	slog.Info("initializing storage", "storage type", *storageType)

	var repo storage.Repository

	switch *storageType {
	case "memory":
		var err error
		repo, err = memory.NewMemory()
		if err != nil {
			slog.Error("failed to initialize memory storage", "error", err)
			os.Exit(1)
		}
	case "postgres":
		var err error
		repo, err = postgres.NewPostgres(cfg)
		if err != nil {
			slog.Error("failed to initialize storage", "error", err)
			os.Exit(1)
		}
	default:
		slog.Error("unknown storage type", "error", *storageType)
		os.Exit(1)
	}

	defer repo.Close()
	slog.Info("storage successfully intialized")

}
