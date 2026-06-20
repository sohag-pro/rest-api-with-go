package main

import "os"

// Config holds runtime configuration sourced from environment variables.
type Config struct {
	Port   string // PORT, default "3000"
	DBPath string // DB_PATH, default "books.db"
	APIKey string // API_KEY, empty disables auth
}

func loadConfig() Config {
	return Config{
		Port:   getenv("PORT", "3000"),
		DBPath: getenv("DB_PATH", "books.db"),
		APIKey: os.Getenv("API_KEY"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
