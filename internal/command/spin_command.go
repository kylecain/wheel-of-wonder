package command

import (
	"fmt"
	"log/slog"
	"math/rand"

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
		Name:        "spin",
		Description: "spin the wheel and get a random movie",
	}
}

func (h *SpinCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	movies, err := h.MovieRepository.GetAll()
	if err != nil {
		slog.Error("failed to get all movies for spin", "error", err)
		return
	}

	if len(movies) == 0 {
		response := "No movies available to spin."
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: response,
			},
		})
		if err != nil {
			slog.Error("failed to respond to spin command", "error", err)
		}
		return
	}

	selectedMovie := movies[rand.Intn(len(movies))]
	response := fmt.Sprintf("You spun the wheel and got: %s", selectedMovie)

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
	if err != nil {
		slog.Error("failed to respond to spin command", "error", err)
	}
}
