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
	MigrationUrl     string
	DatabaseUrl      string
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		slog.Info(".env file not found", "error", err)
	}

	botToken := os.Getenv("BOT_TOKEN")
	guildId := os.Getenv("GUILD_ID")
	generalChannelID := os.Getenv("GENERAL_CHANNEL_ID")
	migationUrl := os.Getenv("MIGRATION_URL")
	databaseUrl := os.Getenv("DATABASE_URL")

	return &Config{
		BotToken:         botToken,
		GuildID:          guildId,
		GeneralChannelID: generalChannelID,
		MigrationUrl:     migationUrl,
		DatabaseUrl:      databaseUrl,
	}
}
