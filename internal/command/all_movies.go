package command

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
)

type AllMovies struct {
	MovieRepository *repository.MovieRepository
}

func NewAllMovies(movieRepository *repository.MovieRepository) *AllMovies {
	return &AllMovies{
		MovieRepository: movieRepository,
	}
}

func (c *AllMovies) Definition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        commandNameAllMovies,
		Description: "Get all movies in the wheel",
	}
}

func (c *AllMovies) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	movies, err := c.MovieRepository.GetAll(i.GuildID)
	if err != nil {
		InteractionResponseError(s, i, err, "Failed to get all movies.")
		return
	}

	response := "Movies in the wheel:\n"
	for _, movie := range movies {
		response += fmt.Sprintf("- ID: %d - Title: %s\n", movie.ID, movie.Title)
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
