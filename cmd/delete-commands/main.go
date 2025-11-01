package main

import (
	"log/slog"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/config"
)

func main() {
	handlerOptions := slog.HandlerOptions{Level: slog.LevelInfo}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &handlerOptions))

	configLogger := logger.With(slog.String("component", "config"))
	config, err := config.NewConfig(configLogger)
	if err != nil {
		configLogger.Error("failed to load config", slog.Any("err", err))
		os.Exit(1)
	}

	s, err := discordgo.New("Bot " + config.BotToken)
	if err != nil {
		slog.Error("error creating bot", "error", err)
		os.Exit(1)
	}

	commands, err := s.ApplicationCommands(config.ApplicationID, config.GuildID)
	if err != nil {
		slog.Error("error getting commands", "error", err)
	}

	for _, cmd := range commands {
		err := s.ApplicationCommandDelete(config.ApplicationID, config.GuildID, cmd.ID)
		if err != nil {
			slog.Error("error deleting command", "error", err)
		}
	}

}
