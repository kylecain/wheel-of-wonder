package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		slog.Error("error loading .env file", "error", err)
		os.Exit(1)
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		slog.Error("BOT_TOKEN is required")
		os.Exit(1)
	}

	return &Config{
		BotToken: botToken,
	}
}
