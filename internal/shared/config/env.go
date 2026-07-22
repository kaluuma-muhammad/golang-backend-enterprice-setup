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
	emailPort, _ := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	cfg := &Config{
		App: AppConfig{
			Name:    os.Getenv("APP_NAME"),
			Port:    os.Getenv("APP_PORT"),
			BaseURL: os.Getenv("APP_BASE_URL"),
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
			Host:        os.Getenv("EMAIL_HOST"),
			Port:        emailPort,
			Username:    os.Getenv("EMAIL_USERNAME"),
			Password:    os.Getenv("EMAIL_PASSWORD"),
			FromAddress: os.Getenv("EMAIL_FROM_ADDRESS"),
			FromName:    os.Getenv("EMAIL_FROM_NAME"),
			Encryption:  os.Getenv("EMAIL_ENCRYPTION"),
		},
		Admin: AdminConfig{
			Email:     os.Getenv("SEED_ADMIN_EMAIL"),
			Password:  os.Getenv("SEED_ADMIN_PASSWORD"),
			FirstName: os.Getenv("SEED_ADMIN_FIRST_NAME"),
			LastName:  os.Getenv("SEED_ADMIN_LAST_NAME"),
		},
		Redis: RedisConfig{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       redisDB,
		},
		RateLimiter: RateLimitConfig{
			Store: os.Getenv("RATE_LIMIT_STORE"),
		},
	}

	return cfg, nil
}
