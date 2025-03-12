package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// type Config struct {
// 	DBConfig `yaml:"db"`
// }

// type DBConfig struct {
// 	Path     string `yaml:"DB_URL" env:"DB_URL" env-default:"postgresql://root:password@localhost:5432/auth?sslmode=disable"`
// 	User     string `yaml:"DB_USER" env:"DB_USER" env-default:"root"`
// 	Password string `yaml:"DB_PASSWORD" env:"DB_PASSWORD" env-default:"password"`
// 	Host     string `yaml:"DB_HOST" env:"DB_HOST" env-default:"localhost"`
// 	Port     string `yaml:"DB_PORT" env:"DB_PORT" env-default:"5432"`
// 	Name     string `yaml:"DB_NAME" env:"DB_NAME" env-default:"url-shortener"`
// }

type Config struct {
	DBConfig
}

type DBConfig struct {
	Path     string
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

func InitConfig() *Config {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	cfg := &Config{
		DBConfig: DBConfig{
			User: dbUser,
			Password: dbPassword,
			Host: dbHost,
			Port: dbPort,
			Name: dbName,
		},
	}

	setUpDBUrl(cfg)

	return cfg
}

func setUpDBUrl(cfg *Config) {
	url := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host,
		cfg.Port, cfg.Name)
	cfg.Path = url
}
