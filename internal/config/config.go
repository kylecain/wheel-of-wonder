package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken      string
	ApplicationID string
	GuildID       string
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		slog.Info(".env file not found")
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		slog.Error("Missing botToken. Exiting...")
		os.Exit(1)
	}

	application_id := os.Getenv("APPLICATION_ID")
	if application_id == "" {
		slog.Warn("Missing applicationID. You will not be able to delete commands.")
	}

	guildId := os.Getenv("GUILD_ID")
	if guildId == "" {
		slog.Warn("Missing guildID. Commands will be installed globally. This can take up to 2 hours.")
	}

	return &Config{
		BotToken:      botToken,
		ApplicationID: application_id,
		GuildID:       guildId,
	}
}
