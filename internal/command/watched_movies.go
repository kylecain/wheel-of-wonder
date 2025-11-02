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
	movieRepository *repository.Movie
	logger          *slog.Logger
}

func NewWatchedMovies(movieRepository *repository.Movie, logger *slog.Logger) *WatchedMovies {
	return &WatchedMovies{
		movieRepository: movieRepository,
		logger:          logger,
	}
}

func (c *WatchedMovies) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        commandNameWatchedMovies,
		Description: "Get previously watched movies",
	}
}

func (c *WatchedMovies) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	l := c.logger.
		With(slog.String("command_name", i.ApplicationCommandData().Name)).
		With(util.InteractionGroup(i))
	l.Info("received command interaction")

	movies, err := c.movieRepository.GetAllWatched(i.GuildID)
	if err != nil {
		l.Error("failed to get all watched movies", slog.Any("err", err))
		util.RespondError(s, i, "Something went wrong getting watched movie data.")
		return
	}

	l.Info("got all watched movies from the database", slog.Int("watched_movie_count", len(movies)))

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
		l.Error("failed to respond to interaction", slog.Any("err", err))
	} else {
		l.Info("successfully responded to command", slog.Int("movie_count", len(movies)))
	}
}
