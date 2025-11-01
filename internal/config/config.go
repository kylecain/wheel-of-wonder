package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken      string
	ApplicationID string
	GuildID       string
}

func NewConfig(logger *slog.Logger) (*Config, error) {
	configLogger := logger.With(slog.String("component", "config"))

	if err := godotenv.Load(); err != nil {
		configLogger.Debug(".env file not found", slog.Any("err", err))
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		return nil, fmt.Errorf("missing required env: BOT_TOKEN")
	}

	application_id := os.Getenv("APPLICATION_ID")
	if application_id == "" {
		configLogger.Warn("missing APPLICATION_ID; you will not be able to delete commands", slog.String("env", "APPLICATION_ID"))
	}

	guildId := os.Getenv("GUILD_ID")
	if guildId == "" {
		configLogger.Warn("missing GUILD_ID; commands will be installed/removed globally; this can take up to 2 hours", slog.String("env", "GUILD_ID"))
	}

	return &Config{
		BotToken:      botToken,
		ApplicationID: application_id,
		GuildID:       guildId,
	}, nil
}
