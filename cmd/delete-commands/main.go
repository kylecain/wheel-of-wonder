package main

import (
	"log/slog"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/config"
)

func main() {
	config := config.NewConfig()

	s, err := discordgo.New("Bot " + config.BotToken)
	if err != nil {
		slog.Error("error creating bot", "error", err)
		os.Exit(1)
	}

	commands, err := s.ApplicationCommands(config.ApplicationID, config.GuildId)
	if err != nil {
		slog.Error("error getting commands", "error", err)
	}

	for _, cmd := range commands {
		err := s.ApplicationCommandDelete(config.ApplicationID, config.GuildId, cmd.ID)
		if err != nil {
			slog.Error("error deleting command", "error", err)
		}
	}

}
