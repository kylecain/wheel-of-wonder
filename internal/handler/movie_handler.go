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
	if i.ApplicationCommandData().Name == "addmovie" {
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

func (h *MovieHandler) GetAll(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}
	if i.ApplicationCommandData().Name == "allmovies" {
		movies, err := h.MovieRepository.GetAll()
		if err != nil {
			slog.Error("failed to get all movies", "error", err)
			return
		}

		response := "Movies in the wheel:\n"
		for _, movie := range movies {
			response += fmt.Sprintf("- %s\n", movie)
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: response,
			},
		})
		if err != nil {
			slog.Error("failed to respond to getall command", "error", err)
		}
	}
}
