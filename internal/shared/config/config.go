package config

type Config struct {
	App         AppConfig
	DB          DBConfig
	JWT         JWTConfig
	Email       EmailConfig
	Admin       AdminConfig
	Redis       RedisConfig
	RateLimiter RateLimitConfig
}

type AppConfig struct {
	Name    string
	Port    string
	BaseURL string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret               string
	Issuer               string
	AccessTokenMinutes   int
	RefreshTokenDays     int
	VerificationMinutes  int
	PasswordResetMinutes int
	MagicLinkMinutes     int
}

type EmailConfig struct {
	Host        string
	Port        int
	Username    string
	Password    string
	FromAddress string
	FromName    string
	Encryption  string
}

type AdminConfig struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type RateLimitConfig struct {
	Store string
}
