package command

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
)

type AllMoviesCommand struct {
	MovieRepository *repository.MovieRepository
}

func NewAllMoviesCommand(movieRepository *repository.MovieRepository) *AllMoviesCommand {
	return &AllMoviesCommand{
		MovieRepository: movieRepository,
	}
}

func (c *AllMoviesCommand) Definition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "allmovies",
		Description: "Get all movies in the wheel",
	}
}

func (c *AllMoviesCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	movies, err := c.MovieRepository.GetAll(i.GuildID)
	if err != nil {
		slog.Error("failed to get all movies", "error", err)
		return
	}

	response := "Movies in the wheel:\n"
	for _, movie := range movies {
		response += fmt.Sprintf("- %s\n", movie.Title)
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
	if err != nil {
		slog.Error("failed to respond to getall command", "error", err)
	}
}
