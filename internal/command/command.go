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

func RegisterAll(
	s *discordgo.Session,
	config *config.Config,
	movieRepository *repository.Movie,
	userRepository *repository.User,
) {
	commands[commandNameAddMovie] = NewAddMovie(movieRepository)
	commands[commandNameAllMovies] = NewAllMovies(movieRepository)
	commands[commandNameSpin] = NewSpin(movieRepository)
	commands[commandNameActiveMovie] = NewActiveMovie(movieRepository)
	commands[commandNameSetActive] = NewSetActive(movieRepository)
	commands[commandNameSetWatched] = NewSetWatched(movieRepository)
	commands[commandNameSetPreferredEventTime] = NewSetPreferredEventTime(userRepository)

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
