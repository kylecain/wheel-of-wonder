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
							CustomID: fmt.Sprintf("%s:%d:%s", customIdSetActiveMovie, selectedMovie.ID, selectedMovie.Title),
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
		customIdSetActiveMovie:   c.setActiveMovieHandler,
		customIdCreateEventModal: c.createEventModalHandler,
		// customIdCreateEvent:      c.createEventHandler,
	}
}

func (c *AddMovieCommand) setActiveMovieHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	args := strings.Split(i.MessageComponentData().CustomID, ":")
	movieIDStr, movieTitle := args[1], args[2]

	movieID, err := strconv.Atoi(movieIDStr)
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

	activeMovie, err := c.MovieRepository.GetMovieByID(movieID)
	if err != nil || activeMovie == nil {
		InteractionResponseError(s, i, err, "failed to retrieve active movie")
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("You set %s as the active movie.", activeMovie.Title),
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Style:    discordgo.PrimaryButton,
							Label:    "Create Event",
							CustomID: fmt.Sprintf("%s:%s:%s", customIdCreateEventModal, movieIDStr, movieTitle),
						},
					},
				},
			},
		},
	})
	if err != nil {
		slog.Error("failed to respond to set active movie command", "error", err)
		return
	}
}

func (c *AddMovieCommand) createEventModalHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	args := strings.Split(i.MessageComponentData().CustomID, ":")
	movieIDStr, movieTitle := args[1], args[2]

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: fmt.Sprintf("%s:%s:%s", customIdCreateEventModal, movieIDStr, movieTitle),
			Title:    "Create Event",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "date",
							Label:       "Enter the date of the event",
							Style:       discordgo.TextInputShort,
							Placeholder: "YYYY-MM-DD",
							Required:    true,
							MaxLength:   10,
							MinLength:   10,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "time",
							Label:       "Enter the time of the event (24-hour format)",
							Style:       discordgo.TextInputShort,
							Placeholder: "HH:mm:ss",
							Required:    true,
							MaxLength:   8,
							MinLength:   8,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "timezone",
							Label:       "Enter the timezone of the event",
							Style:       discordgo.TextInputShort,
							Placeholder: "America/Chicago",
							Required:    true,
							MaxLength:   50,
							MinLength:   8,
						},
					},
				},
			},
		},
	})
	if err != nil {
		InteractionResponseError(s, i, err, "failed to create event modal")
		return
	}
}

// func (c *SpinCommand) createEventHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
// 	data := i.ModalSubmitData()
// }
