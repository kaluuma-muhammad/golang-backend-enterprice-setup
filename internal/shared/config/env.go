package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Load() (*Config, error) {
	_ = godotenv.Load()

	accessMinutes, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_TOKEN_MINUTES"))
	refreshDays, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_TOKEN_DAYS"))
	verificationMinutes, _ := strconv.Atoi(os.Getenv("EMAIL_VERIFICATION_MINUTES"))
	passwordResetMinutes, _ := strconv.Atoi(os.Getenv("PASSWORD_RESET_MINUTES"))
	magicLinkMinutes, _ := strconv.Atoi(os.Getenv("MAGIC_LINK_MINUTES"))

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
			Secret:               os.Getenv("JWT_SECRET"),
			Issuer:               os.Getenv("JWT_ISSUER"),
			AccessTokenMinutes:   int(accessMinutes),
			RefreshTokenDays:     int(refreshDays),
			VerificationMinutes:  int(verificationMinutes),
			PasswordResetMinutes: int(passwordResetMinutes),
			MagicLinkMinutes:     int(magicLinkMinutes),
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
