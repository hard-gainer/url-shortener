package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hard-gainer/url-shortener/internal/api"
	"github.com/hard-gainer/url-shortener/internal/config"
	"github.com/hard-gainer/url-shortener/internal/logger"
	"github.com/hard-gainer/url-shortener/internal/service"
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

	urlService := service.NewURLService(repo)

	server := api.NewServer("0.0.0.0:" + cfg.AppConfig.Port)
    server.WithMiddleware(api.LoggingMiddleware)

	urlHandler := api.NewURLHandler(urlService, cfg.AppConfig.URL)
	urlHandler.RegisterRoutes(server.Mux())

	go func() {
		if err := server.Run(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}

	slog.Info("server exited properly")
}
