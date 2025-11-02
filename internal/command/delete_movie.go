package command

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/component"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
	"github.com/kylecain/wheel-of-wonder/internal/util"
)

type DeleteMovie struct {
	movieRepository *repository.Movie
	logger          *slog.Logger
}

func NewDeleteMovie(movieRepository *repository.Movie, logger *slog.Logger) *DeleteMovie {
	return &DeleteMovie{
		movieRepository: movieRepository,
		logger:          logger,
	}
}

func (c *DeleteMovie) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        commandNameDeleteMovie,
		Description: "Select movie from list to delete",
	}
}

func (c *DeleteMovie) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	l := c.logger.
		With(slog.String("command_name", i.ApplicationCommandData().Name)).
		With(util.InteractionGroup(i))

	l.Info("recieved command interaction")
	movies, err := c.movieRepository.GetAllUnwatched(i.GuildID)
	if err != nil {
		l.Error("failed to get unwatched movies", slog.Any("err", err))
		util.RespondError(s, i, "Something went wrong fetching unwatched movies for guild.")
		return
	} else if len(movies) == 0 {
		l.Warn("no unwatched movies found")
		util.RespondError(s, i, "No unwatched movies were found for guild.")
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    "Pick one of the options below:",
			Components: component.MovieSelectMenu(movies, component.CustomIdDeleteMovie, i),
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		l.Error("failed to respond to interaction", slog.Any("err", err))
	} else {
		l.Info("successfully responded to command")
	}
}
