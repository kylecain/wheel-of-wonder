package command

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
	"github.com/kylecain/wheel-of-wonder/internal/util"
)

type ActiveMovie struct {
	movieRepository *repository.Movie
	logger          *slog.Logger
}

func NewActiveMovie(movieRepository *repository.Movie, logger *slog.Logger) *ActiveMovie {
	return &ActiveMovie{
		movieRepository: movieRepository,
		logger:          logger,
	}
}

func (c *ActiveMovie) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        commandNameActiveMovie,
		Description: "Show the active movie",
	}
}

func (c *ActiveMovie) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	l := util.WithInteractionLogger(c.logger, i)
	l.Info("received command interaction")

	activeMovie, err := c.movieRepository.GetActive(i.GuildID)
	if err != nil {
		l.Error("failed to get active movie", slog.Any("err", err))
		util.RespondError(s, i, "Something went wrong fetching the active movie.")
		return
	} else if activeMovie == nil {
		l.Warn("no active movie found")
		util.RespondError(s, i, "No active movie found.")
		return
	}

	l.Info("responding with active movie", slog.String("movie_title", activeMovie.Title))

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: util.MovieEmbedSlice(activeMovie),
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		l.Error("failed to respond to interaction", slog.Any("err", err))
	}
}
