package config

import (
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	for _, k := range []string{"PORT", "DB_PATH", "API_KEY", "LOG_LEVEL", "READ_TIMEOUT", "WRITE_TIMEOUT", "SHUTDOWN_TIMEOUT"} {
		t.Setenv(k, "")
	}
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Port != "3000" || cfg.DBPath != "books.db" || cfg.LogLevel != "info" {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
	if cfg.ReadTimeout != 10*time.Second {
		t.Fatalf("want 10s read timeout, got %s", cfg.ReadTimeout)
	}
}

func TestLoadOverrides(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("DB_PATH", "/tmp/x.db")
	t.Setenv("API_KEY", "abc")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("READ_TIMEOUT", "5s")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Port != "8080" || cfg.APIKey != "abc" || cfg.LogLevel != "debug" || cfg.ReadTimeout != 5*time.Second {
		t.Fatalf("overrides not applied: %+v", cfg)
	}
}

func TestLoadInvalid(t *testing.T) {
	cases := map[string]map[string]string{
		"bad port":  {"PORT": "0"},
		"nan port":  {"PORT": "abc"},
		"bad level": {"LOG_LEVEL": "trace"},
	}
	for name, env := range cases {
		t.Run(name, func(t *testing.T) {
			for _, k := range []string{"PORT", "LOG_LEVEL"} {
				t.Setenv(k, "")
			}
			for k, v := range env {
				t.Setenv(k, v)
			}
			if _, err := Load(); err == nil {
				t.Fatalf("expected error for %s", name)
			}
		})
	}
}
