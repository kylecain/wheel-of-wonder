package command

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
	"github.com/kylecain/wheel-of-wonder/internal/model"
	"github.com/kylecain/wheel-of-wonder/internal/service"
	"github.com/kylecain/wheel-of-wonder/internal/util"
)

type AddMovie struct {
	movieRepository *repository.Movie
	movieService    *service.Movie
	logger          *slog.Logger
}

func NewAddMovie(movieRepository *repository.Movie, movieService *service.Movie, logger *slog.Logger) *AddMovie {
	return &AddMovie{
		movieRepository: movieRepository,
		movieService:    movieService,
		logger:          logger,
	}
}

func (c *AddMovie) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        commandNameAddMovie,
		Description: "Add a movie",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "movie",
				Description: "movie that will be added to the wheel",
				Required:    true,
			},
		},
	}
}

func (c *AddMovie) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	input := i.ApplicationCommandData().Options[0].StringValue()
	l := c.logger.
		With(slog.String("command_name", i.ApplicationCommandData().Name)).
		With("input", input).
		With(util.InteractionGroup(i))

	l.Info("received command interaction")

	movieInfo, err := c.movieService.FetchMovie(input)
	if err != nil {
		l.Error("failed to fetch movie data", slog.Any("err", err))
		util.RespondError(s, i, "Something went wrong fetching movie data.")
		return
	}

	l.Info("fetched movie info", util.MovieInfoGroup(movieInfo))

	movie := &model.Movie{
		GuildID:     i.GuildID,
		UserID:      i.Member.User.ID,
		Username:    i.Member.User.Username,
		Title:       movieInfo.Title,
		Description: movieInfo.Description,
		Duration:    movieInfo.Duration,
		ImageURL:    movieInfo.ImageURL,
		ContentURL:  movieInfo.ContentURL,
	}

	_, err = c.movieRepository.AddMovie(movie)
	if err != nil {
		l.Error("failed to add movie to database", slog.Any("err", err))
		util.RespondError(s, i, "Something went wrong adding movie.")
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: util.MovieEmbedSlice(movie),
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		l.Error("failed to respond to interaction", slog.Any("err", err))
	} else {
		l.Info("successfully responded to command", slog.String("movie_title", movie.Title))
	}
}
