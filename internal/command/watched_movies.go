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

	var b strings.Builder
	for _, movie := range movies {
		dateStr := movie.UpdatedAt.Format("2006-01-02")
		b.WriteString(fmt.Sprintf("**%s** — %s\n", movie.Title, dateStr))
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Watched Movies",
		Description: b.String(),
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%d watched", len(movies)),
		},
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		slog.Error("failed to respond to watched movies command", "error", err)
	}
}
