package command

import (
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
	MovieRepository *repository.Movie
	UserRepository  *repository.User
}

func NewSpin(movieRepository *repository.Movie, userRepository *repository.User) *Spin {
	return &Spin{
		MovieRepository: movieRepository,
		UserRepository:  userRepository,
	}
}

func (c *Spin) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        commandNameSpin,
		Description: "Spin the wheel and get a random movie",
	}
}

func (c *Spin) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	movies, err := c.MovieRepository.GetAll(i.GuildID)
	if err != nil || len(movies) == 0 {
		util.InteractionResponseError(s, i, err, "failed to get all movies for spin")
		return
	}

	selectedMovie := movies[rand.Intn(len(movies))]

	err = c.setActive(&selectedMovie, i)
	if err != nil {
		util.InteractionResponseError(s, i, err, "failed to set active movie during spin")
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     util.MovieEmbedSlice(&selectedMovie),
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: c.createComponents(&selectedMovie, i),
		},
	})

	if err != nil {
		slog.Error("Failed to respond to spin command", "error", err)
	}
}

func (c *Spin) createComponents(movie *model.Movie, i *discordgo.InteractionCreate) []discordgo.MessageComponent {
	movieIDStr := strconv.FormatInt(movie.ID, 10)
	components := []discordgo.MessageComponent{
		component.CreateEventButton(movieIDStr, movie.Title),
		component.AnnounceMovieButton(movieIDStr, movie.Title),
	}

	user, _ := c.UserRepository.UserByUserId(i.Member.User.ID)
	if user != nil {
		components = append(
			[]discordgo.MessageComponent{
				component.CreateEventPreferredtimeButton(movieIDStr, movie.Title),
			},
			components...,
		)
	}
	return components
}

func (c *Spin) setActive(movie *model.Movie, i *discordgo.InteractionCreate) error {
	currentlyActiveMovie, err := c.MovieRepository.GetActive(i.GuildID)
	if err != nil {
		return err
	}

	if currentlyActiveMovie != nil {
		err = c.MovieRepository.UpdateActive(currentlyActiveMovie.ID, false)
		if err != nil {
			return err
		}

		err = c.MovieRepository.UpdateWatched(currentlyActiveMovie.ID, true)
		if err != nil {
			return err
		}
	}

	err = c.MovieRepository.UpdateActive(movie.ID, true)
	if err != nil {
		return err
	}

	return nil
}
