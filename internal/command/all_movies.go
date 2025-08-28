package command

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
)

type AllMovies struct {
	MovieRepository *repository.Movie
}

func NewAllMovies(movieRepository *repository.Movie) *AllMovies {
	return &AllMovies{
		MovieRepository: movieRepository,
	}
}

func (c *AllMovies) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        commandNameAllMovies,
		Description: "Get all movies in the wheel",
	}
}

func (c *AllMovies) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	movies, err := c.MovieRepository.GetAll(i.GuildID)
	if err != nil {
		InteractionResponseError(s, i, err, "Failed to get all movies.")
		return
	}

	maxIDLen := len("ID")
	maxTitleLen := len("Movie")
	for _, movie := range movies {
		idLen := len(fmt.Sprintf("%d", movie.ID))
		if idLen > maxIDLen {
			maxIDLen = idLen
		}
		if len(movie.Title) > maxTitleLen {
			maxTitleLen = len(movie.Title)
		}
	}

	response := "```" +
		fmt.Sprintf("| %-*s | %-*s |\n", maxIDLen, "ID", maxTitleLen, "Movie") +
		fmt.Sprintf("|%s|%s|\n", strings.Repeat("-", maxIDLen+2), strings.Repeat("-", maxTitleLen+2))
	for _, movie := range movies {
		response += fmt.Sprintf("| %-*d | %-*s |\n", maxIDLen, movie.ID, maxTitleLen, movie.Title)
	}
	response += "```"

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		slog.Error("failed to respond to getall command", "error", err)
	}
}
