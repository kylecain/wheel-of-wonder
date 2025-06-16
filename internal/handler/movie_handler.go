package handler

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
)

type MovieHandler struct {
	MovieRepository *repository.MovieRepository
}

func NewMovieHandler(movieRepository *repository.MovieRepository) *MovieHandler {
	return &MovieHandler{
		MovieRepository: movieRepository,
	}
}

func (h *MovieHandler) Add(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}
	if i.ApplicationCommandData().Name == "add" {
		input := i.ApplicationCommandData().Options[0].StringValue()

		h.MovieRepository.Create(input)

		response := fmt.Sprintf("You added: %s", input)
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: response,
			},
		})
		if err != nil {
			slog.Error("failed to respond to add command", "error", err)
		}
	}
}
