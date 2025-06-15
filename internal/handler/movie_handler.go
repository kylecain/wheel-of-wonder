package handler

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

type MovieHandler struct{}

func NewMovieHandler() *MovieHandler {
	return &MovieHandler{}
}

func (h *MovieHandler) Add(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}
	if i.ApplicationCommandData().Name == "add" {
		input := i.ApplicationCommandData().Options[0].StringValue()
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
