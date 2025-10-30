package command

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/component"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
)

type SetPreferredEventTime struct {
	UserRepository *repository.User
}

func NewSetPreferredEventTime(userRepository *repository.User) *SetPreferredEventTime {
	return &SetPreferredEventTime{
		UserRepository: userRepository,
	}
}

func (c *SetPreferredEventTime) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        commandNameSetPreferredEventTime,
		Description: "Set your preferred event time",
	}
}

func (c *SetPreferredEventTime) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID:   component.CustomIDSetPreferredTimeModal,
			Title:      "Set Preferred Time Modal",
			Components: component.SetPreferredEventTimeModal(),
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		slog.Error("Failed to respond to set-preferred-event-time command", "error", err)
	}
}
