package config

import (
	"os"

	"github.com/joho/godotenv"
)

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		App: AppConfig{
			Name: os.Getenv("APP_NAME"),
			Port: os.Getenv("APP_PORT"),
		},
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		},
		JWT: JWTConfig{
			Secret: os.Getenv("JWT_SECRET"),
		},
		Email: EmailConfig{
			Host:     os.Getenv("EMAIL_HOST"),
			Port:     os.Getenv("EMAIL_PORT"),
			Username: os.Getenv("EMAIL_USERNAME"),
			Password: os.Getenv("EMAIL_PASSWORD"),
			From:     os.Getenv("EMAIL_FROM"),
		},
	}

	return cfg, nil
}
