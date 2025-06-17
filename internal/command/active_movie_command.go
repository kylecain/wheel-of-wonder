package command

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
)

type ActiveMovieCommand struct {
	MovieRepository *repository.MovieRepository
}

func NewActiveMovieCommand(movieRepository *repository.MovieRepository) *ActiveMovieCommand {
	return &ActiveMovieCommand{
		MovieRepository: movieRepository,
	}
}

func (c *ActiveMovieCommand) Definition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "activemovie",
		Description: "Show the active movie",
	}
}

func (c *ActiveMovieCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	activeMovie, err := c.MovieRepository.GetActive(i.GuildID)

	if err != nil {
		slog.Error("failed to get active movie", "error", err)
		return
	}

	response := fmt.Sprintf("The active movie is: %s", activeMovie.Title)
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
	if err != nil {
		slog.Error("failed to respond to add command", "error", err)
	}
}
