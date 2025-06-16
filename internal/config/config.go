package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken     string
	GuildId      string
	MigrationUrl string
	DatabaseUrl  string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		slog.Error("error loading .env file", "error", err)
		os.Exit(1)
	}

	botToken := os.Getenv("BOT_TOKEN")
	guildId := os.Getenv("GUILD_ID")
	migationUrl := os.Getenv("MIGRATION_URL")
	databaseUrl := os.Getenv("DATABASE_URL")

	return &Config{
		BotToken:     botToken,
		GuildId:      guildId,
		MigrationUrl: migationUrl,
		DatabaseUrl:  databaseUrl,
	}
}
