package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken         string
	ApplicationID    string
	GuildID          string
	GeneralChannelID string
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		slog.Info(".env file not found", "error", err)
	}

	botToken := os.Getenv("BOT_TOKEN")
	guildId := os.Getenv("GUILD_ID")
	generalChannelID := os.Getenv("GENERAL_CHANNEL_ID")

	return &Config{
		BotToken:         botToken,
		GuildID:          guildId,
		GeneralChannelID: generalChannelID,
	}
}
