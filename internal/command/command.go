package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/config"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
	"github.com/kylecain/wheel-of-wonder/internal/service"
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
	movieSearchService *service.MovieSearch,
) {
	commands[commandNameActiveMovie] = NewActiveMovie(movieRepository)
	commands[commandNameAddMovie] = NewAddMovie(movieRepository, movieSearchService)
	commands[commandNameAllMovies] = NewAllMovies(movieRepository)
	commands[commandNameDeleteMovie] = NewDeleteMovie(movieRepository)
	commands[commandNameSetPreferredEventTime] = NewSetPreferredEventTime(userRepository)
	commands[commandNameSpin] = NewSpin(movieRepository, userRepository)
	commands[commandNameWatchedMovies] = NewWatchedMovies(movieRepository)

	for _, cmd := range commands {
		s.ApplicationCommandCreate(s.State.User.ID, config.GuildID, cmd.ApplicationCommand())
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
