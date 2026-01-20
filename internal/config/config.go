package config

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	App           AppConfig
	HTTP          HTTPConfig
	DB            DBConfig
	Redis         RedisConfig
	JWT           JWTConfig
	OTP           OTPConfig
	SMS           SMSConfig
	RateLimit     RateLimitConfig
	CORS          CORSConfig
	Observability ObservabilityConfig
}

type AppConfig struct {
	Name     string `env:"APP_NAME" envDefault:"stin-smart-care-be"`
	Env      string `env:"APP_ENV" envDefault:"local"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}

type HTTPConfig struct {
	Host          string        `env:"HTTP_HOST" envDefault:"0.0.0.0"`
	Port          int           `env:"HTTP_PORT" envDefault:"8080"`
	ReadTimeout   time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"15s"`
	WriteTimeout  time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"15s"`
	IdleTimeout   time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"60s"`
	BasePath      string        `env:"HTTP_BASE_PATH" envDefault:"/api/v1"`
	EnableSwagger bool          `env:"HTTP_ENABLE_SWAGGER" envDefault:"true"`
	EnableMetrics bool          `env:"HTTP_ENABLE_METRICS" envDefault:"true"`
	TLSCertFile   string        `env:"TLS_CERT_FILE"`
	TLSKeyFile    string        `env:"TLS_KEY_FILE"`
}

type DBConfig struct {
	Host            string        `env:"DB_HOST" envDefault:"localhost"`
	Port            int           `env:"DB_PORT" envDefault:"5432"`
	Name            string        `env:"DB_NAME" envDefault:"stin_smart_care"`
	User            string        `env:"DB_USER" envDefault:"stin"`
	Password        string        `env:"DB_PASSWORD" envDefault:"stin_pass"`
	SSLMode         string        `env:"DB_SSLMODE" envDefault:"disable"`
	MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" envDefault:"10"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"30m"`
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" envDefault:"localhost"`
	Port     int    `env:"REDIS_PORT" envDefault:"6379"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB" envDefault:"0"`
}

type JWTConfig struct {
	Issuer     string        `env:"JWT_ISSUER" envDefault:"stin-smart-care"`
	Secret     string        `env:"JWT_SECRET" envDefault:"change_me"`
	AccessTTL  time.Duration `env:"JWT_ACCESS_TTL" envDefault:"15m"`
	RefreshTTL time.Duration `env:"JWT_REFRESH_TTL" envDefault:"720h"`
}

type OTPConfig struct {
	TTL           time.Duration `env:"OTP_TTL" envDefault:"5m"`
	Digits        int           `env:"OTP_DIGITS" envDefault:"6"`
	RefCodeLength int           `env:"OTP_REF_CODE_LENGTH" envDefault:"6"`
}

type SMSConfig struct {
	Provider    string `env:"SMS_PROVIDER" envDefault:"console"`
	ThaiBulkSMS ThaiBulkSMSConfig
}

type ThaiBulkSMSConfig struct {
	BaseURL     string        `env:"THAIBULKSMS_BASE_URL" envDefault:"https://api.thaibulksms.com"`
	Endpoint    string        `env:"THAIBULKSMS_ENDPOINT" envDefault:"/sms"`
	APIKey      string        `env:"THAIBULKSMS_API_KEY"`
	APISecret   string        `env:"THAIBULKSMS_API_SECRET"`
	SenderID    string        `env:"THAIBULKSMS_SENDER_ID"`
	AuthMode    string        `env:"THAIBULKSMS_AUTH_MODE" envDefault:"basic"`
	OTPTemplate string        `env:"THAIBULKSMS_OTP_TEMPLATE" envDefault:"Your OTP is {{otp}} (ref: {{ref}})"`
	Timeout     time.Duration `env:"THAIBULKSMS_TIMEOUT" envDefault:"10s"`
}

type RateLimitConfig struct {
	OTPPerPhone int           `env:"OTP_RATE_LIMIT_PER_PHONE" envDefault:"5"`
	OTPPerIP    int           `env:"OTP_RATE_LIMIT_PER_IP" envDefault:"5"`
	Window      time.Duration `env:"OTP_RATE_LIMIT_WINDOW" envDefault:"1m"`
}

type CORSConfig struct {
	AllowedOrigins   string `env:"CORS_ALLOWED_ORIGINS" envDefault:"*"`
	AllowedMethods   string `env:"CORS_ALLOWED_METHODS" envDefault:"GET,POST,PATCH,DELETE,OPTIONS"`
	AllowedHeaders   string `env:"CORS_ALLOWED_HEADERS" envDefault:"Authorization,Content-Type,X-Request-Id"`
	AllowCredentials bool   `env:"CORS_ALLOW_CREDENTIALS" envDefault:"false"`
}

type ObservabilityConfig struct {
	OtelEnabled bool   `env:"OTEL_ENABLED" envDefault:"false"`
	ServiceName string `env:"OTEL_SERVICE_NAME" envDefault:"stin-smart-care-be"`
}

func Load() (Config, error) {
	_ = godotenv.Load()

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Name,
		c.SSLMode,
	)
}

func (c DBConfig) URL() string {
	user := url.UserPassword(c.User, c.Password)
	return fmt.Sprintf(
		"postgres://%s@%s:%d/%s?sslmode=%s",
		user.String(),
		c.Host,
		c.Port,
		c.Name,
		c.SSLMode,
	)
}

func (c RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c CORSConfig) OriginsList() []string {
	if strings.TrimSpace(c.AllowedOrigins) == "" {
		return nil
	}
	if c.AllowedOrigins == "*" {
		return []string{"*"}
	}
	parts := strings.Split(c.AllowedOrigins, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func (c CORSConfig) MethodsList() []string {
	parts := strings.Split(c.AllowedMethods, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func (c CORSConfig) HeadersList() []string {
	parts := strings.Split(c.AllowedHeaders, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
