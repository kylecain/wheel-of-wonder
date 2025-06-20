package command

import (
	"fmt"
	"log/slog"
	"math/rand"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
)

type SpinCommand struct {
	MovieRepository *repository.MovieRepository
}

func NewSpinCommand(movieRepository *repository.MovieRepository) *SpinCommand {
	return &SpinCommand{
		MovieRepository: movieRepository,
	}
}

func (c *SpinCommand) Definition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        commandNameSpin,
		Description: "Spin the wheel and get a random movie",
	}
}

func (c *SpinCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	movies, err := c.MovieRepository.GetAll(i.GuildID)
	if err != nil || len(movies) == 0 {
		InteractionResponseError(s, i, err, "failed to get all movies for spin")
		return
	}

	selectedMovie := movies[rand.Intn(len(movies))]
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("You spun the wheel and got: %s", selectedMovie.Title),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Style:    discordgo.PrimaryButton,
							Label:    "Set as Active",
							CustomID: fmt.Sprintf("%s:%d", customIdSetActiveMovie, selectedMovie.ID),
						},
					},
				},
			},
		},
	})

	if err != nil {
		slog.Error("failed to respond to spin command", "error", err)
	}
}

func (c *AddMovieCommand) ComponentHandlers() map[string]func(*discordgo.Session, *discordgo.InteractionCreate) {
	return map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		customIdSetActiveMovie: c.setActiveMovieHandler,
	}
}

func (c *AddMovieCommand) setActiveMovieHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	value := strings.Split(i.MessageComponentData().CustomID, ":")[1]

	movieID, err := strconv.Atoi(value)
	if err != nil {
		InteractionResponseError(s, i, err, "Failed to convert movieID")
		return
	}

	currentlyActiveMovie, err := c.MovieRepository.GetActive(i.GuildID)
	if err != nil || currentlyActiveMovie == nil {
		InteractionResponseError(s, i, err, "failed to get currently active movie")
		return
	}

	err = c.MovieRepository.UpdateActive(currentlyActiveMovie.ID, false)
	if err != nil {
		InteractionResponseError(s, i, err, "failed to update currently active movie")
		return
	}

	err = c.MovieRepository.UpdateActive(movieID, true)
	if err != nil {
		InteractionResponseError(s, i, err, "failed to update currently active movie")
		return
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("You set the movie with ID %d as active.", movieID),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
