package config

import (
	"os"
	"testing"
)

func TestLoad_BotTokenRequired(t *testing.T) {
	os.Setenv("TELEGRAM_BOT_TOKEN", "test-token")
	os.Setenv("ALLOWED_USER_IDS", "12345,67890")
	defer func() {
		os.Unsetenv("TELEGRAM_BOT_TOKEN")
		os.Unsetenv("ALLOWED_USER_IDS")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.BotToken != "test-token" {
		t.Errorf("expected bot token 'test-token', got %q", cfg.BotToken)
	}
	if len(cfg.AllowedUserIDs) != 2 {
		t.Fatalf("expected 2 allowed user IDs, got %d", len(cfg.AllowedUserIDs))
	}
	if !cfg.IsAllowed(12345) {
		t.Error("expected 12345 to be allowed")
	}
	if cfg.IsAllowed(99999) {
		t.Error("expected 99999 to NOT be allowed")
	}
}

func TestLoad_MissingToken(t *testing.T) {
	os.Unsetenv("TELEGRAM_BOT_TOKEN")
	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestLoad_MissingUserIDs(t *testing.T) {
	os.Setenv("TELEGRAM_BOT_TOKEN", "test-token")
	os.Unsetenv("ALLOWED_USER_IDS")
	defer os.Unsetenv("TELEGRAM_BOT_TOKEN")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing user IDs")
	}
}

func TestIsAllowed(t *testing.T) {
	os.Setenv("TELEGRAM_BOT_TOKEN", "test")
	os.Setenv("ALLOWED_USER_IDS", "42")
	defer func() {
		os.Unsetenv("TELEGRAM_BOT_TOKEN")
		os.Unsetenv("ALLOWED_USER_IDS")
	}()

	cfg, _ := Load()
	if !cfg.IsAllowed(42) {
		t.Error("expected 42 to be allowed")
	}
	if cfg.IsAllowed(1) {
		t.Error("expected 1 to NOT be allowed")
	}
}
