package command

import (
	"fmt"
	"log/slog"
	"math/rand"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/component"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
	"github.com/kylecain/wheel-of-wonder/internal/model"
	"github.com/kylecain/wheel-of-wonder/internal/util"
)

type Spin struct {
	movieRepository *repository.Movie
	userRepository  *repository.User
	logger          *slog.Logger
}

func NewSpin(movieRepository *repository.Movie, userRepository *repository.User, logger *slog.Logger) *Spin {
	return &Spin{
		movieRepository: movieRepository,
		userRepository:  userRepository,
		logger:          logger,
	}
}

func (c *Spin) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        commandNameSpin,
		Description: "Spin the wheel and get a random movie",
	}
}

func (c *Spin) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	l := c.logger.
		With(slog.String("command_name", i.ApplicationCommandData().Name)).
		With(util.InteractionGroup(i))

	l.Info("received command interaction")
	movies, err := c.movieRepository.GetAll(i.GuildID)
	if err != nil {
		l.Error("failed to get all movies", slog.Any("err", err))
		util.RespondError(s, i, "Something went wrong fetching all movies.")
		return
	} else if len(movies) == 0 {
		l.Warn("no movies found")
		util.RespondError(s, i, "No movies were found.")
		return
	}

	selectedMovie := movies[rand.Intn(len(movies))]
	l = l.With(util.MovieGroup(&selectedMovie))
	l.Info("random movie selected")

	err = c.setActive(&selectedMovie, i, l)
	if err != nil {
		l.Error("failed to set the active movie", slog.Any("err", err))
		util.RespondError(s, i, "Something went wrong setting the active movie.")
		return
	}
	l.Info("active movie set")

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     util.MovieEmbedSlice(&selectedMovie),
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: c.createComponents(&selectedMovie, i, l),
		},
	})
	if err != nil {
		l.Error("failed to respond to interaction", slog.Any("err", err))
	} else {
		l.Info("successfully responded to command")
	}
}

func (c *Spin) createComponents(movie *model.Movie, i *discordgo.InteractionCreate, l *slog.Logger) []discordgo.MessageComponent {
	movieIDStr := strconv.FormatInt(movie.ID, 10)
	components := []discordgo.MessageComponent{
		component.CreateEventButton(movieIDStr, movie.Title),
		component.AnnounceMovieButton(movieIDStr, movie.Title),
	}

	user, err := c.userRepository.UserByUserId(i.Member.User.ID)
	if err != nil {
		l.Warn("user does not exist in database", slog.String("user_id", i.Member.User.ID))
	}
	if user != nil {
		l = l.With(util.UserGroup(user))
		l.Info("retrieved user information")
	}

	if user != nil {
		components = append(
			[]discordgo.MessageComponent{
				component.CreateEventPreferredtimeButton(movieIDStr, movie.Title),
			},
			components...,
		)
	}
	actionsRow := discordgo.ActionsRow{
		Components: components,
	}
	return []discordgo.MessageComponent{actionsRow}
}

func (c *Spin) setActive(movie *model.Movie, i *discordgo.InteractionCreate, l *slog.Logger) error {
	currentlyActiveMovie, err := c.movieRepository.GetActive(i.GuildID)
	if err != nil {
		l.Error("failed to get currently active movie", slog.Any("err", err))
		return fmt.Errorf("failed to get active movie: %w", err)
	}

	if currentlyActiveMovie != nil {
		l.Info("retrieved currently active movie", slog.String("currently_active_movie", currentlyActiveMovie.Title))
		err = c.movieRepository.UpdateActive(currentlyActiveMovie.ID, false)
		if err != nil {
			return fmt.Errorf("failed to update activity of currently active movie: %w", err)
		} else {
			l.Info("currently active movie set to false", slog.String("currently_active_movie", currentlyActiveMovie.Title))
		}

		err = c.movieRepository.UpdateWatched(currentlyActiveMovie.ID, true)
		if err != nil {
			return fmt.Errorf("failed to update watched status of currently active movie: %w", err)
		} else {

			l.Info("currently active watched set to true", slog.String("currently_active_movie", currentlyActiveMovie.Title))
		}
	} else {
		l.Warn("there is no currently active movie")
	}

	err = c.movieRepository.UpdateActive(movie.ID, true)
	if err != nil {
		return fmt.Errorf("failed to update spun movie active status: %w", err)
	}

	return nil
}
