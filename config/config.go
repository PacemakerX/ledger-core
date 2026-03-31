package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type Config struct {
	App                      AppConfig
	Database                 DatabaseConfig
	Log                      LogConfig
	PlatformFloatAccountID   uuid.UUID
	PlatformCashAccountID    uuid.UUID
	PlatformRevenueAccountID uuid.UUID
}

type AppConfig struct {
	Env  string
	Port string
}

type DatabaseConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	Name         string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	ConnTimeout  time.Duration
}

type LogConfig struct {
	Level string
}

func Load() (*Config, error) {
	// Only load .env file in development
	// In production, env vars are injected directly
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	maxOpenConns, err := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "25"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_OPEN_CONNS: %w", err)
	}

	maxIdleConns, err := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "5"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_IDLE_CONNS: %w", err)
	}

	connTimeout, err := time.ParseDuration(getEnv("DB_CONN_TIMEOUT", "5s"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_CONN_TIMEOUT: %w", err)
	}

	platformFloatAccountID, err := uuid.Parse(getEnv("PLATFORM_FLOAT_ACCOUNT_ID", ""))
	if err != nil {
		return nil, fmt.Errorf("invalid PLATFORM_FLOAT_ACCOUNT_ID: %w", err)
	}

	platformCashAccountID, err := uuid.Parse(getEnv("PLATFORM_CASH_ACCOUNT_ID", ""))
	if err != nil {
		return nil, fmt.Errorf("invalid PLATFORM_CASH_ACCOUNT_ID: %w", err)
	}

	platformRevenueAccountID, err := uuid.Parse(getEnv("PLATFORM_REVENUE_ACCOUNT_ID", ""))
	if err != nil {
		return nil, fmt.Errorf("invalid PLATFORM_REVENUE_ACCOUNT_ID: %w", err)
	}
	return &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnv("DB_PORT", "5432"),
			User:         getEnv("DB_USER", "ledger_user"),
			Password:     getEnv("DB_PASSWORD", "ledger_password"),
			Name:         getEnv("DB_NAME", "ledger_db"),
			SSLMode:      getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns: maxOpenConns,
			MaxIdleConns: maxIdleConns,
			ConnTimeout:  connTimeout,
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "debug"),
		},
		PlatformFloatAccountID:   platformFloatAccountID,
		PlatformCashAccountID:    platformCashAccountID,
		PlatformRevenueAccountID: platformRevenueAccountID,
	}, nil
}

// DSN returns the PostgreSQL connection string
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d",
		d.Host,
		d.Port,
		d.User,
		d.Password,
		d.Name,
		d.SSLMode,
		int(d.ConnTimeout.Seconds()),
	)
}

// getEnv reads an env variable with a fallback default
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
