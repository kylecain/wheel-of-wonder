package command

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
	"github.com/kylecain/wheel-of-wonder/internal/util"
)

type WatchedMovies struct {
	MovieRepository *repository.Movie
}

func NewWatchedMovies(movieRepository *repository.Movie) *WatchedMovies {
	return &WatchedMovies{
		MovieRepository: movieRepository,
	}
}

func (c *WatchedMovies) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        commandNameWatchedMovies,
		Description: "Get previously watched movies",
	}
}

func (c *WatchedMovies) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	movies, err := c.MovieRepository.GetAllWatched(i.GuildID)
	if err != nil {
		util.InteractionResponseError(s, i, err, "Failed to get watched movies.")
		return
	}

	maxTitleLen := len("Movie")
	maxDateLen := len("Watched At")
	for _, movie := range movies {
		if len(movie.Title) > maxTitleLen {
			maxTitleLen = len(movie.Title)
		}
		dateStr := movie.UpdatedAt.Format("2006-01-02")
		if len(dateStr) > maxDateLen {
			maxDateLen = len(dateStr)
		}
	}

	response := "```" +
		fmt.Sprintf("| %-*s | %-*s |\n", maxTitleLen, "Movie", maxDateLen, "Watched At") +
		fmt.Sprintf("|%s|%s|\n", strings.Repeat("-", maxTitleLen+2), strings.Repeat("-", maxDateLen+2))
	for _, movie := range movies {
		response += fmt.Sprintf("| %-*s | %-*s |\n", maxTitleLen, movie.Title, maxDateLen, movie.UpdatedAt.Format("2006-01-02"))
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
		slog.Error("failed to respond to watched movies command", "error", err)
	}
}
