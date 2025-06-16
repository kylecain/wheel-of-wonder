package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/config"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
)

type Command interface {
	Definition() *discordgo.ApplicationCommand
	Handle(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var registeredCommands = map[string]Command{}

func RegisterAll(s *discordgo.Session, config *config.Config, repository *repository.MovieRepository) {
	registeredCommands["addmovie"] = NewAddMovieCommand(repository)
	registeredCommands["allmovies"] = NewAllMoviesCommand(repository)
	registeredCommands["spin"] = NewSpinCommand(repository)

	for _, cmd := range registeredCommands {
		s.ApplicationCommandCreate(s.State.User.ID, config.GuildId, cmd.Definition())
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		cmd := i.ApplicationCommandData().Name
		if handler, ok := registeredCommands[cmd]; ok {
			handler.Handle(s, i)
		}
	})
}
