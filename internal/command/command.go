package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/config"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
)

type Command interface {
	ApplicationCommand() *discordgo.ApplicationCommand
	Handler(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var commands = map[string]Command{}

func RegisterAll(s *discordgo.Session, config *config.Config, repository *repository.Movie) {
	commands[commandNameAddMovie] = NewAddMovie(repository)
	commands[commandNameAllMovies] = NewAllMovies(repository)
	commands[commandNameSpin] = NewSpin(repository)
	commands[commandNameActiveMovie] = NewActiveMovie(repository)
	commands[commandNameSetActive] = NewSetActive(repository)
	commands[commandNameSetWatched] = NewSetWatched(repository)

	for _, cmd := range commands {
		s.ApplicationCommandCreate(s.State.User.ID, config.GuildId, cmd.ApplicationCommand())
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			cmd := i.ApplicationCommandData().Name
			if handler, ok := commands[cmd]; ok {
				handler.Handler(s, i)
			}
		}
	})
}
