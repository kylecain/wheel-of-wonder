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
var components = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}

func RegisterAll(s *discordgo.Session, config *config.Config, repository *repository.MovieRepository) {
	commands["addmovie"] = NewAddMovieCommand(repository)
	commands["allmovies"] = NewAllMoviesCommand(repository)
	commands["spin"] = NewSpinCommand(repository)
	commands["activemovie"] = NewActiveMovieCommand(repository)
	commands["setwatched"] = NewSetWatchedCommand(repository)

	for _, cmd := range commands {
		s.ApplicationCommandCreate(s.State.User.ID, config.GuildId, cmd.Definition())

		if ch, ok := cmd.(ComponentHandler); ok {
			for customID, handler := range ch.ComponentHandlers() {
				components[customID] = handler
			}
		}
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			cmd := i.ApplicationCommandData().Name
			if handler, ok := commands[cmd]; ok {
				handler.HandleCommand(s, i)
			}
		case discordgo.InteractionMessageComponent:
			customID := i.MessageComponentData().CustomID
			if handler, ok := components[customID]; ok {
				handler(s, i)
			}
		}
	})
}
