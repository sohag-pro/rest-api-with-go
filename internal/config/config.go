// Package config loads and validates runtime configuration from the environment.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all runtime settings.
type Config struct {
	Port            string        // PORT, default "3000"
	DBPath          string        // DB_PATH, default "books.db"
	APIKey          string        // API_KEY, empty disables write auth
	LogLevel        string        // LOG_LEVEL: debug|info|warn|error, default "info"
	ReadTimeout     time.Duration // READ_TIMEOUT, default 10s
	WriteTimeout    time.Duration // WRITE_TIMEOUT, default 10s
	ShutdownTimeout time.Duration // SHUTDOWN_TIMEOUT, default 10s
}

// Load reads configuration from the environment and validates it.
func Load() (Config, error) {
	cfg := Config{
		Port:            getenv("PORT", "3000"),
		DBPath:          getenv("DB_PATH", "books.db"),
		APIKey:          os.Getenv("API_KEY"),
		LogLevel:        getenv("LOG_LEVEL", "info"),
		ReadTimeout:     getDuration("READ_TIMEOUT", 10*time.Second),
		WriteTimeout:    getDuration("WRITE_TIMEOUT", 10*time.Second),
		ShutdownTimeout: getDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
	}
	if err := cfg.validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) validate() error {
	port, err := strconv.Atoi(c.Port)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("invalid PORT %q: must be 1-65535", c.Port)
	}
	if c.DBPath == "" {
		return fmt.Errorf("DB_PATH must not be empty")
	}
	switch c.LogLevel {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("invalid LOG_LEVEL %q: want debug|info|warn|error", c.LogLevel)
	}
	if c.ReadTimeout <= 0 || c.WriteTimeout <= 0 || c.ShutdownTimeout <= 0 {
		return fmt.Errorf("timeouts must be positive durations")
	}
	return nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
