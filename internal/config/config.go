package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConfig
}

type DBConfig struct {
	Url      string
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

	dbUrl := os.Getenv("DB_URL")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	cfg := &Config{
		DBConfig: DBConfig{
			Url:      dbUrl,
			User:     dbUser,
			Password: dbPassword,
			Host:     dbHost,
			Port:     dbPort,
			Name:     dbName,
		},
	}

	// setUpDBUrl(cfg)
	return cfg
}

// func setUpDBUrl(cfg *Config) {
// 	url := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
// 		cfg.User, cfg.Password, cfg.Host,
// 		cfg.Port, cfg.Name)
// 	cfg.Path = url
// }
