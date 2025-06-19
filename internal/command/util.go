package command

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func InteractionResponseError(s *discordgo.Session, i *discordgo.InteractionCreate, err error, message string) {
	slog.Error(message, "error", err)

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})

	if err != nil {
		slog.Error("failed to respond to interaction", "error", err)
	}
}
