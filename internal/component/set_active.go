package component

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
	"github.com/kylecain/wheel-of-wonder/internal/model"
)

type SetActive struct {
	MovieRepository *repository.Movie
	UserRepository  *repository.User
}

func NewSetActive(movieRepository *repository.Movie, userRepository *repository.User) *SetActive {
	return &SetActive{
		MovieRepository: movieRepository,
		UserRepository:  userRepository,
	}
}

func SetActiveButton(movie model.Movie) discordgo.Button {
	return discordgo.Button{
		Style:    discordgo.PrimaryButton,
		Label:    "Set as Active",
		CustomID: fmt.Sprintf("%s:%d:%s", CustomIdSetActiveMovie, movie.ID, movie.Title),
	}
}

func (c *SetActive) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	args := strings.Split(i.MessageComponentData().CustomID, ":")
	movieIDStr, movieTitle := args[1], args[2]

	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		InteractionResponseError(s, i, err, "Failed to convert movieID")
		return
	}

	currentlyActiveMovie, err := c.MovieRepository.GetActive(i.GuildID)
	if err != nil {
		InteractionResponseError(s, i, err, "failed to get currently active movie")
		return
	}

	if currentlyActiveMovie != nil {
		err = c.MovieRepository.UpdateActive(currentlyActiveMovie.ID, false)
		if err != nil {
			InteractionResponseError(s, i, err, "failed to update currently active movie")
			return
		}

		err = c.MovieRepository.UpdateWatched(currentlyActiveMovie.ID, true)
		if err != nil {
			InteractionResponseError(s, i, err, "failed to update previous movie to watched")
			return
		}
	}

	err = c.MovieRepository.UpdateActive(int64(movieID), true)
	if err != nil {
		InteractionResponseError(s, i, err, "failed to update currently active movie")
		return
	}

	activeMovie, err := c.MovieRepository.GetMovieByID(movieID)
	if err != nil || activeMovie == nil {
		InteractionResponseError(s, i, err, "failed to retrieve active movie")
		return
	}

	components := []discordgo.MessageComponent{
		CreateEventButton(movieIDStr, movieTitle),
	}

	user, _ := c.UserRepository.UserByUserId(i.Member.User.ID)
	if user != nil {
		components = append(
			[]discordgo.MessageComponent{
				CreateEventPreferredtimeButton(movieIDStr, movieTitle),
			},
			components...,
		)
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("You set %s as the active movie.", activeMovie.Title),
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: components,
				},
			},
		},
	})
	if err != nil {
		slog.Error("failed to respond to set active movie command", "error", err)
		return
	}
}
