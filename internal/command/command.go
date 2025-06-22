package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/config"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
)

type Command interface {
	Definition() *discordgo.ApplicationCommand
	HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate)
}

type ComponentHandler interface {
	ComponentHandlers() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var commands = map[string]Command{}

func RegisterAll(s *discordgo.Session, config *config.Config, repository *repository.MovieRepository) {
	commands[commandNameAddMovie] = NewAddMovieCommand(repository)
	commands[commandNameAllMovies] = NewAllMoviesCommand(repository)
	commands[commandNameSpin] = NewSpinCommand(repository)
	commands[commandNameActiveMovie] = NewActiveMovieCommand(repository)
	commands[commandNameSetActive] = NewSetActiveCommand(repository)
	commands[commandNameSetWatched] = NewSetWatchedCommand(repository)

	for _, cmd := range commands {
		s.ApplicationCommandCreate(s.State.User.ID, config.GuildId, cmd.Definition())
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			cmd := i.ApplicationCommandData().Name
			if handler, ok := commands[cmd]; ok {
				handler.HandleCommand(s, i)
			}
		}
	})
}
