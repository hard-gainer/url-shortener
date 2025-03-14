package config

import (
	"os"
)

// Config is a main config
type Config struct {
	AppConfig
	DBConfig
}

// AppConfig is a config with specific app information
type AppConfig struct {
	URL  string
	Port string
}

// DBConfig  is a config with specific database information
type DBConfig struct {
	URL      string
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

// InitConfig creates a new Config
func InitConfig() *Config {
	// if err := godotenv.Load(); err != nil {
	// 	panic("No .env file found")
	// }
	// _ = godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	appPort := os.Getenv("APP_PORT")

	cfg := &Config{
		AppConfig: AppConfig{
			Port: appPort,
		},
		DBConfig: DBConfig{
			URL:      dbURL,
			User:     dbUser,
			Password: dbPassword,
			Host:     dbHost,
			Port:     dbPort,
			Name:     dbName,
		},
	}

	return cfg
}
