package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	BotToken       string
	AllowedUserIDs map[int64]bool
}

func Load() (*Config, error) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}

	userIDsStr := os.Getenv("ALLOWED_USER_IDS")
	if userIDsStr == "" {
		return nil, fmt.Errorf("ALLOWED_USER_IDS is required (comma-separated)")
	}

	allowedUserIDs := make(map[int64]bool)
	for _, idStr := range strings.Split(userIDsStr, ",") {
		idStr = strings.TrimSpace(idStr)
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid user ID %q: %w", idStr, err)
		}
		allowedUserIDs[id] = true
	}

	return &Config{
		BotToken:       token,
		AllowedUserIDs: allowedUserIDs,
	}, nil
}

func (c *Config) IsAllowed(userID int64) bool {
	return c.AllowedUserIDs[userID]
}
