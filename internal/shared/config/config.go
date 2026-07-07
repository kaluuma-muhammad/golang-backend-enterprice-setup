package config

type Config struct {
	App   AppConfig
	DB    DBConfig
	JWT   JWTConfig
	Email EmailConfig
}

type AppConfig struct {
	Name string
	Port string
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
