package main

import (
	"fmt"
	"log/slog"

	"github.com/hard-gainer/url-shortener/internal/config"
	"github.com/hard-gainer/url-shortener/internal/logger"
)

func main() {
	logger.InitLogger()
	cfg := config.InitConfig()

	fmt.Println(cfg.Path)
	// storageType := flag.String("storage", "memory", "Storage type")

	slog.Info("initializing storage")
	// storage :=
}
